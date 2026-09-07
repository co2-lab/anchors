package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func noDeSpecAutocontida() mapx.Node { return mapx.Node{Kind: mapx.KindSpec, ID: "x/Y.spec.md"} }

// grafo com um plano e a própria spec — é dele que sai o que conta como referência.
func grafoComPlano() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{
		{ID: "plans/0014-alertas.md", Kind: mapx.KindPlan, Code: "ALRTS"},
		{ID: "x/Y.spec.md", Kind: mapx.KindSpec, Code: "XPTOX"},
	}}
}

// A REFERÊNCIA COM O TEXTO passa. É a distinção que o gate inteiro faz: com o trecho
// junto, o leitor tem o argumento em mãos e não precisa sair da página.
func TestDocSelfContained_referenciaComTextoPassa(t *testing.T) {
	c := "# X\n\n## Regras\n\nO documento `plans/0014-alertas.md` diz o que torna " +
		"este canal diferente:\n\n> *\"O desenho evita expor rota de ingestão.\"*\n\n" +
		"E é o que esta regra herda.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a referência traz o texto junto", v, msg)
	}
}

// A referência que SÓ APONTA é acusada: manda o leitor a um arquivo que ele não tem.
func TestDocSelfContained_referenciaVaziaFalha(t *testing.T) {
	casos := []struct{ nome, linha string }{
		{"caminho completo", "O critério está em `plans/0014-alertas.md` e vale aqui."},
		{"caminho sem crase", "O critério está em plans/0014-alertas.md e vale aqui."},
		{"só o nome do arquivo", "Ver 0014-alertas.md para o contexto."},
		{"revisão", "A decisão mudou na `PLTFR-R0005`."},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			c := "# X\n\n## Regras\n\n" + caso.linha + "\n"
			v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil)
			if v != Fail {
				t.Fatalf("verdict = %v — a referência não traz o texto", v)
			}
			if !strings.Contains(msg, "TRAGA O TEXTO") {
				t.Errorf("a mensagem não diz o conserto: %s", msg)
			}
		})
	}
}

// SEM MATCH DE IDIOMA. Um gate que procurasse "plano" ou "ver" passaria em silêncio numa
// spec em inglês — e silêncio é pior que ausência, porque a spec parece protegida.
func TestDocSelfContained_pegaEmQualquerIdioma(t *testing.T) {
	casos := map[string]string{
		"inglês":    "The decision lives in `plans/0014-alertas.md` and holds here.",
		"espanhol":  "El criterio está en `plans/0014-alertas.md` y vale aquí.",
		"alemão":    "Die Entscheidung steht in `plans/0014-alertas.md`.",
		"japonês":   "詳細は `plans/0014-alertas.md` を参照。",
		"sem prosa": "`plans/0014-alertas.md`",
	}
	for nome, linha := range casos {
		t.Run(nome, func(t *testing.T) {
			c := "# X\n\n## Regras\n\n" + linha + "\n"
			if v, _ := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Fail {
				t.Errorf("verdict = %v — o gate não pode depender do idioma", v)
			}
		})
	}
}

// E a CITAÇÃO também é reconhecida em qualquer idioma: um projeto em francês cita com « »
// e um em alemão com „ “ — um gate que só conhece `"` acusaria os dois injustamente.
func TestDocSelfContained_reconheceAspasDeQualquerIdioma(t *testing.T) {
	casos := map[string]string{
		"aspas retas":         `O documento diz "o desenho evita expor rota".`,
		"aspas tipográficas":  `O documento diz “o desenho evita expor rota”.`,
		"guillemets":          `Le document dit « le design évite d'exposer la route ».`,
		"aspas alemãs":        `Das Dokument sagt „das Design vermeidet die Route".`,
		"cantoneiras (ja/zh)": `文書には「設計はルートの公開を避ける」とある。`,
	}
	for nome, linha := range casos {
		t.Run(nome, func(t *testing.T) {
			c := "# X\n\n## Regras\n\n`plans/0014-alertas.md`: " + linha + "\n"
			if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
				t.Errorf("verdict = %v (%s) — há citação, e o gate não pode exigir aspas latinas",
					v, msg)
			}
		})
	}
}

// Dentro de bloco de código o caminho é EXEMPLO — um comando a rodar, não uma referência
// que manda o leitor embora.
func TestDocSelfContained_ignoraBlocoDeCodigo(t *testing.T) {
	c := "# X\n\n## Regras\n\nRode:\n\n```\nanchors check --changed plans/0014-alertas.md\n```\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — dentro da cerca é exemplo", v, msg)
	}
}

// A spec que cita o CAMINHO DE SI MESMA está se identificando, não mandando ninguém a
// lugar nenhum.
func TestDocSelfContained_oProprioCaminhoNaoEhReferencia(t *testing.T) {
	c := "# X\n\n## Regras\n\nEsta unidade vive em `x/Y.spec.md` e governa o módulo.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — é o caminho dela mesma", v, msg)
	}
}

// Um código de regra (`GLCGL-B01`) não é revisão: a forma `-R000N` é a que a doutrina
// reserva, e acusar a outra cobraria a spec por nomear o próprio assunto.
func TestDocSelfContained_codigoDeRegraNaoEhRevisao(t *testing.T) {
	c := "# X\n\n## Regras\n\nA GLCGL-B01 já provou o campo obrigatório, e esta regra o herda.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — `GLCGL-B01` é regra, não revisão", v, msg)
	}
}

// Só a SPEC é cobrada: o plano referencia outros planos por função, e a feature não vira
// prosa de documentação.
func TestDocSelfContained_soASpec(t *testing.T) {
	c := "A decisão está em `plans/0014-alertas.md`."
	for _, k := range []mapx.Kind{mapx.KindPlan, mapx.KindFeature, mapx.KindTest, mapx.KindCode} {
		if v, _ := checkDocSelfContained(c, mapx.Node{Kind: k}, "", grafoComPlano(), nil); v != Skip {
			t.Errorf("kind %s: verdict = %v, queria Skip", k, v)
		}
	}
}
