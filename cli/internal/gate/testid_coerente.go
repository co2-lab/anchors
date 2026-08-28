package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkTestIDCoerente — o testID é UM contrato com QUATRO pontas, e este gate as
// confronta de uma vez.
//
// Os gates que este substitui (`testid-declared`, `testid-honored`) cobriam três
// arestas e deixavam a quarta descoberta:
//
//	código → spec        exposto sem declarar      (testid-declared)
//	spec   → código      declarado sem expor       (testid-declared)
//	spec   → consumidor  declarado sem consumidor  (testid-honored)
//	consumidor → código  CONSULTADO SEM EXISTIR    ← ninguém
//
// A quarta faltava por uma decisão explícita, escrita no gate antigo: "um teste que
// consulta um id INEXISTENTE já falha ao rodar, ruidosamente — cobrar isso aqui seria
// duplicar um sinal que o runner dá melhor".
//
// A premissa foi REFUTADA por medição (app de referência, 2026-08-24): 11 testIDs inexistentes em
// quatro commits numa tarde. "Falha ao rodar" pressupõe que o flow CHEGUE à linha do
// id. Numa suíte cujo primeiro passo já falhava, os sete ids fantasma seguintes eram
// invisíveis; e quando o primeiro foi corrigido, o erro apontou a TELA DE DESTINO
// ("`:ivle-button-category` is visible"), não o id errado — mandando procurar defeito
// onde não havia. O pior caso não dá sinal nenhum: um id de confirmação trocado pelo
// botão vizinho faz o teste CANCELAR a ação que deveria confirmar, e passar verde.
//
// Por que UM gate e não mais um parcial: um testID renomeado no código produzia dois
// achados, em gates diferentes ("exposto sem declarar" para o nome novo, "declarado sem
// expor" para o velho), quando o fato é um só — alguém renomeou e não propagou. O laudo
// aqui é por testID, com as quatro pontas na mesma linha, porque é assim que se lê o
// que aconteceu em vez de reconstruí-lo.
//
// Parte da spec (`on: [spec]`) porque é ela que declara o inventário — o mesmo ponto de
// partida dos gates que substitui.
func checkTestIDCoerente(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "o inventário de testID é declarado na spec — é dela que o confronto parte"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}
	attr := handleDeTeste(cfg)
	if attr == "" {
		// Inferir `testID` por default faria o gate reportar VERDE sobre o que não
		// conferiu — a pior falha possível num medidor.
		return Skip, "o projeto não declara `derived.test_handle` — não há atributo de ancoragem a confrontar"
	}

	// PONTA 1 — o código: o que a unidade realmente expõe.
	var expostos []string
	var arquivo string
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		if b, err := os.ReadFile(filepath.Join(root, e.To)); err == nil {
			expostos = append(expostos, testIDsExpostos(string(b), attr)...)
			arquivo = filepath.Base(e.To)
		}
	}
	if arquivo == "" {
		return Skip, "spec sem código ligado — nada a confrontar"
	}

	// PONTA 2 — a spec: o inventário declarado.
	declarados := testIDsDeclarados(content, attr)

	if len(expostos) == 0 && len(declarados) == 0 {
		// Nem toda unidade tem superfície de teste. Cobrar inventário de quem não marca
		// elemento algum transformaria o gate em ruído sobre toda unidade de
		// apresentação pura.
		return Skip, "a unidade não expõe `" + attr + "` — não há contrato a declarar"
	}

	// PONTAS 3 e 4 — os consumidores, separados por natureza porque respondem a
	// perguntas diferentes: a feature DESCREVE o handle (documentação executável), o
	// teste/flow o CONSULTA (verificação). Um id descrito e nunca consultado é cenário
	// que ninguém automatizou; consultado e nunca descrito é verificação sem contrato.
	blobConsulta, blobFeature, temConsumidor := superficiesConsumidoras(root, n.ID, g, cfg)

	// O confronto é por ID — a união de tudo que qualquer ponta menciona.
	// A chave é NORMALIZADA (sem a marca): o código escreve `:idep-x` e a spec grava
	// `idep-x` — mesmo handle, duas grafias. Indexar pela string crua fazia o mesmo id
	// aparecer DUAS VEZES no laudo, uma por grafia, dando a impressão de dois defeitos
	// onde há um. O valor guarda a forma como o CÓDIGO a escreve, que é a que o leitor
	// vai procurar no arquivo.
	universo := map[string]string{}
	for _, id := range declarados {
		universo[strings.TrimPrefix(id, ":")] = id
	}
	// O exposto entra por último para a grafia do código prevalecer na exibição.
	for _, id := range expostos {
		universo[strings.TrimPrefix(id, ":")] = id
	}
	// A quarta aresta (consumidor → código) NÃO entra aqui, e a razão é de escopo, não
	// de importância: a superfície e2e é a árvore INTEIRA de flows do projeto — um id de
	// outra tela aparece nela do mesmo jeito. Cobrar da spec do InventoryMove um
	// `:acnw-account-new-*` seria acusá-la por um handle de contas bancárias. Medido ao
	// tentar: 829 achados numa spec só, todos de outras telas.
	//
	// Não há como saber, a partir de UMA spec, qual dela deveria responder por um id
	// que só o flow menciona. A pergunta é do PROJETO ("existe flow procurando handle
	// que ninguém expõe?"), e é respondida por um gate de `scope: project` que varre os
	// dois lados de uma vez — ver `testid-consultado-existe`.

	// A comparação passa por `cobre` (do gate irmão), NUNCA por igualdade de string: a
	// spec grava `ivmv-screen` e o código expõe `:ivmv-screen` — a marca é convenção de
	// escrita, não parte da identidade. E o curinga casa dos dois lados (`item-*` cobre
	// `item-3`). Medido ao usar igualdade crua: 14 falsos positivos numa spec que
	// declara tudo corretamente, com o gate irmão passando limpo na mesma spec.
	var linhas []string
	for _, id := range universo {
		noCodigo, naSpec := cobertoPor(expostos, id), cobertoPor(declarados, id)
		consultado := temConsumidor && idConsultado(blobConsulta, id)
		// A feature é OPCIONAL: nem todo handle precisa aparecer num cenário escrito
		// (um contêiner de tela raramente aparece). Ausência aqui não é ofensa — por
		// isso entra no laudo como informação, nunca como o motivo da reprovação.
		naFeature := blobFeature != "" && idConsultado(blobFeature, id)

		if noCodigo && naSpec && (consultado || !temConsumidor) {
			continue // coerente nas pontas que dá para conferir
		}
		linhas = append(linhas, fmt.Sprintf("  %s\n      código %s   spec %s   feature %s   consultado %s",
			id, marca(noCodigo), marca(naSpec), marca(naFeature), marcaConsulta(consultado, temConsumidor)))
	}

	if len(linhas) == 0 {
		return Pass, ""
	}
	sort.Strings(linhas)
	return Fail, fmt.Sprintf(
		"testID incoerente na trinca (%d) — `%s`:\n\n%s\n\n"+
			"O testID é UM contrato com quatro pontas: o código o expõe, a spec o declara, a "+
			"feature o descreve e o teste/flow o consulta. Cada ✗ é uma ponta que não fecha:\n\n"+
			"  código ✗  o consumidor procura um handle que a tela não emite — e isso NÃO "+
			"aparece sozinho: se um passo anterior falha, o erro nomeia a tela de destino, "+
			"não o id errado.\n"+
			"  spec ✗    superfície não-contratada — por onde a divergência de identidade "+
			"entra sem ninguém ver.\n"+
			"  consultado ✗  contrato sem apoio: parece cobertura, porque a superfície está lá "+
			"e o inventário está completo.\n\n"+
			"Corrija o id na ponta divergente, ou remova-o das outras",
		len(linhas), arquivo, strings.Join(linhas, "\n"))
}

