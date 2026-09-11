package migra

// FORMATO 2 — as chaves e os nomes de gate passam a inglês.
//
// O `lang` traduz o que se LÊ, nunca o que se ESCREVE na configuração: uma chave de YAML e
// um nome de gate são IDENTIFICADORES, e um arquivo com eles em português deixa de
// funcionar num time de outro idioma pelo motivo errado.
//
// O QUE ESTE PASSO CONVERTE, e por que cada coisa importa:
//
//   · as quatro CHAVES que nasceram em português. A pior de perder é `julgamentos`: sem a
//     conversão, os carimbos de julgamento de IA evaporam e o `check` refaz todos,
//     cobrando de novo o que alguém já respondeu. Medido no blue-eyes: 120 arestas
//     carimbadas, 135 vereditos.
//
//   · os NOMES DE GATE gravados dentro dos carimbos. Estes vinham de uma tabela de alias
//     que o `config` consultava para sempre — e ela tinha um defeito de assimetria: a
//     LEITURA normalizava (`mesmoGate`), a ESCRITA não. Um projeto que renomeasse o gate
//     ganhava um SEGUNDO carimbo em vez de atualizar o primeiro. Medido: 40 julgamentos
//     como `regra-cumprida` convivendo com 2 como `rule-fulfilled`, o mesmo gate contado
//     duas vezes.
//
// Com o contrato de formato, o alias deixa de ser necessário: converte-se uma vez, e o
// formato 2 passa a ter só o nome canônico.

func init() {
	Register(Step{
		To:  2,
		Why: "chaves e nomes de gate passam a inglês (o vocabulário do produto é fixo)",
		RenameKeys: map[string]map[string]string{
			"anchors.graph.yaml": {
				"gerado_por":     "generated_by",
				"code_declarado": "code_declared",
				"julgamentos":    "judgments",
			},
			"anchors.yaml": {
				"trinca_opcional": "triad_optional",
			},
		},
		RenameValues: map[string]map[string]map[string]string{
			// Os nomes de gate aparecem como VALOR de `gate:` (dentro dos julgamentos e
			// dos carimbos) e de `id:`/`check:` (na declaração dos gates do projeto).
			"anchors.graph.yaml": {
				"gate": {
					"codigo-catalogado":           "code-cataloged",
					"dependencia-vulneravel":      "dependency-vulnerable",
					"fase-existe":                 "phase-exists",
					"fase-ordenada":               "phase-ordered",
					"feature-nao-vazia":           "feature-not-empty",
					"mock-carimbado":              "mock-stamped",
					"mock-detect-cobre-o-dialeto": "mock-detect-covers-dialect",
					"mock-tipado":                 "mock-typed",
					"no-test-prova-real":          "no-test-proof-real",
					"parent-valido":               "parent-valid",
					"plano-alterado-justificado":  "plan-change-justified",
					"plano-revisado":              "plan-revised",
					"prova-cruza-fronteira":       "proof-crosses-boundary",
					"regra-cumprida":              "rule-fulfilled",
					"sbom-gerado":                 "sbom-generated",
					"secret-nao-vazado":           "no-secret-leaked",
					"sem-duplicacao":              "no-duplication",
					"spec-completa":               "spec-complete",
					"spec-tem-codigo":             "spec-has-code",
					"teste-rastreavel":            "test-traceable",
					"trinca-completa":             "triad-complete",
				},
			},
			"anchors.yaml": {
				"name": {
					"codigo-catalogado":           "code-cataloged",
					"dependencia-vulneravel":      "dependency-vulnerable",
					"fase-existe":                 "phase-exists",
					"fase-ordenada":               "phase-ordered",
					"feature-nao-vazia":           "feature-not-empty",
					"mock-carimbado":              "mock-stamped",
					"mock-detect-cobre-o-dialeto": "mock-detect-covers-dialect",
					"mock-tipado":                 "mock-typed",
					"no-test-prova-real":          "no-test-proof-real",
					"parent-valido":               "parent-valid",
					"plano-alterado-justificado":  "plan-change-justified",
					"plano-revisado":              "plan-revised",
					"prova-cruza-fronteira":       "proof-crosses-boundary",
					"regra-cumprida":              "rule-fulfilled",
					"sbom-gerado":                 "sbom-generated",
					"secret-nao-vazado":           "no-secret-leaked",
					"sem-duplicacao":              "no-duplication",
					"spec-completa":               "spec-complete",
					"spec-tem-codigo":             "spec-has-code",
					"teste-rastreavel":            "test-traceable",
					"trinca-completa":             "triad-complete",
				},
				"id": {
					"codigo-catalogado":           "code-cataloged",
					"dependencia-vulneravel":      "dependency-vulnerable",
					"fase-existe":                 "phase-exists",
					"fase-ordenada":               "phase-ordered",
					"feature-nao-vazia":           "feature-not-empty",
					"mock-carimbado":              "mock-stamped",
					"mock-detect-cobre-o-dialeto": "mock-detect-covers-dialect",
					"mock-tipado":                 "mock-typed",
					"no-test-prova-real":          "no-test-proof-real",
					"parent-valido":               "parent-valid",
					"plano-alterado-justificado":  "plan-change-justified",
					"plano-revisado":              "plan-revised",
					"prova-cruza-fronteira":       "proof-crosses-boundary",
					"regra-cumprida":              "rule-fulfilled",
					"sbom-gerado":                 "sbom-generated",
					"secret-nao-vazado":           "no-secret-leaked",
					"sem-duplicacao":              "no-duplication",
					"spec-completa":               "spec-complete",
					"spec-tem-codigo":             "spec-has-code",
					"teste-rastreavel":            "test-traceable",
					"trinca-completa":             "triad-complete",
				},
				"check": {
					"codigo-catalogado":           "code-cataloged",
					"dependencia-vulneravel":      "dependency-vulnerable",
					"fase-existe":                 "phase-exists",
					"fase-ordenada":               "phase-ordered",
					"feature-nao-vazia":           "feature-not-empty",
					"mock-carimbado":              "mock-stamped",
					"mock-detect-cobre-o-dialeto": "mock-detect-covers-dialect",
					"mock-tipado":                 "mock-typed",
					"no-test-prova-real":          "no-test-proof-real",
					"parent-valido":               "parent-valid",
					"plano-alterado-justificado":  "plan-change-justified",
					"plano-revisado":              "plan-revised",
					"prova-cruza-fronteira":       "proof-crosses-boundary",
					"regra-cumprida":              "rule-fulfilled",
					"sbom-gerado":                 "sbom-generated",
					"secret-nao-vazado":           "no-secret-leaked",
					"sem-duplicacao":              "no-duplication",
					"spec-completa":               "spec-complete",
					"spec-tem-codigo":             "spec-has-code",
					"teste-rastreavel":            "test-traceable",
					"trinca-completa":             "triad-complete",
				},
			},
		},
	})
}
