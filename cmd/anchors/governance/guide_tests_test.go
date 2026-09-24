package governance

import (
	"strings"
	"testing"
)

// A boundary rule applies to every input; the guide says which instrument proves it for
// each shape of input space (blue-eyes #1016: nine survivors in four units, each a boundary
// closed for one case and open for the neighbour).
func TestTestGuideNamesTheInstrumentPerInputShape(t *testing.T) {
	for _, want := range []string{"SMALL AND CLOSED", "EXHAUSTIVE", "LARGE BUT STRUCTURED", "TABLE OF CLASSES", "OPEN", "never a\n  list of the forbidden"} {
		if !strings.Contains(testGuide, want) {
			t.Errorf("the test guide should say %q", want)
		}
	}
}

// The test guide teaches the refresh with the stamp: whoever changes a doubled function
// runs it and adjusts the doubles it lists.
func TestTestGuideTeachesTheStampRefresh(t *testing.T) {
	for _, want := range []string{"anchors stamp --refresh", "mock-stamped", "same commit"} {
		if !strings.Contains(testGuide, want) {
			t.Errorf("the test guide should say %q", want)
		}
	}
}
