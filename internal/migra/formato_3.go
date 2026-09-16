package migra

// FORMATO 3 — os OITO nomes de gate que o formato 2 deixou para trás.
//
// O formato 2 traduziu 21 nomes e parou. Ficaram `cenario-identidade`,
// `cenario-letra-declarada`, `cenario-tipo-alinhado`, `regra-implementada`,
// `header-conforme`, `testid-coerente`, `testid-consultado-existe` e
// `placeholder-preenchido` — metade do vocabulário num idioma e metade no outro.
//
// O custo não é estético. Um nome de gate é IDENTIFICADOR: aparece no `anchors.yaml` de
// todo projeto, na saída do `check`, nos carimbos de julgamento e nas issues que o
// pipeline abre. Enquanto metade está em português, cada gate novo herda a dúvida sobre
// qual convenção seguir — e foi assim que estes oito nasceram, depois da primeira
// tradução.
//
// Dois deles não eram nem português: `header-conforme` é meia palavra em cada idioma, e
// "conforme" traduzido ao pé da letra não diz o que o gate faz (ele confere que o header
// é VÁLIDO). Traduzir por palavra, e não por sentido, produz nome pior que o original.
//
// A ALTERNATIVA seria uma tabela de alias no `runInternal`, lida para sempre. O próprio
// `migra.go` já a condena: "ela não fecha — o arquivo nunca se conserta, cada leitor
// precisa saber dos dois, e a terceira renomeação...". Converte-se uma vez.
//
// O QUE PRECISA SER CONVERTIDO, e por que cada lugar:
//
//   · `anchors.yaml` — `name`, `id` e `check`. Sem isso o `check` responde "checker
//     interno desconhecido" e o gate deixa de rodar: falha barulhenta, mas o projeto para.
//
//   · `anchors.graph.yaml` — o `gate:` dentro dos julgamentos. Esta é a perda SILENCIOSA:
//     um carimbo gravado como `regra-implementada` deixa de casar com o gate
//     `rule-implemented`, o `check` não o encontra e refaz o julgamento — cobrando de novo
//     o que alguém já respondeu. Foi exatamente o defeito que o formato 2 registrou ter
//     medido: 40 julgamentos convivendo com 2 do mesmo gate sob outro nome.

var gates3 = map[string]string{
	"cenario-identidade":       "scenario-identity",
	"cenario-letra-declarada":  "scenario-letter-declared",
	"cenario-tipo-alinhado":    "scenario-type-aligned",
	"header-conforme":          "header-valid",
	"placeholder-preenchido":   "placeholder-filled",
	"regra-implementada":       "rule-implemented",
	"testid-coerente":          "testid-consistent",
	"testid-consultado-existe": "testid-queried-exists",
}

func init() {
	Register(Step{
		To:  3,
		Why: "os oito nomes de gate que ficaram em português na tradução anterior",
		RenameValues: map[string]map[string]map[string]string{
			"anchors.yaml": {
				// `name` é como o PROJETO chama a instância do gate, e `check` é o
				// verificador interno que ela invoca. São os dois lugares onde estes
				// oito aparecem hoje — medido: três deles no `anchors.yaml` do
				// blue-eyes, nas duas chaves.
				"name":  gates3,
				"check": gates3,
			},
			// `gate:` (nos julgamentos do grafo) e `id:` ficam de FORA, e não por
			// esquecimento: nenhum destes oito tem gate nos defaults, então nenhum
			// projeto os declara por `id` nem carimba julgamento com eles. Conferido no
			// blue-eyes: zero ocorrências no `anchors.graph.yaml`.
			//
			// Migrar chave onde o valor não existe não é inofensivo — `id`/`gate`
			// apontam para gates DECLARADOS, e a régua `TestVocabularioAntigoAponta…`
			// reprova um de-para cujo destino não está nos defaults. Ela está certa: uma
			// migração que renomeia para um gate inexistente troca um nome morto por
			// outro, e o projeto fica sem o gate do mesmo jeito.
		},
	})
}
