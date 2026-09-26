package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// marker-parity: a MESMA regra tem de aparecer nas DUAS pontas que a cumprem.
//
// A classe de defeito é o de-para que se desfaz de um lado só. Uma regra que vive em
// dois lugares — a promessa na interface e o cumprimento no servidor — não tem como ser
// conferida por nenhum gate de arquivo: cada lado, olhado sozinho, está impecável. O que
// falha é a RELAÇÃO, e ela some sem deixar erro.
//
// O caso que motivou (MIF, `EXSC-Q01`): a página de exclusão de dados LISTA o que será
// apagado em cada escopo, e o backend APAGA. As duas listas nasceram juntas e nada as
// liga. Se um escopo passar a apagar mais (ou menos), a página segue exibindo a versão
// antiga — e o titular consente com base nela. É consentimento informado sobre exercício
// de direito, então o desencontro não é cosmético.
//
// O gate NÃO é canônico de propósito: ele é genérico, e um projeto pode declarar vários,
// cada um com o seu prefixo e o seu id. O que ele confronta é uma decisão do projeto —
// quais regras vivem em duas pontas — e essa decisão não cabe num catálogo universal.
//
// Declaração (no `anchors.yaml` do projeto):
//
//	gates:
//	  - name: paridade-exclusao-de-dados
//	    id: data-purge-parity           # ÚNICO: é por ele que se dispensa e se cita
//	    check: marker-parity
//	    scope: project
//	    blocking: true
//	    marker_prefix: data-purge-rule  # o prefixo das marcações no código
//	    marker_count: 2                 # quantas ocorrências CADA regra precisa ter
//	    marker_scopes:                  # e ONDE elas têm de estar (uma por escopo)
//	      - apps/landing-page/**
//	      - packages/backend/**
//
// E no código, dos dois lados:
//
//	// @data-purge-rule-conta-completa: a lista do que a página promete apagar
//	// @data-purge-rule-conta-completa: o que o handler de fato apaga
//
// O sufixo depois do prefixo é o NOME da regra, e é ele que emparelha as pontas. Duas
// marcações do mesmo lado NÃO satisfazem o gate: `marker_scopes` existe justamente
// porque contar sem olhar onde deixaria passar o caso que o gate existe para pegar.
func checkMarkerParity(g config.Gate, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	prefixo := strings.TrimSpace(g.MarkerPrefix)
	if prefixo == "" {
		return Pending, i18n.T("gate.marker_parity.pending_no_prefix")
	}

	escopos := g.MarkerScopes
	esperado := g.MarkerCount
	if esperado == 0 {
		esperado = len(escopos)
	}
	if esperado == 0 {
		return Pending, i18n.T("gate.marker_parity.pending_no_count_or_scopes")
	}

	ocorr := scanMarkings(root, prefixo, escopos, cfg)
	if len(ocorr) == 0 {
		// AUSÊNCIA TOTAL não é aprovação. Um prefixo que não aparece em lugar nenhum é
		// quase sempre erro de digitação na declaração — e devolver Pass ali faria o
		// gate parecer vigilante enquanto não vigia coisa alguma.
		return Pending, fmt.Sprintf(i18n.T("gate.marker_parity.pending_no_markers_found"), prefixo, prefixo)
	}

	var faltas []string
	nomes := make([]string, 0, len(ocorr))
	for nome := range ocorr {
		nomes = append(nomes, nome)
	}
	sort.Strings(nomes)

	for _, nome := range nomes {
		locais := ocorr[nome]
		if len(escopos) > 0 {
			// COM escopos declarados, a pergunta é "cada lado tem a sua?" — e é a
			// pergunta certa: duas marcações do mesmo lado somam 2 e esconderiam
			// exatamente o desencontro que se quer pegar.
			var vazios []string
			for _, esc := range escopos {
				if len(locais[esc]) == 0 {
					vazios = append(vazios, esc)
				}
			}
			if len(vazios) > 0 {
				faltas = append(faltas, fmt.Sprintf(i18n.T("gate.marker_parity.missing_in_scopes"),
					prefixo, nome, strings.Join(vazios, ", ")))
			}
			continue
		}
		// SEM escopos, resta a contagem — mais fraca, e por isso o gate a documenta
		// como o modo menos seguro.
		total := 0
		for _, arqs := range locais {
			total += len(arqs)
		}
		if total != esperado {
			faltas = append(faltas, fmt.Sprintf(i18n.T("gate.marker_parity.count_mismatch"),
				prefixo, nome, total, esperado))
		}
	}

	if len(faltas) > 0 {
		return Fail, fmt.Sprintf(i18n.T("gate.marker_parity.parity_failure"),
			len(faltas), prefixo, strings.Join(faltas, "\n      "))
	}

	return Pass, fmt.Sprintf(i18n.T("gate.marker_parity.pass_parity"),
		len(nomes), prefixo, esperado)
}

