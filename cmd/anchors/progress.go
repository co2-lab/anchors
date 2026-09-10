package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

// --- o PROGRESSO mora fora do plano ---
//
// Um plano é DECISÃO: o que vai ser feito, em que ordem, e por quê. Alterá-lo tem de
// significar que a decisão mudou — é sobre isso que o `plano-alterado-justificado` cobra
// uma revisão (`{CODIGO}-R0001`), e é a única defesa contra o projeto derivar em silêncio.
//
// Enquanto os checkboxes de fase viviam no plano, marcar `- [x]` era ALTERAR o plano. E
// aí três coisas quebravam de uma vez:
//
//  1. o gate não distinguia "mudei a direção" de "terminei uma fase", e cobrava revisão
//     das duas. Uma revisão que diz "concluí a fase 1" é ruído — e ruído em gate
//     bloqueante é o que faz alguém desligá-lo.
//
//  2. o carimbo de julgamento da aresta `plano → spec` guarda a rev das DUAS pontas.
//     Marcar a fase mudava a rev do plano, e o julgamento da SPEC caía — um julgamento
//     sobre spec e código, derrubado por um checkbox. O efeito é circular: concluir a
//     fase invalida a verificação da spec que a fase entregou.
//
//  3. medido no blue-eyes: `plans/0017-mutacao.md` tinha DOIS commits — o que o criou (83
//     linhas) e um que mudou 1 linha, `- [ ]` para `- [x]`. Cem por cento das alterações
//     pós-criação eram progresso.
//
// O arquivo de progresso fica FORA DO MAPA de propósito. Se entrasse como nó com camada,
// os gates voltariam a confrontá-lo e o problema renasceria com outro nome — inclusive o
// `plano-alterado-justificado` cobrando revisão de um arquivo cuja única função é mudar.
// Ele é estado, não decisão: ninguém precisa justificar por que o estado avançou.

// progressSuffix liga o arquivo de estado ao plano: mesmo nome, sufixo fixo.
//
// A definição CANÔNICA é a do `scan` — é ele que precisa manter o arquivo fora do mapa,
// e uma segunda constante aqui poderia divergir dela em silêncio. Esta é a mesma string,
// e o teste `TestProgresso_sufixoBateComOScan` confronta as duas.
//
// A colocação ao lado (em vez de uma pasta `progress/`) é o que faz os dois serem lidos
// juntos: quem abre o plano vê o companheiro na mesma listagem, e um plano cujo progresso
// não existe fica visível pela ausência.
const progressSuffix = "-progress.md"

// progressPath devolve o arquivo de progresso de um plano.
func progressPath(plano string) string {
	ext := filepath.Ext(plano)
	return strings.TrimSuffix(plano, ext) + progressSuffix
}

// fasesDoPlano lê os códigos de fase declarados nos cabeçalhos do plano.
//
// A fonte é o CABEÇALHO (`### PLTFR-F01 — ...`), a mesma que os gates `fase-existe` e
// `fase-ordenada` já usam. Ler daqui em vez de manter uma segunda lista é o que garante
// que o progresso fale das fases que existem: uma fase renomeada aparece, uma inventada
// não.
//
// O TAMANHO DO CÓDIGO vem de `config.CodeLengthPattern()`, e não de um `{4,5}` escrito à
// mão: `code_lengths` é configurável por projeto. Um literal aqui daria a paridade que
// este comentário afirma apenas para os projetos no default — nos outros, a fase seria
// reconhecida pelos gates e ignorada por este comando, e o progresso nasceria vazio sem
// nada acusar.
//
// Construído por CHAMADA, não em `var`: a config é carregada depois da inicialização do
// pacote, e um regex montado no init congelaria o default.
func phaseInHeaderRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^#{2,4}[^\S\n]+([A-Z0-9]` +
		config.CodeLengthPattern() + `-F\d{2})\b[^\S\n]*—?[^\S\n]*(.*)$`)
}

type planPhase struct {
	Codigo string
	Titulo string
}

func planPhases(conteudo string) []planPhase {
	var out []planPhase
	for _, m := range phaseInHeaderRE().FindAllStringSubmatch(conteudo, -1) {
		out = append(out, planPhase{Codigo: m[1], Titulo: strings.TrimSpace(m[2])})
	}
	return out
}

