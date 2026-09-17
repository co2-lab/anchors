package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// `anchors board serve` — o board que NÃO é foto.
//
// O board publicado no GitHub Pages é gerado por pipeline: ele roda, produz o
// `board.json`, publica. Entre uma execução e outra, o que se vê é passado — e num
// projeto com nove agentes entregando, minutos bastam para a foto mentir sobre quem tem
// qual card.
//
// ESTE COMANDO SERVE O MESMO HTML. O "mesmo" não é economia de trabalho: duas
// renderizações do board divergiriam com o tempo, e a local passaria a mentir de outro
// jeito — que é pior que atrasar, porque ninguém sabe qual das duas está certa.
//
// A ESTRATÉGIA É INCREMENTAL, com FALLBACK para varredura completa:
//
//	incremental  pergunta "o que mudou desde X?" — o custo acompanha a MUDANÇA
//	completa     relê tudo — o custo acompanha o TEMPO
//
// O incremental existe porque o rate limit é real: 5000 chamadas/hora compartilhadas
// entre todos os agentes, e uma varredura completa custa ~15. Com refresh de 30s seriam
// 1800/hora só do board — e este projeto já estourou o limite hoje, com seis agentes.
//
// O FALLBACK existe porque o incremental pode não estar disponível: sem permissão de
// leitura de eventos, API recusando, ou a primeira carga (que não tem "desde quando").
// Sem ele o comando falharia onde a varredura funcionaria — e um board que não sobe é
// pior que um board com alguns segundos de atraso.
func newBoardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "board",
		Short: "O board do projeto — publicado pelo pipeline, ou servido ao vivo",
	}
	cmd.AddCommand(newBoardServeCmd())
	return cmd
}