// scanMarkings devolve, por NOME de regra, os arquivos onde ela aparece, agrupados
// pelo escopo declarado que os contém.
// Never fails: an unreadable directory is skipped (MRPRM-E01), so the walk has no error
// to return — the "scan error" verdict it once fed could not happen.
func scanMarkings(root, prefixo string, escopos []string, cfg *config.Config) map[string]map[string][]string {
	// O nome da regra é o que vem depois do prefixo: letras, dígitos, `-` e `_`. Para
	// aqui de propósito — a marcação costuma ser seguida de `:` e da prosa que explica.
	re := regexp.MustCompile(`@` + regexp.QuoteMeta(prefixo) + `-([A-Za-z0-9_-]+)`)
	ocorr := map[string]map[string][]string{}

	_ = filepath.WalkDir(root, func(caminho string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // diretório ilegível não derruba a varredura inteira
		}
		if d.IsDir() {
			if ignored(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !likelyText(d.Name()) {
			return nil
		}
		rel, relErr := filepath.Rel(root, caminho)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)

		b, readErr := os.ReadFile(caminho)
		if readErr != nil {
			return nil
		}
		achados := re.FindAllStringSubmatch(string(b), -1)
		if len(achados) == 0 {
			return nil
		}
		esc := scopeOf(rel, escopos)
		for _, m := range achados {
			nome := m[1]
			if ocorr[nome] == nil {
				ocorr[nome] = map[string][]string{}
			}
			if !containsFile(ocorr[nome][esc], rel) {
				ocorr[nome][esc] = append(ocorr[nome][esc], rel)
			}
		}
		return nil
	})
	_ = cfg
	return ocorr
}

// scopeOf diz a QUAL escopo declarado o arquivo pertence. Sem escopos declarados (ou
// sem casar nenhum), cai num balde único — é o modo "só contagem".
func scopeOf(rel string, escopos []string) string {
	for _, e := range escopos {
		if ok, _ := doublestar.Match(e, rel); ok {
			return e
		}
	}
	return "(fora dos escopos)"
}

func containsFile(lista []string, alvo string) bool {
	for _, x := range lista {
		if x == alvo {
			return true
		}
	}
	return false
}

// ignored: diretórios que nunca contêm marcação de regra e cuja varredura só custa.
func ignored(nome string) bool {
	switch nome {
	case ".git", "node_modules", "dist", "build", ".next", "coverage", "vendor", ".anchors":
		return true
	}
	return false
}

// likelyText evita ler binário. A lista é por EXTENSÃO porque é onde a marcação vive:
// código e documentação. Um arquivo sem extensão conhecida é pulado — o custo de errar
// para menos aqui é o gate não ver uma marcação em lugar exótico, e o de errar para mais
// é ler megabytes de imagem a cada varredura.
func likelyText(nome string) bool {
	switch strings.ToLower(filepath.Ext(nome)) {
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".go", ".py", ".rb", ".java",
		".kt", ".swift", ".rs", ".php", ".cs", ".md", ".yaml", ".yml", ".sql", ".sh":
		return true
	}
	return false
}
