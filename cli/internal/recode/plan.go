package recode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/scan"
)

// FileChange é o que o recode faria num arquivo: o conteúdo novo e as ocorrências que
// motivaram a troca (para o dry-run mostrar).
type FileChange struct {
	Path        string
	Occurrences []Occurrence
	NewContent  string
	count       int
}

// FileRename é um arquivo cujo NOME contém o código: o recode faz git mv (dialeto).
type FileRename struct {
	From, To string
}

// Plan é o resultado do planejamento: os arquivos a mudar + as validações.
type Plan struct {
	Old, New string
	Files    []FileChange
	Total    int // total de substituições de CONTEÚDO

	// Dialeto de projeto (só quando o anchors.yaml declara `recode:`):
	TestIDs      int          // substituições de prefixo de testID (parte do Total? não — separado)
	Renames      []FileRename // arquivos a renomear (git mv)
	TestIDLegacy string       // aviso: prefixo esperado ausente, mas este candidato existe
}

// BuildPlan varre o projeto e monta o plano de recode de OLD→NEW. Valida:
//   - OLD e NEW bem-formados;
//   - NEW não está em uso por OUTRA unidade (colisão) — a menos que seja o próprio OLD;
//   - OLD realmente existe em algum arquivo (senão não há o que renomear).
//
// Não escreve nada — só lê. O comando decide dry-run vs apply.
func BuildPlan(root string, cfg *config.Config, old, new string) (*Plan, error) {
	if !ValidCode(old) {
		return nil, fmt.Errorf("código de origem %q inválido (esperado 4 chars A-Z0-9)", old)
	}
	if !ValidCode(new) {
		return nil, fmt.Errorf("código de destino %q inválido (esperado 4 chars A-Z0-9)", new)
	}
	if old == new {
		return nil, fmt.Errorf("origem e destino são o mesmo código (%s)", old)
	}

	files, err := scan.Walk(root, cfg)
	if err != nil {
		return nil, fmt.Errorf("varrer o projeto: %w", err)
	}

	// Colisão: NEW já é o code DONO de alguma unidade? (o scan carimba Codes por header)
	for _, f := range files {
		for _, c := range f.Codes {
			if c == new {
				return nil, fmt.Errorf("o código de destino %s já é usado por %s — escolha outro", new, f.Path)
			}
		}
	}

	// Dialeto de projeto (opcional): prefixo de testID + patterns de nome de arquivo.
	oldTID, newTID := "", ""
	var filePatterns []string
	if cfg.Recode != nil {
		oldTID = TestIDPrefix(old, cfg.Recode.TestID)
		newTID = TestIDPrefix(new, cfg.Recode.TestID)
		filePatterns = cfg.Recode.FilePatterns
	}

	plan := &Plan{Old: old, New: new}
	legacyHits := 0 // ocorrências de um prefixo de testID divergente (aviso de legado)
	for _, f := range files {
		content, rerr := os.ReadFile(filepath.Join(root, f.Path))
		if rerr != nil {
			continue
		}
		occ := Find(string(content), old)
		newContent, n := Rewrite(string(content), old, new)

		// dialeto: troca também o prefixo de testID derivado (se declarado).
		tn := 0
		if oldTID != "" {
			newContent, tn = RewriteTestIDs(newContent, oldTID, newTID)
			plan.TestIDs += tn
			// aviso de legado: só olhamos arquivos DA TRINCA (que têm header/scenario do
			// código) — se este arquivo tem testIDs mas o prefixo esperado não bateu, é
			// um prefixo divergente (recode manual anterior). Não conhecemos o prefixo
			// antigo; só que o esperado sumiu e há testIDs. Conta independente das
			// outras trocas no arquivo.
			if tn == 0 && n > 0 {
				legacyHits += CountAnyTestID(string(content))
			}
		}

		if n == 0 && tn == 0 {
			continue
		}
		plan.Files = append(plan.Files, FileChange{
			Path: f.Path, Occurrences: occ, NewContent: newContent, count: n + tn,
		})
		plan.Total += n
	}

	// renames de arquivo cujo NOME contém o código (dialeto).
	if len(filePatterns) > 0 {
		if renames, rerr := findRenames(root, old, new, filePatterns); rerr == nil {
			plan.Renames = renames
		}
	}

	// aviso de legado: o prefixo de testID esperado (oldTID) não apareceu em lugar
	// nenhum, mas há um candidato divergente (ex.: código-antigo minúsculo). Reporta.
	if oldTID != "" && plan.TestIDs == 0 && legacyHits > 0 {
		plan.TestIDLegacy = fmt.Sprintf("o prefixo de testID esperado %q não foi encontrado; "+
			"há %d testID(s) com um prefixo divergente (provável recode manual anterior). "+
			"O recode NÃO os toca (não adivinha) — padronize-os à mão para o prefixo derivado do código", oldTID, legacyHits)
	}

	if plan.Total == 0 && plan.TestIDs == 0 && len(plan.Renames) == 0 {
		return nil, fmt.Errorf("o código %s não aparece em nenhum arquivo do projeto", old)
	}
	sort.Slice(plan.Files, func(i, j int) bool { return plan.Files[i].Path < plan.Files[j].Path })
	return plan, nil
}

