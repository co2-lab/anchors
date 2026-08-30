// Superfícies CONSUMIDORAS do testID: onde procurar quem se apoia no handle — o teste
// ligado pela trinca, o teste vizinho (compartilhado ou do pai) e os flows de ponta a
// ponta, que vivem fora do grafo.
//
// Este arquivo já foi `testid_honored.go`. O gate foi aposentado por `testid-coerente`
// — que também cobre a ponta que o antigo declarava fora de escopo (id consultado e
// inexistente), depois de a premissa "isso já falha ao rodar, ruidosamente" ter sido
// refutada por medição.
package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// idConsultado — o id aparece em alguma superfície consumidora? Para o template
// (`bdgt-item-*`), basta o PREFIXO ser consultado: o sufixo é dado de runtime, e o
// flow que casa `bdgt-item-.*` ou constrói `bdgt-item-${id}` está usando o contrato.
func idConsultado(blob, id string) bool {
	nu := strings.TrimPrefix(id, ":")
	if strings.HasSuffix(nu, "-*") {
		return strings.Contains(blob, strings.TrimSuffix(nu, "*"))
	}
	// Confronta as duas formas (com e sem a marca), porque a superfície consumidora
	// pode ter sido escrita antes da convenção de marcação.
	return strings.Contains(blob, nu) || strings.Contains(blob, ":"+nu)
}

// lerSuperficieE2E lê os flows de ponta a ponta declarados pelo projeto.
//
// A resolução é em DOIS passos, e confundi-los custa caro: `surfaces[e2e]` devolve a
// CHAVE da superfície (ex.: "e2e"), não um caminho — quem tem o caminho é
// `files[chave]`. Ler `surfaces` como se fosse path faz o gate procurar um diretório
// que não existe, achar zero consumidores e reportar como ÓRFÃO todo id usado apenas
// pelo E2E — acusando exatamente quem cumpre o contrato.
//
// Quando o projeto declara o regime mas não o arquivo (o app de referência hoje: `e2e: e2e` sem
// `files.e2e`), não há onde procurar. Devolver vazio aqui é honesto; o gate trata a
// falta de consumidor conhecido como Skip, não como reprovação.
func lerSuperficieE2E(root string, cfg *config.Config) []string {
	if cfg == nil || cfg.Derived == nil {
		return nil
	}
	chave := cfg.Derived.Surfaces["e2e"]
	if chave == "" {
		return nil
	}
	// O PRIMEIRO padrão da camada: este gate resolve UM caminho de superfície, e uma
	// camada com vários padrões não muda o que ele pergunta.
	var padrao string
	if ps := cfg.Derived.PadroesDe()[chave]; len(ps) > 0 {
		padrao = ps[0]
	}
	if padrao == "" {
		// Sem template de arquivo para a superfície, procuramos por override — a
		// mesma precedência que o resto do framework usa.
		for _, ov := range cfg.Derived.Overrides {
			if ps := ov.PadroesDe()[chave]; len(ps) > 0 {
				padrao = ps[0]
				break
			}
		}
	}
	if padrao == "" {
		return nil
	}
	// A superfície é declarada como template de região (ex.: `apps/mobile/.maestro`);
	// aqui interessa a RAIZ dela, varrida por inteiro.
	dir := filepath.Join(root, primeiroSegmentoEstatico(padrao))
	var out []string
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if b, e := os.ReadFile(p); e == nil {
			out = append(out, string(b))
		}
		return nil
	})
	return out
}

// lerTestesVizinhos lê os testes da MESMA pasta e os da pasta-IRMÃ dentro do módulo.
//
// A pasta sozinha não basta, e o caso que mostra isso é o componente: `CategoryCard`
// vive em `trends/components/` e quem o exercita é `TrendsScreen.test.tsx`, em
// `trends/screens/` — a tela que o renderiza. Sem subir um nível, o gate acusa de
// órfão um handle que o teste consulta, e o achado aponta para o lugar errado.
//
// Sobe UM nível só (o módulo da feature), não a árvore inteira: ler tudo tornaria
// qualquer menção do repositório uma prova de consumo, e o gate deixaria de medir.
func lerTestesVizinhos(root, specID string) []string {
	dir := filepath.Join(root, filepath.Dir(specID))
	dirs := []string{dir}
	// As irmãs dentro do módulo (`features/trends/{components,screens,hooks}`).
	if entradas, err := os.ReadDir(filepath.Dir(dir)); err == nil {
		for _, e := range entradas {
			if e.IsDir() {
				if d := filepath.Join(filepath.Dir(dir), e.Name()); d != dir {
					dirs = append(dirs, d)
				}
			}
		}
	}
	var out []string
	for _, d := range dirs {
		entradas, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		for _, e := range entradas {
			if e.IsDir() || !strings.Contains(e.Name(), ".test.") {
				continue
			}
			if b, err := os.ReadFile(filepath.Join(d, e.Name())); err == nil {
				out = append(out, string(b))
			}
		}
	}
	return out
}

// primeiroSegmentoEstatico corta o template no primeiro placeholder (`{{module}}`,
// `{unit}`), devolvendo o prefixo fixo — a raiz que dá para varrer.
var placeholderRE = regexp.MustCompile(`\{\{?[a-zA-Z_]+\}?\}`)

func primeiroSegmentoEstatico(padrao string) string {
	if loc := placeholderRE.FindStringIndex(padrao); loc != nil {
		padrao = padrao[:loc[0]]
	}
	return strings.TrimSuffix(filepath.Clean(padrao), "/")
}
