package gate

import (
	"reflect"
	"strings"
	"testing"
)

func TestExposedTestIDs_collectsEveryOccurrence(t *testing.T) {
	src := "<View testID={`row-${id}`} /><Text testID={`cell-${i}`} />\n" +
		"<Icon testID={on ? 'icon-on' : 'icon-off'} />\n"
	got := exposedTestIDs(src, "testID")
	want := []string{"row-*", "cell-*", "icon-on", "icon-off"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("every template and both ternary branches:\n got %v\nwant %v", got, want)
	}
}

func TestDeclaredTestIDs_everyIdOnALine(t *testing.T) {
	spec := "## Test IDs\n\n- `abcd-screen`, `abcd-close`\n"
	got := declaredTestIDs(spec, "testID")
	if !reflect.DeepEqual(got, []string{"abcd-screen", "abcd-close"}) {
		t.Errorf("got %v", got)
	}
}

func TestList_truncatesAfterEight(t *testing.T) {
	ids := func(n int) []string {
		var out []string
		for i := 0; i < n; i++ {
			out = append(out, string(rune('a'+i)))
		}
		return out
	}
	if got := list(ids(8)); got != "`a`, `b`, `c`, `d`, `e`, `f`, `g`, `h`" {
		t.Errorf("eight ids are listed whole: %s", got)
	}
	got := list(ids(10))
	if !strings.HasSuffix(got, "`h` (+2)") || strings.Contains(got, "`i`") {
		t.Errorf("ten ids list eight and count two: %s", got)
	}
}
