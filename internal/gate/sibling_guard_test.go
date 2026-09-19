package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestSiblingGuard(t *testing.T) {
	cfg := &config.Config{
		Dialect: &config.Dialect{Family: "ts"},
	}

	t.Run("pula quando não é código", func(t *testing.T) {
		v, _ := checkSiblingGuard("export function f() {}", mapx.Node{Kind: mapx.KindSpec}, "", nil, cfg)
		if v != Skip {
			t.Errorf("esperava Skip, veio %v", v)
		}
	})

	t.Run("pendente quando sem dialeto", func(t *testing.T) {
		v, msg := checkSiblingGuard("export function f() {}", mapx.Node{Kind: mapx.KindCode}, "", nil, &config.Config{})
		if v != Pending {
			t.Errorf("esperava Pending, veio %v (%s)", v, msg)
		}
	})

	t.Run("pula com menos de 3 funções", func(t *testing.T) {
		src := `
export function a(key: string) { return key }
export function b(key: string) { return key }
`
		v, _ := checkSiblingGuard(src, mapx.Node{Kind: mapx.KindCode}, "", nil, cfg)
		if v != Skip {
			t.Errorf("esperava Skip, veio %v", v)
		}
	})

	t.Run("assimetria reprova", func(t *testing.T) {
		src := `
export function save(key: string) {
	if (!key) throw new Error("empty key")
	return key
}
export function update(key: string) {
	if (!key) throw new Error("empty key")
	return key
}
export function get(key: string) {
	return key
}
`
		v, msg := checkSiblingGuard(src, mapx.Node{Kind: mapx.KindCode}, "", nil, cfg)
		if v != Fail {
			t.Fatalf("esperava Fail, veio %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "get") {
			t.Errorf("mensagem deve citar a função sem guarda: %s", msg)
		}
		if !strings.Contains(msg, "save") || !strings.Contains(msg, "update") {
			t.Errorf("mensagem deve citar as funções com guarda: %s", msg)
		}
	})

	t.Run("dispensa com motivo passa", func(t *testing.T) {
		src := `
export function save(key: string) {
	if (!key) throw new Error("empty key")
	return key
}
export function update(key: string) {
	if (!key) throw new Error("empty key")
	return key
}
export function get(key: string) {
	// @no-guard: delega a validação ao storage subjacente
	return key
}
`
		v, msg := checkSiblingGuard(src, mapx.Node{Kind: mapx.KindCode}, "", nil, cfg)
		if v != Pass {
			t.Errorf("esperava Pass com waiver, veio %v (%s)", v, msg)
		}
	})

	t.Run("reconhece guarda em Go sem parenteses", func(t *testing.T) {
		goCfg := &config.Config{Dialect: &config.Dialect{Family: "go"}}
		src := `
func Save(key string) error {
	if key == "" { return errors.New("empty") }
	return nil
}
func Update(key string) error {
	if key == "" { return errors.New("empty") }
	return nil
}
func Get(key string) string {
	return key
}
`
		v, msg := checkSiblingGuard(src, mapx.Node{Kind: mapx.KindCode}, "", nil, goCfg)
		if v != Fail {
			t.Fatalf("esperava Fail por assimetria em Go, veio %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "Get") {
			t.Errorf("deve acusar Get em Go: %s", msg)
		}
	})

	t.Run("usa padrao de guarda customizado", func(t *testing.T) {
		customCfg := &config.Config{
			Dialect: &config.Dialect{
				Family: "go",
				GuardPatterns: []string{
					`validate\({{param}}\)`,
				},
			},
		}
		src := `
func Save(key string) {
	validate(key)
}
func Update(key string) {
	validate(key)
}
func Get(key string) {
}
`
		v, msg := checkSiblingGuard(src, mapx.Node{Kind: mapx.KindCode}, "", nil, customCfg)
		if v != Fail {
			t.Fatalf("esperava Fail com guard_patterns customizado, veio %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "Get") {
			t.Errorf("deve acusar Get: %s", msg)
		}
	})
}
