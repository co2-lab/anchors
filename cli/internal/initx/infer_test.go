package initx

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// writeFixture cria um projeto de brinquedo em dir: uma trinca co-localizada no
// mobile, código no backend, guides, e ruído a ignorar (node_modules).
func writeFixture(t *testing.T, dir string) {
	t.Helper()
	files := map[string]string{
		// trinca co-localizada (screen no mobile)
		"apps/mobile/src/screens/Login.tsx":      "export const Login = () => null // LOGIX-A01",
		"apps/mobile/src/screens/Login.spec.md":  "> **Código**: `LOGIX`\n### LOGIX-A01: entrar",
		"apps/mobile/src/screens/Login.feature":  "@LOGIX-A01\nCenário: entrar",
		"apps/mobile/src/screens/Login.test.tsx": "it('LOGIX-A01: entra', () => {})",
		// mais um alvo, para bater o limiar de co-location (>=3)
		"apps/mobile/src/screens/Home.tsx":      "export const Home = () => null",
		"apps/mobile/src/screens/Home.spec.md":  "> **Código**: `HOMEX`",
		"apps/mobile/src/screens/Home.test.tsx": "it('HOMEX-S01', () => {})",
		"apps/mobile/src/screens/Prof.tsx":      "export const Prof = () => null",
		"apps/mobile/src/screens/Prof.spec.md":  "> **Código**: `PROF`",
		// backend (só código, sem trinca)
		"packages/backend/handlers/auth.ts": "export function auth() {}",
		"packages/backend/repos/user.ts":    "export function getUser() {}",
		// guides
		"guides/FRONTEND_GUIDE.md": "# Frontend",
		"guides/BACKEND_GUIDE.md":  "# Backend",
		// ruído: tem que ser ignorado
		"node_modules/react/index.js": "module.exports = {}",
	}
	// garante massa de código suficiente (codeRoots exige >=10 por dir): enche o backend
	for i := range 12 {
		files["packages/backend/gen/f"+itoa(i)+".ts"] = "export const x = 1"
	}
	for i := range 12 {
		files["apps/mobile/src/extra/g"+itoa(i)+".ts"] = "export const y = 1"
	}
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func itoa(i int) string { return string(rune('0'+i/10)) + string(rune('0'+i%10)) }

func TestInfer(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir)

	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !p.HasSpecMD || !p.HasFeature || !p.HasTest {
		t.Errorf("deveria detectar spec/feature/test; got spec=%v feature=%v test=%v",
			p.HasSpecMD, p.HasFeature, p.HasTest)
	}
	if !p.Colocated {
		t.Error("deveria detectar co-location (há trinca ao lado do código)")
	}
	if p.GuideDir != "guides" {
		t.Errorf("GuideDir = %q, want \"guides\"", p.GuideDir)
	}
	if len(p.GuideFiles) != 2 {
		t.Errorf("esperava 2 guides, got %d (%v)", len(p.GuideFiles), p.GuideFiles)
	}
	if !contains(p.CodeDirs, "apps/mobile") || !contains(p.CodeDirs, "packages/backend") {
		t.Errorf("CodeDirs deveria conter apps/mobile e packages/backend; got %v", p.CodeDirs)
	}
	if !contains(p.CodeExts, ".ts") || !contains(p.CodeExts, ".tsx") {
		t.Errorf("CodeExts deveria conter .ts e .tsx; got %v", p.CodeExts)
	}
}

func TestInfer_ignoresNodeModules(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir)
	p, err := Infer(dir)
	if err != nil {
		t.Fatal(err)
	}
	// node_modules não pode virar um dir de código
	for _, d := range p.CodeDirs {
		if d == "node_modules" || filepath.Base(d) == "react" {
			t.Errorf("node_modules deveria ser ignorado, mas apareceu em CodeDirs: %v", p.CodeDirs)
		}
	}
}

func TestInfer_buildsConfig(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir)
	p, _ := Infer(dir)
	cfg := p.Config

	// as layers de ARTEFATO NÃO são pré-criadas pela inferência — vêm da escolha
	// do usuário no init (ApplyArtifactChoice), que é sempre perguntada.
	for _, name := range []string{"spec", "feature", "test", "guide"} {
		if _, ok := cfg.Layers[name]; ok {
			t.Errorf("a inferência não deveria pré-criar a layer de artefato %q (é escolha do usuário)", name)
		}
	}
	// mas os artefatos DETECTADOS são reportados (para pré-marcar as opções)
	det := p.DetectedArtifacts()
	for _, name := range []string{"spec", "feature", "test", "guide"} {
		if !det[name] {
			t.Errorf("DetectedArtifacts deveria conter %q (está no fixture)", name)
		}
	}
	// há ao menos uma layer de código (essas a inferência propõe como default)
	if len(CodeLayerNames(cfg)) == 0 {
		t.Error("config proposta deveria ter camadas de código")
	}
	// governs começa vazio (é preenchido na P&R)
	if len(cfg.Governs) != 0 {
		t.Errorf("governs deveria começar vazio, got %+v", cfg.Governs)
	}
}

func contains(s []string, v string) bool {
	return slices.Contains(s, v)
}
