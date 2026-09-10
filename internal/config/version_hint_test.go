package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Um campo NOVO no config quebra todo projeto cujo binário é anterior a ele, e a mensagem
// do yaml — "field docs not found in type config.Config" — descreve o sintoma em
// vocabulário de implementação: `config.Config` é um tipo Go, e quem lê o `anchors.yaml`
// não o conhece. Ela manda procurar um erro de digitação que pode não existir.
//
// Medido: o `docs:` entrou na v0.1.58 e o hook do projeto governado, com binário anterior,
// barrou o commit com essa mensagem. O arquivo estava certo; a ferramenta é que era velha.
func TestLoad_chaveDesconhecidaNomeiaAHipoteseDaVersao(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	os.WriteFile(p, []byte("version: 1\nlayers: {}\nchaveQueNaoExiste: 1\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("chave desconhecida passou — o KnownFields existe para pegá-la")
	}
	msg := err.Error()
	if !strings.Contains(msg, "chaveQueNaoExiste") {
		t.Errorf("a mensagem não nomeia a chave: %s", msg)
	}
	if !strings.Contains(msg, "ANTIGO") {
		t.Errorf("a mensagem não oferece a hipótese da versão — quem tem binário velho\n"+
			"procura um erro de digitação que não existe:\n%s", msg)
	}
	if !strings.Contains(msg, "go install") {
		t.Errorf("a mensagem diz o problema e não o conserto: %s", msg)
	}
}

// A dica só aparece quando cabe: um YAML mal formado não tem nada a ver com versão, e
// sugerir atualizar o binário ali mandaria o autor pelo caminho errado.
func TestLoad_erroDeSintaxeNaoGanhaADicaDeVersao(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	os.WriteFile(p, []byte("version: 1\n  layers: [\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("YAML quebrado passou")
	}
	if strings.Contains(err.Error(), "ANTIGO") {
		t.Errorf("erro de sintaxe ganhou a dica de versão: %s", err)
	}
}
