package gate

import "testing"

// DOCUMENTAR é abrir uma seção para a unidade. CITAR é escrever o nome no corpo de outra.
//
// Duas armadilhas reais, no mesmo arquivo do projeto de referência:
//
//   · uma NOTA listava os componentes "que o gate `doc-required` deve sinalizar", e
//     `MetricCard` estava nela. O gate lia, achava o nome, e se dava por satisfeito —
//     silenciado pela própria nota que dizia que faltava documentá-lo.
//
//   · a seção de `StatusBadge` citava `MetricCard` ao explicar quando NÃO se usa um em vez
//     do outro. Prosa legítima — e não é documentação da unidade citada.
//
// Nos dois casos a unidade aparecia no arquivo e não tinha entrada. Ficou sem documentação
// por semanas, com o gate verde.
func TestMencaoNoCorpoNaoContaComoDocumentacao(t *testing.T) {
	doc := "# Catálogo\n\n" +
		"## `StatusBadge.spec.md` (`STBDS`)\n\n" +
		"Use para o estado agregado; nunca para o ponto de status — isso é `MTCRM`.\n\n" +
		"## Nota — este catálogo cresce por card\n\n" +
		"Faltam entrada para `MTCRM` e `CHRTS`; cada um recebe a sua quando o gate sinalizar.\n"

	if !mentionsUnit(doc, "STBDS", "StatusBadge.spec.md") {
		t.Error("a unidade COM seção própria deixou de contar — falso negativo")
	}
	if mentionsUnit(doc, "MTCRM", "MetricCard.spec.md") {
		t.Error("citação no corpo e menção em nota contaram como documentação — " +
			"é o defeito que deixou o MetricCard sem entrada por semanas")
	}
	if mentionsUnit(doc, "CHRTS", "Charts.spec.md") {
		t.Error("menção só na nota contou como documentação")
	}
}

// O TÍTULO pode nomear a unidade pelo ARQUIVO, não só pelo código — é como o catálogo do
// projeto de referência escreve, e era o que a rede secundária já cobria no corpo.
func TestTituloPeloNomeDoArquivoConta(t *testing.T) {
	doc := "## `apps/mobile/src/components/MetricCard.spec.md`\n\nO cartão de número.\n"
	if !mentionsUnit(doc, "MTCRM", "apps/mobile/src/components/MetricCard.spec.md") {
		t.Error("título que nomeia o ARQUIVO não contou — o catálogo real escreve assim")
	}
}

// DOCUMENTO SEM SEÇÕES continua valendo pela menção.
//
// Um `openapi.yaml` não tem títulos Markdown, e exigir uma seção por unidade ali seria
// cobrar uma estrutura que o formato não tem. Foi por isso que a menção existe como rede
// secundária desde o início, e este caso precisa seguir funcionando como antes.
func TestDocumentoSemSecoesSegueValendoPelaMencao(t *testing.T) {
	yaml := "openapi: 3.1.0\npaths:\n  /services:\n    get:\n      operationId: ServiceList\n"
	if !mentionsUnit(yaml, "SRLSS", "ServiceList.spec.md") {
		t.Error("num documento sem seções a menção deixou de contar — o OpenAPI não tem " +
			"títulos Markdown, e cobrar seção ali é exigir estrutura que o formato não tem")
	}
}

// O TÍTULO DE OUTRA UNIDADE não conta para esta, mesmo contendo o nome dela.
func TestTituloDeOutraUnidadeNaoConta(t *testing.T) {
	doc := "## `MetricCardList.spec.md` (`MTCLS`)\n\nA lista de cartões.\n"
	if mentionsUnit(doc, "MTCRM", "MetricCard.spec.md") {
		t.Error("o título de `MetricCardList` contou como documentação de `MetricCard`")
	}
}
