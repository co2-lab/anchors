package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/change"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// `anchors work <artefato> --for <alvo>` emite o PROMPT DE TRABALHO de uma etapa.
//
// Por que existe: o `anchors guide` ensina a DOUTRINA (o que é uma spec, como se escreve
// um teste) — é a régua, permanente e agnóstica de alvo. Mas quem vai executar UMA etapa
// sobre UM arquivo precisa de outra coisa: o que ler AGORA, nesta ordem; qual é a camada
// deste alvo e o que ela exige; onde nascem as peças da trinca; o que NÃO é escopo desta
// etapa; e como verificar no fim.
//
// Sem isso, cada orquestrador reescreve esse prompt à mão a cada vez — e sai diferente
// toda vez (o observado: um agente reinventou o molde de `.feature` de backend lendo os
// vizinhos, embora o guide do projeto tivesse a seção certa; outro gastou uma rodada
// descobrindo que uma camada declarativa não tem spec).
//
// O comando NÃO inventa: ele COMPÕE o que já está declarado na Estrutura —
//   - `governs:` → quais guides regem a camada do alvo (a régua a ler)
//   - `layers:`  → o que é a camada, e se ela é regida ou reconhecida
//   - `derived:` → onde nascem spec/feature/teste desse alvo
//   - `work:`    → passos extra que só o projeto sabe (override opcional por camada)
func newWorkCmd() *cobra.Command {
	var root, target string

	cmd := &cobra.Command{
		Use:   "work <artefato>",
		Short: "Emite o prompt de trabalho de uma etapa (spec|code|feature|test|review) para um alvo",
		Long: `Compõe o PROMPT DE TRABALHO de uma etapa sobre um alvo concreto — o que um
agente (ou subagente) precisa para executar aquele pedaço sem reinventar o roteiro.

  anchors work spec   --for packages/backend/repositories/metadata.ts
  anchors work test   --for apps/mobile/src/hooks/useAuth.ts
  anchors work review --for packages/backend/business-logic/pricing.ts

O ` + "`review`" + ` é a última etapa do ciclo, e a diferente: não produz artefato —
CONFRONTA o que foi entregue. Existe porque o gate verifica o que é DECLARÁVEL, e o que
sobra só aparece para quem ataca de fora (mutando a regra, rodando com entrada de borda).
Rode-o depois de ` + "`anchors deliver`" + `, que é o que lhe dá escopo.

Diferença para o ` + "`anchors guide`" + `: o guide ensina a DOUTRINA (permanente, sem alvo);
o work entrega a TAREFA (o que ler agora, qual a camada deste alvo, onde nascem as
peças, o que não é seu escopo, como verificar).

O conteúdo é COMPOSTO do anchors.yaml — nada é inventado aqui.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			artifact := strings.ToLower(args[0])
			if !queue.ArtefatoDeTrabalhoValido(artifact) {
				return fmt.Errorf("artefato desconhecido %q — use: %s", artifact,
					strings.Join(queue.ArtefatosDeTrabalho, ", "))
			}
			if target == "" {
				return fmt.Errorf("informe o alvo com --for <caminho> " +
					"(ex.: --for packages/backend/repositories/metadata.ts)")
			}
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar config: %w", err)
			}
			rel := relTo(absRoot, target)
			g, _ := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath)) // sem mapa: cai na convenção

			// O alvo é a UNIDADE (o arquivo de código), não uma peça derivada dela. Apontar
			// para a spec é o engano previsível — é o artefato que já existe, então é o que
			// vem à mão. O resultado era silenciosamente absurdo: os caminhos saíam
			// `x.spec.spec.md`, `x.spec.feature`, e a spec aparecia como peça a produzir de
			// uma spec que já estava no disco. Nada avisava.
			//
			// Redirecionar (em vez de recusar) porque a intenção é inequívoca: quem pediu
			// `work feature --for x.spec.md` quer a feature da unidade que x.spec.md
			// descreve. Dizer o que se fez mantém o usuário no controle.
			if alvo, achou := unidadeDaPecaDerivada(absRoot, rel, cfg, g); achou {
				fmt.Fprintf(os.Stderr, "nota: `%s` é uma peça derivada, não a unidade. "+
					"Usando `%s` como alvo.\n\n", rel, alvo)
				rel = alvo
			}
			out, err := composeWorkPrompt(absRoot, rel, artifact, cfg, g)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&target, "for", "", "o arquivo-alvo desta etapa (obrigatório)")
	return cmd
}

// composeWorkPrompt monta o prompt. A ordem das seções é deliberada: PAPEL (quem você
// é) → ALVO (sobre o quê) → RÉGUA (o que ler antes) → FRONTEIRA (o que não é seu) →
// PROCEDIMENTO (como) → VERIFICAÇÃO (como saber que acabou).
func composeWorkPrompt(root, rel, artifact string, cfg *config.Config, g *mapx.Graph) (string, error) {
	layer, _ := scan.Classify(rel, cfg)
	l, hasLayer := cfg.Layers[layer]

	var b strings.Builder
	// O REVIEW não produz artefato — ele CONFRONTA. Chamar de "trabalho de review" e
	// "você vai produzir review" empurra o revisor para o modo autor, que é o oposto do
	// que se quer: quem escreve tende a explicar o que vê, quem confronta tenta quebrar.
	if artifact == "review-plan" {
		fmt.Fprintf(&b, "# Review de CONJUNTO — `%s`\n\n", rel)
	} else if artifact == "review" {
		fmt.Fprintf(&b, "# Review de `%s`\n\n", rel)
	} else {
		fmt.Fprintf(&b, "# Trabalho: %s de `%s`\n\n", artifact, rel)
	}

	// PAPEL + a recusa antecipada. Uma camada RECONHECIDA não tem spec — dizer isso
	// AQUI evita a rodada de descoberta (e o arquivo errado nascendo).
	// A recusa vale para TODA peça da trinca, não só a spec. Antes cobria apenas `spec`, e
	// `anchors work feature --for <arquivo declarativo>` abria com "Você vai produzir
	// feature" para, duas seções abaixo, dizer "não tem spec, feature nem teste próprios"
	// — e prescrever a verificação de um arquivo que o mesmo prompt proíbe criar. Medido
	// num E2E real: a task teve de ser descartada à mão.
	// PARE também quando a PEÇA desta etapa é dispensada pela camada (`trinca_opcional`),
	// ainda que a camada seja REGIDA. Antes o prompt abria com "Você vai produzir
	// **feature**" e o roteiro completo de produção, e a dispensa aparecia quatro linhas
	// abaixo, como item de uma lista. Um worker que segue a manchete cria o arquivo
	// proibido — e a régua tinha dito as duas coisas.
	if hasLayer && pecasDispensadas(layer, cfg)[artifact] {
		fmt.Fprintf(&b, "## PARE\n\n`%s` — a camada **%s** DISPENSA a peça `%s` "+
			"(`trinca_opcional` no anchors.yaml).\n\nA dispensa é declarada, não um "+
			"esquecimento: esta camada não prova comportamento com esta peça. Criá-la "+
			"produziria o artefato vazio que a declaração existe para evitar.\n\n"+
			"Não crie a %s. Se a fila lhe deu esta task, ela é ruído — descarte-a "+
			"(`anchors drop`) e reporte.\n", rel, layer, artifact, artifact)
		return b.String(), nil
	}
	if hasLayer && l.Regime == "declarativo" && (artifact == "spec" || artifact == "feature" || artifact == "test") {
		peca := map[string]string{"spec": "spec", "feature": "feature", "test": "teste"}[artifact]
		fmt.Fprintf(&b, "## PARE\n\n`%s` pertence à camada **%s**, declarada como RECONHECIDA "+
			"(`regime: declarativo`) no anchors.yaml.\n\nCamadas assim **não têm spec, feature "+
			"nem teste próprios**: elas não originam regra (só traduzem/configuram), então não há "+
			"o que especificar nem que provar em cenário. Quem prova o comportamento é a camada "+
			"que DECIDE, com a trinca dela.\n\nNão crie a %s. Se há uma DECISÃO a documentar, "+
			"ela pertence à camada que decide — reporte a contradição a quem pediu esta etapa.\n",
			rel, layer, peca)
		return b.String(), nil
	}

	if artifact == "review-plan" {
		b.WriteString("Você é o REVISOR DO CONJUNTO. Cada unidade deste plano já passou por " +
			"review próprio, e cada uma está correta **sozinha** — não repita esse trabalho.\n\n" +
			"O seu alvo é a COSTURA: o que só aparece quando as peças se encontram. Medido " +
			"num E2E real: de 6 achados de um review adversarial, **4 atravessavam unidades**, " +
			"e em todos eles cada peça passava isolada. É a classe que nenhum gate pega e que " +
			"o review por unidade não alcança por definição de escopo.\n\n" +
			"Sua saída é um RELATÓRIO com evidência executada, não uma correção.\n\n")
	} else if artifact == "review" {
		b.WriteString("Você é o REVISOR desta unidade. O trabalho já foi entregue e **passou " +
			"por todos os gates** — é esse o estado normal de tudo que chega aqui, e é " +
			"exatamente por isso que você existe: o gate confronta o que é DECLARÁVEL, e o " +
			"que sobra só aparece para quem ataca de fora.\n\n" +
			"Sua saída é um RELATÓRIO com evidência executada, não uma correção.\n\n")
	} else {
		fmt.Fprintf(&b, "Você vai produzir **%s** para o alvo abaixo, seguindo a régua do projeto.\n\n", artifact)
	}
	fmt.Fprintf(&b, "## Alvo\n\n- Arquivo: `%s`\n", rel)
	if hasLayer {
		fmt.Fprintf(&b, "- Camada: **%s**", layer)
		if l.Regime != "" {
			fmt.Fprintf(&b, " (regime: %s)", l.Regime)
		}
		b.WriteString("\n")
		if len(l.Tags) > 0 {
			fmt.Fprintf(&b, "- Tags: %s\n", strings.Join(l.Tags, ", "))
		}
	} else {
		b.WriteString("- Camada: **não classificada** — confirme se o caminho pertence a alguma " +
			"camada do `layers:` no anchors.yaml antes de prosseguir.\n")
	}

	// RÉGUA: os guides que regem esta camada, lidos do `governs:`.
	if guides := guidesFor(l, cfg, artifact); len(guides) > 0 {
		b.WriteString("\n## Leia primeiro (nesta ordem)\n\n")
		fmt.Fprintf(&b, "1. `anchors guide %s` — a doutrina do artefato (o que é, o que não é)\n", artifact)
		for i, g := range guides {
			fmt.Fprintf(&b, "%d. `%s` — a régua deste projeto para esta camada\n", i+2, g)
		}
		fmt.Fprintf(&b, "%d. O alvo (`%s`) e seus vizinhos de camada, para seguir o dialeto local\n",
			len(guides)+2, rel)
	} else {
		b.WriteString("\n## Leia primeiro\n\n")
		fmt.Fprintf(&b, "1. `anchors guide %s` — a doutrina do artefato\n", artifact)
		fmt.Fprintf(&b, "2. O alvo (`%s`) e seus vizinhos de camada\n", rel)
		b.WriteString("\n> Nenhum guide do projeto rege esta camada (`governs:` no anchors.yaml). " +
			"Siga o dialeto dos vizinhos e registre a lacuna.\n")
	}

	// TRINCA: onde nascem as peças, derivado do `derived:`.
	//
	// Camada RECONHECIDA (`regime: declarativo`) não tem trinca: listar spec/feature/test
	// ali contradizia a própria linha de cima do prompt, que declara o regime. Um agente
	// que confiasse nesta seção criaria a spec que a camada proíbe — e o `anchors new`
	// depois a recusaria, sem que nada explicasse a contradição.
	if artifact == "review" || artifact == "review-plan" {
		// O revisor não vê "onde as peças nascem" — ele vê O QUE EXISTE, que é o material
		// do confronto. E vê o registro de entrega, que dá escopo e traz a intenção
		// declarada pelo autor para ser confrontada contra o disco.
		b.WriteString("\n## O que você vai confrontar\n\n")
		writeTrincaPaths(&b, rel, artifact, layer, cfg, g)
		writeDeliveryRecord(&b, root, rel)
	} else if hasLayer && l.Regime == "declarativo" {
		b.WriteString("\n## As peças e onde nascem\n\n")
		fmt.Fprintf(&b, "**Só o próprio arquivo.** `%s` está numa camada RECONHECIDA "+
			"(`regime: declarativo`): ela não origina regra — traduz, transporta ou "+
			"declara. Por isso **não tem spec, feature nem teste próprios**.\n\n", rel)
		b.WriteString("A regra que este arquivo serve mora na camada de quem DECIDE (a que " +
			"o consome). Se você sentir falta de uma spec aqui, o mais provável é que " +
			"a decisão esteja faltando lá — não que esta camada precise de uma.\n")
	} else {
		b.WriteString("\n## As peças e onde nascem\n\n")
		writeTrincaPaths(&b, rel, artifact, layer, cfg, g)
	}

	// REGIMES: as tags de nível que os cenários da feature DEVEM declarar. Estão no
	// anchors.yaml (`derived.regimes`) e nenhum comando as listava — um agente escreveu
	// `@nivel-integracao` (português, coerente com o resto do Gherkin) quando o certo
	// era `@nivel-integration`, e só descobriu garimpando o YAML.
	if artifact == "feature" || artifact == "test" {
		writeRegimes(&b, cfg)
	}

	// FRONTEIRA: o que NÃO é desta etapa. É o que o polvo chama de "belongs to the
	// layer agents" — sem isso o agente extrapola o escopo.
	b.WriteString("\n## O que NÃO é seu escopo\n\n")
	for _, s := range outOfScope(artifact) {
		fmt.Fprintf(&b, "- %s\n", s)
	}

	// PROCEDIMENTO: universal + override do projeto (híbrido).
	b.WriteString("\n## Procedimento\n\n")
	// Em camada RECONHECIDA o procedimento padrão não serve: ele manda "leia a spec
	// inteira; ela é a régua", e ali não existe spec. Mandar ler o que não existe deixa o
	// agente sem chão — ou pior, o convence a criar a spec proibida.
	passos := procedureFor(artifact, cfg)
	if hasLayer && l.Regime == "declarativo" {
		passos = procedureDeclarativa(rel)
	}
	for i, s := range passos {
		fmt.Fprintf(&b, "%d. %s\n", i+1, s)
	}
	if extra := l.Work[artifact]; len(extra) > 0 {
		b.WriteString("\n**Passos específicos desta camada** (declarados no anchors.yaml):\n\n")
		for _, s := range extra {
			fmt.Fprintf(&b, "- %s\n", s)
		}
	}

	// O QUE OS GATES VÃO COBRAR — dito ANTES, não depois.
	//
	// O relato foi este: "o guide diz 'rode `anchors check` para saber o formato'. Isso é
	// feedback pós-fato: eu escrevo, rodo, e só então descubro." Deu certo por copiar os
	// vizinhos, não porque a régua tenha dito. Um requisito que o autor só conhece depois
	// de reprovar é um requisito mal comunicado — o gate SABE o que exige, então dizer
	// antes custa nada e economiza uma rodada.
	if regras := exigenciasDosGates(artifact, cfg); len(regras) > 0 {
		b.WriteString("\n## O que os gates vão cobrar (leia ANTES de escrever)\n\n")
		for _, r := range regras {
			fmt.Fprintf(&b, "- %s\n", r)
		}
	}

	// REGRA NOVA / REGRA VIOLADA PEDE RÉGUA. O buraco simétrico ao do §5.1 do QUALITY: lá
	// se caça a âncora sem gate; aqui, a REGRA sem gate — um dever declarado em prosa que
	// nenhum medidor confronta. Enquanto vive só no guide, vale enquanto alguém lembra.
	//
	// O caso que originou isto: um projeto tinha o combinado "todo sheet usa a lib de
	// bottom-sheet, nunca o modal nativo". Não era gate. Uma varredura achou 20 usos do
	// modal nativo em 11 arquivos, invisíveis por meses — e o custo não era estético: as
	// duas famílias aninhadas produzem falha SILENCIOSA (o sheet não abre, sem erro).
	// Virou uma entrada declarativa de fronteira em minutos. A regra existia; a régua não.
	b.WriteString("\n## Se você ADICIONAR uma regra — ou PEGAR uma regra falhando\n\n")
	b.WriteString("Nos dois casos, pergunte: **isto cabe num gate?** Corrigir a violação é " +
		"metade do trabalho; a outra metade é a régua que impede a reincidência. Uma regra " +
		"que foi violada uma vez prova, por construção, que a memória humana não a sustenta.\n\n")
	b.WriteString("Ordem de preferência do medidor:\n\n")
	b.WriteString("1. **Determinístico declarativo** — a regra é reconhecível por padrão no " +
		"artefato (import proibido, campo obrigatório, nome fora de convenção): cabe numa " +
		"entrada de configuração do framework, sem código novo. Custo quase zero — tente isto primeiro.\n")
	b.WriteString("2. **Determinístico com código** — a regra exige atravessar estrutura " +
		"(resolver uma aresta, comparar dois artefatos): pede um gate próprio.\n")
	b.WriteString("3. **Julgamento por IA** — a regra é semântica e sem forma fixa " +
		"(\"o comentário descreve o que o código faz\").\n\n")
	b.WriteString("Só se as três não servirem a regra fica em prosa — e aí é declarada como " +
		"dívida explícita, não esquecida em silêncio. Ao terminar de corrigir uma violação: " +
		"*\"o que impede isto de voltar amanhã?\"*. Se a resposta é \"alguém vai lembrar\", " +
		"o trabalho não terminou. (Ver QUALITY.md §5.1, \"A regra sem régua\".)\n")

	// ISSUES ABERTAS DESTA UNIDADE. O ciclo escreve o achado do revisor numa issue e depois
	// nunca mais o lê: medido num E2E, os cinco achados mais valiosos de uma spec vieram de
	// uma issue que o próprio Anchors tinha escrito, e sem alguém trazê-la à mão a etapa
	// seguinte sairia com "✗0 bloqueantes" exatamente igual.
	//
	// É a falha de memória do framework: o defeito foi encontrado, registrado, e some do
	// caminho de quem poderia corrigi-lo. Trazê-lo para o prompt é o que fecha o laço entre
	// quem acha e quem conserta.
	writeIssuesAbertas(&b, root, rel)

	// SINAIS DE EXECUÇÃO. Os gates que medem o que só a EXECUÇÃO revela (o teste passa? o
	// teste PROVA a linha, ou só a executa?) leem de sinais ingeridos — e ficam `~` para
	// sempre se ninguém os ingerir. Medido num E2E real: `mutation-score` saiu `~828` no
	// repositório inteiro, e era exatamente o gate que teria pego 3 cenários decorativos
	// (o teste citava a restrição e não a provava). O instrumento existia, configurado, e
	// nunca tinha sido rodado.
	if artifact == "test" || artifact == "review" {
		b.WriteString("\n## Sinais de execução (o que o gate não vê sozinho)\n\n")
		b.WriteString("Gate verde não prova que o teste PROVA a regra — prova que ele roda. " +
			"Quem separa as duas coisas é a mutação: altere a linha e veja se o teste cai.\n\n")
		b.WriteString("```sh\n")
		b.WriteString("# 1) a suíte, e o resultado para o mapa\n")
		b.WriteString("#    (o gate `testes-passam` lê daqui: sem isto ele fica ~, e o gate que\n")
		b.WriteString("#     conta `it()` no fonte dá ✓ mesmo num arquivo que não COMPILA)\n")
		b.WriteString("anchors ingest --junit <relatório JUnit da suíte>\n\n")
		fmt.Fprintf(&b, "# 2) mutação DESTA unidade (o projeto inteiro é caro; escope no alvo)\n"+
			"#    rode o mutador do projeto sobre `%s`, depois:\nanchors ingest --mutation <relatório.json>\n```\n\n", rel)
		b.WriteString("> Um mutante SOBREVIVENTE é uma alteração no código que nenhum teste " +
			"percebeu: ali a regra não está provada. É o achado que os gates determinísticos " +
			"não alcançam.\n")
	}

	// COMO REGISTRAR A ENTREGA. Ausente de todo prompt até aqui: a seção "Como verificar
	// que terminou" acabava no `check`, e o `anchors deliver` — que é o que ENFILEIRA a
	// próxima etapa — não era citado em lugar nenhum. Medido em três execuções: o worker
	// terminava, nada acontecia, e o orquestrador tinha de descobrir o comando e suas
	// quatro flags obrigatórias por tentativa e erro (o CLI as revela uma por execução).
	//
	// Pior: o `deliver` repreende quem omite `--decision`/`--uncovered` ("você declarou
	// ZERO decisões livres… isso AFIRMA que a régua decidiu tudo") — a régua penalizando
	// uma omissão que ela própria induziu ao não pedir.
	if artifact != "review" && artifact != "review-plan" && artifact != "review-plan-draft" {
		b.WriteString("\n## Como REGISTRAR a entrega (é isto que enfileira a próxima etapa)\n\n")
		b.WriteString("```sh\n")
		fmt.Fprintf(&b, "anchors deliver --stage %s --unit %s \\\n", artifact, rel)
		b.WriteString("  --date <AAAA-MM-DD> \\\n")
		b.WriteString("  --intent \"o que você se propôs a fazer, em uma frase\" \\\n")
		b.WriteString("  --decision \"toda escolha que a régua NÃO decidiu por você\" \\\n")
		b.WriteString("  --uncovered \"o que ficou sem prova, e por quê\"\n```\n\n")
		b.WriteString("`--decision` e `--uncovered` aceitam a flag repetida (uma por item). " +
			"Declarar ZERO decisões AFIRMA que a régua decidiu tudo — se você escolheu " +
			"qualquer coisa sozinho (um nome, um limite, uma ordem), é uma decisão. É o que " +
			"o revisor vai confrontar contra o disco.\n")
	}

	// FECHAMENTO DO REVIEW. O review é a única etapa que não produz artefato — ela produz
	// ACHADOS —, e por isso era a única que terminava sem dizer o que fazer com o
	// resultado. Medido num E2E real: o prompt repetia três vezes "NÃO corrija — quem
	// acha não conserta" e não emitia comando de registro. O ciclo `deliver → review → ?`
	// não fechava: o achado ficava no relatório do subagente, e sumia com ele.
	//
	// Pior era o ACHADO CRUZADO: o revisor de uma unidade encontrou o bug crítico numa
	// unidade IRMÃ e, corretamente, não o corrigiu. Nenhuma task foi enfileirada, porque
	// nada no framework recebia um achado sobre outro alvo. Sem alguém ler o relatório à
	// mão e rotear, o bug seguiria verde e entregue.
	if artifact == "review" || artifact == "review-plan" || artifact == "review-plan-draft" {
		b.WriteString("\n## Como FECHAR o review (obrigatório)\n\n")
		b.WriteString("O review não produz artefato — produz ACHADOS, e achado que fica no seu " +
			"relatório morre com a sessão. Registre CADA um, no alvo a que ele pertence:\n\n")
		b.WriteString("```sh\n# nenhum achado:\n")
		fmt.Fprintf(&b, "anchors judge %s --gate review --verdict pass --reason \"o que você confrontou, e por que passou\"\n\n", rel)
		b.WriteString("# com achados (o --reason É o corpo da issue: laudo completo, não frase de veredito):\n")
		fmt.Fprintf(&b, "anchors judge %s --gate review --verdict fail --reason \"$(cat <<'EOF'\n", rel)
		b.WriteString("## Laudo\n### 1. <o quê>  (arquivo:linha)\n- **por quê:** a regra que isso quebra\n" +
			"- **como corrigir:** o caminho\nEOF\n)\"\n```\n\n")
		b.WriteString("**Achado numa unidade DIFERENTE da que você revisou?** Registre-o no alvo " +
			"dela, não no seu — `anchors judge <a-outra-unidade> --gate review --verdict fail`. " +
			"É o único jeito de o achado virar trabalho: quem revisou a unidade A não conserta " +
			"a unidade B, mas precisa deixar o defeito onde alguém o encontre.\n\n")
		b.WriteString("> Você continua **não corrigindo** — quem acha não conserta. Registrar não " +
			"é corrigir: é fazer o achado sobreviver a você.\n")
	}

	// VERIFICAÇÃO: sempre o mesmo comando, sempre citado.
	b.WriteString("\n## Como verificar que terminou\n\n")
	b.WriteString("```sh\nanchors map build\n")
	// O `--changed` recebe a peça que ESTA etapa produz, não o alvo da unidade: na etapa
	// `spec`, o `.ts` do alvo em geral ainda não existe, e o `check` aborta sem conferir
	// nada. Prescrever um comando que a própria etapa impede de funcionar é a régua se
	// contradizendo dentro do mesmo prompt.
	alvo := alvoDaVerificacao(rel, artifact, layer, cfg)
	fmt.Fprintf(&b, "anchors check --changed %s --no-record --deterministic\n```\n\n", alvo)
	b.WriteString("Todos os gates **bloqueantes** devem sair com `✗0`. Um `~` não é falha — é o " +
		"gate dizendo que não teve o que confrontar (o motivo vem escrito). Se um gate reprovar, " +
		"corrija antes de considerar a etapa concluída.\n\n")
	// O `--no-record` acima é para ITERAR: confere sem efeito colateral, quantas vezes
	// precisar. Mas a etapa que FECHA tem de registrar — o carimbo é a memória de que
	// aquela aresta foi validada, e é dele que o `anchors stale` vive.
	//
	// Medido: com o `work` prescrevendo só a forma `--no-record`, o carimbo nunca
	// acontecia durante o ciclo, e `anchors stale` reportava 9.975 de 9.975 arestas
	// "nunca validada" — um comando inutilizável, que parecia bug do framework e era o
	// prompt nunca mandando registrar. Um `check --all` bastou para derrubar o número a
	// 19.
	b.WriteString("Quando os gates estiverem verdes, rode uma última vez **sem** " +
		"`--no-record`:\n\n```sh\n")
	fmt.Fprintf(&b, "anchors check --changed %s --deterministic\n```\n\n", alvo)
	b.WriteString("É o que CARIMBA as arestas desta unidade como validadas. Sem esse passo o " +
		"mapa não guarda que o confronto aconteceu: o trabalho fica pronto e o `anchors stale` " +
		"segue dizendo que nunca ninguém olhou.\n")

	// OPT-OUTS: a saída LEGÍTIMA quando um gate reprova algo que é decisão consciente.
	// Sem esta seção, o agente descobre a existência da dispensa pela mensagem de erro do
	// gate — foi o que aconteceu num teste real: `@no-scenario` não estava em guide algum,
	// e o autor teve de adivinhar a sintaxe por tentativa e erro.
	b.WriteString("\n## Se um gate reprovar algo que é decisão consciente\n\n")
	b.WriteString("Existe uma porta legítima, e ela **exige a razão escrita** — marcador nu " +
		"não dispensa nada. A dispensa fica na linha do que ela dispensa (ou no comentário " +
		"logo acima), datada pelo git e visível para quem ler depois:\n\n")
	switch artifact {
	case "spec":
		b.WriteString("- `@no-code: <razão>` na linha de uma regra que NÃO terá trecho de código " +
			"(restrição satisfeita pela ausência, invariante provada por teste). " +
			"**Se ao implementar a razão não se sustentar** — a regra afinal precisa de código, " +
			"ou não deveria existir —, isso não é detalhe de formatação: é a spec estar errada " +
			"sobre a própria unidade. Abra issue (`anchors judge <alvo> --gate review " +
			"--verdict fail --reason \"…\"`) em vez de apagar a regra ou inventar um " +
			"comentário para calar o gate.\n")
		b.WriteString("- `@no-scenario: <razão>` na linha de um requisito que não terá cenário " +
			"— para o que é mesmo não-observável por cenário, não para o que dá trabalho testar.\n")
	case "code":
		b.WriteString("- `@no-paginate: <razão>` numa função que promete o conjunto e não pagina " +
			"(o limite é deliberado e você sabe o teto).\n" +
			"- `@allow-boundary: <razão>` numa linha que cruza uma fronteira de camada — dívida " +
			"reconhecida, visível no código em vez de numa lista distante.\n")
	}
	b.WriteString("Ver a lista completa em `anchors guide header`. Dispensar é decisão " +
		"registrada; apagar o problema não.\n")

	// REGRA DE OURO: o que fazer quando a régua não decide.
	b.WriteString("\n## Se a régua não decidir\n\n")
	b.WriteString("**Não chute.** Se a spec/guide admite duas leituras com consequências " +
		"diferentes, ou não decide algo que o código vai precisar, o registro tem LUGAR — " +
		"não é um comentário no PR, que morre no merge:\n\n")
	if artifact == "spec" {
		b.WriteString("- **Escreva a pergunta em `## Decisões em aberto`**, na spec que você " +
			"está produzindo. Enquanto houver item ali, a spec não passa por pronta (gate " +
			"`open-questions-resolved`) — que é o efeito desejado: implementar com pergunta " +
			"aberta é adivinhar, e a adivinhação não é confrontada por gate nenhum, porque " +
			"todas as peças existem e se referenciam.\n" +
			"- **Não invente a resposta para esvaziar a seção.** Quando ela vier, PROMOVA-A " +
			"a regra (com código) e marque o item como resolvido — a pergunta fica no rastro.\n" +
			"- **Se não há dúvida, escreva `nenhuma`.** Afirmar que se olhou é diferente de " +
			"omitir a seção.\n")
	} else {
		b.WriteString("- **Pare e leia `## Decisões em aberto` na spec.** Se houver item ali, " +
			"a decisão ainda não foi tomada: não a tome sozinho no código.\n" +
			"- **Se a lacuna é nova** (você descobriu implementando), ela pertence à SPEC. " +
			"Registre-a lá e reporte — não resolva no código, onde ninguém vai encontrá-la.\n" +
			"- Escolha a leitura mais defensável para não travar, mas deixe a escolha " +
			"VISÍVEL onde se decide. Uma decisão silenciosa vira bug que ninguém rastreia.\n")
	}
	return b.String(), nil
}

