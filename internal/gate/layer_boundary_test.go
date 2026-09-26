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
	t.Run("LYBNL-B02: Content matching a forbidden pattern fails, naming line and reason", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{
		Layer: "screens", Forbid: `from '@/repositories`,
		Because: "tela fala com hook, não com dado",
	})
	código := "import { useState } from 'react'\nimport { getUser } from '@/repositories/user'\n"

	v, d := rodaFronteira(t, "src/screens/Home.ts", código, cfg)
	if v != Fail {
		t.Fatalf("violação deveria reprovar, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "linha 2") && !strings.Contains(d, "line 2") {
		t.Errorf("não apontou ONDE: %s", d)
	}
	if !strings.Contains(d, "tela fala com hook") {
		t.Errorf("não disse o PORQUÊ (uma proibição sem motivo vira ritual): %s", d)
	}
}

// A regra vale para a camada declarada, e só. O mesmo import num hook é legítimo.
func TestLayerBoundaryEscopadaNaCamada(t *testing.T) {
	t.Run("LYBNL-B03: A rule scoped to a layer charges only that layer", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})
	código := "import { getUser } from '@/repositories/user'\n"

	if v, d := rodaFronteira(t, "src/hooks/useUser.ts", código, cfg); v == Fail {
		t.Fatalf("hook PODE importar repositório — a regra é de screens (%s)", d)
	}
}

// Regra sem `layer` vale para todo código: é assim que se declara uma proibição global
// (relógio cru, cor literal, console.log…).
func TestLayerBoundaryGlobal(t *testing.T) {
	t.Run("LYBNL-B04: A rule with no layer holds for all code", func(t *testing.T) {})
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
	t.Run("LYBNL-B05: Severity warn records without failing, and the default is error", func(t *testing.T) {})
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

// Opt-out honesto: a dispensa na própria linha, COM razão escrita, vale.
func TestLayerBoundaryDispensaInline(t *testing.T) {
	t.Run("LYBNL-B06: A waiver with a written reason on the line waives that line", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})

	inline := "import { getUser } from '@/repositories/user' // @allow-boundary: migração da tela pendente, ticket app de referência-412\n"
	if v, d := rodaFronteira(t, "src/screens/Home.ts", inline, cfg); v != Pass {
		t.Errorf("dispensa inline COM razão deveria passar, foi %s (%s)", v, d)
	}
}

// O import não tem onde receber um comentário de fim de linha legível em várias
// linguagens: obrigar a marcação inline empurraria o autor a não marcar.
func TestLayerBoundaryDispensaNaLinhaAcima(t *testing.T) {
	t.Run("LYBNL-B07: The waiver also holds in the comment on the line above", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})

	acima := "// @allow-boundary: migração da tela pendente, ticket app de referência-412\nimport { getUser } from '@/repositories/user'\n"
	if v, d := rodaFronteira(t, "src/screens/Home.ts", acima, cfg); v != Pass {
		t.Errorf("dispensa na linha ACIMA deveria passar, foi %s (%s)", v, d)
	}
}

// Marcador NU não dispensa — senão vira um jeito silencioso de calar o gate, e some o
// rastro de que houve decisão.
func TestLayerBoundaryDispensaNuaNaoVale(t *testing.T) {
	t.Run("LYBNL-B08: A bare marker with no reason does not waive", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})

	nu := "import { getUser } from '@/repositories/user' // @allow-boundary:\n"
	if v, _ := rodaFronteira(t, "src/screens/Home.ts", nu, cfg); v != Fail {
		t.Errorf("marcador NU deveria continuar reprovando, foi %s", v)
	}
}

// Sem fronteiras declaradas o gate NÃO passa: ele não verificou nada, e um ✓ seria
// mentira. Pendente diz o que declarar.
func TestLayerBoundarySemDeclaracaoEhPendente(t *testing.T) {
	t.Run("LYBNL-B09: With no boundary declared the verdict is Pending, never Pass", func(t *testing.T) {})
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
	t.Run("LYBNL-B10: An invalid forbid pattern fails visibly", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `[inválido(`})
	v, d := rodaFronteira(t, "src/screens/Home.ts", "qualquer coisa", cfg)
	if v != Fail {
		t.Fatalf("padrão inválido deveria falhar visivelmente, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "inválido") && !strings.Contains(d, "invalid") {
		t.Errorf("não explicou o problema de config: %s", d)
	}
}

