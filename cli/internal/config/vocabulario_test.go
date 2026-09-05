package config

import (
	"os"
	"path/filepath"
	"testing"
)

// O NOME ANTIGO AINDA CARREGA — e este é o teste que impede quebrar quem já usa.
//
// Os nomes de gate foram do português para o inglês. Um projeto que declarou
// `codigo-catalogado` no `anchors.yaml` não pode parar de funcionar porque atualizou o
// binário: ele migra quando quiser.
func TestVocabulario_nomeAntigoCarregaEViraCanonico(t *testing.T) {
	dir := t.TempDir()
	yaml := `version: 1
layers:
    spec:
        pattern: '**/*.spec.md'
        kind: spec
gates:
    - name: codigo-catalogado
      on: [spec]
      check: codigo-catalogado
      blocking: true
      measures: 'todo símbolo tem regra'
`
	p := filepath.Join(dir, DefaultFile)
	if err := os.WriteFile(p, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("um anchors.yaml com o nome ANTIGO deixou de carregar: %v", err)
	}
	if len(cfg.Gates) != 1 {
		t.Fatalf("gates: %d", len(cfg.Gates))
	}
	g := cfg.Gates[0]
	if g.Name != "code-cataloged" {
		t.Errorf("Name = %q, esperava o canônico 'code-cataloged'", g.Name)
	}
	if g.Check != "code-cataloged" {
		t.Errorf("Check = %q, esperava o canônico", g.Check)
	}
	// O ID acompanha: deixá-lo no nome antigo faria o `check` reportar um nome e o mapa
	// gravar outro, e o carimbo nunca casaria.
	if g.ID != "code-cataloged" {
		t.Errorf("ID = %q, esperava acompanhar o nome", g.ID)
	}
}

// O NOME NOVO carrega igual — senão a migração seria um caminho só de ida sem destino.
func TestVocabulario_nomeNovoCarregaDireto(t *testing.T) {
	dir := t.TempDir()
	yaml := `version: 1
layers:
    spec:
        pattern: '**/*.spec.md'
        kind: spec
gates:
    - name: code-cataloged
      on: [spec]
      check: code-cataloged
      blocking: true
      measures: 'every symbol has a rule'
`
	p := filepath.Join(dir, DefaultFile)
	if err := os.WriteFile(p, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Gates[0].Name != "code-cataloged" {
		t.Errorf("Name = %q", cfg.Gates[0].Name)
	}
}

// Nenhum nome antigo é igual ao novo: uma entrada assim é ruído que faz o `doctor`
// avisar sobre um nome que já está certo.
func TestVocabulario_nenhumaEntradaEhIdentidade(t *testing.T) {
	for a, n := range nomesAntigos {
		if a == n {
			t.Errorf("%q mapeia para si mesmo", a)
		}
	}
	for a, n := range checksAntigos {
		if a == n {
			t.Errorf("check %q mapeia para si mesmo", a)
		}
	}
}

// O DOCTOR precisa ver o que migrar — e lê o arquivo CRU, porque a carga já converteu.
func TestVocabulario_obsoletoEhDetectadoNoArquivoCru(t *testing.T) {
	yaml := "gates:\n  - name: trinca-completa\n    check: mock-carimbado\n"
	achados := VocabularioObsoleto(yaml)
	if len(achados) != 2 {
		t.Fatalf("achou %d obsoletos, esperava 2: %v", len(achados), achados)
	}
	if achados["trinca-completa"] != "triad-complete" {
		t.Errorf("trinca-completa → %q", achados["trinca-completa"])
	}
}

func TestVocabulario_arquivoJaMigradoNaoAcusaNada(t *testing.T) {
	yaml := "gates:\n  - name: triad-complete\n    check: triad-complete\n"
	if a := VocabularioObsoleto(yaml); len(a) != 0 {
		t.Errorf("arquivo já migrado acusou %v", a)
	}
}

// Um nome desconhecido passa intacto: o projeto pode ter gates próprios, e converter o
// que não está no de-para renomearia gate de terceiro.
func TestVocabulario_nomeDesconhecidoPassaIntacto(t *testing.T) {
	if n, mudou := CanonicalizaNome("meu-gate-proprio"); mudou || n != "meu-gate-proprio" {
		t.Errorf("nome próprio foi convertido: %q (mudou=%v)", n, mudou)
	}
}