// guidesFor devolve os guides a ler para produzir `artifact` sobre um alvo da camada
// `l`. São DUAS origens, e ambas importam:
//   - o guide do ARTEFATO (quem rege a tag `spec`/`feature`/`test`) — como se escreve
//     uma spec neste projeto, qualquer que seja a camada do alvo;
//   - os guides da CAMADA do alvo (quem rege as tags dela) — o que esta camada exige.
//
// Considerar só a camada esconderia o SPEC_GUIDE de quem vai escrever a spec de um
// repository (o alvo tem a tag `repository`, não `spec`) — que é justamente o guide
// mais relevante da tarefa.
func guidesFor(l config.Layer, cfg *config.Config, artifact string) []string {
	tags := map[string]bool{artifact: true}
	for _, t := range l.Tags {
		tags[t] = true
	}
	var out []string
	seen := map[string]bool{}
	for _, g := range cfg.Governs {
		if tags[g.Governs] && !seen[g.From] {
			seen[g.From] = true
			out = append(out, g.From)
		}
	}
	sort.Strings(out)
	return out
}

// writeTrincaPaths mostra onde cada peça da trinca nasce para este alvo, usando o
// `derived:` do projeto (co-location por padrão, overrides por camada).
func writeTrincaPaths(b *strings.Builder, rel, artifact, layer string, cfg *config.Config, g *mapx.Graph) {
	if cfg.Derived == nil {
		b.WriteString("> O projeto não declara `derived:` — confirme onde as peças moram " +
			"olhando os vizinhos da camada.\n")
		return
	}
	files, overridden := caminhosDerivados(rel, layer, cfg)
	dispensadas := pecasDispensadas(layer, cfg)
	order := []string{"spec", "feature", "test"}
	for _, k := range order {
		tpl, ok := files[k]
		if !ok {
			continue
		}
		p := tpl
		mark := " "
		if k == artifact {
			mark = "→" // a peça desta etapa
		}
		exists := ""
		if _, err := os.Stat(p); err == nil {
			exists = "  (já existe)"
		}
		note := ""
		if overridden[k] {
			note = "  ← override da camada (não é co-localizado)"
		}
		// A camada pode DISPENSAR uma peça (`trinca_opcional`), e o prompt tem de dizer
		// isso: listá-la como "onde nasce" faz o executor obediente criar exatamente o
		// artefato que a régua declarou não querer. Aconteceu num E2E real — a camada
		// dizia "opt-out honesto, em vez de 49 features/testes vazios", e o mesmo prompt
		// mandava criar os dois.
		if dispensadas[k] {
			fmt.Fprintf(b, "  `%s` — %s  ← DISPENSADA por esta camada (`trinca_opcional`): "+
				"NÃO crie\n", p, k)
			continue
		}
		fmt.Fprintf(b, "%s `%s` — %s%s%s\n", mark, p, k, exists, note)
	}
	b.WriteString("\n> Caminhos já resolvidos para ESTA camada (co-location + overrides do " +
		"`derived:`). Não crie a peça em outro lugar.\n")
}

