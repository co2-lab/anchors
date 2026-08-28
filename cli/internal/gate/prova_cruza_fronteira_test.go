package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaProvaCruzada(t *testing.T, spec, codigo string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}, {ID: "u.ts", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "u.spec.md", To: "u.ts", Type: mapx.EdgeSpecifies}},
	}
	return checkProvaCruzaFronteira(spec, mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}, root, g, nil)
}

// O CASO QUE MOTIVOU O GATE, e ele está vivo no app de referência sem ter divergido ainda.
// `SEATX-B01` declara que o preço espelha o backend; `orgPlans.ts` cita
// `orgBilling.ts` DUAS VEZES em comentário e não importa uma linha; o teste faz
// `toBe(15)`. Prova que é 15 — não que é o mesmo que o backend cobra. Mudar o
// backend para 18 mantinha tudo verde.
func TestFronteira_regraQueEspelhaSemImportarReprova(t *testing.T) {
	spec := "| `SEATX-B01` | O preço vem do backend ***{fonte-unica}*** `orgBilling.ts` — divergir aqui mente sobre a cobrança. |\n"
	codigo := `// SEATX-B01: estes preços ESPELHAM o orgBilling.ts do backend.
export const SEAT_PRICE = { individual: 15, familia: 30 }`
	v, msg := rodaProvaCruzada(t, spec, codigo)
	if v != Fail {
		t.Fatalf("quer Fail, obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "SEATX-B01") || !strings.Contains(msg, "orgBilling") {
		t.Errorf("a mensagem devia nomear a regra e o alvo; msg = %q", msg)
	}
	if !strings.Contains(msg, "CÓPIA") {
		t.Errorf("a mensagem devia explicar que descreve uma cópia; msg = %q", msg)
	}
}

// O comentário é justamente o que NÃO conta. `orgPlans.ts` menciona o alvo duas
// vezes em prosa; se o gate lesse comentário, aprovaria o caso que o motivou.
func TestFronteira_citacaoSoEmComentarioNaoSatisfaz(t *testing.T) {
	spec := "| `XPTOX-B01` | ***{fonte-unica}*** `outra.ts` |\n"
	codigo := `/**
 * Ver outra.ts — a fonte real vive lá.
 */
// outra.ts define os mesmos valores
export const A = 1`
	if v, msg := rodaProvaCruzada(t, spec, codigo); v != Fail {
		t.Fatalf("menção só em comentário não satisfaz a relação; obteve %v (%s)", v, msg)
	}
}

// A contraparte, e é o `BLBLX-B06` do app de referência depois da correção: a regra afirma
// "fonte única com o backend" e o código IMPORTA a função de lá. Passa.
func TestFronteira_codigoQueImportaPassa(t *testing.T) {
	spec := "| `BLBLX-B06` | Σ lançamentos que contam no saldo (`contaNoSaldo`, ***{fonte-unica}*** `balanceReconciliation.ts`) |\n"
	codigo := `import { contaNoSaldo } from '@backend/business-logic/balanceReconciliation'
export const somar = (txs) => txs.filter(contaNoSaldo)`
	if v, msg := rodaProvaCruzada(t, spec, codigo); v != Pass {
		t.Fatalf("o código importa a unidade citada — devia passar; obteve %v (%s)", v, msg)
	}
}

// O alias do projeto não pode ser exigido igual ao caminho da spec: a spec cita
// `packages/backend/business-logic/orgBilling.ts` e o código importa
// `@backend/business-logic/orgBilling`. É a mesma unidade — a comparação é pelo
// nome do módulo, não pelo caminho literal.
func TestFronteira_aliasDiferenteDoCaminhoDaSpecAindaCasa(t *testing.T) {
	spec := "| `XPTOX-B01` | ***{fonte-unica}*** `packages/backend/business-logic/orgBilling.ts` |\n"
	codigo := `import { SEAT_PRICE } from '@backend/business-logic/orgBilling'
export const preco = SEAT_PRICE`
	if v, msg := rodaProvaCruzada(t, spec, codigo); v != Pass {
		t.Fatalf("alias e extensão variam por projeto; o módulo é o mesmo. obteve %v (%s)", v, msg)
	}
}

// A dispensa é a saída honesta: há relação, e ela NÃO é importável. Um contrato de
// rede, um arquivo gerado, um valor que vive num provedor externo. Sem esta porta o
// gate empurraria o autor a inventar um import falso para calá-lo.
func TestFronteira_dispensaDeclaradaPassa(t *testing.T) {
	spec := "| `XPTOX-B01` | ***{fonte-unica}*** `schema.json` do provedor @no-cross: contrato de rede — o arquivo não existe no repo, a conferência é no teste de integração |\n"
	codigo := "export const A = 1"
	// Skip, não Pass: sem exigência a confrontar o gate não afirma conformidade —
	// ele se cala. E a mensagem diz QUAL silêncio é: houve regra, foi dispensada.
	v, msg := rodaProvaCruzada(t, spec, codigo)
	if v != Skip {
		t.Fatalf("dispensa declarada devia ser Skip; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "no-cross") {
		t.Errorf("o Skip devia dizer que a relação foi dispensada, não que não existia; msg = %q", msg)
	}
}

// "É fonte única" em PROSA, sem arquivo citado: não há alvo a confrontar e não há
// marcação. O gate se cala — cobrar seria falso-positivo (o app de referência tem
// `ALLOWED_IMPORT_EXTENSIONS é fonte única` exatamente assim, e é legítimo: a
// unidade se declara dona do conceito, só não usou o carimbo ainda).
func TestFronteira_fonteUnicaEmProsaSemArquivoEhSkip(t *testing.T) {
	spec := "| `OCRBX-X02` | `ALLOWED_IMPORT_EXTENSIONS` é fonte única (picker/câmera/share validam por ela) |\n"
	if v, msg := rodaProvaCruzada(t, spec, "export const A = 1"); v != Skip {
		t.Fatalf("sem alvo e sem marcação — devia ser Skip, obteve %v (%s)", v, msg)
	}
}

// O CARIMBO DE DONO: a unidade que DETÉM o conceito não espelha ninguém, então não
// há import a exigir. É a assimetria que o usuário desenhou — `{fonte-unica}` diz
// "isto é um conceito compartilhado"; `(@fonte-unica)` diz "e a fonte sou eu".
func TestFronteira_carimboDeDonoNaoExigeImport(t *testing.T) {
	spec := "| `RDSRX-B04` | `contaNoSaldo`: só item de fatura fica fora ***{fonte-unica}*** (@fonte-unica) |\n"
	v, msg := rodaProvaCruzada(t, spec, "export function contaNoSaldo(tx) { return tx.sourceType !== 'invoice' }")
	if v != Pass {
		t.Fatalf("o dono não espelha ninguém — devia passar; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "É a fonte") {
		t.Errorf("a mensagem devia dizer que esta unidade é a fonte; msg = %q", msg)
	}
}

// A prosa não é mais cobrança, é AVISO — o que faz a convenção ser adotada em vez
// de esquecida. Quem escreveu "espelha o X.ts" provavelmente não conhecia o
// marcador; o aviso ensina sem barrar (nasce informativo).
func TestFronteira_prosaSemMarcacaoAvisa(t *testing.T) {
	spec := "| `OGIFX-R02` | Toda a copy vem de `presentation/copy.ts`, os mesmos valores do backend |\n"
	v, msg := rodaProvaCruzada(t, spec, "export const A = 1")
	if v != Fail {
		t.Fatalf("prosa que fala de relação devia avisar; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "SEM a marcação") {
		t.Errorf("a mensagem devia pedir a marcação; msg = %q", msg)
	}
	if !strings.Contains(msg, "OGIFX-R02") {
		t.Errorf("devia nomear a regra; msg = %q", msg)
	}
}

// Marcou mas não disse de qual arquivo, e não carimbou como dona: a marcação fica
// sem sentido — ou falta o alvo, ou falta o carimbo.
func TestFronteira_marcacaoSemAlvoNemCarimboAvisa(t *testing.T) {
	spec := "| `XPTOX-B01` | o valor é ***{fonte-unica}*** do domínio |\n"
	v, msg := rodaProvaCruzada(t, spec, "export const A = 1")
	if v != Fail {
		t.Fatalf("marcação sem alvo nem carimbo devia avisar; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "sem dizer de QUAL arquivo") {
		t.Errorf("a mensagem devia explicar o que falta; msg = %q", msg)
	}
}

// Regra que não afirma relação nenhuma não é assunto deste gate.
func TestFronteira_regraSemAfirmacaoDeRelacaoEhSkip(t *testing.T) {
	spec := "| `XPTOX-B01` | soma os lançamentos do mês e devolve o total |\n"
	if v, _ := rodaProvaCruzada(t, spec, "export const A = 1"); v != Skip {
		t.Fatalf("regra sem afirmação de relação devia ser Skip, obteve %v", v)
	}
}

// A afirmação e o arquivo têm de estar na MESMA linha de regra. Um parágrafo de
// prosa dizendo "fonte única" longe da tabela não vira exigência — senão o gate
// cobraria a introdução da spec, onde "fonte única" aparece como intenção geral.
func TestFronteira_afirmacaoForaDaLinhaDeRegraNaoConta(t *testing.T) {
	spec := `## Visão Geral
Esta unidade é fonte única e espelha o ` + "`outra.ts`" + ` do backend.

## Rules
| Código | Efeito |
| --- | --- |
| ` + "`XPTOX-B01`" + ` | soma os valores |
`
	if v, _ := rodaProvaCruzada(t, spec, "export const A = 1"); v != Skip {
		t.Fatalf("afirmação em prosa fora da linha de regra não é exigência, obteve %v", v)
	}
}

func TestImportaUnidade(t *testing.T) {
	casos := []struct {
		nome  string
		code  string
		alvo  string
		match bool
	}{
		{"import nomeado", "import { x } from '@backend/orgBilling'\n", "packages/backend/orgBilling.ts", true},
		{"require", "const b = require('./orgBilling')\n", "orgBilling.ts", true},
		{"só menção em string", "const msg = 'ver orgBilling para detalhes'\n", "orgBilling.ts", false},
		{"nome ausente", "import { y } from './outro'\n", "orgBilling.ts", false},
		{"import multilinha", "import {\n  contaNoSaldo,\n} from '@backend/balanceReconciliation'\n", "balanceReconciliation.ts", true},
	}
	for _, c := range casos {
		if got := importaUnidade(c.code, c.alvo); got != c.match {
			t.Errorf("%s: importaUnidade = %v, quer %v", c.nome, got, c.match)
		}
	}
}

// rodaComMapa monta um grafo com DUAS unidades, para o alvo por código de regra
// ter o que resolver.
func rodaComMapa(t *testing.T, spec, codigo string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "dono.ts"), []byte("export const X = 1"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "u.spec.md", Kind: mapx.KindSpec, Code: "SEATX"},
			{ID: "u.ts", Kind: mapx.KindCode, Code: "SEATX"},
			{ID: "dono.spec.md", Kind: mapx.KindSpec, Code: "RDSRX"},
			{ID: "dono.ts", Kind: mapx.KindCode, Code: "RDSRX"},
		},
		Edges: []mapx.Edge{{From: "u.spec.md", To: "u.ts", Type: mapx.EdgeSpecifies}},
	}
	return checkProvaCruzaFronteira(spec, mapx.Node{ID: "u.spec.md", Kind: mapx.KindSpec}, root, g, nil)
}

// A FORMA PREFERIDA: o alvo é o CÓDIGO DA REGRA, não o caminho do arquivo. O código
// é identidade estável — já usado em spec, feature, teste e comentário — e sobrevive
// à refatoração que move o arquivo. O repo tem a prova de que o caminho não
// sobrevive: o comentário de `orgPlans.ts` citava `functions/_shared/orgBilling.ts`,
// caminho que não existia mais.
func TestFronteira_alvoPorCodigoDeRegraResolvePeloMapa(t *testing.T) {
	spec := "| `SEATX-B01` | o preço é ***{fonte-unica}*** de `RDSRX-B04` |\n"
	// O código NÃO importa `dono.ts` — reprova.
	v, msg := rodaComMapa(t, spec, "export const SEAT_PRICE = { individual: 15 }")
	if v != Fail {
		t.Fatalf("alvo por código de regra devia ser cobrado; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "dono.ts") {
		t.Errorf("a mensagem devia nomear o arquivo RESOLVIDO pelo mapa; msg = %q", msg)
	}

	// Agora importando: passa.
	v2, msg2 := rodaComMapa(t, spec, "import { X } from './dono'\nexport const SEAT_PRICE = X")
	if v2 != Pass {
		t.Fatalf("com o import do alvo resolvido devia passar; obteve %v (%s)", v2, msg2)
	}
}

// Código de regra que não existe no mapa avisa em vez de passar em silêncio — um
// alvo que não resolve é uma relação que ninguém confronta.
func TestFronteira_codigoDeRegraInexistenteAvisa(t *testing.T) {
	spec := "| `SEATX-B01` | ***{fonte-unica}*** de `ZZZZX-B99` |\n"
	v, msg := rodaComMapa(t, spec, "export const A = 1")
	if v != Fail {
		t.Fatalf("alvo que não resolve devia avisar; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "não resolve") {
		t.Errorf("a mensagem devia dizer que o alvo não resolve; msg = %q", msg)
	}
}

// Auto-referência não é relação: uma regra que cita OUTRA regra da MESMA unidade
// (`SEATX-B02` dentro da spec `SEATX`) não aponta para o outro lado de nada.
func TestFronteira_autoReferenciaNaoEhAlvo(t *testing.T) {
	spec := "| `SEATX-B01` | o preço é ***{fonte-unica}*** e `SEATX-B02` o consome |\n"
	v, msg := rodaComMapa(t, spec, "export const A = 1")
	// Sem alvo externo e sem carimbo de dono → a marcação fica sem sentido.
	if v != Fail || !strings.Contains(msg, "sem dizer de QUAL arquivo") {
		t.Fatalf("auto-referência não é alvo; esperava aviso de marcação sem alvo, obteve %v (%s)", v, msg)
	}
}