func marca(ok bool) string {
	if ok {
		return "✓"
	}
	return "✗"
}

// marcaConsulta distingue "ninguém consulta" de "não há onde procurar". Um traço onde
// o projeto não declarou superfície consumidora é honesto; um ✗ ali acusaria o autor
// por uma configuração ausente.
func marcaConsulta(ok, temConsumidor bool) string {
	if !temConsumidor {
		return "—"
	}
	return marca(ok)
}

// cobertoPor: algum item da lista cobre este id? Delega a `cobre`, que trata a marca e
// o curinga — as duas formas em que o MESMO handle se escreve diferente entre pontas.
func cobertoPor(xs []string, id string) bool {
	for _, x := range xs {
		if cobre(x, id) || cobre(id, x) {
			return true
		}
	}
	return false
}

// superficiesConsumidoras separa quem CONSULTA (teste ligado, teste vizinho, flows e2e)
// de quem DESCREVE (a feature). Devolve também se havia onde procurar — sem isso o gate
// não distingue "ninguém consulta" de "o projeto não declarou superfície", e acusaria o
// autor por uma configuração ausente.
func superficiesConsumidoras(root, specID string, g *mapx.Graph, cfg *config.Config) (consulta, feature string, temConsumidor bool) {
	var cs, fs []string
	for _, e := range g.Neighbors(specID).Out {
		if e.Type != mapx.EdgeCoveredBy {
			continue
		}
		// A feature em si: DESCREVE o handle nos cenários.
		if b, err := os.ReadFile(filepath.Join(root, e.To)); err == nil {
			fs = append(fs, string(b))
		}
		// O teste ligado nasce na FEATURE, não na spec — dois saltos.
		for _, fe := range g.Neighbors(e.To).Out {
			if fe.Type == mapx.EdgeTestedBy {
				if b, err := os.ReadFile(filepath.Join(root, fe.To)); err == nil {
					cs = append(cs, string(b))
				}
			}
		}
	}
	cs = append(cs, lerSuperficieE2E(root, cfg)...)
	// O teste COMPARTILHADO ou do PAI: um arquivo prova várias unidades e a aresta
	// `tested-by` não alcança todas. Sem isto, handle consultado aparece como órfão.
	cs = append(cs, lerTestesVizinhos(root, specID)...)
	return strings.Join(cs, "\n"), strings.Join(fs, "\n"), len(cs) > 0
}
