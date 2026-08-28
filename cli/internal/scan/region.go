package scan

import (
	"crypto/sha256"
	"encoding/hex"
	"path"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// Região — a delimitação DECLARADA de um código de cenário no fonte (TRACEABILITY §3).
//
// As outras três formas da identidade são delimitadas pela sintaxe onde vivem (a linha
// da spec, o `Scenario` da feature, o bloco do caso de teste). O código era a exceção: o
// marcador dizia onde o requisito começa e nada dizia onde termina, então a identidade no
// fonte era um CONJUNTO ("este arquivo realiza CODEX-A03") e não um INTERVALO ("de 471 a
// 545"). Basta para perguntas de existência e falha nas de mudança: um arquivo com catorze
// requisitos dá catorze respostas quando a verdade é uma.
//
// A região usa a sintaxe que os editores já colapsam, com o código repetido no fecho:
//
//	// #region [CODEX-A03]: confirmar persiste o lançamento e fecha.
//	…
//	// #endregion [CODEX-A03]
//
// O fecho repete o código porque com fecho ANÔNIMO um aninhamento trocado é indetectável:
// a contagem de abre/fecha bate e o intervalo medido é o errado, em silêncio. Com o código,
// o confronto é nome contra nome e a mensagem aponta a linha.
type Regiao struct {
	Code  string // o código de cenário delimitado
	Start int    // linha do #region (1-indexada)
	End   int    // linha do #endregion (1-indexada)
	Rev   string // sha256 curto do conteúdo do intervalo, incluindo os marcadores
}

// ErroRegiao — um defeito de pareamento. Não é o mesmo que ausência de região: um arquivo
// sem nenhuma região é legítimo (a delimitação é OPCIONAL, e vale o rev do arquivo).
// Região aberta e nunca fechada, fecho sem abertura, ou fecho com código diferente do que
// abriu são erros de escrita — e é o que o gate `region-pair-honored` reprova.
type ErroRegiao struct {
	Linha int
	Kind  string // "sem-fecho" | "fecho-orfao" | "fecho-trocado"
	Code  string // o código da abertura (ou do fecho órfão)
	Achou string // no fecho-trocado: o código que o fecho traz
}

// A gramática do marcador. `#region`/`#endregion` é o par que VS Code e Visual Studio
// colapsam nativamente; aceitamos com ou sem espaço após o marcador de comentário, e o
// código entre colchetes — que o delimita sem ambiguidade com o texto da descrição.
//
// Deliberadamente NÃO amarramos ao prefixo de comentário da linguagem: `config.CommentMarkers`
// varia (`//`, `#`, `--`, `<!--`) e o que importa aqui é reconhecer a anotação, não validar a
// sintaxe do arquivo. Exigir o prefixo certo faria o parser calar num arquivo de extensão
// desconhecida — o modo de falha que este engine já mediu duas vezes: vocabulário errado não
// dá erro, dá silêncio.
// Compilados por CHAMADA, não em `var`: o comprimento do código vem da config do projeto
// (`code_lengths`), que é carregada DEPOIS da inicialização dos globais. Um `var` aqui
// congelaria o padrão default e a declaração do projeto não teria efeito — o modo de falha
// que este campo existe para consertar.
func reRegiaoAbre() *regexp.Regexp {
	return regexp.MustCompile(`#region\s*\[\s*([A-Z0-9]` + config.CodeLengthPattern() + `-[A-Za-z0-9-]+)\s*\]`)
}

func reRegiaoFecha() *regexp.Regexp {
	return regexp.MustCompile(`#endregion\s*(?:\[\s*([A-Z0-9]` + config.CodeLengthPattern() + `-[A-Za-z0-9-]+)\s*\])?`)
}

// Regioes extrai as regiões de um conteúdo, junto dos erros de pareamento.
//
// Aninhamento é suportado por uma pilha: uma região interna fecha antes da externa, e o
// fecho tem de casar com o topo. O `rev` de cada região cobre o intervalo INCLUINDO os
// marcadores — mudar a descrição do próprio requisito é mudança do requisito.
func Regioes(content string) ([]Regiao, []ErroRegiao) {
	linhas := strings.Split(content, "\n")
	type aberta struct {
		code  string
		linha int
	}
	var pilha []aberta
	var out []Regiao
	var erros []ErroRegiao

	for i, l := range linhas {
		n := i + 1
		// o fecho é testado ANTES da abertura: `#endregion` contém a substring `region`,
		// e testar na ordem inversa faria todo fecho parecer uma abertura sem código.
		if m := reRegiaoFecha().FindStringSubmatch(l); m != nil {
			if len(pilha) == 0 {
				erros = append(erros, ErroRegiao{Linha: n, Kind: "fecho-orfao", Achou: m[1]})
				continue
			}
			topo := pilha[len(pilha)-1]
			if m[1] != "" && m[1] != topo.code {
				// Erro REPORTADO mas não fatal: fechamos o topo de todo jeito, para que um
				// fecho trocado não faça toda região seguinte parecer desbalanceada e
				// produza uma cascata de erros derivados de um só defeito.
				erros = append(erros, ErroRegiao{
					Linha: n, Kind: "fecho-trocado", Code: topo.code, Achou: m[1],
				})
			}
			pilha = pilha[:len(pilha)-1]
			out = append(out, Regiao{
				Code:  topo.code,
				Start: topo.linha,
				End:   n,
				Rev:   hashLinhas(linhas[topo.linha-1 : n]),
			})
			continue
		}
		if m := reRegiaoAbre().FindStringSubmatch(l); m != nil {
			pilha = append(pilha, aberta{code: m[1], linha: n})
		}
	}
	// o que sobrou na pilha nunca fechou
	for _, a := range pilha {
		erros = append(erros, ErroRegiao{Linha: a.linha, Kind: "sem-fecho", Code: a.code})
	}
	return out, erros
}

// hashLinhas — sha256 truncado em 12 hex, o mesmo comprimento do `rev` de nó, para que os
// dois sejam comparáveis a olho no mapa e no mesmo formato para quem lê.
func hashLinhas(linhas []string) string {
	h := sha256.Sum256([]byte(strings.Join(linhas, "\n")))
	return hex.EncodeToString(h[:])[:12]
}

// composeRefRE — a COMPOSIÇÃO de um roteiro: `runFlow: <caminho>` e a forma com
// `file:` (que é a mesma coisa quando o runFlow leva `env:`). O caminho é literal no
// YAML, então a aresta é DECLARADA — não há inferência aqui.
//
// Aceita valor entre quotes ou nu, e ignora comentário à direita. Só `.yaml`/`.yml`:
// `runScript:` aponta para um `.js` que é dado de entrada, não composição de roteiro.
var composeRefRE = regexp.MustCompile(`(?m)^\s*(?:-\s*)?(?:runFlow|file):\s*["']?([^"'#\s]+\.ya?ml)["']?`)

// ComposeRefs devolve os roteiros que ESTE roteiro compõe, como caminhos relativos à
// raiz — a base do `depends-on` de composição.
//
// Existe porque a composição de flows era invisível ao mapa, e a invisibilidade tinha um
// custo medido: `utils/login.yaml` é usado por 286 roteiros, e alterá-lo não vencia a
// evidência de nenhum deles. Numa sessão real toquei esse util DEPOIS de 8 suítes já
// terem passado; revalidei as 13 por lembrança minha, não por ferramenta — e somar
// medições de códigos diferentes produz um placar que nunca existiu.
//
// O alvo é resolvido contra o DIRETÓRIO do roteiro (os caminhos são relativos a ele,
// no estilo `../../utils/login.yaml`), não contra a raiz: resolver na raiz produziria
// caminhos que não existem e a aresta seria descartada em silêncio.
func ComposeRefs(content, rel string) []string {
	ms := composeRefRE.FindAllStringSubmatch(content, -1)
	if len(ms) == 0 {
		return nil
	}
	dir := path.Dir(rel)
	seen := map[string]bool{}
	var out []string
	for _, m := range ms {
		alvo := path.Clean(path.Join(dir, m[1]))
		if alvo == rel || seen[alvo] { // auto-referência não é dependência
			continue
		}
		seen[alvo] = true
		out = append(out, alvo)
	}
	return out
}
