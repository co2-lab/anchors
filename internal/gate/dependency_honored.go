package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkDependencyHonored — gate RELACIONAL da aresta spec→código: confronta a Tabela de
// Dependências de uma spec (os MÉTODOS que ela promete consumir) com o uso REAL no
// código da unidade. Pega a divergência que o feature-test-match não vê: a spec declara
// `DEP2 → metadataVersioning · resolveVersion, applyEdit`, mas o código só usa
// resolveVersion — o applyEdit prometido nunca é chamado (o bug que passou por 65 testes
// verdes no E2E do KV).
//
// Régua (estática, sem execução):
//   - Só confronta SÍMBOLOS: os identificadores entre `backticks` do campo Método. Uma
//     descrição em prosa ("CRUD + consultas") não é confrontável e é ignorada.
//   - Cada símbolo declarado deve aparecer no conteúdo (não-comentário) do código que a
//     spec `specifies`. Símbolo declarado + ausente no código → reprova.
func checkDependencyHonored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "não é uma spec — só spec tem Tabela de Dependências"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	out := g.Neighbors(n.ID).Out

	// símbolos prometidos, por dependência (só os entre backticks — os confrontáveis).
	type promise struct {
		dep     string // DEPn (para a mensagem)
		target  string // arquivo-dependência (para a mensagem)
		symbols []string
	}
	var promises []promise
	for _, e := range out {
		if e.Type != mapx.EdgeDependsOn {
			continue
		}
		syms := backtickedSymbols(e.Method)
		if len(syms) == 0 {
			continue // Método em prosa (ex.: "CRUD + consultas") — não confrontável
		}
		promises = append(promises, promise{dep: e.Dep, target: e.To, symbols: syms})
	}
	if len(promises) == 0 {
		return Skip, "a Tabela de Dependências não promete nenhum SÍMBOLO confrontável " +
			"(só descrição em prosa, ou nenhuma dependência) — nada a verificar"
	}

	// o código que esta spec `specifies` (o consumidor real dos métodos).
	var codePaths []string
	for _, e := range out {
		if e.Type == mapx.EdgeSpecifies {
			codePaths = append(codePaths, e.To)
		}
	}
	if len(codePaths) == 0 {
		return Pending, "spec sem código ligado (specifies) — nada a confrontar ainda"
	}

	// une o conteúdo (não-comentário) do código regido.
	var codeBody strings.Builder
	for _, cp := range codePaths {
		b, err := os.ReadFile(filepath.Join(root, cp))
		if err != nil {
			continue
		}
		codeBody.WriteString(stripLineComments(string(b)))
		codeBody.WriteString("\n")
	}
	code := codeBody.String()

	// confronta cada símbolo prometido.
	var missing []string
	for _, p := range promises {
		for _, sym := range p.symbols {
			if symbolUsed(code, sym) {
				continue
			}
			entry := fmt.Sprintf("%s (%s)", sym, depTargetName(p.target))
			// A promessa quebrada é MUITAS VEZES uma renomeação (ping→pingHeartbeat).
			// Temos como sugerir: procuramos no código um símbolo parecido, para a
			// mensagem dizer o que fazer em vez de só apontar o erro.
			if near := nearestSymbol(code, sym); near != "" {
				entry += " — o código usa `" + near + "`; renomeação?"
			}
			missing = append(missing, entry)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Fail, "a Tabela de Dependências promete símbolo(s) que o código NÃO usa: " +
			strings.Join(missing, ", ") +
			" — ou o código deveria usá-los, ou a spec não deveria declará-los."
	}
	return Pass, ""
}

// backtickedSymbols extrai os identificadores entre `crase` de um campo Método. O scan
// só remove os backticks das BORDAS, então um Método com 2+ símbolos preserva os
// backticks internos (ex.: "resolveVersion`, `applyEdit"). Pegamos cada trecho crasado;
// se não há crase, é prosa (devolve vazio).
var backtickSymRE = regexp.MustCompile("`([A-Za-z_$][A-Za-z0-9_$]*)`")

func backtickedSymbols(method string) []string {
	m := strings.TrimSpace(method)
	if m == "" {
		return nil
	}
	// PROSA não é contrato. Uma célula sem NENHUMA crase é descrição livre ("solicitações",
	// "CRUD de IRDependent") — o autor não está prometendo um símbolo, está explicando a
	// dependência. Só o que está entre crases é promessa verificável.
	if !strings.Contains(m, "`") {
		return nil
	}
	// O scan PRESERVA as crases do campo Método (ver scan.go), então casamos direto.
	// Compatibilidade: se as bordas vierem sem crase (mapa antigo), o wrapped recompõe.
	wrapped := m
	if !strings.HasPrefix(m, "`") {
		wrapped = "`" + wrapped
	}
	if !strings.HasSuffix(wrapped, "`") {
		wrapped += "`"
	}
	var out []string
	seen := map[string]bool{}
	for _, mm := range backtickSymRE.FindAllStringSubmatch(wrapped, -1) {
		s := mm[1]
		// ignora "tipo"/anotações que não são símbolos de valor (ex.: "(tipo)").
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// symbolUsed diz se um identificador aparece como TOKEN no código (fronteira de palavra
// — não casa substring de um nome maior).
func symbolUsed(code, sym string) bool {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(sym) + `\b`)
	return re.MatchString(code)
}

func depTargetName(path string) string {
	return filepath.Base(path)
}

// identInCodeRE lista os identificadores do código, para achar o "quase igual".
var identInCodeRE = regexp.MustCompile(`\b[A-Za-z_$][A-Za-z0-9_$]{2,}\b`)

// nearestSymbol procura no código um identificador que seja plausivelmente o MESMO
// símbolo renomeado: um contém o outro como prefixo/sufixo (ping→pingHeartbeat,
// computeSeatsAmount→computeSeatsAmountCents). Conservador de propósito — só sugere
// quando a relação é de extensão do nome, não por distância de edição genérica (que
// produziria sugestões ruins). Vazio = sem palpite honesto.
func nearestSymbol(code, sym string) string {
	if len(sym) < 4 {
		return "" // curto demais p/ prefixo ser sinal
	}
	lower := strings.ToLower(sym)
	best := ""
	for _, cand := range identInCodeRE.FindAllString(code, -1) {
		if cand == sym {
			continue
		}
		lc := strings.ToLower(cand)
		if !strings.HasPrefix(lc, lower) && !strings.HasSuffix(lc, lower) {
			continue
		}
		// prefere o candidato mais curto (a extensão mínima do nome).
		if best == "" || len(cand) < len(best) {
			best = cand
		}
	}
	return best
}
