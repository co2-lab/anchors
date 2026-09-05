package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// `anchors new <kind> <nome>` emite o ESQUELETO de um artefato da trinca (spec,
// feature, test) já com o cabeçalho @anchors e a identidade (code/ref) resolvidos —
// o piso que o `check` exige. O anchors "não gera conteúdo", mas gerar a MOLDURA
// correta (header + seções obrigatórias) é estrutura, não conteúdo: tira o atrito de
// começar um arquivo e garante que ele nasce conforme.
//
// Parametrizável por kind: cada kind tem seções DEFAULT (entram sempre) e OPCIONAIS
// (entram com --with). --without remove uma default. --list-sections mostra o menu.
func newNewCmd() *cobra.Command {
	var root, codeStr, out, preset string
	var with, without []string
	var listSections bool

	cmd := &cobra.Command{
		Use:   "new <kind> <nome>",
		Short: "Emite o esqueleto de um artefato (spec/feature/test) conforme a régua",
		Long: `Gera a moldura de um artefato da trinca com o cabeçalho @anchors e a
identidade já resolvidos:

  anchors new spec Login --out src/screens/Login.spec.md      → code novo + seções default
  anchors new spec Calc --preset backend-logic --out src/business-logic/Calc.spec.md
  anchors new spec Login --out <path> --with auth,deps
  anchors new spec Card  --out <path> --without overview
  anchors new feature Login --out <path> --code ABCD
  anchors new spec --list-sections    → seções (com "quando usar") e presets do kind

O --out é obrigatório: o artefato nasce JUNTO da unidade que descreve (a spec ao lado
do código, a feature ao lado da spec) — nunca na raiz do repositório.

Kinds: spec, feature, test. A spec ganha um CODE novo (único no mapa); feature e test
ganham um REF (apontam para a spec). Use --code para fixar a identidade à mão.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := strings.ToLower(args[0])
			tpl, ok := templates[kind]
			if !ok {
				return fmt.Errorf("kind desconhecido %q — use: spec, feature ou test", kind)
			}

			if listSections {
				printSections(kind, tpl)
				return nil
			}
			if len(args) < 2 {
				return fmt.Errorf("informe o nome da unidade (ex: `anchors new %s Login`)", kind)
			}
			name := args[1]

			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}

			// identidade: spec é dona (code:), feature/test REFERENCIAM (ref:).
			//
			// A distinção decide de onde a identidade vem, e errá-la produz um órfão por
			// construção: gerar um código NOVO para uma feature faz o `ref:` apontar para
			// uma spec que não existe. Aconteceu — `anchors new feature metadataVersioning`
			// cunhou `MTVA` em vez de referenciar a spec irmã, que declarava `MTVR`. Quem
			// confiasse no scaffold criaria a peça já desconectada da trinca.
			//
			// Então: para feature/test, a identidade é LIDA da spec irmã (o `--out` diz
			// onde a peça nasce; a spec correspondente está no lugar que a Estrutura
			// manda). Só se não houver spec é que se gera — e aí com aviso, porque o mais
			// provável é que o caminho esteja errado, não que a spec deva nascer depois.
			id := strings.ToUpper(codeStr)
			if id == "" && kind != "spec" {
				if ref, origem := refDaSpecIrma(absRoot, out, name); ref != "" {
					id = ref
					fmt.Printf("identidade: `%s` (lida de %s)\n", id, origem)
				}
			}
			if id == "" {
				id, err = resolveNewCode(absRoot, name)
				if err != nil {
					return err
				}
				if kind != "spec" {
					fmt.Printf("aviso: nenhuma spec encontrada para este alvo — gerei o código `%s`.\n"+
						"  Uma %s referencia a spec (`ref:`); sem spec, ela nasce ÓRFÃ e o gate\n"+
						"  trinca-completa vai acusar. Confira o --out, ou crie a spec antes.\n", id, kind)
				}
			}

			sections, ordem, err := resolveSectionsWithPreset(tpl, preset, with, without)
			if err != nil {
				return err
			}

			outPath := out
			if outPath == "" {
				// A RAIZ do repo é o lugar errado em praticamente todo caso: o artefato
				// mora COM a unidade que ele descreve (a spec ao lado do .ts, a feature
				// ao lado da spec). Gravar na raiz por default produz um arquivo que o
				// próprio `check` depois acusa de fora de lugar — então recusamos e
				// dizemos o que fazer, em vez de adivinhar o diretório errado em silêncio.
				return fmt.Errorf("informe onde o artefato deve nascer com --out "+
					"(ex.: --out packages/backend/business-logic/%s%s). O artefato mora "+
					"junto da unidade que ele descreve, não na raiz do repositório", name, tpl.ext)
			}
			if !filepath.IsAbs(outPath) {
				outPath = filepath.Join(absRoot, outPath)
			}

			// Só agora: o dialeto de comentário do header vem da EXTENSÃO do arquivo de
			// destino, então o render precisa do caminho já resolvido.
			// A config traz o DIALETO (para o esqueleto de teste nascer na linguagem certa)
			// e as camadas (para recusar o artefato em camada declarativa). Sem config, o
			// `new` continua funcionando com os defaults — é comando de bootstrap.
			cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile))
			content := renderArtifact(tpl, name, id, outPath, absRoot, sections, ordem, cfg)
			// A camada do ALVO decide se essa spec pode existir. Recusar aqui é o
			// ponto mais barato: antes do arquivo nascer.
			if cfg != nil {
				if err := refuseIfRecognizedLayer(absRoot, outPath, cfg, kind); err != nil {
					return err
				}
			}
			if _, statErr := os.Stat(outPath); statErr == nil {
				return fmt.Errorf("%s já existe — não sobrescrevo (use --out para outro caminho)", outPath)
			}
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
				return err
			}
			fmt.Printf("✓ %s criado (%s: %s)\n", outPath, tpl.idField, id)

			// O PLANO nasce com o companheiro de progresso. Ver `progress.go`: o
			// checkbox no plano fazia "terminei uma fase" ser indistinguível de "mudei a
			// decisão", e derrubava o julgamento da spec que a fase entregou.
			//
			// Criar aqui, e não pedir ao usuário, é o que faz a separação ser o default.
			// Um plano sem companheiro convida a marcar `[x]` no próprio plano — que é o
			// caminho que já existia e o que se está removendo.
			if kind == "plan" {
				prog, err := escreveProgressoInicial(outPath, content, id)
				if err != nil {
					// Não é falha do `new`: o plano nasceu. O progresso se cria à mão.
					fmt.Printf("  ⚠ progresso não criado: %v\n", err)
				} else {
					fmt.Printf("✓ %s criado (o ESTADO do plano; marque `[x]` aqui, nunca no plano)\n",
						relTo(absRoot, prog))
				}
			}

			fmt.Println("  preencha as seções e rode `anchors check --changed " + relTo(absRoot, outPath) + "`")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&codeStr, "code", "", "usa este código de identidade (senão gera um único)")
	cmd.Flags().StringVar(&out, "out", "", "caminho de saída (default: <nome><ext> na raiz)")
	cmd.Flags().StringSliceVar(&with, "with", nil, "adiciona seções OPCIONAIS (csv)")
	cmd.Flags().StringSliceVar(&without, "without", nil, "remove seções DEFAULT (csv)")
	cmd.Flags().StringVar(&preset, "preset", "", "conjunto de seções de um tipo de unidade (ver --list-sections)")
	cmd.Flags().BoolVar(&listSections, "list-sections", false, "lista as seções do kind e sai")
	return cmd
}

// section é um bloco nomeado do artefato. Default entra sempre (salvo --without);
// !Default entra só com --with. body pode conter os placeholders {name} e {id}.
type section struct {
	Key     string
	Title   string
	Default bool
	Body    string
	// Purpose diz QUANDO usar esta seção — o critério de escolha, não o que ela contém.
	// Seções VARIANTES (ex.: `contract` vs `signature`) descrevem o mesmo papel em
	// dialetos diferentes; quem escreve a spec decide qual cabe, e o Purpose torna essa
	// decisão informada em vez de chute.
	Purpose string
	// Variants nomeia as outras seções que cobrem o MESMO papel (alternativas). O
	// --list-sections as mostra juntas, para a escolha ser consciente.
	Variants []string
	// Realizes é a LETRA de tipo de regra que esta seção cataloga (`B` para efeitos,
	// `X` para restrições…). Serve para traduzir o título GENÉRICO do template para o
	// nome que o PROJETO usa: `rule_types` já declara `X → "Modelo de Dado"`, e sem essa
	// ligação o preset emite "Restrições" num projeto onde 50 specs vizinhas escrevem
	// "Modelo de Dado".
	//
	// O problema que resolve, medido: o preset `schema` emitia Contrato/Regras/Auth
	// enquanto os 50 vizinhos usavam Modelo de Dado/Comportamentos/Permissões e Acesso.
	// Duas fontes da própria régua discordando — e o CLI não avisava. Renomear as seções
	// no framework seria pior: poria a convenção de um projeto dentro do engine.
	Realizes string
}

// template descreve a moldura de um kind: dialeto de comentário do header, campo de
// identidade (code|ref), extensão e o catálogo de seções (ordenado).
type template struct {
	kind     string
	ext      string
	idField  string // "code" (dono) ou "ref" (referencia a spec)
	headerFn func(id, outPath string) string
	sections []section
}

// renderArtifact monta o arquivo: header @anchors + seções escolhidas, na ordem do
// catálogo. Substitui {name}/{id} no corpo de cada seção.
// outPath entra porque o dialeto de comentário do header depende da EXTENSÃO real do
// arquivo — um `.py` não leva `//`.
func renderArtifact(t template, name, id, outPath, root string, chosen map[string]bool, ordem []string, cfg *config.Config) string {
	lang, kw := cfg.DialectFor().GherkinFor()
	var b strings.Builder
	// a feature recebe o IDIOMA (para a linha `# language:`); os demais, o caminho (para
	// o dialeto de comentário). Cada headerFn usa o que lhe interessa.
	headerArg := outPath
	if t.kind == "feature" {
		headerArg = lang
	}
	b.WriteString(t.headerFn(id, headerArg))
	// A ordem das seções é o FIO DE LEITURA da spec, então ela é do preset quando há um:
	// numa spec de store, "Shape do Estado" precisa vir antes de "Invariantes" (não se
	// pode enunciar invariante sobre um estado que o leitor ainda não conhece). Sem
	// preset, vale a ordem do catálogo.
	// O léxico de seções é da CAMADA do alvo (ver `section_titles`): resolver uma vez,
	// fora do loop.
	camadaDoArtefato := camadaDoAlvo(root, outPath, cfg)
	for _, s := range ordenaSecoes(t, chosen, ordem) {
		body := strings.NewReplacer("{name}", name, "{id}", id).Replace(s.Body)
		// Traduz o título GENÉRICO para o nome que ESTE projeto usa, quando `rule_types`
		// declara um. Sem isso o preset emitia "Restrições" num projeto cujas 50 specs
		// vizinhas escrevem "Modelo de Dado" — duas fontes da própria régua discordando,
		// e o autor tendo de escolher entre obedecer o template ou os vizinhos.
		body = traduzTitulo(body, s, cfg, camadaDoArtefato)
		// {TEST_BODY} é resolvido pelo DIALETO: o esqueleto de caso de teste tem sintaxe,
		// e sintaxe é do projeto (ver testBody em new_templates.go).
		if strings.Contains(body, "{TEST_BODY}") {
			body = strings.ReplaceAll(body, "{TEST_BODY}", testBody(cfg.DialectFor().Family, name, id))
		}
		// A tag de REGIME é do projeto (`derived.regimes` faz o de-para para o regime
		// canônico). Cravar `@nivel-unit` gerava um cenário que o gate de correspondência
		// não confronta em projeto nenhum que use outro vocabulário — e o `work` já
		// ensinava que "a tag é do PROJETO e não é traduzível". O template contradizia
		// a própria régua.
		body = strings.ReplaceAll(body, "{UNIT_TAG}", unitRegimeTag(cfg))
		// palavras-chave do Gherkin no idioma do projeto (default `en`).
		body = strings.NewReplacer(
			"{FEATURE}", kw.Feature, "{SCENARIO}", kw.Scenario, "{OUTLINE}", kw.Outline,
			"{GIVEN}", kw.Given, "{WHEN}", kw.When, "{THEN}", kw.Then, "{EXAMPLES}", kw.Examples,
		).Replace(body)
		b.WriteString(body)
	}
	return b.String()
}

// resolveSections aplica with/without ao default do template, validando os nomes.
func resolveSections(t template, with, without []string) (map[string]bool, error) {
	valid := map[string]bool{}
	for _, s := range t.sections {
		valid[s.Key] = true
	}
	for _, w := range append(append([]string{}, with...), without...) {
		if w != "" && !valid[w] {
			return nil, fmt.Errorf("seção desconhecida %q — veja `anchors new %s --list-sections`", w, t.kind)
		}
	}
	chosen := map[string]bool{}
	for _, s := range t.sections {
		chosen[s.Key] = s.Default
	}
	for _, w := range with {
		if w != "" {
			chosen[w] = true
		}
	}
	for _, w := range without {
		if w != "" {
			chosen[w] = false
		}
	}
	return chosen, nil
}

// resolveSectionsWithPreset resolve as seções partindo de um PRESET (se dado) em vez
// dos defaults do kind. --with/--without seguem valendo por cima, para refinar.
func resolveSectionsWithPreset(t template, preset string, with, without []string) (map[string]bool, []string, error) {
	if preset == "" {
		c, err := resolveSections(t, with, without)
		return c, nil, err // sem preset: ordem do catálogo
	}
	if t.kind != "spec" {
		return nil, nil, fmt.Errorf("--preset só se aplica a `spec` (kind atual: %s)", t.kind)
	}
	def, ok := specPresets[preset]
	if !ok {
		return nil, nil, fmt.Errorf("preset desconhecido %q — veja `anchors new spec --list-sections`", preset)
	}
	valid := map[string]bool{}
	for _, s := range t.sections {
		valid[s.Key] = true
	}
	for _, w := range append(append([]string{}, with...), without...) {
		if w != "" && !valid[w] {
			return nil, nil, fmt.Errorf("seção desconhecida %q — veja `anchors new %s --list-sections`", w, t.kind)
		}
	}
	chosen := map[string]bool{}
	for _, s := range t.sections {
		chosen[s.Key] = false
	}
	for _, k := range def.Sections {
		chosen[k] = true
	}
	for _, w := range with {
		if w != "" {
			chosen[w] = true
		}
	}
	for _, w := range without {
		if w != "" {
			chosen[w] = false
		}
	}
	return chosen, def.Sections, nil
}

// ordenaSecoes devolve as seções escolhidas na ORDEM do preset (quando há um), com as
// adicionadas por --with logo depois — na ordem do catálogo, que é o único critério
// disponível para quem o preset não previu.
func ordenaSecoes(t template, chosen map[string]bool, ordem []string) []section {
	porChave := map[string]section{}
	for _, s := range t.sections {
		porChave[s.Key] = s
	}
	var out []section
	visto := map[string]bool{}
	for _, k := range ordem {
		if chosen[k] && !visto[k] {
			if s, ok := porChave[k]; ok {
				out = append(out, s)
				visto[k] = true
			}
		}
	}
	for _, s := range t.sections {
		if chosen[s.Key] && !visto[s.Key] {
			out = append(out, s)
			visto[s.Key] = true
		}
	}
	return out
}

func printSections(kind string, t template) {
	fmt.Printf("Seções de `%s` (default entra sempre; opcional entra com --with):\n\n", kind)
	for _, s := range t.sections {
		tag := "opcional"
		if s.Default {
			tag = "default "
		}
		fmt.Printf("  [%s] %-12s %s\n", tag, s.Key, s.Title)
		if s.Purpose != "" {
			fmt.Printf("             ↳ quando usar: %s\n", s.Purpose)
		}
		if len(s.Variants) > 0 {
			fmt.Printf("             ↳ alternativa a: %s (escolha UMA)\n", strings.Join(s.Variants, ", "))
		}
	}
	if kind == "spec" && len(specPresets) > 0 {
		fmt.Println("\nPresets (conjuntos prontos — `--preset <nome>`):")
		names := make([]string, 0, len(specPresets))
		for n := range specPresets {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			p := specPresets[n]
			fmt.Printf("  %-15s %s\n", n, p.Desc)
			fmt.Printf("  %-15s   seções: %s\n", "", strings.Join(p.Sections, ", "))
		}
	}
	fmt.Println("\n  --with a,b     adiciona opcionais    --without a,b   remove defaults")
	fmt.Println("  --preset nome  parte de um conjunto pronto (refina com --with/--without)")
}

// resolveNewCode gera um código único para a unidade, ligado à Estrutura (usa o
// prefixo de módulo da layer se houver) — mesma lógica de `anchors code`.
func resolveNewCode(root, nameOrPath string) (string, error) {
	mapPath := filepath.Join(root, mapx.DefaultPath)
	taken, _, err := takenCodes(mapPath)
	if err != nil {
		// sem mapa (projeto novo): gera sem checar unicidade — melhor que travar.
		taken = map[string]bool{}
	}
	// A config vem ANTES de qualquer geração, e fora do ramo do caminho: é o `Load` que
	// dispara o SetSlotsHook, e sem ele o gerador emite no comprimento canônico (5) mesmo
	// num projeto que declara `code_lengths: [4]`. O `Load` já morou dentro do `if` abaixo,
	// e o efeito era `anchors new spec receiptMatch` → RCMTR (5) enquanto o MESMO comando
	// com caminho → RCMT (4): a identidade dependia da forma do argumento, não do projeto.
	cfg, cfgErr := config.Load(filepath.Join(root, config.DefaultFile))
	// Camada 1 (code_prefix da layer) tem precedência; senão o caminho decide via
	// GenerateFromPath (dir-pai p/ basename genérico ou desempate de colisão).
	if strings.ContainsRune(nameOrPath, '/') {
		prefix := ""
		if cfgErr == nil {
			rel := relTo(root, nameOrPath)
			if layer, _ := scan.Classify(rel, cfg); layer != "" {
				if l, ok := cfg.Layers[layer]; ok {
					prefix = l.CodePrefix
				}
			}
		}
		if prefix != "" {
			return code.GenerateUniqueWithPrefix(unitName(nameOrPath), prefix, taken), nil
		}
		return code.GenerateFromPath(nameOrPath, taken), nil
	}
	return code.GenerateUnique(nameOrPath, taken), nil
}

// refuseIfRecognizedLayer barra a criação de uma SPEC para um alvo de camada
// RECONHECIDA (`regime: declarativo` — dao/infra/resource). Essas camadas são
// declaradas justamente para SAIR do escrutínio: elas não originam regra, logo não têm
// spec. Sem esta guarda a ferramenta gera alegremente o arquivo que a doutrina proíbe —
// e a contradição só aparece muito depois, na revisão (ou nunca).
//
// A checagem é sobre o ALVO que a spec descreveria (o `.ts` irmão), não sobre o
// caminho da própria spec.
func refuseIfRecognizedLayer(root, outPath string, cfg *config.Config, kind string) error {
	if kind != "spec" || cfg == nil {
		return nil
	}
	rel, err := filepath.Rel(root, outPath)
	if err != nil {
		return nil // fora da raiz: não é nosso caso
	}
	// o alvo é o irmão sem o sufixo .spec.md — tentamos as extensões usuais.
	base := strings.TrimSuffix(rel, ".spec.md")
	for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
		layer, _ := scan.Classify(base+ext, cfg)
		if layer == "" {
			continue
		}
		l, ok := cfg.Layers[layer]
		if !ok || l.Regime != "declarativo" {
			continue
		}
		return fmt.Errorf("`%s` é camada RECONHECIDA (regime: declarativo) — não tem spec.\n"+
			"  Camadas assim são declaradas no anchors.yaml justamente para sair do escrutínio: "+
			"elas não originam regra (só traduzem/configuram), então não há o que especificar.\n"+
			"  Se há uma DECISÃO a documentar, ela pertence à camada que decide (ex.: o schema "+
			"decide o modelo; o DAO só o traduz). Ver `layers.%s` no anchors.yaml.", layer, layer)
	}
	return nil
}

// unitRegimeTag devolve a tag que ESTE projeto usa para o regime unitário, lida do
// de-para `derived.regimes`. Sem de-para declarado, emite um marcador TODO em vez de
// chutar uma tag: um cenário com a tag errada passa despercebido (nenhum gate o
// confronta), enquanto um TODO visível é corrigido na hora.
func unitRegimeTag(cfg *config.Config) string {
	if cfg != nil && cfg.Derived != nil {
		// ordena para o resultado ser estável entre execuções (mapa não tem ordem).
		var tags []string
		for tag, canônico := range cfg.Derived.Regimes {
			if canônico == "unit" {
				tags = append(tags, tag)
			}
		}
		if len(tags) > 0 {
			sort.Strings(tags)
			return "@" + tags[0]
		}
	}
	return "@TODO-tag-de-regime (declare o de-para em `derived.regimes` do anchors.yaml)"
}

// refDaSpecIrma acha a identidade que uma feature/teste deve REFERENCIAR: o `code:` da
// spec da mesma unidade. Procura onde a Estrutura manda a spec morar — ao lado da peça
// (co-location) — tentando os sufixos usuais sobre o mesmo tronco de nome.
//
// Devolve (código, caminho-de-onde-veio) ou ("", "") se não achar.
func refDaSpecIrma(root, outPath, name string) (string, string) {
	if outPath == "" {
		return "", ""
	}
	abs := outPath
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(root, outPath)
	}
	dir := filepath.Dir(abs)
	// tronco: o nome do arquivo sem a extensão composta (`x.feature`, `x.test.ts`).
	base := filepath.Base(abs)
	for _, suf := range []string{".feature", ".test.ts", ".test.tsx", ".spec.ts", "_test.go", "_test.py", "_spec.rb"} {
		base = strings.TrimSuffix(base, suf)
	}
	if i := strings.Index(base, "."); i > 0 {
		base = base[:i]
	}
	if base == "" {
		base = name
	}
	cand := filepath.Join(dir, base+".spec.md")
	b, err := os.ReadFile(cand)
	if err != nil {
		return "", ""
	}
	if c := codeDoHeaderSpec(string(b)); c != "" {
		rel, rerr := filepath.Rel(root, cand)
		if rerr != nil {
			rel = cand
		}
		return c, rel
	}
	return "", ""
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func specHeaderCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

func codeDoHeaderSpec(content string) string {
	if m := specHeaderCodeRE().FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// traduzTitulo substitui o cabeçalho `## <genérico>` pelo nome que o projeto usa para a
// mesma letra de tipo de regra.
//
// A fonte é `rule_types[].sections` — que o projeto JÁ declara para o gate `rule-types`.
// Nada novo a configurar: a informação existia e não estava ligada ao `new`.
//
// Usa o PRIMEIRO título declarado para a letra, que é a convenção dominante do projeto.
// Se o projeto não declara a letra, o título genérico fica — o framework não inventa nome.
func traduzTitulo(body string, s section, cfg *config.Config, camada string) string {
	if cfg == nil {
		return body
	}
	// 1) `section_titles` renomeia POR CHAVE — cobre qualquer seção do catálogo.
	//
	// A tradução por `rule_types` (abaixo) só alcança seções que realizam uma letra de
	// regra, e era a única que existia. As três seções que divergiram num E2E real —
	// `contract`, `domain`, `effects` — não têm `Realizes`, então nunca eram traduzidas:
	// o preset emitia "Contrato/Domínio/Efeitos" onde as 50 specs vizinhas escreviam
	// "Modelo de Dado/Comportamentos/Notas de Implementação", e NENHUMA das três
	// coincidia. Dois agentes escreveram em dialetos opostos, cada um obedecendo a uma
	// fonte da régua, os dois verdes nos gates.
	local := cfg.TituloDaSecao(s.Key, "", camada)
	// 2) `rule_types.sections` continua valendo para as seções de regra — é onde o
	// projeto já declarava o nome da seção junto com a letra que ela cataloga.
	if local == "" && s.Realizes != "" {
		for _, rt := range cfg.RuleTypes {
			if strings.EqualFold(rt.Letter, s.Realizes) && len(rt.Sections) > 0 {
				local = rt.Sections[0]
				break
			}
		}
	}
	if local == "" {
		return body
	}
	m := tituloSecaoRE.FindStringSubmatch(body)
	if m == nil || strings.EqualFold(strings.TrimSpace(m[2]), local) {
		return body
	}
	return strings.Replace(body, m[0], m[1]+" "+local+"\n", 1)
}

var tituloSecaoRE = regexp.MustCompile(`(?m)^(#{2,4})\s+([^\n]+)\n`)

// camadaDoAlvo resolve a camada da unidade que este artefato descreve — usada para achar
// o léxico de seções DAQUELA camada (`section_titles`). O alvo é o irmão sem o sufixo de
// peça derivada; sem alvo reconhecível, devolve vazio e vale o padrão do framework.
func camadaDoAlvo(root, outPath string, cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	rel, err := filepath.Rel(root, outPath)
	if err != nil {
		return ""
	}
	// A camada do ALVO primeiro, sempre. Classificar o próprio caminho da peça devolve a
	// camada genérica da peça (`spec`, `feature`) — que não tem léxico próprio e não é a
	// unidade que a spec descreve. Tentar o alvo depois nunca acontecia, porque a peça
	// sempre casa: a spec de um modelo era classificada como `spec`, não `schema-model`.
	base := strings.TrimSuffix(strings.TrimSuffix(rel, ".spec.md"), ".feature")
	if base != rel {
		exts := []string{".ts", ".tsx", ".go", ".py", ".js"}
		// O alvo que EXISTE no disco decide. Sem isto, a primeira extensão que casa
		// QUALQUER camada vence — e um catch-all casa tudo: a spec de uma tela
		// (`P.spec.md`, irmã de `P.tsx`) era classificada como `mobile-code`, porque
		// `P.ts` casava o catch-all e a busca parava ali. A camada errada traz o léxico
		// errado, e a spec nasce em dialeto que nenhuma vizinha usa.
		for _, ext := range exts {
			if _, err := os.Stat(filepath.Join(root, base+ext)); err != nil {
				continue
			}
			if l, _ := scan.Classify(base+ext, cfg); l != "" {
				return l
			}
		}
		// Nenhum alvo no disco (a spec nasce ANTES do código, que é o caso normal):
		// a extensão mais específica vence a mais genérica — `.tsx` antes de `.ts`,
		// porque quem tem camada própria costuma ser a específica.
		for _, ext := range []string{".tsx", ".ts", ".go", ".py", ".js"} {
			if l, _ := scan.Classify(base+ext, cfg); l != "" {
				return l
			}
		}
	}
	if l, _ := scan.Classify(rel, cfg); l != "" {
		return l
	}
	return ""
}
