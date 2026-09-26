package telemetry

import (
	"strings"
	"testing"
)

// O OPT-OUT precisa ser fácil de ACERTAR. Quem procura como desligar tenta uma palavra
// óbvia, e exigir a forma exata seria uma pegadinha — a pessoa acharia que desligou.
func TestDesligada_aceitaAsFormasQueAlguemTentaria(t *testing.T) {
	for _, v := range []string{"off", "OFF", "Off", "0", "false", "FALSE", "no", " off "} {
		t.Setenv(EnvVar, v)
		if !Disabled("") {
			t.Errorf("%q deveria desligar — quem escreve isso acha que desligou", v)
		}
	}
}

func TestDesligada_oPadraoEhLigada(t *testing.T) {
	t.Setenv(EnvVar, "")
	if Disabled("") {
		t.Error("sem declaração nenhuma, a telemetry é ligada — é a decisão do mantenedor")
	}
}

// O AMBIENTE VENCE O ARQUIVO, e não é detalhe: quem roda num CI precisa desligar sem
// commitar. Commitar para desligar telemetry faria a decisão de uma pessoa virar mudança
// no repositório do time.
func TestDesligada_oAmbienteVenceOArquivo(t *testing.T) {
	t.Setenv(EnvVar, "on")
	if Disabled("off") {
		t.Error("o ambiente dizendo `on` tem de vencer o arquivo dizendo `off`")
	}
	t.Setenv(EnvVar, "off")
	if !Disabled("on") {
		t.Error("o ambiente dizendo `off` tem de vencer o arquivo dizendo `on`")
	}
}

// Sem o ambiente, o arquivo decide.
func TestDesligada_semAmbienteOArquivoDecide(t *testing.T) {
	t.Setenv(EnvVar, "")
	if !Disabled("off") {
		t.Error("`telemetry: off` no anchors.yaml precisa desligar")
	}
	if Disabled("on") {
		t.Error("`telemetry: on` mantém ligada")
	}
}

// O AVISO aparece UMA VEZ por máquina. Repetir seria pior que não avisar: quem lê a mesma
// coisa toda vez para de ler, e o aviso deixa de cumprir o que justifica o opt-out.
func TestAvisa_umaVezPorMaquina(t *testing.T) {
	dir := t.TempDir()
	var buf strings.Builder

	Notice(&buf, dir)
	if buf.Len() == 0 {
		t.Fatal("a primeira execução precisa avisar")
	}
	primeiro := buf.Len()

	buf.Reset()
	Notice(&buf, dir)
	if buf.Len() != 0 {
		t.Errorf("a segunda execução não deveria avisar; escreveu %d bytes", buf.Len())
	}
	_ = primeiro
}

// O TEXTO diz três coisas, e a terceira é a que torna o opt-out honesto: quem lê
// "coletamos dados" e não encontra COMO DESLIGAR na mesma tela assume o pior.
func TestAvisa_dizOQueColetaOQueNaoEComoDesligar(t *testing.T) {
	var buf strings.Builder
	Notice(&buf, t.TempDir())
	// A comparação normaliza a caixa: o aviso escreve DECISÃO em maiúscula para destacar,
	// e prender a asserção à forma exata testaria a tipografia, não a informação.
	texto := strings.ToLower(buf.String())

	for _, quer := range []string{
		"decisão",               // o que coleta
		"não envia",             // o que não coleta
		"conteúdo de arquivo",   // o medo mais provável de quem lê
		"anchors_telemetry=off", // como desligar, sem editar arquivo
		"telemetry: off",
	} {
		if !strings.Contains(texto, strings.ToLower(quer)) {
			t.Errorf("o aviso deveria conter %q", quer)
		}
	}
}

func TestAvisa_oMarcadorFicaForaDoGit(t *testing.T) {
	if !strings.HasPrefix(noticeFile, ".anchors/") {
		t.Errorf("o marcador precisa ficar em `.anchors/` (não versionado); está em %q",
			noticeFile)
	}
}
