package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func cfgComFronteiras(bs ...config.Boundary) *config.Config {
	return &config.Config{
		Layers: map[string]config.Layer{
			"screens":      {Pattern: "src/screens/**/*.ts"},
			"hooks":        {Pattern: "src/hooks/**/*.ts"},
			"repositories": {Pattern: "src/repositories/**/*.ts"},
		},
		Boundaries: bs,
	}
}

func rodaFronteira(t *testing.T, arquivo, conteúdo string, cfg *config.Config) (Verdict, string) {
	t.Helper()
	return checkLayerBoundary(conteúdo, mapx.Node{ID: arquivo, Kind: mapx.KindCode}, "", nil, cfg)
}

// O caso que motivou o gate: a tela falando direto com o repositório, pulando o hook.
func TestLayerBoundaryViolacao(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{
		Layer: "screens", Forbid: `from '@/repositories`,
		Because: "tela fala com hook, não com dado",
	})
	código := "import { useState } from 'react'\nimport { getUser } from '@/repositories/user'\n"

	v, d := rodaFronteira(t, "src/screens/Home.ts", código, cfg)
	if v != Fail {
		t.Fatalf("violação deveria reprovar, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "linha 2") {
		t.Errorf("não apontou ONDE: %s", d)
	}
	if !strings.Contains(d, "tela fala com hook") {
		t.Errorf("não disse o PORQUÊ (uma proibição sem motivo vira ritual): %s", d)
	}
}

// A regra vale para a camada declarada, e só. O mesmo import num hook é legítimo.
func TestLayerBoundaryEscopadaNaCamada(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})
	código := "import { getUser } from '@/repositories/user'\n"

	if v, d := rodaFronteira(t, "src/hooks/useUser.ts", código, cfg); v == Fail {
		t.Fatalf("hook PODE importar repositório — a regra é de screens (%s)", d)
	}
}

// Regra sem `layer` vale para todo código: é assim que se declara uma proibição global
// (relógio cru, cor literal, console.log…).
func TestLayerBoundaryGlobal(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{
		Forbid: `new Date\(\)|Date\.now\(\)`, Because: "relógio vem do util (testabilidade)",
	})
	for _, arquivo := range []string{"src/screens/Home.ts", "src/hooks/useX.ts", "src/repositories/user.ts"} {
		t.Run(arquivo, func(t *testing.T) {
			v, _ := rodaFronteira(t, arquivo, "const agora = new Date()\n", cfg)
			if v != Fail {
				t.Fatalf("proibição global deveria valer aqui também, foi %s", v)
			}
		})
	}
}

// `severity: warn` é a maturação POR REGRA: trava a fronteira nova sem desligar o gate
// por causa do backlog da antiga.
func TestLayerBoundarySeveridade(t *testing.T) {
	código := "import { x } from '@/legacy/thing'\n"

	warn := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `@/legacy`, Severity: "warn"})
	if v, d := rodaFronteira(t, "src/screens/Home.ts", código, warn); v != Pending {
		t.Fatalf("regra em migração registra sem reprovar, foi %s (%s)", v, d)
	}

	erro := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `@/legacy`})
	if v, _ := rodaFronteira(t, "src/screens/Home.ts", código, erro); v != Fail {
		t.Fatalf("severidade default é error, foi %s", v)
	}
}

// Opt-out honesto: vale com razão (na linha OU no comentário acima), não vale nu.
func TestLayerBoundaryDispensa(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})

	inline := "import { getUser } from '@/repositories/user' // @allow-boundary: migração da tela pendente, ticket app de referência-412\n"
	if v, d := rodaFronteira(t, "src/screens/Home.ts", inline, cfg); v != Pass {
		t.Errorf("dispensa inline COM razão deveria passar, foi %s (%s)", v, d)
	}

	acima := "// @allow-boundary: migração da tela pendente, ticket app de referência-412\nimport { getUser } from '@/repositories/user'\n"
	if v, d := rodaFronteira(t, "src/screens/Home.ts", acima, cfg); v != Pass {
		t.Errorf("dispensa na linha ACIMA deveria passar (import não tem onde receber comentário legível), foi %s (%s)", v, d)
	}

	nu := "import { getUser } from '@/repositories/user' // @allow-boundary:\n"
	if v, _ := rodaFronteira(t, "src/screens/Home.ts", nu, cfg); v != Fail {
		t.Errorf("marcador NU deveria continuar reprovando, foi %s", v)
	}
}

// Sem fronteiras declaradas o gate NÃO passa: ele não verificou nada, e um ✓ seria
// mentira. Pendente diz o que declarar.
func TestLayerBoundarySemDeclaracaoEhPendente(t *testing.T) {
	for nome, cfg := range map[string]*config.Config{
		"config nula":       nil,
		"sem boundaries":    {},
		"boundaries vazias": {Boundaries: []config.Boundary{}},
	} {
		t.Run(nome, func(t *testing.T) {
			v, d := rodaFronteira(t, "src/screens/Home.ts", "qualquer coisa", cfg)
			if v != Pending {
				t.Fatalf("veredito = %s, queria Pending (%s)", v, d)
			}
			if !strings.Contains(d, "boundaries") {
				t.Errorf("Pendente não diz o que declarar: %s", d)
			}
		})
	}
}