func newBoardServeCmd() *cobra.Command {
	var porta int
	var intervalo time.Duration
	var repo string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Sobe o board local, com o estado VIVO em vez da foto do pipeline",
		Long: `Sobe o board num servidor local, lendo o estado atual do GitHub.

O board publicado é uma FOTO: o pipeline roda, gera o board.json, publica. Entre uma
execução e outra o que se vê é passado — e num projeto com muitos agentes entregando,
minutos bastam para a foto mentir sobre quem tem qual card.

  anchors board serve                 # http://localhost:7777
  anchors board serve --port 8080
  anchors board serve --interval 15s  # o piso entre duas leituras

COMO ELE SE MANTÉM FRESCO. A leitura é INCREMENTAL: ele pergunta ao GitHub o que mudou
desde a última vez, e o custo acompanha a MUDANÇA em vez do tempo. Quando o incremental
não está disponível — na primeira carga, sem permissão, ou com a API recusando — ele
RECUA para a varredura completa (fallback), que relê tudo.

A diferença importa pelo rate limit: são 5000 chamadas/hora compartilhadas entre todos
os agentes do projeto, e uma varredura completa custa ~15. Recarregar a cada 30s sem
incremental seriam 1800/hora só do board.

O --interval é o PISO entre duas leituras, não a frequência: se nada mudou, ele não
relê. Serve para o board não sobrecarregar a API quando muita coisa acontece de uma vez.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if repo == "" {
				r, err := repoAtual()
				if err != nil {
					return fmt.Errorf("descobrir o repositório: %w\n\n"+
						"declare com --repo <dono/nome> se este diretório não for um clone", err)
				}
				repo = r
			}

			html, err := initx.BoardHTML()
			if err != nil {
				return fmt.Errorf("ler o HTML do board: %w", err)
			}

			f := &boardSource{repo: repo, piso: intervalo}

			mux := http.NewServeMux()
			mux.HandleFunc("/board.json", func(w http.ResponseWriter, r *http.Request) {
				dados, err := f.leia()
				if err != nil {
					// O ERRO VAI PARA O CLIENTE, e não para um log que ninguém lê. Um
					// board que mostra dado velho sem dizer que falhou é a foto de novo,
					// agora sem ninguém saber.
					http.Error(w, fmt.Sprintf(`{"erro":%q}`, err.Error()), http.StatusBadGateway)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-store")
				w.Write(dados)
			})
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(html))
			})

			addr := fmt.Sprintf("localhost:%d", porta)
			fmt.Printf("board ao vivo em http://%s\n", addr)
			fmt.Printf("  repositório: %s\n", repo)
			fmt.Printf("  piso entre leituras: %s\n\n", intervalo)
			fmt.Println("Ctrl-C para parar.")
			return http.ListenAndServe(addr, mux)
		},
	}

	cmd.Flags().IntVar(&porta, "port", 7777, "porta do servidor local")
	cmd.Flags().DurationVar(&intervalo, "interval", 15*time.Second,
		"piso entre duas leituras do GitHub (não é a frequência: sem mudança, não relê)")
	cmd.Flags().StringVar(&repo, "repo", "", "dono/nome (padrão: o do clone atual)")
	return cmd
}

// boardSource guarda a última leitura e decide quando reler.
//
// O estado aqui não é cache por otimização — é o que torna o incremental possível: sem
// guardar "desde quando", não há o que perguntar ao GitHub.
type boardSource struct {
	repo string
	piso time.Duration

	mu      sync.Mutex
	ultimo  []byte
	lidoEm  time.Time
	desdeAt string // o `updatedAt` mais recente que já vimos
}

func (f *boardSource) leia() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// O PISO. Sem ele, cada aba aberta do board vira uma leitura — e o rate limit é
	// compartilhado com os agentes que estão trabalhando.
	if f.ultimo != nil && time.Since(f.lidoEm) < f.piso {
		return f.ultimo, nil
	}

	// INCREMENTAL primeiro: só vale a pena se já temos uma base e um "desde quando".
	if f.ultimo != nil && f.desdeAt != "" {
		mudou, err := f.algoMudou()
		if err == nil && !mudou {
			// Nada mudou: devolve o que temos e nem relê. É aqui que o custo deixa de
			// acompanhar o tempo.
			f.lidoEm = time.Now()
			return f.ultimo, nil
		}
		// Erro no incremental NÃO é erro do comando: recua para a varredura completa.
		// É o fallback, e o silêncio aqui é deliberado — quem usa o board quer o board,
		// não um diagnóstico da estratégia de leitura.
	}

	dados, err := coletaCompleta(f.repo)
	if err != nil {
		if f.ultimo != nil {
			// Falhou, mas temos leitura anterior. Devolver o velho calado seria a foto
			// de novo; devolver erro derrubaria um board que estava funcionando. O meio
			// é devolver o velho E dizer que é velho — o HTML mostra o `generated`.
			return f.ultimo, nil
		}
		return nil, err
	}

	f.ultimo = dados
	f.lidoEm = time.Now()
	f.desdeAt = maisRecente(dados)
	return dados, nil
}

// algoMudou pergunta ao GitHub se há issue tocada depois do que já vimos.
//
// Uma chamada, um campo. É a diferença entre custo por MUDANÇA e custo por TEMPO.
//
// PELO REST, COM ORDENAÇÃO EXPLÍCITA, e os dois detalhes são o conserto de um defeito
// medido: o `gh issue list --limit 1` devolve a issue de maior NÚMERO, não a tocada por
// último. Onze cards foram fechados e o board não releu, porque o mais recente por
// número (`#683`, 19:22) era mais velho que o tocado de fato (`#630`, 19:29).
//
// O REST também escapa do limite SECUNDÁRIO do GraphQL, que derruba o `gh issue list`
// enquanto `gh api rate_limit` ainda reporta a cota cheia — medido neste projeto com
// nove agentes, e o motivo de o board ter ficado sem dado antes.
func (f *boardSource) algoMudou() (bool, error) {
	out, err := exec.Command("gh", "api",
		"/repos/"+f.repo+"/issues?state=all&labels=anchors&sort=updated&direction=desc&per_page=1",
		"--jq", ".[0].updated_at // \"\"").Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) > f.desdeAt, nil
}

// maisRecente extrai o `updated` mais novo do board, que vira o "desde quando" seguinte.
func maisRecente(dados []byte) string {
	var b struct {
		Items []struct {
			Updated string `json:"updated"`
		} `json:"items"`
	}
	if json.Unmarshal(dados, &b) != nil {
		return ""
	}
	maior := ""
	for _, i := range b.Items {
		if i.Updated > maior {
			maior = i.Updated
		}
	}
	return maior
}

// repoAtual descobre `dono/nome` pelo clone, para o comando funcionar sem flag.
func repoAtual() (string, error) {
	out, err := exec.Command("gh", "repo", "view", "--json", "nameWithOwner",
		"--jq", ".nameWithOwner").Output()
	if err != nil {
		return "", err
	}
	r := strings.TrimSpace(string(out))
	if r == "" {
		return "", fmt.Errorf("o `gh` não devolveu o repositório")
	}
	return r, nil
}

