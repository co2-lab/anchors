// Package change materializa o REGISTRO DE ENTREGA: o que um agente diz ter feito, ao
// terminar uma etapa. É a peça que faltava para fechar o ciclo de vida — sem ela, o
// trabalho acabava quando o código nascia, e ninguém confrontava o que foi entregue.
//
// Por que um arquivo, e não uma chamada direta ao revisor:
//
//  1. O review precisa de ESCOPO. Sem o registro, quem revisa tem de adivinhar o que
//     mudou — e revisar "o repositório" é revisar nada. O change diz exatamente quais
//     arquivos entraram, sob qual spec, com qual intenção declarada.
//  2. O registro sobrevive à sessão. Um review pedido a quente depende de alguém lembrar
//     de pedir; medido num E2E real, foi preciso o usuário cobrar. Arquivo no disco é
//     visto pelo watcher, e o watcher não esquece.
//  3. A intenção DECLARADA vira material de confronto. O agente escreve o que acha que
//     fez; o revisor confronta contra o que está no disco. A divergência entre os dois é,
//     por si só, um achado — e não existe forma de detectá-la sem ter as duas versões.
//
// Diferente de um patch: o Anchors não intermedia a escrita (o agente edita o repo
// direto). O change não CARREGA a mudança — ele a DESCREVE e aponta. Quem quiser o diff
// usa git; o valor aqui é a intenção + o escopo, que o git não tem.
//
// Changes vivem em `changes/` (conteúdo de PROJETO, versionável, auditável), não em
// `.anchors/` (estado efêmero do daemon) — pela mesma razão das issues: é registro do que
// aconteceu com o produto, não do que a ferramenta está fazendo agora.
package change

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Dir é a raiz dos registros de entrega, relativa à raiz do projeto.
const Dir = "changes"

// ReviewedDir guarda os que já foram revisados — o histórico. Mover para cá é o que
// fecha o registro, do mesmo jeito que a pasta codifica o estado de uma issue.
const ReviewedDir = "changes/reviewed"

// Change — uma entrega declarada.
type Change struct {
	// Etapa que produziu (spec|code|feature|test|plan) — o mesmo vocabulário de artefato
	// do `anchors work`, para o revisor saber que régua se aplica.
	Stage string
	// Unit é o alvo: o arquivo de CÓDIGO que identifica a unidade de propósito. É por ele
	// que o revisor alcança a trinca inteira.
	Unit string
	// Files são os arquivos tocados nesta entrega.
	Files []string
	// Intent é o que o autor diz ter feito — em uma ou duas frases, na língua do projeto.
	// É a metade declarada do confronto.
	Intent string
	// Decisions são as escolhas que a régua NÃO decidiu e o autor tomou. Este campo existe
	// porque a decisão silenciosa é a origem da maior parte dos defeitos que atravessam
	// gates: tudo existe, tudo se referencia, e alguém escolheu sozinho.
	Decisions []string
	// Uncovered é o que o autor SABE que não está provado (regra sem teste, caso de borda
	// não coberto). Declarar é o que separa dívida assumida de esquecimento.
	Uncovered []string
	Date      string // AAAA-MM-DD — carimbado por quem cria
	Agent     string // quem entregou (opcional; útil quando há vários workers)
}

// Key é a identidade estável: <stage>--<unidade>. Duas entregas da mesma etapa sobre a
// mesma unidade colidem de propósito — a segunda SUBSTITUI a primeira, porque o que
// importa é o estado atual da entrega, não o histórico de tentativas.
func (c Change) Key() string {
	return c.Stage + "--" + slug(c.Unit)
}

// Path é onde o registro vive enquanto não foi revisado.
func (c Change) Path(root string) string {
	return filepath.Join(root, Dir, c.Key()+".md")
}

// Render escreve o registro em markdown. O formato é legível por pessoa e parseável por
// máquina — mesma escolha das issues: quem abre o arquivo entende sem ferramenta.
func (c Change) Render() string {
	var b strings.Builder
	b.WriteString("<!-- @anchors\n  layer: change\n")
	fmt.Fprintf(&b, "  stage: %s\n  unit: %s\n  date: %s\n", c.Stage, c.Unit, c.Date)
	if c.Agent != "" {
		fmt.Fprintf(&b, "  agent: %s\n", c.Agent)
	}
	b.WriteString("-->\n")
	fmt.Fprintf(&b, "# Entrega: %s de `%s`\n\n", c.Stage, c.Unit)

	b.WriteString("## O que foi feito\n\n" + strings.TrimSpace(c.Intent) + "\n\n")

	b.WriteString("## Arquivos\n\n")
	for _, f := range c.Files {
		fmt.Fprintf(&b, "- `%s`\n", f)
	}

	// As duas seções seguintes são o que torna o registro útil ao revisor. Vazias, elas
	// AFIRMAM algo (não houve decisão livre; nada ficou sem prova) — o que é diferente de
	// omiti-las. Por isso são sempre emitidas.
	b.WriteString("\n## Decisões que a régua não decidiu\n\n")
	if len(c.Decisions) == 0 {
		b.WriteString("nenhuma — tudo o que foi feito estava decidido na spec/guides.\n")
	} else {
		for _, d := range c.Decisions {
			fmt.Fprintf(&b, "- %s\n", d)
		}
	}

	b.WriteString("\n## O que NÃO está provado\n\n")
	if len(c.Uncovered) == 0 {
		b.WriteString("nada — cada regra desta entrega tem cenário e teste que a exercita.\n")
	} else {
		for _, u := range c.Uncovered {
			fmt.Fprintf(&b, "- %s\n", u)
		}
	}
	return b.String()
}

// Save grava o registro, criando o diretório se preciso.
func Save(root string, c Change) (string, error) {
	p := c.Path(root)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, []byte(c.Render()), 0o644); err != nil {
		return "", err
	}
	return p, nil
}

// Pending lista os registros ainda não revisados, em ordem estável.
func Pending(root string) ([]string, error) {
	dir := filepath.Join(root, Dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // sem entregas registradas ainda
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	sort.Strings(out)
	return out, nil
}

// MarkReviewed move o registro para o histórico. Mover (em vez de apagar) preserva o
// rastro: dá para responder "esta unidade já passou por review, e quando".
func MarkReviewed(root, path string) (string, error) {
	dest := filepath.Join(root, ReviewedDir, filepath.Base(path))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(path, dest); err != nil {
		return "", err
	}
	return dest, nil
}

// slug normaliza um caminho para caber num nome de arquivo, preservando a legibilidade.
func slug(s string) string {
	s = strings.TrimSuffix(s, filepath.Ext(s))
	r := strings.NewReplacer("/", "-", "\\", "-", " ", "-", ".", "-")
	return strings.Trim(r.Replace(s), "-")
}
