package main

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
	fonte := leFonte(t, "backfill_labels.go")
	for _, peca := range []string{"PrefixoLabelSob", "PrefixoLabelSobAntigo"} {
		if !strings.Contains(fonte, peca) {
			t.Errorf("o backfill não lê %q — o vínculo dos cards mais antigos está na "+
				"grafia anterior, e ignorá-la os deixa sem bloqueio para sempre", peca)
		}
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
	if !strings.Contains(fonte, "Estava bloqueado por:") {
		t.Error("o `decided` remove o bloqueio sem registrar por quem o card esperou — " +
			"some a resposta de por que ele ficou parado")
	}
}
