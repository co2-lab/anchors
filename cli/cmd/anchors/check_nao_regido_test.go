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
		var nr errNaoRegido
		if errors.As(err, &nr) {
			t.Fatalf("classificado como NÃO-REGIDO (sairia %d, o hook faria continue): %v", ExitNaoRegido, err)
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
		var nr errNaoRegido
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
		var nr errNaoRegido
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
		var nr errNaoRegido
		if !errors.As(err, &nr) {
			t.Fatalf("issues/ deveria ser não-regido (exit %d), veio: %v", ExitNaoRegido, err)
		}
	})
}
