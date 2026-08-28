package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newImpactCmd() *cobra.Command {
	var root, mapPath string
	cmd := &cobra.Command{
		Use:   "impact <arquivo>",
		Short: "Análise de impacto: o que uma alteração neste arquivo atinge",
		Long: `Consulta o mapa e mostra, para um arquivo alterado, DUAS dimensões:

  • Propaga para  — os filhos que dependem dele e precisam ser refeitos (a onda,
    descendo; para em nós @noPropagation).
  • Valida contra — os pais contra os quais ele deve ser confrontado (subindo;
    uma divergência aqui vira issue, não propagação).

É só consulta — não abre issue nem altera nada.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}

			target := relTo(absRoot, args[0])
			if !nodeExists(g, target) {
				return fmt.Errorf("arquivo %q não está no mapa", target)
			}

			imp := g.AnalyzeImpact(target)
			printImpact(imp)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa (default: <root>/anchors.graph.yaml)")
	return cmd
}

// relTo converte o arg (que pode ser absoluto ou relativo ao cwd) para o caminho
// relativo à raiz do projeto — a forma como os nós são identificados no mapa.
// relTo converte o argumento do usuário num caminho relativo à RAIZ do projeto.
//
// O argumento pode vir de duas origens legítimas, e elas se resolvem de formas opostas:
// relativo ao diretório onde o comando foi chamado (como todo shell faz) ou relativo à
// raiz do projeto (que é a forma que os próprios prompts do Anchors imprimem, e a que os
// agentes copiam). Enquanto a raiz era sempre o CWD, as duas coincidiam. Desde que os
// comandos passaram a SUBIR até a raiz, deixaram de coincidir — e um `--for` copiado do
// prompt, executado de um subdiretório, virava `packages/backend/packages/backend/…`.
//
// A régua: se o caminho relativo à raiz EXISTE, é ele. É o desempate certo porque a origem
// do argumento não é declarada — só o disco sabe qual leitura corresponde a um arquivo.
//
// A saída é SEMPRE com barra normal, porque o que se devolve é o id de um nó, e o id é o
// mesmo em toda máquina: o mapa é versionado e trafega entre macOS, Linux e Windows. Sem
// isto, `filepath.Clean`/`Rel` devolvem "packages\backend\x.spec.md" no Windows, a busca
// no mapa (gravado com "/") não casa, e o arquivo aparece como "REGIDO mas fora do mapa" —
// pedindo um `map build` que reescreveria o mapa inteiro no dialeto da máquina. Converter
// só na fronteira do disco é o que mantém o mapa intercambiável.
func relTo(root, arg string) string {
	if !filepath.IsAbs(arg) {
		if _, err := os.Stat(filepath.Join(root, arg)); err == nil {
			return filepath.ToSlash(filepath.Clean(arg))
		}
	}
	abs, err := filepath.Abs(arg)
	if err != nil {
		return filepath.ToSlash(arg)
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return filepath.ToSlash(arg)
	}
	return filepath.ToSlash(rel)
}

func nodeExists(g *mapx.Graph, id string) bool {
	for _, n := range g.Nodes {
		if n.ID == id {
			return true
		}
	}
	return false
}

func printImpact(imp mapx.Impact) {
	fmt.Printf("impacto de: %s\n\n", imp.Origin)

	fmt.Printf("↓ Propaga para (%d) — dependem desta mudança, refazer:\n", len(imp.Propagate))
	if len(imp.Propagate) == 0 {
		fmt.Println("   (nada — nenhum filho depende dele)")
	}
	for _, n := range imp.Propagate {
		fmt.Printf("   %s\n", n)
	}

	fmt.Printf("\n↑ Valida contra (%d) — confrontar; divergência vira issue:\n", len(imp.Validate))
	if len(imp.Validate) == 0 {
		fmt.Println("   (nada — não é regido por ninguém)")
	}
	for _, n := range imp.Validate {
		fmt.Printf("   %s\n", n)
	}
}
