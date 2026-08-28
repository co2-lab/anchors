package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// sibling-guard: funções IRMÃS (exportadas do mesmo módulo, recebendo o mesmo
// parâmetro) devem tratar esse parâmetro de forma consistente. Quando duas guardam e
// uma não, a que não guarda é quase sempre esquecimento — não decisão.
//
// O caso que motivou (real): um módulo de versionamento exportava três funções sobre o
// mesmo histórico. As DUAS de escrita faziam `versions.filter(v => v.key === key)`
// antes de decidir; a de LEITURA não filtrava — e era justamente ela que recebia o
// array multi-chave direto do repositório. Resultado: pedir a versão de "fornecedor"
// devolvia a de "centro-custo", em silêncio. Passou por 11 gates verdes e 16 testes.
//
// A assimetria é o sinal, e é o que torna isto detectável sem entender o domínio: o
// gate não sabe o que a guarda faz, só que as irmãs a aplicam e uma não.
//
// DELIBERADAMENTE conservador — só acusa quando:
//   - há 3+ funções exportadas recebendo o MESMO nome de parâmetro, e
//   - a MAIORIA aplica uma guarda reconhecível sobre ele, e
//   - ao menos uma não aplica NENHUMA.
//
// Com 2 funções não há maioria (é 1×1, ambíguo), e uma guarda minoritária pode ser a
// exceção legítima. O gate erra para o lado de não incomodar.
func checkSiblingGuard(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, "não é código — a assimetria se lê entre funções de um módulo"
	}
	d := cfg.DialectFor()
	if d.ExportedFunc == "" {
		// Pendente, não Pass: sem saber reconhecer uma função, o gate não verificou nada —
		// e um ✓ aqui seria mentira (QUALITY §7, o terceiro estado).
		return Pending, "o projeto não declarou como reconhecer uma função exportada. " +
			"Declare `dialect.family` (" + strings.Join(config.KnownDialectFamilies(), ", ") +
			") ou `dialect.exported_func` no anchors.yaml — sem isso este gate não vê o código"
	}
	fns := exportedFuncs(content, d)
	if len(fns) < 3 {
		return Skip, "menos de 3 funções exportadas — sem maioria para comparar"
	}

	// agrupa por PARÂMETRO comum: qual nome aparece na assinatura de várias.
	porParam := map[string][]exportedFunc{}
	for _, f := range fns {
		for _, p := range f.params {
			porParam[p] = append(porParam[p], f)
		}
	}

	var achados []string
	for param, irmas := range porParam {
		if len(irmas) < 3 {
			continue
		}
		var comGuarda, semGuarda []string
		for _, f := range irmas {
			switch {
			case guardaSobre(f.body, param):
				comGuarda = append(comGuarda, f.name)
			case noGuardRE.MatchString(f.body):
				// Opt-out HONESTO, com razão escrita. Esta função não precisa da guarda,
				// e o autor diz por quê — tipicamente porque DELEGA a uma irmã que já
				// guarda, e a "assimetria" é aparente.
				//
				// Medido: das 3 assimetrias de um repositório real, uma era exatamente
				// isso — `isProbableDuplicate` chamava `similarity`, que guarda uma linha
				// abaixo. Sem a saída declarada, o autor só tinha duas opções ruins:
				// duplicar a guarda para calar o gate, ou conviver com o achado para
				// sempre. O gate irmão (`pagination-honored`) já oferecia `@no-paginate`;
				// este não oferecia nada.
			default:
				semGuarda = append(semGuarda, f.name)
			}
		}
		// maioria guarda, minoria não → a minoria é suspeita
		if len(comGuarda) > len(semGuarda) && len(semGuarda) > 0 {
			sort.Strings(comGuarda)
			sort.Strings(semGuarda)
			achados = append(achados, fmt.Sprintf(
				"`%s` recebe `%s` sem nenhuma guarda, enquanto %s guardam (%s)",
				strings.Join(semGuarda, "`, `"), param,
				plural(len(comGuarda), "a irmã", "as irmãs"),
				strings.Join(comGuarda, ", ")))
		}
	}
	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	return Fail, "assimetria entre funções irmãs: " + strings.Join(achados, "; ") +
		". Ou a guarda falta na que não tem, ou ela é desnecessária nas que têm — " +
		"as duas leituras não podem estar certas ao mesmo tempo. Se a função não precisa " +
		"da guarda (porque DELEGA a uma irmã que já guarda, por exemplo), marque " +
		"`// @no-guard: <razão>` no corpo dela — dispensar é decisão registrada; duplicar " +
		"a guarda só para calar o gate esconde a razão."
}

