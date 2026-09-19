package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/settings"
)

// AgentID devolve o identificador deste agente (<host>/<sessão>).
func AgentID() string {
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	sessao := strings.TrimSpace(os.Getenv("ANCHORS_SESSION"))
	if sessao == "" {
		sessao = os.Getenv("USER")
		if sessao == "" {
			sessao = "default"
		}
	}
	return host + "/" + sessao
}

// RoleList lista os perfis conhecidos e suas atribuições.
func RoleList() string {
	var b strings.Builder
	b.WriteString(i18n.T("role.list_header"))
	for _, r := range settings.KnownRoles() {
		fmt.Fprintf(&b, "  %-22s %s\n", string(r), r.Does())
	}
	return b.String()
}

// PrintRole mostra o que o perfil declarado significa.
func PrintRole(r settings.Role) {
	fmt.Println(i18n.T("role.declared_title", r.Title()))
	fmt.Printf("  %s\n", r.Does())

	if r.Can(settings.CapDecideProduct) {
		fmt.Println(i18n.T("role.decides_product_notice"))
	} else {
		fmt.Println(i18n.T("role.does_not_decide_product_notice"))
	}
	if lente := r.Lens(); lente != "" {
		fmt.Println(i18n.T("role.lens_notice", lente))
	}
}

// AskRole pergunta o perfil no terminal.
func AskRole() (settings.Role, error) {
	in := bufio.NewReader(os.Stdin)
	for tentativa := 0; tentativa < 3; tentativa++ {
		fmt.Println(i18n.T("role.prompt_which"))
		fmt.Print(RoleList())
		fmt.Print(i18n.T("role.prompt_field"))

		linha, err := in.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("%s: %w", i18n.T("role.read_response_error"), err)
		}
		if r := settings.ParseRole(linha); r != "" {
			return r, nil
		}
		fmt.Print(i18n.T("role.not_recognized", strings.TrimSpace(linha)))
	}
	return "", fmt.Errorf("%s", i18n.T("role.unrecognized_error"))
}

// InteractiveTerminal diz se há alguém do outro lado para responder.
func InteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	nul, err := os.Stat(os.DevNull)
	if err != nil {
		return false
	}
	return !os.SameFile(fi, nul)
}

// DecidesProduct diz se este agente declarou que decide o rumo do produto.
func DecidesProduct(root string) bool {
	s, err := settings.Load(root)
	if err != nil {
		return false
	}
	return s.HandlesUserIssues()
}
