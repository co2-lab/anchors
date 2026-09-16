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
func (f *boardSource) algoMudou() (bool, error) {
	out, err := exec.Command("gh", "issue", "list",
		"--repo", f.repo, "--state", "all", "--label", "anchors",
		"--limit", "1", "--json", "updatedAt",
		"--jq", ".[0].updatedAt // \"\"").Output()
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

	out, err := exec.Command("gh", "issue", "list",
		"--repo", repo, "--state", "all", "--label", "anchors", "--limit", "500",
		"--json", "number,title,url,labels,comments,updatedAt,body,createdAt,closedAt,author",
		"--jq", jq).Output()
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
