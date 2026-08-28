package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// O contrato central de `needs_tool`: ferramenta ausente PULA, não reprova.
//
// A regressão que este teste barra é a que motivou o campo — sem ele o `sh` sai 127 e
// o gate vira Fail, dizendo "o projeto violou algo" quando o que falta é o binário. Um
// Fail aqui reprovaria todo projeto recém-criado que ainda não instalou a ferramenta.
func TestFerramentaAusenteVirandoSkipNaoFail(t *testing.T) {
	g := config.Gate{
		Name: "gate-fantasma", On: []string{"code"},
		Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
		// `false` reprovaria (exit 1) SE chegasse a rodar — é isso que torna o teste
		// conclusivo: só o Skip explica o verde, e não um comando que por acaso passou.
		Run:       "false",
		NeedsTool: "binario-que-nao-existe-em-lugar-nenhum-xyz",
		Blocking:  config.Bool(true),
	}
	nodes := []mapx.Node{{ID: "a.ts", Kind: mapx.Kind("code")}}

	res := RunCompleto([]config.Gate{g}, nodes, t.TempDir(), nil, &config.Config{}, false)
	if len(res) != 1 {
		t.Fatalf("esperava 1 veredito, veio %d", len(res))
	}
	if res[0].Verdict != Skip {
		t.Fatalf("ferramenta ausente devia dar Skip, veio %q (detalhe: %s)", res[0].Verdict, res[0].Detail)
	}
	if !strings.Contains(res[0].Detail, "binario-que-nao-existe") {
		t.Errorf("o laudo deve NOMEAR a ferramenta que falta; veio: %q", res[0].Detail)
	}
}

// A contraparte: com a ferramenta presente o gate roda de verdade. Sem este caso, um
// `needs_tool` que sempre pulasse passaria no teste acima e desligaria o gate em silêncio
// — a falha mais cara possível, porque some com a medição parecendo saudável.
func TestFerramentaPresenteContinuaExecutando(t *testing.T) {
	g := config.Gate{
		Name: "gate-real", On: []string{"code"},
		Scope: config.ScopeProject, ScopeFull: config.ScopeProject,
		Run: "false", // reprova de propósito
		// `sh` existe em qualquer POSIX — é a ferramenta presente mais segura de assumir.
		NeedsTool: "sh",
		Blocking:  config.Bool(true),
	}
	nodes := []mapx.Node{{ID: "a.ts", Kind: mapx.Kind("code")}}

	res := RunCompleto([]config.Gate{g}, nodes, t.TempDir(), nil, &config.Config{}, false)
	if len(res) != 1 {
		t.Fatalf("esperava 1 veredito, veio %d", len(res))
	}
	if res[0].Verdict != Fail {
		t.Fatalf("ferramenta presente devia deixar o gate RODAR (e reprovar), veio %q", res[0].Verdict)
	}
}
