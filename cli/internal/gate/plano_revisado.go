package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- o plano REVISADO avisa quem o lê ---
//
// Planejar erra, e o erro aparece quando o trabalho começa. Editar o plano antigo perderia
// o registro: um plano é ÂNCORA, e editá-lo depois de implementado faz o documento
// descrever algo que não foi o que aconteceu.
//
// Por isso a revisão é um plano NOVO com `revises:`. O custo dessa escolha é a leitura
// fora de ordem — quem abre o plano antigo não sabe que existe um mais novo, e segue uma
// decisão que foi revista. É o defeito que este gate impede.
//
// A régua é simples: se algum plano declara `revises: X`, então X tem de dizer isso no
// topo. O aviso não é cosmético — é a única coisa que alcança quem chegou pelo caminho
// errado.

// O VOCABULÁRIO já existe, e o gate usa o que existe.
//
// A primeira versão casava prosa em português ("revisado por"), e um plano em inglês não
// casaria — o gate reprovaria um aviso que EXISTE. A segunda inventou `@revised-by`,
// e inventar marcador tem outro custo: mais uma convenção para alguém decorar.
//
// O que já está no vocabulário:
//
//   - o ALERTA do Markdown (`> [!IMPORTANT]`, `> [!WARNING]`) — sintaxe que o GitHub
//     renderiza, independente de idioma, e que quem escreve markdown já conhece;
//   - o CÓDIGO do plano e da fase (`FNDTN`, `FNDTN-F01`) — a identidade que o Anchors já
//     usa em todo lugar, e que não muda quando alguém reescreve o título.
//
// O aviso de topo é um `[!IMPORTANT]` citando o código do plano que revisa; a marcação de
// parte é um `[!WARNING]` citando o mesmo código. A diferença entre os dois é o nível do
// alerta, e ela é semântica: o de topo muda como se lê o documento inteiro, o de parte
// muda um trecho.
// As DUAS formas são aceitas, e cada uma tem sua razão:
//
//   - `> [!IMPORTANT]` com o código — renderiza no GitHub como caixa de destaque, e quem
//     escreve markdown já conhece a sintaxe. É o que salta aos olhos de quem abre o
//     arquivo na web;
//   - `@revised-by: <código ou caminho>` — explícito e buscável com grep, na mesma
//     família do `@no-mark` que a doutrina já usa. É o que se encontra numa varredura.
//
// Aceitar as duas não é indecisão: elas servem a leitores diferentes (o humano na web e a
// ferramenta na linha de comando), e exigir uma delas obrigaria metade dos projetos a
// escrever no formato que não é o seu.
func avisoDeTopoRE() *regexp.Regexp {
	return regexp.MustCompile(`(?i)(>\s*\[!IMPORTANT\][\s\S]{0,400}?[A-Z0-9]` +
		config.CodeLengthPattern() + `\b|@revised-by[^\S\n]*:?[^\S\n]*\S)`)
}

// alteradoPorRE casa a marcação de UMA PARTE atingida pela revisão.
//
// Distinta do aviso de topo de propósito: aquele diz que o plano mudou, este diz ONDE. Um
// plano revisado sem nenhuma marcação de parte obriga quem lê a reler tudo procurando o
// que mudou — o mesmo custo de não haver aviso.
func alteradoPorRE() *regexp.Regexp {
	return regexp.MustCompile(`(?i)(>\s*\[!WARNING\][\s\S]{0,400}?[A-Z0-9]` +
		config.CodeLengthPattern() + `\b|@amended-by[^\S\n]*:?[^\S\n]*\S)`)
}

