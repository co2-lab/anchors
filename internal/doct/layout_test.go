// @anchors
//   code: LYTSL
//   ref: DCOXX

package doct

import (
	"reflect"
	"strings"
	"testing"
)

func TestLayoutBig(t *testing.T) {
	t.Run("DCOXX-B01: Either measure passing its cut-off makes the selection big", func(t *testing.T) {
		l := Layout{MaxUnits: 20, MaxLines: 2000}
		if !l.Big(Size{Units: 21, Lines: 10}) {
			t.Error("21 units is above the unit cut-off and must be big")
		}
		if !l.Big(Size{Units: 1, Lines: 2001}) {
			t.Error("2001 lines is above the line cut-off and must be big")
		}
		if l.Big(Size{Units: 1, Lines: 10}) {
			t.Error("a selection below both cut-offs is not big")
		}
	})
	t.Run("DCOXX-I01: A selection exactly at the limit is not big", func(t *testing.T) {
		l := Layout{MaxUnits: 20, MaxLines: 2000}
		if l.Big(Size{Units: 20, Lines: 2000}) {
			t.Error("exactly at both limits is not big: the cut-off is strict")
		}
		if !l.Big(Size{Units: 21, Lines: 2000}) || !l.Big(Size{Units: 20, Lines: 2001}) {
			t.Error("one past either limit must be big")
		}
	})
}

func TestDefaultLayout(t *testing.T) {
	t.Run("DCOXX-B02: The default cut-off is 20 units or 2000 lines", func(t *testing.T) {
		if got := DefaultLayout(); got != (Layout{MaxUnits: 20, MaxLines: 2000}) {
			t.Errorf("DefaultLayout() = %+v, want 20 units / 2000 lines", got)
		}
	})
	t.Run("DCOXX-X01: The layout never splits a layer into pages", func(t *testing.T) {
		// The decision's whole surface is a yes/no answer and a sentence: no method of
		// the layout returns a page count or a partition.
		typ := reflect.TypeOf(Layout{})
		var names []string
		for i := 0; i < typ.NumMethod(); i++ {
			names = append(names, typ.Method(i).Name)
		}
		if strings.Join(names, ",") != "Big,Describe" {
			t.Errorf("Layout methods = %v, want only Big and Describe", names)
		}
	})
}

func TestLayoutDescribe(t *testing.T) {
	t.Run("DCOXX-B03: The summary sentence exists only for a big selection and names its numbers", func(t *testing.T) {
		l := Layout{MaxUnits: 20, MaxLines: 2000}
		if got := l.Describe(Size{Units: 3, Rules: 9, Lines: 100}); got != "" {
			t.Errorf("a small selection gets no sentence, got %q", got)
		}
		got := l.Describe(Size{Units: 25, Rules: 90, Lines: 100})
		for _, want := range []string{"25 unidades", "90 regras", "20", "2000"} {
			if !strings.Contains(got, want) {
				t.Errorf("the sentence does not name %q: %q", want, got)
			}
		}
	})
}
