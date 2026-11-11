// @anchors
//   code: STRQT
//   ref: STRQS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const alertsScreen = "<!-- @anchors\n  code: ALRTS\n  layer: screen\n-->\n# Alerts\n\nWhen the list is empty a hint shows.\n\n## States\n\n" +
	"### ALRTS-S01: Loading\n\n### ALRTS-S02: Error, with a retry\n\n@no-state: offline — the screen reads only local data\n"

func TestStatesRequired_eachRequiredState(t *testing.T) {
	t.Run("STRQS-B01: A required state the spec neither registers nor waives fails", func(t *testing.T) {})
	cfg := &config.Config{RequiredStates: map[string][]string{"screen": {"loading", "error", "empty", "offline"}}}
	n := mapx.Node{ID: "Alerts.spec.md", Kind: mapx.KindSpec}
	v, msg := checkStatesRequired(alertsScreen, n, "", nil, cfg)
	if v != Fail || !strings.Contains(msg, "empty") || strings.Contains(msg, "loading") || strings.Contains(msg, "offline,") {
		t.Errorf("empty is named only in prose: %v %s", v, msg)
	}
	if v, msg := checkStatesRequired(alertsScreen+"\n### ALRTS-S03: Empty\n", n, "", nil, cfg); v != Pass {
		t.Errorf("every state registered or waived: %v %s", v, msg)
	}
}

func TestStatesRequired_outsideItsSubject(t *testing.T) {
	t.Run("STRQS-B02: Outside its subject the gate is skipped", func(t *testing.T) {})
	n := mapx.Node{ID: "Alerts.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkStatesRequired(alertsScreen, n, "", nil, &config.Config{}); v != Skip {
		t.Errorf("nothing required: %v", v)
	}
	if v, _ := checkStatesRequired(alertsScreen, n, "", nil, &config.Config{RequiredStates: map[string][]string{"api": {"error"}}}); v != Skip {
		t.Errorf("another layer: %v", v)
	}
	if v, _ := checkStatesRequired("x", mapx.Node{ID: "a.ts", Kind: mapx.KindCode}, "", nil, &config.Config{RequiredStates: map[string][]string{"screen": {"error"}}}); v != Skip {
		t.Errorf("no spec: %v", v)
	}
}
