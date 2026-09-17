package main

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// --- quais arquivos são DERIVADOS, e por isso se reconstroem em vez de se mesclar ---
//
// Quem resolve conflito precisa distinguir dois casos que parecem iguais no `git status`:
//
//	gerado     o conflito é RUÍDO — reconstruir produz a resposta certa
//	de trabalho  o conflito é DIVERGÊNCIA — duas pessoas escreveram coisas diferentes
//
// A lista dos gerados estava escrita À MÃO no script de um projeto:
//
//	grep -vE '^anchors\.graph\.yaml$|^docs/|-progress\.md$'
//
// Três caminhos literais, e os três já existem no produto — `mapx.DefaultPath`,
// `doct.OutDir` e `scan.SufixoProgresso`. Duplicá-los em texto quebra o agnosticismo: o
// `docs/` é configurável, e um projeto que o mudasse veria o script tratar documentação
// gerada como divergência de conteúdo.
//
// Este comando é a resposta do PRODUTO à pergunta. Quem resolve pergunta; ninguém escreve
// a lista de novo.

func newGeneratedPathsCmd() *cobra.Command {
	var root, formato string
	cmd := &cobra.Command{
		Use:   "generated-paths",
		Short: "Os caminhos dos arquivos DERIVADOS — quem conflita neles reconstrói",
		Long: `Imprime os padrões dos arquivos que o Anchors gera.

Quem resolve conflito usa isto para distinguir RUÍDO de DIVERGÊNCIA: num arquivo
gerado, reconstruir produz a resposta certa; num de trabalho, duas pessoas
escreveram coisas diferentes e só elas decidem.

    anchors generated-paths              # um por linha
    anchors generated-paths --format re  # como alternância de regex, para o grep

A lista sai do PRODUTO, não de constantes repetidas no script de cada projeto.
É uma fonte só: quando um caminho mudar, todos os projetos acompanham.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			// A CONFIG é carregada para o comando falhar cedo num projeto sem ela — a
			// pergunta "o que é gerado?" só faz sentido num projeto regido.
			if _, err := config.Load(absRoot + "/" + config.DefaultFile); err != nil {
				return err
			}

			// O MAPA é o gerado por excelência: derivado da árvore, com merge driver
			// próprio (`anchors map merge`).
			padroes := []string{"^" + regexEscape(mapx.DefaultPath) + "$"}

			// A DOCUMENTAÇÃO COMPILADA. Hoje o diretório é a constante `doct.OutDir` — não
			// é declarável, e por isso o ganho aqui não é "ler a config": é a lista sair de
			// UM lugar. Quando o diretório virar declarável, muda-se aqui e todos os
			// projetos acompanham; hoje cada script que o repete tem de ser caçado.
			padroes = append(padroes, "^"+regexEscape(doct.OutDir)+"/")

			// O PROGRESSO de cada plano. É sufixo e não diretório: o companheiro vive ao
			// lado do plano, onde quer que ele esteja.
			padroes = append(padroes, regexEscape(scan.SufixoProgresso)+"$")

			if formato == "re" {
				fmt.Fprintln(cmd.OutOrStdout(), strings.Join(padroes, "|"))
				return nil
			}
			for _, p := range padroes {
				fmt.Fprintln(cmd.OutOrStdout(), p)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&formato, "format", "lines",
		"`lines` (um por linha) ou `re` (alternância, para o grep -E)")
	return cmd
}

// regexEscape protege os metacaracteres de um caminho usado dentro de um regex.
//
// `regexp.QuoteMeta` escaparia a barra também em alguns dialetos; aqui só o ponto importa,
// e é o que aparece em `anchors.graph.yaml` e `-progress.md`. Escapar de menos deixaria o
// `.` casar qualquer caractere — `anchorsXgraph.yaml` passaria por gerado.
func regexEscape(s string) string {
	return strings.ReplaceAll(s, ".", `\.`)
}