// A PROVA do agnosticismo: as mesmas regras arquiteturais, expressas no dialeto de import
// de cada linguagem. O engine não conhece nenhum deles — quem escreve o padrão é o projeto.
func TestLayerBoundaryAgnosticoEntreLinguagens(t *testing.T) {
	t.Run("LYBNL-I01: The same rule is expressible in six language dialects", func(t *testing.T) {})
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
	t.Run("LYBNL-B11: The pattern is matched against the whole file, catching a multi-line import", func(t *testing.T) {})
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
	if !strings.Contains(d, "linha 1") && !strings.Contains(d, "line 1") {
		t.Errorf("não apontou o início do import: %s", d)
	}
}

// A contrapartida: o mesmo import de UMA linha continua sendo pego (não é regressão).
func TestLayerBoundaryPegaImportDeUmaLinha(t *testing.T) {
	t.Run("LYBNL-I02: A single-line import of the same shape is still caught", func(t *testing.T) {})
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
	t.Run("LYBNL-I03: The waiver holds on any line of the matched stretch", func(t *testing.T) {})
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
	t.Run("LYBNL-I04: A line anchor keeps holding per line", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `(?m)^import 'proibido'$`})
	código := "const x = \"import 'proibido'\"\n"

	if v, d := rodaFronteira(t, "src/screens/Home.ts", código, cfg); v != Pass {
		t.Fatalf("âncora de linha não deveria casar no meio da linha, foi %s (%s)", v, d)
	}
}

// O gate é de CÓDIGO: uma spec ou uma feature não tem import a proibir, e confrontá-las
// acusaria o texto que fala sobre o padrão.
func TestLayerBoundarySoConfrontaCodigo(t *testing.T) {
	t.Run("LYBNL-B01: An artifact that is not code leaves without a verdict", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})
	conteudo := "A spec fala do import `from '@/repositories/user'` como o que NÃO se faz.\n"

	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindFeature, mapx.KindTest} {
		v, d := checkLayerBoundary(conteudo, mapx.Node{ID: "src/screens/Home.spec.md", Kind: k}, "", nil, cfg)
		if v != Skip {
			t.Errorf("kind %s deveria sair sem veredito, foi %s (%s)", k, v, d)
		}
	}
}

// O gate NÃO inventa fronteira: sem declaração ele não reprova a violação que qualquer
// revisor apontaria. A arquitetura é decisão do projeto.
func TestLayerBoundaryNaoInventaFronteira(t *testing.T) {
	t.Run("LYBNL-X01: The gate does not decide which boundaries exist", func(t *testing.T) {})
	violacaoObvia := "import { getUser } from '@/repositories/user'\n"

	if v, d := rodaFronteira(t, "src/screens/Home.ts", violacaoObvia, &config.Config{}); v == Fail {
		t.Errorf("sem declaração o gate não pode reprovar — inventaria dever: %s", d)
	}
}

// O gate casa TEXTO, não parseia import: o padrão do projeto pega o termo onde quer que
// ele apareça. É o que o torna agnóstico de linguagem.
func TestLayerBoundaryCasaTextoNaoImport(t *testing.T) {
	t.Run("LYBNL-X02: The gate does not parse the language, it matches text", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `@/repositories`})
	// O termo NÃO está num import — está numa string de log.
	codigo := "const origem = 'veio de @/repositories/user'\n"

	if v, d := rodaFronteira(t, "src/screens/Home.ts", codigo, cfg); v != Fail {
		t.Fatalf("o gate casa TEXTO: entender o grafo de import de cada linguagem prenderia "+
			"o engine a um conjunto de ecossistemas. Foi %s (%s)", v, d)
	}
}

