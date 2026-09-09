// Package settings guarda a configuração LOCAL do agente — o que é dele e não do projeto.
//
// A distinção é a razão de o pacote existir. O `anchors.yaml` é a Estrutura: versionada,
// revisada, igual para todo mundo. O que vive aqui é o oposto — vale para UM agente numa
// máquina, e não deve chegar ao git.
//
// O caso que a criou: num projeto com vários devs, cada um pode rodar o seu agente. Nem
// todos podem decidir pelo produto. Um agente que pega um card `needs-user` e pergunta ao
// dev que o está rodando obtém uma resposta — e a resposta pode não ser a do dono do
// projeto. O escalonamento existe justamente para levar a pergunta a quem decide, e um
// agente prestativo demais o curto-circuita.
//
// Por isso o padrão é NÃO atuar. Quem pode decidir declara que pode; quem não declarou
// segue trabalhando nos cards comuns, e os escalonados esperam quem os resolve.
package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Dir é onde o estado local do Anchors vive — o mesmo do daemon, e no `.gitignore`.
const Dir = ".anchors"

// File é o arquivo de configuração do agente.
const File = "settings.yaml"

// Settings é o que o agente decidiu sobre si mesmo neste projeto.
type Settings struct {
	// Role é o PERFIL de quem opera este agente, e é dele que as capacidades derivam.
	//
	// Substituiu o `user_issues` booleano, que respondia uma pergunta só — quem decide o
	// produto — e deixava o resto implícito: o revisor de segurança e o de performance
	// liam a mesma régua, e o QA recebia o card do dev.
	//
	// O papel é o que a pessoa sabe dizer sobre si: "sou dev" é uma resposta, "atuo em
	// needs-user, não escrevo plano, reviso código" é um formulário que ninguém preenche
	// com cuidado. E o papel sobrevive ao Anchors ganhar capacidades novas — elas nascem
	// mapeadas aos perfis existentes, sem ninguém redeclarar nada.
	Role Role `yaml:"role,omitempty"`
	// UserIssues diz se ESTE agente atua nos cards escalonados (`needs-user`).
	//
	// É um ponteiro para distinguir três estados, e a distinção é o mecanismo: `nil` é
	// "nunca perguntei" — e é o único caso em que o agente deve perguntar. `false` e
	// `true` são decisões tomadas, e perguntar de novo seria ignorar a resposta.
	//
	// Sem o ponteiro, "não declarado" e "declarado como não" seriam o mesmo valor, e o
	// agente perguntaria a cada sessão a quem já disse não.
	UserIssues *bool `yaml:"user_issues"`
	// Agent é quem declarou, para o registro fazer sentido quando alguém o lê.
	Agent string `yaml:"agent,omitempty"`
	// DecidedAt é quando. Carimbado por quem escreve — o Anchors não lê o relógio.
	DecidedAt string `yaml:"decided_at,omitempty"`
}

// Path devolve o caminho do arquivo de configuração.
func Path(root string) string { return filepath.Join(root, Dir, File) }

// Load lê a configuração local. Ausência NÃO é erro: um projeto recém-clonado não tem o
// arquivo, e é exatamente o caso em que o agente precisa perguntar.
func Load(root string) (Settings, error) {
	b, err := os.ReadFile(Path(root))
	if err != nil {
		if os.IsNotExist(err) {
			return Settings{}, nil
		}
		return Settings{}, err
	}
	var s Settings
	if err := yaml.Unmarshal(b, &s); err != nil {
		return Settings{}, fmt.Errorf("%s: %w", Path(root), err)
	}
	return s, nil
}

