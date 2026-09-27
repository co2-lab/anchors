package gate

import (
	"os"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// O self-healer: alguns checks são REPARÁVEIS — a correção é mecânica e segura (ex.:
// escrever a data do commit no `updated_at`). `anchors check --fix` aplica esses
// reparos. Só checks com um fixer REGISTRADO são reparáveis; os demais só reportam.

// FixResult descreve um reparo aplicado (ou tentado).
type FixResult struct {
	Gate   string
	Target string
	Fixed  bool
	Detail string
}

// fixers: check → função que repara o arquivo. Recebe o conteúdo e o contexto,
// devolve o novo conteúdo e se mudou algo. Só reparos determinísticos e seguros.
var fixers = map[string]func(content string, n mapx.Node, root string) (string, bool){
	"updated-at-atual": fixUpdatedAt,
}

// Fixable diz se um check tem reparo automático.
func Fixable(check string) bool {
	_, ok := fixers[check]
	return ok
}

// Fix roda o FIXER de cada gate reparável sobre todo nó a que o gate se aplica e grava
// as correções em disco. Devolve o que foi consertado (ou tentado). Não consulta o
// veredito do gate: é o fixer quem decide se há o que corrigir, e um nó já correto
// não muda e não entra no resultado.
func Fix(gates []config.Gate, nodes []mapx.Node, root string) []FixResult {
	var out []FixResult
	for _, g := range gates {
		fix, ok := fixers[g.Check]
		if !ok {
			continue
		}
		for _, n := range nodes {
			if !applies(g, n, root) {
				continue
			}
			path := filepath.Join(root, n.ID)
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			newContent, changed := fix(string(content), n, root)
			if !changed {
				continue
			}
			if err := os.WriteFile(path, []byte(newContent), 0o644); err != nil {
				// Os detalhes passam por i18n: eram literais em português, impressos
				// assim em projeto de qualquer `lang:`.
				out = append(out, FixResult{g.Name, n.ID, false, i18n.T("gate.fix.write_failed", err)})
				continue
			}
			out = append(out, FixResult{g.Name, n.ID, true, i18n.T("gate.fix.updated_at_fixed")})
		}
	}
	return out
}

// fixUpdatedAt reescreve o `updated_at` do header com a data do último commit — só se
// diverge e o arquivo está commitado (nada a corrigir num arquivo com edição pendente,
// cuja data ainda vai mudar). Não cria o campo se ele não existe (isso é trabalho do
// autor + o gate header-conforme); só CORRIGE um valor errado.
func fixUpdatedAt(content string, n mapx.Node, root string) (string, bool) {
	m := updatedAtRE.FindStringSubmatchIndex(content)
	if m == nil {
		return content, false // sem campo — nada a corrigir aqui (o autor cria)
	}
	// a data CORRETA: hoje se há edição não-commitada (alteração em curso), senão a
	// data do último commit. Espelha a regra do checkUpdatedAt.
	var correct string
	mudou, sabido := gitmeta.UncommittedChanges(root, n.ID)
	if !sabido {
		return content, false // sem git: não há data correta a apurar, e chutar seria pior
	}
	if mudou {
		correct = gitmeta.Today()
	} else if d, ok := gitmeta.LastCommitDate(root, n.ID); ok {
		correct = d
	} else {
		return content, false // sem commit e sem edição — nada a comparar
	}
	// m[2]:m[3] é o grupo capturado (a data). Substitui só ela se diverge.
	if content[m[2]:m[3]] == correct {
		return content, false
	}
	return content[:m[2]] + correct + content[m[3]:], true
}