// Regex inválido é erro de CONFIG e precisa APARECER: ignorado em silêncio, desligaria a
// regra sem ninguém saber — o pior desfecho possível para um gate.
func TestLayerBoundaryRegexInvalidoNaoSilencia(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `[inválido(`})
	v, d := rodaFronteira(t, "src/screens/Home.ts", "qualquer coisa", cfg)
	if v != Fail {
		t.Fatalf("padrão inválido deveria falhar visivelmente, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "inválido") {
		t.Errorf("não explicou o problema de config: %s", d)
	}
}

// A PROVA do agnosticismo: as mesmas regras arquiteturais, expressas no dialeto de import
// de cada linguagem. O engine não conhece nenhum deles — quem escreve o padrão é o projeto.
func TestLayerBoundaryAgnosticoEntreLinguagens(t *testing.T) {
	casos := []struct{ nome, forbid, viola, ok string }{
		{"TypeScript", `from ['"]@/repositories`, `import { x } from '@/repositories/user'`, `import { x } from '@/hooks/useUser'`},
		{"Python", `^\s*from\s+app\.repositories`, `from app.repositories.user import get`, `from app.hooks.user import use`},
		{"Go", `"myapp/internal/repositories"`, `	"myapp/internal/repositories"`, `	"myapp/internal/hooks"`},
		{"Java", `^import\s+com\.app\.repositories`, `import com.app.repositories.UserRepo;`, `import com.app.hooks.UserHook;`},
		{"Rust", `^\s*use\s+crate::repositories`, `use crate::repositories::user;`, `use crate::hooks::user;`},
		{"Ruby", `require ['"]repositories/`, `require 'repositories/user'`, `require 'hooks/user'`},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: c.forbid})
			if v, d := rodaFronteira(t, "src/screens/Home.ts", c.viola+"\n", cfg); v != Fail {
				t.Errorf("não pegou a violação em %s: %s", c.nome, d)
			}
			if v, d := rodaFronteira(t, "src/screens/Home.ts", c.ok+"\n", cfg); v != Pass {
				t.Errorf("falso positivo no import legítimo em %s: %s", c.nome, d)
			}
		})
	}
}

// O furo que o gate tinha: casava LINHA A LINHA, então todo padrão cujo alvo se espalha
// por várias linhas escapava. Medido no app de referência: das 7 telas que importavam `Modal` do
// react-native, o gate acusava UMA — a única com o import numa linha só. As outras 6
// quebram em várias porque passam de 100 colunas (Prettier), e ficavam invisíveis.
func TestLayerBoundaryPegaImportMultilinha(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{
		Layer:   "screens",
		Forbid:  `import\s*\{[^}]*\bModal\b[^}]*\}\s*from ['"]react-native['"]`,
		Because: "sheet/modal usa AppBottomSheet",
	})
	código := "import {\n  View,\n  Modal,\n  Text,\n} from 'react-native'\n"

	v, d := rodaFronteira(t, "src/screens/Membro.ts", código, cfg)
	if v != Fail {
		t.Fatalf("import multilinha deveria reprovar, foi %s (%s)", v, d)
	}
	// Aponta a linha onde o casamento COMEÇA — o `import {`.
	if !strings.Contains(d, "linha 1") {
		t.Errorf("não apontou o início do import: %s", d)
	}
}

// A contrapartida: o mesmo import de UMA linha continua sendo pego (não é regressão).
func TestLayerBoundaryPegaImportDeUmaLinha(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{
		Layer:  "screens",
		Forbid: `import\s*\{[^}]*\bModal\b[^}]*\}\s*from ['"]react-native['"]`,
	})
	código := "import { View, Modal, Text } from 'react-native'\n"

	if v, d := rodaFronteira(t, "src/screens/Membro.ts", código, cfg); v != Fail {
		t.Fatalf("import de uma linha deveria reprovar, foi %s (%s)", v, d)
	}
}

// A dispensa num import multilinha: a marcação natural fica na linha do `from`, não na do
// `import` — cobrar que ela esteja na primeira linha do casamento seria exigir que o autor
// soubesse onde o regex começou a casar. É o arranjo real do WelcomeModal.tsx no app de referência.
func TestLayerBoundaryDispensaEmQualquerLinhaDoTrecho(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{
		Layer:  "screens",
		Forbid: `import\s*\{[^}]*\bModal\b[^}]*\}\s*from ['"]react-native['"]`,
	})
	código := "import {\n  View,\n  Modal,\n} from 'react-native' // @allow-boundary: overlay de tela cheia, não um sheet\n"

	if v, d := rodaFronteira(t, "src/screens/Welcome.ts", código, cfg); v != Pass {
		t.Fatalf("dispensa na linha do `from` deveria passar, foi %s (%s)", v, d)
	}
}

// Um projeto que QUEIRA ancorar o padrão numa linha só continua podendo: `^`/`$` seguem
// valendo por linha, porque `(?s)` muda o `.`, não o significado das âncoras.
func TestLayerBoundaryAncoraDeLinhaSegueValendo(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `(?m)^import 'proibido'$`})
	código := "const x = \"import 'proibido'\"\n"

	if v, d := rodaFronteira(t, "src/screens/Home.ts", código, cfg); v != Pass {
		t.Fatalf("âncora de linha não deveria casar no meio da linha, foi %s (%s)", v, d)
	}
}
