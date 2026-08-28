package scan

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
)

// Ignore decide o que o Anchors NÃO deve enxergar.
//
// A lista embutida cobre o que é universal (node_modules, .git, dist…). O resto vem do
// `.gitignore` do projeto — que é a fonte que a equipe JÁ mantém dizendo o que é
// descartável, e que o Anchors ignorava.
//
// O caso que motivou: um revisor criou sondas (`probe1.ts`) para atacar a unidade por
// execução, e o watcher as enfileirou como trabalho a fazer. O orquestrador teve de
// descartar cinco tasks à mão. As sondas eram instrumento, não entrega — e o
// `.gitignore` do projeto já dizia isso.
//
// Por que não inventar um `ignore:` no anchors.yaml: seria uma segunda lista sobre o
// mesmo assunto, que envelheceria em relação à primeira. O padrão medido nesta sessão é
// que **duas fontes da régua discordando custa mais que uma fonte ausente** — e aqui a
// fonte já existe. O `anchors.yaml` pode ESTENDER (via `comments`-style override), nunca
// substituir.
type Ignore struct {
	dirs     map[string]bool
	patterns []ignorePattern
}

type ignorePattern struct {
	glob    string
	negate  bool // linha `!padrão`: re-inclui o que uma regra anterior excluiu
	dirOnly bool
	// ancorado: a linha começava com `/` (ou já continha barra no meio). No git isso
	// significa "a partir da RAIZ", e não "em qualquer nível" — `/data/` ignora `data/`
	// na raiz e NÃO ignora `amplify/data/`.
	//
	// Descartar essa distinção custou caro: o `TrimPrefix(glob, "/")` fazia a linha
	// `/data/` do projeto casar o segmento intermediário de `amplify/data/models`, e 129
	// unidades — as specs e os modelos do schema — sumiram do mapa em silêncio. O mapa
	// não fica errado com barulho: ele fica menor, e nada acusa a ausência.
	ancorado bool
}

// universalIgnored são diretórios cujo nome quase nunca é material do projeto. Ficam
// embutidos porque a varredura deles é cara e, na esmagadora maioria dos casos, inútil.
//
// "Quase nunca" não é "nunca", e a diferença é o ponto: `build` é saída de compilação na
// maioria dos projetos e é PASTA DE CÓDIGO em alguns; `dist` idem. Uma lista embutida que
// não pode ser contestada é o Anchors decidindo, por um projeto que ele não conhece, o que
// nele é descartável — e o custo do erro é silencioso: a camada inteira some do mapa e
// nenhum gate acusa, porque para eles aquelas unidades nunca existiram.
//
// Por isso a lista é DERROTÁVEL: uma negação explícita no `.gitignore` (`!build/`) ou uma
// camada da Estrutura cujo pattern aponte para dentro dela vencem o default. O projeto que
// diz "aqui é código" tem a palavra final sobre a suposição do framework.
//
// `.git` e `.anchors` NÃO são derrotáveis: são a maquinaria, não o material.
var universalIgnored = map[string]bool{
	"node_modules": true, "dist": true, "build": true,
	"vendor": true, ".next": true, "coverage": true, ".expo": true,
}

// efemeros são ARQUIVOS que nascem do editor/sistema e desaparecem sozinhos — nunca são
// material do projeto, em projeto nenhum.
//
// Diferente de `build/` ou `dist/`, aqui não há ambiguidade a respeitar: ninguém versiona
// o swap do vim nem o `.!21662!` que o editor escreve durante um salvamento atômico. E o
// `.gitignore` do projeto não cobre isso justamente porque não é decisão do projeto — é
// ruído do sistema de arquivos.
//
// Medido num E2E real: o watcher enfileirou uma task para
// `amplify/data/.!21662!resource.spec.md`, arquivo que existiu por milissegundos durante
// um salvamento. O orquestrador teve de descartá-la à mão, junto com outras duas.
var efemeros = []string{
	".!*!*",                  // salvamento atômico (o padrão que apareceu no E2E)
	"*.swp",                  // vim
	"*.swo",                  // vim
	"*~",                     // backup de editor (emacs, gedit)
	".#*",                    // lock do emacs
	"#*#",                    // autosave do emacs
	"*.tmp",                  // genérico
	".DS_Store", "Thumbs.db", // sistema
}

