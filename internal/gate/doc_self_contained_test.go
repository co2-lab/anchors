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
	t.Run("DSCDC-B01: A reference that brings the passage it announces passes", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nO documento `plans/0014-alertas.md` diz o que torna " +
		"este canal diferente:\n\n> *\"O desenho evita expor rota de ingestão.\"*\n\n" +
		"E é o que esta regra herda.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a referência traz o texto junto", v, msg)
	}
}

// A referência que SÓ APONTA é acusada: manda o leitor a um arquivo que ele não tem.
func TestDocSelfContained_referenciaVaziaFalha(t *testing.T) {
	t.Run("DSCDC-B02: A reference that only points is accused", func(t *testing.T) {})
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
			if !strings.Contains(msg, "TRAGA O TEXTO") && !strings.Contains(msg, "BRING THE TEXT") {
				t.Errorf("a mensagem não diz o conserto: %s", msg)
			}
		})
	}
}

// SEM MATCH DE IDIOMA. Um gate que procurasse "plano" ou "ver" passaria em silêncio numa
// spec em inglês — e silêncio é pior que ausência, porque a spec parece protegida.
func TestDocSelfContained_pegaEmQualquerIdioma(t *testing.T) {
	t.Run("DSCDC-I01: The ruler matches structure, never vocabulary", func(t *testing.T) {})
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
	t.Run("DSCDC-B03: A quotation counts in any written tradition", func(t *testing.T) {})
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
	t.Run("DSCDC-B04: A path inside a code fence is an example, not a reference", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nRode:\n\n```\nanchors check --changed plans/0014-alertas.md\n```\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — dentro da cerca é exemplo", v, msg)
	}
}

// A spec que cita o CAMINHO DE SI MESMA está se identificando, não mandando ninguém a
// lugar nenhum.
func TestDocSelfContained_oProprioCaminhoNaoEhReferencia(t *testing.T) {
	t.Run("DSCDC-B05: The spec citing its own path is identifying itself", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nEsta unidade vive em `x/Y.spec.md` e governa o módulo.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — é o caminho dela mesma", v, msg)
	}
}

// Um código de regra (`GLCGL-B01`) não é revisão: a forma `-R000N` é a que a doutrina
// reserva, e acusar a outra cobraria a spec por nomear o próprio assunto.
func TestDocSelfContained_codigoDeRegraNaoEhRevisao(t *testing.T) {
	t.Run("DSCDC-B06: A rule code is not a revision", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nA GLCGL-B01 já provou o campo obrigatório, e esta regra o herda.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — `GLCGL-B01` é regra, não revisão", v, msg)
	}
}

// Só a SPEC é cobrada: o plano referencia outros planos por função, e a feature não vira
// prosa de documentação.
func TestDocSelfContained_soASpec(t *testing.T) {
	t.Run("DSCDC-B07: Only the spec is charged", func(t *testing.T) {})
	c := "A decisão está em `plans/0014-alertas.md`."
	for _, k := range []mapx.Kind{mapx.KindPlan, mapx.KindFeature, mapx.KindTest, mapx.KindCode} {
		if v, _ := checkDocSelfContained(c, mapx.Node{Kind: k}, "", grafoComPlano(), nil); v != Skip {
			t.Errorf("kind %s: verdict = %v, queria Skip", k, v)
		}
	}
}

// A REVISÃO com a EXPLICAÇÃO na própria linha passa: o código é o rótulo e a frase é o
// conteúdo. Exigir a forma da citação obrigaria a escrever pior para passar — 14 dos 35
// achados da primeira versão eram assim, acusados só por não usarem aspas.
func TestDocSelfContained_revisaoComExplicacaoNaLinhaPassa(t *testing.T) {
	t.Run("DSCDC-B08: A revision cited with an explanation on the same line passes", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nA `PLTFR-R0003` corrigiu o escopo desta spec: o contador de " +
		"conexões saiu, porque a fonte não o expõe por réplica.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a linha diz O QUE a revisão mudou", v, msg)
	}
	// A contraparte que prova que a régua é a explicação, e não a mera presença da
	// revisão: o rótulo curto, sem dizer o que mudou, continua acusado.
	nu := "# X\n\n## Regras\n\nDecidido pelo usuário (`PLTFR-R0003`).\n"
	if v, _ := checkDocSelfContained(nu, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Fail {
		t.Errorf("verdict = %v — o rótulo não diz O QUE foi decidido", v)
	}
}

// SEM MAPA não há como saber quais caminhos são nós. A lista do que conta como referência
// sai do mapa e de nenhum outro lugar — inventar um padrão `plans/*.md` à mão só valeria
// para projetos que chamam a pasta de `plans`.
func TestDocSelfContained_semMapaEhSkip(t *testing.T) {
	t.Run("DSCDC-B09: With no map the confrontation is skipped", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nO critério está em `plans/0014-alertas.md` e vale aqui.\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", nil, nil); v != Skip {
		t.Errorf("verdict = %v (%s) — sem mapa não há lista de nós", v, msg)
	}
}

