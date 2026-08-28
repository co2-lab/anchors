package gate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// mock-carimbado: o dublê carrega a marca do trecho que ele substitui, e o gate
// RECALCULA essa marca contra o módulo real.
//
// É a camada agnóstica do par que ataca mock drift. O `mock-tipado` resolve o problema
// onde a linguagem ajuda (tipos estruturais fazem o compilador conferir a superfície do
// módulo), e só ali: a maior parte dos ecossistemas não tem `Partial<typeof X>`, e nem
// em TypeScript ele alcança a forma do VALOR devolvido — medido, 202 testes dublavam um
// hook do React Query com dois campos onde o real devolve ~15, e o tipo não distingue
// "parcial deliberado" de "defasado".
//
// O carimbo não interpreta o código: ele o LÊ. Hash de texto pega qualquer mudança —
// assinatura, corpo, tipo, constante — e funciona em Python, Ruby, Go ou JS puro, sem
// extrator por linguagem.
//
// O gate RECALCULA em vez de validar formato. Um carimbo que ninguém confronta é
// teatro: quem edita o teste o regeneraria para casar com o próprio mock, e ele passaria
// a certificar a si mesmo. Recalcular é o que torna o mecanismo à prova de quem escreve.
func checkMockCarimbado(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, "o carimbo vive no teste, ao lado do dublê que ele cobre"
	}
	detector, err := detectorDeDuble(cfg)
	if err != nil {
		// Regex inválido é erro de CONFIGURAÇÃO, não do arquivo sob análise — e precisa
		// falhar alto: silenciá-lo faria o gate varrer zero dublês e reportar verde.
		return Fail, fmt.Sprintf("`derived.mock_detect` não compila: %v", err)
	}
	if detector == nil {
		return Skip, "o projeto não declara `derived.mock_detect` — o Anchors não sabe " +
			"reconhecer um dublê neste ecossistema, e adivinhar reportaria verde sobre o que não conferiu"
	}

	carimbos := carimbosDeclarados(content)

	// AUSÊNCIA de carimbo é acusada, não pulada — e isto é a metade mais importante do
	// gate.
	//
	// Um carimbo divergente ACUSA; um carimbo ausente é SILÊNCIO, que é o mesmo "falha
	// aberto" que o `trinca-completa` existe para fechar. Se a ausência passasse, o
	// carimbo viraria opcional na prática: quem não põe nunca é cobrado, e o mecanismo
	// protegeria apenas quem já escolheu ser protegido — isso é convenção, não gate.
	//
	// A tentação de pular vem de olhar projeto LEGADO, onde ligar a cobrança produz
	// centenas de achados de uma vez. Mas desenhar por esse caso transformaria um
	// problema de migração numa propriedade permanente do framework: todo projeto
	// futuro herdaria a frouxidão para acomodar um projeto antigo. Num projeto que
	// NASCE com o Anchors não há dívida — o primeiro dublê é escrito depois do gate
	// existir, e cobrar custa um carimbo.
	//
	// O legado se resolve com o vocabulário que já existe: `blocking: false` durante a
	// adoção, e opt-out por unidade com razão escrita onde a decisão for deliberada.
	var semCarimbo []string
	for _, modulo := range dublesDetectados(content, detector) {
		if !ehModuloRegido(modulo, g) {
			continue // biblioteca de terceiro não é cobrada (ver `ehModuloRegido`)
		}
		if !moduloTemCarimbo(modulo, carimbos) {
			semCarimbo = append(semCarimbo, modulo)
		}
	}
	if len(semCarimbo) > 0 {
		return Fail, fmt.Sprintf(
			"%d dublê(s) sem carimbo de contrato: %s.\n\nSem carimbo não há o que "+
				"recalcular — o dublê pode reproduzir um contrato que já mudou e nada "+
				"acusa. Carimbe o trecho que ele substitui: "+
				"`// @contract: <arquivo> | <linha que abre o trecho> | <qtd de linhas> | <hash>`",
			len(semCarimbo), strings.Join(semCarimbo, ", "))
	}
	if len(carimbos) == 0 {
		return Skip, "o teste não dubla módulo regido — não há carimbo a cobrar"
	}

	var divergentes []string
	for _, c := range carimbos {
		atual, err := recalculaCarimbo(root, c)
		if err != nil {
			divergentes = append(divergentes, fmt.Sprintf("%s (%v)", c.ancora, err))
			continue
		}
		if atual != c.hash {
			divergentes = append(divergentes, fmt.Sprintf(
				"`%s` — carimbo %s, contrato hoje %s", c.ancora, c.hash, atual))
		}
	}
	if len(divergentes) == 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(
		"%d carimbo(s) não batem com o módulo real:\n  - %s\n\n"+
			"O trecho mudou depois que este dublê foi escrito — ele agora reproduz um "+
			"contrato que não existe mais, e o teste passa sobre a versão antiga. "+
			"Confira o que mudou, ajuste o dublê e atualize o carimbo (ou registre a "+
			"divergência deliberada com uma razão escrita ao lado)",
		len(divergentes), strings.Join(divergentes, "\n  - "))
}

