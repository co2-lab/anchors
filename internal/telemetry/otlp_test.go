package telemetry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func fixedClock() time.Time { return time.Date(2026, 9, 14, 18, 30, 0, 0, time.UTC) }

// O CORPO precisa ser OTLP válido — e a única forma de saber é montá-lo e conferir a forma
// que o protocolo exige. Um body malformado é aceito com 200 por alguns coletores e
// descartado em silêncio: o defeito só apareceria quando ninguém achasse os dados.
func TestEmite_produzOTLPValido(t *testing.T) {
	var recebido []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recebido, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	body, err := e.build(New(ClaimServed, map[string]any{
		"candidatos": 14, "pulados": 2, "motivo": "pr-aberto",
	}, fixedClock))
	if err != nil {
		t.Fatal(err)
	}

	var p map[string]any
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatalf("o body não é JSON: %v", err)
	}
	rl, ok := p["resourceLogs"].([]any)
	if !ok || len(rl) != 1 {
		t.Fatalf("OTLP exige `resourceLogs`; veio %v", p)
	}
	_ = recebido
}

// O RECURSO identifica o PRODUTO, nunca quem o usa. Sem isto, o primeiro atributo que
// alguém acrescentasse por conveniência — hostname, usuário, repositório — entraria sem
// nada acusar.
func TestEmite_oRecursoNaoIdentificaQuemUsa(t *testing.T) {
	e := NewEmitter(Config{Enabled: true}, "0.1.99")
	body, _ := e.build(New(CheckFinished, nil, fixedClock))
	s := string(body)

	if !strings.Contains(s, `"anchors"`) || !strings.Contains(s, `"0.1.99"`) {
		t.Error("o recurso deveria identificar o produto e a versão")
	}
	for _, proibido := range []string{"hostname", "user", "repo", "host.name", "user.name"} {
		if strings.Contains(s, proibido) {
			t.Errorf("o recurso não pode carregar %q — ele identifica quem USA, não o produto", proibido)
		}
	}
}

// TELEMETRIA NUNCA BLOQUEIA. Um endpoint fora do ar não pode travar o `anchors check` de
// ninguém — perder um evento é infinitamente melhor.
func TestEmite_naoBloqueiaNemFalhaComServidorMorto(t *testing.T) {
	e := NewEmitter(Config{Enabled: true, Endpoint: "http://127.0.0.1:1"}, "0.1.99")
	feito := make(chan struct{})
	go func() {
		e.Emit(New(TurnEnded, map[string]any{"estado": "in-progress"}, fixedClock))
		close(feito)
	}()
	select {
	case <-feito:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("`Emit` bloqueou — telemetry não pode atrasar o trabalho de ninguém")
	}
}

// DESLIGADA devolve nil, e o nil aceita `Emit` sem fazer nada. É o que permite ao chamador
// não precisar de `if` em torno de cada chamada — e um `if` esquecido seria um panic.
func TestNovoEmissor_desligadaDevolveNilQueAceitaEmite(t *testing.T) {
	e := NewEmitter(Config{Enabled: false}, "0.1.99")
	if e != nil {
		t.Fatal("desligada, o emitter deveria ser nil")
	}
	e.Emit(New(ClaimEmpty, nil, fixedClock)) // não pode entrar em panic
}

// O VALOR só aceita número, booleano e string de vocabulário. Qualquer outra coisa vira a
// DESCRIÇÃO DO TIPO — é a guarda que impede um `error` com caminho de arquivo, ou um `map`
// com conteúdo, de vazar numa chamada futura escrita sem atenção.
func TestValorOTLP_oQueNaoEhVocabularioViraOTipo(t *testing.T) {
	v := otlpValue(map[string]string{"caminho": "/home/alguem/projeto/spec.md"})
	s, _ := v["stringValue"].(string)
	if strings.Contains(s, "/home/alguem") {
		t.Errorf("um valor não previsto vazou conteúdo: %q", s)
	}
	if s != "(map[string]string)" {
		t.Errorf("esperava a descrição do tipo, veio %q", s)
	}
	// E os previstos passam.
	if otlpValue(14)["intValue"] != "14" {
		t.Error("número deveria virar intValue")
	}
	if otlpValue("pr-aberto")["stringValue"] != "pr-aberto" {
		t.Error("vocabulário do produto deveria passar")
	}
}

