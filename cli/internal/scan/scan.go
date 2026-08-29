// Package scan percorre o repositório lendo TEXTO (nunca parseando código) e
// extrai os fatos de que o mapa precisa: que arquivos existem, de que CAMADA (pela
// Estrutura declarada, não por regra fixa), e que código de cenário cada um
// carrega. Ver DECISIONS.md D2.
package scan

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
)

// scenarioCodeRE reconhece um código de cenário (TRACEABILITY §3): `<UNIDADE>-<LETRA><NN>`.
//
// A classe de LETRAS é o vocabulário de tipos de regra do projeto (`rule_types`), não uma
// constante. As canônicas são só o DEFAULT de quem não declara nada — um projeto que
// declare a sua (uma letra para "Invariante", por exemplo) precisa ser enxergado aqui.
//
// Este regex já esteve soldado com as letras canônicas enquanto `gate` e `testsig` se
// reconfiguravam pela config. A consequência não seria um erro: seria SILÊNCIO. Os gates
// enxergariam o código novo e o scan não, então a unidade simplesmente não entraria no
// mapa — e nenhum gate acusaria a ausência, porque para eles ela nunca existiu. É o mesmo
// modo de falha já medido duas vezes neste engine (um dialeto de assinatura que só
// reconhecia `export function` cegava 97,6% do código; um regex só em inglês escondia 108
// cenários de um gate BLOQUEANTE): **vocabulário errado não dá erro, dá silêncio.**
var scenarioCodeRE = scenarioCodeREFor(config.DefaultRuleLetters)

func scenarioCodeREFor(letters string) *regexp.Regexp {
	if letters == "" {
		letters = config.DefaultRuleLetters
	}
	return regexp.MustCompile(`\b[A-Z0-9]` + config.CodeLengthPattern() + `-(?:[` + regexp.QuoteMeta(letters) +
		`]\d{2}(?:-[a-z][a-z0-9-]*)?|DS-[A-Za-z0-9-]+|VR)\b`)
}

// SetRuleLetters reconfigura a gramática de código deste pacote para o vocabulário do
// projeto. Espelha as funções de mesmo nome em `gate` e `testsig` — os três leem os mesmos
// códigos e precisam concordar sobre o que é um código.
func SetRuleLetters(letters string) {
	scenarioCodeRE = scenarioCodeREFor(letters)
}

// ScenarioCodeRE expõe a gramática VIGENTE de código de cenário — já com as
// `rule_types` do projeto aplicadas por SetRuleLetters.
//
// Existe para que outros pacotes reconheçam código sem reescrever o regex: uma cópia com
// as letras canônicas fixas rejeitaria o código de uma letra que o projeto declarou, e o
// Anchors passaria a exigir uma referência que ele mesmo se recusa a reconhecer.
func ScenarioCodeRE() *regexp.Regexp { return scenarioCodeRE }

// noPropRE detecta a anotação @noPropagation (opt-out honesto no filho).
var noPropRE = regexp.MustCompile(`@noPropagation\b`)

// sharedCodeRE detecta @anchors-shared-code — o opt-out honesto da colisão de
// identidade: "os códigos de cenário aqui pertencem a OUTRA unidade de propósito"
// (ex.: um teste de util/handler que prova o cenário da tela que ele serve). Deixa o
// registro no código (datado, localizado) em vez de silenciar a regra escondido —
// como `// eslint-disable`. O detector de identidade-duplicada ignora quem a declara.
var sharedCodeRE = regexp.MustCompile(`@anchors-shared-code\b`)

