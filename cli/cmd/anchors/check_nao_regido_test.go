package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func cfgRegencia() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{
		"hook":    {Pattern: "src/hooks/**/*.ts", Kind: "code"},
		"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
		"feature": {Pattern: "**/*.feature", Kind: "feature"},
		"test":    {Pattern: "**/*.test.ts", Kind: "test"},
		"doc":     {Pattern: "**/*.md", Kind: "doc"},
	}}
}

// O FURO que este teste tranca: `selectNodes` respondia a MESMA coisa para duas
// situações opostas — "arquivo regido novo, ainda fora do mapa" e "arquivo que o
// Anchors nem rege" — e o pre-commit, sem ter como distingui-las, tratava as duas
// como benignas. O resultado era o pior caso possível: um hook/tela/service NOVO
// commitava limpo sem spec, sem feature e sem teste, porque fora do mapa nenhum
// gate o confronta. Exatamente o trabalho que o framework existe para barrar.
func TestSelectNodesDistingueRegidoDeNaoRegido(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src", "hooks"), 0o755)

	regido := "src/hooks/useNovo.ts"
	naoRegido := "package.json"
	os.WriteFile(filepath.Join(dir, regido), []byte("export const x = 1\n"), 0o644)
	os.WriteFile(filepath.Join(dir, naoRegido), []byte("{}\n"), 0o644)

	// Mapa VAZIO: é a condição real logo após criar um arquivo, antes do `map build`.
	g := &mapx.Graph{}
	cfg := cfgRegencia()

	t.Run("regido fora do mapa é OFENSA (barra)", func(t *testing.T) {
		_, _, err := selectNodes(g, cfg, false, []string{regido}, dir)
		if err == nil {
			t.Fatal("passou sem erro — arquivo regido fora do mapa tem de barrar o commit")
		}
		var nr errNotGoverned
		if errors.As(err, &nr) {
			t.Fatalf("classificado como NÃO-REGIDO (sairia %d, o hook faria continue): %v", ExitNotGoverned, err)
		}
		if !strings.Contains(err.Error(), "REGIDO") {
			t.Fatalf("mensagem não diz que o arquivo é regido: %v", err)
		}
	})

	t.Run("não-regido sai com código próprio (benigno)", func(t *testing.T) {
		_, _, err := selectNodes(g, cfg, false, []string{naoRegido}, dir)
		if err == nil {
			t.Fatal("esperava o sinal de não-regido, veio nil")
		}
		var nr errNotGoverned
		if !errors.As(err, &nr) {
			t.Fatalf("não sinalizou não-regido — o hook barraria package.json: %v", err)
		}
	})

	t.Run("inexistente não vira não-regido", func(t *testing.T) {
		// Erro de digitação no caminho não pode virar "benigno" e sumir: sairia 3 e o
		// pre-commit daria continue, silenciando o engano.
		_, _, err := selectNodes(g, cfg, false, []string{"src/hooks/naoExiste.ts"}, dir)
		if err == nil {
			t.Fatal("esperava erro para caminho inexistente")
		}
		var nr errNotGoverned
		if errors.As(err, &nr) {
			t.Fatalf("caminho inexistente classificado como não-regido: %v", err)
		}
	})

	t.Run("registro do Anchors (issues/) não é regido", func(t *testing.T) {
		// `issues/` casa a camada `doc` (`**/*.md`) mas o scanner NUNCA o indexa: é a
		// SAÍDA do próprio Anchors. Sem consultar o ignore, o caminho virava "regido
		// fora do mapa" — e `map build` não o acrescentava nunca, então o commit ficava
		// barrado para sempre. Medido ao commitar as issues que o `check` resolveu.
		os.MkdirAll(filepath.Join(dir, "issues", "done"), 0o755)
		iss := "issues/done/2026-08-15--violation--x.md"
		os.WriteFile(filepath.Join(dir, iss), []byte("# issue\n"), 0o644)
		_, _, err := selectNodes(g, cfg, false, []string{iss}, dir)
		var nr errNotGoverned
		if !errors.As(err, &nr) {
			t.Fatalf("issues/ deveria ser não-regido (exit %d), veio: %v", ExitNotGoverned, err)
		}
	})
}

// O `-progress.md` é o mesmo impasse que o `issues/`/`changes/` já resolvia, por outra
// porta: ele casa a camada `plan` (`plans/*.md`) e o scanner NUNCA o indexa — de
// propósito, porque um arquivo que existe para MUDAR não pode ser confrontado por gates
// que cobram justificativa de mudança.
//
// Sem a exclusão, o resultado era o pior caso: o arquivo dito "regido", ausente do mapa,
// e `map build` não o acrescentando nunca — commit barrado para sempre pelo próprio
// mecanismo que separou decisão de estado.
//
// Medido no blue-eyes ao commitar os 17 progressos que o `anchors new progress` acabara
// de criar.
func TestSelectNodes_progressoNaoEhRegido(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "plans"), 0o755)

	plano := "plans/0002-plataforma.md"
	progresso := "plans/0002-plataforma-progress.md"
	os.WriteFile(filepath.Join(dir, plano), []byte("# Plano\n"), 0o644)
	os.WriteFile(filepath.Join(dir, progresso), []byte("# Progresso\n\n- [x] feito\n"), 0o644)

	cfg := cfgRegencia()
	cfg.Layers["plan"] = config.Layer{Pattern: "plans/*.md", Kind: "plan"}
	g := &mapx.Graph{} // mapa vazio: a condição real logo depois de criar o arquivo

	_, _, err := selectNodes(g, cfg, false, []string{progresso}, dir)
	if err == nil {
		t.Fatal("passou sem erro — o esperado é errNotGoverned, não silêncio")
	}
	var nr errNotGoverned
	if !errors.As(err, &nr) {
		t.Fatalf("o progresso foi tratado como REGIDO e o commit ficaria barrado "+
			"para sempre (o `map build` nunca o acrescenta): %v", err)
	}

	// e o PLANO continua regido: a exclusão é do companheiro, não da camada.
	_, _, err = selectNodes(g, cfg, false, []string{plano}, dir)
	if err == nil {
		t.Fatal("o plano fora do mapa devia barrar")
	}
	if errors.As(err, &nr) {
		t.Fatal("o plano foi classificado como não-regido — a exclusão vazou para a camada")
	}
}
