package health

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// --- a aprovação que o próprio autor não pode dar ---
//
// O fluxo pressupõe que OUTRO agente revise, e `required_approvals: 1` é o que impede o
// merge sem revisão. Mas o GitHub RECUSA aprovação do próprio PR — é regra de plataforma,
// sem configuração que a contorne — e agentes na mesma máquina compartilham a conta.
//
// O resultado é um fluxo travado: o card chega a `ready-to-review`, o revisor confronta e
// aprova… e não consegue. Nenhum PR do projeto avança pelo caminho normal.
//
// Isso não é defeito do Anchors nem do projeto: é uma restrição conhecida, com dois
// contornos que dependem de quem opera. O papel do doctor é DETECTAR a situação e dizer
// qual dos dois cabe — em vez de deixar alguém descobrir no meio de um merge.

// PodeIgnorarProtecao diz se a conta atual consegue mesclar por cima da exigência de
// aprovação: precisa ser admin do repositório E a proteção não pode alcançar admins.
//
// São duas condições porque `enforce_admins: true` é o modo em que o dono do repositório
// escolheu não ter escape — e nesse caso ser admin não ajuda.
func PodeIgnorarProtecao(repo, branch string) (bool, string) {
	if _, err := exec.LookPath("gh"); err != nil {
		return false, "o `gh` não está no PATH"
	}
	login, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return false, "não deu para descobrir a conta atual"
	}
	perm, err := exec.Command("gh", "api",
		"repos/"+repo+"/collaborators/"+strings.TrimSpace(string(login))+"/permission",
		"--jq", ".permission").Output()
	if err != nil || strings.TrimSpace(string(perm)) != "admin" {
		return false, "a conta atual não é admin do repositório"
	}
	out, err := exec.Command("gh", "api",
		"repos/"+repo+"/branches/"+branch+"/protection",
		"--jq", ".enforce_admins.enabled").Output()
	if err != nil {
		return false, "não deu para ler a proteção do branch"
	}
	if strings.TrimSpace(string(out)) == "true" {
		return false, "`enforce_admins` está ligado — nem admin ignora a exigência"
	}
	return true, ""
}

// checkAprovacaoAlcancavel avisa quando a exigência de aprovação não tem como ser
// cumprida: o autor não pode aprovar o próprio PR, e não há escape configurado.
func checkAprovacaoAlcancavel(cfg *config.Config) []Finding {
	if cfg == nil || cfg.Workflow == nil || cfg.Workflow.AprovacoesExigidas() == 0 {
		return nil // zero exigido: não há o que ficar inalcançável
	}
	if _, err := exec.LookPath("gh"); err != nil {
		return nil // sem `gh` o doctor já reclama noutro achado
	}
	repo := cfg.Workflow.Repo
	branch := cfg.Workflow.BranchDeIntegracao()

	if ok, _ := PodeIgnorarProtecao(repo, branch); ok {
		// Admin com escape: o fluxo funciona, e o merge usa `gh pr merge --admin`. Não é
		// achado — é o contorno previsto, e repeti-lo a cada `doctor` viraria ruído.
		return nil
	}
	return []Finding{{"aprovacao-inalcancavel", Warn, repo,
		fmt.Sprintf("`required_approvals: %d` exige aprovação, e o GitHub RECUSA que o autor "+
			"aprove o próprio PR — regra de plataforma, sem configuração que a contorne. "+
			"Como os agentes desta máquina usam a mesma conta, nenhum PR avança pelo caminho "+
			"normal. Duas saídas: (1) uma conta de serviço para os agentes, ou (2) "+
			"`required_approvals: 0` no anchors.yaml, deixando a revisão ser cobrada pelo "+
			"estado do card em vez da aprovação do GitHub. Rode `anchors doctor --fix` para "+
			"a segunda", cfg.Workflow.AprovacoesExigidas())}}
}

// desligaExigenciaDeAprovacao aplica a saída (2): zera a exigência no GitHub.
//
// O `anchors.yaml` continua sendo a fonte da verdade — quem quiser a exigência de volta
// declara `required_approvals: 1` e roda o fix de novo. Aqui só se ajusta o GitHub para
// não travar um fluxo que ele mesmo impede de cumprir.
func DesligaExigenciaDeAprovacao(repo, branch string) error {
	body := `{"required_status_checks":null,"enforce_admins":false,` +
		`"required_pull_request_reviews":{"required_approving_review_count":0},` +
		`"restrictions":null}`
	cmd := exec.Command("gh", "api", "--method", "PUT",
		"repos/"+repo+"/branches/"+branch+"/protection", "--input", "-")
	cmd.Stdin = strings.NewReader(body)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", branch, strings.TrimSpace(string(out)))
	}
	return nil
}

// aprovacaoAtual lê quantas aprovações o branch exige hoje. Devolve -1 quando não dá para
// saber — que é diferente de zero, e o chamador não pode confundir os dois.
func aprovacaoAtual(repo, branch string) int {
	out, err := exec.Command("gh", "api",
		"repos/"+repo+"/branches/"+branch+"/protection",
		"--jq", ".required_pull_request_reviews.required_approving_review_count").Output()
	if err != nil {
		return -1
	}
	var n int
	if json.Unmarshal([]byte(strings.TrimSpace(string(out))), &n) != nil {
		return -1
	}
	return n
}