// File é um arquivo escaneado, já classificado numa camada da Estrutura.
type File struct {
	Path  string // relativo à raiz
	Layer string // nome da camada que casou (ex.: "spec", "code")
	Kind  string // kind do nó (vem da camada)
	Rev   string // sha256 curto do conteúdo
	Codes []string
	// HeaderCode é a identidade DECLARADA no header (`code: XXXX`) — a fonte de verdade
	// quando existe. Sem ela, a identidade é inferida do primeiro código de cenário do
	// texto, o que erra quando a spec CITA outra unidade antes de definir a sua: uma spec
	// de modelo que abre referenciando `DTAXX-B11` era registrada no mapa como dona de
	// `DTAX`, e todo gate relacional passava a olhar a unidade errada.
	HeaderCode    string
	NoPropagation bool  // o texto contém a anotação @noPropagation
	SharedCode    bool  // o texto contém @anchors-shared-code (opt-out de colisão)
	Deps          []Dep // dependências de reúso declaradas (Tabela de Dependências, SPEC_TYPES §5)
	// Seeds: os caminhos de spec que um PLANO semeia. Extraídos aqui porque o scan já tem
	// o conteúdo em mãos — o mapa não precisa reabrir o arquivo só para isso.
	Seeds []string
	// Needs: os planos que precisam TERMINAR antes deste começar (`needs:` no header).
	//
	// É diferente de `Deps`: aquele diz "consumo este arquivo", este diz "não comece
	// antes daquele trabalho acabar". A distinção existe porque com vários agentes
	// puxando cards de uma fila comum, alguém pegaria o plano de uma feature antes de a
	// fundação existir — e descobriria isso ao tentar criar o primeiro arquivo.
	Needs []string
	// Parent é o CÓDIGO do artefato que contém este — a fase de um card, o plano de uma
	// fase. Um só: pertencimento é uma relação de um pai, ao contrário de `needs`, que é
	// lista porque uma coisa pode depender de várias.
	Parent string
}

// Dep é uma linha da Tabela de Dependências de uma spec consumidora (SPEC_TYPES §5):
// a declaração de que esta unidade consome outra camada. Vira uma aresta `depends-on`
// da spec para o Arquivo, com Method como metadado (impacto fino).
type Dep struct {
	Code   string // DEPn — local à spec (ponteiro que o data contract usa)
	File   string // arquivo da camada consumida (o alvo da aresta, relativo à raiz)
	Method string // método/símbolo consumido (metadado da aresta)
	Layer  string // tipo/camada declarado da dependência (p/ gate de limite de camada)
}

var ignored = map[string]bool{
	"node_modules": true, ".git": true, "dist": true, "build": true,
	"vendor": true, ".next": true, "coverage": true, ".expo": true,
}