// A ASSIMETRIA é deliberada: o escape da explicação na linha vale para a REVISÃO, não
// para o CAMINHO. "O critério está em `plans/0014.md` e vale para todo handler novo" é
// uma linha longa que não diz QUAL é o critério — o caminho não é rótulo de nada, é o
// lugar aonde a pessoa teria de ir.
func TestDocSelfContained_prosaLongaNaoSalvaOCaminho(t *testing.T) {
	t.Run("DSCDC-I02: The on-line explanation escape belongs to the revision, not to the path", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nO critério de alarme está definido em `plans/0014-alertas.md` " +
		"e vale para todo handler novo deste módulo, sem exceção documentada até aqui.\n"
	v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s) — a prosa é longa e não diz QUAL é o critério", v, msg)
	}
	// E a prova de que a assimetria é o que separa os dois: a MESMA quantidade de prosa
	// em volta de uma revisão passa.
	comRevisao := "# X\n\n## Regras\n\nO critério de alarme foi fixado pela `PLTFR-R0003` " +
		"e vale para todo handler novo deste módulo, sem exceção documentada até aqui.\n"
	if v, _ := checkDocSelfContained(comRevisao, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v — a revisão é rótulo, e a frase em volta é o conteúdo", v)
	}
}

// O veredito NOMEIA a linha e mostra o que ela diz. Um gate que acusa sem apontar
// transfere o trabalho de diagnóstico para quem lê.
func TestDocSelfContained_oAchadoApontaALinha(t *testing.T) {
	t.Run("DSCDC-I03: The verdict names the line and shows what it says", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nVer `plans/0014-alertas.md`.\n"
	v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s)", v, msg)
	}
	// A referência vazia está na 5ª linha (1-based) do conteúdo acima.
	if !strings.Contains(msg, "5") {
		t.Errorf("o achado devia trazer o NÚMERO da linha; msg = %q", msg)
	}
	if !strings.Contains(msg, "Ver `plans/0014-alertas.md`.") {
		t.Errorf("o achado devia mostrar o TRECHO da linha; msg = %q", msg)
	}
}

// O gate não julga se o conteúdo acompanhante é FIEL ao que a referência anuncia. A régua
// é se o leitor fica com algo para ler — julgar fidelidade é outra classe de gate.
func TestDocSelfContained_naoJulgaFidelidadeDoTrecho(t *testing.T) {
	t.Run("DSCDC-X01: The gate does not judge whether the accompanying content is faithful", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nO documento `plans/0014-alertas.md` fixa o limiar de alarme:\n\n" +
		"> *\"O relatório mensal sai no primeiro dia útil.\"*\n"
	if v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a citação não bate com o anúncio, e isso é outra régua",
			v, msg)
	}
}

// INFORMATIVO: o gate marca e não oferece comando que conserte. A mensagem tem de dizer o
// CONSERTO em prosa — "traga o texto" —, porque reescrever a frase é trabalho de quem a
// escreveu, e nunca mandar rodar um comando que geraria o texto no lugar dela.
func TestDocSelfContained_marcaSemOferecerComando(t *testing.T) {
	t.Run("DSCDC-X02: The gate marks and does not block", func(t *testing.T) {})
	c := "# X\n\n## Regras\n\nO critério está em `plans/0014-alertas.md`.\n"
	v, msg := checkDocSelfContained(c, noDeSpecAutocontida(), "", grafoComPlano(), nil)
	if v != Fail {
		t.Fatalf("verdict = %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "TRAGA O TEXTO") && !strings.Contains(msg, "BRING THE TEXT") {
		t.Errorf("a mensagem devia dizer o conserto em prosa; msg = %q", msg)
	}
	if strings.Contains(msg, "anchors ") {
		t.Errorf("não há comando que conserte — reescrever a frase é de quem a escreveu; msg = %q", msg)
	}
}

// A régua da explicação erra para o lado de DEIXAR PASSAR. É um gate informativo, e falso
// positivo em massa é o que faz alguém desligá-lo — então a prosa que mal ultrapassa o
// limiar passa, mesmo dizendo pouco.
func TestDocSelfContained_aReguaErraParaODeixarPassar(t *testing.T) {
	t.Run("DSCDC-X03: Measuring explanation errs on the permissive side", func(t *testing.T) {})
	// Prosa logo acima do limiar, e dizendo pouco: passa de propósito.
	acima := "# X\n\n## Regras\n\nA `PLTFR-R0003` mudou o escopo desta seção aqui.\n"
	if v, msg := checkDocSelfContained(acima, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Pass {
		t.Errorf("verdict = %v (%s) — a régua erra para o lado permissivo", v, msg)
	}
	// E o rótulo curto continua acusado: a régua é grosseira, não inexistente.
	abaixo := "# X\n\n## Regras\n\nDecidido (`PLTFR-R0003`).\n"
	if v, _ := checkDocSelfContained(abaixo, noDeSpecAutocontida(), "", grafoComPlano(), nil); v != Fail {
		t.Errorf("verdict = %v — o rótulo curto não é explicação", v)
	}
}
