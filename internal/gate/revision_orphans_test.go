package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// THE FIXTURES ARE IN PORTUGUESE ON PURPOSE. They reproduce the spec of the reference app
// — a project written in Portuguese — and the ruler has to work on it: the negation test
// below only proves anything because it feeds "não" through `extraStopwords`. They are
// input data, not prose of this project.

// The model spec reproduces the shape that produced the gate: a revision that changed the
// badge from COUNT to DOT, naming B03 and B04 — and leaving the invariant I02, which speaks
// of the same badge, unmentioned.
const specWithRevision = `<!-- @anchors
  code: NTCNN
-->
# NotificationCenter

> **NTCNN-R0002:** o badge é um PONTO, não um número.
>
> **Revises:** ` + "`B03`, `B04`" + `

### NTCNN-B01 — A lista é de entregas, e uma entrega não muda de estado

### NTCNN-B03 — O badge é um PONTO, e não conta ao abrir

### NTCNN-B04 — Sem notificação nova, o sino fica limpo

### NTCNN-I02 — O badge nunca conta o que a lista não mostra
`

func specNodeRev() mapx.Node {
	return mapx.Node{ID: "n.spec.md", Kind: mapx.KindSpec, Code: "NTCNN"}
}

func TestRevisionOrphans_notASpec(t *testing.T) {
	t.Run("RVORP-B01: Non-spec artifacts skip confrontation", func(t *testing.T) {})
	n := mapx.Node{ID: "p.md", Kind: mapx.KindPlan}
	if v, _ := checkRevisionOrphans(specWithRevision, n, "", nil, nil); v != Skip {
		t.Errorf("not a spec and the verdict was %v", v)
	}
}

func TestRevisionOrphans_noRevision(t *testing.T) {
	t.Run("RVORP-B02: A spec with no revision has nothing to confront", func(t *testing.T) {})
	noRev := "### NTCNN-B01 — uma regra\n\n### NTCNN-B02 — outra\n"
	if v, _ := checkRevisionOrphans(noRev, specNodeRev(), "", nil, nil); v != Skip {
		t.Errorf("with no revision Skip was expected, got %v", v)
	}
}

// A revision that names nothing cannot be confronted against any sibling — and it is the
// sibling that keeps asserting what the revision revoked.
func TestRevisionOrphans_noRevisesField(t *testing.T) {
	t.Run("RVORP-B03: A revision naming no rule cannot be confronted", func(t *testing.T) {})
	noField := strings.Replace(specWithRevision, "> **Revises:** `B03`, `B04`\n", "", 1)
	// PENDING, not Fail: the field is new, and 439 revisions written before it cannot all
	// be born accused. Pending states the truth — there is no way to confront — without
	// charging whoever wrote before the ruler existed.
	v, msg := checkRevisionOrphans(noField, specNodeRev(), "", nil, nil)
	if v != Pending {
		t.Fatalf("a revision without `Revises:` should be Pending, got %v", v)
	}
	if !strings.Contains(msg, "Revises") {
		t.Errorf("the verdict does not name the missing field: %q", msg)
	}
}

func TestRevisionOrphans_unknownCode(t *testing.T) {
	t.Run("RVORP-B04: A revision naming a rule the spec does not define fails", func(t *testing.T) {})
	withGhost := strings.Replace(specWithRevision, "`B03`, `B04`", "`B03`, `B99`", 1)
	v, msg := checkRevisionOrphans(withGhost, specNodeRev(), "", nil, nil)
	if v != Fail {
		t.Fatalf("unknown code and the verdict was %v", v)
	}
	if !strings.Contains(msg, "B99") {
		t.Errorf("the verdict does not carry the code: %q", msg)
	}
}

// THE CASE THAT PRODUCED THE GATE: I02 speaks of the same badge and was not mentioned.
func TestRevisionOrphans_orphanedSiblingIsReported(t *testing.T) {
	t.Run("RVORP-B05: A sibling sharing vocabulary and left unmentioned is reported", func(t *testing.T) {})
	v, msg := checkRevisionOrphans(specWithRevision, specNodeRev(), "", nil, nil)
	if v != Fail {
		t.Fatalf("I02 speaks of the badge and was not mentioned — expected Fail, got %v", v)
	}
	if !strings.Contains(msg, "I02") {
		t.Errorf("the verdict does not name the orphan: %q", msg)
	}
	if !strings.Contains(msg, "badge") {
		t.Errorf("the verdict does not say WHICH vocabulary links the two: %q", msg)
	}
	// B01 speaks of deliveries and the list, not of the badge: it must not come in.
	if strings.Contains(msg, "B01") {
		t.Errorf("it reported a rule that does not share the subject: %q", msg)
	}
}