// efemero diz se o BASENAME é um arquivo transitório de editor/sistema.
func efemero(rel string) bool {
	base := filepath.Base(rel)
	for _, p := range efemeros {
		if ok, _ := doublestar.Match(p, base); ok {
			return true
		}
	}
	return false
}

// maquinaria são os diretórios que nunca são material, sob nenhuma configuração — o
// versionamento e o estado do próprio Anchors. Não há projeto em que varrê-los faça
// sentido, então nenhuma declaração os reabilita.
var maquinaria = map[string]bool{".git": true, ".anchors": true}

// registroDoAnchors são os diretórios que o PRÓPRIO Anchors escreve como saída: issues,
// registros de entrega, relatórios. São material versionado do projeto (por isso não vivem
// em `.anchors/`), mas não são material a TRABALHAR — e a diferença importa.
//
// Medido: o watcher enfileirou tasks de `triage` sobre os arquivos de `issues/todo/` que
// o `anchors judge` tinha acabado de escrever. O ciclo passou a se alimentar da própria
// saída — cada review gerava uma issue, que gerava uma task, que ninguém podia executar.
// Três execuções descartaram essas tasks à mão.
var registroDoAnchors = map[string]bool{"issues": true, "changes": true}

// LoadIgnore lê o `.gitignore` da raiz e compõe com a lista universal.
// Ausência do arquivo não é erro — só significa que o projeto não declarou nada.
func LoadIgnore(root string) *Ignore {
	return LoadIgnoreFor(root, nil)
}

// LoadIgnoreFor compõe a lista de exclusão sabendo o que a ESTRUTURA declara. Uma camada
// cujo pattern aponte para dentro de um diretório da lista embutida (`build/**/*.ts` num
// projeto onde `build` é código-fonte) DERROTA o default: o projeto declarou que ali há
// material, e a suposição do framework cede à declaração.
func LoadIgnoreFor(root string, cfg *config.Config) *Ignore {
	ig := &Ignore{dirs: map[string]bool{}}
	reabilitados := diretoriosReabilitadosPelaEstrutura(cfg)
	for d := range universalIgnored {
		if !reabilitados[d] {
			ig.dirs[d] = true
		}
	}
	for d := range maquinaria {
		ig.dirs[d] = true
	}
	f, err := os.Open(filepath.Join(root, ".gitignore"))
	if err != nil {
		return ig
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	// Uma negação no `.gitignore` (`!build/`) também derrota a lista embutida — é a forma
	// que a equipe já conhece para dizer "isto aqui, apesar do nome, não é descartável".
	// Sem isso a negação seria ignorada em silêncio, que é a falha que a lista causa.
	for sc.Scan() {
		linha := strings.TrimSpace(sc.Text())
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}
		p := ignorePattern{glob: linha}
		if strings.HasPrefix(linha, "!") {
			p.negate, p.glob = true, strings.TrimPrefix(linha, "!")
		}
		if strings.HasSuffix(p.glob, "/") {
			p.dirOnly, p.glob = true, strings.TrimSuffix(p.glob, "/")
		}
		// Barra inicial OU barra no meio ancoram na raiz (regra do gitignore). Só o
		// padrão sem barra alguma (`*.log`, `node_modules`) casa em qualquer nível.
		if strings.HasPrefix(p.glob, "/") {
			p.ancorado, p.glob = true, strings.TrimPrefix(p.glob, "/")
		} else if strings.Contains(strings.TrimSuffix(p.glob, "/"), "/") {
			p.ancorado = true
		}
		if p.negate {
			// a maquinaria resiste: `!.git/` não faz o Anchors varrer o repositório git.
			if nome := filepath.Base(p.glob); !maquinaria[nome] {
				delete(ig.dirs, nome)
			}
		}
		ig.patterns = append(ig.patterns, p)
	}
	return ig
}

