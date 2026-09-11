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

// O ARQUIVO ATRASADO é a terceira hipótese, e faltava — era justamente a que mais acontece.
//
// Medido num PR do projeto de referência: o branch trazia `trinca_opcional` (formato 1), o
// CI rodava o binário novo, e a mensagem mandava ATUALIZAR O BINÁRIO. O conselho era o
// INVERSO do conserto — quem o seguisse ficaria voltando versão até desistir.
//
// A hipótese é DECIDÍVEL, não adivinhada: o `version:` do arquivo responde sozinho.
func TestLoad_arquivoAntigoComChaveRenomeadaMandaMigrar(t *testing.T) {
	// O registro de migração é injetado pelo `main`; aqui o teste o simula, porque é o
	// que separa "chave renomeada" de "typo".
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(k string) bool { return k == "trinca_opcional" }

	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	os.WriteFile(p, []byte("version: 1\nlayers:\n    spec:\n        pattern: \"**/*.spec.md\"\n"+
		"        kind: spec\n        trinca_opcional: [tested-by]\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("a chave antiga deveria ser recusada — ela foi renomeada")
	}
	msg := err.Error()
	if !strings.Contains(msg, "anchors migrate") {
		t.Errorf("a mensagem deveria mandar MIGRAR; veio:\n%s", msg)
	}
	// E NÃO deve mandar atualizar o binário: é o conselho inverso, e foi o que aconteceu.
	if strings.Contains(msg, "go install") {
		t.Errorf("o binário está certo — mandar atualizá-lo é o conselho INVERSO:\n%s", msg)
	}
}

// AS DUAS CONDIÇÕES são necessárias. Só o formato não basta: um arquivo antigo com um typo
// de verdade receberia "rode `anchors migrate`" — o comando roda, não conserta o typo, e
// quem lê perde a confiança na mensagem seguinte.
func TestLoad_arquivoAntigoComTypoNaoMandaMigrar(t *testing.T) {
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(k string) bool { return k == "trinca_opcional" }

	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	os.WriteFile(p, []byte("version: 1\nlayers: {}\nlayerz: 1\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("o typo deveria ser recusado")
	}
	msg := err.Error()
	if strings.Contains(msg, "anchors migrate") {
		t.Errorf("a migração não conserta typo — mandar migrar gasta a confiança:\n%s", msg)
	}
	if !strings.Contains(msg, "escrita errada") {
		t.Errorf("o typo deveria receber a hipótese da digitação; veio:\n%s", msg)
	}
}

// Sem o registro injetado (nil), a mensagem cai na versão genérica — o comportamento certo
// para quem não tem a tabela de migração à mão.
func TestLoad_semRegistroDeMigracaoCaiNaMensagemGenerica(t *testing.T) {
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = nil

	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	os.WriteFile(p, []byte("version: 1\nlayers: {}\nqualquerCoisa: 1\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("esperava erro")
	}
	if strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("sem registro não há como afirmar que a chave foi renomeada:\n%s", err)
	}
}

// Um arquivo JÁ no formato atual com chave desconhecida é typo ou binário velho — nunca
// "arquivo antigo". A condição do formato é o que separa.
func TestLoad_arquivoNoFormatoAtualNuncaMandaMigrar(t *testing.T) {
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(string) bool { return true } // diria "renomeada" para tudo

	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	os.WriteFile(p, []byte("version: 2\nlayers: {}\ntrinca_opcional: 1\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("esperava erro")
	}
	if strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("o arquivo já está no formato atual — migrar não muda nada:\n%s", err)
	}
}

// ARQUIVO SEM `version:` é formato 1, não formato atual.
//
// Os primeiros `anchors.yaml` do produto não declaravam o campo. Tratá-los como "já
// atualizados" os deixaria com a mensagem genérica — mandando atualizar um binário que já
// é o mais novo, quando o conserto é migrar.
//
// A mutação que devolve `FormatoAtualDeConfig` na ausência do campo passava por todos os
// outros testes: eles declaram `version:` sempre.
func TestLoad_semCampoVersionEhFormatoUm(t *testing.T) {
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(k string) bool { return k == "trinca_opcional" }

	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.yaml")
	// SEM `version:` — como os primeiros arquivos do produto.
	os.WriteFile(p, []byte("layers:\n    spec:\n        pattern: \"**/*.spec.md\"\n"+
		"        kind: spec\n        trinca_opcional: [tested-by]\n"), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("a chave antiga deveria ser recusada")
	}
	if !strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("arquivo sem `version:` é formato 1 e precisa MIGRAR; veio:\n%s", err)
	}
}