// checkPlanoRevisado confronta os dois lados do `revises`.
func checkPlanoRevisado(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, "só plano revisa plano"
	}
	if g == nil {
		return Pending, "sem mapa não dá para achar quem revisa quem"
	}

	// Quem revisa ESTE plano?
	var revisores []string
	for _, o := range g.Nodes {
		if o.Kind != mapx.KindPlan || o.ID == n.ID {
			continue
		}
		for _, alvo := range o.Revises {
			if alvo == n.ID {
				revisores = append(revisores, o.ID)
			}
		}
	}
	if len(revisores) == 0 {
		// Se ESTE declara supersedes, o alvo precisa existir: apontar para um plano que
		// não está no mapa é uma revisão que não revisa nada.
		var ausentes []string
		for _, alvo := range n.Revises {
			achou := false
			for _, o := range g.Nodes {
				if o.ID == alvo {
					achou = true
				}
			}
			if !achou {
				ausentes = append(ausentes, alvo)
			}
		}
		if len(ausentes) > 0 {
			return Fail, fmt.Sprintf("declara `revises: %s`, e esse plano não está no "+
				"mapa. Uma revisão que aponta para o nada não avisa ninguém — e o plano que "+
				"se pretendia revisar continua sendo lido como vigente",
				strings.Join(ausentes, ", "))
		}
		// QUEM REVISA é lembrado do passo que falta, e ele é no OUTRO arquivo.
		//
		// O gate cobra o aviso do lado do plano revisado, e quem escreve o revisor só
		// descobre isso quando o outro reprova — num arquivo que ele talvez nem tenha
		// aberto. Dizer aqui alcança a pessoa no momento em que ela está fazendo a
		// revisão.
		// Só cobra os que AINDA não avisam: um lembrete que não some vira ruído
		// permanente, e ruído permanente é o que treina a equipe a ignorar o gate.
		var semAviso []string
		for _, alvo := range n.Revises {
			b, err := os.ReadFile(filepath.Join(root, alvo))
			if err != nil || !avisoDeTopoRE().MatchString(topo(string(b))) {
				semAviso = append(semAviso, alvo)
			}
		}
		if len(semAviso) > 0 {
			return Pending, fmt.Sprintf("revisa %s, e o plano revisado precisa AVISAR quem "+
				"o lê. Escreva no topo de %s:\n    > `@revised-by: %s` — <o que mudou e "+
				"por quê>\nSem isso, quem abrir o plano antigo segue uma decisão que foi "+
				"revista — ele continua parecendo coerente, porque É o registro do que se "+
				"decidiu na época",
				strings.Join(semAviso, ", "), semAviso[0], n.ID)
		}
		return Skip, "nenhum plano revisa este"
	}

	// AS PARTES atingidas, e não só o topo.
	//
	// O aviso de topo diz que o plano mudou; ele não diz ONDE. Quem lê a fase 3 não sabe
	// se ela é uma das que mudaram, e o custo de descobrir é reler o plano inteiro
	// procurando — que é o mesmo custo de não ter aviso nenhum.
	//
	// A cobrança é PENDENTE, não falha: nem toda revisão atinge uma parte nomeável (um
	// plano pode ser revisado por inteiro), e transformar isso em reprovação obrigaria a
	// inventar marcação onde ela não cabe.
	if avisoDeTopoRE().MatchString(topo(content)) && !alteradoPorRE().MatchString(content) {
		return Pending, fmt.Sprintf("avisa que foi revisado por %s, e não marca QUAIS partes "+
			"mudaram. Quem lê a fase 3 não sabe se ela é uma delas, e descobrir custa reler "+
			"o plano inteiro. Marque cada parte atingida:\n"+
			"    · já implementada — o texto FICA (descreve o que foi feito), e abaixo:\n"+
			"      > `@amended-by: %s` — <o que muda daqui em diante>\n"+
			"    · ainda não implementada — reescreva, preservando o original:\n"+
			"      > `@amended-by: %s` — **Era:** <texto original>",
			strings.Join(revisores, ", "), revisores[0], revisores[0])
	}

	// Revisado: o aviso é obrigatório, e no TOPO — quem lê de cima para baixo precisa
	// saber antes de seguir a primeira decisão.
	if !avisoDeTopoRE().MatchString(topo(content)) {
		return Fail, fmt.Sprintf("é revisado por %s e não avisa quem o lê. Um plano é "+
			"registro do que se decidiu, e continua valendo como registro — mas quem chega "+
			"nele sem saber da revisão segue uma decisão que foi revista. Escreva no TOPO: "+
			"`> @revised-by: %s — <o que mudou e por quê>`",
			strings.Join(revisores, ", "), revisores[0])
	}
	return Pass, ""
}

// topo devolve o começo do documento — onde o aviso tem de estar para ser visto.
//
// 40 linhas cobre o header, o título e a primeira seção: um aviso depois disso é
// encontrado por quem já leu o plano inteiro, e aí ele não avisou nada.
func topo(content string) string {
	linhas := strings.Split(content, "\n")
	if len(linhas) > 40 {
		linhas = linhas[:40]
	}
	return strings.Join(linhas, "\n")
}
