package flow

import (
	"os"
	"strings"
	"testing"
)

// leFonte lê um arquivo deste pacote para o teste confrontar o que ele declara.
//
// É o recurso de último caso: prefira medir COMPORTAMENTO. Aqui ele se justifica porque o
// comando conversa com a API do GitHub — exercitá-lo exigiria um repositório real, e o que
// se guarda é quais decisões o código toma, não a saída dele.
func leFonte(t *testing.T, nome string) string {
	t.Helper()
	b, err := os.ReadFile(nome)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// O BACKFILL NÃO INVENTA VÍNCULO.
//
// Uma label nova não alcança o passado: os cards que pararam antes de o `blocked-by-<n>`
// existir só têm o `needs-user`, e o board diz que esperam sem dizer por quem.
//
// O vínculo é RECUPERÁVEL porque o `escalate` sempre escreveu o número do card de origem na
// label `under-<n>` do card novo — a relação existe na direção contrária.
//
// O QUE NÃO PODE: adivinhar. Card em `needs-user` sem ninguém apontando para ele fica como
// está. Inventar um bloqueador é pior que não ter a label — o `claim` passaria a segurar um
// card por causa de uma decisão que ninguém ligou a ele, e o card ficaria parado sem que
// resolver a decisão o solte.
//
// MEDIDO no projeto de referência: 29 decisões abertas, 5 vínculos recuperáveis, 4 pulados.
// Os 4 pulados são o comportamento certo, não uma falha de cobertura.
func TestBackfillRecuperaOVinculoSemAdivinhar(t *testing.T) {
	cmd := newBackfillLabelsCmd()

	// O DRY-RUN tem de existir: escrever label em massa num board de produção sem poder
	// ver antes é o tipo de comando que se roda uma vez e se lamenta.
	if cmd.Flags().Lookup("dry-run") == nil {
		t.Error("falta `--dry-run`: este comando escreve em N issues de uma vez, e quem o " +
			"roda precisa ver o que ele faria antes de fazer")
	}

	// A FONTE do vínculo é a label `under-<n>`, e o legado dela.
	//
	// Confronta o FONTE porque o comando fala com a API do GitHub: um teste de
	// comportamento exigiria um repositório real. O que se guarda aqui é que as duas
	// grafias são lidas — a migração para inglês não terminou nos boards que já existem, e
	// ler só a nova perderia silenciosamente todo card anterior a ela.
	// A GRAFIA É UMA SÓ. Houve aqui a exigência de ler também a anterior (`sob-<n>`), e
	// ela saiu: o projeto está em beta fechado, não há board de terceiro a respeitar, e
	// aceitar duas formas da mesma coisa é superfície de divergência — a primeira vez que
	// uma delas mudasse, a outra ficaria para trás sem ninguém notar.
	fonte := leFonte(t, "backfill_labels.go")
	if !strings.Contains(fonte, "PrefixoLabelSob") {
		t.Error("o backfill não lê `PrefixoLabelSob` — é de onde o vínculo sai, e sem ele " +
			"o comando não tem o que recuperar")
	}
	if strings.Contains(fonte, "PrefixoLabelSobAntigo") {
		t.Error("o backfill voltou a ler a grafia anterior — uma forma só, e a migração " +
			"não tem board de terceiro para respeitar")
	}

	// A DECISÃO NÃO SEGURA A SI MESMA: um card pode ser decisão e origem de outro achado,
	// e ligá-lo a si o tornaria eternamente bloqueado.
	if !strings.Contains(fonte, "origem != n") {
		t.Error("o backfill pode ligar um card a si mesmo — um card que é decisão E origem " +
			"de outro achado ficaria bloqueado por si, sem nada que o solte")
	}

	// SÓ CARD ABERTO recebe o bloqueio: pôr label em card fechado não muda nada.
	if !strings.Contains(fonte, `"OPEN"`) {
		t.Error("o backfill não confere se o card de origem está aberto — rotular fechado " +
			"não muda nada e enche o histórico")
	}
}

// O DESBLOQUEIO NÃO APAGA A HISTÓRIA.
//
// `blocked-by-<n>` sai junto com o `needs-user` — senão fica pendurada num card já livre,
// e quem lê o board depois vê bloqueio que não existe mais.
//
// MAS A LABEL É O ÚNICO LUGAR ONDE O VÍNCULO ESTÁ. Removê-la sem registrar apaga a resposta
// de "por que isto ficou parado três dias?" e o caminho de volta para a decisão que
// segurava. O comentário é imutável, datado, e sobrevive a qualquer mexida posterior nas
// labels — o contrário da label, que é estado do AGORA.
func TestDesbloqueioRegistraOVinculoAntesDeRemover(t *testing.T) {
	fonte := leFonte(t, "decided.go")

	// A LEITURA vem antes da remoção: depois, a informação não existe em lugar nenhum.
	iLe := strings.Index(fonte, "labelsDeBloqueio(cfg.Workflow.Repo, card)")
	iRemove := strings.Index(fonte, `"--remove-label"`)
	if iLe < 0 {
		t.Fatal("o `decided` não lê as labels de bloqueio — elas ficariam penduradas num " +
			"card já livre, e o board mostraria bloqueio que não existe mais")
	}
	if iRemove < 0 {
		t.Fatal("não achei a remoção de label no `decided` — o comando mudou de forma")
	}
	if iLe > iRemove {
		t.Error("o `decided` lê as labels de bloqueio DEPOIS de removê-las: nesse ponto a " +
			"informação já não existe, e o rastro sai vazio")
	}

	// O RASTRO em comentário, com os números: dizer "estava bloqueado" sem dizer por quem
	// registra que houve espera e perde por onde rastrear a decisão para trás.
	if !strings.Contains(fonte, "**Was blocked by:**") {
		t.Error("o `decided` remove o bloqueio sem registrar por quem o card esperou — " +
			"some a resposta de por que ele ficou parado")
	}
}

// PRECEDÊNCIA ENTRE DECISÕES NÃO SE INFERE — mas PODE SER DECLARADA.
//
// A primeira versão desta régua proibia `needs-user` + `blocked-by-<n>` no mesmo card, e
// estava errada: as duas juntas são LEGÍTIMAS. Uma decisão pode depender de outra, e isso é
// ordem de corretude — "responda o #701 primeiro, porque a resposta dele condiciona a do
// #647". Proibir apagaria informação real.
//
// O DEFEITO ERA INFERIR. O `backfill` lê o `under-<n>`, que significa "nasci do trabalho do
// #n" — PROCEDÊNCIA. Quando o #n é um card de trabalho, o vínculo é inequívoco: o trabalho
// para até a decisão sair. Quando o #n é outra decisão, `under` não distingue "uma
// condiciona a outra" de "as duas nasceram no mesmo lugar", e são perguntas diferentes.
//
// MEDIDO no projeto de referência: quatro cards (#647, #706, #721, #555) receberam bloqueio
// por essa inferência, e nos quatro o bloqueador era decisão que apenas nasceu junto. O
// único que o `under` acertou foi o #198 → #647 — card em `in-review`, trabalho de verdade.
//
// A ASSIMETRIA entre os dois comandos é o ponto:
//
//	backfill    INFERE do `under-<n>` → recupera só o caso inequívoco (origem não escalada)
//	escalate    DECLARA no ato       → vale sempre, inclusive decisão que precede decisão
//
// Quem quiser ordem entre duas decisões existentes põe a label à mão: é julgamento de quem
// leu as duas, igual a decidir se são o mesmo assunto.
func TestBackfillNaoInfereprecedenciaEntreDecisoes(t *testing.T) {
	// A GUARDA vive no `backfill`, que infere. O `escalate` NÃO pode tê-la.
	if g := "jaTem[initx.LabelNeedsUser]"; !strings.Contains(leFonte(t, "backfill_labels.go"), g) {
		t.Errorf("o backfill perdeu a guarda %q — ele voltaria a inferir precedência do "+
			"`under-<n>`, que diz procedência: medido, erra em 4 de 5 casos", g)
	}
	if strings.Contains(leFonte(t, "escalate.go"), "!jaEra") {
		t.Error("o `escalate` ganhou a guarda do backfill — ali o vínculo é DECLARADO no " +
			"ato por quem escalou, e suprimi-lo apagaria a ordem que ele estabeleceu " +
			"quando a decisão nova precede uma que já estava aberta")
	}
}

// O ACHADO NASCE COM PROCEDÊNCIA, e quando não nasce o comando DIZ.
//
// O `--card` era opcional e silencioso. Sem ele o achado não tem de onde veio, e ninguém
// consegue perguntar "o que este trabalho gerou?" — a doutrina do `under-<n>` diz que o
// vínculo é LABEL justamente porque frase no corpo não se consulta, e sem o `--card` não há
// nem frase.
//
// MEDIDO no projeto de referência: 18 de 27 decisões abertas sem amarração alguma. Todas
// traziam `**Onde:** <arquivo>` — procedência pelo ARQUIVO, que responde onde o problema
// está e não de qual trabalho ele saiu. Uma chegou a escrever "o achado veio do #690" em
// prosa, exatamente o que a doutrina condena.
//
// A ASSIMETRIA não tinha razão: o `pr-body` já descobria o card pelo `ANCHORS_AGENT`, e o
// `escalate` não. Os dois perguntam a mesma coisa — "qual card este agente pegou?".
func TestEscalateDescobreOCardDeOrigem(t *testing.T) {
	fonte := leFonte(t, "escalate.go")

	// A DESCOBERTA reusa a peça do `pr-body`: duas implementações da mesma pergunta
	// divergiriam na primeira vez que uma delas mudasse.
	if !strings.Contains(fonte, `requestedCards("", cfg)`) {
		t.Error("o `escalate` não descobre o card de origem — o achado nasce sem " +
			"procedência, e o `pr-body` já sabia fazer isso pelo `anchors-owner`")
	}

	// AMBIGUIDADE recusa em vez de chutar: com dois cards em mãos, só o agente sabe em
	// qual estava mexendo. Escolher o primeiro amarraria o achado ao trabalho errado.
	if !strings.Contains(fonte, "provide `--card <n>`") {
		t.Error("o `escalate` não trata o caso de vários cards em mãos — escolher um por " +
			"conta própria amarra o achado ao trabalho errado, e o vínculo errado é pior " +
			"que vínculo nenhum")
	}

	// E O SILÊNCIO é o que se corrige: nascer solto é legítimo (alguém pode escalar fora
	// de um trabalho), acontecer sem aviso não.
	//
	// A primeira versão desta asserção tinha um `||` que a tornava sempre verdadeira —
	// uma das metades casava a palavra "procedência" em qualquer lugar do arquivo,
	// inclusive nos comentários que eu mesmo acabara de escrever. Ela passava sem medir.
	if !strings.Contains(fonte, "anchors-owner` — the finding is born WITHOUT ") {
		t.Error("o `escalate` aceita achado sem procedência em silêncio — quem o rodou " +
			"não fica sabendo que abriu um card que ninguém vai poder rastrear")
	}
}
