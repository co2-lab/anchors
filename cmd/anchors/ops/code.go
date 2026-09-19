package ops

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// `anchors code` fecha o furo #7 do exercício C: a IA precisava de um código de
// identidade ÚNICO antes de escrever a spec, e o CLI não ajudava — daí a colisão
// SPCR. Agora o CLI é o autor+juiz da identidade, lendo os códigos já em uso no mapa.
func newCodeCmd() *cobra.Command {
	var root, mapPath, check string
	cmd := &cobra.Command{
		Use:   "code <name>",
		Short: "Generate a unique identity code for a new unit",
		Long: `Generates the scenario code (the identity, TRACEABILITY §3) of a unit,
GUARANTEEING uniqueness against the codes already in use in the map. Use it BEFORE writing the
spec, so as not to collide with another unit (which would generate wrong cross propagation).

  anchors code Spacer            → suggests a free code for "Spacer"
  anchors code --check SPCR      → says whether SPCR is free or already belongs to someone
  anchors code list              → lists ALL the codes in use, with who uses them
  anchors code list --in apps/mobile     → only those of the workspace

The "code list" reads the identity field of the MAP, it does not search for a pattern in the text: a grep for
ABCDX-X01 matches a mention in prose, a comment and a file name, and depends on the length of the
code (which varies per project, see code_lengths). The output is "code<TAB>where", one per
line; the summary goes to stderr, so "anchors code list | ..." pipes without dirt.

The algorithm (Layer 2 of the SPEC_GUIDE: name compression; Layer 3: collision
resolution) is agnostic. Module prefixes (Layer 1) are project dialect — if yours
uses them, choose the code by hand and validate it with --check.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			taken, owners, err := takenCodes(mapPath)
			if err != nil {
				return fmt.Errorf("read the map: %w (run `anchors map build`)", err)
			}

			// modo --check: valida um código proposto
			if check != "" {
				c := strings.ToUpper(check)
				if who, ok := owners[c]; ok {
					fmt.Printf("✗ %s is already used by: %s\n", c, strings.Join(who, ", "))
					fmt.Println("  choose another (or run `anchors code <name>` for a free suggestion)")
					return errCollision
				}
				fmt.Printf("✓ %s is free\n", c)
				return nil
			}

			// modo gerar: precisa do nome (ou caminho)
			if len(args) != 1 {
				return fmt.Errorf("provide the unit name (e.g.: `anchors code Spacer`) or use --check <code>")
			}
			// Ligação com a ESTRUTURA. Precedência:
			//  1. layer com code_prefix declarado (Camada 1) → usa esse prefixo curado.
			//  2. senão, caminho → GenerateFromPath: usa o dir-pai quando o basename é
			//     genérico (handler/index) ou quando há colisão (desempate semântico).
			//  3. nome puro (sem '/') → algoritmo canônico (Camada 2).
			name := unitName(args[0])
			prefix := ""
			// A config é lida SEMPRE, não só quando o argumento é um caminho.
			//
			// Ela carrega o `code_lengths`, e é o `Load` que ajusta o comprimento gerado (o
			// gancho no pacote `code`). Lendo-a apenas para resolver o prefixo de camada, um
			// `anchors code AlertsScreen` (nome puro) gerava no default do framework: 5
			// slots num projeto de 4, sugerindo identidade que o próprio engine do projeto
			// não reconhece.
			cfg, cerr := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if cerr == nil && strings.ContainsRune(args[0], '/') {
				rel := common.RelTo(absRoot, args[0])
				if layer, _ := scan.Classify(rel, cfg); layer != "" {
					if l, ok := cfg.Layers[layer]; ok {
						prefix = l.CodePrefix
					}
				}
			}

			var canonical, unique string
			switch {
			case prefix != "":
				canonical = code.GenerateWithPrefix(name, prefix)
				unique = code.GenerateUniqueWithPrefix(name, prefix, taken)
			case strings.ContainsRune(args[0], '/'):
				canonical = code.Generate(name)
				unique = code.GenerateFromPath(args[0], taken)
			default:
				canonical = code.Generate(name)
				unique = code.GenerateUnique(name, taken)
			}

			fmt.Printf("✓ free code: %s\n", unique)
			if prefix != "" {
				fmt.Printf("  (module prefix '%s' from the Structure + distinctive of '%s')\n", prefix, name)
			}
			if unique != canonical {
				if owner, ok := owners[canonical]; ok {
					fmt.Printf("  (the canonical %s already belongs to %s — adjusted to a free code)\n",
						canonical, strings.Join(owner, ", "))
				} else {
					// basename genérico (handler/index/…): a identidade veio do dir-pai.
					fmt.Printf("  (generic basename — the identity came from the parent dir, not from the file)\n")
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringVar(&check, "check", "", "validate whether a proposed code is free")
	cmd.AddCommand(newCodeListCmd())
	return cmd
}

// unitName extrai o nome da unidade de um caminho: o basename sem extensão nem
// sufixos de artefato (Login.spec.md → Login; auth/screens/NewLogin.tsx → NewLogin).
func unitName(path string) string {
	base := filepath.Base(path)
	for _, suf := range []string{".spec.md", ".feature"} {
		if b, ok := strings.CutSuffix(base, suf); ok {
			return b
		}
	}
	if b, _, ok := strings.Cut(base, ".test."); ok {
		return b
	}
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// errCollision faz o --check sair com código ≠ 0 quando o código proposto colide —
// útil para a IA/script ramificar sem parsear texto.
var errCollision = &collisionError{}

type collisionError struct{}

func (*collisionError) Error() string { return "code already in use" }

// newCodeListCmd — `anchors code list`, subcomando e não flag.
//
// Segue a forma que o CLI já usa (`anchors map build` / `map show`): quando o verbo muda o
// que o comando FAZ (gerar um código vs. enumerar os existentes), subcomando é mais honesto
// que flag — o help de cada um fica separado e as flags de um não poluem o outro.
func newCodeListCmd() *cobra.Command {
	var root, mapPath, filtro string
	var conferir, corrigir, emJSON bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the identity codes IN USE in the project (from the map, not by regex)",
		Long: `Enumerates the codes the map knows, with each unit's folder.

Reads the identity field of the MAP — it does not search for a pattern in the text. A grep for ABCDX-X01 matches
a mention in prose, a comment and a file name, and depends on the length of the code, which
varies per project (see code_lengths in anchors.yaml). Here the answer is exact.

  anchors code list                      every code in use
  anchors code list --in apps/mobile     only those of one workspace of the monorepo

Output: "code<TAB>where", one per line; the summary goes to stderr, so the list pipes
without dirt. A code with SEVERAL folders is an identity collision — the doctor
reports it, and here it becomes visible for free.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			// A config é lida ANTES de qualquer decisão sobre comprimento: é ela que carrega
			// o `code_lengths` (e, pelo gancho, ajusta o comprimento que a geração emite).
			// Sem esta linha o comando responde no default do FRAMEWORK — reportou "fora do
			// comprimento declarado (5)" num projeto que declara 4, marcando os 601 códigos
			// dele como errados. É o mesmo erro que o comando pai já tinha e que consertei
			// uma hora antes: ler a config condicionalmente é ler no default.
			if _, cerr := config.Load(filepath.Join(absRoot, config.DefaultFile)); cerr != nil {
				return fmt.Errorf("load %s: %w", config.DefaultFile, cerr)
			}
			_, owners, err := takenCodes(mapPath)
			if err != nil {
				return fmt.Errorf("read the map: %w (run `anchors map build`)", err)
			}

			// --check: confere o COMPRIMENTO de cada código contra o `code_lengths` do
			// projeto, e propõe o conserto.
			//
			// Deliberadamente NÃO regenera do nome do arquivo. Testado no app de referência: com o
			// comprimento certo, zero de cinco códigos batiam com o gerado — e nenhum dos
			// declarados estava errado. `ARSC` para AlertsScreen é escolha humana melhor que
			// o `LRTS` do algoritmo, e a própria doutrina diz que o gerado é SUGESTÃO ("o
			// importante é a unicidade no namespace global, não a estética"). Um check por
			// regeneração marcaria os 659 códigos do app de referência como errados.
			//
			// Comprimento é diferente: é INVARIANTE, não estética. Código fora do
			// `code_lengths` não é reconhecido pelo engine — a trinca fica invisível e os
			// gates de identidade não têm o que confrontar, sem nada acusar a causa. Foi o
			// que custou uma investigação hoje (o mapa do app de referência tinha 8 códigos para 275
			// specs) e o que deixou 105 testes vermelhos no Anchors.
			if conferir {
				aceitos := map[int]bool{}
				for _, l := range config.CodeLengths {
					aceitos[l] = true
				}
				// Os códigos ERRADOS não ocupam espaço na resolução de colisão: eles VÃO
				// mudar, então reservar o lugar deles congelaria o erro. Sem esta regra, o
				// primeiro código fora de comprimento a ocupar `LOGI` desviaria o correto
				// para `LOGJ` — propagando o defeito em vez de corrigi-lo. (Regra do Adriel.)
				semOsErrados := map[string]bool{}
				for c := range owners {
					if aceitos[len([]rune(c))] {
						semOsErrados[c] = true
					}
				}
				alvo := code.Slots // o comprimento que a geração emite: o conserto vai para ele
				type divergencia struct {
					atual, correto string
					onde           []string
				}
				var divs []divergencia
				var ok, citados int
				for c, donos := range owners {
					if filtro != "" {
						var m []string
						for _, d := range donos {
							if strings.HasPrefix(d, filtro) {
								m = append(m, d)
							}
						}
						if len(m) == 0 {
							continue
						}
						donos = m
					}
					if aceitos[len([]rune(c))] {
						ok++
						continue
					}
					// Sem identidade DECLARADA não há o que conferir: o código veio de uma
					// citação (fixture de teste, exemplo em doc), e citação não promete
					// identidade. Medido no Anchors: sem este filtro, 29 falsos positivos e
					// zero achados — o projeto tem 0 specs, e todo "código em uso" vinha de
					// string de teste.
					if !codeDecl[c] {
						citados++
						continue
					}
					sort.Strings(donos)
					// O conserto é o código CANÔNICO — o que o algoritmo geraria para o nome
					// da unidade — e não padding com X.
					//
					// O padding parecia conservador (`MTVR` → `MTVRX` preserva o prefixo
					// escolhido), mas produz código que o próprio `anchors code` não geraria:
					// o algoritmo extrai mais uma LETRA DO NOME, e o X é só último recurso
					// para nome curto. `AuditDetailScreen` canônico é `ADDTD`, não `ADDTX`.
					// Deixar o X seria plantar divergência para o próximo check acusar.
					forem := unitName(codeFile[c])
					canonico := code.GenerateUnique(forem, semOsErrados)
					divs = append(divs, divergencia{atual: c, correto: canonico, onde: donos})
				}
				sort.Slice(divs, func(i, j int) bool { return divs[i].atual < divs[j].atual })
				if len(divs) == 0 {
					fmt.Printf("✓ %d declared code(s) with conforming length (%s)\n",
						ok, joinLens(config.CodeLengths))
					if citados > 0 {
						fmt.Printf("  (%d code(s) only CITED — test fixture or doc example — not checked)\n", citados)
					}
					return nil
				}
				fmt.Printf("%d code(s) outside the declared length (%s):\n\n",
					len(divs), joinLens(config.CodeLengths))
				for _, d := range divs {
					fmt.Printf("  %s → %s\t%s\n", d.atual, d.correto, strings.Join(d.onde, ", "))
				}
				fmt.Printf("\n  %d conforming. The proposal is the CANONICAL code (%d chars) — the one the\n", ok, alvo)
				fmt.Println("  algorithm would generate for the unit name, with collision resolved against the")
				fmt.Println("  codes that ALREADY conform (the wrong ones reserve no place: they will change).")
				fmt.Println("  Apply with `anchors code list --check --fix` (uses `anchors recode`, which propagates).")
				return errCollision
			}
			// modo --list: enumera os códigos EM USO, com quem os usa.
			//
			// Existe porque a alternativa que as pessoas usam é grep por padrão de código, e
			// ela traz mais do que códigos: casa menção em prosa, comentário, nome de arquivo
			// e qualquer string com a forma `ABCDX-X01`. O mapa já tem o campo estruturado
			// (`Node.Code`, preenchido pelo scan a partir do header/identidade), então
			// enumerar dali é exato — sem falso positivo e sem depender do comprimento do
			// código, que varia por projeto (`code_lengths`).
			//
			// `--in <prefixo>` filtra por caminho, para monorepo: `--in apps/mobile` responde
			// "os códigos deste workspace", que é a pergunta real de quem trabalha num deles.
			type linha struct {
				code string
				onde []string
			}
			var linhas []linha
			for c, donos := range owners {
				if filtro != "" {
					var mantidos []string
					for _, d := range donos {
						if strings.HasPrefix(d, filtro) {
							mantidos = append(mantidos, d)
						}
					}
					if len(mantidos) == 0 {
						continue
					}
					donos = mantidos
				}
				sort.Strings(donos)
				linhas = append(linhas, linha{code: c, onde: donos})
			}
			sort.Slice(linhas, func(i, j int) bool { return linhas[i].code < linhas[j].code })
			if len(linhas) == 0 {
				// Em JSON, "nenhum" é uma lista vazia — não uma frase. Um consumidor que
				// recebesse prosa aqui quebraria no parse justamente no projeto novo, que é
				// onde ele mais roda.
				if emJSON {
					fmt.Println("[]")
					return nil
				}
				if filtro != "" {
					fmt.Printf("no code in use under %q\n", filtro)
					return nil
				}
				fmt.Println("no code in use — the map has no node with identity (run `anchors map build`)")
				return nil
			}
			// --json: a mesma lista, para quem CONSOME em vez de ler. O pipeline de
			// identificação (`.github/workflows/anchors-identify.yml`) precisa do código e
			// das pastas separados — parsear o TSV exigiria dividir o campo `onde` por ", ",
			// que se confunde com uma vírgula dentro de um caminho.
			//
			// Vai para stdout inteiro, sem o resumo em stderr: um consumidor quer um
			// documento JSON válido, não um fluxo com anexo.
			if emJSON {
				// `arquivo`, `kind` e `titulo` existem porque quem consome isto precisa
				// NOMEAR o trabalho — e `onde` é a PASTA da unidade, não o arquivo: 16
				// planos em `plans/` têm o mesmo `onde`. Sem estes campos, o consumidor
				// teria de parsear o grafo por conta própria (medido: um pipeline tentou,
				// com `grep -B1`, e falhou porque o `code:` está 3 linhas abaixo do `id:`).
				type saida struct {
					Code    string   `json:"code"`
					Onde    []string `json:"onde"`
					Arquivo string   `json:"arquivo,omitempty"`
					Kind    string   `json:"kind,omitempty"`
					Titulo  string   `json:"titulo,omitempty"`
					// Needs é a ORDEM DE TRABALHO: as fases de plano que precisam fechar
					// antes desta unidade poder ser trabalhada.
					//
					// Sai aqui porque é o pipeline de claim que precisa dela, e a
					// alternativa seria ele fazer grep no header de cada markdown — o
					// mesmo caminho que já falhou antes com `grep -B1` e que motivou os
					// campos acima.
					Needs []string `json:"needs,omitempty"`
					// Parent é o PERTENCIMENTO — quem contém este artefato. Sai junto com
					// `needs` porque quem monta uma árvore precisa das duas: `parent` dá
					// a estrutura, `needs` dá a ordem dentro dela.
					Parent string `json:"parent,omitempty"`
					// Revises: os planos que ESTE revisa. Sai aqui porque o pipeline de
					// identificação precisa saber quais planos foram revisados para avisar
					// os cards que a revisão atinge — e ele não lê o mapa por conta própria.
					Revises []string `json:"revises,omitempty"`
				}
				kinds := kindByFile(mapPath)
				needs := needsByFile(mapPath)
				parents := parentByFile(mapPath)
				revs := revisesByFile(mapPath)
				out := make([]saida, 0, len(linhas))
				for _, l := range linhas {
					arq := codeFile[l.code]
					out = append(out, saida{
						Code:    l.code,
						Onde:    l.onde,
						Arquivo: arq,
						Kind:    kinds[arq],
						Titulo:  fileTitle(absRoot, arq),
						Needs:   needs[arq],
						Parent:  parents[arq],
						Revises: revs[arq],
					})
				}
				b, jerr := json.MarshalIndent(out, "", "  ")
				if jerr != nil {
					return fmt.Errorf("serialize to JSON: %w", jerr)
				}
				fmt.Println(string(b))
				return nil
			}
			for _, l := range linhas {
				// Uma unidade por linha, e o código PRIMEIRO: quem lê está procurando um
				// código, não um caminho. Várias pastas no mesmo código é colisão de
				// identidade — o `doctor` a reporta, e aqui ela fica visível de graça.
				fmt.Printf("%s\t%s\n", l.code, strings.Join(l.onde, ", "))
			}
			fmt.Fprintf(os.Stderr, "\n%d code(s) in use%s\n", len(linhas),
				map[bool]string{true: " under " + filtro, false: ""}[filtro != ""])
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringVar(&filtro, "in", "", "only the codes under this path prefix (e.g.: apps/mobile)")
	cmd.Flags().BoolVar(&conferir, "check", false, "check the LENGTH of each code against code_lengths and show the fix")
	cmd.Flags().BoolVar(&corrigir, "fix", false, "with --check: apply the fix via `anchors recode` (propagates through spec/feature/test/map)")
	cmd.Flags().BoolVar(&emJSON, "json", false, "emit the list as JSON (for script/pipeline consumption)")
	return cmd
}

// kindByFile devolve o kind de cada nó do mapa. Serve ao `--json`: quem consome
// precisa dizer que TIPO de artefato o trabalho é ("Implementar plan — Fundação").
func kindByFile(mapPath string) map[string]string {
	out := map[string]string{}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return out
	}
	for _, n := range g.Nodes {
		out[n.ID] = string(n.Kind)
	}
	return out
}

