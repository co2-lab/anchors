package health

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
)

// --- ambiente do modo `github` (BOOTSTRAP.md §7.13/§7.14) ---
//
// O fluxo do modo `github` pressupõe um board configurado e três pipelines no lugar. Se
// qualquer peça faltar, nada falha ruidosamente: o fluxo simplesmente NÃO ACONTECE. Um
// pipeline de identificação ausente não gera erro — gera silêncio, e os artefatos ficam
// sem card para sempre.
//
// É a mesma classe de problema da falta de git, e recebe o mesmo tratamento: o doctor
// avisa com antecedência, antes de alguém descobrir no meio de um trabalho.
//
// Só roda no modo `github`. Cobrar board de um projeto que declarou `mode: local` seria
// ruído garantido — e ruído recorrente treina a equipe a ignorar o doctor.
func checkAmbienteGitHub(cfg *config.Config, root string) []Finding {
	if cfg == nil || !cfg.ModoGitHub() {
		return nil
	}
	var out []Finding
	out = append(out, checkPipelines(root)...)
	out = append(out, checkBoard(cfg)...)
	return out
}

// checkPipelines confere o que dá para conferir lendo o disco: os três workflows existem,
// e os que precisam de serialização a declaram.
func checkPipelines(root string) []Finding {
	var out []Finding
	for _, w := range initx.FaltaWorkflow(root) {
		out = append(out, Finding{"pipeline-ausente", Warn, w.Arquivo,
			"pipeline do fluxo não existe em " + initx.DirWorkflows + " — sem ele, " +
				w.Papel + " não acontece (e não falha: só não acontece). " +
				"Rode `anchors doctor --fix`"})
	}
	// PRESENTE mas sem serialização é o pior caso: parece configurado, e devolve a
	// corrida que o pipeline existia para eliminar — em silêncio, porque o arquivo ESTÁ
	// lá. Achado próprio, distinto de "está faltando".
	for _, w := range initx.SemConcurrency(root) {
		out = append(out, Finding{"pipeline-sem-serializacao", Warn, w.Arquivo,
			"o pipeline existe mas não declara `concurrency` com `cancel-in-progress: false` — " +
				"duas execuções simultâneas podem atribuir o mesmo card a dois agentes, " +
				"ou criar o mesmo card duas vezes"})
	}
	return out
}

// checkBoard confere o GitHub Project: existe, está acessível, e tem o campo `Status` com
// as colunas que os pipelines movem.
//
// Quando falta escopo no token, ele DIZ que não conseguiu conferir. Calar seria pior do
// que não verificar: um doctor silencioso sobre o board é lido como "board OK", e é o tipo
// de silêncio que tranquiliza sem ter olhado — a mesma régua do `DirtyCount`.
func checkBoard(cfg *config.Config) []Finding {
	if _, err := exec.LookPath("gh"); err != nil {
		return []Finding{{"board-nao-verificado", Warn, "gh",
			"o `gh` não está no PATH — não deu para conferir o Project. " +
				"O modo `github` inteiro depende dele (WORKFLOW.md §4)"}}
	}
	dono, _, ok := strings.Cut(cfg.Workflow.Repo, "/")
	if !ok {
		return nil // repo malformado já é erro de config, cobrado no Load
	}

	saida, err := exec.Command("gh", "project", "list", "--owner", dono, "--format", "json").Output()
	if err != nil {
		// O erro mais comum aqui é escopo faltando, e a mensagem do `gh` já o nomeia.
		// Repassá-la é melhor do que traduzir: ela traz o comando exato do conserto.
		det := "não deu para listar os Projects de " + dono
		if strings.Contains(erroDe(err), "read:project") || strings.Contains(erroDe(err), "scope") {
			det += ": falta escopo no token — rode `gh auth refresh -s read:project`"
		}
		return []Finding{{"board-nao-verificado", Warn, dono,
			det + ". O board pode estar correto ou não — isto NÃO é um \"está tudo certo\""}}
	}

	var lista struct {
		Projects []struct {
			Number int    `json:"number"`
			Title  string `json:"title"`
		} `json:"projects"`
	}
	if json.Unmarshal(saida, &lista) != nil || len(lista.Projects) == 0 {
		return []Finding{{"board-ausente", Warn, dono,
			"nenhum GitHub Project encontrado em " + dono + " — o estado do trabalho é a " +
				"COLUNA do board (BOOTSTRAP.md §7.13), então sem ele os pipelines não têm " +
				"onde registrar o andamento"}}
	}
	return nil
}

// erroDe extrai o stderr de um erro de exec, onde o `gh` escreve a causa. Sem isto a
// mensagem chegaria como "exit status 1", que não diz nada a ninguém.
func erroDe(err error) string {
	var ee *exec.ExitError
	if ok := asExitError(err, &ee); ok {
		return string(ee.Stderr)
	}
	return err.Error()
}

func asExitError(err error, target **exec.ExitError) bool {
	if ee, ok := err.(*exec.ExitError); ok {
		*target = ee
		return true
	}
	return false
}