// Save grava a configuração local.
//
// Cria o `.anchors/` se preciso — o diretório é do daemon e pode não existir num projeto
// onde o watcher nunca rodou.
func Save(root string, s Settings) error {
	if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	cabecalho := "# Configuração LOCAL deste agente — não vai para o git (`.anchors/` está\n" +
		"# no `.gitignore`). O que está aqui vale para UMA máquina, e não para o projeto:\n" +
		"# dois devs no mesmo repositório podem ter perfis diferentes.\n" +
		"#\n" +
		"# O PERFIL diz que trabalho é seu, e as capacidades derivam dele. Não é\n" +
		"# hierarquia — o arquiteto não manda no dev, ele responde outra pergunta — e não\n" +
		"# é permissão de repositório, que o git já controla.\n" +
		"#\n" +
		"# A diferença que mais aparece: só `product-owner` e `architect` atuam nos cards\n" +
		"# ESCALONADOS (`needs-user`), os que esperam decisão de quem conhece o produto.\n" +
		"# Os outros perfis escalam e seguem para o próximo card — porque a resposta de\n" +
		"# quem conhece o código é razoável, e vira decisão de produto sem passar pelo\n" +
		"# plano nem deixar rastro.\n" +
		"#\n" +
		"# Declare com `anchors settings role`; veja com `anchors settings show`.\n"
	return os.WriteFile(Path(root), append([]byte(cabecalho), b...), 0o644)
}

// Pode diz se o perfil declarado tem a capacidade.
//
// Sem perfil, NADA — o padrão fechado é o mesmo do booleano que este mecanismo substituiu, e
// pela mesma razão: o custo de errar para o lado aberto é alguém decidir sem autoridade, e
// isso é invisível depois do fato.
func (s Settings) Can(c Capability) bool {
	if s.Role == "" {
		return s.legacyUserIssues(c)
	}
	return s.Role.Can(c)
}

// legacyUserIssues lê o campo antigo, para quem já declarou.
//
// O `user_issues: true` de um `settings.yaml` escrito antes dos perfis continua valendo — e
// o comando pede o perfil na próxima vez. Ignorar o campo antigo faria quem já declarou
// perder a capacidade sem nada avisar, no meio de um trabalho.
func (s Settings) legacyUserIssues(c Capability) bool {
	return c == CapDecideProduct && s.UserIssues != nil && *s.UserIssues
}

// HandlesUserIssues diz se o agente atua nos escalonados.
//
// Não declarado = NÃO. O padrão fechado é deliberado: quem pode decidir declara que pode, e
// o custo de errar para o lado aberto é alguém decidir o produto sem autoridade — que é
// invisível depois do fato.
func (s Settings) HandlesUserIssues() bool {
	return s.Can(CapDecideProduct)
}

// Decided diz se a escolha já foi feita — é o que separa "perguntar" de "seguir".
//
// O perfil OU o campo antigo: quem declarou `user_issues` antes dos perfis não é perguntado
// de novo por causa da migração. O comando pede o perfil quando houver outra razão para
// perguntar.
func (s Settings) Decided() bool { return s.Role != "" || s.UserIssues != nil }

// Describe a decisão, para o comando reportar.
func (s Settings) Describe() string {
	if !s.Decided() {
		return "nenhum perfil declarado — o agente não sabe que trabalho é dele"
	}
	if s.Role == "" {
		// Declarou pelo campo antigo. Diz o que vale E que falta o perfil, senão a
		// migração ficaria invisível para quem lê.
		if s.HandlesUserIssues() {
			return "atua nos escalonados (pelo `user_issues` antigo) — declare o perfil " +
				"com `anchors settings role`"
		}
		return "não atua nos escalonados (pelo `user_issues` antigo) — declare o perfil " +
			"com `anchors settings role`"
	}
	por := ""
	if s.Agent != "" {
		por = " · " + s.Agent
	}
	return s.Role.Title() + por
}

// Bool devolve o ponteiro que o campo pede, para quem monta o Settings.
func Bool(v bool) *bool { return &v }

// ParseAnswer lê o que uma pessoa digitou e devolve a decisão.
//
// Aceita as duas línguas e as formas curtas porque a pergunta é feita no terminal, e quem
// responde digita o que lhe vem primeiro. Devolve `nil` para o que não entendeu — e aí o
// comando pergunta de novo, em vez de assumir.
func ParseAnswer(r string) *bool {
	switch strings.ToLower(strings.TrimSpace(r)) {
	case "s", "sim", "y", "yes":
		return Bool(true)
	case "n", "nao", "não", "no":
		return Bool(false)
	}
	return nil
}
