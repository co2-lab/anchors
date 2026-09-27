package initx

import (
	"maps"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
)

// ApplyPreset grava no cfg as layers de CÓDIGO de um preset de stack. Para presets
// MODULARES, deduz o code_prefix de cada módulo real encontrado sob o ModuleGlob
// (usando os diretórios detectados), ligando o gerador de identidade (Camada 1) à
// Estrutura — em vez de uma tabela hardcoded. Puro: recebe os módulos detectados,
// não varre disco. Devolve os prefixos deduzidos (módulo → prefixo) para exibição.
func ApplyPreset(cfg *config.Config, p Preset, modules []string) map[string]string {
	if cfg.Layers == nil {
		cfg.Layers = map[string]config.Layer{}
	}
	prefixes := DeduceModulePrefixes(modules)

	// O prefixo de identidade é POR módulo, não por layer (o `anchors code` resolve o
	// prefixo pelo caminho → módulo → prefixo na hora de gerar). As layers do preset
	// entram como estão; a tag 'module' que já vem do preset marca as modulares.
	maps.Copy(cfg.Layers, p.ToLayers())
	return prefixes
}

// DeduceModulePrefixes deriva um prefixo de 2 chars para cada módulo (o basename do
// diretório do módulo), garantindo UNICIDADE entre eles — se dois módulos gerariam o
// mesmo prefixo (ex.: "auth" e "audit" → AU), o segundo é ajustado. Determinístico:
// processa em ordem alfabética.
//
// The key is the basename, unless two modules share it: then each of them is keyed by
// its path. Keyed by basename alone, `apps/auth` and `packages/auth` collapsed into one
// entry and one module vanished from the mapping.
func DeduceModulePrefixes(modules []string) map[string]string {
	sorted := make([]string, 0, len(modules))
	for _, m := range modules {
		sorted = append(sorted, strings.TrimRight(m, "/"))
	}
	sort.Strings(sorted)

	sameName := map[string]int{}
	for _, m := range sorted {
		sameName[filepath.Base(m)]++
	}

	out := map[string]string{}
	taken := map[string]bool{}
	for _, m := range sorted {
		name := filepath.Base(m)
		key := name
		if sameName[name] > 1 {
			key = m
		}
		pfx := code.ModulePrefix(name)
		// unicidade entre módulos: se colide, varia a 2ª letra; com as 26 da mesma
		// inicial tomadas, varia também a 1ª. Parar na 2ª letra devolvia o prefixo JÁ
		// tomado ao 27º módulo da mesma inicial — dois módulos com a mesma identidade.
		if taken[pfx] {
			pfx = firstFreePrefix(pfx[0], taken)
		}
		taken[pfx] = true
		out[key] = pfx
	}
	return out
}

// firstFreePrefix returns the first two-letter prefix not taken: the same first letter
// with the second from A to Z, then every other first letter in order from A. It
// returns "" only when all 676 are taken.
func firstFreePrefix(first byte, taken map[string]bool) string {
	firsts := []byte{first}
	for c := byte('A'); c <= 'Z'; c++ {
		if c != first {
			firsts = append(firsts, c)
		}
	}
	for _, f := range firsts {
		for c := byte('A'); c <= 'Z'; c++ {
			if cand := string([]byte{f, c}); !taken[cand] {
				return cand
			}
		}
	}
	return ""
}