// writeInitialProgress cria o `-progress.md` de um plano, com uma linha por fase.
//
// Não sobrescreve: o arquivo guarda o estado do trabalho, e regravá-lo apagaria o que já
// foi registrado. Um plano que ganha fase nova tem a linha acrescentada à mão — o comando
// não reescreve estado que não é dele.
func writeInitialProgress(planoPath, conteudoPlano, codigo string) (string, error) {
	destino := progressPath(planoPath)
	if _, err := os.Stat(destino); err == nil {
		return "", fmt.Errorf("%s já existe — o progresso é estado, e não sobrescrevo", destino)
	}

	var b strings.Builder
	b.WriteString("<!-- anchors:progress — o ESTADO do plano `" +
		filepath.Base(planoPath) + "`.\n\n")
	b.WriteString("Este arquivo fica FORA DO MAPA de propósito: ele existe para MUDAR, e um\n")
	b.WriteString("arquivo que muda por natureza não pode ser confrontado pelos gates que\n")
	b.WriteString("cobram justificativa de mudança. O plano ao lado é a DECISÃO; alterá-lo tem\n")
	b.WriteString("de significar que a decisão mudou.\n\n")
	b.WriteString("Marque `[x]` aqui, nunca no plano.\n-->\n\n")
	b.WriteString("# Progresso — " + codigo + "\n\n")

	fases := planPhases(conteudoPlano)
	if len(fases) == 0 {
		b.WriteString("TODO: o plano ainda não declara fases. Quando declarar, acrescente uma\n")
		b.WriteString("seção por fase aqui, com um item por spec semeada.\n")
	}
	for _, f := range fases {
		b.WriteString("## " + f.Codigo)
		if f.Titulo != "" {
			b.WriteString(" — " + f.Titulo)
		}
		b.WriteString("\n\n- [ ] TODO: um item por spec que esta fase semeia\n\n")
	}

	if err := os.WriteFile(destino, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return destino, nil
}

// newProgressCmd cria o `-progress.md` de um plano que JÁ EXISTE.
//
// O `anchors new plan` cria o companheiro junto com o plano — mas só ele. Um projeto que
// adotou o Anchors antes deste mecanismo tem todos os planos sem companheiro, para
// sempre, e nada acusa.
//
// Medido no blue-eyes: 17 planos, ZERO com `-progress.md`, e 17 com checkbox dentro do
// plano — que é exatamente o que este arquivo existe para tirar de lá. O mecanismo
// existia, estava testado, e não alcançava um único plano do projeto.
//
// O caminho feliz (plano novo) funcionava; era a ADOÇÃO que não tinha caminho.
func newProgressCmd() *cobra.Command {
	var root, para string
	cmd := &cobra.Command{
		Use:   "progress --for <plano>",
		Short: "Cria o `-progress.md` de um plano existente (o ESTADO, fora do mapa)",
		Long: `Cria o companheiro de progresso de um plano que já existe.

Um plano é DECISÃO; o progresso é ESTADO. Enquanto os checkboxes de fase viviam no
plano, marcar ` + "`- [x]`" + ` era ALTERAR o plano — e isso cobrava revisão de quem só
terminou uma fase, além de derrubar o julgamento da spec que a fase entregou (o carimbo
da aresta guarda a rev das duas pontas).

O ` + "`anchors new plan`" + ` já cria o companheiro. Este comando é para os planos que
nasceram antes:

  anchors new progress --for plans/0002-plataforma.md

Não sobrescreve: o arquivo guarda estado, e regravá-lo apagaria o que já foi registrado.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if para == "" {
				return fmt.Errorf("informe o plano com --for (ex: `anchors new progress --for plans/0002-x.md`)")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			planoPath := para
			if !filepath.IsAbs(planoPath) {
				planoPath = filepath.Join(absRoot, para)
			}
			conteudo, err := os.ReadFile(planoPath)
			if err != nil {
				return fmt.Errorf("ler o plano: %w", err)
			}
			// O CÓDIGO vem do header do próprio plano, e não de um argumento: pedi-lo
			// abriria a porta para o progresso nascer com identidade divergente da do
			// plano que ele acompanha — e o par deixaria de ser localizável.
			codigo := codeDoHeaderSpec(string(conteudo))
			if codigo == "" {
				return fmt.Errorf("o plano %s não declara `code:` no header @anchors — "+
					"sem ele o progresso nasceria sem identidade", para)
			}
			prog, err := writeInitialProgress(planoPath, string(conteudo), codigo)
			if err != nil {
				return err
			}
			fmt.Printf("✓ %s criado (o ESTADO do plano %s; marque `[x]` aqui, nunca no plano)\n",
				relTo(absRoot, prog), codigo)
			if fases := planPhases(string(conteudo)); len(fases) > 0 {
				fmt.Printf("  %d fase(s) do plano, uma seção para cada\n", len(fases))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&para, "for", "", "o plano cujo progresso será criado")
	return cmd
}
