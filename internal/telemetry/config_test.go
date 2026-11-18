// @anchors
//   code: CNTSE
//   ref: TLCNT

package telemetry

import (
	"testing"
)

// The OPT-OUT must be easy to GET RIGHT. Whoever looks for how to turn it off tries an
// obvious word, and requiring the exact form would be a trap — the person would think it was off.
func TestDisabled_acceptsTheFormsSomeoneWouldTry(t *testing.T) {
	t.Run("TLCNT-B01: Every form someone would try turns telemetry off", func(t *testing.T) {})
	for _, v := range []string{"off", "OFF", "Off", "0", "false", "FALSE", "no", " off "} {
		t.Setenv(EnvVar, v)
		if !Disabled("") {
			t.Errorf("%q should turn it off — whoever writes this believes it is off", v)
		}
	}
}

func TestDisabled_theDefaultIsOn(t *testing.T) {
	t.Run("TLCNT-B02: With nothing declared telemetry is on", func(t *testing.T) {})
	t.Setenv(EnvVar, "")
	if Disabled("") {
		t.Error("with nothing declared, telemetry is on — it is the maintainer's decision")
	}
}

// THE ENVIRONMENT OVERRIDES THE FILE, and it is not a detail: whoever runs in CI must turn it
// off without committing. Committing to turn telemetry off would make one person's decision a
// change in the team's repository.
func TestDisabled_theEnvironmentOverridesTheFile(t *testing.T) {
	t.Run("TLCNT-B03: The environment overrides the file in both directions", func(t *testing.T) {})
	t.Setenv(EnvVar, "on")
	if Disabled("off") {
		t.Error("the environment saying `on` must override the file saying `off`")
	}
	t.Setenv(EnvVar, "off")
	if !Disabled("on") {
		t.Error("the environment saying `off` must override the file saying `on`")
	}
}

// Without the environment, the file decides.
func TestDisabled_withoutTheEnvironmentTheFileDecides(t *testing.T) {
	t.Run("TLCNT-B04: Without the environment the file decides", func(t *testing.T) {})
	t.Setenv(EnvVar, "")
	if !Disabled("off") {
		t.Error("`telemetry: off` in anchors.yaml must turn it off")
	}
	if Disabled("on") {
		t.Error("`telemetry: on` keeps it on")
	}
}

// The same word means the same thing whichever source carries it.
func TestDisabled_aValueMeansTheSameInBothSources(t *testing.T) {
	t.Run("TLCNT-I01: A value means the same in the environment and in the file", func(t *testing.T) {})
	for _, v := range []string{"off", "0", "false", "no", " No ", "on"} {
		t.Setenv(EnvVar, v)
		fromEnv := Disabled("")
		t.Setenv(EnvVar, "")
		fromFile := Disabled(v)
		if fromEnv != fromFile {
			t.Errorf("%q: environment says disabled=%v, file says disabled=%v", v, fromEnv, fromFile)
		}
	}
}

// Only the closed list turns it off: anything else keeps the default.
func TestDisabled_aWordOutsideTheListKeepsItOn(t *testing.T) {
	t.Run("TLCNT-X01: A word outside the closed list keeps telemetry on", func(t *testing.T) {})
	for _, v := range []string{"yes", "disabled", "of"} {
		t.Setenv(EnvVar, v)
		if Disabled("") {
			t.Errorf("environment %q should keep telemetry on", v)
		}
		t.Setenv(EnvVar, "")
		if Disabled(v) {
			t.Errorf("file %q should keep telemetry on", v)
		}
	}
}
