package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkIdentityConsistent — o CÓDIGO da spec é a identidade da unidade. Este gate
// confronta essa identidade contra as outras superfícies onde ela reaparece: o
// testID exposto pelo código e o nome do baseline de regressão visual.
//
// POR QUE ELE EXISTE. Os demais gates de identidade verificam PRESENÇA (a spec tem
// código? o cenário tem código?), nunca CONCORDÂNCIA entre superfícies. Uma unidade
// pode declarar `code: BGET`, expor `testID=":bdge-screen"` e guardar o baseline como
// `BDEDX-VR-*.png` — três identidades para a mesma coisa — e o pipeline fica verde,
// porque cada gate olha a sua superfície isoladamente.
//
// O CUSTO REAL não é estético. O dicionário de códigos que o spellcheck consome é
// GERADO do mapa; uma sigla que só existe no testID não está lá, o corretor a acusa,
// e a correção natural (dicionarizar a sigla) CRISTALIZA a divergência — o sintoma
// vira vocabulário e a causa some. Foi assim que `bdgc`/`bdge` entraram no dicionário
// do app de referência, e é o que motivou este gate.
//
// O QUE NÃO É DIVERGÊNCIA, e este é o discernimento que o gate precisa ter: um
// componente pode legitimamente usar o código de OUTRA unidade como prefixo, quando o
// testID nomeia ONDE o elemento aparece. `SpendingMonthCard` (SMCD) expõe
// `:home-spending-card` porque vive dentro da `HomeScreen` (HOME) — e é assim que o
// teste de ponta a ponta o localiza, partindo da tela. Cobrar unificação aí seria
// impor uma convenção que quebraria a navegação dos flows.
//
// A regra que separa os dois casos: um prefixo que é código de ALGUMA unidade do mapa
// é identidade conhecida (reuso deliberado); um prefixo com forma de código que não
// pertence a ninguém é identidade ÓRFÃ — não existe no mapa, logo não existe no
// dicionário gerado, logo é exatamente o caso que empurra siglas para o vocabulário
// manual. Só esse o gate cobra.
func checkIdentityConsistent(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a identidade é declarada na spec — é dela que o confronto parte"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}
	code := strings.ToUpper(strings.TrimSpace(n.Code))
	if code == "" {
		// Ausência de código é ofensa de `spec-tem-codigo`. Cobrar de novo aqui
		// transformaria um achado em dois, e o segundo sem nada a acrescentar.
		return Skip, "spec sem código declarado — a ausência é cobrada por `spec-tem-codigo`"
	}

	conhecidos := codigosDoMapa(g)
	var orfas []string

	// ── Superfície 1: o testID exposto pela unidade de código ────────────────
	// A unidade se alcança pela aresta `specifies` (spec → código). Ler o irmão
	// pelo nome do arquivo funcionaria no caso co-localizado e falharia em
	// silêncio no resto — o grafo já sabe a resposta.
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.To))
		if err != nil {
			continue
		}
		for _, sigla := range siglasDeTestID(string(b)) {
			if strings.EqualFold(sigla, code) || conhecidos[strings.ToUpper(sigla)] {
				continue
			}
			orfas = append(orfas,
				fmt.Sprintf("testID `%s-…` (%s)", sigla, filepath.Base(e.To)))
		}
	}

	// ── Superfície 2: o baseline de regressão visual ─────────────────────────
	// Convenção: `<Unidade>.<CODE>-VR-<variante>.png`, ao lado da unidade — a mesma
	// que o gate `vr-baseline` usa para encontrá-las. Aqui NÃO vale a dispensa do
	// código-de-outra-unidade: o baseline é a prova DESTA unidade, não um ponteiro
	// para onde ela aparece.
	base := strings.TrimSuffix(n.ID, ".spec.md")
	pngs, _ := doublestar.Glob(os.DirFS(root), base+".*-VR*.png")
	for _, p := range pngs {
		if sigla := siglaDeBaseline(filepath.Base(p)); sigla != "" && !strings.EqualFold(sigla, code) {
			orfas = append(orfas, fmt.Sprintf("baseline `%s` (%s)", sigla, filepath.Base(p)))
		}
	}

	if len(orfas) == 0 {
		return Pass, ""
	}
	sort.Strings(orfas)
	orfas = dedupOrdenado(orfas)
	return Fail, fmt.Sprintf(
		"identidade divergente: a spec declara `%s`, mas a mesma unidade aparece como %s. "+
			"Nenhuma dessas siglas é código de unidade alguma do mapa — então o dicionário de "+
			"códigos, que é GERADO do mapa, não as conhece: o spellcheck as acusa e "+
			"dicionarizá-las cristaliza a divergência em vez de resolvê-la.\n\n"+
			"Unifique pelo código da spec, ou mude o `code:` se a sigla certa for a outra. "+
			"Prefixo que é código de OUTRA unidade não cai aqui — nomear a tela onde o "+
			"elemento aparece (`home-…` num componente da HomeScreen) é reuso legítimo",
		code, strings.Join(orfas, "; "))
}

// codigosDoMapa — todas as identidades declaradas no projeto. É o conjunto que
// distingue reuso deliberado (prefixo que É código de alguém) de identidade órfã.
func codigosDoMapa(g *mapx.Graph) map[string]bool {
	out := map[string]bool{}
	for _, n := range g.Nodes {
		if c := strings.ToUpper(strings.TrimSpace(n.Code)); c != "" {
			out[c] = true
		}
	}
	return out
}

// siglaTestIDRE captura o prefixo de um testID literal ou de template. O `:` de
// marcação é opcional: a convenção é do projeto, e um projeto sem ela deve ser
// confrontado do mesmo jeito.
var siglaTestIDRE = regexp.MustCompile("testID=\\{?[`\"']:?([A-Za-z]{4,5})-")

// siglasDeTestID devolve os prefixos com FORMA DE CÓDIGO (4-5 letras) usados como
// testID. A forma é o primeiro filtro; quem decide se a sigla é órfã é o chamador,
// confrontando-a com os códigos do mapa.
func siglasDeTestID(src string) []string {
	visto := map[string]bool{}
	var out []string
	for _, m := range siglaTestIDRE.FindAllStringSubmatch(src, -1) {
		s := strings.ToUpper(m[1])
		if !visto[s] {
			visto[s] = true
			out = append(out, s)
		}
	}
	return out
}

// siglaDeBaseline extrai o código de `<Unidade>.<CODE>-VR-<variante>.png`.
func siglaDeBaseline(nome string) string {
	partes := strings.Split(nome, ".")
	if len(partes) < 2 {
		return ""
	}
	seg := partes[len(partes)-2] // o segmento imediatamente antes de ".png"
	if i := strings.Index(seg, "-VR"); i > 0 {
		return seg[:i]
	}
	return ""
}

func dedupOrdenado(in []string) []string {
	visto := map[string]bool{}
	out := in[:0]
	for _, s := range in {
		if !visto[s] {
			visto[s] = true
			out = append(out, s)
		}
	}
	return out
}
