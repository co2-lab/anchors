package health

import (
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
)

// --- ambiente do modo `github` (BOOTSTRAP.md §7.13/§7.14) ---
//
// O fluxo do modo `github` pressupõe quatro pipelines no lugar. Se
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
	// O BOARD não é conferido, e isso é deliberado: o estado do trabalho é uma LABEL, e
	// o Project é só espelho opcional (ver BOOTSTRAP.md §7.13). Cobrar um board que o
	// fluxo não precisa produziria um achado que ninguém precisa resolver — e ruído
	// recorrente treina a equipe a ignorar o doctor.
	out := checkPipelines(root, cfg)
	out = append(out, checkProtecaoDeBranch(cfg)...)
	return append(out, checkAprovacaoAlcancavel(cfg)...)
}

// checkPipelines confere o que dá para conferir lendo o disco: os três workflows existem,
// e os que precisam de serialização a declaram.
func checkPipelines(root string, cfg *config.Config) []Finding {
	var out []Finding
	for _, w := range initx.FaltaWorkflow(root) {
		out = append(out, Finding{"pipeline-ausente", Warn, w.Arquivo,
			"pipeline do fluxo não existe em " + initx.DirWorkflows + " — sem ele, " +
				w.Papel + " não acontece (e não falha: só não acontece). " +
				"Rode `anchors doctor --fix`"})
	}
	// DESATUALIZADO: o pipeline é do Anchors (marcador intacto) e o template mudou desde
	// que ele foi instalado. É como uma correção chega a quem já instalou — sem isto, um
	// defeito corrigido na fonte continua rodando no projeto para sempre, e nada avisa.
	//
	// Só vale para pipeline NÃO editado: um que o time customizou é dele, e a diferença
	// em relação ao template é a customização, não atraso.
	for _, w := range initx.WorkflowsDesatualizados(root, cfg) {
		out = append(out, Finding{"pipeline-desatualizado", Warn, w.Arquivo,
			"o pipeline é o template do Anchors e ficou para trás — uma correção no " +
				"desenho do fluxo não chegou a este projeto. Rode `anchors doctor --fix` " +
				"(se você editou o arquivo, remova a linha `# anchors:template` e ele " +
				"passa a ser seu)"})
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

// checkProtecaoDeBranch confere que a `main` exige PR.
//
// É a regra que o fluxo inteiro pressupõe, e a única cuja ausência não produz erro em
// lugar nenhum: sem proteção, um push direto na main funciona — e pula o card, pula a
// revisão, e pula o pipeline de identificação (que dispara na ABERTURA do PR, porque o
// push na main acontece depois do merge, quando o trabalho já terminou).
//
// O silêncio aqui é o mais caro do fluxo: tudo parece funcionar, e o ciclo de governança
// simplesmente não acontece.
func checkProtecaoDeBranch(cfg *config.Config) []Finding {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil // sem `gh` o doctor já reclama noutro achado; não duplicar
	}
	repo := cfg.Workflow.Repo
	out, err := exec.Command("gh", "api",
		"repos/"+repo+"/branches/main/protection",
		"--jq", ".required_pull_request_reviews != null").Output()
	if err != nil {
		// 404 é a resposta para branch sem proteção — e é o achado, não um erro.
		return []Finding{{"main-sem-protecao", Warn, repo,
			"a `main` não exige pull request — um push direto nela pula o card, a revisão " +
				"e o pipeline de identificação (que dispara na ABERTURA do PR). Nada falha: " +
				"o ciclo de governança simplesmente não acontece. Rode `anchors doctor --fix`"}}
	}
	if strings.TrimSpace(string(out)) != "true" {
		return []Finding{{"main-sem-protecao", Warn, repo,
			"a `main` tem proteção, mas não exige pull request — ver `anchors doctor --fix`"}}
	}
	return nil
}
