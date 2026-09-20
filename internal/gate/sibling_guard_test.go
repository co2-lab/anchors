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
		t.Run("SBGRD-B01: An artifact that is not code leaves without a verdict", func(t *testing.T) {})
		v, _ := checkSiblingGuard("export function f() {}", mapx.Node{Kind: mapx.KindSpec}, "", nil, cfg)
		if v != Skip {
			t.Errorf("esperava Skip, veio %v", v)
		}
	})

	t.Run("pendente quando sem dialeto", func(t *testing.T) {
		t.Run("SBGRD-B02: Without a declared dialect the verdict is undetermined", func(t *testing.T) {})
		v, msg := checkSiblingGuard("export function f() {}", mapx.Node{Kind: mapx.KindCode}, "", nil, &config.Config{})
		if v != Pending {
			t.Errorf("esperava Pending, veio %v (%s)", v, msg)
		}
	})

	t.Run("pula com menos de 3 funções", func(t *testing.T) {
		t.Run("SBGRD-B03: Fewer than three siblings on the same parameter is left alone", func(t *testing.T) {})
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
		t.Run("SBGRD-B04: The sibling that does not guard is accused", func(t *testing.T) {})
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
		t.Run("SBGRD-B08: A waiver with a written reason silences the accusation", func(t *testing.T) {})
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
		t.Run("SBGRD-X01: The gate does not invent what an exported function looks like", func(t *testing.T) {})
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

// Os casos que faltavam: as duas pontas da simetria, a citação do nome, a guarda cujo
// CONTEÚDO difere, e a função sozinha.
func TestSiblingGuard_CasosDaSimetria(t *testing.T) {
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts"}}
	roda := func(src string) (Verdict, string) {
		return checkSiblingGuard(src, mapx.Node{Kind: mapx.KindCode}, "", nil, cfg)
	}

	t.Run("SBGRD-B05: When every sibling guards, nothing is accused", func(t *testing.T) {
		src := "\nexport function save(key: string) {\n\tif (!key) throw new Error(\"x\")\n\treturn key\n}\n" +
			"export function update(key: string) {\n\tif (!key) throw new Error(\"x\")\n\treturn key\n}\n" +
			"export function get(key: string) {\n\tif (!key) throw new Error(\"x\")\n\treturn key\n}\n"
		if v, msg := roda(src); v == Fail {
			t.Errorf("todas guardam — não há assimetria a reportar: %v (%s)", v, msg)
		}
	})

	t.Run("SBGRD-B06: When no sibling guards, nothing is accused either", func(t *testing.T) {
		src := "\nexport function save(key: string) { return key }\n" +
			"export function update(key: string) { return key }\n" +
			"export function get(key: string) { return key }\n"
		if v, msg := roda(src); v == Fail {
			t.Errorf("nenhuma guarda é DECISÃO, não esquecimento: %v (%s)", v, msg)
		}
	})

	t.Run("SBGRD-B07: The verdict names the function and the parameter", func(t *testing.T) {
		src := "\nexport function save(chave: string) {\n\tif (!chave) throw new Error(\"x\")\n\treturn chave\n}\n" +
			"export function update(chave: string) {\n\tif (!chave) throw new Error(\"x\")\n\treturn chave\n}\n" +
			"export function get(chave: string) { return chave }\n"
		v, msg := roda(src)
		if v != Fail {
			t.Fatalf("assimetria tem de reprovar: %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "get") {
			t.Errorf("o veredito não nomeia a função sem guarda: %s", msg)
		}
		if !strings.Contains(msg, "chave") {
			t.Errorf("o veredito não nomeia o parâmetro em jogo: %s", msg)
		}
	})

	// O gate não entende o que a guarda FAZ. Aqui cada irmã guarda de um jeito
	// diferente — e é justamente por não julgar o conteúdo que ele consegue medir
	// isto sem conhecer o domínio.
	t.Run("SBGRD-I01: The gate never judges what the guard does", func(t *testing.T) {
		src := "\nexport function save(key: string) {\n\tif (!key) throw new Error(\"vazio\")\n\treturn key\n}\n" +
			"export function update(key: string) {\n\tif (key.length > 64) throw new Error(\"longo\")\n\treturn key\n}\n" +
			"export function get(key: string) { return key }\n"
		v, msg := roda(src)
		if v != Fail {
			t.Fatalf("guardas de conteúdo DIFERENTE ainda são guardas: %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "get") {
			t.Errorf("só a AUSÊNCIA devia ser reportada: %s", msg)
		}
	})

	t.Run("SBGRD-X02: A single function in isolation is not accused", func(t *testing.T) {
		if v, msg := roda("\nexport function get(key: string) { return key }\n"); v == Fail {
			t.Errorf("sem irmãs não há assimetria — e a assimetria é a prova toda: %v (%s)", v, msg)
		}
	})

	// As três condições valem JUNTAS: cada módulo abaixo satisfaz só parte delas.
	t.Run("SBGRD-I02: The three conservatism conditions hold together", func(t *testing.T) {
		parciais := []string{
			// 3 irmãs, mas nenhuma guarda (falta "a maioria guarda")
			"\nexport function a(k: string) { return k }\nexport function b(k: string) { return k }\nexport function c(k: string) { return k }\n",
			// 3 funções no arquivo, mas só 2 compartilham o parâmetro `k` — o piso de três
			// vale por PARÂMETRO, não por arquivo. Sem isso, a terceira função (que usa
			// outro parâmetro) faria um par virar "maioria".
			"\nexport function a(k: string) {\n\tif (!k) throw new Error(\"x\")\n\treturn k\n}\nexport function b(k: string) { return k }\nexport function c(outro: string) { return outro }\n",
			// 3 irmãs e todas guardam (falta "ao menos uma não guarda")
			"\nexport function a(k: string) {\n\tif (!k) throw new Error(\"x\")\n\treturn k\n}\nexport function b(k: string) {\n\tif (!k) throw new Error(\"x\")\n\treturn k\n}\nexport function c(k: string) {\n\tif (!k) throw new Error(\"x\")\n\treturn k\n}\n",
		}
		for i, src := range parciais {
			if v, msg := roda(src); v == Fail {
				t.Errorf("caso %d satisfaz só parte das condições e foi acusado: %v (%s)", i, v, msg)
			}
		}
	})
}
