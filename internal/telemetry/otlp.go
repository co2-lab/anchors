package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// --- OTLP/HTTP escrito à mão, e por quê ---
//
// O SDK oficial do OpenTelemetry traz ~40 módulos. O Anchors tem SEIS dependências diretas,
// e isso não é acaso: ele instala com um `go install` e roda em máquina de agente, de CI e
// de quem só quer experimentar.
//
// O protocolo OTLP/HTTP+JSON é um POST com um body documentado. São as ~80 linhas abaixo,
// e o que se ganha em troca é o destino ser TROCÁVEL: a mesma emissão vai para Honeycomb,
// Grafana, Elastic ou um coletor local — muda a URL, não o código.
//
// O que NÃO se ganha: retry sofisticado, batching adaptativo, compressão. Nada disso vale
// para o volume aqui (poucos eventos por comando), e cada um seria mais código para manter.

// DefaultEndpoint é o Honeycomb. Trocável por config — o protocolo é o mesmo em todos.
const DefaultEndpoint = "https://api.honeycomb.io/v1/logs"

// sendDeadline é curto DE PROPÓSITO.
//
// Telemetria nunca pode atrasar o trabalho. Se o servidor demora, o evento se perde — e
// perder um evento é infinitamente melhor que o `anchors check` de alguém travar por causa
// de um endpoint fora do ar.
const sendDeadline = 2 * time.Second

// Emitter send os eventos. Nil-safe: um `*Emitter` nulo aceita `Emit` sem fazer nada, e é
// o que permite ao chamador não precisar de `if` em torno de cada chamada.
type Emitter struct {
	cfg    Config
	client *http.Client
	// service identifica o PRODUTO, não quem o usa: "anchors" e a versão.
	version string
	// emVoo conta os envios que ainda não terminaram.
	//
	// Sem isto o evento se perde SEMPRE num CLI: `Emit` dispara uma goroutine, o comando
	// termina, e o processo morre antes de o POST sair. Medido — o coletor local não
	// recebeu nada até este campo existir, e o defeito é invisível em teste de unidade,
	// onde o processo continua vivo depois da chamada.
	emVoo sync.WaitGroup
}

// NewEmitter devolve nil quando a telemetry está desligada.
//
// Nil e não um emitter inerte: um objeto que existe e não faz nada é o tipo de coisa que
// alguém liga por engano numa refatoração. O nil é uma afirmação.
func NewEmitter(cfg Config, version string) *Emitter {
	if !cfg.Enabled {
		return nil
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = DefaultEndpoint
	}
	return &Emitter{
		cfg:     cfg,
		version: version,
		client:  &http.Client{Timeout: sendDeadline},
	}
}

// Emit send um evento. NUNCA devolve erro, e nunca bloqueia o chamador.
//
// As duas coisas são a mesma decisão: telemetry que falha não pode virar problema de quem
// está trabalhando. O comando segue; o evento se perde em silêncio.
func (e *Emitter) Emit(ev Event) {
	if e == nil {
		return
	}
	body, err := e.build(ev)
	if err != nil {
		return
	}
	e.emVoo.Add(1)
	go func() {
		defer e.emVoo.Done()
		e.send(body)
	}()
}

// Flush espera os envios em voo, até o prazo.
//
// Chamado quando o processo vai terminar. O prazo é o mesmo do envio: telemetria não pode
// atrasar a saída de um comando — se o servidor demora, o evento se perde, e perder um
// evento é infinitamente melhor que fazer alguém esperar por ele.
func (e *Emitter) Flush() {
	if e == nil {
		return
	}
	pronto := make(chan struct{})
	go func() { e.emVoo.Wait(); close(pronto) }()
	select {
	case <-pronto:
	case <-time.After(sendDeadline):
	}
}

func (e *Emitter) send(body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), sendDeadline)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.cfg.Endpoint, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range e.cfg.Headers {
		req.Header.Set(k, v)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return
	}
	_ = resp.Body.Close()
}

// build produz o body OTLP/HTTP para UM registro de log.
//
// Logs e não spans: os eventos do Anchors são fatos pontuais ("o claim serviu o card 433"),
// não intervalos com início e fim. Forçá-los em span exigiria inventar duração.
func (e *Emitter) build(ev Event) ([]byte, error) {
	attrs := make([]map[string]any, 0, len(ev.Attrs))
	for k, v := range ev.Attrs {
		attrs = append(attrs, map[string]any{"key": k, "value": otlpValue(v)})
	}
	payload := map[string]any{
		"resourceLogs": []map[string]any{{
			"resource": map[string]any{
				// O RECURSO identifica o PRODUTO e a versão — nunca quem o usa. Nada de
				// hostname, usuário ou repositório.
				"attributes": []map[string]any{
					{"key": "service.name", "value": map[string]any{"stringValue": "anchors"}},
					{"key": "service.version", "value": map[string]any{"stringValue": e.version}},
				},
			},
			"scopeLogs": []map[string]any{{
				"logRecords": []map[string]any{{
					"timeUnixNano": fmt.Sprint(ev.At.UnixNano()),
					"body":         map[string]any{"stringValue": string(ev.Name)},
					"attributes":   attrs,
				}},
			}},
		}},
	}
	return json.Marshal(payload)
}

// otlpValue converte um valor Go para o formato do protocolo.
//
// Só três tipos são aceitos, e a lista é curta de propósito: números, booleanos, e string
// de VOCABULÁRIO (nome de gate, estado de card). Qualquer outra coisa vira a descrição do
// tipo, nunca o conteúdo — é a guarda que impede um `map` ou um `error` com caminho de
// arquivo de vazar por descuido numa chamada futura.
func otlpValue(v any) map[string]any {
	switch t := v.(type) {
	case string:
		return map[string]any{"stringValue": t}
	case int:
		return map[string]any{"intValue": fmt.Sprint(t)}
	case int64:
		return map[string]any{"intValue": fmt.Sprint(t)}
	case bool:
		return map[string]any{"boolValue": t}
	default:
		return map[string]any{"stringValue": fmt.Sprintf("(%T)", v)}
	}
}
