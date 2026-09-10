package config

import (
	"os"
	"path/filepath"
	"testing"
)

// O merge existe para que a redação de um gate canônico viva em UM lugar. Sem ele, o
// projeto precisa repetir `check`/`measures`/`ask` inteiros, a cópia diverge no primeiro
// ajuste do framework, e os projetos seguem perguntando a versão velha sem ninguém ver.
func TestLoadMergeCanonicalPreencheOmitidos(t *testing.T) {
	SetCanonicalGateResolver(func(name string) (Gate, bool) {
		if name != "gate-canonico" {
			return Gate{}, false
		}
		return Gate{
			Name: "gate-canonico", On: []string{"spec"}, Check: "algo",
			Measures: "o que ele mede", Ask: "a pergunta canônica",
			Scope: ScopeBatch, Cost: "slow", Category: "test",
			When: []string{"ci"},
		}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	p := filepath.Join(t.TempDir(), "anchors.yaml")
	// só o NOME: tudo o mais deve vir do canônico.
	if err := os.WriteFile(p, []byte("version: 1\ngates:\n  - name: gate-canonico\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	g := c.Gates[0]
	if g.Check != "algo" || g.Measures != "o que ele mede" || g.Ask != "a pergunta canônica" {
		t.Errorf("campos textuais não herdaram: %+v", g)
	}
	if g.EffectiveScope() != ScopeBatch || g.Cost != "slow" || g.Category != "test" {
		t.Errorf("escopo/custo/categoria não herdaram: %+v", g)
	}
	if len(g.On) != 1 || g.On[0] != "spec" || len(g.When) != 1 {
		t.Errorf("listas não herdaram: %+v", g)
	}
}

// O projeto SEMPRE vence onde declara — o merge preenche buraco, não sobrescreve. Sem
// isto, customizar um gate canônico seria impossível.
func TestLoadMergeCanonicalNaoSobrescreveODeclarado(t *testing.T) {
	SetCanonicalGateResolver(func(string) (Gate, bool) {
		return Gate{Name: "g", Measures: "canônico", Ask: "pergunta canônica"}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p,
		[]byte("version: 1\ngates:\n  - name: g\n    measures: \"do projeto\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if c.Gates[0].Measures != "do projeto" {
		t.Errorf("o declarado no projeto deve vencer: %q", c.Gates[0].Measures)
	}
	if c.Gates[0].Ask != "pergunta canônica" {
		t.Errorf("o omitido ainda deve herdar: %q", c.Gates[0].Ask)
	}
}

// Gate que não é canônico passa intacto — o merge não pode inventar campo para gate
// de projeto (um `run:` herdado de outro gate seria um comando executando sozinho).
func TestLoadMergeCanonicalIgnoraGateDeProjeto(t *testing.T) {
	SetCanonicalGateResolver(func(string) (Gate, bool) { return Gate{}, false })
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p,
		[]byte("version: 1\ngates:\n  - name: meu-gate\n    run: \"bash x.sh\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if c.Gates[0].Check != "" || c.Gates[0].Ask != "" {
		t.Errorf("gate de projeto não deve herdar nada: %+v", c.Gates[0])
	}
}

// O caso que motivou o `*bool`: o canônico é BLOQUEANTE e o projeto declarou
// `blocking: false` de propósito (maturação — a unidade ainda tem débito e o gate não
// pode barrar). Com `bool`, o merge lia `false`, concluía "não declarou" e PROMOVIA o
// gate a bloqueante: commits que o autor liberou conscientemente passavam a ser barrados,
// sem nada no yaml explicando por quê.
func TestLoadMergeCanonicalBlockingFalseExplicitoVence(t *testing.T) {
	SetCanonicalGateResolver(func(string) (Gate, bool) {
		return Gate{Name: "g", Blocking: Bool(true)}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p,
		[]byte("version: 1\ngates:\n  - name: g\n    blocking: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if c.Gates[0].IsBlocking() {
		t.Error("`blocking: false` explícito é DECISÃO e não pode ser sobrescrito pelo canônico")
	}
}

// A contraparte: OMITIDO herda mesmo. É o que distingue "não declarei" de "decidi", e
// sem esta metade o campo simplesmente sairia do merge.
func TestLoadMergeCanonicalBlockingOmitidoHerda(t *testing.T) {
	SetCanonicalGateResolver(func(string) (Gate, bool) {
		return Gate{Name: "g", Blocking: Bool(true)}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 1\ngates:\n  - name: g\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if !c.Gates[0].IsBlocking() {
		t.Error("campo OMITIDO deve herdar do canônico")
	}
}

// Sem declaração em lugar nenhum, o gate é INFORMATIVO: o silêncio não pode barrar
// commit — a maturação é escolha explícita (QUALITY §7).
func TestGateSemBlockingEhInformativo(t *testing.T) {
	if (Gate{Name: "g"}).IsBlocking() {
		t.Error("gate sem severidade declarada não pode barrar")
	}
}
