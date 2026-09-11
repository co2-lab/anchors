package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// --- a DOCUMENTAÇÃO AGREGADA que a unidade deve alimentar ---
//
// O `anchors.yaml` declara quais documentos o projeto tem como contrato, e QUANDO cada um
// precisa ser tocado:
//
//	docs:
//	    required:
//	        - kind: openapi
//	          path: docs/contratos/openapi.yaml
//	          trigger: [lambdas]
//
// `config.RequiredFor` resolve isso desde que o campo existe. Dois lugares o consultavam —
// o `anchors work` (que informa o dever a quem pega o card) e o `docs duties` (que lista).
// NENHUM confrontava.
//
// O QUE ISSO CUSTOU, medido no projeto de referência: um PR entregou uma lambda nova — com
// spec, código, feature e teste — sem tocar o OpenAPI nem o esquema de dados. Os dois são
// declarados com `trigger: [lambdas]`, os dois são obrigatórios, e TODOS os checks ficaram
// verdes. Só apareceu porque dois agentes colidiram no mesmo card e havia com que comparar:
// a branch do segundo tinha os 267 linhas de OpenAPI que a do primeiro não tinha.
//
// Sem a colisão, ninguém teria notado. Não há com que comparar quando só um trabalho existe.
//
// AS DUAS PERGUNTAS, e a segunda é a que importa:
//
//	1. o documento EXISTE?
//	2. ele MENCIONA esta unidade?
//
// A primeira sozinha aprovaria um arquivo vazio criado para calar o gate. A segunda é
// grosseira de propósito — ela procura o código de identidade e o nome do arquivo, não
// entende o conteúdo. Um OpenAPI que cita a rota mas descreve o formato errado passa aqui,
// e isso é aceitável: o gate separa "não documentado" de "documentado", e a qualidade do
// que está escrito é trabalho de revisão.
//
// O que ele NÃO faz é aprovar silêncio.

// checkDocRequired confere se as documentações que esta unidade dispara existem e a citam.
func checkDocRequired(_ string, n mapx.Node, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	// Parte da SPEC, e não do código: a spec é a âncora, tem a camada no header e o código
	// de identidade. Partir do código faria a mesma unidade ser cobrada uma vez por
	// arquivo — três avisos idênticos para uma trinca.
	if n.Kind != mapx.KindSpec {
		return Skip, "o dever de documentação é da UNIDADE, e a spec é quem a representa"
	}
	if cfg == nil || len(cfg.AllRequiredDocs()) == 0 {
		return Skip, "o projeto não declara `docs.required` — não há documentação contratada"
	}

	// A CAMADA DA UNIDADE, e não a do nó — e a diferença derruba o gate inteiro.
	//
	// O nó da spec tem `layer: spec` (a camada do ARTEFATO); a unidade que ela descreve
	// mora em `lambdas`. Um `trigger: [lambdas]` comparado contra `n.Layer` nunca casa, e
	// o gate responde "nenhuma documentação é disparada" para justamente a unidade que
	// deve duas.
	//
	// `scan.LayerOfUnit` resolve o que o autor escreveu no header, e cai no `ClassifyPath`
	// quando não há header. É a mesma função que o `anchors work` usa para informar o
	// dever — e o comentário dela registra este mesmo defeito, medido antes no
	// `docs duties`.
	camada := scan.LayerOfUnit(root, n.ID, cfg)
	duties := cfg.RequiredFor(camada, n.Code)
	if len(duties) == 0 {
		return Skip, fmt.Sprintf("nenhuma documentação é disparada por `%s`", camada)
	}

	var missing, silent []string
	for _, d := range duties {
		b, err := os.ReadFile(filepath.Join(root, d.Path))
		if err != nil {
			missing = append(missing, fmt.Sprintf("%s (%s)", d.Path, d.Kind))
			continue
		}
		if n.Code != "" && !mentionsUnit(string(b), n.Code, n.ID) {
			silent = append(silent, fmt.Sprintf("%s (%s)", d.Path, d.Kind))
		}
	}

	if len(missing) == 0 && len(silent) == 0 {
		return Pass, ""
	}

	var msg strings.Builder
	if len(missing) > 0 {
		msg.WriteString(fmt.Sprintf("%d documentação(ões) que esta unidade dispara não existe(m): %s.\n\n",
			len(missing), strings.Join(missing, ", ")))
	}
	if len(silent) > 0 {
		msg.WriteString(fmt.Sprintf("%d documentação(ões) existe(m) e não menciona(m) `%s`: %s.\n\n",
			len(silent), n.Code, strings.Join(silent, ", ")))
		// A frase que explica o modo de falha, porque ele é silencioso: o arquivo existe,
		// tem conteúdo real, e quem o lê não tem como saber que falta uma entrada.
		msg.WriteString("Um contrato que existe e não descreve a unidade nova é pior que a\n" +
			"ausência dele: quem consome encontra o documento, confia, e descobre o\n" +
			"formato em produção.\n\n")
	}
	msg.WriteString("Cada uma é declarada em `docs.required` com o `trigger` que a aciona.\n")
	msg.WriteString("`anchors docs duties --unit <arquivo>` lista o que esta unidade deve.")
	return Fail, msg.String()
}

// mentionsUnit procura a unidade no documento, e procura de dois jeitos.
//
// O CÓDIGO de identidade é o alvo principal: ele é único no projeto e é o que a doutrina
// usa para ligar as coisas. O NOME DO ARQUIVO é a rede secundária — um OpenAPI costuma
// citar a rota e o operationId, não o código do Anchors, e exigir o código ali faria o gate
// cobrar uma convenção que o formato do documento não tem.
//
// Deliberadamente grosseiro: não entende YAML, não entende Markdown, não valida conteúdo.
// Ele separa "não documentado" de "documentado" — e a qualidade do que está escrito é
// trabalho de revisão, não de gate.
func mentionsUnit(doc, code, id string) bool {
	if code != "" && strings.Contains(doc, code) {
		return true
	}
	// O nome do arquivo SEM extensão: `ServiceList.spec.md` vira `ServiceList`, que é o que
	// aparece num `operationId` ou num título de seção.
	base := filepath.Base(id)
	for _, sufixo := range []string{".spec.md", ".md", ".ts", ".tsx", ".go"} {
		base = strings.TrimSuffix(base, sufixo)
	}
	return base != "" && strings.Contains(doc, base)
}
