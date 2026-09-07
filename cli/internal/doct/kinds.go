package doct

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// --- o que cada tipo de documentação PEDE ---
//
// Nomear o arquivo não basta. "Atualize `docs/openapi.yaml`" leva um agente a acrescentar
// o endpoint que acabou de escrever e parar ali — sem os erros, sem os exemplos, sem o
// esquema do corpo. O resultado passa em qualquer verificação de existência e é inútil
// para quem consome: um contrato que só descreve o caminho feliz não é contrato.
//
// Cada tipo aqui diz o que a documentação daquele tipo precisa RESPONDER. É a mesma
// escolha que o Anchors faz nos gates: a instrução carrega o porquê, porque instrução sem
// porquê vira ritual — o agente cumpre a forma e erra o conteúdo.

// Instruction é o que o agente precisa saber para produzir uma documentação de um tipo.
type Instruction struct {
	Titulo string
	// Pede são os itens que a documentação daquele tipo tem de cobrir.
	Pede []string
	// Armadilha é o modo de falha típico daquele tipo — o jeito de "cumprir" a tarefa
	// produzindo algo que não serve.
	Armadilha string
}

var instructions = map[string]Instruction{
	config.KindOpenAPI: {
		Titulo: "OpenAPI — o contrato dos endpoints",
		Pede: []string{
			"cada rota com método, caminho e o que ela faz",
			"o esquema do CORPO de entrada e da resposta, com os tipos reais",
			"os códigos de ERRO e quando cada um acontece — não só o 200",
			"autenticação: qual esquema a rota exige, e o que ela devolve sem ele",
			"exemplo de requisição e de resposta para cada rota",
		},
		Armadilha: "acrescentar a rota nova e parar aí. Um contrato que descreve só o " +
			"caminho feliz não é contrato: quem consome descobre os erros em produção.",
	},
	config.KindC4: {
		Titulo: "C4 — a arquitetura em quatro níveis",
		Pede: []string{
			"Nível 1, Contexto: o sistema, quem o usa, e os sistemas externos com que " +
				"fala — as FONTES que os planos nomeiam entram aqui",
			"Nível 2, Contêineres: o que roda separado (app, API, banco, fila) e como " +
				"cada par se comunica, com o protocolo",
			"Nível 3, Componentes: dentro de um contêiner, as peças e suas " +
				"responsabilidades — as CAMADAS do `anchors.yaml` são o esqueleto disto",
			"Nível 4, Código: só onde há algo não óbvio; o mapa (`anchors map show`) já " +
				"é a versão de máquina",
			"em cada nível, o que está FORA dele — o C4 vale pelo que cada nível omite",
		},
		Armadilha: "um diagrama só, com tudo. É a falha clássica do desenho de " +
			"arquitetura, e é o que o C4 existe para evitar: separar em níveis é o " +
			"mecanismo, não a decoração.",
	},
	config.KindSchema: {
		Titulo: "Esquema de dados — as tabelas e o que as liga",
		Pede: []string{
			"cada tabela/coleção com os campos, tipos e o que é obrigatório",
			"as CHAVES e os índices, e a consulta que justifica cada índice",
			"as relações, e o que acontece com a filha quando a mãe morre",
			"o que é PII, e como o projeto a trata — anonimização, retenção",
			"as migrações: o esquema é um estado, e o caminho até ele importa",
		},
		Armadilha: "documentar as colunas e omitir os índices. O esquema sem os " +
			"índices descreve o que se guarda e não o que se consegue perguntar.",
	},
	config.KindComponent: {
		Titulo: "Catálogo de componentes",
		Pede: []string{
			"cada componente com as props, os tipos, e quais são obrigatórias",
			"os ESTADOS: vazio, carregando, erro, e o cheio — não só o cheio",
			"acessibilidade: papel, rótulo, e o que o teclado alcança",
			"quando NÃO usar este componente, e qual usar em vez dele",
		},
		Armadilha: "listar as props e mostrar um exemplo no estado cheio. Os estados " +
			"vazio e de erro são onde a interface realmente decide, e são os que ninguém desenha.",
	},
	config.KindADR: {
		Titulo: "ADR — as decisões de arquitetura",
		Pede: []string{
			"o contexto: o que estava em jogo quando se decidiu",
			"a decisão, em uma frase",
			"as alternativas consideradas e por que foram recusadas",
			"as consequências, inclusive as ruins — uma ADR só com vantagens não foi decisão",
			"o estado: proposta, aceita, substituída (e por qual)",
		},
		Armadilha: "escrever a ADR depois, justificando o que já foi feito. A que " +
			"registra a escolha vale; a que a defende é publicidade.",
	},
	config.KindRunbook: {
		Titulo: "Runbook — o que fazer quando quebra",
		Pede: []string{
			"os sintomas, como quem está de plantão os vê (o alarme, não a causa)",
			"o diagnóstico: o que olhar, em que ordem, e o que cada resposta significa",
			"a ação, com os comandos reais — quem lê isto está sob pressão",
			"o que NÃO fazer, e por quê",
			"como saber que voltou ao normal",
		},
		Armadilha: "escrever o runbook a partir da causa. Quem está de plantão " +
			"começa pelo sintoma e não sabe a causa — é o que ele está tentando descobrir.",
	},
}

// InstructionFor devolve o que o agente precisa saber sobre uma documentação.
//
// Tipo desconhecido devolve uma instrução mínima em vez de erro: o projeto pode declarar
// uma documentação que o Anchors não conhece, e recusá-la faria o mecanismo servir só ao
// que já foi previsto. O que ele não sabe instruir, ele ao menos nomeia.
func InstructionFor(d config.DocArtifact) Instruction {
	if i, ok := instructions[strings.ToLower(strings.TrimSpace(d.Kind))]; ok {
		return i
	}
	titulo := d.Kind
	if titulo == "" {
		titulo = d.Path
	}
	return Instruction{Titulo: titulo}
}

// KnownKinds são os tipos que o Anchors sabe instruir.
func KnownKinds() []string {
	return []string{config.KindOpenAPI, config.KindC4, config.KindSchema,
		config.KindComponent, config.KindADR, config.KindRunbook}
}

// Duty escreve o que uma documentação exige, para o card do agente.
func Duty(d config.DocArtifact) string {
	i := InstructionFor(d)
	var b strings.Builder
	fmt.Fprintf(&b, "    %s — `%s`\n", i.Titulo, d.Path)
	if d.Why != "" {
		fmt.Fprintf(&b, "      %s\n", d.Why)
	}
	for _, p := range i.Pede {
		fmt.Fprintf(&b, "      · %s\n", p)
	}
	if i.Armadilha != "" {
		fmt.Fprintf(&b, "      ARMADILHA: %s\n", i.Armadilha)
	}
	return b.String()
}