// coletaCompleta monta o board inteiro — a varredura que o fallback usa.
//
// Ela chama o MESMO `gh issue list` do pipeline, com o mesmo `--jq`: o contrato do
// `board.json` é único, e duas implementações dele divergiriam no primeiro campo novo.
func coletaCompleta(repo string) ([]byte, error) {
	jq, err := initx.BoardCollectJQ()
	if err != nil {
		return nil, err
	}

	// PELO REST, e não pelo `gh issue list`.
	//
	// O porcelain vai por GraphQL, cujo limite SECUNDÁRIO derruba a consulta enquanto o
	// `gh api rate_limit` ainda reporta cota cheia. Medido neste projeto com nove agentes:
	// o board ficou sem dado nenhum, servindo só o erro, enquanto o REST respondia
	// normalmente a tudo.
	//
	// O REST nomeia os campos em snake_case (`updated_at`, `created_at`, `closed_at`) e
	// não traz `comments` embutido — o jq do pipeline espera camelCase e a lista de
	// comentários. A ponte abaixo renomeia e injeta `comments: []` para o contrato
	// continuar o mesmo.
	//
	// OS COMENTÁRIOS DOS ABERTOS VÊM POR GRAPHQL, em lote.
	//
	// Isto injetava `comments: []` fixo, e o custo era o `owner`/`ownership` — que saem do
	// comentário `anchors-owner:`. Sem eles o board não mostrava de quem é o card, e a
	// modal (que também lê comentário) abria vazia.
	//
	// A razão original era boa: buscar comentário por REST é uma chamada POR CARD, e com
	// centenas de cards isso é lento e bate no limite. Mas o GraphQL traz issue COM
	// comentários numa chamada só — medido: 3 issues com os últimos comentários de cada
	// numa requisição.
	//
	// SÓ OS ABERTOS, e é o que torna o custo aceitável: o `owner` de um card fechado não
	// interessa a ninguém, e os abertos são uma fração do total (medido: 87 de 421).
	//
	// FALHA AQUI NÃO DERRUBA O BOARD. Se o GraphQL recusar — é o limite secundário que
	// motivou o REST —, os comentários ficam vazios e o board serve o resto. Um board sem
	// o dono continua sendo melhor que um board sem nada.
	// O MAPA VAI POR ARQUIVO, e não embutido na expressão.
	//
	// A primeira versão o interpolava no jq, e o `exec` recusou: `argument list too long`.
	// Com 87 cards e trinta comentários cada, o literal passa do limite do sistema — e o
	// erro aparecia como board VAZIO, não como falha de coleta.
	//
	// `--slurpfile` lê o arquivo e o expõe como variável. O `[0]` porque ele sempre entrega
	// um array dos documentos do arquivo, e o nosso é um objeto só.
	arqCom, limpaCom := comentariosEmArquivo(repo)
	defer limpaCom()
	ponte := `[.[] | select(.pull_request == null) | {
	    number, title, url, body, labels, author: .user,
	    updatedAt: .updated_at, createdAt: .created_at, closedAt: .closed_at,
	    comments: ($coment[0][.number|tostring] // [])
	  }] | ` + jq

	// O `--paginate` sem `--jq` devolve os arrays de cada página concatenados
	// (`[...][...]`), que não é JSON válido. O `--slurp` os juntaria, mas o `gh` recusa
	// combiná-lo com `--jq`. A saída é pedir o cru e costurar aqui.
	cru, err := exec.Command("gh", "api", "--paginate",
		"/repos/"+repo+"/issues?state=all&labels=anchors&per_page=100").Output()
	if err == nil {
		//  achata o  que a costura das páginas produz.
		argvJq := []string{"-c"}
		if arqCom != "" {
			argvJq = append(argvJq, "--slurpfile", "coment", arqCom)
		} else {
			// SEM O ARQUIVO o jq não conheceria `$coment`, e a expressão inteira falharia
			// — levando o board a zero itens por falta do DONO, que é o menos importante
			// do que ele mostra.
			argvJq = append(argvJq, "--argjson", "coment", "[{}]")
		}
		argvJq = append(argvJq, "add | "+ponte)
		filtro := exec.Command("jq", argvJq...)
		filtro.Stdin = strings.NewReader("[" + strings.ReplaceAll(string(cru), "][", ",") + "]")
		var saida []byte
		saida, err = filtro.Output()
		if err == nil {
			// `[[...]]` — a costura aninha um nível a mais; o jq já devolve o array final.
			out := saida
			items := strings.TrimSpace(string(out))
			if items == "" {
				items = "[]"
			}
			return []byte(fmt.Sprintf(`{"generated":%q,"items":%s}`,
				time.Now().UTC().Format(time.RFC3339), items)), nil
		}
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("filtrar o board: %w\n  %s", err, strings.TrimSpace(string(ee.Stderr)))
		}
	}
	var out []byte
	if err != nil {
		var stderr string
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = strings.TrimSpace(string(ee.Stderr))
		}
		return nil, fmt.Errorf("coletar do GitHub: %w%s", err, func() string {
			if stderr != "" {
				return "\n  " + stderr
			}
			return ""
		}())
	}

	items := strings.TrimSpace(string(out))
	if items == "" {
		items = "[]"
	}
	board := fmt.Sprintf(`{"generated":%q,"items":%s}`,
		time.Now().UTC().Format(time.RFC3339), items)
	return []byte(board), nil
}

