package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

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
		Use:   "code <nome>",
		Short: "Gera um código de identidade único para uma nova unidade",
		Long: `Gera o código de cenário (a identidade, TRACEABILITY §3) de uma unidade,
GARANTINDO unicidade contra os códigos já em uso no mapa. Use ANTES de escrever a
spec, para não colidir com outra unidade (o que geraria propagação cruzada errada).

  anchors code Spacer            → sugere um código livre para "Spacer"
  anchors code --check SPCR      → diz se SPCR está livre ou já é de alguém
  anchors code list              → lista TODOS os códigos em uso, com quem os usa
  anchors code list --in apps/mobile     → só os do workspace

O "code list" lê o campo de identidade do MAPA, não procura por padrão no texto: um grep de
ABCDX-X01 casa menção em prosa, comentário e nome de arquivo, e depende do comprimento do
código (que varia por projeto, ver code_lengths). A saída é "código<TAB>onde", uma por
linha; o resumo vai para stderr, então "anchors code list | ..." encadeia sem sujeira.

O algoritmo (Camada 2 do SPEC_GUIDE: compressão do nome; Camada 3: resolução de
colisão) é agnóstico. Prefixos de módulo (Camada 1) são dialeto de projeto — se o seu
usa, escolha o código à mão e valide com --check.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			taken, owners, err := takenCodes(mapPath)
			if err != nil {
				return fmt.Errorf("ler o mapa: %w (rode `anchors map build`)", err)
			}

			// modo --check: valida um código proposto
			if check != "" {
				c := strings.ToUpper(check)
				if who, ok := owners[c]; ok {
					fmt.Printf("✗ %s já é usado por: %s\n", c, strings.Join(who, ", "))
					fmt.Println("  escolha outro (ou rode `anchors code <nome>` para uma sugestão livre)")
					return errCollision
				}
				fmt.Printf("✓ %s está livre\n", c)
				return nil
			}

			// modo gerar: precisa do nome (ou caminho)
			if len(args) != 1 {
				return fmt.Errorf("informe o nome da unidade (ex: `anchors code Spacer`) ou use --check <código>")
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
				rel := relTo(absRoot, args[0])
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

			fmt.Printf("✓ código livre: %s\n", unique)
			if prefix != "" {
				fmt.Printf("  (prefixo de módulo '%s' da Estrutura + distintivo de '%s')\n", prefix, name)
			}
			if unique != canonical {
				if owner, ok := owners[canonical]; ok {
					fmt.Printf("  (o canônico %s já é de %s — ajustado para um código livre)\n",
						canonical, strings.Join(owner, ", "))
				} else {
					// basename genérico (handler/index/…): a identidade veio do dir-pai.
					fmt.Printf("  (basename genérico — a identidade veio do dir-pai, não do arquivo)\n")
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().StringVar(&check, "check", "", "valida se um código proposto está livre")
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

func (*collisionError) Error() string { return "código já em uso" }

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
		Short: "Lista os códigos de identidade EM USO no projeto (do mapa, não por regex)",
		Long: `Enumera os códigos que o mapa conhece, com a pasta de cada unidade.

Lê o campo de identidade do MAPA — não procura padrão no texto. Um grep por ABCDX-X01 casa
menção em prosa, comentário e nome de arquivo, e depende do comprimento do código, que
varia por projeto (ver code_lengths no anchors.yaml). Aqui a resposta é exata.

  anchors code list                      todos os códigos em uso
  anchors code list --in apps/mobile     só os de um workspace do monorepo

Saída: "código<TAB>onde", uma por linha; o resumo vai para stderr, então a lista encadeia
em pipe sem sujeira. Um código com VÁRIAS pastas é colisão de identidade — o doctor a
reporta, e aqui ela fica visível de graça.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
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
				return fmt.Errorf("carregar %s: %w", config.DefaultFile, cerr)
			}
			_, owners, err := takenCodes(mapPath)
			if err != nil {
				return fmt.Errorf("ler o mapa: %w (rode `anchors map build`)", err)
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
					fmt.Printf("✓ %d código(s) declarado(s) com comprimento conforme (%s)\n",
						ok, joinLens(config.CodeLengths))
					if citados > 0 {
						fmt.Printf("  (%d código(s) apenas CITADO(s) — fixture de teste ou exemplo em doc — não conferidos)\n", citados)
					}
					return nil
				}
				fmt.Printf("%d código(s) fora do comprimento declarado (%s):\n\n",
					len(divs), joinLens(config.CodeLengths))
				for _, d := range divs {
					fmt.Printf("  %s → %s\t%s\n", d.atual, d.correto, strings.Join(d.onde, ", "))
				}
				fmt.Printf("\n  %d conforme(s). O proposto é o código CANÔNICO (%d chars) — o que o\n", ok, alvo)
				fmt.Println("  algoritmo geraria para o nome da unidade, com colisão resolvida contra os")
				fmt.Println("  códigos JÁ conformes (os errados não reservam lugar: eles vão mudar).")
				fmt.Println("  Aplique com `anchors code list --check --fix` (usa `anchors recode`, que propaga).")
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
					fmt.Printf("nenhum código em uso sob %q\n", filtro)
					return nil
				}
				fmt.Println("nenhum código em uso — o mapa não tem nó com identidade (rode `anchors map build`)")
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
				type saida struct {
					Code string   `json:"code"`
					Onde []string `json:"onde"`
				}
				out := make([]saida, 0, len(linhas))
				for _, l := range linhas {
					out = append(out, saida{Code: l.code, Onde: l.onde})
				}
				b, jerr := json.MarshalIndent(out, "", "  ")
				if jerr != nil {
					return fmt.Errorf("serializar em JSON: %w", jerr)
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
			fmt.Fprintf(os.Stderr, "\n%d código(s) em uso%s\n", len(linhas),
				map[bool]string{true: " sob " + filtro, false: ""}[filtro != ""])
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().StringVar(&filtro, "in", "", "só os códigos sob este prefixo de caminho (ex.: apps/mobile)")
	cmd.Flags().BoolVar(&conferir, "check", false, "confere o COMPRIMENTO de cada código contra o code_lengths e mostra o conserto")
	cmd.Flags().BoolVar(&corrigir, "fix", false, "com --check: aplica o conserto via `anchors recode` (propaga por spec/feature/teste/mapa)")
	cmd.Flags().BoolVar(&emJSON, "json", false, "emite a lista em JSON (para consumo por script/pipeline)")
	return cmd
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
		return "não declarado"
	}
	return strings.Join(partes, " ou ")
}