// Walk percorre root e classifica cada arquivo pelas CAMADAS da config. Um arquivo
// que não casa nenhuma camada é ignorado. Quando casa mais de uma, vence a de
// pattern mais específico (heurística: maior comprimento do pattern).
func Walk(root string, cfg *config.Config) ([]File, error) {
	// O vocabulário é ligado AQUI, na porta de entrada, e não em cada chamador: um
	// chamador novo que esquecesse a ligação não daria erro — daria um mapa sem as
	// unidades cujo código usa letra declarada pelo projeto, silenciosamente.
	SetRuleLetters(cfg.RuleLetters())
	var out []File
	// O que o projeto declarou descartável no `.gitignore` também não é trabalho para o
	// Anchors. Antes, só o WATCHER respeitava essa lista, e o `Walk` — que alimenta
	// `check`, `map` e `recode` — usava a lista embutida: a sonda de um revisor não era
	// enfileirada pelo watcher, mas o `check` a contava como unidade a fazer. Duas fontes
	// da régua discordando sobre o mesmo assunto, que é o modo de falha mais caro medido
	// neste engine — regra ausente produz pergunta; regra contraditória produz
	// desobediência sem consciência.
	ig := LoadIgnoreFor(root, cfg)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// `filepath.Rel` devolve o separador do SO — no Windows, `\`. O resto do Anchors
		// fala em `/`: os patterns do `layers:`, os IDs dos nós do mapa, as arestas, os
		// caminhos que o usuário digita. Normalizar aqui, na fronteira, é o que impede a
		// forma nativa de vazar para dentro do grafo.
		rel := filepath.ToSlash(mustRel(root, path))
		if d.IsDir() {
			if ig.SkipDir(d.Name(), rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if ig.SkipFile(rel) {
			return nil
		}
		layer, kind := classify(rel, cfg)
		if layer == "" {
			return nil // não pertence a nenhuma camada declarada
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		out = append(out, File{
			Path:          rel,
			Layer:         layer,
			Kind:          kind,
			Rev:           shortHash(content),
			Codes:         extractCodes(content),
			HeaderCode:    extractHeaderCode(string(content)),
			Seeds:         extractSeeds(kind, string(content)),
			Needs:         needsFor(kind, content, root, rel),
			Parent:        parentDe(content),
			NoPropagation: noPropRE.Match(content),
			SharedCode:    sharedCodeRE.Match(content),
			Deps:          depsFor(kind, content, root, rel),
		})
		return nil
	})
	return out, err
}

// Classify decide a que camada (e kind) um caminho relativo pertence, pelos
// patterns da Estrutura. Devolve layer="" se não pertence a nenhuma. Exportada
// para o watcher classificar um único arquivo alterado sem re-escanear tudo.
func Classify(rel string, cfg *config.Config) (layer, kind string) {
	return classify(rel, cfg)
}

// classify decide a que camada um arquivo pertence, pelos patterns da Estrutura.
// Respeita `exclude`.
//
// Quando DOIS patterns casam o mesmo arquivo, alguém tem de desempatar — e o `doublestar`
// não ajuda: uma biblioteca de glob responde sim/não, porque glob não é uma linguagem com
// precedência (diferente de CSS ou de roteamento HTTP, onde a especificidade faz parte da
// especificação). O desempate é, necessariamente, decisão do Anchors.
//
// A régua é COMPRIMENTO DO PATTERN, e ela é um chute — mede verbosidade, não precisão. O
// caso que expôs isso: `packages/backend/{infra/**,functions/**/index.ts,…}` (108 chars)
// vencia `packages/backend/infra/models/**/*.ts` (37 chars), embora o segundo aponte para
// um subconjunto do primeiro. A alternativa ganhou por ser prolixa.
//
// Foram tentadas três réguas melhores (prefixo literal; prefixo+sufixo; expansão de
// alternativas medindo o pior caso). Todas acertaram o caso que as motivou e quebraram
// outro: a terceira reclassificou 72 nós do repositório real, porque um `**` no MEIO do
// pattern (`apps/mobile/src/**/business-logic/**/*.ts`) esconde justamente a parte
// discriminante. A lição medida: **toda heurística de especificidade acerta o exemplo à
// vista e erra o que não está sendo olhado** — ela adivinha intenção a partir de sintaxe.
//
// Então o comprimento fica como DEFAULT (não regride nada no repo real), e a intenção,
// quando existe, é DECLARADA: `priority: N` na camada vence qualquer heurística. É o mesmo
// princípio de `identified_as_form` e `dialect` — onde o framework não pode saber, ele
// pergunta em vez de adivinhar, e a resposta fica à vista na Estrutura.
//
// A iteração sobre `cfg.Layers` é um map, cuja ordem Go não define. Empate TOTAL (mesma
// prioridade e mesmo comprimento) sortearia a camada a cada execução — e uma classificação
// que muda entre duas rodadas do mesmo comando envenena todo gate que depende dela. Por
// isso o desempate final é o NOME da camada: arbitrário, mas estável.
func classify(rel string, cfg *config.Config) (layer, kind string) {
	// Normaliza AQUI, e não só em quem chama: são mais de dez chamadores (`check`, `code`,
	// `new`, `watch`…), e basta um passar a forma nativa do Windows para o `doublestar.Match`
	// abaixo não casar NENHUM pattern que tenha diretório — `apps/mobile/**` contra
	// `apps\mobile\App.tsx`. O modo de falha é silencioso e caro: o arquivo não fica com a
	// camada errada, ele fica SEM camada, some do mapa, e nenhum gate volta a confrontá-lo.
	//
	// MEDIDO em 24/08 no app de referência: `anchors map build` no Windows produzia 1857 nós e 5416
	// arestas onde o mesmo commit dava 2799/11302 — sumiam TODOS os `depends-on`,
	// `specifies`, `covered-by` e `seeds`, e metade dos `governs`. Sem uma linha de erro: o
	// mapa encolhia 4× e o `check` seguinte acusava bloqueante em arquivo intocado.
	//
	// `strings.ReplaceAll`, NÃO `filepath.ToSlash`: este último troca pelo separador da
	// plataforma ONDE RODA, então no macOS/Linux é no-op e a `\` passa intacta. O teste
	// que cobre isto foi escrito no Windows (onde `ToSlash` resolvia) e reprovava aqui —
	// o mesmo caminho tem de casar nas duas formas em QUALQUER máquina, porque o mapa é
	// versionado e trafega entre elas.
	rel = strings.ReplaceAll(rel, `\`, "/")
	melhor := prioridade{prio: math.MinInt, tam: -1}
	for name, l := range cfg.Layers {
		if !matchGlob(l.Pattern, rel) {
			continue
		}
		if excluded(rel, l.Exclude) {
			continue
		}
		p := prioridade{prio: l.Priority, tam: len(l.Pattern), nome: name}
		if layer == "" || p.venceContra(melhor) {
			melhor, layer, kind = p, name, l.Kind
		}
	}
	return layer, kind
}

// prioridade é a régua de desempate, em ordem: declarada > comprimento > nome.
type prioridade struct {
	prio int
	tam  int
	nome string
}

func (a prioridade) venceContra(b prioridade) bool {
	if a.prio != b.prio {
		return a.prio > b.prio
	}
	if a.tam != b.tam {
		return a.tam > b.tam
	}
	return a.nome < b.nome // estabilidade: o map não tem ordem
}

// AmbiguidadeDeCamada é um arquivo cuja camada foi decidida por HEURÍSTICA — dois patterns
// casaram, nenhum declarou `priority`, e o desempate foi o comprimento do pattern.
//
// Não é erro: na maioria das vezes o comprimento acerta. É um AVISO, porque é o ponto onde
// o Anchors adivinhou a intenção do projeto, e adivinhação silenciosa é o que produz a
// classificação errada que ninguém vê. Quem quiser resolver, declara `priority`.
type AmbiguidadeDeCamada struct {
	Arquivo    string
	Vencedora  string
	Perdedoras []string
}

// Ambiguidades devolve os arquivos classificados por desempate heurístico — o material do
// alerta em `check`/`doctor`.
func Ambiguidades(files []File, cfg *config.Config) []AmbiguidadeDeCamada {
	var out []AmbiguidadeDeCamada
	for _, f := range files {
		var casam []prioridade
		for name, l := range cfg.Layers {
			if !matchGlob(l.Pattern, f.Path) || excluded(f.Path, l.Exclude) {
				continue
			}
			casam = append(casam, prioridade{prio: l.Priority, tam: len(l.Pattern), nome: name})
		}
		if len(casam) < 2 {
			continue
		}
		sort.Slice(casam, func(i, j int) bool { return casam[i].venceContra(casam[j]) })
		// Se a vencedora se distingue por prioridade DECLARADA, o projeto já decidiu —
		// não há adivinhação a denunciar.
		if casam[0].prio > casam[1].prio {
			continue
		}
		a := AmbiguidadeDeCamada{Arquivo: f.Path, Vencedora: casam[0].nome}
		for _, p := range casam[1:] {
			a.Perdedoras = append(a.Perdedoras, p.nome)
		}
		out = append(out, a)
	}
	return out
}

func excluded(rel string, globs []string) bool {
	for _, g := range globs {
		if matchGlob(g, rel) {
			return true
		}
	}
	return false
}

// matchGlob casa um glob (com ** e {a,b}) contra um caminho relativo.
// mustRel devolve o caminho relativo à raiz. O erro só ocorre quando `path` não está sob
// `root` — impossível vindo de `WalkDir(root, …)`. Se ocorresse, engolir o erro devolveria
// o caminho ABSOLUTO como se fosse relativo, e todo pattern deixaria de casar: o arquivo
// sumiria do mapa em silêncio, que é exatamente o modo de falha que o `ToSlash` acima
// existe para fechar. Devolver vazio é ruidoso de propósito.
func mustRel(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	return rel
}

func matchGlob(pattern, path string) bool {
	ok, err := doublestar.Match(pattern, path)
	return err == nil && ok
}

// extractCodes devolve os códigos que o arquivo POSSUI — ignorando os que aparecem
// apenas em COMENTÁRIOS. Um código citado num comentário (`// ...vê CODEX-VR`) é uma
// REFERÊNCIA a outra unidade, não posse; contá-lo como posse gera colisão de
// identidade falsa (TRACEABILITY §3). Specs (markdown) não têm comentário de código,
// então seus códigos de header/seção contam normalmente.
func extractCodes(content []byte) []string {
	stripped := stripComments(string(content))
	matches := scenarioCodeRE.FindAllString(stripped, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := map[string]bool{}
	var codes []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			codes = append(codes, m)
		}
	}
	return codes
}

// stripComments remove o conteúdo de comentários de linha (`//`, `#`) e de bloco
// (`/* */`) para que códigos citados neles não contem como posse. Conservador: não
// tenta entender strings/regex (raro conter um código de cenário); o custo de um
// falso-negativo aqui é ínfimo perto do falso-positivo que resolve. Markdown (`#` de
// heading) NÃO é tratado como comentário — só `#` seguido de espaço em início de
// linha de código; para não quebrar specs, `#` só conta como comentário quando NÃO
// é um cabeçalho markdown (heurística: linha que começa com `#` + espaço fica).
func stripComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inBlock := false
	for _, line := range strings.Split(s, "\n") {
		if inBlock {
			if i := strings.Index(line, "*/"); i >= 0 {
				inBlock = false
				line = line[i+2:]
			} else {
				b.WriteByte('\n')
				continue
			}
		}
		// comentário de bloco iniciado nesta linha
		for {
			open := strings.Index(line, "/*")
			if open < 0 {
				break
			}
			close := strings.Index(line[open:], "*/")
			if close < 0 {
				line = line[:open]
				inBlock = true
				break
			}
			line = line[:open] + line[open+close+2:]
		}
		// comentário de linha `//` (JS/TS/Go) — corta a partir dele
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// depHeadingRE casa o cabeçalho da seção de dependências (markdown), tolerante a
// nível de heading e acentuação/caixa: "## Dependências", "### Dependencias de Dados".
var depHeadingRE = regexp.MustCompile(`(?im)^#{1,6}\s+depend[eê]ncias?\b`)

// depCodeRE valida o código de uma linha de dependência: DEP seguido de dígitos.
var depCodeRE = regexp.MustCompile(`^DEP\d+$`)

// headerDepRE captura a linha `dep:` do cabeçalho @anchors (camadas reconhecidas sem
// spec declaram aqui os ARQUIVOS de que dependem — a "Tabela de Dependências inline",
// já que não têm .spec.md). Aceita 1+ caminhos separados por vírgula.
var headerDepRE = regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*dep:\s*(.+)$`)

// headerNeedsRE captura a linha `needs:` do cabeçalho @anchors de um PLANO: os planos
// que precisam terminar antes deste começar. Aceita 1+ caminhos separados por vírgula.
var headerNeedsRE = regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*needs:\s*(.+)$`)

// headerParentRE captura a linha `parent:` — QUEM É O PAI deste artefato.
//
// `parent` e `needs` são relações DIFERENTES, e confundi-las produz uma árvore errada:
// `needs` é ORDEM ("vem depois de"), `parent` é PERTENCIMENTO ("está dentro de"). A fase 2
// vem depois da fase 1 e NÃO está dentro dela — as duas pertencem ao mesmo plano.
//
// Enquanto só havia `needs`, a árvore inferia o pai a partir da ordem, e desenhava as
// quatro fases de um plano encaixadas uma na outra como uma escada.
var headerParentRE = regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*parent:\s*(.+)$`)

// needsFor lê `needs:` de um PLANO (caminhos de outros planos) ou de uma SPEC (o CÓDIGO
// da fase que precisa fechar antes).
//
// Na spec o `needs:` era recusado, com o argumento de que "uma spec não espera trabalho
// terminar, ela É o trabalho". Isso vale para a spec como DOCUMENTO — mas não para o
// trabalho de implementá-la, que é o que vira card. A ordem dentro de um plano vivia em
// prosa ("Fase 3 — depende da Fase 2"), e prosa não é confrontável: as specs de um plano
// nasciam todas disponíveis, e o agente pegava a da fase 3 com a fase 1 em aberto.
//
// Na spec o valor é um CÓDIGO (`FNDTN-F02`), não um caminho — a fase não é um arquivo, é
// um item catalogado dentro do plano. Por isso a resolução de caminho não se aplica aqui.
func needsFor(kind string, content []byte, root, rel string) []string {
	switch kind {
	case "plan":
		return extractNeeds(content, root, rel)
	case "spec":
		return extractNeedsCodigo(content)
	}
	return nil
}

// parentDe lê o `parent:` do header — o código de quem contém este artefato.
//
// Vale para qualquer kind: uma spec pertence a uma fase, uma fase a um plano, e um plano
// pode pertencer a outro num projeto que agrupe por épico. A hierarquia é livre porque a
// forma de organizar trabalho varia, e o Anchors não tem por que impor uma.
func parentDe(content []byte) string {
	m := headerParentRE.FindSubmatch(content)
	if m == nil {
		return ""
	}
	raw := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(string(m[1])), "-->"))
	return stripInlineCode(raw)
}

// extractNeedsCodigo lê o `needs:` de uma spec, onde o valor é o código da fase.
func extractNeedsCodigo(content []byte) []string {
	m := headerNeedsRE.FindSubmatch(content)
	if m == nil {
		return nil
	}
	raw := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(string(m[1])), "-->"))
	var out []string
	for _, part := range strings.Split(raw, ",") {
		p := stripInlineCode(strings.TrimSpace(part))
		// Só o que PARECE código de fase. Um caminho aqui é engano de quem escreveu (a
		// spec depende de uma FASE, não de um arquivo), e aceitá-lo em silêncio deixaria
		// a dependência sem efeito — o gate não a encontraria no plano.
		if p == "" || !codigoDeFaseRE.MatchString(p) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// codigoDeFaseRE casa `FNDTN-F02` — o código de uma fase de plano.
var codigoDeFaseRE = regexp.MustCompile(`^[A-Z0-9]` + config.CodeLengthPattern() + `-F\d{2}$`)

// extractNeeds lê a linha `needs:` e resolve cada caminho relativo à raiz. Só faz
// sentido em plano — um `needs:` numa spec seria a pergunta errada: spec não espera
// trabalho terminar, ela É o trabalho.
func extractNeeds(content []byte, root, rel string) []string {
	m := headerNeedsRE.FindSubmatch(content)
	if m == nil {
		return nil
	}
	raw := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(string(m[1])), "-->"))
	var out []string
	for _, part := range strings.Split(raw, ",") {
		p := stripInlineCode(strings.TrimSpace(part))
		if p == "" {
			continue
		}
		out = append(out, resolveDepPath(root, rel, p))
	}
	return out
}

// depsFor extrai as dependências de reúso conforme o papel do arquivo:
//   - SPEC → a Tabela de Dependências markdown (SPEC_TYPES §5).
//   - demais (código de camada reconhecida sem spec) → a linha `dep:` do header @anchors.
//
// Um guide/doc pode ter uma seção "Dependências" documental (tabela Arquivo|Descrição)
// que NÃO é dep de reúso; por isso a tabela é lida só p/ spec (o header `dep:` não colide).
func depsFor(kind string, content []byte, root, rel string) []Dep {
	if kind == "spec" {
		return extractDeps(content, root, rel)
	}
	// Roteiro de teste (.yaml): a dependência é a COMPOSIÇÃO — o `runFlow:` que ele
	// executa. Um roteiro não tem header `@anchors` com linha `dep:`, então sem este
	// ramo ele ficaria sem nenhuma aresta de saída, e a composição (um util usado por
	// centenas de flows) seria invisível ao impacto e ao frescor de evidência.
	if kind == "test" && (strings.HasSuffix(rel, ".yaml") || strings.HasSuffix(rel, ".yml")) {
		var deps []Dep
		for _, alvo := range ComposeRefs(string(content), rel) {
			deps = append(deps, Dep{File: alvo, Method: "runFlow"})
		}
		return deps
	}
	return extractHeaderDeps(content, root, rel)
}

// extractHeaderDeps lê a linha `dep:` do cabeçalho e resolve cada caminho para um alvo
// relativo à raiz. Sem código local (DEPn) — a identidade da dep é o próprio arquivo.
func extractHeaderDeps(content []byte, root, rel string) []Dep {
	m := headerDepRE.FindSubmatch(content)
	if m == nil {
		return nil
	}
	raw := strings.TrimSpace(string(m[1]))
	// corta comentário-fecho de markdown e um "porquê" após " -- " ou " # "
	raw = strings.TrimSuffix(strings.TrimSpace(strings.TrimSuffix(raw, "-->")), "")
	var deps []Dep
	for _, part := range strings.Split(raw, ",") {
		p := stripInlineCode(strings.TrimSpace(part))
		if p == "" {
			continue
		}
		deps = append(deps, Dep{File: resolveDepPath(root, rel, p)})
	}
	return deps
}

// extractDeps lê a Tabela de Dependências de uma spec (SPEC_TYPES §5) e devolve suas
// linhas. Só specs (markdown) têm essa tabela; para os demais devolve nil. O parsing
// é de TEXTO (coerente com D2 — nunca parseia código): acha a 1ª seção cujo heading
// casa `Dependências`, lê a 1ª tabela markdown ali, mapeia as colunas pelo header
// (Cód/Arquivo/Método/Camada, em qualquer ordem) e extrai as linhas válidas.
//
// O Arquivo declarado na tabela é resolvido para caminho relativo à raiz: aceita-se
// tanto um caminho já relativo à raiz quanto um relativo ao src comum (o autor
// escreve `stores/auth.store.ts`); a resolução tenta casar contra um arquivo existente.
func extractDeps(content []byte, root, specRel string) []Dep {
	s := string(content)
	loc := depHeadingRE.FindStringIndex(s)
	if loc == nil {
		return nil
	}
	// varre da seção em diante, pega a primeira tabela (linhas iniciando com '|')
	lines := strings.Split(s[loc[1]:], "\n")
	var header []string
	var rows [][]string
	inTable := false
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if !inTable {
			// nova seção antes da tabela → não há tabela nesta seção
			if strings.HasPrefix(t, "#") {
				break
			}
			if strings.HasPrefix(t, "|") {
				header = splitRow(t)
				inTable = true
			}
			continue
		}
		if !strings.HasPrefix(t, "|") {
			break // fim da tabela
		}
		if isDividerRow(t) {
			continue // a linha |---|---| do markdown
		}
		rows = append(rows, splitRow(t))
	}
	if header == nil {
		return nil
	}
	ci := columnIndex(header) // mapeia papel→índice pela semântica do header
	// A Tabela de Dependências (SPEC_TYPES §5) EXIGE as colunas Cód (DEPn) E Arquivo.
	// A coluna Cód é o que a distingue de uma tabela documental de "Dependências" (ex.:
	// um guide que lista arquivos relacionados numa tabela Arquivo|Descrição) — sem
	// código local não há a referência que o data contract usa, e não é dep de reúso.
	if ci.code < 0 || ci.file < 0 {
		return nil
	}
	var deps []Dep
	for _, r := range rows {
		get := func(i int) string {
			if i < 0 || i >= len(r) {
				return ""
			}
			return stripInlineCode(strings.TrimSpace(r[i]))
		}
		getRaw := func(i int) string {
			if i < 0 || i >= len(r) {
				return ""
			}
			return strings.TrimSpace(r[i])
		}
		code := get(ci.code)
		file := get(ci.file)
		if file == "" || (ci.code >= 0 && !depCodeRE.MatchString(code)) {
			continue // linha sem arquivo, ou com código malformado, não é dep válida
		}
		deps = append(deps, Dep{
			Code: code,
			File: resolveDepPath(root, specRel, file),
			// Method PRESERVA as crases: elas são o sinal de que o autor prometeu um
			// SÍMBOLO (verificável pelo gate dependency-honored), não prosa descritiva.
			// Medido: 81% das células já usavam crase — remover apagava o contrato.
			Method: getRaw(ci.method),
			Layer:  get(ci.layer),
		})
	}
	return deps
}

type depCols struct{ code, file, method, layer int }

// columnIndex mapeia cada papel para o índice da coluna pela SEMÂNTICA do header
// (tolerante a acento/caixa e a colunas extras/reordenadas).
func columnIndex(header []string) depCols {
	ci := depCols{code: -1, file: -1, method: -1, layer: -1}
	for i, h := range header {
		switch norm(h) {
		case "cod", "codigo", "code":
			ci.code = i
		case "arquivo", "file", "origem":
			if ci.file < 0 {
				ci.file = i
			}
		case "metodo", "method", "simbolo", "funcao":
			ci.method = i
		case "camada", "layer", "tipo":
			ci.layer = i
		}
	}
	return ci
}

// resolveDepPath resolve o arquivo declarado na tabela para um caminho relativo à
// raiz. O autor escreve o path na convenção que lhe é natural (relativo à raiz, ou
// relativo ao `src/` do pacote com ou sem o próprio `src/` — como os imports `@/src/…`
// do projeto). Tenta as variantes em ordem e devolve a 1ª que existe no disco:
//  1. decl cru (já relativo à raiz)
//  2. <base até e incluindo /src/> + decl   (decl SEM o `src/`, ex.: "hooks/x.ts")
//  3. <base antes de /src/> + decl          (decl COM o `src/`, ex.: "src/hooks/x.ts")
//
// Se nada existir, devolve o valor original (o anti-drift do mapa acusa a aresta morta
// depois — melhor uma aresta a resolver que o silêncio).
func resolveDepPath(root, specRel, decl string) string {
	decl = strings.TrimPrefix(filepath.ToSlash(decl), "./")
	try := []string{decl}
	if i := strings.Index(specRel, "/src/"); i >= 0 {
		base := filepath.ToSlash(specRel[:i+len("/src/")]) // ".../src/"
		try = append(try,
			base+decl,                         // decl sem "src/"
			base[:len(base)-len("src/")]+decl, // decl com "src/" (evita "src/src/")
		)
	}
	for _, cand := range try {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(cand))); err == nil {
			return cand
		}
	}
	return decl
}

func splitRow(t string) []string {
	t = strings.Trim(t, "|")
	parts := strings.Split(t, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func isDividerRow(t string) bool {
	for _, c := range strings.Trim(t, "|") {
		if c != '-' && c != ':' && c != ' ' {
			return false
		}
	}
	return strings.Contains(t, "-")
}

// stripInlineCode tira o markdown de código inline (backticks) de uma célula.
func stripInlineCode(s string) string { return strings.Trim(s, "`") }

// norm normaliza um header de coluna: minúsculo, sem acento comum, só a 1ª palavra.
func norm(h string) string {
	h = strings.ToLower(strings.TrimSpace(stripInlineCode(h)))
	r := strings.NewReplacer("ó", "o", "é", "e", "ê", "e", "í", "i", "á", "a", "ç", "c", "ú", "u")
	h = r.Replace(h)
	if i := strings.IndexAny(h, " ("); i >= 0 {
		h = h[:i] // "cod (arquivo·metodo)" → "cod"
	}
	return h
}

// shortHash é o `rev` de um nó: a identidade do CONTEÚDO, e é ela que decide se um sinal
// ingerido (teste passou? mutante sobreviveu?) ainda vale e se uma aresta carimbada ficou
// velha.
//
// O `\r\n` sai ANTES do hash, e isso não é cosmética: com `core.autocrlf=true` o mesmo
// commit tem bytes diferentes no Windows e no macOS, então o rev de TODO arquivo divergia
// entre as duas máquinas. O efeito aparecia longe da causa — um `map build` no Windows
// reescrevia os 2799 revs, o `preservarSinais` não reconhecia nenhum nó e DESCARTAVA os
// 1061 sinais acumulados; commitado, o mapa fazia o mesmo estrago do outro lado.
//
// MEDIDO em 24/08 no app de referência: `entitlementDecision.feature` dava `9d743c6d9abb` em bytes
// brutos e `daa7cef87644` normalizado — e é `daa7cef87644` que está no mapa gerado no
// macOS. Ou seja, normalizar não migra nada: faz o Windows CONCORDAR com o que já existe.
func shortHash(b []byte) string {
	sum := sha256.Sum256(bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n")))
	return hex.EncodeToString(sum[:])[:12]
}

// ShortHash expõe o hash de revisão (mesma fórmula do scan) para o watcher
// recalcular o rev de um único arquivo alterado.
func ShortHash(b []byte) string { return shortHash(b) }

// headerCodeRE captura a identidade declarada no header `@anchors`. Os prefixos cobrem os
// dialetos de comentário dos artefatos (HTML/Markdown, `//`, `#`).
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func headerCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

// extractHeaderCode devolve a identidade DECLARADA, ou vazio se o header não a declara.
func extractHeaderCode(content string) string {
	if m := headerCodeRE().FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// planSeedRE acha os caminhos de spec que um plano SEMEIA (cita entre crases).
var planSeedRE = regexp.MustCompile("`([^`]+\\.spec\\.md)`")

// extractSeeds devolve as specs que um plano semeia. Vazio para qualquer outro kind — só
// plano semeia.
func extractSeeds(kind, content string) []string {
	if kind != "plan" {
		return nil
	}
	visto := map[string]bool{}
	var out []string
	for _, m := range planSeedRE.FindAllStringSubmatch(content, -1) {
		// Um caminho com metacaractere de glob (`*.spec.md`) é PROSA sobre um conjunto
		// ("todo `*.spec.md` precisa de header"), não a citação de um arquivo que vai
		// nascer. Tratá-lo como seed fazia o plano parecer eternamente não-cumprido,
		// porque o "arquivo" nunca existiria.
		if visto[m[1]] || strings.HasPrefix(filepath.Base(m[1]), "_TEMPLATE") ||
			strings.ContainsAny(m[1], "*?[{") {
			continue
		}
		visto[m[1]] = true
		out = append(out, m[1])
	}
	return out
}