// moduloTemCarimbo liga o dublê ao carimbo pelo CAMINHO do arquivo carimbado.
//
// O especificador do dublê (`@/src/hooks/useX`) e o caminho do carimbo
// (`apps/mobile/src/hooks/useX.ts`) descrevem o mesmo arquivo por vias diferentes —
// alias e caminho de disco. Casar pelo sufixo sem extensão resolve os dois sem
// precisar de um resolvedor de alias, que seria específico do ecossistema.
func moduloTemCarimbo(modulo string, carimbos []carimboDeclarado) bool {
	alvo := semExtensao(strings.TrimPrefix(modulo, "./"))
	for strings.HasPrefix(alvo, "../") {
		alvo = strings.TrimPrefix(alvo, "../")
	}
	if strings.HasPrefix(alvo, "@") {
		if i := strings.Index(alvo, "/"); i >= 0 {
			alvo = alvo[i+1:]
		}
	}
	for _, c := range carimbos {
		arq := semExtensao(c.arquivo)
		if arq == alvo || strings.HasSuffix(arq, "/"+alvo) {
			return true
		}
	}
	return false
}

// carimboDeclarado — o que o teste afirma sobre o trecho que dubla.
type carimboDeclarado struct {
	arquivo string // caminho do módulo, relativo à raiz
	ancora  string // a LINHA INTEIRA que abre o trecho (conteúdo, nunca número)
	qtd     int    // quantas linhas a partir da âncora entram no hash
	hash    string // o hash gravado
}

// carimboRE casa a anotação:
//
//	// @contract: caminho/do/modulo.ts | export function useX( | 10 | 361280fb
//
// O separador é `|` porque a âncora é uma linha de código e pode conter vírgula, dois
// pontos e parênteses — qualquer separador mais comum a partiria no meio.
var carimboRE = regexp.MustCompile(
	`@contract:\s*([^|\n]+?)\s*\|\s*(.+?)\s*\|\s*(\d+)\s*\|\s*([0-9a-f]+)`)

func carimbosDeclarados(content string) []carimboDeclarado {
	var out []carimboDeclarado
	for _, m := range carimboRE.FindAllStringSubmatch(content, -1) {
		qtd, err := strconv.Atoi(m[3])
		if err != nil || qtd <= 0 {
			continue
		}
		out = append(out, carimboDeclarado{
			arquivo: m[1], ancora: m[2], qtd: qtd, hash: m[4],
		})
	}
	return out
}

// recalculaCarimbo lê o módulo real e devolve o hash do trecho HOJE.
//
// A âncora é procurada por CONTEÚDO — é o que torna o carimbo imune a deslocamento.
// Duas ocorrências da mesma linha tornam o alvo ambíguo, e o gate prefere acusar a
// escolher uma: um carimbo que aponta para "alguma das duas" não prova nada.
func recalculaCarimbo(root string, c carimboDeclarado) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, c.arquivo))
	if err != nil {
		return "", fmt.Errorf("módulo não encontrado: %s", c.arquivo)
	}
	linhas := strings.Split(string(b), "\n")

	idx := -1
	ocorrencias := 0
	for i, l := range linhas {
		if strings.TrimRight(l, "\r") == c.ancora {
			ocorrencias++
			if idx < 0 {
				idx = i
			}
		}
	}
	switch {
	case ocorrencias == 0:
		// A âncora sumiu: renome, remoção ou reescrita da linha. É achado, não erro de
		// ferramenta — e falha EXPLÍCITA é melhor que silêncio, porque o dublê
		// certamente está desatualizado.
		return "", fmt.Errorf("âncora não encontrada — o trecho foi renomeado ou removido")
	case ocorrencias > 1:
		return "", fmt.Errorf("âncora ambígua (%d ocorrências) — use uma linha única do trecho", ocorrencias)
	}

	fim := idx + c.qtd
	if fim > len(linhas) {
		fim = len(linhas)
	}
	return hashDoTrecho(strings.Join(linhas[idx:fim], "\n")), nil
}

// hashDoTrecho — sha256 truncado em 8 hex. Truncado porque o carimbo mora numa linha de
// comentário e é lido por humano; 32 bits bastam para detectar mudança acidental, que é
// o que este gate persegue (não há adversário forjando colisão contra o próprio teste).
func hashDoTrecho(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])[:8]
}

// detectorDeDuble compila o regex que ESTE projeto usa para escrever um dublê.
//
// Devolve (nil, nil) quando o projeto não declara — e é o que desliga o gate. A
// alternativa (embutir o padrão jest/vitest como default) faria o gate rodar num
// projeto Python, casar zero dublês e reportar VERDE sobre o que não conferiu.
func detectorDeDuble(cfg *config.Config) (*regexp.Regexp, error) {
	if cfg == nil || cfg.Derived == nil {
		return nil, nil
	}
	p := strings.TrimSpace(cfg.Derived.MockDetect)
	if p == "" {
		return nil, nil
	}
	re, err := regexp.Compile(p)
	if err != nil {
		return nil, err
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("o padrão precisa de um grupo de captura com o módulo dublado")
	}
	return re, nil
}

// dublesDetectados aplica o padrão do projeto e devolve os módulos dublados.
func dublesDetectados(content string, re *regexp.Regexp) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			out = append(out, m[1])
		}
	}
	return out
}
