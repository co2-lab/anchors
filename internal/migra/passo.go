package migra

import (
	"fmt"
	"sort"
)

// UM PASSO POR VERSÃO DE FORMATO, e a corrente que os liga.
//
// O migrador não sabe converter "do formato antigo para o atual" — ele sabe converter de N
// para N+1, e aplica os passos em sequência. Um projeto parado no formato 1 com o binário
// no 4 atravessa 1→2, 2→3, 3→4, nessa ordem.
//
// POR QUE ASSIM, E NÃO UMA FUNÇÃO QUE OLHA O ESTADO E CONSERTA.
//
// A função única parece mais simples e envelhece mal: ela precisa reconhecer TODOS os
// estados intermediários possíveis, e cada formato novo acrescenta uma combinação a mais
// para ela distinguir. Quem escreve o passo 5 precisa lembrar do 1.
//
// O passo isolado só precisa saber de UMA transição. O do formato 2 foi escrito uma vez,
// está correto, e nunca mais muda — o formato 3 é um arquivo novo que não o toca.
//
// É a mesma razão pela qual migração de banco de dados se escreve assim, e o motivo é o
// mesmo: o estado de partida de quem atualiza não é conhecido por quem escreve.

// Step é uma transição de formato.
type Step struct {
	// To: o formato que este passo PRODUZ. O passo que produz 2 converte de 1.
	To int
	// Why: o que mudou, em uma linha. Vai para a saída do `anchors migrate` — quem
	// revisa o diff precisa saber o que esperar antes de abri-lo.
	Why string
	// RenameKeys: chaves YAML que mudaram de nome, por arquivo (basename → de → para).
	RenameKeys map[string]map[string]string
	// RenameValues: valores que mudaram, por arquivo e por chave. Diferente de renomear
	// a chave: aqui a chave fica e o VALOR dela vira outro — é o caso dos nomes de gate
	// gravados dentro dos carimbos de julgamento.
	RenameValues map[string]map[string]map[string]string
}

// steps são todas as transições conhecidas, uma por formato.
//
// Registrar aqui é o que torna a atualização barata: acrescentar um formato é acrescentar
// uma entrada, e nada mais precisa saber que ela existe.
var steps []Step

// Register adiciona um passo. Chamado do `init()` de cada arquivo `formato_N.go`.
//
// Um arquivo por formato, e não uma lista central: quem escreve o formato 3 cria
// `formato_3.go` e não toca em linha nenhuma do 2. Um arquivo central seria editado a cada
// versão, e é onde os conflitos de merge aconteceriam.
func Register(p Step) {
	steps = append(steps, p)
	sort.Slice(steps, func(i, j int) bool { return steps[i].To < steps[j].To })
}

// StepsFrom devolve os passos que levam de `origem` a `destino`, em ordem.
//
// Um buraco na corrente é ERRO, não silêncio: se falta o passo que produz o formato 3, um
// projeto no 2 não pode ir para o 4 fingindo que o 3 não existia — o arquivo ficaria com o
// número novo e o conteúdo velho.
func StepsFrom(origem, destino int) ([]Step, error) {
	var out []Step
	esperado := origem + 1
	for _, p := range steps {
		if p.To <= origem || p.To > destino {
			continue
		}
		if p.To != esperado {
			return nil, fmt.Errorf("falta o passo de migração que produz o formato %d "+
				"(o próximo registrado é o %d) — sem ele o arquivo receberia o número "+
				"novo com o conteúdo antigo", esperado, p.To)
		}
		out = append(out, p)
		esperado++
	}
	if esperado != destino+1 {
		return nil, fmt.Errorf("falta o passo de migração que produz o formato %d", esperado)
	}
	return out, nil
}

// AllSteps devolve os passos registrados, em ordem.
//
// Existe para as RÉGUAS: um teste confere que todo destino de conversão corresponde a um
// gate que existe — um de-para que aponta para nome inexistente converteria o projeto para
// um gate que o `check` não encontra, e ele sumiria sem nada acusar.
func AllSteps() []Step {
	out := make([]Step, len(steps))
	copy(out, steps)
	return out
}