// outOfScope é a FRONTEIRA do papel: o que pertence a outra etapa. Universal por
// artefato — a doutrina do Anchors, não do projeto.
func outOfScope(artifact string) []string {
	switch artifact {
	case "review-plan-draft":
		return []string{
			"**Não revise o CÓDIGO** — ele ainda não existe. Seu alvo é o plano como roteiro: " +
				"ele é executável por quem não estava na conversa que o gerou?",
			"**Não confunda com o review de conjunto** (`review-plan`), que roda no FIM e " +
				"confronta a costura do que foi entregue. Este roda ANTES, e confronta o roteiro.",
			"**Não aceite \"decide-se depois\" dentro de um item.** Fora do checklist, em " +
				"prosa, com quem decide — é o comportamento correto e não é achado.",
		}
	case "review-plan":
		return []string{
			"**Não revise a peça isolada** — isso já foi feito. Seu alvo é o que só aparece " +
				"quando as peças se encontram.",
			"**Não corrija nada.** Reporte; quem acha não conserta.",
			"**Não aceite \"cada parte passou\" como prova do todo.** É a premissa que os 4 " +
				"achados que atravessam unidades desmentem.",
		}
	case "review":
		return []string{
			"**Não corrija o código** — nem o \"obviamente errado\". Reporte; quem acha não conserta.",
			"**Não reescreva a spec.** Se ela não decide o suficiente, isso É o achado.",
			"**Não confie no verde.** Todos os gates passando é o estado NORMAL de um trabalho " +
				"que chega até aqui — é justamente por isso que o review existe.",
			"**Não relate o que não executou.** \"Parece que\" e \"provavelmente\" não são achados; " +
				"comando + saída são.",
		}
	case "spec":
		return []string{
			"**Não escreva código, feature nem teste** — esta etapa produz só a spec.",
			"**Não descreva implementação** (estruturas internas, algoritmo, nomes de variável): a spec descreve COMPORTAMENTO observável e contrato.",
			"**Não decida por outra camada**: se a decisão pertence ao schema/à regra/à interface, aponte para lá em vez de replicar.",
		}
	case "code":
		return []string{
			"**Não reescreva a spec para caber no código** — se a spec está errada, reporte; não a ajuste em silêncio.",
			"**Não invente regra que a spec não tem**: comportamento novo exige spec antes.",
			"**Não misture camadas**: respeite o que a camada do alvo pode importar (ver `layers:`).",
		}
	case "feature":
		return []string{
			"**Não crie código de cenário novo**: todo `@CODE-X##` da feature deve existir na spec.",
			"**Não descreva implementação** nos passos — Gherkin fala de comportamento observável.",
			"**Não escreva o teste aqui** — a feature declara os cenários; o teste os prova.",
		}
	case "test":
		return []string{
			"**Não teste o que a feature não declara**: cada teste rastreia a um cenário com código.",
			"**Não mocke a lógica sob teste** (mock que replica a regra prova a si mesmo, não o código).",
			"**Não altere o código para o teste passar** sem entender se o defeito é do código ou do teste.",
		}
	}
	return nil
}

