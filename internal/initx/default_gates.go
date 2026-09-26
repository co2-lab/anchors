package initx

import "github.com/co2-lab/anchors/internal/config"

// dependOnIngestedSignal são os gates que só têm o que medir depois de `anchors ingest`
// receber um relatório de teste, cobertura ou mutação. Bloquear com base num sinal que
// ainda não existe barraria o commit por ausência de dado — não por defeito.
//
// Ficam informativos mesmo em projeto novo, e o usuário os promove quando a suíte
// estiver rodando no CI.
var dependOnIngestedSignal = map[string]bool{
	"tests-green": true, "line-coverage": true, "coverage-delta": true,
	"mutation-score": true, "scenario-coverage": true, "sbom-generated": true,
	"dependency-vulnerable": true, "no-duplication": true,
	"license-compatible": true, "circular": true, "deadcode": true, "spellcheck": true,
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
				Name: "spec-complete", ID: "spec-complete", On: []string{"spec"}, Check: "spec-sections",
				Blocking: config.Bool(false), Measures: "the spec has at least one state/rule, with no placeholder",
			},
			config.Gate{
				Name: "spec-has-code", ID: "spec-has-code", On: []string{"spec"}, Check: "has-code",
				Blocking: config.Bool(false), Measures: "the spec carries a scenario code (identity)",
			},
			// A spec sozinha atravessa TODOS os gates relacionais — eles falham ABERTO
			// (sem teste ligado, não há o que confrontar) e o pipeline conclui "pode
			// promover" sobre trabalho que não existe. Este gate pergunta o oposto: as
			// peças EXISTEM? Nasce informativo porque quase todo projeto tem débito.
			config.Gate{
				Name: "triad-complete", ID: "triad-complete", On: []string{"spec"}, Check: "triad-complete",
				Blocking: config.Bool(false), Measures: "the spec has code, feature and test that realise it",
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
				Name: "no-test-proof-real", ID: "no-test-proof-real", On: []string{"spec"},
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
				Name: "rule-fulfilled", ID: "rule-fulfilled", On: []string{"spec"},
				Blocking: config.Bool(false), Measures: config.MeasuresJudgment,
				Ask: "Cada regra desta spec está marcada no código por um comentário com o " +
					"código dela (`// ABCDX-B01: …`). Leia a regra e o trecho que ela marca e " +
					"responda: o trecho REALIZA o que a regra descreve? Reprove quando o " +
					"comentário estiver sobre código que faz outra coisa, quando a regra " +
					"descrever um caso que o trecho não trata, ou quando a marcação estiver " +
					"num lugar genérico (topo do arquivo, import) em vez do trecho que decide. " +
					"A pergunta é sobre o que o código EXECUTA, não sobre o que o comentário " +
					"afirma. Ao reprovar, PROPONHA a correção como patch. " +
					tbdInstruction("o código"),
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
				Name: "code-cataloged", ID: "code-cataloged", On: []string{"spec"}, Check: "code-cataloged",
				Blocking: config.Bool(false),
				Measures: "every exported symbol has a rule in the spec or a written waiver",
			},
			config.Gate{
				Name: "rule-types", ID: "rule-types", On: []string{"spec"}, Check: "rule-types",
				Blocking: config.Bool(false), Measures: "every code letter is declared in the vocabulary, with no conflict",
			},
			// A Tabela de Dependências promete símbolos; o código precisa usá-los.
			// Pega a divergência que o feature-test-match não vê (a aresta spec→código).
			config.Gate{
				Name: "dependency-honored", ID: "dependency-honored", On: []string{"spec"}, Check: "dependency-honored",
				Blocking: config.Bool(false), Measures: "the methods declared in the Dependencies Table are used in the code",
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
				Name: "proof-crosses-boundary", ID: "proof-crosses-boundary", On: []string{"spec"}, Check: "proof-crosses-boundary",
				Blocking: config.Bool(false), Measures: "a rule asserting a relation with another unit has the code importing that unit",
			},
			config.Gate{
				Name: "contract-status-declared", ID: "contract-status-declared", On: []string{"spec"}, Check: "contract-status-declared",
				Blocking: config.Bool(false), Measures: "the Output Contract statuses are the ones the handler returns — and only those",
			},
			// O dever que um artefato contrai com um lugar que ele NÃO conhece (LGPD,
			// i18n, a11y, auditoria). Só age se o projeto declarar `obligations:` — sem
			// declaração, o gate é inerte. Nasce ligado para que a categoria exista.
			config.Gate{
				Name: "obligation-honored", ID: "obligation-honored", On: []string{"spec"}, Check: "obligation-honored",
				Blocking: config.Bool(false), Measures: "the artifact fulfils the cross-cutting obligations it takes on",
			},
		)
	}
	// FUNÇÕES IRMÃS: quando a maioria guarda um parâmetro e uma não, a que não guarda é
	// quase sempre esquecimento. Detectável sem entender o domínio — é assimetria.
	if chosen["test"] || chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "sibling-guard", ID: "sibling-guard", On: []string{"code"}, Check: "sibling-guard",
			Blocking: config.Bool(false), Measures: "sibling functions handle the same parameter consistently",
		})
	}
	// PLANO: um plano semeia specs, e semear numa camada que não tem spec faz quem
	// executa gastar uma rodada descobrindo a contradição. A origem é checável.
	if chosen["plan"] {
		gates = append(gates, config.Gate{
			Name: "plan-seeds-valid", ID: "plan-seeds-valid", On: []string{"plan"}, Check: "plan-seeds-valid",
			Blocking: config.Bool(false), Measures: "the plan only seeds specs in a layer that has specs",
		})
		// A FONTE que o plano nomeia tem de ter dono declarado.
		//
		// Medido no blue-eyes: o plano 0008 dizia "Fonte: **GA4**" e declarava
		// `needs:` só do 0005. O adaptador vinha do 0002, e a dependência existia SÓ
		// NA PROSA — até o 0002 remover o adaptador numa revisão. Os dois planos
		// seguiram internamente coerentes, e a contradição só apareceu ao começar o
		// 0008. Nenhum gate pegava: o `dependency-honored` confronta o `needs:`
		// declarado, e aqui o defeito é o `needs:` que FALTA.
		gates = append(gates, config.Gate{
			Name: "plan-source-declared", ID: "plan-source-declared", On: []string{"plan"},
			Check: "plan-source-declared", Blocking: config.Bool(false),
			Measures: "the source the plan names has the adapter's plan in `needs:`",
		})
		// A DOC COMPILADA envelhece em silêncio. O conteúdo mora na spec; o `docs/*.md`
		// é derivado dela por template, e quem altera uma regra e esquece de recompilar
		// deixa a documentação afirmando a versão ANTIGA — com conteúdo real, e por isso
		// convincente. Uma doc obviamente incompleta manda procurar a fonte; uma doc
		// desatualizada não manda procurar nada.
		//
		// INFORMATIVO por decisão: o desvio é resolvido por um `anchors docs build`, e o
		// pipeline roda o build no merge. Reprovar o autor por trabalho que a máquina faz
		// sozinha gastaria a atenção da revisão com o que ela menos precisa.
		gates = append(gates, config.Gate{
			Name: "docs-fresh", ID: "docs-fresh", On: []string{"spec"},
			Check: "docs-fresh", Blocking: config.Bool(false),
			Measures: "the compiled `docs/*.md` reflects the spec it came from",
		})
		// O EIXO VERTICAL — a doutrina de produto e quem a realiza.
		//
		// Quatro perguntas distintas, cada uma pegando um silencio que as outras nao
		// veem. As tres de referencia reprovam (a referencia quebrada parece
		// rastreabilidade e nao resolve para nada); a de realizacao INFORMA, porque uma
		// regra decidida hoje para ser implementada no proximo ciclo e' trabalho
		// legitimo, nao defeito.
		gates = append(gates, config.Gate{
			Name: "plan-doctrine-exists", ID: "plan-doctrine-exists", On: []string{"plan"},
			Check: "plan-doctrine-exists", Blocking: config.Bool(false),
			Measures: "the product doctrine the plan seeds exists",
		})
		gates = append(gates, config.Gate{
			Name: "doctrine-realized", ID: "doctrine-realized", On: []string{"product"},
			Check: "doctrine-realized", Blocking: config.Bool(false),
			Measures: "every product rule is realised by some spec",
		})
		gates = append(gates, config.Gate{
			Name: "spec-doctrine-exists", ID: "spec-doctrine-exists", On: []string{"spec"},
			Check: "spec-doctrine-exists", Blocking: config.Bool(false),
			Measures: "the doctrine the spec references with `@realizes` exists",
		})

		// O EIXO DAS FEATURE FLAGS. Uma flag multiplica os caminhos do codigo sem
		// multiplicar a spec, e os tres custos (revisao, teste, remocao) sao silencios que
		// nenhum outro gate enxerga.
		// THE WAY BACK OF EACH PAIR OF THE TRIAD.
		//
		// `spec-feature-match` and `feature-test-match` walk the ORIGIN looking for the
		// destination. These walk the destination asking whether the origin still exists —
		// and that is what catches what a revert leaves behind: a scenario with no rule, a
		// test proving a reverted rule. Measured: the real case passed with an EMPTY message.
		//
		// Informational at birth: every historical revert left leftovers, and a large batch
		// of identical findings becomes noise people learn to scroll past.
		gates = append(gates, config.Gate{
			Name: "feature-spec-match", ID: "feature-spec-match", On: []string{"feature"},
			Check: "feature-spec-match", Blocking: config.Bool(false),
			Measures: "every scenario of the feature matches a rule the spec declares",
		})
		gates = append(gates, config.Gate{
			Name: "test-feature-match", ID: "test-feature-match", On: []string{"test"},
			Check: "test-feature-match", Blocking: config.Bool(false),
			Measures: "every code the test proves matches a declared scenario",
		})

		// A REVISAO QUE MUDOU O SIGNIFICADO DE UMA PALAVRA, e nao disse a quem.
		//
		// A BLOCKING-class gate (RVORP-B08, decided by the user): a new project is born with
		// it blocking. `false` here is the maturation state of an EXISTING project, like
		// every structural gate: its old revisions carry no `Checked:`, so each is born a
		// finding, and the project measures that queue before promoting the gate.
		gates = append(gates, config.Gate{
			Name: "revision-orphans", ID: "revision-orphans", On: []string{"spec"},
			Check: "revision-orphans", Blocking: config.Bool(false),
			Measures: "the revision names the sibling rules that speak of the same subject",
		})
		gates = append(gates, config.Gate{
			Name: "flag-scenario-grammar", ID: "flag-scenario-grammar", On: []string{"flag"},
			Check: "flag-scenario-grammar", Blocking: config.Bool(true),
			Measures: "each scenario's condition is written in the grammar",
		})
		gates = append(gates, config.Gate{
			Name: "flag-scenarios-complete", ID: "flag-scenarios-complete", On: []string{"flag"},
			Check: "flag-scenarios-complete", Blocking: config.Bool(true),
			Measures: "the flag declares the ABSENT case",
		})
		gates = append(gates, config.Gate{
			Name: "flag-scenario-exists", ID: "flag-scenario-exists", On: []string{"spec"},
			Check: "flag-scenario-exists", Blocking: config.Bool(true),
			Measures: "the scenario the spec cites with `@gated-by` exists",
		})
		// The WAY BACK of the previous one: that one confronts the spec citing a scenario
		// that does not exist; this one, the scenario nobody cites. Informational because
		// the flag is born BEFORE the specs that cite it — demanding both ends in the same
		// commit is not how the work happens.
		gates = append(gates, config.Gate{
			Name: "flag-scenario-governs", ID: "flag-scenario-governs", On: []string{"flag"},
			Check: "flag-scenario-governs", Blocking: config.Bool(false),
			Measures: "every flag scenario governs at least one rule",
		})
		// INFORMATIVO: teste por cenario e' a regra certa e a mais cara de cumprir. O
		// cenario escrito hoje e testado no commit seguinte e' trabalho normal, nao defeito.
		gates = append(gates, config.Gate{
			Name: "flag-covered", ID: "flag-covered", On: []string{"flag"},
			Check: "flag-covered", Blocking: config.Bool(false),
			Measures: "each flag scenario has at least one green test",
		})
		gates = append(gates, config.Gate{
			Name: "doctrine-not-duplicated", ID: "doctrine-not-duplicated", On: []string{"spec"},
			Check: "doctrine-not-duplicated", Blocking: config.Bool(false),
			Measures: "the spec references the doctrine instead of copying its text",
		})

		// A EXIGENCIA POR CAMADA: opcional por padrao, e quem decide e' a Estrutura.
		// Sem camada declarando `requires_doctrine`, este gate e' Skip em tudo — o que
		// e' o certo para um framework, e o oposto do certo para um produto.
		gates = append(gates, config.Gate{
			Name: "spec-realizes-doctrine", ID: "spec-realizes-doctrine", On: []string{"spec"},
			Check: "spec-realizes-doctrine", Blocking: config.Bool(false),
			Measures: "the spec's rules declare the product doctrine they realise",
		})

		// A FALHA declarada, tratada e registrada — a camada estatica do conceito.
		//
		// A spec ja' catalogava como a unidade falha (a secao `-E`), e NADA confrontava
		// isso: as regras atravessavam o pipeline inteiro sem que se perguntasse se eram
		// tratadas, se logavam, ou se aconteciam.
		//
		// `failure-logged` e' o que sustenta as camadas seguintes: um tratamento que
		// engole a falha sem registrar nada e' o silencio perfeito — ela acontece, nada
		// sabe, e nenhuma ferramenta a jusante tem o que ler.
		gates = append(gates, config.Gate{
			Name: "failure-handled", ID: "failure-handled", On: []string{"spec"},
			Check: "failure-handled", Blocking: config.Bool(false),
			Measures: "the declared failure has a path that handles it in the code",
		})
		gates = append(gates, config.Gate{
			Name: "failure-logged", ID: "failure-logged", On: []string{"spec"},
			Check: "failure-logged", Blocking: config.Bool(false),
			Measures: "the failure handling RECORDS the occurrence",
		})
		gates = append(gates, config.Gate{
			Name: "failure-declared", ID: "failure-declared", On: []string{"spec"},
			Check: "failure-declared", Blocking: config.Bool(false),
			Measures: "the handling that exists in the code answers some declared failure",
		})

		// A COBERTURA, que o gate acima nao tem como ver: ele confere as paginas que
		// EXISTEM, e o defeito silencioso e' a spec que nao chega a pagina nenhuma.
		// Os templates filtram por camada, entao uma spec que nenhum filtro seleciona
		// compila para lugar nenhum — e todas as paginas seguem corretas.
		gates = append(gates, config.Gate{
			Name: "docs-covered", ID: "docs-covered", On: []string{"spec"},
			Check: "docs-covered", Blocking: config.Bool(false),
			Measures: "every spec reaches some documentation page",
		})
		// A DOCUMENTAÇÃO AGREGADA que a unidade alimenta. O `docs.required` declara qual
		// documento é contrato e qual camada ou unidade o dispara; até aqui a declaração
		// era resolvida (`RequiredFor`) e nunca CONFRONTADA.
		//
		// Medido no projeto de referência: um PR entregou uma lambda nova — spec, código,
		// feature e teste — sem tocar o OpenAPI nem o esquema de dados, os dois declarados
		// com `trigger: [lambdas]`. Todos os checks ficaram verdes. Só apareceu porque
		// dois agentes colidiram no mesmo card e havia com que comparar.
		//
		// BLOQUEANTE, diferente do `docs-fresh` ao lado, e a diferença é o que a máquina
		// consegue fazer sozinha: um `docs/*.md` desatualizado se conserta com um
		// `docs build`, e reprovar o autor por isso gasta a atenção da revisão. Um
		// contrato que não descreve a rota nova NÃO se conserta sozinho — alguém precisa
		// escrever o que a rota faz.
		// ESCOPO `batch`: um veredito por DOCUMENTO, e não por unidade.
		//
		// A primeira versão rodava por nó, e cada spec cujo contrato não a mencionava
		// virava um card. No projeto de referência isso produziu 37 cards de uma vez — e
		// todos escreviam nos MESMOS dois arquivos.
		//
		// O efeito foi uma fila que não anda: cada PR mergeado invalidava os outros,
		// porque o ponto de inserção mudava. Medido: de 6 PRs, 2 passavam e 4 conflitavam.
		//
		// O defeito era do gate, que fatiou por unidade um trabalho que é por documento.
		gates = append(gates, config.Gate{
			Name: "doc-required", ID: "doc-required", On: []string{"spec"},
			Check: "doc-required", Blocking: config.Bool(true),
			Scope:    config.ScopeBatch,
			Measures: "the contracted documents exist and mention the units that trigger them",
		})
		// A SPEC BASTA POR SI. O corpo dela vira documentação palavra por palavra, e
		// quem lê o `docs/` não tem o repositório aberto: uma frase que só APONTA para um
		// plano manda essa pessoa a um arquivo que ela não vai abrir — que é exatamente o
		// que o mecanismo de documentação existe para eliminar.
		//
		// A referência COM o trecho citado junto não é acusada: ali o leitor tem o
		// argumento em mãos. O que o gate marca é o andaime que anuncia e não entrega.
		//
		// Casa ESTRUTURA, não vocabulário: os caminhos que o mapa conhece e a forma
		// `{CODIGO}-R000N`. Um gate que procurasse "plano" ou "ver" passaria em silêncio
		// no projeto escrito noutra língua — e silêncio é pior que ausência, porque a
		// spec pareceria protegida.
		gates = append(gates, config.Gate{
			Name: "doc-self-contained", ID: "doc-self-contained", On: []string{"spec"},
			Check: "doc-self-contained", Blocking: config.Bool(false),
			Measures: "the spec carries the text it cites, instead of pointing at another file",
		})
		// A ORDEM dentro do plano. Um plano sem fases catalogadas passa (elas são
		// opcionais); o gate só cobra a coerência de quem as declarou.
		// O PLANO REVISADO avisa quem o lê. Sem isso, quem abre um plano antigo segue uma
		// decisão que foi revista — e o plano parece coerente, porque ele É o registro
		// coerente do que se decidiu na época.
		gates = append(gates, config.Gate{
			Name: "plan-revised", ID: "plan-revised", On: []string{"plan"},
			Check: "plan-revised", Blocking: config.Bool(projetoNovo),
			Measures: "a plan revised by another tells whoever reads it",
		})
		// O PLANO ALTERADO diz por que mudou. Quem implementa é quem descobre o erro do
		// plano, e corrigi-lo em silêncio faz o projeto caminhar para um destino que
		// ninguém escolheu — deriva que nenhum gate de ESTADO vê, porque o plano
		// corrigido fica válido.
		//
		// `skip_on: [all]` não é detalhe: em `--all` não existe "alterado", e cobrar
		// revisão de todo plano acusaria quem acertou de primeira.
		gates = append(gates, config.Gate{
			Name: "plan-change-justified", ID: "plan-change-justified",
			On: []string{"plan", "spec"}, Check: "plan-change-justified",
			Blocking: config.Bool(projetoNovo), SkipOn: []string{"all"},
			Measures: "the changed plan/spec records the revision that says why it changed",
		})
		gates = append(gates, config.Gate{
			Name: "phase-ordered", ID: "phase-ordered", On: []string{"plan"}, Check: "phase-ordered",
			Blocking: config.Bool(projetoNovo),
			Measures: "the plan's phases do not depend on what comes after them",
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
			Name: "parent-valid", ID: "parent-valid", On: onde,
			Check: "parent-valid", Blocking: config.Bool(projetoNovo),
			Measures: "the declared `parent:` points at something that exists, with no cycle",
		})
	}
	if chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "phase-exists", ID: "phase-exists", On: []string{"spec"}, Check: "phase-exists",
			Blocking: config.Bool(projetoNovo),
			Measures: "the spec's `needs:` points at a phase some plan catalogues",
		})
	}
	if chosen["feature"] {
		gates = append(gates, config.Gate{
			Name: "feature-not-empty", ID: "feature-not-empty", On: []string{"feature"}, Check: "non-empty",
			Blocking: config.Bool(false), Measures: "the feature is not an empty skeleton",
		})
		// scenario-asserts vai além do 'não está vazio': o passo de RESULTADO precisa
		// afirmar um resultado. "Então o efeito XXXXX-B01 se verifica" satisfaz todos os
		// outros gates e não diz nada — a definição do que 'se verifica' migra para o
		// teste, que passa a ser escrito sem saber qual era o desfecho. Medido: 151 de 442
		// features de um projeto real, e foi a razão estrutural pela qual uma regra ficou
		// sem caso discriminante.
		gates = append(gates, config.Gate{
			Name: "scenario-asserts", ID: "scenario-asserts", On: []string{"feature"}, Check: "scenario-asserts",
			Blocking: config.Bool(false), Measures: "the outcome step asserts an observable outcome",
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
				Name: "test-traceable", ID: "test-traceable", On: []string{"test"}, Check: "test-traceable",
				Blocking: config.Bool(false),
				Measures: "the test cites the code of what it proves (visible to the relational gates)",
			},
			config.Gate{
				Name: "tests-green", ID: "tests-green", On: []string{"test"}, Check: "tests-pass",
				Blocking: config.Bool(false), Measures: "this file's tests pass (from the ingested result)",
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
				Name: "mock-detect-covers-dialect", ID: "mock-detect-covers-dialect", On: []string{"test"},
				Blocking: config.Bool(false), Measures: config.MeasuresJudgment,
				Ask: "O projeto declara `derived.mock_detect` — o regex que reconhece um " +
					"dublê de teste neste ecossistema. Leia este arquivo de teste e responda: " +
					"o padrão declarado alcança TODAS as formas de dublê que ele usa? " +
					"Reprove se houver dublê que o regex não casa (outra função, outro " +
					"dialeto, decorator em vez de chamada, dublê de módulo inteiro vs de " +
					"membro). Um regex que casa zero faz o gate `mock-carimbado` reportar " +
					"verde sem ter conferido nada. Ao reprovar, PROPONHA o padrão corrigido " +
					"como patch do `anchors.yaml`." +
					tbdInstruction("o teste"),
			},
			config.Gate{
				Name: "mock-stamped", ID: "mock-stamped", On: []string{"test"}, Check: "mock-stamped",
				Blocking: config.Bool(false),
				Measures: "the double's stamp matches the real snippet (recomputed, not just validated)",
			},
			config.Gate{
				Name: "mock-typed", ID: "mock-typed", On: []string{"test"}, Check: "mock-typed",
				Blocking: config.Bool(false),
				Measures: "the test double derives from the real module (not a frozen copy of the contract)",
			},
			config.Gate{
				Name: "line-coverage", ID: "line-coverage", On: []string{"code"}, Check: "line-coverage",
				Blocking: config.Bool(false), Measures: "line coverage >= threshold (from the ingested lcov)",
			},
			config.Gate{
				Name: "coverage-delta", ID: "coverage-delta", On: []string{"code"}, Check: "coverage-delta",
				Blocking: config.Bool(false), Measures: "line coverage did not drop vs. the previous ingestion",
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
				Blocking: config.Bool(false), Measures: "the test PROVES the line: killed mutants >= threshold (from the ingested mutation report)",
			},
		)
		// pagination-honored é da mesma família: pega o que NENHUM teste pega, porque o
		// teste roda com 3 registros e a produção com 3 mil. Medido no app de referência: 11 arquivos,
		// entre eles uma listagem de importação de extrato que perde silenciosamente as
		// transações além do 1º MB — o usuário importa o arquivo e parte some sem erro.
		gates = append(gates, config.Gate{
			Name: "pagination-honored", ID: "pagination-honored", On: []string{"code"}, Check: "pagination-honored",
			Blocking: config.Bool(false), Measures: "a function promising the whole set does not silently return the first page",
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
			Blocking: config.Bool(false), Measures: "the layer does not reach what is not its own (`boundaries:`)",
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
			Blocking: config.Bool(false), Measures: "the number the spec asserts matches the code",
		})

		gates = append(gates, config.Gate{
			Name: "domain-declared", ID: "domain-declared", On: []string{"spec"}, Check: "domain-declared",
			Blocking: config.Bool(false), Measures: "the spec declares what it accepts and who guards the boundary",
		})

		gates = append(gates, config.Gate{
			Name: "open-questions-resolved", ID: "open-questions-resolved", On: []string{"spec"}, Check: "open-questions-resolved",
			Blocking: config.Bool(false), Measures: "the spec has no open question — implementing is not guessing",
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
	// ref-resolves: o `ref:` aponta para a spec CERTA. `header-valid` confere que o
	// campo existe; ninguém conferia que ele RESOLVE. Um ref errado parece rastreabilidade
	// e atribui a unidade à spec errada, com todo gate relacional confrontando o par
	// errado. Medido: 49 arquivos de modelo apontando para a identidade de antes de uma
	// desfusão que ninguém propagou, e mais um teste citando duas specs fantasma.
	// Condicionado a `spec`: sem spec no projeto não há `code:` para o `ref:` resolver.
	if chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "ref-resolves", ID: "ref-resolves", On: []string{"code", "feature", "test"}, Check: "ref-resolves",
			Blocking: config.Bool(false), Measures: "the `ref:` points at the sibling spec's `code:`",
		})
	}

	if chosen["spec"] {
		gates = append(gates, config.Gate{
			Name: "code-reference-valid", ID: "code-reference-valid", On: []string{"spec"}, Check: "code-reference-valid",
			Blocking: config.Bool(false), Measures: "every code the spec cites exists in the project",
		})
	}

	if chosen["spec"] && chosen["feature"] {
		gates = append(gates, config.Gate{
			Name: "spec-feature-match", ID: "spec-feature-match", On: []string{"spec"}, Check: "spec-feature-match",
			Blocking: config.Bool(false), Measures: "every requirement declared in the spec has a scenario in the feature",
		})
		if chosen["spec"] {
			gates = append(gates, config.Gate{
				Name: "scenario-coverage", ID: "scenario-coverage", On: []string{"spec"}, Check: "scenario-coverage",
				Blocking: config.Bool(false), Measures: "each scenario of the spec has a passing test (semantic coverage)",
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
				Name: "no-secret-leaked", ID: "no-secret-leaked", On: []string{"code"},
				Scope: config.ScopeBatch, ScopeFull: config.ScopeProject,
				Run:       "gitleaks git --no-banner --redact -v",
				NeedsTool: "gitleaks", InstallHint: "brew install gitleaks",
				Blocking: config.Bool(true), When: []string{"pre-commit", "ci"}, Cost: "fast",
				Category: "security",
				Measures: "no secret (key, token, credential) enters the history",
			},
			// Informativo: a CVE nova aparece sem ninguém mexer no código, então bloquear
			// pararia um merge por algo que o autor não causou nem pode resolver na hora.
			config.Gate{
				Name: "dependency-vulnerable", ID: "dependency-vulnerable", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Run:       "osv-scanner scan source -r .",
				NeedsTool: "osv-scanner", InstallHint: "brew install osv-scanner",
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "security",
				Measures: "the project's dependencies have no known vulnerability (OSV)",
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
				Name: "no-duplication", ID: "no-duplication", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Run:       "npx --yes jscpd . --reporters console --silent",
				NeedsTool: "npx", InstallHint: "instale Node.js (npx acompanha)",
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "quality",
				Measures: "no code block appears copied in two places",
			},
			// Procedência: responde "o que exatamente foi entregue" — a pergunta que só
			// se faz depois do incidente, quando reconstruir a resposta já é impossível.
			config.Gate{
				Name: "sbom-generated", ID: "sbom-generated", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Run:       "syft scan dir:. -o cyclonedx-json=sbom.json -q",
				NeedsTool: "syft", InstallHint: "brew install syft",
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "provenance",
				Measures: "the inventory of shipped components (SBOM) is generated and versionable",
			},
			// Ortografia: erros de grafia no código, testes e documentação minam a
			// confiança e quebram buscas. Agnóstico: typos (binário nativo ultra-rápido)
			// ou cspell (Node). Ver `guides/GATES_ECOSYSTEM_GUIDE.md`.
			config.Gate{
				Name: "spellcheck", ID: "spellcheck", On: []string{"code", "test", "spec"},
				Scope: config.ScopeBatch, ScopeFull: config.ScopeProject,
				Run:       "typos",
				NeedsTool: "typos", InstallHint: "brew install typos",
				Blocking: config.Bool(false), When: []string{"pre-commit", "ci"}, Cost: "fast",
				Category: "style",
				Measures: "no spelling mistake in the text or the identifiers",
			},
			// Conformidade de licenças: garante que nenhuma dependência de produção traga
			// licenças com copyleft forte (AGPL, SSPL, etc.) incompatíveis com o projeto.
			// O comando varia conforme o gerenciador de pacotes do ecossistema.
			config.Gate{
				Name: "license-compatible", ID: "license-compatible", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "fast",
				Category: "legal",
				Measures: "no dependency with strong copyleft or an incompatible licence",
			},
			// Dependência circular: pergunta de arquitetura de projeto (o ciclo é do grafo,
			// não de um arquivo isolado). Ferramenta varia por dialeto (madge, go vet, etc).
			config.Gate{
				Name: "circular", ID: "circular", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Blocking: config.Bool(false), When: []string{"pre-push", "ci"}, Cost: "slow",
				Category: "architecture",
				Measures: "there is no import cycle between modules",
			},
			// Código morto: exportações, símbolos e arquivos órfãos sem consumidor.
			// Ferramenta varia por dialeto (knip, deadcode, vulture, cargo-udeps).
			config.Gate{
				Name: "deadcode", ID: "deadcode", On: []string{"code"},
				Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
				Blocking: config.Bool(false), When: []string{"ci"}, Cost: "slow",
				Category: "maintenance",
				Measures: "there is no orphan export, symbol or file in the project",
			},
		)
	}

	// O guide, se houver, deve destilar seus pontos de conformidade (determinístico).
	if chosen["guide"] {
		gates = append(gates, config.Gate{
			Name: "guide-checklist", ID: "guide-checklist", On: []string{"guide"}, Check: "guide-has-checklist",
			Blocking: config.Bool(false), Measures: "the guide has the conformance-points section (CKn)",
		})
	}
	if projetoNovo {
		for i := range gates {
			if !dependOnIngestedSignal[gates[i].Name] {
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

// init liga a lista de nomes de gate ao pacote `config`, para o teste que confronta o
// de-para do vocabulário antigo contra os gates que existem de verdade.
//
// Por injeção e não por import: `config` não pode importar `initx` (seria ciclo), e sem
// a lista o teste do de-para não teria contra o que confrontar — um destino errado
// passaria despercebido.
func init() {
	config.RegisterGateNames(func() []string {
		todos := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, false)
		out := make([]string, 0, len(todos))
		for _, g := range todos {
			out = append(out, g.Name)
		}
		return out
	})
}
