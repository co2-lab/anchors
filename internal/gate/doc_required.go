package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// checkDocRequired confronta os documentos declarados contra TODAS as unidades que os
// disparam — um veredito por DOCUMENTO, não por unidade.
//
// POR QUE AGREGADO, e a razão é medida:
//
// A primeira versão rodava por nó: cada spec cujo contrato não a mencionava virava um card.
// No projeto de referência isso produziu 37 cards de uma vez — e todos os 37 escreviam nos
// MESMOS dois arquivos (`openapi.yaml`, `dados.md`).
//
// O efeito foi uma fila que não anda: cada PR mergeado invalidava os outros 36, porque o
// ponto de inserção mudava. Medi: de 6 PRs, 2 passavam e 4 conflitavam. Resolver os 21
// restantes exigiria 21 rodadas de resolver-verificar-mergear, uma por PR.
//
// O defeito não é dos agentes nem dos PRs: é do GATE, que fatiou por unidade um trabalho
// que é por documento. Um card "o `openapi.yaml` não descreve estas 18 lambdas" é um PR só,
// com uma seção por unidade — e nenhum conflito.
//
// O QUE NÃO MUDA: as duas perguntas continuam as mesmas (o documento existe, e menciona a
// unidade), e a mensagem continua nomeando cada unidade que falta. O que muda é onde o
// veredito é ancorado.
func checkDocRequiredAggregate(_ config.Gate, root string, graph *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if cfg == nil || len(cfg.AllRequiredDocs()) == 0 {
		return Skip, "o projeto não declara `docs.required` — não há documentação contratada"
	}
	if graph == nil {
		return Skip, "sem mapa: não há como saber quais unidades disparam cada documento"
	}

	// Por DOCUMENTO: quais unidades ele deveria mencionar, e quais ele não menciona.
	type pending struct {
		doc         config.DocArtifact
		missing     bool
		content     string
		unmentioned []string
	}
	byDoc := map[string]*pending{}
	order := []string{}

	for _, n := range graph.Nodes {
		if n.Kind != mapx.KindSpec || n.Code == "" {
			continue
		}
		layer := scan.LayerOfUnit(root, n.ID, cfg)
		for _, d := range cfg.RequiredFor(layer, n.Code) {
			pend, visto := byDoc[d.Path]
			if !visto {
				pend = &pending{doc: d}
				b, err := os.ReadFile(filepath.Join(root, d.Path))
				if err != nil {
					pend.missing = true
				} else {
					pend.content = string(b)
				}
				byDoc[d.Path] = pend
				order = append(order, d.Path)
			}
			if !pend.missing && !mentionsUnit(pend.content, n.Code, n.ID) {
				pend.unmentioned = append(pend.unmentioned, n.Code)
			}
		}
	}

	var msg strings.Builder
	failed := false
	for _, path := range order {
		p := byDoc[path]
		switch {
		case p.missing:
			failed = true
			fmt.Fprintf(&msg, "`%s` (%s) NÃO EXISTE, e é disparado por %d unidade(s).\n\n",
				path, p.doc.Kind, len(p.unmentioned)+1)
		case len(p.unmentioned) > 0:
			failed = true
			sort.Strings(p.unmentioned)
			fmt.Fprintf(&msg, "`%s` (%s) não menciona %d unidade(s): %s.\n\n",
				path, p.doc.Kind, len(p.unmentioned), strings.Join(p.unmentioned, ", "))
		}
	}
	if !failed {
		return Pass, ""
	}

	// A frase que explica o modo de falha, porque ele é silencioso: o arquivo existe, tem
	// conteúdo real, e quem o lê não tem como saber que falta uma entrada.
	msg.WriteString("Um contrato que existe e não descreve a unidade é pior que a ausência\n" +
		"dele: quem consome encontra o documento, confia, e descobre o formato em\n" +
		"produção.\n\n")
	// UM CARD POR DOCUMENTO, e o PR também: a alternativa (um card por unidade) fatiou
	// por unidade um trabalho que é por documento, e produziu 37 PRs que conflitavam
	// entre si porque todos escreviam no mesmo arquivo.
	msg.WriteString("Cada documento é declarado em `docs.required` com o `trigger` que o\n" +
		"aciona. Trate UM documento por vez — as unidades que faltam nele são seções\n" +
		"do mesmo arquivo, e separá-las em PRs diferentes produz conflito a cada merge.")
	return Fail, msg.String()
}