// fileTitle lê o primeiro `# título` de um markdown. É o nome que o AUTOR deu ao
// artefato — melhor do que qualquer coisa derivada do caminho. Vazio para arquivo que
// não é markdown ou não tem título, e aí quem consome decide o que fazer.
func fileTitle(root, rel string) string {
	if rel == "" || !strings.HasSuffix(rel, ".md") {
		return ""
	}
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		return ""
	}
	for _, linha := range strings.Split(string(b), "\n") {
		t, ok := strings.CutPrefix(strings.TrimSpace(linha), "# ")
		if !ok {
			continue
		}
		// Um título costuma repetir o tipo e o número que já estão em outro lugar
		// ("Plano 0001 — Fundação"). Quem consome já tem o kind e o código, e a
		// repetição só faz o texto crescer: sobra a parte que de fato nomeia.
		if _, depois, achou := strings.Cut(t, "—"); achou {
			t = depois
		}
		return strings.TrimSpace(t)
	}
	return ""
}

// takenCodes lê o mapa e devolve o conjunto de códigos em uso + quem os usa (por
// unidade), para sugestão e para --check.
// codeFile guarda, por código, o arquivo que melhor representa a unidade (a spec quando
// existe). Preenchido por takenCodes; lido pelo --check para regenerar o código canônico.
var codeFile = map[string]string{}

