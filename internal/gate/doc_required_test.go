package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// O CASO QUE MOTIVOU O GATE, medido no projeto de referência:
//
// Um PR entregou uma lambda nova — spec, código, feature e teste — sem tocar o OpenAPI nem
// o esquema de dados. Os dois são declarados com `trigger: [lambdas]`, os dois são
// obrigatórios, e TODOS os checks ficaram verdes.
//
// Só apareceu porque dois agentes colidiram no mesmo card: a branch do segundo tinha 267
// linhas de OpenAPI que a do primeiro não tinha. Sem a colisão, ninguém teria notado — não
// há com que comparar quando só um trabalho existe.

func cfgComDocs() *config.Config {
	return &config.Config{
		// O `path` é o que faz `scan.LayerOfUnit` classificar o arquivo quando ele não tem
		// header — e sem ele a camada sai vazia, o `trigger` não casa, e o gate pula
		// justamente a unidade que deveria cobrar.
		Layers: map[string]config.Layer{"lambdas": {Pattern: "packages/lambdas/**", Kind: "code"}},
		Docs: &config.Docs{Required: []config.DocArtifact{
			{Kind: "openapi", Path: "docs/contratos/openapi.yaml", Trigger: []string{"lambdas"}},
			{Kind: "schema", Path: "docs/contratos/dados.md", Trigger: []string{"lambdas"}},
		}},
	}
}

func noDeLambda() mapx.Node {
	return mapx.Node{
		ID: "packages/lambdas/backend/ServiceList.spec.md", Kind: mapx.KindSpec,
		Layer: "lambdas", Code: "SRLSS",
	}
}

