// @anchors
//   ref: ATGDT

package governance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/settings"
)

// THE GUIDE CHANGES with the local declaration, and the difference is the point.
//
// `settings user-issues` closes the claim's door: whoever does not decide the product gets
// no escalated card. But the door used most is another, and it has no label — the agent
// ASKING the dev who is running it. The dev knows the code and will answer; the answer is
// reasonable, and becomes a product decision taken by someone with no authority.
func TestAutonomyGuideChangesWithTheDeclaration(t *testing.T) {
	t.Run("ATGDT-B01: A role that decides the product is told to record its decisions", func(t *testing.T) {})
	t.Run("ATGDT-X01: A role that decides the product is not forbidden to ask", func(t *testing.T) {})
	withoutAuthority := t.TempDir()
	if err := settings.Save(withoutAuthority, settings.Settings{Role: "dev", DecidedAt: "2026-09-26"}); err != nil {
		t.Fatal(err)
	}
	withAuthority := t.TempDir()
	if err := settings.Save(withAuthority, settings.Settings{Role: "product-owner", DecidedAt: "2026-09-26"}); err != nil {
		t.Fatal(err)
	}

	text := autonomyGuide(withoutAuthority)
	if !strings.Contains(text, "Do not ask whoever is running you.") {
		t.Errorf("whoever does not decide the product must read the explicit ban:\n%s", text)
	}
	if !strings.Contains(text, "--for-user") {
		t.Error("the instruction does not say what to do instead of asking")
	}
	// the reason matters more than the ban: a rule without a why is the first to be dodged
	if !strings.Contains(text, "authority") {
		t.Error("the instruction bans without saying why")
	}

	other := autonomyGuide(withAuthority)
	if strings.Contains(other, "Do not ask whoever is running you.") {
		t.Error("whoever decides the product does not get the ban")
	}
	if !strings.Contains(other, "Your role (Product Owner) decides the direction of this product.") {
		t.Errorf("the deciding role should be told it decides:\n%s", other)
	}
	if !strings.Contains(other, "is written, not") || !strings.Contains(other, "--for-user") {
		t.Error("even whoever decides must read that the decision is RECORDED, through an escalation")
	}
}

// WITHOUT A DECLARATION, the ruler is the closed one — the same default as the claim, for
// the same reason: erring on the open side lets someone decide the product with no
// authority, and that is invisible after the fact.
// A declaration by the old `user_issues` flag has no role: the guide printed "Your role ()".
func TestAutonomyGuideLegacyFlagNamesNoEmptyRole(t *testing.T) {
	t.Run("ATGDT-B07: The old user_issues flag with no role is named as such, never as an empty role", func(t *testing.T) {})
	yes, no := true, false
	for _, c := range []struct {
		flag *bool
		want string
		ban  bool
	}{
		{&yes, "Your declaration (the old `user_issues` flag, with no role yet — declare one with `anchors settings role`) says you decide the direction of this product.", false},
		{&no, "**Your declaration (the old `user_issues` flag, with no role yet) says you do NOT decide the direction\nof this product**", true},
	} {
		root := t.TempDir()
		if err := settings.Save(root, settings.Settings{UserIssues: c.flag, DecidedAt: "2026-09-26"}); err != nil {
			t.Fatal(err)
		}
		text := autonomyGuide(root)
		if strings.Contains(text, "()") || strings.Contains(text, "Your role") {
			t.Errorf("user_issues=%v: no role was declared, none may be named:\n%s", *c.flag, text)
		}
		if !strings.Contains(text, c.want) {
			t.Errorf("user_issues=%v: the guide must say what was declared (%q):\n%s", *c.flag, c.want, text)
		}
		if got := strings.Contains(text, "Do not ask whoever is running you."); got != c.ban {
			t.Errorf("user_issues=%v: the ban on asking is %v, want %v", *c.flag, got, c.ban)
		}
	}
}

func TestAutonomyGuideWithoutDeclarationIsClosed(t *testing.T) {
	t.Run("ATGDT-B03: With no role declared the guide is the closed one", func(t *testing.T) {})
	text := autonomyGuide(t.TempDir())
	if !strings.Contains(text, "Do not ask whoever is running you.") {
		t.Errorf("with no declaration the default has to be the closed one:\n%s", text)
	}
	// and the text cannot CLAIM a declaration that did not happen: whoever never declared
	// would look in `.anchors/settings.yaml` for a record that is not there
	if strings.Contains(text, "Your role (") {
		t.Error("the text claims a declaration that did not happen")
	}
	if !strings.Contains(text, "**You did not declare a role**") {
		t.Error("the text does not say that no role was declared")
	}
}