// codeDecl marca os códigos cuja identidade é DECLARADA — os únicos que o --check confere.
var codeDecl = map[string]bool{}

func takenCodes(mapPath string) (taken map[string]bool, owners map[string][]string, err error) {
	g, err := mapx.Load(mapPath)
	if err != nil {
		return nil, nil, err
	}
	taken = map[string]bool{}
	owners = map[string][]string{}
	seen := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Code == "" {
			continue
		}
		taken[n.Code] = true
		// `path.Dir`, e não `filepath.Dir`: o ID do nó é sempre em `/`, e no Windows o
		// `filepath` devolveria `packages\backend`. Como o `--in packages/backend` compara
		// por PREFIXO DE TEXTO, a forma nativa fazia o filtro não casar nada — o comando
		// respondia "nenhum código em uso" para um workspace cheio deles.
		dir := path.Dir(n.ID)
		key := n.Code + "\x00" + dir
		if !seen[key] {
			seen[key] = true
			owners[n.Code] = append(owners[n.Code], dir)
		}
		// O ARQUIVO-DONO da identidade, para quem precisa regenerar o código a partir do
		// nome da unidade. A spec vence: ela é a dona (TRACEABILITY §3), e o nome dela é o
		// da unidade. Sem spec, fica o primeiro arquivo visto — melhor que nada, e o
		// chamador sabe distinguir pelo sufixo.
		if _, já := codeFile[n.Code]; !já || n.Kind == mapx.KindSpec {
			codeFile[n.Code] = n.ID
		}
		// Identidade DECLARADA (header `code:`) vs. inferida do texto. Só a declarada é
		// conferível: um fixture que usa `AAAAX-B01` como exemplo não PROMETE ser a unidade
		// AAAA. Ver Node.CodeDeclarado.
		if n.CodeDeclarado {
			codeDecl[n.Code] = true
		}
	}
	return taken, owners, nil
}