// procedureFor é o procedimento UNIVERSAL do artefato. O projeto acrescenta os seus
// passos via `layers.<camada>.work.<artefato>` (híbrido).
func procedureFor(artifact string, cfg *config.Config) []string {
	switch artifact {
	case "review-plan-draft":
		// O review do PLANO RECÉM-ESCRITO, antes de alguém executá-lo.
		//
		// É o único ponto do ciclo em que corrigir custa uma edição de texto. Depois, o
		// defeito do plano já se espalhou pelas specs que ele semeou — e o custo é
		// assimétrico: um erro na spec afeta uma unidade; um erro no plano afeta todas.
		//
		// A régua vem de defeitos medidos em planos reais: item cuja existência depende de
		// decisão não tomada (o executor decidiu sozinho), ambiguidade que três rodadas
		// seguidas tropeçaram na mesma pedra, e caminho que a estrutura já mudou.
		return []string{
			"**Cace o item cuja EXISTÊNCIA depende de decisão não tomada.** Distinga: " +
				"delegar ONDE se decide é legítimo (\"os índices ficam na spec do schema, " +
				"decididos ao ver o schema atual\" — o item existe, o conteúdo mora noutro " +
				"artefato). Adiar SE o item existe é o defeito (\"nasce se o mobile chamar " +
				"por X; senão marcar fora de escopo\") — quem executa vai decidir sozinho ou " +
				"parar o fluxo. Aconteceu: o agente decidiu, e ninguém tinha pedido.",
			"**Confronte cada caminho citado contra o DISCO.** O plano cita arquivos e " +
				"diretórios; a estrutura muda. Um plano que manda editar um arquivo que foi " +
				"dividido guia o executor para o lugar errado — e três rodadas seguidas " +
				"tropeçaram na mesma ambiguidade (`models/` em dois lugares com regimes opostos).",
			"**Confronte o plano contra a Estrutura declarada.** Ele semeia spec em camada " +
				"que TEM spec? Manda criar peça em camada declarativa? O `anchors work spec " +
				"--for <alvo>` responde por alvo — use-o nos itens de que você duvidar.",
			"**Leia cada item perguntando \"o que quem executa isto precisa saber e não " +
				"está aqui?\"** — não \"está bem escrito?\". O plano é executável ou não é.",
			"**Verifique a ordem.** Um item que depende de outro precisa vir depois — e o " +
				"plano precisa DIZER a dependência, não deixá-la implícita na numeração.",
			"**A Definição de Pronto é verificável?** Cada critério precisa ser confrontável " +
				"por alguém que não escreveu o plano. \"Funciona bem\" não é critério.",
			"**NÃO reescreva o plano.** Reporte; quem acha não conserta — e um plano " +
				"reescrito pelo revisor perde o dono.",
		}
	case "review-plan":
		// O review de CONJUNTO. Derivado da medição: dos 6 achados de um review adversarial
		// real, apenas 2 cabiam no escopo de uma unidade. Os outros 4 atravessavam — e em
		// todos eles CADA PEÇA estava correta sozinha, que é exatamente por que o review
		// por unidade não os alcança e nenhum gate os pega.
		return []string{
			"**Leia o plano inteiro e os registros em `changes/`.** Seu escopo é a COSTURA, " +
				"não a peça: o review por unidade já garantiu cada uma. Comece listando as " +
				"unidades entregues e como elas se conectam.",
			"**Siga o DADO de ponta a ponta.** Quem grava, quem lê, quem apaga. O caso que " +
				"motivou: a spec do modelo mandava plantar o campo de propriedade, o DAO não " +
				"plantava e o handler não preenchia — três peças corretas, e o registro criado " +
				"por Lambda ficava invisível ao dono na tela.",
			"**Confronte PROMESSA entre camadas.** O que a camada de baixo entrega é o que a " +
				"de cima assume? O caso: o DAO truncava a consulta sem paginar, e a regra que " +
				"o consome escolhe \"a versão de maior mês\" — truncar não dá erro, dá a " +
				"resposta ERRADA em silêncio.",
			"**Confronte REGRA contra REGRA**, entre specs de camadas diferentes e dentro da " +
				"mesma spec. Nenhum gate lê duas regras juntas. Uma invariante declarada numa " +
				"spec pode ser violada pelo código que outra spec autoriza.",
			"**Prove a obrigação transversal por EFEITO, não por estrutura.** Um dever pode " +
				"passar em toda peça (a tabela aparece no handler) e não funcionar no todo. " +
				"Mute o bloco que o cumpre e rode a suíte: se nada cai, o dever não está provado.",
			"**Confronte os registros de entrega contra o disco.** Cada `changes/*.md` afirma " +
				"o que foi feito, o que foi decidido sozinho e o que não está provado. " +
				"Divergência entre o declarado e o real é achado — e já apareceu: um registro " +
				"afirmava ter atualizado contagens de um arquivo que não foi tocado.",
			"**Verifique a Definição de Pronto do plano**, item a item, com evidência de " +
				"execução. O que ficou por fazer é achado; o que ficou por fazer e foi marcado " +
				"como feito é achado grave.",
			"**Classifique por severidade e NÃO corrija.** Sua saída é o relatório; quem acha " +
				"não conserta.",
		}
	case "review":
		// Esta régua é DERIVADA de defeitos reais. Em três rodadas de um E2E, 7 defeitos
		// graves passaram com todos os gates verdes; nenhum foi achado lendo código. Os
		// que apareceram vieram de ATAQUE — mutar a regra, rodar com entrada de borda,
		// confrontar campo do modelo contra a regra que o interpreta. Por isso o
		// procedimento manda executar, não revisar no sentido usual.
		return []string{
			"**Leia o registro de entrega** (`changes/`) e confronte a INTENÇÃO declarada " +
				"contra o que está no disco. Divergência entre o que o autor diz ter feito e o " +
				"que ele fez é, por si só, um achado.",
			"**Confronte spec ↔ código, regra a regra.** Cada regra catalogada está " +
				"implementada? Alguma foi implementada ao contrário? Há comportamento no código " +
				"que nenhuma regra governa? (campo do modelo que ninguém lê, ramo que a spec " +
				"não prevê).",
			"**MUTE a regra e rode a suíte.** Apague a linha que implementa cada regra; se os " +
				"testes continuarem verdes, aquele teste NÃO prova aquela linha. É o único " +
				"instrumento que mede o poder do teste — e já pegou regra com 12 testes verdes " +
				"e zero prova. Restaure o arquivo depois de cada mutação.",
			"**ATAQUE por execução, com entrada adversarial.** Escreva um script que chame a " +
				"unidade com: vazio, limite, valor que a spec permite mas o código não trata, " +
				"ordem estranha, caractere especial na chave. Relate com o comando e a SAÍDA " +
				"real — não com o que você deduziu lendo.",
			"**Confronte as regras entre si.** Duas regras da mesma spec podem se contradizer " +
				"(uma diz ausência, outra diz valor, para o mesmo caso). Uma regra de baixo pode " +
				"contradizer a de cima. Nenhum gate confronta regra↔regra — este é o lugar.",
			"**Verifique o que a peça promete ao MUNDO.** Recurso declarado existe na infra? " +
				"Env var citada é provida? Índice que a consulta precisa foi criado? O typecheck " +
				"passa por acidente (`!`, índice de objeto) escondendo algo que só quebra em " +
				"produção?",
			"**Classifique por severidade** (CRÍTICO/SÉRIO/MENOR) e escreva a evidência de cada " +
				"achado: comando + saída, ou arquivo:linha. Achado sem evidência não é acionável.",
			"**NÃO corrija.** Sua saída é o relatório e, se houver achado crítico ou sério, uma " +
				"issue. Quem acha e quem conserta são papéis diferentes — é o que preserva a " +
				"independência do próximo review.",
		}
	case "spec":
		return []string{
			"Confirme a camada do alvo e se ela é REGIDA (tem spec) ou RECONHECIDA (`regime: declarativo` → não tem).",
			"Leia 2–3 specs VIZINHAS da mesma camada: elas mostram o dialeto real do projeto (que pode divergir do template).",
			"Gere a moldura com `anchors new spec <Nome> --preset <preset> --out <caminho>` (veja os presets em `anchors new spec --list-sections`).",
			"Preencha as seções catalogando cada regra com seu código (`{CODE}-<letra><NN>`); as letras válidas estão em `rule_types` no anchors.yaml.",
			"Na Tabela de Dependências, use `crases` só em símbolos que o código VAI usar — isso vira contrato verificável. Descrição livre fica em prosa.",
		}
	case "code":
		return []string{
			"Leia a spec inteira antes de escrever a primeira linha; ela é a régua.",
			"Leia 1–2 arquivos vizinhos da mesma camada para seguir o padrão local (imports, erro, estilo).",
			marcacaoDaRegra(cfg),
			"Se a spec prometer um símbolo na Tabela de Dependências, USE esse símbolo — o gate `dependency-honored` confronta.",
		}
	case "feature":
		return []string{
			"Liste os códigos de cenário que a spec declara — a feature cobre esses, e só esses.",
			"Gere a moldura com `anchors new feature <Nome> --out <caminho>`.",
			"Escreva um cenário por comportamento observável, com a tag do código e a tag de regime (`@nivel-unit` etc., ver `derived.regimes`).",
			"O TÍTULO do cenário deve descrever o comportamento — o teste vai espelhá-lo (o gate confronta código E descrição).",
		}
	case "test":
		return []string{
			"Leia a feature: cada cenário dela vira um teste, com o código no título.",
			"INVARIANTE (`-I##`) se prova em CICLO FECHADO: aplique o resultado do produtor e verifique pelo CONSUMIDOR. Montar o estado à mão e verificar só a leitura testa o consumidor, não a composição — e é onde os bugs se escondem.",
			"Depois de escrever, MUTE a regra: apague a linha que ela implementa e rode. Se o teste continuar verde, ele não prova aquela linha — conserte o teste, não o código.",
			"Descubra ONDE o teste desta camada mora (co-location ou `derived.overrides`) antes de criar o arquivo.",
			"Espelhe o título do cenário no `it(...)` — descrição divergente é drift que o gate acusa.",
			"Rode os testes e confirme que passam por afirmarem o comportamento certo, não por acaso.",
		}
	}
	return nil
}

