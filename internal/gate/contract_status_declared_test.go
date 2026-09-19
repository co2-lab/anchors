package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaContrato(t *testing.T, spec, codigo string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "handler.ts"), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "handler.spec.md", Kind: mapx.KindSpec},
			{ID: "handler.ts", Kind: mapx.KindCode},
		},
		Edges: []mapx.Edge{{From: "handler.spec.md", To: "handler.ts", Type: mapx.EdgeSpecifies}},
	}
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts"}}
	return checkContractStatusDeclared(spec, mapx.Node{ID: "handler.spec.md", Kind: mapx.KindSpec}, root, g, cfg)
}

// O CORPUS deste gate são os 8 achados reais da auditoria de 2026-08 no app de referência. Cada
// teste abaixo reproduz um deles: se o gate existisse antes, teria pegado todos.

// accept-org-invite declarava 3 status e emitia 8. Os dois 403 e o 409 ausentes
// eram a defesa contra sequestro de convite — o cliente programado pela tabela não
// tratava a recusa.
func TestContrato_statusDeSegurancaOmitidoReprova(t *testing.T) {
	t.Run("CSDCN-B01: A status emitted and not declared is accused by number", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | assento aceito |
| 404 | seat não encontrado |
| 5xx | falha |
`
	codigo := `export const handler = async () => {
  if (!caller) return jsonResponse(401, { error: 'Unauthorized' })
  if (!caller.email) return jsonResponse(403, { error: 'sem e-mail verificado' })
  if (!seat) return jsonResponse(404, { error: 'não encontrado' })
  if (seat.taken) return jsonResponse(409, { error: 'já aceito por outra conta' })
  return jsonResponse(200, { success: true })
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Fail {
		t.Fatalf("quer Fail, obteve %v (%s)", v, msg)
	}
	for _, s := range []string{"401", "403", "409"} {
		if !strings.Contains(msg, s) {
			t.Errorf("a mensagem devia acusar o %s omitido; msg = %q", s, msg)
		}
	}
}

// reanalyse-metadata declarava 402 para cota excedida e NENHUM caminho o emitia (a
// cota responde 429). Um cliente que tratasse 402 como "precisa pagar" nunca
// dispararia esse ramo — código morto que ninguém percebe.
func TestContrato_statusFantasmaReprova(t *testing.T) {
	t.Run("CSDCN-B02: A status declared and emitted by no path is accused as a phantom", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | reanálise concluída |
| 402 | cota de IA excedida |
| 429 | cota de IA excedida |
`
	codigo := `export const handler = async () => {
  if (semCota) return { statusCode: 429, body: '{}' }
  return { statusCode: 200, body: '{}' }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Fail {
		t.Fatalf("quer Fail, obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "402") {
		t.Errorf("a mensagem devia acusar o 402 fantasma; msg = %q", msg)
	}
	if !strings.Contains(msg, "morto") && !strings.Contains(msg, "dead") {
		t.Errorf("a mensagem devia explicar que é código morto no cliente; msg = %q", msg)
	}
}

// import-transactions declarava "200 com transações" — uma API síncrona que não
// existe. O caminho normal responde 202 `processing`.
func TestContrato_statusRealDiferenteDoDeclaradoReprova(t *testing.T) {
	t.Run("CSDCN-X02: The gate does not judge when each status is right", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | importado (batch + transações + resumo) |
| 400 | payload inválido |
| 401 | não autenticado |
| 5xx | falha |
`
	codigo := `export const handler = async () => {
  if (!caller) return { statusCode: 401, body: '{}' }
  if (foraDoPrefixo) return { statusCode: 403, body: '{}' }
  return { statusCode: 202, body: JSON.stringify({ status: 'processing' }) }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Fail {
		t.Fatalf("quer Fail, obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "202") || !strings.Contains(msg, "403") {
		t.Errorf("devia acusar 202 e 403 emitidos e não declarados; msg = %q", msg)
	}
	if !strings.Contains(msg, "400") {
		t.Errorf("devia acusar o 400 declarado que o código não emite; msg = %q", msg)
	}
}

// O contrato correto passa. Sem isto o gate seria um gerador de ruído.
func TestContrato_tabelaFielPassa(t *testing.T) {
	t.Run("CSDCN-B03: A faithful table passes", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | recibo confirmado |
| 400 | payload inválido |
| 401 | não autenticado |
| 403 | batch inexistente ou de outro usuário (anti-probing) |
| 409 | batch não está em draft |
| 5xx | falha |
`
	codigo := `export const handler = async () => {
  if (!caller) return { statusCode: 401, body: '{}' }
  if (!body) return { statusCode: 400, body: '{}' }
  if (!batch || batch.userId !== userId) return { statusCode: 403, body: '{}' }
  if (batch.status !== 'draft') return { statusCode: 409, body: '{}' }
  return { statusCode: 200, body: '{}' }
}`
	if v, msg := rodaContrato(t, spec, codigo); v != Pass {
		t.Fatalf("contrato fiel devia passar; obteve %v (%s)", v, msg)
	}
}

// O 500 do try/catch do topo não é decisão do handler: quase todo um tem, e as
// specs declaram `5xx`. Cobrá-lo produziria falso-positivo em massa.
func TestContrato_quinhentosDoTryCatchNaoEhCobrado(t *testing.T) {
	t.Run("CSDCN-B04: The 500 of the top-level try/catch is not charged", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
| 5xx | falha |
`
	codigo := `export const handler = async () => {
  try {
    return { statusCode: 200, body: '{}' }
  } catch {
    return { statusCode: 500, body: '{}' }
  }
}`
	if v, msg := rodaContrato(t, spec, codigo); v != Pass {
		t.Fatalf("o 500 do try/catch não devia ser cobrado; obteve %v (%s)", v, msg)
	}
}

// A faixa `5xx` cobre 502/503, mas NÃO serve de guarda-chuva para os 4xx: um `4xx`
// genérico esconderia justamente as recusas de acesso que este gate persegue.
func TestContrato_faixaCincoXXCobreMasQuatroXXNao(t *testing.T) {
	t.Run("CSDCN-B05: The 5xx range covers, the 4xx range does not", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
| 4xx | erro do cliente |
| 5xx | falha |
`
	codigo := `export const handler = async () => {
  if (a) return { statusCode: 503, body: '{}' }
  if (b) return { statusCode: 403, body: '{}' }
  return { statusCode: 200, body: '{}' }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Fail {
		t.Fatalf("o 403 sob um `4xx` genérico devia reprovar; obteve %v (%s)", v, msg)
	}
	if strings.Contains(msg, "503") {
		t.Errorf("o 503 é coberto por `5xx` e não devia ser acusado; msg = %q", msg)
	}
	if !strings.Contains(msg, "403") {
		t.Errorf("o 403 devia ser acusado; msg = %q", msg)
	}
}

// Comentário não é comportamento: um `// devolve 404` descreve o que a função FAZIA.
// Contá-lo faria o gate aprovar um contrato que o código já não cumpre.
func TestContrato_statusEmComentarioNaoConta(t *testing.T) {
	t.Run("CSDCN-B06: A status that lives only in a comment is not emitted", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
| 404 | não encontrado |
`
	codigo := `export const handler = async () => {
  // Antes devolvia { statusCode: 404 } aqui; agora unifica em 403 (anti-probing).
  if (!x) return { statusCode: 403, body: '{}' }
  return { statusCode: 200, body: '{}' }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Fail {
		t.Fatalf("quer Fail (403 emitido, 404 fantasma); obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "404") {
		t.Errorf("o 404 só existe em comentário — devia ser acusado como fantasma; msg = %q", msg)
	}
}

// Sem seção de contrato, o gate não tem o que confrontar — Skip, não Fail. Cobrar a
// existência da seção é trabalho do `spec-completa`.
func TestContrato_semSecaoEhSkip(t *testing.T) {
	t.Run("CSDCN-B07: Without the contract section there is nothing to confront", func(t *testing.T) {})
	if v, _ := rodaContrato(t, "## Efeitos\n| `X-B01` | faz algo |\n", "return { statusCode: 200 }"); v != Skip {
		t.Fatalf("spec sem Contrato de Saída devia ser Skip, obteve %v", v)
	}
}

// `void` é o contrato de um handler de cron/trigger: não devolve status. Se o código
// também não devolve, não há nada a confrontar.
func TestContrato_codigoSemStatusEhSkip(t *testing.T) {
	t.Run("CSDCN-B08: Code that returns no status is skipped", func(t *testing.T) {})
	spec := "## Contrato de Saída\n`void` (efeito: persiste itens).\n"
	codigo := "export const handler = async () => { await gravar() }"
	if v, _ := rodaContrato(t, spec, codigo); v != Skip {
		t.Fatalf("código sem status devia ser Skip, obteve %v", v)
	}
}

// FALSO-POSITIVO REALX, achado ao rodar o gate contra o app de referência antes de ligá-lo:
// `log-consent` define `fail(status, message)` e chama `fail(400, …)`. A primeira
// versão do gate só reconhecia `jsonResponse(`/`cors(`, então via apenas o 200 e
// acusava o 400 como fantasma. A forma genérica `helper(400, …)` cobre os helpers
// locais que cada handler define.
func TestContrato_helperLocalComStatusLiteralEhReconhecido(t *testing.T) {
	t.Run("CSDCN-B09: A literal status passed to a local helper counts as emitted", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | consentimento registrado |
| 400 | payload inválido |
| 5xx | falha |
`
	codigo := `function ok(body: unknown) {
  return { statusCode: 200, body: JSON.stringify(body) }
}
function fail(status: number, message: string) {
  return { statusCode: status, body: JSON.stringify({ error: message }) }
}
export const handler = async () => {
  if (!corpo) return fail(400, 'Invalid JSON')
  return ok({ done: true })
}`
	if v, msg := rodaContrato(t, spec, codigo); v != Pass {
		t.Fatalf("o 400 vem de `fail(400, …)` — devia passar; obteve %v (%s)", v, msg)
	}
}

// A contraparte: com um helper que recebe o status por PARÂMETRO, o gate não pode
// afirmar que um status declarado é fantasma — os valores passados podem incluí-lo
// por um caminho que a leitura textual não alcança. Ele segue cobrando o que
// encontrou de literal e se cala sobre o resto.
func TestContrato_comStatusDinamicoNaoAcusaFantasma(t *testing.T) {
	t.Run("CSDCN-B10: With a dynamic status the phantom side goes quiet and the literals still count", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
| 404 | não encontrado |
`
	codigo := `function resp(status: number) {
  return { statusCode: status, body: '{}' }
}
export const handler = async () => {
  if (!x) return resp(403)
  return { statusCode: 200, body: '{}' }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Fail {
		t.Fatalf("o 403 literal ainda devia ser cobrado; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "403") {
		t.Errorf("devia acusar o 403 emitido e não declarado; msg = %q", msg)
	}
	if strings.Contains(msg, "404") {
		t.Errorf("com status dinâmico o gate NÃO pode acusar fantasma; msg = %q", msg)
	}
}

// rodaContratoComDialeto é a mesma bancada, com um dialeto declarado — é o que
// prova que o gate é AGNÓSTICO: o léxico vem da config, não do gate.
func rodaContratoComDialeto(t *testing.T, spec, codigo, arquivo string, cfg *config.Config) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, arquivo), []byte(codigo), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "handler.spec.md", Kind: mapx.KindSpec},
			{ID: arquivo, Kind: mapx.KindCode},
		},
		Edges: []mapx.Edge{{From: "handler.spec.md", To: arquivo, Type: mapx.EdgeSpecifies}},
	}
	return checkContractStatusDeclared(spec, mapx.Node{ID: "handler.spec.md", Kind: mapx.KindSpec}, root, g, cfg)
}

// O GATE NÃO SABE JAVASCRIPT. Este teste é a régua da agnosticidade: o mesmo gate,
// sobre um handler Go `net/http`, com o léxico vindo de `dialect.family: go`.
// Se alguém voltar a embutir `statusCode:` na lógica, este teste reprova.
func TestContrato_dialetoGoLeStatusDoNetHTTP(t *testing.T) {
	t.Run("CSDCN-I01: The lexicon comes from the project's dialect, not from the gate", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
`
	// Go não tem `statusCode:` em lugar nenhum — o status sai de `WriteHeader`,
	// e por constante NOMEADA, que o gate traduz para número.
	codigo := `package handler

func Handle(w http.ResponseWriter, r *http.Request) {
	if !autorizado(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(200)
}`
	cfg := &config.Config{Dialect: &config.Dialect{Family: "go"}}
	v, msg := rodaContratoComDialeto(t, spec, codigo, "handler.go", cfg)
	if v != Fail {
		t.Fatalf("o 403 do `http.StatusForbidden` devia ser cobrado; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "403") {
		t.Errorf("devia traduzir StatusForbidden → 403; msg = %q", msg)
	}
}

// Dialeto declarado à mão (sem família): um projeto Rails/Python que ensina o gate
// a ler o seu próprio léxico. Prova que não é preciso ter família embutida.
func TestContrato_dialetoExplicitoEnsinaLexicoProprio(t *testing.T) {
	t.Run("CSDCN-I02: A dialect declared by hand teaches the gate its own lexicon", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
`
	codigo := `def create
  return render_error(422) unless valido?
  render json: {}, status: 200
end`
	cfg := &config.Config{Dialect: &config.Dialect{
		HTTPStatus: `status:\s*(\d{3})|render_error\(\s*(\d{3})\s*\)`,
	}}
	v, msg := rodaContratoComDialeto(t, spec, codigo, "handler.rb", cfg)
	if v != Fail {
		t.Fatalf("o 422 devia ser cobrado pelo léxico declarado; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "422") {
		t.Errorf("devia acusar o 422; msg = %q", msg)
	}
}

func TestContrato_semDialetoEhPendente(t *testing.T) {
	t.Run("CSDCN-B11: Without a declared dialect the verdict is Pending", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
`
	codigo := `export const handler = async () => ({ statusCode: 200 })`
	// Sem dialeto configurado (cfg nil ou Dialect vazio) -> Pending
	v, _ := rodaContratoComDialeto(t, spec, codigo, "handler.ts", nil)
	if v != Pending {
		t.Fatalf("sem dialeto configurado devia ser Pending; obteve %v", v)
	}
}

func TestContrato_optOutEhSkip(t *testing.T) {
	t.Run("CSDCN-B12: An explicit opt-out of the http_status field is honoured", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
`
	codigo := `export const handler = async () => ({ statusCode: 200 })`
	cfg := &config.Config{
		Dialect: &config.Dialect{
			OptOut: []string{"http_status"},
		},
	}
	v, _ := rodaContratoComDialeto(t, spec, codigo, "handler.ts", cfg)
	if v != Skip {
		t.Fatalf("com opt_out de http_status devia ser Skip; obteve %v", v)
	}
}

func TestContrato_secaoOutputContractEmIngles(t *testing.T) {
	spec := `## Output Contract
| Status | When |
| --- | --- |
| 200 | success |
| 404 | not found |
`
	codigo := `export const handler = async () => {
  if (!item) return { statusCode: 404, body: '{}' }
  return { statusCode: 200, body: '{}' }
}`
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts"}}
	v, msg := rodaContratoComDialeto(t, spec, codigo, "handler.ts", cfg)
	if v != Pass {
		t.Fatalf("seção ## Output Contract devia passar; obteve %v (%s)", v, msg)
	}
}

// A CONSTANTE NOMEADA vale o número que significa. Um handler escrito inteiramente com
// `http.StatusForbidden` não tem um único dígito no corpo — ler só dígitos aprovaria
// todos eles em silêncio.
func TestContrato_constanteNomeadaValeONumero(t *testing.T) {
	t.Run("CSDCN-I03: A named constant is worth the number it means", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
`
	codigo := `package handler

func Handle(w http.ResponseWriter, r *http.Request) {
	if !dono(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}`
	cfg := &config.Config{Dialect: &config.Dialect{Family: "go"}}
	v, msg := rodaContratoComDialeto(t, spec, codigo, "handler.go", cfg)
	if v != Fail {
		t.Fatalf("o 403 escrito só como constante devia ser cobrado; obteve %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "403") {
		t.Errorf("devia traduzir StatusForbidden → 403; msg = %q", msg)
	}
	if strings.Contains(msg, "StatusOK") || strings.Contains(msg, "StatusForbidden") {
		t.Errorf("a acusação é pelo NÚMERO, não pelo nome da constante; msg = %q", msg)
	}
}

// O gate NÃO EXIGE a faixa genérica. Uma tabela que declara só o 200, sobre um handler
// cujo único outro caminho é o 500 do catch, passa: tratar a faixa como status exigiria
// adivinhar quais números ela cobre, e a adivinhação seria cobrada como declaração.
func TestContrato_naoExigeFaixaGenerica(t *testing.T) {
	t.Run("CSDCN-X01: The gate does not demand the generic ranges", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
`
	codigo := `export const handler = async () => {
  try {
    return { statusCode: 200, body: '{}' }
  } catch {
    return { statusCode: 500, body: '{}' }
  }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Pass {
		t.Fatalf("sem faixa declarada e só com o 500 do catch devia passar; obteve %v (%s)", v, msg)
	}
}

// A contraparte do dinâmico, pelo lado do que o gate SE CALA: o status declarado que
// nenhum literal mostra não pode ser acusado quando o envelope recebe o código por
// parâmetro. Aqui não há nenhum literal a cobrar — o veredito tem de ser limpo.
func TestContrato_dinamicoSilenciaOFantasmaPorCompleto(t *testing.T) {
	t.Run("CSDCN-X03: The gate does not charge the phantom side under a dynamic status", func(t *testing.T) {})
	spec := `## Contrato de Saída
| Status | Quando |
| --- | --- |
| 200 | ok |
| 409 | conflito |
`
	codigo := `function resp(status: number) {
  return { statusCode: status, body: '{}' }
}
export const handler = async () => {
  if (conflito) return resp(codigoDeConflito)
  return { statusCode: 200, body: '{}' }
}`
	v, msg := rodaContrato(t, spec, codigo)
	if v != Pass {
		t.Fatalf("com status dinâmico e nenhum literal a cobrar, o veredito devia ser limpo; obteve %v (%s)", v, msg)
	}
	if strings.Contains(msg, "409") {
		t.Errorf("o 409 pode ser emitido por um caminho que a leitura textual não alcança; msg = %q", msg)
	}
}
