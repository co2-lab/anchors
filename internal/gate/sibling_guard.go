package gate

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
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
		return Skip, i18n.T("gate.sibling_guard.skip_not_code")
	}
	d := cfg.DialectFor()
	if d.ExportedFunc == "" {
		// Pendente, não Pass: sem saber reconhecer uma função, o gate não verificou nada —
		// e um ✓ aqui seria mentira (QUALITY §7, o terceiro estado).
		return Pending, i18n.T("gate.sibling_guard.pending_dialect_missing", strings.Join(config.KnownDialectFamilies(), ", "))
	}
	fns := exportedFuncs(content, d)
	if len(fns) < 3 {
		return Skip, i18n.T("gate.sibling_guard.skip_few_functions")
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
		// O piso de tres vale por PARAMETRO, e nao so por arquivo: tres funcoes exportadas
		// onde so duas compartilham o parametro ainda sao um par, e par nao tem maioria.
		//
		// A condicao e' REDUNDANTE com a da maioria logo abaixo (`comGuarda > semGuarda`
		// com duas irmas da' 1 > 1, falso) — achado por mutacao: troca-la por `< 2` deixa a
		// suite verde, e nao por falta de teste. Fica como guarda EXPLICITA porque a
		// intencao ("tres e' o piso") nao se le na aritmetica da maioria, e quem mexer nela
		// depois nao tem como saber que estava apoiada nisso.
		if len(irmas) < 3 {
			continue
		}
		var comGuarda, semGuarda []string
		for _, f := range irmas {
			switch {
			case guardOver(f.body, param, d):
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
			sisterWord := i18n.T("gate.sibling_guard.sister_plural")
			if len(comGuarda) == 1 {
				sisterWord = i18n.T("gate.sibling_guard.sister_singular")
			}
			achados = append(achados, i18n.T("gate.sibling_guard.item_asymmetry",
				strings.Join(semGuarda, "`, `"), param,
				sisterWord,
				strings.Join(comGuarda, ", ")))
		}
	}
	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	return Fail, i18n.T("gate.sibling_guard.fail_asymmetry", strings.Join(achados, "; "))
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

// guardOver diz se o corpo aplica alguma guarda RECONHECÍVEL sobre o parâmetro:
// filtrar, validar, lançar, ou retornar cedo por causa dele. Não entende semântica —
// procura o parâmetro perto de uma construção de guarda.
func guardOver(body, param string, d config.Dialect) bool {
	p := regexp.QuoteMeta(param)
	if len(d.GuardPatterns) > 0 {
		for _, pat := range d.GuardPatterns {
			reStr := strings.ReplaceAll(pat, "{{param}}", p)
			re := d.Compile(reStr)
			if re != nil && re.MatchString(body) {
				return true
			}
		}
		return false
	}
	padroes := []string{
		`\.filter\([^)]*` + p + `\b`,
		`\b` + p + `\s*\.\s*filter\(`,
		`if\s*\(?[^\n{:]*\b` + p + `\b[^\n{:]*\)?\s*[:{]?\s*(throw|return|raise|panic)`,
		`\b` + p + `\b[^\n]*\?\?`,
		`(throw|assert|raise|panic)[^\n]*\b` + p + `\b`,
	}
	for _, pat := range padroes {
		re := d.Compile(pat)
		if re != nil && re.MatchString(body) {
			return true
		}
	}
	return false
}

// noGuardRE: a dispensa DECLARADA da guarda, com razão escrita — marcador nu não dispensa.
var noGuardRE = regexp.MustCompile(`@no-guard[^\S\n]*:[^\S\n]*[^\s|]\S*`)