// writeRegimes imprime o vocabulário de tags de regime do projeto — a tradução da tag
// que o cenário declara para o regime canônico, e a superfície onde ele é provado.
func writeRegimes(b *strings.Builder, cfg *config.Config) {
	if cfg.Derived == nil || len(cfg.Derived.Regimes) == 0 {
		return
	}
	b.WriteString("\n## Tags de regime válidas (use EXATAMENTE estas)\n\n")
	tags := make([]string, 0, len(cfg.Derived.Regimes))
	for t := range cfg.Derived.Regimes {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	for _, t := range tags {
		canon := cfg.Derived.Regimes[t]
		surface := ""
		if cfg.Derived.Surfaces != nil {
			if sfc, ok := cfg.Derived.Surfaces[canon]; ok {
				surface = fmt.Sprintf(" → provado na superfície `%s`", sfc)
			}
		}
		fmt.Fprintf(b, "- `@%s` (regime %s)%s\n", t, canon, surface)
	}
	b.WriteString("\nA tag é do PROJETO e não é traduzível: escrever `@nivel-integracao` " +
		"em vez de `@nivel-integration` faz o cenário não ser confrontado por nenhum gate.\n")
}

// unidadeDaPecaDerivada resolve a UNIDADE quando o alvo dado é uma peça derivada
// (spec/feature/test). Prefere o mapa — a aresta `specifies` diz exatamente qual código a
// spec descreve; sem mapa, cai na convenção de nome (tronco + extensões usuais).
func unidadeDaPecaDerivada(root, rel string, cfg *config.Config, g *mapx.Graph) (string, bool) {
	layer, _ := scan.Classify(rel, cfg)
	l, ok := cfg.Layers[layer]
	if !ok {
		return "", false
	}
	// É peça derivada? O kind da camada decide (spec/feature/test), não a extensão.
	switch l.Kind {
	case "spec", "feature", "test":
	default:
		return "", false
	}

	// 1) pelo mapa: a spec APONTA o código (`specifies`); feature/test chegam via a spec.
	if g != nil {
		if alvo := alvoPorAresta(g, rel, "specifies"); alvo != "" {
			return alvo, true
		}
		// feature/test → sobe até a spec (`covered-by`/`tested-by` chegam NELES)
		for _, e := range g.Edges {
			if e.To != rel {
				continue
			}
			if e.Type == "covered-by" || e.Type == "tested-by" {
				if alvo := alvoPorAresta(g, e.From, "specifies"); alvo != "" {
					return alvo, true
				}
			}
		}
	}

	// 2) pela convenção: mesmo tronco, extensões de código usuais no mesmo diretório.
	dir := filepath.Dir(rel)
	base := filepath.Base(rel)
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".spec.ts", "_test.go", "_test.py", "_spec.rb"} {
		base = strings.TrimSuffix(base, suf)
	}
	if i := strings.Index(base, "."); i > 0 {
		base = base[:i]
	}
	for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".java", ".rs", ".rb", ".php", ".kt", ".cs"} {
		// ToSlash porque o que se devolve é o id de um nó, não um caminho de disco: no
		// Windows o Join daria "src\x.ts", e o id não casaria com o mapa (gravado com "/").
		cand := filepath.ToSlash(filepath.Join(dir, base+ext))
		if _, err := os.Stat(filepath.Join(root, cand)); err == nil {
			return cand, true
		}
	}
	return "", false
}

