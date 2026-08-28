package health

import (
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
)

// --- ambiente do modo `github` (BOOTSTRAP.md §7.13/§7.14) ---
//
// O fluxo do modo `github` pressupõe três pipelines no lugar. Se
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
	return checkPipelines(root)
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
