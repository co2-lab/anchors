package initx

import "github.com/co2-lab/anchors/internal/config"

// dependemDeSinalIngerido são os gates que só têm o que medir depois de `anchors ingest`
// receber um relatório de teste, cobertura ou mutação. Bloquear com base num sinal que
// ainda não existe barraria o commit por ausência de dado — não por defeito.
//
// Ficam informativos mesmo em projeto novo, e o usuário os promove quando a suíte
// estiver rodando no CI.
var dependemDeSinalIngerido = map[string]bool{
	"tests-green": true, "line-coverage": true, "coverage-delta": true,
	"mutation-score": true, "scenario-coverage": true, "sbom-gerado": true,
	"dependencia-vulneravel": true, "sem-duplicacao": true,
}

// DefaultGates devolve os gates que um projeto deve NASCER com, conforme os artefatos
// que ele usa. É o que amarra os sinais de teste ao ciclo de vida: se o projeto tem
// testes, o ciclo já nasce cobrando execução + cobertura, em vez de o usuário escrever
// tudo à mão.
//
// `projetoNovo` decide o estado de maturação, e a distinção é o ponto:
//
//   - projeto EXISTENTE → informativo (QUALITY §7). A doutrina é explícita: "impor o
//     gate como bloqueante imediatamente pararia o projeto", porque um projeto real
//     quase nunca cumpre, no dia um, o limiar que quer atingir.
//   - projeto NOVO → BLOQUEANTE. Ali a premissa se inverte: não há débito a acomodar,
//     e o gate não para nada — ele impede o PRIMEIRO desvio, que é quando corrigir
//     custa menos. Nascer informativo adiaria uma cobrança que nunca vai ser mais
//     barata do que agora.
//
// A exceção são os gates que dependem de sinal ingerido: sem relatório, eles barrariam
// por ausência de dado, não por defeito.
//
// `chosen` são os artefatos escolhidos no init (spec/feature/test/guide/plan).
func DefaultGates(chosen map[string]bool, projetoNovo bool) []config.Gate {
	var gates []config.Gate

	if chosen["spec"] {
		gates = append(gates,
			config.Gate{
				Name: "spec-completa", ID: "spec-completa", On: []string{"spec"}, Check: "spec-sections",
				Blocking: config.Bool(false), Measures: "a spec tem ao menos um estado/regra, sem placeholder",
			},
			config.Gate{
				Name: "spec-tem-codigo", ID: "spec-tem-codigo", On: []string{"spec"}, Check: "has-code",
				Blocking: config.Bool(false), Measures: "a spec carrega um código de cenário (identidade)",
			},
			// A spec sozinha atravessa TODOS os gates relacionais — eles falham ABERTO
			// (sem teste ligado, não há o que confrontar) e o pipeline conclui "pode
			// promover" sobre trabalho que não existe. Este gate pergunta o oposto: as
			// peças EXISTEM? Nasce informativo porque quase todo projeto tem débito.
			config.Gate{
				Name: "trinca-completa", ID: "trinca-completa", On: []string{"spec"}, Check: "trinca-completa",
				Blocking: config.Bool(false), Measures: "a spec tem código, feature e teste que a realizam",
			},
			// A metade que o determinístico NÃO alcança.
			//
			// O `trinca-completa` confere que a referência do `@no-test` RESOLVE: existe
			// um teste que menciona aquele código. Isso é binário e ele decide sozinho.
			//
			// O que sobra é julgamento — aquele teste prova mesmo ESTE comportamento? A
			// busca do gate determinístico é TEXTUAL de propósito: casar `it('CODE: …')`
			// amarraria o Anchors ao dialeto de um framework de teste, e um projeto Go ou
			// Python perderia a dispensa. O preço dessa neutralidade é uma folga — um
			// código citado só em comentário resolve a referência sem provar nada.
			//
			// Sem este gate a folga vira rigor APARENTE: a referência resolve, o pipeline
			// fica verde, e a prova não existe. Determinístico prova que a referência
			// RESOLVE; julgamento prova que ela VALE.
			config.Gate{
				Name: "no-test-prova-real", ID: "no-test-prova-real", On: []string{"spec"},
				// Só as specs que DISPENSARAM o teste têm o que julgar. Sem o filtro, o
				// gate enfileirava uma pergunta de IA para toda spec do projeto — 583
				// alvos onde 16 dispensaram — e o contador de pendências passava a medir
				// o tamanho do projeto em vez do tamanho do trabalho.
				Requires: "@no-test",
				Blocking: config.Bool(false), Measures: config.MeasuresJudgment,
				Ask: "Esta spec dispensa o teste próprio com `@no-test` e aponta um código " +
					"de cenário como prova em outro lugar. Encontre o teste que menciona esse " +
					"código e responda: ele EXERCITA o comportamento que esta spec descreve, " +
					"ou apenas cita o código? Reprove se a menção for de comentário, se o " +
					"teste exercitar outro comportamento, ou se a prova for indireta a ponto " +
					"de uma quebra nesta unidade não derrubá-lo. A pergunta é sobre o que o " +
					"teste EXECUTA, não sobre o que o nome dele promete.",
			},
			// O PAR do `regra-implementada`, e a divisão entre eles é de NATUREZA, não de
			// rigor:
			//
			//   determinístico  o código de referência ESTÁ no trecho? — barato, roda a
			//                   cada commit, e é falsificável por construção (basta
			//                   escrever o comentário).
			//   julgamento      o trecho marcado FAZ o que a regra diz? — caro, roda em
			//                   revisão, e é o único que pega o comentário mentiroso.
			//
			// Um sozinho não basta. Só o determinístico certifica quem marcou sem
			// implementar; só o julgamento seria caro demais para cada commit. Rodando em
			// momentos diferentes, eles se complementam.
			config.Gate{
				Name: "regra-cumprida", ID: "regra-cumprida", On: []string{"spec"},
				Blocking: config.Bool(false), Measures: config.MeasuresJudgment,
				Ask: "Cada regra desta spec está marcada no código por um comentário com o " +
					"código dela (`// ABCDX-B01: …`). Leia a regra e o trecho que ela marca e " +
					"responda: o trecho REALIZA o que a regra descreve? Reprove quando o " +
					"comentário estiver sobre código que faz outra coisa, quando a regra " +
					"descrever um caso que o trecho não trata, ou quando a marcação estiver " +
					"num lugar genérico (topo do arquivo, import) em vez do trecho que decide. " +
					"A pergunta é sobre o que o código EXECUTA, não sobre o que o comentário " +
					"afirma. Ao reprovar, PROPONHA a correção como patch.",
			},
			// O vocabulário de letras do código é do PROJETO (`rule_types`), mas o gate
			// que impede conflito e letra não declarada é universal: uma letra fora do
			// vocabulário torna a regra INVISÍVEL à rastreabilidade — ela aparece na
			// spec, na feature e no teste, e nenhum gate a enxerga.
			// O INVERSO do `regra-implementada`: aquele parte da spec e pergunta "esta
			// regra tem código?"; este parte do código e pergunta "este símbolo tem
			// regra?". Sem os dois, a divergência escapa por um dos lados — medido: uma
			// spec catalogava 2 regras para 7 funções exportadas e ninguém perguntou
			// pelas 5 restantes. A dispensa por símbolo (`@no-rule: <razão>`) cobre o
			// caso legítimo: nem toda exportação carrega decisão de negócio.
			config.Gate{
				Name: "codigo-catalogado", ID: "codigo-catalogado", On: []string{"spec"}, Check: "codigo-catalogado",
				Blocking: config.Bool(false),
				Measures: "todo símbolo exportado tem regra na spec ou dispensa escrita",
			},
			config.Gate{
				Name: "rule-types", ID: "rule-types", On: []string{"spec"}, Check: "rule-types",
				Blocking: config.Bool(false), Measures: "toda letra de código é declarada no vocabulário, sem conflito",
			},
			// A Tabela de Dependências promete símbolos; o código precisa usá-los.
			// Pega a divergência que o feature-test-match não vê (a aresta spec→código).
			config.Gate{
				Name: "dependency-honored", ID: "dependency-honored", On: []string{"spec"}, Check: "dependency-honored",
				Blocking: config.Bool(false), Measures: "os métodos declarados na Tabela de Dependências são usados no código",
			},
			// O padrão mais repetido de uma auditoria de 51 divergências spec×código
			// (app de referência, 2026-08): oito handlers declaravam um Contrato de Saída que o
			// código não cumpria, e o status omitido era quase sempre o de SEGURANÇA
			// (403 de ownership, 409 de conflito) — quem escreve a tabela pensa no
			// caminho felizardo. O inverso também conta: status declarado que nenhum
			// caminho emite é código morto no cliente, e some sem ninguém notar.
			// Os três achados MAIS GRAVES da mesma auditoria tinham forma idêntica:
			// dois lados definiam a mesma coisa, cada um tinha teste, e cada teste
			// confrontava a PRÓPRIA cópia. A trinca ficava completa e os gates verdes
			// porque nenhum perguntava se a PROVA alcança o outro lado. Este gate
			// cobra: regra que afirma "espelha"/"fonte única" com um arquivo citado
			// exige que o código IMPORTE aquele arquivo.
			config.Gate{
				Name: "prova-cruza-fronteira", ID: "prova-cruza-fronteira", On: []string{"spec"}, Check: "prova-cruza-fronteira",
				Blocking: config.Bool(false), Measures: "regra que afirma relação com outra unidade tem o código importando aquela unidade",
			},
			config.Gate{
				Name: "contract-status-declared", ID: "contract-status-declared", On: []string{"spec"}, Check: "contract-status-declared",
				Blocking: config.Bool(false), Measures: "os status do Contrato de Saída são os que o handler devolve — e só eles",
			},
			// O dever que um artefato contrai com um lugar que ele NÃO conhece (LGPD,
			// i18n, a11y, auditoria). Só age se o projeto declarar `obligations:` — sem
			// declaração, o gate é inerte. Nasce ligado para que a categoria exista.
			config.Gate{
				Name: "obligation-honored", ID: "obligation-honored", On: []string{"spec"}, Check: "obligation-honored",
				Blocking: config.Bool(false), Measures: "o artefato cumpre as obrigações transversais que contrai",
			},
		)
	}
	// FUNÇÕES IRMÃS: quando a maioria guarda um parâmetro e uma não, a que não guarda é
	// quase sempre esquecimento. Detectável sem entender o domínio — é assimetria.
	if chosen["test"] || chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "sibling-guard", ID: "sibling-guard", On: []string{"code"}, Check: "sibling-guard",
			Blocking: config.Bool(false), Measures: "funções irmãs tratam o mesmo parâmetro de forma consistente",
		})
	}
	// PLANO: um plano semeia specs, e semear numa camada que não tem spec faz quem
	// executa gastar uma rodada descobrindo a contradição. A origem é checável.
	if chosen["plan"] {
		gates = append(gates, config.Gate{
			Name: "plan-seeds-valid", ID: "plan-seeds-valid", On: []string{"plan"}, Check: "plan-seeds-valid",
			Blocking: config.Bool(false), Measures: "o plano só semeia spec em camada que tem spec",
		})
		// A ORDEM dentro do plano. Um plano sem fases catalogadas passa (elas são
		// opcionais); o gate só cobra a coerência de quem as declarou.
		gates = append(gates, config.Gate{
			Name: "fase-ordenada", ID: "fase-ordenada", On: []string{"plan"}, Check: "fase-ordenada",
			Blocking: config.Bool(projetoNovo),
			Measures: "as fases do plano não dependem do que vem depois delas",
		})
	}
	// PERTENCIMENTO vale para qualquer artefato com header, não só para spec: uma fase
	// pertence a um plano, e um plano pode pertencer a outro. Mas SÓ ENTRA se o projeto
	// tem algum desses artefatos — um projeto vazio não semeia gate nenhum, e um gate que
	// nasce sem nada para medir é a impressão de defesa que não existe.
	if chosen["spec"] || chosen["plan"] || chosen["code"] {
		var onde []string
		for _, k := range []string{"spec", "plan", "code"} {
			if chosen[k] {
				onde = append(onde, k)
			}
		}
		gates = append(gates, config.Gate{
			Name: "parent-valido", ID: "parent-valido", On: onde,
			Check: "parent-valido", Blocking: config.Bool(projetoNovo),
			Measures: "o `parent:` declarado aponta para algo que existe, sem ciclo",
		})
	}
	if chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "fase-existe", ID: "fase-existe", On: []string{"spec"}, Check: "fase-existe",
			Blocking: config.Bool(projetoNovo),
			Measures: "o `needs:` da spec aponta para uma fase que algum plano cataloga",
		})
	}
	if chosen["feature"] {
		gates = append(gates, config.Gate{
			Name: "feature-nao-vazia", ID: "feature-nao-vazia", On: []string{"feature"}, Check: "non-empty",
			Blocking: config.Bool(false), Measures: "a feature não é um esqueleto vazio",
		})
		// scenario-asserts vai além do 'não está vazio': o passo de RESULTADO precisa
		// afirmar um resultado. "Então o efeito XXXXX-B01 se verifica" satisfaz todos os
		// outros gates e não diz nada — a definição do que 'se verifica' migra para o
		// teste, que passa a ser escrito sem saber qual era o desfecho. Medido: 151 de 442
		// features de um projeto real, e foi a razão estrutural pela qual uma regra ficou
		// sem caso discriminante.
		gates = append(gates, config.Gate{
			Name: "scenario-asserts", ID: "scenario-asserts", On: []string{"feature"}, Check: "scenario-asserts",
			Blocking: config.Bool(false), Measures: "o passo de resultado afirma um resultado observável",
		})
	}
	// Os gates de TESTE — o que amarra os sinais ingeridos ao ciclo. Só fazem sentido
	// se o projeto tem testes E specs (a cobertura por cenário cruza os dois).
	if chosen["test"] {
		gates = append(gates,
			// Um teste que existe, passa e cobre o comportamento certo — mas não cita
			// código — é INVISÍVEL para os gates relacionais: o `feature-test-match`
			// reporta os cenários como não implementados, e quem for consertar escreve um
			// segundo teste do mesmo comportamento. Medido duas vezes num projeto real
			// (códigos invertidos em 3 telas; store com 5 casos sem código nenhum).
			config.Gate{
				Name: "teste-rastreavel", ID: "teste-rastreavel", On: []string{"test"}, Check: "teste-rastreavel",
				Blocking: config.Bool(false),
				Measures: "o teste cita o código do que prova (é visível aos gates relacionais)",
			},
			config.Gate{
				Name: "tests-green", ID: "tests-green", On: []string{"test"}, Check: "tests-pass",
				Blocking: config.Bool(false), Measures: "os testes deste arquivo passam (do resultado ingerido)",
			},
			// O PAR que ataca PROVA FALSA, e não ausência de prova: um dublê que não
			// deriva do módulo real segue verde depois que o módulo muda. `tests-green`
			// diz que passou, `feature-test-match` que o cenário casa, `trinca-completa`
			// que o teste existe — os três respondem "sim" sobre um teste que mente.
			//
			// São camadas, não alternativas. O `mock-carimbado` é o AGNÓSTICO: hash de
			// texto, pega assinatura, corpo, tipo e constante em qualquer linguagem, e é
			// o único que alcança a forma do VALOR devolvido. O `mock-tipado` só age
			// onde a linguagem tem tipos estruturais, mas ali sai de graça — o
			// compilador já confere, sem carimbo a manter.
			// O PONTO CEGO do `mock-carimbado`: o `mock_detect` é um regex escrito à
			// mão, e um padrão SINTATICAMENTE válido mas errado para o dialeto do
			// projeto casa nada — o gate varre zero dublês e passa em silêncio. É a
			// falha mais perigosa da família, porque o verde vem de não ter olhado.
			//
			// Nenhuma checagem determinística cobre isto: para saber se o padrão alcança
			// o que o projeto escreve é preciso LER os testes e reconhecer a forma de
			// dublê do ecossistema. Daí ser julgamento — e daí a sugestão vir como
			// patch: quem julga já leu o suficiente para propor o `mock_detect` correto.
			config.Gate{
				Name: "mock-detect-cobre-o-dialeto", ID: "mock-detect-cobre-o-dialeto", On: []string{"test"},
				Blocking: config.Bool(false), Measures: config.MeasuresJudgment,
				Ask: "O projeto declara `derived.mock_detect` — o regex que reconhece um " +
					"dublê de teste neste ecossistema. Leia este arquivo de teste e responda: " +
					"o padrão declarado alcança TODAS as formas de dublê que ele usa? " +
					"Reprove se houver dublê que o regex não casa (outra função, outro " +
					"dialeto, decorator em vez de chamada, dublê de módulo inteiro vs de " +
					"membro). Um regex que casa zero faz o gate `mock-carimbado` reportar " +
					"verde sem ter conferido nada. Ao reprovar, PROPONHA o padrão corrigido " +
					"como patch do `anchors.yaml`.",
			},
			config.Gate{
				Name: "mock-carimbado", ID: "mock-carimbado", On: []string{"test"}, Check: "mock-carimbado",
				Blocking: config.Bool(false),
				Measures: "o carimbo do dublê corresponde ao trecho real (recalculado, não só validado)",
			},
			config.Gate{
				Name: "mock-tipado", ID: "mock-tipado", On: []string{"test"}, Check: "mock-tipado",
				Blocking: config.Bool(false),
				Measures: "o dublê de teste deriva do módulo real (não é cópia congelada do contrato)",
			},
			config.Gate{
				Name: "line-coverage", ID: "line-coverage", On: []string{"code"}, Check: "line-coverage",
				Blocking: config.Bool(false), Measures: "cobertura de linha >= limiar (do lcov ingerido)",
			},
			config.Gate{
				Name: "coverage-delta", ID: "coverage-delta", On: []string{"code"}, Check: "coverage-delta",
				Blocking: config.Bool(false), Measures: "a cobertura de linha não caiu vs. a ingestão anterior",
			},
			// mutation-score nasce com o projeto porque é o ÚNICO gate que responde "o
			// teste prova a linha?" — todos os outros respondem "a linha executou?" ou
			// "a peça existe?". Medido: um arquivo com 16 testes verdes e cobertura
			// cheia onde apagar uma guarda do código não derrubava teste nenhum.
			//
			// Nasce informativo e fica Pending até o projeto ingerir o sinal — o Anchors
			// não pode EXIGIR a ferramenta (depende do stack), mas também não deve
			// esconder o que se perde sem ela: o doctor alerta, e o projeto decide
			// quando virar blocking.
			config.Gate{
				Name: "mutation-score", ID: "mutation-score", On: []string{"code"}, Check: "mutation-score",
				Blocking: config.Bool(false), Measures: "o teste PROVA a linha: mutantes mortos >= limiar (do relatório de mutação ingerido)",
			},
		)
		// pagination-honored é da mesma família: pega o que NENHUM teste pega, porque o
		// teste roda com 3 registros e a produção com 3 mil. Medido no app de referência: 11 arquivos,
		// entre eles uma listagem de importação de extrato que perde silenciosamente as
		// transações além do 1º MB — o usuário importa o arquivo e parte some sem erro.
		gates = append(gates, config.Gate{
			Name: "pagination-honored", ID: "pagination-honored", On: []string{"code"}, Check: "pagination-honored",
			Blocking: config.Bool(false), Measures: "função que promete o conjunto não devolve a primeira página em silêncio",
		})
	}
	// layer-boundary confronta `boundaries:` — o que cada camada NÃO alcança. Nasce com
	// TODO projeto que declarou camadas, porque declarar `layers:` sem defendê-las é
	// desenhar a arquitetura e deixá-la como documentação. Fica Pendente (dizendo o que
	// declarar) até o projeto escrever suas fronteiras.
	//
	// A lacuna era visível: um projeto real mantinha 366 linhas de shell com 15 regras
	// arquiteturais escritas à mão — todas da mesma forma, reimplementadas porque o
	// framework não oferecia o mecanismo. E o script tinha um escopo implícito no caminho
	// do grep (só varria o mobile) que ninguém via ao ler a regra.
	// Condicionado a `code`: um projeto que não rastreia código não tem camada a defender
	// (e o teste de que projeto sem artefato não semeia gate algum pegou isso).
	if chosen["code"] {
		gates = append(gates, config.Gate{
			Name: "layer-boundary", ID: "layer-boundary", On: []string{"code"}, Check: "layer-boundary",
			Blocking: config.Bool(false), Measures: "a camada não alcança o que não é dela (`boundaries:`)",
		})
	}

	// open-questions-resolved é o mais barato dos gates e pega a classe mais cara: a
	// AMBIGUIDADE que ninguém resolveu. A spec não decide; quem implementa escolhe;
	// a escolha nunca chega a quem tinha a resposta. Nenhum outro gate pega, porque todas
	// as peças existem e se referenciam — o defeito é uma decisão que ninguém tomou.
	// Só cobra quem ABRIU a seção, então não vira ritual em spec simples.
	if chosen["spec"] {
		// domain-declared: a spec diz o que a unidade ACEITA, e QUEM garante que o
		// inválido não chega. Medido: 71% das specs de um projeto real catalogavam efeitos
		// e só 14% declaravam entrada — e todo defeito de borda de três rodadas de review
		// morava no que ninguém declarou. `constraints` não substitui: dizer "não valido"
		// empurra o dever para fora e CRIA órfão; a coluna do dono obriga a nomear alguém.
		// count-honored: a mentira que envelhece SOZINHA. "Os 50 modelos do produto" vira
		// falso quando alguém adiciona o 51º — ninguém precisa errar. Medido: um arquivo
		// de spec com 7 afirmações numéricas, onde adicionar UM modelo tornava 10 frases
		// obsoletas. A spec declara COMO conferir; o engine não adivinha o que contar.
		gates = append(gates, config.Gate{
			Name: "count-honored", ID: "count-honored", On: []string{"spec"}, Check: "count-honored",
			Blocking: config.Bool(false), Measures: "o número que a spec afirma bate com o código",
		})

		gates = append(gates, config.Gate{
			Name: "domain-declared", ID: "domain-declared", On: []string{"spec"}, Check: "domain-declared",
			Blocking: config.Bool(false), Measures: "a spec declara o que aceita e quem garante a fronteira",
		})

		gates = append(gates, config.Gate{
			Name: "open-questions-resolved", ID: "open-questions-resolved", On: []string{"spec"}, Check: "open-questions-resolved",
			Blocking: config.Bool(false), Measures: "a spec não tem pergunta em aberto — implementar não é adivinhar",
		})
	}
	// spec-feature-match fecha a ponta que faltava na trinca: feature→test já era
	// confrontado, spec→feature não. Um requisito declarado e sem cenário atravessa o
	// pipeline com TODOS os gates verdes — a spec tem código, a feature existe, a feature
	// bate com o teste. Medido: 43 specs de um projeto real.
	// code-reference-valid pega a ÂNCORA QUE MENTE numa forma que nenhum outro gate via:
	// a spec cita o requisito de uma unidade que não existe. A citação parece
	// rastreabilidade e aponta para o vazio — e um leitor futuro a toma como registro do
	// que foi feito. Medido: uma spec afirmava "índices que o schema criou" citando 4
	// códigos inexistentes, com todos os gates verdes.
	// ref-resolves: o `ref:` aponta para a spec CERTA. `header-conforme` confere que o
	// campo existe; ninguém conferia que ele RESOLVE. Um ref errado parece rastreabilidade
	// e atribui a unidade à spec errada, com todo gate relacional confrontando o par
	// errado. Medido: 49 arquivos de modelo apontando para a identidade de antes de uma
	// desfusão que ninguém propagou, e mais um teste citando duas specs fantasma.
	// Condicionado a `spec`: sem spec no projeto não há `code:` para o `ref:` resolver.
	if chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "ref-resolves", ID: "ref-resolves", On: []string{"code", "feature", "test"}, Check: "ref-resolves",
			Blocking: config.Bool(false), Measures: "o `ref:` aponta para o `code:` da spec irmã",
		})
	}

	if chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "code-reference-valid", ID: "code-reference-valid", On: []string{"spec"}, Check: "code-reference-valid",
			Blocking: config.Bool(false), Measures: "todo código citado pela spec existe no projeto",
		})
	}

	if chosen["spec"] && chosen["feature"] {
		gates = append(gates, config.Gate{
			Name: "spec-feature-match", ID: "spec-feature-match", On: []string{"spec"}, Check: "spec-feature-match",
			Blocking: config.Bool(false), Measures: "todo requisito declarado na spec tem cenário na feature",
		})
		if chosen["spec"] {
			gates = append(gates, config.Gate{
				Name: "scenario-coverage", ID: "scenario-coverage", On: []string{"spec"}, Check: "scenario-coverage",
				Blocking: config.Bool(false), Measures: "cada cenário da spec tem um teste que passou (cobertura semântica)",
			})
		}
	}
	// --- SEGURANÇA E PROCEDÊNCIA (ferramenta externa, agnósticos de linguagem) ---
	//
	// São os primeiros canônicos com `run:`, e só são seguros como canônicos porque
	// declaram `needs_tool`: sem a ferramenta o motor devolve Skip e o `doctor` avisa,
	// em vez de reprovar um projeto recém-criado por algo que ninguém escreveu errado.
	//
	// Agnósticos por construção — nenhum dos três sabe o que é TypeScript, Go ou Rust:
	// gitleaks lê o git, osv-scanner reconhece os lockfiles de 15+ ecossistemas e syft
	// detecta sozinho o que há no projeto. Um gate de licença ficaria de fora aqui: as
	// ferramentas de licença são todas atadas a um gerenciador de pacotes.
	//
	// `on: [code]` porque é o kind que todo projeto tem, qualquer que seja a trinca
	// escolhida — e o escopo é o projeto, então o `on` só decide SE roda, não sobre o quê.
	if chosen["code"] {
		gates = append(gates,
			// Segredo vazado é o único achado sem volta: rotacionar a chave é caro e o
			// histórico do git guarda a original. Por isso é o único da leva que nasce
			// bloqueante — os outros dois informam.
			config.Gate{
				Name: "secret-nao-vazado", ID: "secret-nao-vazado", On: []string{"code"},
				Scope: config.ScopeBatch, ScopeFull: config.ScopeProject,
				Run:       "gitleaks git --no-banner --redact -v",
				NeedsTool: "gitleaks", InstallHint: "brew install gitleaks",
				Blocking: config.Bool(true), When: []string{"pre-commit", "ci"}, Cost: "fast",
				Category: "security",
				Measures: "nenhum segredo (chave, token, credencial) entra no histórico",
			},
			// Informativo: a CVE nova aparece sem ninguém mexer no código, então bloquear
			// pararia um merge por algo que o autor não causou nem pode resolver na hora.
			config.Gate{
				Name: "dependencia-vulneravel", ID: "dependencia-vulneravel", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Run:       "osv-scanner scan source -r .",
				NeedsTool: "osv-scanner", InstallHint: "brew install osv-scanner",
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "security",
				Measures: "as dependências do projeto não têm vulnerabilidade conhecida (OSV)",
			},
			// Cópia-e-cola é dívida que NENHUM outro gate vê: o typecheck passa, o lint passa,
			// os testes passam — cada cópia está correta. O defeito só aparece quando uma delas
			// muda e as outras não, e aí a divergência é silenciosa: a tela continua
			// funcionando, só faz a coisa um pouco diferente da irmã.
			//
			// Medido no projeto que originou o gate: duas telas tinham o MESMO bloco de scroll
			// infinito copiado linha por linha — mesma conta, mesmo limiar, mesmo guard. A
			// terceira cópia teria divergido sem ninguém notar. O jscpd achou em 5ms o que
			// nenhum dos 60 gates existentes viu.
			//
			// Informativo, e é deliberado: nem toda duplicata é dívida. Duas implementações
			// parecidas de coisas que EVOLUEM SEPARADO devem mesmo ficar separadas — unificar
			// por semelhança acidental acopla o que o domínio quer solto. O gate mostra; quem
			// conhece o domínio decide.
			//
			// Sem flags: a calibragem vive no `.jscpd.json` do projeto, que é o mecanismo
			// OFICIAL da ferramenta (limiares, `ignore`, e os marcadores `jscpd:ignore-start`
			// no código). Melhor um mecanismo que a ferramenta mantém do que um nosso, que
			// teríamos de sustentar e explicar — e que envelheceria sozinho.
			//
			// A calibragem NÃO é detalhe: medido no projeto que originou o gate, o default
			// (5 linhas) dava 487 achados, dos quais 375 eram teste × teste — setup e mocks
			// repetidos, semelhança de forma e não de intenção. Em 20 linhas sobram 37, e o
			// que sobra é dívida de verdade: 249 linhas entre dois `seal.ts`, e o par
			// Edit/New de uma mesma entidade. Um gate que nasce com 487 achados não é
			// medição, é ruído — e ruído treina a equipe a ignorar o gate (levando os outros
			// junto).
			config.Gate{
				Name: "sem-duplicacao", ID: "sem-duplicacao", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Run:       "npx --yes jscpd . --reporters console --silent",
				NeedsTool: "npx", InstallHint: "instale Node.js (npx acompanha)",
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "quality",
				Measures: "nenhum bloco de código aparece copiado em dois lugares",
			},
			// Procedência: responde "o que exatamente foi entregue" — a pergunta que só
			// se faz depois do incidente, quando reconstruir a resposta já é impossível.
			config.Gate{
				Name: "sbom-gerado", ID: "sbom-gerado", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Run:       "syft scan dir:. -o cyclonedx-json=.anchors/sbom.json -q",
				NeedsTool: "syft", InstallHint: "brew install syft",
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "provenance",
				Measures: "o inventário de componentes entregues (SBOM) é gerado e versionável",
			},
		)
	}

	// O guide, se houver, deve destilar seus pontos de conformidade (determinístico).
	if chosen["guide"] {
		gates = append(gates, config.Gate{
			Name: "guide-checklist", ID: "guide-checklist", On: []string{"guide"}, Check: "guide-has-checklist",
			Blocking: config.Bool(false), Measures: "o guide tem a seção de pontos de conformidade (CKn)",
		})
	}
	if projetoNovo {
		for i := range gates {
			if !dependemDeSinalIngerido[gates[i].Name] {
				gates[i].Blocking = config.Bool(true)
			}
		}
	}
	return gates
}