var _ = os.Getenv // mantém o import quando o corpo muda

// comentariosDosAbertos devolve, como literal jq, um mapa de número → comentários.
//
// O `owner` e o `ownership` do board saem do comentário `anchors-owner:`, e a modal do card
// mostra a conversa. Sem eles o board local ficava mudo sobre quem está com o quê.
//
// GRAPHQL PORQUE O REST PEDIRIA UMA CHAMADA POR CARD. A consulta abaixo traz cem issues com
// os últimos comentários de cada, e pagina pelo cursor.
//
// `last: 30` e não `first`: o que interessa é o comentário MAIS RECENTE — o `anchors-owner`
// atual, não o primeiro que o card recebeu. Um card com muita conversa teria o dono no fim.
//
// DEVOLVE `{}` em qualquer erro, e isso é deliberado: o limite secundário do GraphQL é o
// que motivou a coleta por REST, e um board sem dono é melhor que um board sem nada.
// comentariosEmArquivo grava o mapa num temporário e devolve o caminho.
//
// O jq recebe por `--slurpfile` porque o literal não cabe na linha de comando: medido com
// 87 cards, o `exec` recusou com `argument list too long` — e o sintoma era o board vazio.
//
// Devolve caminho vazio quando não há o que gravar; quem chama trata passando `[{}]`.
func comentariosEmArquivo(repo string) (string, func()) {
	m := comentariosDosAbertos(repo)
	if m == "" || m == "{}" {
		return "", func() {}
	}
	f, err := os.CreateTemp("", "anchors-board-coment-*.json")
	if err != nil {
		return "", func() {}
	}
	if _, err := f.WriteString(m); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", func() {}
	}
	f.Close()
	return f.Name(), func() { os.Remove(f.Name()) }
}

func comentariosDosAbertos(repo string) string {
	partes := strings.SplitN(repo, "/", 2)
	if len(partes) != 2 {
		return "{}"
	}
	q := `query($owner:String!,$name:String!,$cursor:String){
	  repository(owner:$owner,name:$name){
	    issues(first:100,states:OPEN,labels:["anchors"],after:$cursor){
	      pageInfo{hasNextPage endCursor}
	      nodes{number comments(last:30){nodes{body}}}
	    }
	  }
	}`
	acc := map[string][]map[string]string{}
	cursor := ""
	for i := 0; i < 10; i++ { // teto: 1000 cards abertos é muito além do real
		argv := []string{"api", "graphql", "-f", "query=" + q,
			"-F", "owner=" + partes[0], "-F", "name=" + partes[1]}
		if cursor != "" {
			argv = append(argv, "-F", "cursor="+cursor)
		}
		out, err := exec.Command("gh", argv...).Output()
		if err != nil {
			// PARCIAL VALE MAIS QUE NADA: se a segunda página falhar, o board mostra o
			// dono dos cards da primeira em vez de nenhum.
			break
		}
		var resp struct {
			Data struct {
				Repository struct {
					Issues struct {
						PageInfo struct {
							HasNextPage bool   `json:"hasNextPage"`
							EndCursor   string `json:"endCursor"`
						} `json:"pageInfo"`
						Nodes []struct {
							Number   int `json:"number"`
							Comments struct {
								Nodes []struct {
									Body string `json:"body"`
								} `json:"nodes"`
							} `json:"comments"`
						} `json:"nodes"`
					} `json:"issues"`
				} `json:"repository"`
			} `json:"data"`
		}
		if json.Unmarshal(out, &resp) != nil {
			break
		}
		for _, n := range resp.Data.Repository.Issues.Nodes {
			lista := make([]map[string]string, 0, len(n.Comments.Nodes))
			for _, c := range n.Comments.Nodes {
				lista = append(lista, map[string]string{"body": c.Body})
			}
			acc[fmt.Sprintf("%d", n.Number)] = lista
		}
		if !resp.Data.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		cursor = resp.Data.Repository.Issues.PageInfo.EndCursor
	}
	if len(acc) == 0 {
		return "{}"
	}
	b, err := json.Marshal(acc)
	if err != nil {
		return "{}"
	}
	return string(b)
}