func TestRevisionOrphans_everythingMentionedPasses(t *testing.T) {
	t.Run("RVORP-B06: Every vocabulary-sharing sibling accounted for passes", func(t *testing.T) {})
	complete := strings.Replace(specWithRevision,
		"> **Revises:** `B03`, `B04`", "> **Revises:** `B03`, `B04`\n>\n> **Checked:** `I02`", 1)
	if v, msg := checkRevisionOrphans(complete, specNodeRev(), "", nil, nil); v != Pass {
		t.Errorf("every sibling accounted for and the verdict was %v: %s", v, msg)
	}
}

// `Checked` does not assert that the rule is right — it asserts that somebody looked. It
// is cheap to write AFTER reading, and impossible to write honestly without reading.
func TestRevisionOrphans_checkedClearsTheAccusation(t *testing.T) {
	t.Run("RVORP-B07: Checked clears the accusation without asserting correctness", func(t *testing.T) {})
	withChecked := strings.Replace(specWithRevision,
		"> **Revises:** `B03`, `B04`", "> **Revises:** `B03`, `B04`\n>\n> **Checked:** `I02`", 1)
	_, msg := checkRevisionOrphans(withChecked, specNodeRev(), "", nil, nil)
	if strings.Contains(msg, "I02") {
		t.Errorf("I02 was checked and is still reported: %q", msg)
	}
}

func TestRevisionOrphans_ruleNeverAccusesItself(t *testing.T) {
	t.Run("RVORP-I01: A revised rule is never its own orphan", func(t *testing.T) {})
	_, msg := checkRevisionOrphans(specWithRevision, specNodeRev(), "", nil, nil)
	// B03 is the revised rule itself, and it speaks of the badge: if it accused itself, it
	// would show up here.
	for _, line := range strings.Split(msg, "\n") {
		if strings.Contains(line, "B03 —") {
			t.Errorf("the revised rule accused itself: %q", line)
		}
	}
}

// NEGATION does not link two rules.
//
// Half the rules of a well-written spec state what the unit does NOT do, and "não" linked
// every rule to every other. Measured on the real spec: with negation counting, three
// rules were reported; without it, one — the one that actually contradicted.
func TestRevisionOrphans_negationDoesNotLink(t *testing.T) {
	t.Run("RVORP-X01: Terms shared by the whole unit do not discriminate", func(t *testing.T) {})
	fixture := `> **NTCNN-R0001:** mudou.
>
> **Revises:** ` + "`B01`" + `

### NTCNN-B01 — O badge não some sozinho
### NTCNN-B02 — A ordem não muda por filtro
### NTCNN-I02 — O badge engana quando conta
`
	v, msg := checkRevisionOrphans(fixture, specNodeRev(), "", nil, nil)
	if v != Fail || !strings.Contains(msg, "I02") {
		t.Fatalf("I02 shares `badge` with the revised rule: %v / %s", v, msg)
	}
	// B02 shares only the negation with B01.
	if strings.Contains(msg, "B02 —") {
		t.Errorf("negation linked two rules: %q", msg)
	}
}

// A WAIVER IN A COMMENT is not part of what the rule asserts.
//
// Measured on the real spec: `B07` carries two `@no-*` with a written reason and came out
// with 34 terms, against 4 to 6 for its siblings — sharing vocabulary with almost all of
// them by accident of prose, not by subject.
func TestRevisionOrphans_commentStaysOutOfTheVocabulary(t *testing.T) {
	fixture := `> **NTCNN-R0001:** mudou.
>
> **Revises:** ` + "`B01`" + `

### NTCNN-B01 — O badge fica limpo
### NTCNN-B02 — A ordem segue a chegada <!-- @no-scenario: o badge decide, e esta unidade não desenha nada -->
`
	_, msg := checkRevisionOrphans(fixture, specNodeRev(), "", nil, nil)
	if strings.Contains(msg, "B02 —") {
		t.Errorf("the `badge` cited in the WAIVER became the rule's subject: %q", msg)
	}
}

// And the REAL case, reproduced: one domain word is enough when the title is clean.
func TestRevisionOrphans_oneDomainWordIsEnough(t *testing.T) {
	fixture := `> **NTCNN-R0001:** o badge é um ponto.
>
> **Revises:** ` + "`B03`" + `

### NTCNN-B03 — O badge é um PONTO, e some ao abrir
### NTCNN-B04 — Sem notificação nova, o sino fica limpo
### NTCNN-I02 — O badge nunca conta o que a lista não mostra
`
	v, msg := checkRevisionOrphans(fixture, specNodeRev(), "", nil, nil)
	if v != Fail || !strings.Contains(msg, "I02") {
		t.Fatalf("`badge` alone links when the title is clean: %v / %s", v, msg)
	}
	if strings.Contains(msg, "B04 —") {
		t.Errorf("B04 does not speak of the badge and was reported: %q", msg)
	}
}
