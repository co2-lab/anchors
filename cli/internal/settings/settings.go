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
		"# no `.gitignore`). O que está aqui vale para UMA máquina, e não para o projeto.\n" +
		"#\n" +
		"# `user_issues` diz se este agente atua nos cards escalonados (`needs-user`), que\n" +
		"# são os que esperam decisão de quem conhece o produto. O padrão é não: um agente\n" +
		"# que os pega e pergunta a quem o está rodando obtém uma resposta que pode não ser\n" +
		"# a do dono do projeto — e o escalonamento existe para levar a pergunta a quem\n" +
		"# decide.\n"
	return os.WriteFile(Path(root), append([]byte(cabecalho), b...), 0o644)
}

// HandlesUserIssues diz se o agente atua nos escalonados.
//
// Não declarado = NÃO. O padrão fechado é deliberado: quem pode decidir declara que pode, e
// o custo de errar para o lado aberto é alguém decidir o produto sem autoridade — que é
// invisível depois do fato.
func (s Settings) HandlesUserIssues() bool {
	return s.UserIssues != nil && *s.UserIssues
}

// Decided diz se a escolha já foi feita — é o que separa "perguntar" de "seguir".
func (s Settings) Decided() bool { return s.UserIssues != nil }

// Describe a decisão, para o comando reportar.
func (s Settings) Describe() string {
	if !s.Decided() {
		return "não declarado — o agente ainda não decidiu se atua nos escalonados"
	}
	if s.HandlesUserIssues() {
		por := ""
		if s.Agent != "" {
			por = " (declarado por " + s.Agent + ")"
		}
		return "ATUA nos cards escalonados" + por
	}
	return "NÃO atua nos cards escalonados — eles esperam quem decide o produto"
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
