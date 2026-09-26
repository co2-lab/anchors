package initx

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAvailablePacks(t *testing.T) {
	got := AvailablePacks()
	want := map[string][]string{
		"accessibility": {"accessibility/wcag"},
		"health":        {"health/hipaa"},
		"payment":       {"payment/pci-dss"},
		"privacy":       {"privacy/ccpa", "privacy/gdpr", "privacy/lgpd"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AvailablePacks = %v, want %v", got, want)
	}
}

func TestSeedPacks_copiesAllAndPreservesAdapted(t *testing.T) {
	root := t.TempDir()

	// The project already adapted its LGPD pack: seeding must not overwrite it.
	adapted := filepath.Join(root, "packs", "privacy", "lgpd.yaml")
	os.MkdirAll(filepath.Dir(adapted), 0o755)
	os.WriteFile(adapted, []byte("name: lgpd # adapted\n"), 0o644)

	criados, preservados, err := SeedPacks(root)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"packs/privacy/lgpd.yaml"}; !reflect.DeepEqual(preservados, want) {
		t.Fatalf("preserved = %v, want %v", preservados, want)
	}
	wantCreated := []string{
		"packs/accessibility/wcag.yaml",
		"packs/health/hipaa.yaml",
		"packs/payment/pci-dss.yaml",
		"packs/privacy/ccpa.yaml",
		"packs/privacy/gdpr.yaml",
	}
	if !reflect.DeepEqual(criados, wantCreated) {
		t.Fatalf("created = %v, want %v", criados, wantCreated)
	}
	if b, _ := os.ReadFile(adapted); string(b) != "name: lgpd # adapted\n" {
		t.Fatalf("the adapted pack was overwritten: %q", b)
	}
	for _, p := range wantCreated {
		disk, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			t.Fatalf("%s was not written: %v", p, err)
		}
		embedded, _ := packsFS.ReadFile(p)
		if string(disk) != string(embedded) {
			t.Errorf("%s on disk differs from the embedded pack", p)
		}
	}

	// A second run creates nothing and preserves everything.
	criados, preservados, err = SeedPacks(root)
	if err != nil || len(criados) != 0 || len(preservados) != 6 {
		t.Fatalf("re-seed = created %v, preserved %v, err %v", criados, preservados, err)
	}
}

func TestSeedPacks_reportsWriteFailure(t *testing.T) {
	root := t.TempDir()
	// `packs` is a file: the directories cannot be created.
	os.WriteFile(filepath.Join(root, "packs"), []byte("x"), 0o644)
	if _, _, err := SeedPacks(root); err == nil {
		t.Fatal("SeedPacks must report that it could not write the packs")
	}
}