// A declared role that does not decide reads who does, and the three instructions every
// non-decider needs: do not ask, move on after escalating, and what not to escalate.
func TestAutonomyGuideForADeclaredNonDecider(t *testing.T) {
	t.Run("ATGDT-B02: A declared role that does not decide is told who decides", func(t *testing.T) {})
	t.Run("ATGDT-B06: Whoever does not decide is told not to ask, to move on and what not to escalate", func(t *testing.T) {})
	dir := t.TempDir()
	if err := settings.Save(dir, settings.Settings{Role: "dev", DecidedAt: "2026-09-26"}); err != nil {
		t.Fatal(err)
	}
	for _, g := range []string{autonomyGuide(dir), autonomyGuide(t.TempDir())} {
		for _, want := range []string{"Do not ask whoever is running you.", "move on to the next card", "### What is NOT to be escalated"} {
			if !strings.Contains(g, want) {
				t.Errorf("a non-decider should read %q:\n%s", want, g)
			}
		}
	}
	g := autonomyGuide(dir)
	if !strings.Contains(g, "**Your role (Dev) does NOT decide the direction of this product**") ||
		!strings.Contains(g, "`product-owner` or the `architect`") {
		t.Errorf("a declared non-decider should read that its role does not decide, and who does:\n%s", g)
	}
}

// "Reviewing" is not one thing: whoever hunts a data leak and whoever hunts a query in a
// loop read the same code with different questions. A role with a lens reads it.
func TestAutonomyGuideShowsTheRoleLens(t *testing.T) {
	t.Run("ATGDT-B05: A role with a lens reads it", func(t *testing.T) {})
	qa := t.TempDir()
	if err := settings.Save(qa, settings.Settings{Role: "qa", DecidedAt: "2026-09-26"}); err != nil {
		t.Fatal(err)
	}
	if g := autonomyGuide(qa); !strings.Contains(g, "## This role's lens (QA)\n\n"+settings.Role("qa").Lens()+".") {
		t.Errorf("the qa role should read its lens:\n%s", g)
	}
	dev := t.TempDir()
	if err := settings.Save(dev, settings.Settings{Role: "dev", DecidedAt: "2026-09-26"}); err != nil {
		t.Fatal(err)
	}
	if g := autonomyGuide(dev); strings.Contains(g, "This role's lens") {
		t.Errorf("a role with no lens gets no lens section:\n%s", g)
	}
}

