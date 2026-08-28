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
func DeduceModulePrefixes(modules []string) map[string]string {
	sorted := append([]string(nil), modules...)
	sort.Strings(sorted)

	out := map[string]string{}
	taken := map[string]bool{}
	for _, m := range sorted {
		name := filepath.Base(strings.TrimRight(m, "/"))
		pfx := code.ModulePrefix(name)
		// unicidade entre módulos: se colide, varia a 2ª letra
		if taken[pfx] {
			base := pfx[:1]
			for c := byte('A'); c <= 'Z'; c++ {
				cand := base + string(c)
				if !taken[cand] {
					pfx = cand
					break
				}
			}
		}
		taken[pfx] = true
		out[name] = pfx
	}
	return out
}