// O FLUSH é o que faz o evento CHEGAR num CLI.
//
// `Emit` dispara o POST numa goroutine para não atrasar o comando — e num programa de linha
// de comando isso significa que o processo morre antes de o envio sair.
//
// O defeito é INVISÍVEL em teste de unidade: ali o processo continua vivo depois da
// chamada, e o POST acaba saindo. Só apareceu ao rodar o binário de verdade contra um
// coletor local — que não recebeu NADA até este mecanismo existir.
func TestFlush_esperaOEnvioAntesDeOProcessoSair(t *testing.T) {
	chegou := make(chan struct{}, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// O atraso simula o que acontece de verdade: o POST não é instantâneo, e sem o
		// flush o processo morreria neste intervalo.
		time.Sleep(50 * time.Millisecond)
		chegou <- struct{}{}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	e.Emit(New(TurnEnded, map[string]any{"estado": "in-progress"}, fixedClock))
	e.Flush()

	select {
	case <-chegou:
	default:
		t.Fatal("o `Flush` voltou antes de o envio chegar — num CLI o evento se perderia")
	}
}

// O FLUSH TAMBÉM NÃO PODE BLOQUEAR. Um servidor que não responde não pode atrasar a saída
// de um comando — se ele demora, o evento se perde, e é o desfecho certo.
func TestFlush_naoEsperaParaSempre(t *testing.T) {
	// O handler espera um sinal que só chega no fim do teste. `time.Sleep(10s)` seria o
	// mesmo do ponto de vista do cliente, e faria o TESTE durar dez segundos — o
	// `srv.Close()` espera as conexões abertas.
	solta := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-solta
	}))
	defer func() { close(solta); srv.Close() }()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	e.Emit(New(CheckFinished, nil, fixedClock))

	inicio := time.Now()
	e.Flush()
	if d := time.Since(inicio); d > sendDeadline+time.Second {
		t.Errorf("o `Flush` esperou %v — telemetria não pode atrasar a saída de um comando", d)
	}
}

// O PRAZO DO FLUSH é uma defesa PRÓPRIA, e não o timeout do cliente HTTP.
//
// São duas, e a distinção só apareceu na bateria de mutação: remover o `time.After` do
// `select` deixou o teste acima verde, porque o `http.Client` tem timeout e devolve o
// controle sozinho. A defesa externa escondia a interna.
//
// Por que a interna importa: um envio que não passa pelo cliente HTTP — uma goroutine presa
// em qualquer outro ponto, hoje ou numa mudança futura — travaria a saída do processo para
// sempre. O `Flush` não pode depender de outra camada ter prazo.
func TestFlush_temPrazoProprioIndependenteDoClienteHTTP(t *testing.T) {
	e := NewEmitter(Config{Enabled: true, Endpoint: "http://exemplo.invalido"}, "0.1.99")

	// Um envio que NUNCA termina, e que não passa pelo cliente HTTP: é o que separa a
	// defesa do `Flush` da defesa do `http.Client`.
	e.emVoo.Add(1)
	defer e.emVoo.Done() // solta ao fim do teste, para não vazar a goroutine

	// O `Flush` roda em goroutine com um limite PRÓPRIO do teste: sem o prazo interno ele
	// espera para sempre, e um teste que trava não falha — ele estoura o tempo do pacote
	// inteiro e esconde qual asserção quebrou.
	voltou := make(chan time.Duration, 1)
	go func() {
		inicio := time.Now()
		e.Flush()
		voltou <- time.Since(inicio)
	}()

	select {
	case d := <-voltou:
		if d > sendDeadline+500*time.Millisecond {
			t.Errorf("o `Flush` esperou %v sem o cliente HTTP no caminho", d)
		}
	case <-time.After(sendDeadline + 2*time.Second):
		t.Fatal("o `Flush` não voltou — ele precisa do prazo PRÓPRIO, senão uma goroutine " +
			"presa trava a saída do processo para sempre")
	}
}

// Nil aceita `Flush` — é o que permite ao `main` chamá-lo sem saber se a telemetria está
// ligada, e um `if` esquecido ali seria um panic na saída do processo.
func TestFlush_nilNaoEntraEmPanico(t *testing.T) {
	var e *Emitter
	e.Flush()
}