// A settings file that cannot be read must never turn into "you may ask": that would be
// the most expensive silent failure of the mechanism.
func TestAutonomyGuideUnreadableSettingsIsClosed(t *testing.T) {
	t.Run("ATGDT-I01: An unreadable declaration reads as no role", func(t *testing.T) {})
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(settings.Path(dir)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settings.Path(dir), []byte("role: [product-owner\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := autonomyGuide(dir)
	if strings.Contains(g, "Your role (") || !strings.Contains(g, "Do not ask whoever is running you.") {
		t.Errorf("an unreadable declaration must fall on the closed side:\n%s", g)
	}
}

// `decidesProduct` erra para o lado FECHADO quando não consegue ler.
//
// Um erro de leitura que resultasse em "pode perguntar" seria a falha silenciosa mais cara
// deste mecanismo: o agente perguntaria, e ninguém saberia que a régua não foi aplicada.
func TestDecidesProduct_erraParaOLadoFechado(t *testing.T) {
	if common.DecidesProduct(t.TempDir()) {
		t.Error("sem arquivo, o agente não decide o produto")
	}
	if common.DecidesProduct("/caminho/que/nao/existe") {
		t.Error("com erro de leitura, o agente não decide o produto")
	}
}

// `/dev/null` NÃO é terminal, e o `ModeCharDevice` sozinho não sabe disso.
//
// Foi o defeito que impediu um agente de trabalhar: `/dev/null` também é char device, então
// `anchors next < /dev/null` — o caso normal de execução não-interativa — passava pela
// guarda, caía na pergunta do perfil e morria com `erro: ler a resposta: EOF`.
//
// Um comando que morre por EOF não diz o que fazer: parece defeito da ferramenta, e quem lê
// vai procurar o problema no ambiente.
func TestInteractiveTerminal_devNullNaoEhTerminal(t *testing.T) {
	nul, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer nul.Close()

	original := os.Stdin
	os.Stdin = nul
	defer func() { os.Stdin = original }()

	if common.InteractiveTerminal() {
		t.Error("`/dev/null` foi tratado como terminal — o agente cairia na pergunta e " +
			"morreria por EOF")
	}
}

// E um arquivo comum também não: `anchors next < resposta.txt` não é alguém do outro lado.
func TestInteractiveTerminal_arquivoNaoEhTerminal(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	original := os.Stdin
	os.Stdin = f
	defer func() { os.Stdin = original }()

	if common.InteractiveTerminal() {
		t.Error("um arquivo comum foi tratado como terminal")
	}
}

// A NEW DEV asked their agent to contribute. The agent read CONTRIBUTING, built the
// CORRECT onboarding plan — install, declare the role, `doctor --fix`, `guide work` — and
// stopped at step 4 asking for authorization: "step 4 changes the remote repository
// (branch protection, labels) (...) I want your explicit OK".
//
// The caution was right. What was missing was knowing that `doctor --fix` is IDEMPOTENT:
// on a repository already set up it changes nothing. Without that, the agent picks between
// two errors — ask authorization for everything (and never start) or for nothing.
func TestAutonomyGuidePreparationAsksNoAuthorization(t *testing.T) {
	t.Run("ATGDT-B04: Every profile reads that preparing the environment asks no authorization", func(t *testing.T) {})
	g := autonomyGuide(t.TempDir())

	if !strings.Contains(g, "does not ask for authorization") {
		t.Error("the guide should say that preparing the environment asks no authorization")
	}
	if !strings.Contains(g, "idempotent") && !strings.Contains(g, "IDEMPOTENT") {
		t.Error("the reason has to be said: the commands check before acting")
	}
	if !strings.Contains(g, "doctor --fix") || !strings.Contains(g, "anchors map build") {
		t.Error("the preparation commands have to be named")
	}
	// the ruler has to be NAMED, or the agent generalizes wrongly — "touching the remote"
	// would include `git push`, and it would ask OK again to deliver work
	if !strings.Contains(g, "REVERSIBILITY") {
		t.Error("the guide should say WHICH is the ruler, not only list exceptions")
	}
	// and the counter-example: without it, "preparation asks no authorization" reads as
	// "nothing asks authorization"
	if !strings.Contains(g, "DOES ask for authorization") {
		t.Error("the guide should say what STILL asks for authorization")
	}
}

func TestAutonomyGuideSaysSettingsRoleTakesArguments(t *testing.T) {
	t.Run("ATGDT-B04: Every profile reads that preparing the environment asks no authorization", func(t *testing.T) {})
	// The agent also stopped for this: "step 3 may be interactive (...) I stop and hand the
	// command back to you". The command takes the role and the date as arguments, and no
	// text said so.
	g := autonomyGuide(t.TempDir())
	if !strings.Contains(g, "anchors settings role <role> --date") {
		t.Error("the guide should show the NON-interactive form of `settings role`")
	}
	if !strings.Contains(g, "waits") && !strings.Contains(g, "waiting") {
		t.Error("the consequence has to be said: an agent with no terminal waits")
	}
}

// The section holds for EVERY profile, and that is not a detail: the first version put it
// after the fork that separates who decides the product from who does not, and the
// `architect` — which returns early — never saw it.
func TestAutonomyGuidePreparationHoldsForEveryProfile(t *testing.T) {
	t.Run("ATGDT-B04: Every profile reads that preparing the environment asks no authorization", func(t *testing.T) {})
	for _, role := range []string{"dev", "architect", "product-owner", "qa", "reviewer"} {
		dir := t.TempDir()
		s := settings.Settings{Role: settings.Role(role), DecidedAt: "2026-09-10"}
		if err := settings.Save(dir, s); err != nil {
			t.Fatalf("save the settings of %s: %v", role, err)
		}
		if g := autonomyGuide(dir); !strings.Contains(g, "does not ask for authorization") {
			t.Errorf("role %s does not see the preparation section", role)
		}
	}
}

// A settings file that cannot be read keeps the closed guide and says why, instead of
// claiming no role was declared.
func TestAutonomyGuideNamesAnUnreadableSettingsFile(t *testing.T) {
	t.Run("ATGDT-B08: An unreadable settings file is named, and the guide stays closed", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(settings.Path(root)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settings.Path(root), []byte("role: [not: a role\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	text := autonomyGuide(root)
	if !strings.Contains(text, "could not be read") || strings.Contains(text, "You did not declare a role") {
		t.Errorf("an unreadable settings file must be named, not reported as no declaration:\n%s", text)
	}
	if !strings.Contains(text, "Do not ask whoever is running you.") {
		t.Errorf("the guide must stay on the closed side:\n%s", text)
	}
}