// SkipDir diz se um DIRETÓRIO inteiro deve ser pulado. Pular no diretório é o que torna a
// varredura barata — descer em `node_modules` para depois descartar cada arquivo custaria
// segundos em todo comando.
func (ig *Ignore) SkipDir(nome, rel string) bool {
	// O que o Anchors escreve não é trabalho para o Anchors.
	if registroDoAnchors[nome] {
		return true
	}
	if ig == nil {
		return universalIgnored[nome]
	}
	if ig.dirs[nome] {
		return true
	}
	return ig.match(rel, true)
}

// SkipFile diz se um ARQUIVO deve ser ignorado.
func (ig *Ignore) SkipFile(rel string) bool {
	if efemero(rel) {
		return true // ruído do editor: nunca é material, em projeto nenhum
	}
	if ig == nil {
		return false
	}
	return ig.match(rel, false)
}

// match aplica os padrões na ordem, como o git: a última regra que casa vence, e uma
// negação (`!padrão`) re-inclui o que veio antes.
func (ig *Ignore) match(rel string, isDir bool) bool {
	rel = filepath.ToSlash(rel)
	ignorado := false
	for _, p := range ig.patterns {
		// `dirOnly` (`data/`) restringe o que o padrão CASA, não o que ele cobre: o
		// diretório é ignorado e tudo abaixo dele vai junto. Pular todo arquivo aqui
		// deixava `data/dump.json` visível mesmo com `data/` ignorado — e o custo é o
		// oposto do bug anterior: material descartável voltando a virar trabalho.
		if p.dirOnly && !isDir && !cobreAbaixo(p, rel) {
			continue
		}
		if !casaGitignore(p.glob, rel, p.ancorado) {
			continue
		}
		ignorado = !p.negate
	}
	return ignorado
}

// casaGitignore reproduz o essencial da semântica do git.
//
// A distinção que importa é ANCORAGEM: um padrão sem barra alguma (`*.log`,
// `node_modules`) casa em qualquer nível; um padrão com barra (`/data/`, `docs/build`)
// casa a partir da RAIZ. Tratar os dois igual foi o que fez `/data/` comer
// `amplify/data/models`.
func casaGitignore(glob, rel string, ancorado bool) bool {
	if ancorado {
		if ok, _ := doublestar.Match(glob, rel); ok {
			return true
		}
		// um diretório ignorado cobre tudo abaixo dele
		return strings.HasPrefix(rel, strings.TrimSuffix(glob, "/")+"/")
	}
	// sem âncora: casa o basename em qualquer profundidade — e, quando casa um
	// DIRETÓRIO intermediário, cobre o que está abaixo (`node_modules` pega
	// `a/node_modules/b.js`).
	if ok, _ := doublestar.Match(glob, filepath.Base(rel)); ok {
		return true
	}
	for _, seg := range strings.Split(rel, "/") {
		if ok, _ := doublestar.Match(glob, seg); ok {
			return true
		}
	}
	return false
}

// cobreAbaixo diz se `rel` está DENTRO de um diretório que este padrão ignora. É o que
// faz `data/` cobrir `data/dump.json` sem que o padrão precise casar o arquivo em si.
func cobreAbaixo(p ignorePattern, rel string) bool {
	dir := filepath.ToSlash(filepath.Dir(rel))
	for dir != "." && dir != "/" && dir != "" {
		if casaGitignore(p.glob, dir, p.ancorado) {
			return true
		}
		dir = filepath.ToSlash(filepath.Dir(dir))
	}
	return false
}

// diretoriosReabilitadosPelaEstrutura devolve os nomes da lista embutida que alguma camada
// declara como material. O sinal é o pattern apontar PARA DENTRO do diretório: uma camada
// `build/**/*.ts` diz que `build` é código; uma camada `**/*.ts` que por acaso o alcança,
// não — senão qualquer catch-all reabilitaria `node_modules`.
func diretoriosReabilitadosPelaEstrutura(cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	if cfg == nil {
		return out
	}
	for _, l := range cfg.Layers {
		for _, seg := range strings.Split(filepath.ToSlash(l.Pattern), "/") {
			if universalIgnored[seg] {
				out[seg] = true
			}
		}
	}
	return out
}