// joinLens formata os comprimentos aceitos para a mensagem ("4", "4 ou 5").
func joinLens(ls []int) string {
	partes := make([]string, len(ls))
	for i, l := range ls {
		partes[i] = fmt.Sprint(l)
	}
	if len(partes) == 0 {
		return "not declared"
	}
	return strings.Join(partes, " or ")
}

// needsByFile devolve a ordem de trabalho declarada por cada arquivo do mapa — as
// fases de plano que precisam fechar antes dele.
func needsByFile(mapPath string) map[string][]string {
	out := map[string][]string{}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return out
	}
	for _, n := range g.Nodes {
		if len(n.Needs) > 0 {
			out[n.ID] = n.Needs
		}
	}
	return out
}

// parentByFile devolve o pertencimento declarado por cada arquivo do mapa — o código
// de quem o contém.
func parentByFile(mapPath string) map[string]string {
	out := map[string]string{}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return out
	}
	for _, n := range g.Nodes {
		if n.Parent != "" {
			out[n.ID] = n.Parent
		}
	}
	return out
}

// revisesByFile devolve os planos que cada arquivo do mapa revisa.
func revisesByFile(mapPath string) map[string][]string {
	out := map[string][]string{}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return out
	}
	for _, n := range g.Nodes {
		if len(n.Revises) > 0 {
			out[n.ID] = n.Revises
		}
	}
	return out
}
