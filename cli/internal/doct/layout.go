package doct

import "fmt"

// --- o FORMATO é UM, e é decidido em um lugar ---
//
// Uma camada com três unidades cabe numa página com tudo — regras, invariantes e cenários.
// Com trinta, a mesma página é um documento que ninguém rola até o fim.
//
// A decisão é simples; distribuí-la é que não. Ela aparece em três lugares e os três têm
// de concordar:
//
//	a página de camada   mostra tudo, ou resume?
//	o índice de regras   o link vai à regra, ou à unidade que a contém?
//	o de comportamento   idem, para o cenário
//
// A PRIMEIRA VERSÃO deixou cada um decidir por conta, e o resultado foi medido: 483 de 812
// links quebrados. A página resumia, o índice apontava para a âncora da regra, e a âncora
// não existia — todos clicáveis, todos parando no mesmo lugar. Um índice assim é pior que
// não ter índice, porque parece funcionar.
//
// Aqui a decisão é tomada UMA VEZ, no `docs build`, e todo template recebe a mesma
// resposta. Não há como um discordar do outro: não há dois lugares onde discordar.

// Layout é a decisão de formato, resolvida antes de qualquer template rodar.
type Layout struct {
	// MaxUnits é o número de unidades a partir do qual a página de camada resume em vez
	// de trazer o conteúdo inteiro.
	MaxUnits int
	// MaxLines é o mesmo, medido em LINHAS de spec. Existe além do `MaxUnits` porque
	// contar unidades engana: cinco specs longas geram mais página que vinte curtas, e é
	// o comprimento que faz alguém desistir de rolar. Qualquer um dos dois estourando
	// resume — o critério é "esta página ficou grande", não "por qual das duas razões".
	MaxLines int
}

// DefaultLayout é o corte quando o projeto não declara outro.
//
// Vinte unidades ou duas mil linhas. Os números são uma régua grosseira e assumidamente
// imperfeita — é por isso que são configuráveis. O que eles não podem ser é implícitos:
// um limiar escondido no código faria o formato da documentação mudar sozinho no dia em
// que a vigésima primeira unidade entrasse, sem ninguém ter decidido nada.
func DefaultLayout() Layout { return Layout{MaxUnits: 20, MaxLines: 2000} }

// Big diz se um recorte passou do corte — é a pergunta que o template faz.
func (l Layout) Big(s Size) bool {
	return s.Units > l.MaxUnits || s.Lines > l.MaxLines
}

// Describe diz, em uma frase, o que este layout decidiu para um recorte. Vai no topo da
// página gerada, para que quem a lê saiba que está vendo um resumo e onde está o resto.
func (l Layout) Describe(s Size) string {
	if !l.Big(s) {
		return ""
	}
	return fmt.Sprintf("Esta camada tem %d unidades e %d regras — acima do corte de %d "+
		"unidades / %d linhas, então esta página traz o RESUMO de cada unidade. "+
		"O texto completo está na spec.", s.Units, s.Rules, l.MaxUnits, l.MaxLines)
}