func escreveDoc(t *testing.T, root, rel, conteudo string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDocRequired_documentoAusenteReprova(t *testing.T) {
	root := t.TempDir()
	v, msg := checkDocRequired("", noDeLambda(), root, nil, cfgComDocs())
	if v != Fail {
		t.Fatalf("sem os documentos declarados, esperava Fail; veio %v", v)
	}
	for _, quer := range []string{"openapi.yaml", "dados.md", "docs duties"} {
		if !strings.Contains(msg, quer) {
			t.Errorf("a mensagem deveria nomear %q; veio:\n%s", quer, msg)
		}
	}
}

// A SEGUNDA PERGUNTA é a que importa. Um arquivo vazio criado para calar o gate passaria
// se ele só conferisse existência — e é exatamente o que alguém faz quando o gate barra e
// o prazo aperta.
func TestDocRequired_documentoQueExisteEnaoMencionaReprova(t *testing.T) {
	root := t.TempDir()
	escreveDoc(t, root, "docs/contratos/openapi.yaml", "openapi: 3.0.0\npaths:\n  /outra: {}\n")
	escreveDoc(t, root, "docs/contratos/dados.md", "# Esquema\n\nNada sobre esta unidade.\n")

	v, msg := checkDocRequired("", noDeLambda(), root, nil, cfgComDocs())
	if v != Fail {
		t.Fatalf("documento que existe e não menciona a unidade deveria reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "não menciona") {
		t.Errorf("a mensagem deveria distinguir AUSENTE de MUDO; veio:\n%s", msg)
	}
	// O modo de falha tem de estar DITO: o arquivo existe, tem conteúdo real, e quem o lê
	// não tem como saber que falta uma entrada.
	if !strings.Contains(msg, "descobre o\nformato em produção") {
		t.Errorf("a mensagem deveria explicar por que um contrato mudo é pior que ausente; veio:\n%s", msg)
	}
}

func TestDocRequired_documentoQueMencionaPelaCodigoPassa(t *testing.T) {
	root := t.TempDir()
	escreveDoc(t, root, "docs/contratos/openapi.yaml", "openapi: 3.0.0\n# SRLSS — o ranking\npaths:\n  /servicos: {}\n")
	escreveDoc(t, root, "docs/contratos/dados.md", "# Esquema\n\n## SRLSS\n\ntabela de serviços\n")

	if v, msg := checkDocRequired("", noDeLambda(), root, nil, cfgComDocs()); v != Pass {
		t.Errorf("os dois documentos citam o código; esperava Pass, veio %v: %s", v, msg)
	}
}

// O NOME DO ARQUIVO é a rede secundária, e ela existe por uma razão de formato: um OpenAPI
// cita a rota e o `operationId`, não o código do Anchors. Exigir o código ali faria o gate
// cobrar uma convenção que o formato do documento não tem.
func TestDocRequired_mencaoPeloNomeDoArquivoTambemVale(t *testing.T) {
	root := t.TempDir()
	escreveDoc(t, root, "docs/contratos/openapi.yaml",
		"openapi: 3.0.0\npaths:\n  /servicos:\n    get:\n      operationId: ServiceList\n")
	escreveDoc(t, root, "docs/contratos/dados.md", "# Esquema\n\n## ServiceList\n")

	if v, msg := checkDocRequired("", noDeLambda(), root, nil, cfgComDocs()); v != Pass {
		t.Errorf("`ServiceList` no operationId deveria contar como menção; veio %v: %s", v, msg)
	}
}

// Um documento SÓ, dos dois, não basta — e o teste separa os casos porque o modo de falha
// parcial é o mais provável: quem lembra do OpenAPI esquece do esquema.
func TestDocRequired_umDosDoisNaoBasta(t *testing.T) {
	root := t.TempDir()
	escreveDoc(t, root, "docs/contratos/openapi.yaml", "openapi: 3.0.0\n# SRLSS\n")
	// `dados.md` não existe

	v, msg := checkDocRequired("", noDeLambda(), root, nil, cfgComDocs())
	if v != Fail {
		t.Fatalf("faltando um dos dois, esperava Fail; veio %v", v)
	}
	if !strings.Contains(msg, "dados.md") {
		t.Errorf("a mensagem deveria nomear o que falta; veio:\n%s", msg)
	}
	if strings.Contains(msg, "openapi.yaml") {
		t.Errorf("o documento que ESTÁ certo não deveria aparecer na reprovação; veio:\n%s", msg)
	}
}

// Projeto que não declara `docs:` não é cobrado. Cobrar OpenAPI de quem não tem API seria
// ruído, e ruído ensina a ignorar o gate.
func TestDocRequired_semDeclaracaoNaoCobra(t *testing.T) {
	if v, _ := checkDocRequired("", noDeLambda(), t.TempDir(), nil, &config.Config{}); v != Skip {
		t.Errorf("sem `docs.required` o gate deveria PULAR; veio %v", v)
	}
}

// A camada que não dispara documento nenhum passa sem cobrança — é o que o `trigger`
// existe para permitir. Cobrar toda unidade ensinaria o agente a ignorar o aviso.
func TestDocRequired_camadaSemGatilhoNaoCobra(t *testing.T) {
	// O CAMINHO é o que decide a camada, não o campo do nó — então a unidade fora do
	// gatilho tem de estar fora do diretório também. A primeira versão deste teste só
	// trocava `n.Layer` e continuava apontando para `packages/lambdas/...`: a
	// classificação por caminho vencia, e o teste media o oposto do que dizia medir.
	n := mapx.Node{
		ID: "packages/shared/Formatting.spec.md", Kind: mapx.KindSpec,
		Layer: "spec", Code: "UTLXX",
	}
	if v, _ := checkDocRequired("", n, t.TempDir(), nil, cfgComDocs()); v != Skip {
		t.Errorf("camada fora do `trigger` deveria PULAR; veio %v", v)
	}
}

// O gate parte da SPEC. Partir do código faria a mesma unidade ser cobrada uma vez por
// arquivo — três avisos idênticos para uma trinca, e quem lê aprende a passar por cima.
func TestDocRequired_partiDaSpecNaoDoCodigo(t *testing.T) {
	n := noDeLambda()
	n.Kind = mapx.KindCode
	n.ID = "packages/lambdas/backend/ServiceList.ts"
	if v, _ := checkDocRequired("", n, t.TempDir(), nil, cfgComDocs()); v != Skip {
		t.Errorf("o nó de CÓDIGO deveria pular — a spec representa a unidade; veio %v", v)
	}
}

// A CAMADA DA UNIDADE, e não a do NÓ — e a diferença derruba o gate inteiro.
//
// O nó de uma spec tem `layer: spec` (a camada do ARTEFATO); a unidade que ela descreve
// mora em `lambdas`. Comparar o `trigger: [lambdas]` contra `n.Layer` nunca casa, e o gate
// responde "nenhuma documentação é disparada" para justamente a unidade que deve duas.
//
// Medido contra o projeto de referência: a primeira versão deste gate devolvia `Skip` para
// o `ServiceList` — a unidade cuja falta de OpenAPI motivou o gate existir. Ele teria
// entrado aprovando o caso que veio consertar.
//
// Nenhuma mutação pegaria isso: o código estava consistente consigo mesmo, e os testes
// usavam nós montados à mão com a camada já correta.
func TestDocRequired_usaACamadaDaUnidadeNaoADoNo(t *testing.T) {
	root := t.TempDir()

	// A spec declara a camada no HEADER, que é o que `scan.LayerOfUnit` lê. O nó do mapa
	// diz `spec`, como diz na vida real.
	escreveDoc(t, root, "packages/lambdas/backend/ServiceList.spec.md",
		"<!-- @anchors\n  code: SRLSS\n  layer: lambdas\n-->\n# ServiceList\n")

	n := noDeLambda()
	n.Layer = "spec" // ← o que o mapa realmente guarda

	v, msg := checkDocRequired("", n, root, nil, cfgComDocs())
	if v == Skip {
		t.Fatalf("o gate pulou a unidade que deve duas documentações — leu a camada do NÓ "+
			"em vez da camada da UNIDADE; msg: %s", msg)
	}
	if v != Fail {
		t.Fatalf("os dois documentos não existem; esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, "openapi.yaml") || !strings.Contains(msg, "dados.md") {
		t.Errorf("a mensagem deveria nomear as duas; veio:\n%s", msg)
	}
}