func plural(n int, um, muitos string) string {
	if n == 1 {
		return um
	}
	return muitos
}

type exportedFunc struct {
	name   string
	params []string
	body   string
}

// exportedFuncs extrai as funções exportadas usando o DIALETO do projeto — nenhuma
// sintaxe de linguagem vive aqui. O custo de cravar já foi medido: a primeira versão
// deste gate exigia `export function` e ficava cega para 291 de 298 funções dos models do
// app de referência (97,6% usam `export async function`), passando verde por não ver nada.
// Ver internal/config/dialect.go.
func exportedFuncs(content string, d config.Dialect) []exportedFunc {
	funcRE := d.Compile(d.ExportedFunc)
	if funcRE == nil {
		return nil // dialeto não declarado — quem chama responde Pendente
	}
	paramRE := d.Compile(d.ParamName)
	locs := funcRE.FindAllStringSubmatchIndex(content, -1)
	var out []exportedFunc
	for i, loc := range locs {
		// a regex tem dois grupos alternativos (function / const arrow); vale o que casou.
		name := ""
		for g := 1; g*2+1 < len(loc); g++ {
			if loc[g*2] >= 0 {
				name = content[loc[g*2]:loc[g*2+1]]
				break
			}
		}
		if name == "" {
			continue
		}
		fim := len(content)
		if i+1 < len(locs) {
			fim = locs[i+1][0]
		}
		corpo := content[loc[1]:fim]
		// a assinatura vai até o `)` que fecha os parâmetros
		sig := corpo
		if j := strings.Index(corpo, "{"); j > 0 {
			sig = corpo[:j]
		}
		var params []string
		if paramRE != nil {
			for _, m := range paramRE.FindAllStringSubmatch(sig, -1) {
				params = append(params, m[1])
			}
		}
		out = append(out, exportedFunc{name: name, params: params, body: corpo})
	}
	return out
}

// guardaSobre diz se o corpo aplica alguma guarda RECONHECÍVEL sobre o parâmetro:
// filtrar, validar, lançar, ou retornar cedo por causa dele. Não entende semântica —
// procura o parâmetro perto de uma construção de guarda.
func guardaSobre(body, param string) bool {
	p := regexp.QuoteMeta(param)
	padroes := []string{
		`\.filter\([^)]*` + p + `\b`,                              // versions.filter(v => v.key === key)
		`\b` + p + `\s*\.\s*filter\(`,                             // param.filter(...)
		`if\s*\([^)]*\b` + p + `\b[^)]*\)\s*\{?\s*(throw|return)`, // guarda explícita
		`\b` + p + `\b[^\n]*\?\?`,                                 // default via ??
		`(throw|assert)[^\n]*\b` + p + `\b`,                       // validação que lança
	}
	for _, pat := range padroes {
		if regexp.MustCompile(pat).MatchString(body) {
			return true
		}
	}
	return false
}

// noGuardRE: a dispensa DECLARADA da guarda, com razão escrita — marcador nu não dispensa.
var noGuardRE = regexp.MustCompile(`@no-guard[^\S\n]*:[^\S\n]*[^\s|]\S*`)