// O gate não julga se a fronteira VALE a pena: a régua é a DECLARAÇÃO. Uma proibição que
// um revisor acharia inofensiva reprova do mesmo jeito.
func TestLayerBoundaryNaoJulgaSeAFronteiraFazSentido(t *testing.T) {
	t.Run("LYBNL-X03: The gate does not judge whether the boundary is the right one to draw", func(t *testing.T) {})
	cfg := cfgComFronteiras(config.Boundary{
		Layer: "screens", Forbid: `useState`, Because: "decisão de design deste projeto",
	})

	if v, d := rodaFronteira(t, "src/screens/Home.ts", "import { useState } from 'react'\n", cfg); v != Fail {
		t.Fatalf("o julgamento é de quem escreve a Estrutura, não do gate: %s (%s)", v, d)
	}
}

func TestLayerBoundary_findingDetails(t *testing.T) {
	t.Run("LYBNL-B02: Content matching a forbidden pattern fails, naming line and reason", func(t *testing.T) {
		cfg := cfgComFronteiras(config.Boundary{Layer: "screens", Forbid: `from '@/repositories`})
		content := "import a from '@/repositories/a'\nconst x = 1\nimport b from '@/repositories/b'\n"
		v, d := rodaFronteira(t, "src/screens/Home.ts", content, cfg)
		if v != Fail || !strings.Contains(d, "line 1:") || !strings.Contains(d, "line 3:") {
			t.Errorf("every occurrence is named by its line: %v (%s)", v, d)
		}
		if !strings.Contains(d, "layer `screens` cannot contain") {
			t.Errorf("a scoped rule names its layer: %s", d)
		}
		if strings.Contains(d, "warning(s)") {
			t.Errorf("no warning was found, none is announced: %s", d)
		}
	})

	t.Run("LYBNL-B04: A rule with no layer holds for all code", func(t *testing.T) {
		cfg := cfgComFronteiras(config.Boundary{Forbid: `Date\.now`})
		_, d := rodaFronteira(t, "src/hooks/x.ts", "Date.now()\n", cfg)
		if !strings.Contains(d, "this layer cannot contain") {
			t.Errorf("a global rule speaks of this layer: %s", d)
		}
	})

	t.Run("LYBNL-B05: Severity warn records without failing, and the default is error", func(t *testing.T) {
		cfg := cfgComFronteiras(
			config.Boundary{Forbid: `console\.log`},
			config.Boundary{Forbid: `debugger`, Severity: "warn"},
		)
		v, d := rodaFronteira(t, "src/hooks/x.ts", "console.log(1)\ndebugger\n", cfg)
		if v != Fail || !strings.Contains(d, "and 1 more warning(s)") {
			t.Errorf("the failure still counts the warnings beside it: %v (%s)", v, d)
		}
	})

	t.Run("LYBNL-B06: A waiver with a written reason on the line waives that line", func(t *testing.T) {
		// the match ends at the end of line 1; the waiver on line 2 belongs to line 2 only
		cfg := cfgComFronteiras(config.Boundary{Forbid: `Date\.now\(\)`})
		content := "const t = Date.now()\nconst u = 2 // @allow-boundary: unrelated reason\n"
		if v, d := rodaFronteira(t, "src/hooks/x.ts", content, cfg); v != Fail {
			t.Errorf("a waiver on the line BELOW does not waive: %v (%s)", v, d)
		}
	})
}

func TestLayerBoundary_listsAtMostFiveFindings(t *testing.T) {
	cfg := cfgComFronteiras(config.Boundary{Forbid: `console\.log`})
	five := strings.Repeat("console.log(1)\n", 5)
	if _, d := rodaFronteira(t, "src/hooks/x.ts", five, cfg); strings.Contains(d, "(and ") {
		t.Errorf("five findings are all listed, with no remainder: %s", d)
	}
	seven := strings.Repeat("console.log(1)\n", 7)
	_, d := rodaFronteira(t, "src/hooks/x.ts", seven, cfg)
	if !strings.Contains(d, "(and 2 more)") || strings.Contains(d, "line 6:") {
		t.Errorf("seven findings list five and count the other two: %s", d)
	}
}
