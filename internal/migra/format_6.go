package migra

// FORMAT 6 — the triad became the unit.
//
// A spec stopped regulating three pieces: it regulates the code, the feature, the test
// and the documentation, so "triad" named a count that is no longer true. The vocabulary
// is now the UNIT — the spec and every sibling it regulates — and the gate that charges
// its pieces is `unit-complete`.
//
// What a project carries of the old name, and is converted here:
//
//   · the gate `triad-complete`, as the value of `name`, `id` and `check` in the
//     configuration and of `gate` in the map (its judgments and stamps);
//   · the layer key `optional_triad_edges`, which lists the edges a layer waives.
//
// LOSS IF NOT CONVERTED: SILENT for the gate. A project that still declares
// `name: triad-complete` declares a gate this binary no longer knows — the pieces of the
// unit stop being charged, and nothing says so. So the minimum readable format goes up:
// a project not migrated is refused with the message to run `anchors migrate`, instead
// of being read by half. The key alone would fail loudly (the parser refuses it), and is
// converted for the same reason format 4 converted its keys.

// gatesToUnit is the one table of this step's gate renames, applied wherever a gate name is
// written.
var gatesToUnit = map[string]string{
	"triad-complete": "unit-complete",
}

func init() {
	Register(Step{
		To:  6,
		Why: "the triad became the unit: gate triad-complete → unit-complete, key optional_triad_edges → optional_unit_edges",
		RenameKeys: map[string]map[string]string{
			"anchors.yaml": {
				"optional_triad_edges": "optional_unit_edges",
			},
		},
		RenameValues: map[string]map[string]map[string]string{
			"anchors.graph.yaml": {"gate": gatesToUnit},
			"anchors.yaml": {
				"name":  gatesToUnit,
				"id":    gatesToUnit,
				"check": gatesToUnit,
			},
		},
	})
}