func alvoPorAresta(g *mapx.Graph, from, tipo string) string {
	for _, e := range g.Edges {
		if e.From == from && string(e.Type) == tipo {
			return e.To
		}
	}
	return ""
}

// procedureDeclarativa é o procedimento de uma camada RECONHECIDA (`regime: declarativo`).
// Não há spec para ler — a régua é o CONTRATO da camada vizinha que este arquivo serve, e
// o dialeto dos irmãos. O risco característico aqui não é divergir de uma spec: é a camada
// declarativa DECIDIR alguma coisa, virando regra escondida onde ninguém procura.
func procedureDeclarativa(rel string) []string {
	return []string{
		"Esta camada **não tem spec**: a régua é o contrato de quem consome este arquivo, " +
			"mais o dialeto dos vizinhos. Leia 2–3 irmãos da mesma camada ANTES de escrever.",
		"Traduza, transporte ou declare — **não decida**. Se você se pegar escrevendo um " +
			"`if` de negócio, um default que muda o resultado, ou uma validação de regra, " +
			"pare: isso pertence à camada que decide, e aqui vira regra escondida onde " +
			"ninguém vai procurar.",
		"Mantenha a superfície previsível: mesmas convenções de nome, erro e retorno dos " +
			"irmãos. Quem consome esta camada conta com isso.",
		"Se faltar uma decisão para você conseguir escrever, ela falta na spec de quem " +
			"consome — reporte lá, não resolva aqui.",
	}
}