// checkDocRequired confere se as documentações que esta unidade dispara existem e a citam.
func checkDocRequired(_ string, n mapx.Node, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	// Parte da SPEC, e não do código: a spec é a âncora, tem a layer no header e o código
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
	// O nó da spec tem `layer: spec` (a layer do ARTEFATO); a unidade que ela descreve
	// mora em `lambdas`. Um `trigger: [lambdas]` comparado contra `n.Layer` nunca casa, e
	// o gate responde "nenhuma documentação é disparada" para justamente a unidade que
	// deve duas.
	//
	// `scan.LayerOfUnit` resolve o que o autor escreveu no header, e cai no `ClassifyPath`
	// quando não há header. É a mesma função que o `anchors work` usa para informar o
	// dever — e o comentário dela registra este mesmo defeito, medido antes no
	// `docs duties`.
	layer := scan.LayerOfUnit(root, n.ID, cfg)
	duties := cfg.RequiredFor(layer, n.Code)
	if len(duties) == 0 {
		return Skip, fmt.Sprintf("nenhuma documentação é disparada por `%s`", layer)
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
	// A MENÇÃO NUMA NOTA NÃO CONTA — e o caso que obrigou isto é quase cômico.
	//
	// No projeto de referência o `componentes.md` tinha uma nota listando os componentes
	// "que o gate `doc-required` deve sinalizar", e `MetricCard` estava nela. O gate lia o
	// arquivo, achava o nome, e se dava por satisfeito — SILENCIADO PELA PRÓPRIA NOTA QUE
	// DIZIA QUE FALTAVA DOCUMENTÁ-LO. A unidade ficou sem entrada por semanas.
	//
	// A mesma armadilha vale para a citação cruzada: a entrada de `StatusBadge` mencionava
	// `MetricCard` ao explicar quando NÃO se usa um em vez do outro. É prosa legítima, e
	// não é documentação da unidade citada.
	//
	// O QUE CONTINUA GROSSEIRO, deliberadamente. Isto não valida conteúdo, não entende
	// Markdown além de contar `#`, e não julga se o que está escrito presta — a qualidade
	// segue sendo trabalho de revisão. A distinção é só uma: menção DENTRO de um bloco de
	// nota não conta; em qualquer outro lugar, conta como antes.
	// A SEÇÃO PRÓPRIA é a resposta forte: quem documenta abre uma seção para a unidade.
	if hasOwnSection(doc, code, id) {
		return true
	}
	// SEM SEÇÃO, só conta o documento que não TEM seções — um OpenAPI, um YAML, um
	// arquivo de prosa corrida. Ali a menção é tudo o que existe, e cobrar título seria
	// exigir uma estrutura que o formato não tem.
	if hasAnyHeading(doc) {
		return false
	}
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

// hasOwnSection responde se ALGUMA seção do documento é sobre esta unidade.
//
// A DIFERENÇA entre citar e documentar, e ela decidiu o desenho. Duas armadilhas reais, no
// mesmo arquivo do projeto de referência:
//
//	· uma NOTA listava os componentes "que o gate `doc-required` deve sinalizar", e
//	  `MetricCard` estava nela — o gate lia, achava o nome, e se dava por satisfeito.
//	  SILENCIADO PELA PRÓPRIA NOTA QUE DIZIA QUE FALTAVA DOCUMENTÁ-LO.
//
//	· a seção de `StatusBadge` citava `MetricCard` ao explicar quando NÃO se usa um em vez
//	  do outro. Prosa legítima e útil — e não é documentação da unidade citada.
//
// Nos dois casos a unidade aparecia no arquivo e não tinha entrada. Ficou sem documentação
// por semanas, com o gate verde.
//
// O TÍTULO é o que separa: quem documenta uma unidade abre uma seção para ela. Citar é
// escrever o nome no corpo de outra.
//
// AGNÓSTICO. Não procura a palavra "Nota", nem qualquer outra — o Anchors roda em projetos
// que documentam em qualquer idioma. Procura a IDENTIDADE (o código, o nome do arquivo) na
// linha de título, que é estrutura de Markdown, não vocabulário.
func hasOwnSection(doc, code, id string) bool {
	base := baseName(id)
	for _, l := range strings.Split(doc, "\n") {
		if headingLevel(l) == 0 {
			continue
		}
		// FRONTEIRA, e não substring: `MetricCardList.spec.md` no título contaria como
		// documentação de `MetricCard` — um nome que CONTÉM o outro é outra unidade, e
		// aceitar isso devolveria o defeito por uma porta lateral.
		if code != "" && wholeWordInLine(l, code) {
			return true
		}
		if base != "" && wholeWordInLine(l, base) {
			return true
		}
	}
	return false
}

// baseName é o nome do arquivo sem extensão: `ServiceList.spec.md` vira `ServiceList`.
func baseName(id string) string {
	base := filepath.Base(id)
	for _, sufixo := range []string{".spec.md", ".md", ".ts", ".tsx", ".go"} {
		base = strings.TrimSuffix(base, sufixo)
	}
	return base
}

// headingLevel conta os `#` iniciais de uma linha de título Markdown (0 se não for uma).
func headingLevel(l string) int {
	t := strings.TrimLeft(l, " \t")
	n := 0
	for n < len(t) && t[n] == '#' {
		n++
	}
	if n == 0 || n >= len(t) || t[n] != ' ' {
		return 0
	}
	return n
}

// hasAnyHeading diz se o documento é organizado em seções.
//
// Um `openapi.yaml` não tem títulos Markdown, e exigir uma seção por unidade ali seria
// cobrar uma estrutura que o formato não tem — foi por isso que a menção existe como rede
// secundária desde o início. A distinção mantém esse caso funcionando como antes.
func hasAnyHeading(doc string) bool {
	for _, l := range strings.Split(doc, "\n") {
		if headingLevel(l) > 0 {
			return true
		}
	}
	return false
}

// wholeWordInLine procura o nome com FRONTEIRA nos dois lados.
//
// Um identificador é delimitado por qualquer caractere que não componha nome: crase,
// espaço, barra, parêntese, ponto. `MetricCard` casa em "`MetricCard.spec.md`" e não casa
// em "`MetricCardList.spec.md`" — e é essa distinção que impede o título de uma unidade de
// contar como documentação de outra cujo nome ele contém.
func wholeWordInLine(line, name string) bool {
	i := 0
	for {
		j := strings.Index(line[i:], name)
		if j < 0 {
			return false
		}
		start := i + j
		end := start + len(name)
		if !isNameChar(line, start-1) && !isNameChar(line, end) {
			return true
		}
		i = start + 1
	}
}

// isNameChar diz se a posição carrega um caractere que faz parte de um identificador.
// Fora dos limites da linha conta como fronteira.
func isNameChar(s string, i int) bool {
	if i < 0 || i >= len(s) {
		return false
	}
	c := s[i]
	return c == '_' || c == '-' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