// findRenames varre o filesystem (não só as layers — .png/.yaml podem estar fora do
// grafo) por arquivos cujo NOME casa os patterns para o código old, e devolve o git mv.
func findRenames(root, old, new string, patterns []string) ([]FileRename, error) {
	var out []FileRename
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() && skipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, rerr := filepath.Rel(root, p)
		if rerr != nil {
			return nil
		}
		if FileMatchesCode(rel, old, patterns) {
			to := RenameFilePath(rel, old, new)
			if to != rel {
				out = append(out, FileRename{From: rel, To: to})
			}
		}
		return nil
	})
	sort.Slice(out, func(i, j int) bool { return out[i].From < out[j].From })
	return out, err
}

func skipDir(name string) bool {
	switch name {
	case "node_modules", ".git", "dist", "build", "vendor", ".next", "coverage", ".expo":
		return true
	}
	return false
}

// Apply grava as mudanças do plano no disco: reescreve conteúdo e renomeia arquivos
// (git mv, com fallback para rename simples fora de repo git). Retorna quantos
// arquivos foram escritos + renomeados.
func (p *Plan) Apply(root string) (int, error) {
	written := 0
	for _, fc := range p.Files {
		abs := filepath.Join(root, fc.Path)
		info, err := os.Stat(abs)
		mode := os.FileMode(0o644)
		if err == nil {
			mode = info.Mode()
		}
		if err := os.WriteFile(abs, []byte(fc.NewContent), mode); err != nil {
			return written, fmt.Errorf("escrever %s: %w", fc.Path, err)
		}
		written++
	}
	for _, r := range p.Renames {
		if err := gitMove(root, r.From, r.To); err != nil {
			return written, fmt.Errorf("renomear %s → %s: %w", r.From, r.To, err)
		}
		written++
	}
	return written, nil
}

// gitMove renomeia via `git mv` (preserva history + stage). FORA de um repositório,
// cai num os.Rename simples — ali o rename É o comportamento certo, e não há histórico
// a preservar.
//
// O que ele deliberadamente NÃO faz é usar o rename para contornar uma RECUSA do git.
// A versão anterior tratava as duas coisas como uma só ("não é repo git (ou git mv
// falhou) → rename direto"), e a diferença é grande: quando há repositório e o `git mv`
// falha, ele está recusando por um motivo — conflito em curso, lock do índice, arquivo
// fora do controle de versão. Mover assim mesmo deixa o índice inconsistente com o
// disco, em silêncio.
//
// O estrago era proporcional ao comando: o recode renomeia em MASSA. Bastava a recusa
// no meio de 40 renomeações para metade ir staged e metade ir untracked+deleted, com
// `Apply` retornando sucesso — e o rastro que permitiria desfazer é exatamente o que
// não foi criado.
func gitMove(root, from, to string) error {
	if dir := filepath.Dir(filepath.Join(root, to)); dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}
	if gitmeta.Verifica(root) != gitmeta.Disponível {
		// Sem git ou sem repositório: o rename é o único caminho, e é o correto.
		return os.Rename(filepath.Join(root, from), filepath.Join(root, to))
	}
	// SEM `-k`. O `-k` manda o git PULAR em silêncio o que não consegue mover (um
	// arquivo ainda não rastreado, por exemplo) e sair com status 0 — o rename não
	// acontece, `err` é nil, e o Apply conta como feito. Um projeto recém-inicializado,
	// onde nada foi commitado ainda, teria o recode anunciando "✓ 40 arquivos
	// reescritos" com os 40 parados no lugar.
	cmd := exec.Command("git", "mv", from, to)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	// Arquivo ainda não rastreado não é recusa a respeitar: não há histórico a
	// preservar, então o rename direto faz o que o usuário pediu — e o índice não fica
	// inconsistente, porque o git nunca soube deste arquivo.
	if strings.Contains(string(out), "not under version control") {
		return os.Rename(filepath.Join(root, from), filepath.Join(root, to))
	}
	return fmt.Errorf("`git mv %s %s` recusou: %s\n"+
		"   (há repositório aqui, então o arquivo NÃO foi movido por fora — mover à revelia "+
		"do git deixaria o índice divergindo do disco, e é o que torna difícil desfazer)",
		from, to, strings.TrimSpace(string(out)))
}