// exigenciasDosGates traduz, em requisitos legíveis, o que os gates INTERNOS declarados
// para este artefato vão confrontar. Só descreve gate que o projeto realmente declarou —
// prometer cobrança que não existe é tão ruim quanto esconder a que existe.
func exigenciasDosGates(artifact string, cfg *config.Config) []string {
	// o que cada checker interno exige, em uma frase acionável
	porChecker := map[string]string{
		"spec-sections": "**Toda regra precisa estar CATALOGADA** — código + lugar estruturado. " +
			"Vale cabeçalho (`### ABCDX-B01 — ...`), linha de tabela (`| \\`ABCDX-B01\\` | ... |`) " +
			"ou bullet-negrito (`- **ABCDX-B01** ...`). Menção solta em prosa NÃO conta, e " +
			"placeholder não preenchido reprova.",
		"has-code": "O arquivo carrega ao menos um **código de cenário** (a identidade).",
		"header-conforme": "O **header `@anchors`** no topo, com a identidade: `code:` se este " +
			"artefato é dono, `ref:` se referencia a spec.",
		"rule-types": "Cada **letra de código** (o `B` de `-B01`) precisa estar declarada no " +
			"vocabulário `rule_types` do projeto — e a seção que a define, também.",
		"route-declared":     "Spec de TELA declara a **rota** no cabeçalho.",
		"dependency-honored": "Todo símbolo prometido na **Tabela de Dependências** (entre crases) é usado no código.",
		"spec-feature-match": "Todo requisito declarado tem **cenário na feature** (ou `@no-scenario: <razão>`).",
		"regra-implementada": "Toda regra catalogada aparece **no código** (o trecho que a realiza traz o código dela em comentário) — ou é dispensada na linha dela com `@no-code: <razão>`, para o que é satisfeito pela AUSÊNCIA de código. Declare regra a regra: é o que troca adivinhação por confronto.",
		"open-questions-resolved": "A seção **`## Decisões em aberto`** É OBRIGATÓRIA e nunca " +
			"fica vazia: ou lista o que a spec não decide, ou traz `nenhuma`. Escrever " +
			"`nenhuma` é uma AFIRMAÇÃO — \"olhei e não há dúvida\" —, diferente de omitir a " +
			"seção, que não diz nada. O que estiver aberto ali segura a spec até ser decidido.",
		"code-reference-valid": "Todo **código citado** precisa existir no projeto (citação órfã reprova).",
		"scenario-asserts":     "O passo de RESULTADO afirma um **resultado observável** — não `\\\"o efeito X se verifica\\\"`.",
		"feature-test-match": "Cada cenário da feature tem um teste que o prova, e o teste " +
			"**declara qual cenário prova**: o código vai no TÍTULO do caso (`it('[ABCDX-B01] " +
			"…')` ou `describe('… [ABCDX-B01]')`), não em comentário — o gate remove comentários " +
			"antes de procurar. Marcar é AFIRMAR que aquele teste prova aquele cenário; " +
			"marcador no teste errado cria rastreabilidade falsa, e todo gate relacional passa " +
			"a confrontar o par errado com tudo verde.",
		"non-empty":          "O arquivo **não pode ser um esqueleto vazio**.",
		"trinca-completa":    "A unidade precisa da **trinca completa** (spec + feature + teste).",
		"ref-resolves":       "O `ref:` precisa apontar para o **`code:` da spec irmã** — não para outra.",
		"pagination-honored": "Função que promete um conjunto **não devolve a primeira página** em silêncio.",
		"layer-boundary":     "Respeite as **fronteiras de camada** declaradas em `boundaries:`.",
	}
	var out []string
	visto := map[string]bool{}
	for _, gt := range cfg.Gates {
		if !gateVale(gt, artifact) || visto[gt.Check] {
			continue
		}
		if frase, ok := porChecker[gt.Check]; ok {
			visto[gt.Check] = true
			if !gt.IsBlocking() {
				frase += " *(informativo)*"
			}
			out = append(out, frase)
		}
	}
	return out
}

// gateVale: este gate se aplica ao artefato que está sendo produzido?
//
// O `on:` do gate diz sobre QUAL NÓ ele roda; esta função responde outra pergunta — quem
// precisa CONHECER a exigência ao escrever. Nem sempre é o mesmo.
//
// `feature-test-match` roda sobre a feature e cobra do TESTE: exige que o código do
// cenário apareça no título do caso. Enquanto ele só aparecia na etapa `feature`, quem
// escrevia o teste não era avisado — medido, um agente escreveu 36 features corretas e
// depois esbarrou num gate BLOQUEANTE vermelho sem entender por quê; a exigência só estava
// legível no código-fonte do gate.
func gateVale(gt config.Gate, artifact string) bool {
	for _, k := range gt.On {
		if k == artifact {
			return true
		}
	}
	// Gates que rodam numa peça e cobram de OUTRA: quem produz a peça cobrada tem de saber.
	if artifact == "test" && gt.Check == "feature-test-match" {
		return true
	}
	return false
}