// CanonicalGate devolve a declaração canônica de um gate pelo nome, se existir.
//
// Mesma fonte do `init` — a lista acima, com TODOS os artefatos ligados — para que a
// redação de um gate exista em UM lugar só. Um catálogo escrito à parte divergiria da
// semente no primeiro ajuste, que é exatamente o problema que este acessor resolve.
func CanonicalGate(name string) (config.Gate, bool) {
	for _, g := range canonicalCatalog() {
		if g.Name == name {
			return g, true
		}
	}
	return config.Gate{}, false
}

// canonicalCatalog materializa a lista completa (todos os artefatos). Os gates são
// nomeados de forma única, então ligar tudo não gera colisão — só a união.
func canonicalCatalog() []config.Gate {
	// `false`: aqui só importa QUAIS gates existem, não o estado de maturação deles.
	return DefaultGates(map[string]bool{
		"spec": true, "feature": true, "test": true,
		"code": true, "guide": true, "plan": true,
	}, false)
}

// init registra o catálogo canônico no pacote `config`, que o usa para completar os
// campos omitidos de um gate canônico declarado num anchors.yaml (ver mergeCanonical).
//
// A injeção é no `init` do pacote — e não numa chamada explícita de cada comando —
// porque o merge tem de valer para TODO `config.Load`, inclusive os que rodam antes de
// qualquer setup. Um comando que esquecesse de registrar leria a config sem o merge e se
// comportaria diferente dos outros, que é o tipo de divergência silenciosa que o merge
// existe para eliminar.
func init() { config.SetCanonicalGateResolver(CanonicalGate) }
