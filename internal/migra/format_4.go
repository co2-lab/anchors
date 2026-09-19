package migra

// FORMATO 4 — quatro chaves cujo NOME mentia sobre o que elas guardam.
//
// Os formatos 2 e 3 trataram do IDIOMA das chaves. Este trata da COERÊNCIA: chaves cujo
// nome não diz o que elas são, e que por isso são lidas errado antes mesmo de serem
// usadas. O custo não é estético — é que quem escreve o `anchors.yaml` decide a partir do
// nome, e um nome que sugere a coisa errada produz configuração errada em silêncio.
//
// As quatro, e o que cada nome fazia o leitor supor:
//
//   · `auto_judgment` → parece guardar o MODO de julgamento (qual estratégia), e é um
//     interruptor. O prefixo `enable_` é o que diz "isto liga e desliga" — e a chave que
//     nasceu depois dela (`enforce_section_language`) já seguiu essa forma, deixando as
//     duas convenções convivendo no mesmo arquivo.
//
//   · `triad_optional` → parece booleano ("a trinca é opcional?"), e é uma LISTA de quais
//     arestas da trinca ficam dispensadas (`[tested-by]`). Um booleano e uma lista não se
//     confundem só na leitura: quem supõe booleano escreve `triad_optional: true`, o
//     parser recusa, e o erro aponta para o tipo em vez do entendimento.
//
//   · `requires_code` → não diz de QUEM se exige o código. Fica sob `rule_types`, e o que
//     ela nomeia são as SEÇÕES que cobram código do cenário — não o tipo de regra, não o
//     projeto. Sem o sujeito no nome, a chave é lida como "este tipo de regra exige
//     código", que é outra coisa.
//
//   · `rule_marking` → parece guardar a marcação em si, e guarda a POLÍTICA sobre ela
//     (`required`). O sufixo `_policy` devolve o sujeito: o que se declara é o regime, não
//     o dado.
//
// PERDA SE NÃO CONVERTER: diferente dos nomes de gate (formatos 2 e 3), aqui a falha é
// BARULHENTA — o parser recusa a chave desconhecida e o `check` para com a mensagem que o
// `config` já sabe dar ("o ARQUIVO está no formato antigo… rode `anchors migrate`"). Não
// há perda silenciosa de carimbo. Mesmo assim a conversão é obrigatória: sem ela, todo
// projeto existente para de carregar até alguém editar à mão.
//
// Só `anchors.yaml`: nenhuma das quatro aparece no grafo — são declaração de política, e o
// mapa guarda medição.

func init() {
	Register(Step{
		To:  4,
		Why: "quatro chaves renomeadas para dizerem o que guardam (interruptor, lista, sujeito)",
		RenameKeys: map[string]map[string]string{
			"anchors.yaml": {
				"auto_judgment":  "enable_auto_judgment",
				"triad_optional": "optional_triad_edges",
				"requires_code":  "sections_require_code",
				"rule_marking":   "rule_marking_policy",
			},
		},
	})
}