// writeDeliveryRecord anexa o REGISTRO DE ENTREGA ao prompt do revisor, quando existe.
//
// É a metade declarada do confronto: o autor escreveu o que acha que fez, quais decisões
// tomou sozinho e o que sabe não estar provado. O revisor confronta isso contra o disco —
// e a divergência entre as duas versões é, por si só, um achado que não existe forma de
// detectar sem ter as duas.
//
// Sem registro, o review ainda roda (a unidade está no disco), mas perde o escopo e a
// declaração de intenção. Por isso a ausência é reportada, não silenciada.
func writeDeliveryRecord(b *strings.Builder, root, rel string) {
	pend, _ := change.Pending(root)
	unidade := filepath.ToSlash(rel)
	// TODOS os registros da unidade, não o primeiro. Uma unidade acumula uma entrega por
	// etapa (spec, feature, code, test), e servir só a primeira dá ao revisor a foto mais
	// VELHA: medido, o registro da etapa `code` dizia "NADA está provado: o .test.ts é da
	// próxima etapa" — verdade quando foi escrito, falso quando o review rodou, com o
	// teste já no disco. Um revisor obediente pularia a mutação por achar que não há suíte.
	var achados []string
	for _, p := range pend {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if m := changeUnitRE.FindStringSubmatch(string(raw)); m == nil || m[1] != unidade {
			continue
		}
		relPath, _ := filepath.Rel(root, p)
		achados = append(achados, fmt.Sprintf("### Registro de entrega — `%s`\n\n"+
			"```markdown\n%s\n```\n", relPath, strings.TrimSpace(string(raw))))
	}
	if len(achados) > 0 {
		sort.Strings(achados)
		fmt.Fprintf(b, "\n## Registros de entrega desta unidade (%d)\n\n", len(achados))
		b.WriteString("O autor declarou o que segue, uma vez por etapa. **Confronte cada " +
			"afirmação contra o disco** — sobretudo \"nenhuma decisão livre\" e \"nada sem " +
			"prova\", que são afirmações fortes e é onde os defeitos se escondem. Note que " +
			"um registro descreve o estado NAQUELE momento: o que ele diz que falta pode já " +
			"existir.\n\n")
		for _, a := range achados {
			b.WriteString(a + "\n")
		}
		return
	}
	b.WriteString("\n> **Sem registro de entrega** para esta unidade (`anchors deliver`). " +
		"Você revisa o que está no disco, mas sem a intenção declarada do autor — logo, sem " +
		"como confrontar o que ele ACHA que fez contra o que fez. Registre isso no relatório.\n")
}

// caminhosDerivados resolve, para um alvo e sua camada, ONDE cada peça da trinca nasce —
// aplicando co-location e os overrides do `derived:`, com os placeholders já substituídos.
//
// Existe separada porque duas coisas precisam da mesma resposta: a seção que MOSTRA os
// caminhos e o comando de verificação que o prompt prescreve. Enquanto só a primeira
// resolvia, o comando prescrito apontava para o alvo (`--for <alvo>.ts`) mesmo na etapa
// `spec` — onde o `.ts` ainda não existe e o `check` aborta:
//
//	Error: "packages/.../metadataVersioning.ts" não existe no disco nem no mapa
//
// O mesmo prompt dizia "não escreva código nesta etapa" e prescrevia um comando que só
// funciona com o código escrito. Os dois agentes de spec de um E2E real bateram nisso,
// independentemente.
func caminhosDerivados(rel, layer string, cfg *config.Config) (map[string]string, map[string]bool) {
	files, overridden := map[string]string{}, map[string]bool{}
	if cfg.Derived == nil {
		return files, overridden
	}
	dir := filepath.Dir(rel)
	name := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))
	// O módulo é o diretório-pai — usado por overrides que agrupam por Lambda/módulo
	// (ex.: `packages/backend/__tests__/unit/lambdas/{{module}}.test.ts`).
	module := filepath.Base(dir)
	bruto := map[string]string{}
	for k, v := range cfg.Derived.Files {
		bruto[k] = v
	}
	// OVERRIDE da camada vence a co-location. Resolvê-lo aqui é o ponto do comando:
	// mandar o leitor "conferir os overrides" seria devolver a ele exatamente a
	// descoberta manual que este prompt existe para eliminar.
	for _, ov := range cfg.Derived.Overrides {
		if ov.When != layer {
			continue
		}
		for k, v := range ov.Files {
			bruto[k] = v
			overridden[k] = true
		}
	}
	// `{{ext}}` é a extensão do ALVO, não "ts" fixo. Numa tela (`.tsx`), o valor fixo
	// fazia o prompt prescrever `Tela.test.ts` sob a frase "Não crie a peça em outro
	// lugar" — um arquivo que, criado, é erro de sintaxe (JSX em `.ts`). O arquivo real
	// ao lado era `.test.tsx`, e o mapa o resolvia certo: só o prompt mentia.
	ext := strings.TrimPrefix(filepath.Ext(rel), ".")
	if ext == "" || ext == "md" {
		ext = "ts" // alvo sem extensão de código (ou uma spec): o default do projeto
	}
	r := strings.NewReplacer("{{dir}}", dir, "{{name}}", name, "{{module}}", module, "{{ext}}", ext)
	for k, v := range bruto {
		files[k] = r.Replace(v)
	}
	return files, overridden
}

// alvoDaVerificacao devolve o caminho que o `anchors check --changed` deve receber nesta
// etapa: a peça que ELA produz, não o alvo da unidade. Cai no alvo quando a etapa não
// produz peça derivada (código, review).
func alvoDaVerificacao(rel, artifact, layer string, cfg *config.Config) string {
	switch artifact {
	case "spec", "feature", "test":
		if files, _ := caminhosDerivados(rel, layer, cfg); files[artifact] != "" {
			return files[artifact]
		}
	}
	return rel
}

// pecasDispensadas traduz o `trinca_opcional` da camada (declarado por ARESTA) para as
// PEÇAS que ele dispensa. `covered-by` é a aresta spec→feature, logo dispensa a feature;
// `tested-by` é feature→test, logo dispensa o teste.
func pecasDispensadas(layer string, cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	if cfg == nil || layer == "" {
		return out
	}
	l, ok := cfg.Layers[layer]
	if !ok {
		return out
	}
	porAresta := map[string]string{"covered-by": "feature", "tested-by": "test", "specifies": "spec"}
	for _, aresta := range l.TrincaOpcional {
		if peca := porAresta[aresta]; peca != "" {
			out[peca] = true
		}
	}
	return out
}

// writeIssuesAbertas traz para o prompt os achados ainda em aberto sobre esta unidade —
// de `issues/todo/` e `issues/doing/`, que são os estados vivos.
func writeIssuesAbertas(b *strings.Builder, root, rel string) {
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	for _, suf := range []string{".spec.md", ".feature", ".test", ".spec"} {
		base = strings.TrimSuffix(base, suf)
	}
	slug := strings.ReplaceAll(base, "/", "-")

	var achadas []string
	for _, st := range []issue.State{issue.Todo, issue.Doing} {
		nomes, _ := issue.List(root, st)
		for _, nome := range nomes {
			if !strings.Contains(nome, slug) {
				continue
			}
			achadas = append(achadas, filepath.Join(issue.Dir, string(st), nome))
		}
	}
	if len(achadas) == 0 {
		return
	}
	sort.Strings(achadas)
	b.WriteString("\n## Achados JÁ REGISTRADOS sobre esta unidade (leia ANTES de escrever)\n\n")
	b.WriteString("Alguém já confrontou esta unidade e encontrou o que está abaixo. Não é " +
		"histórico: é trabalho pendente, e os gates NÃO o repetem — eles confrontam o que é " +
		"declarável, e estes achados são justamente o que sobrou fora disso.\n\n")
	for _, a := range achadas {
		fmt.Fprintf(b, "- `%s`\n", a)
	}
	b.WriteString("\n> Leia cada uma. O que ainda valer, resolva nesta etapa; o que já não " +
		"valer, feche (`anchors judge <alvo> --gate review --verdict pass --reason \"...\"`). " +
		"Uma issue que ninguém lê é um defeito que o pipeline já viu e deixou passar.\n")
}

// marcacaoDaRegra emite o passo de ligar regra↔código conforme o projeto a EXIGE ou não.
//
// A ressalva "se o projeto usa esse padrão" existia para não impor a prática a quem não
// a adotou — mas ela também dava saída a quem a adotou: quem implementa lê "se", decide
// que não, e a marcação nunca acontece. O gate `regra-implementada` então cobra algo que
// o procedimento apresentou como opcional, e a dívida só aparece depois de reprovar.
//
// Com `derived.rule_marking: required` declarado, o passo vira obrigação — e o
// procedimento passa a ensinar ANTES o que o gate cobra DEPOIS.
func marcacaoDaRegra(cfg *config.Config) string {
	if cfg != nil && cfg.Derived != nil &&
		strings.EqualFold(strings.TrimSpace(cfg.Derived.RuleMarking), "required") {
		return "Implemente cada regra da spec e MARQUE no código o trecho que a realiza " +
			"(`// {CODE}-B01: …`) — este projeto exige a marcação, e o gate " +
			"`regra-implementada` a confronta. O que não tiver código, dispense na linha " +
			"da regra (`@no-code: <razão>`)."
	}
	return "Implemente cada regra da spec; cite o código da regra (`{CODE}-B01`) no " +
		"comentário do trecho que a realiza, se o projeto usa esse padrão."
}
