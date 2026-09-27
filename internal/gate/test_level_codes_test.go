package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

var levelFeatureNode = mapx.Node{Kind: mapx.KindFeature, ID: "screen/S.feature"}

// levelsCfg declares the levels on the gate entry that runs the check.
func levelsCfg(levels map[string]config.TestLevel) *config.Config {
	return &config.Config{Gates: []config.Gate{{Name: "test-level-codes", Check: "test-level-codes", Levels: levels}}}
}

// levelCfg: the visual level accepts only VR codes; the unit level refuses them.
func levelCfg() *config.Config {
	return levelsCfg(map[string]config.TestLevel{
		"vr-level":   {Allow: []string{`-VR$`}},
		"unit-level": {Exclude: []string{`-VR$`}},
	})
}

func TestTestLevelCodes_NotFeatureSkips(t *testing.T) {
	t.Run("TLVCD-B01: A node that is not a feature is skipped", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-B04 @vr-level\n  Scenario: a\n"
	if v, _ := checkTestLevelCodes(f, mapx.Node{Kind: mapx.KindCode, ID: "S.go"}, "", nil, levelCfg()); v != Skip {
		t.Fatalf("a non-feature must be skipped, got %v", v)
	}
}

func TestTestLevelCodes_NothingDeclaredSkips(t *testing.T) {
	t.Run("TLVCD-B02: A project with no filter per level is skipped", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-B04 @vr-level\n  Scenario: a\n"
	for _, cfg := range []*config.Config{nil, {}, levelsCfg(nil), {Gates: []config.Gate{{Name: "other", Check: "vr-baseline", Levels: map[string]config.TestLevel{"vr-level": {Allow: []string{`-VR$`}}}}}}} {
		if v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, cfg); v != Skip {
			t.Fatalf("without levels on a test-level-codes gate every code is accepted, got %v (%s)", v, msg)
		}
	}
}

func TestTestLevelCodes_NoFilteredLevelSkips(t *testing.T) {
	t.Run("TLVCD-B03: A feature with no filtered level is skipped", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-VR @integration-level\n  Scenario: a\n"
	if v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, levelCfg()); v != Skip {
		t.Fatalf("a level without a filter accepts everything, got %v (%s)", v, msg)
	}
}

func TestTestLevelCodes_AllowAcceptsOnlyMatches(t *testing.T) {
	t.Run("TLVCD-B04: A level with allow accepts only matching codes", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-B04 @vr-level\n  Scenario: a\n"
	v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, levelCfg())
	if v != Fail || !strings.Contains(msg, "SMCSX-B04 (@vr-level)") {
		t.Fatalf("a rule code under an allow-only-VR level must fail, got %v (%s)", v, msg)
	}
}

func TestTestLevelCodes_ExcludeRefusesEvenWhenAllowed(t *testing.T) {
	t.Run("TLVCD-B05: A level with exclude refuses matching codes even when allowed", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-VR @unit-level\n  Scenario: a\n"
	v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, levelCfg())
	if v != Fail || !strings.Contains(msg, "SMCSX-VR (@unit-level)") {
		t.Fatalf("a VR code under a level that excludes it must fail, got %v (%s)", v, msg)
	}
	both := levelsCfg(map[string]config.TestLevel{
		"unit-level": {Allow: []string{`^SMCSX-`}, Exclude: []string{`-VR$`}},
	})
	if v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, both); v != Fail {
		t.Fatalf("exclude must win over allow, got %v (%s)", v, msg)
	}
}

func TestTestLevelCodes_EveryCodeWithoutSuffix(t *testing.T) {
	t.Run("TLVCD-B06: Every code of the scenario is confronted without its suffix", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-VR#02 @SMCSX-B04 @vr-level\n  Scenario: a\n"
	v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, levelCfg())
	if v != Fail || !strings.Contains(msg, "SMCSX-B04") || strings.Contains(msg, "SMCSX-VR") {
		t.Fatalf("the second code must be confronted and the first read without #02, got %v (%s)", v, msg)
	}
}

func TestTestLevelCodes_MessageSortedWithLevels(t *testing.T) {
	t.Run("TLVCD-B07: The failure names each refused code with its level, sorted", func(t *testing.T) {})
	f := "Feature: X\n\n  @ZZZZX-VR @unit-level\n  Scenario: a\n\n  @AAAAX-B01 @vr-level\n  Scenario: b\n"
	v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, levelCfg())
	a, z := strings.Index(msg, "AAAAX-B01 (@vr-level)"), strings.Index(msg, "ZZZZX-VR (@unit-level)")
	if v != Fail || a < 0 || z < 0 || a > z {
		t.Fatalf("both codes with their levels, sorted, got %v (%s)", v, msg)
	}
}

func TestTestLevelCodes_EntriesAddUp(t *testing.T) {
	t.Run("TLVCD-B09: Levels come from the gate's own entries, and two entries add their lists", func(t *testing.T) {})
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "a", Check: "test-level-codes", Levels: map[string]config.TestLevel{"unit-level": {Exclude: []string{`-VR$`}}}},
		{Name: "b", Check: "test-level-codes", Levels: map[string]config.TestLevel{"unit-level": {Exclude: []string{`-B99$`}}}},
	}}
	for _, code := range []string{"SMCSX-VR", "SMCSX-B99"} {
		f := "Feature: X\n\n  @" + code + " @unit-level\n  Scenario: a\n"
		if v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, cfg); v != Fail {
			t.Errorf("%s: both entries' exclusions must hold, got %v (%s)", code, v, msg)
		}
	}
	allows := &config.Config{Gates: []config.Gate{
		{Name: "a", Check: "test-level-codes", Levels: map[string]config.TestLevel{"vr-level": {Allow: []string{`-VR$`}}}},
		{Name: "b", Check: "test-level-codes", Levels: map[string]config.TestLevel{"vr-level": {Allow: []string{`-B04$`}}}},
	}}
	for _, code := range []string{"SMCSX-VR", "SMCSX-B04"} {
		f := "Feature: X\n\n  @" + code + " @vr-level\n  Scenario: a\n"
		if v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, allows); v != Pass {
			t.Errorf("%s: both entries' allowances must hold, got %v (%s)", code, v, msg)
		}
	}
}

func TestTestLevelCodes_AcceptedPass(t *testing.T) {
	t.Run("TLVCD-B08: Accepted codes pass", func(t *testing.T) {})
	f := "Feature: X\n\n  @SMCSX-VR @vr-level\n  Scenario: a\n\n  @SMCSX-B04 @unit-level\n  Scenario: b\n"
	if v, msg := checkTestLevelCodes(f, levelFeatureNode, "", nil, levelCfg()); v != Pass {
		t.Fatalf("accepted codes must pass, got %v (%s)", v, msg)
	}
}
