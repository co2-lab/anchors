package config

// --- o VOCABULÁRIO em inglês ---
//
// Os nomes de gate são IDENTIFICADORES: eles vão para o `anchors.yaml` de cada projeto,
// para tutoriais, para respostas de fórum. São fixos em inglês e NÃO se traduzem — um
// `anchors.yaml` escrito por um time brasileiro tem de funcionar num time espanhol sem
// tradução nenhuma.
//
// A TABELA DE ALIAS QUE VIVIA AQUI FOI REMOVIDA.
//
// Ela aceitava os nomes em português — `regra-cumprida`, `unit-complete`, `spec-completa`
// e outros dezoito — e os convertia na carga, para sempre. Isso resolve e não fecha: o
// arquivo nunca se conserta, e o mapa acumula carimbos nos DOIS formatos.
//
// E tinha um defeito de assimetria: a LEITURA normalizava (`mapx.mesmoGate`), a ESCRITA
// não. Um projeto que renomeasse o gate ganhava um SEGUNDO carimbo em vez de atualizar o
// primeiro. Medido no app de referência: 40 julgamentos gravados como `regra-cumprida` convivendo
// com 2 como `rule-fulfilled` — o mesmo gate contado duas vezes, e o `check` rejulgando o
// que já fora respondido.
//
// Com o contrato de formato, o alias deixou de ser necessário: a conversão virou o passo
// de migração `1→2` (ver `internal/migra/formato_2.go`), roda UMA VEZ, e o formato 2 só
// tem nome canônico. Quem está no formato 1 é recusado com a mensagem que manda migrar —
// não lido pela metade.
// gates default, que vive no `initx`. E `config` não pode importar `initx` — seria ciclo,
// já que o `initx` monta gates a partir de tipos daqui.
//
// A injeção resolve: o `initx` registra a lista no seu `init()`, e o teste a consulta. Sem
// isso o teste não teria contra o que confrontar, e um de-para apontando para o vazio
// passaria — o projeto carregaria, o gate viraria um nome que nenhum verificador conhece,
// e o `check` reportaria "gate sem nada para medir" sem dizer que a causa foi a conversão.

var defaultGateNames func() []string

// RegisterGateNames liga a lista de gates default ao pacote config.
func RegisterGateNames(f func() []string) { defaultGateNames = f }

// DefaultGateNamesForTest devolve os nomes registrados, ou nil se ninguém registrou.
func DefaultGateNamesForTest() []string {
	if defaultGateNames == nil {
		return nil
	}
	return defaultGateNames()
}

// --- the LETTERS of the artifacts that are not specs ---
//
// A code's letter says what the item is: `-B01` a behaviour, `-S01` a state. Specs take
// theirs from `rule_types` (or the canonical `DefaultRuleLetters`); plans, flows and
// actions have fixed letters of their own, declared here and read by everything that
// writes or recognizes them — templates, guides, gates, the flow graph.
//
// WHY THESE THREE, AND WHAT THEY REPLACED (decided on 2026-09-28, migrated by format 5):
//
//   - A plan's PHASE was `F`, from the Portuguese *Fase*. In English it would be `P`, which
//     became the canonical letter of *Presentation validations*. It is now `W` — Wave: a
//     plan delivers in waves, each ordered after the ones it depends on.
//   - A flow's STEP was `P`, from the Portuguese *Passo*. `S` (Step) was not free: it is a
//     spec's *State*. It is now `T` — Task: each step fits one action (`Fits: ACHCK`), so it
//     is the task the flow assigns at that point.
//   - A RESULT — of an action, or of a flow fitted as a piece inside another — was `R`
//     (*Result*), the same letter as a spec's *permission* rule —
//     and as a product doctrine's rule and a plan's revision block (`-R0001`, four digits).
//     It is now `O` — Outcome.
//
// The criterion: the vocabulary is English, and a letter of a non-spec artifact does not
// reuse a canonical spec letter with another meaning. Only letters no spec uses were
// candidates. There is NO alias for the old letters, on the precedent of the gate names
// above: reading both forever never closes, so `anchors migrate` (format 5) rewrites the
// project's codes once, and the old letters are not read after it.
const (
	// PhaseLetter is the letter of a plan's phase: `FNDTN-W01`.
	PhaseLetter = "W"
	// StepLetter is the letter of a flow's step: `WORKR-T04`.
	StepLetter = "T"
	// OutcomeLetter is the letter of a result, of an action (`ACHCK-O01`) or of a flow
	// fitted as a piece (`JUDGE-O01`).
	OutcomeLetter = "O"
)
