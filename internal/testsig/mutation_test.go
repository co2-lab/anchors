package testsig

import (
	"os"
	"path/filepath"
	"testing"
)

// Fixture no formato Mutation Testing Elements — o mesmo JSON que Stryker, PIT,
// Infection e mutmut emitem. Cobre cada status que muda a conta.
const fixtureMT = `{
  "schemaVersion": "1.0",
  "files": {
    "src/business-logic/pricing.ts": {
      "language": "typescript",
      "source": "…",
      "mutants": [
        {"id":"1","status":"Killed","mutatorName":"ConditionalExpression","location":{"start":{"line":10,"column":1},"end":{"line":10,"column":9}}},
        {"id":"2","status":"Killed","mutatorName":"ArithmeticOperator","location":{"start":{"line":11,"column":1},"end":{"line":11,"column":9}}},
        {"id":"3","status":"Timeout","mutatorName":"EqualityOperator","location":{"start":{"line":12,"column":1},"end":{"line":12,"column":9}}},
        {"id":"4","status":"Survived","mutatorName":"BooleanLiteral","location":{"start":{"line":42,"column":1},"end":{"line":42,"column":9}}},
        {"id":"5","status":"NoCoverage","mutatorName":"StringLiteral","location":{"start":{"line":88,"column":1},"end":{"line":88,"column":9}}},
        {"id":"6","status":"CompileError","mutatorName":"BlockStatement","location":{"start":{"line":99,"column":1},"end":{"line":99,"column":9}}}
      ]
    }
  }
}`

func TestParseMutation(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "mutation.json")
	if err := os.WriteFile(p, []byte(fixtureMT), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm, ok := rep.Files["src/business-logic/pricing.ts"]
	if !ok {
		t.Fatalf("arquivo não encontrado; chaves = %v", keys(rep))
	}

	// Timeout conta como MORTE: o teste percebeu a mutação (travou por causa dela).
	if fm.Killed != 3 {
		t.Errorf("killed = %d, queria 3 (2 Killed + 1 Timeout)", fm.Killed)
	}
	// SOBREVIVENTE é só o que o teste EXECUTOU e não percebeu. O NoCoverage já contou
	// como sobrevivente aqui, e a decisão mudou em 25/08: "existe teste que execute esta
	// linha?" é pergunta do gate de COBERTURA. Somando os dois, o gate de mutação passava
	// a reportar 187 arquivos (medido no app de referência) que ele não é dono de resolver — e no meio
	// deles se perdiam os zeros REAIS, em que o teste roda e não verifica nada.
	if fm.Survived != 1 {
		t.Errorf("survived = %d, queria 1 — só o Survived de verdade", fm.Survived)
	}
	if fm.NoCoverage != 1 {
		t.Errorf("noCoverage = %d, queria 1 — contado à parte, fora do score", fm.NoCoverage)
	}
	// CompileError fica FORA do denominador: não diz nada sobre a qualidade do teste.
	if want := 75.0; fm.Score != want {
		t.Errorf("score = %.1f, queria %.1f (3 mortos de 4 EXECUTADOS)", fm.Score, want)
	}
	// As linhas dos sobreviventes são o que torna o score acionável — e só entram as que
	// um teste de fato alcançou. A linha 88, sem cobertura, é assunto do outro gate.
	if len(fm.SurvivedAt) != 1 || fm.SurvivedAt[0] != 42 {
		t.Errorf("linhas sobreviventes = %v, queria [42]", fm.SurvivedAt)
	}
}

func TestParseMutationRecusaLixo(t *testing.T) {
	dir := t.TempDir()
	casos := map[string]string{
		"não é JSON":                    "isto não é json",
		"JSON válido sem a chave files": `{"schemaVersion":"1.0"}`,
	}
	for nome, conteúdo := range casos {
		t.Run(nome, func(t *testing.T) {
			p := filepath.Join(dir, "x.json")
			os.WriteFile(p, []byte(conteúdo), 0o644)
			if _, err := ParseMutation(p); err == nil {
				t.Fatal("aceitou entrada inválida em silêncio — o pior desfecho: score fantasma")
			}
		})
	}
}

func keys(r *MutationReport) []string {
	var out []string
	for k := range r.Files {
		out = append(out, k)
	}
	return out
}

// Fixture de um arquivo cujos mutantes foram TODOS ignorados pela ferramenta — o caso
// da tabela de constantes com `ignoreStatic` ligado. Não é teste faltando: é o
// instrumentador dizendo que ali não há experimento a fazer.
const fixtureTudoIgnorado = `{
  "schemaVersion": "1.0",
  "files": {
    "src/constants/categories.ts": {
      "language": "typescript",
      "source": "…",
      "mutants": [
        {"id":"1","status":"Ignored","mutatorName":"StringLiteral","location":{"start":{"line":6,"column":1},"end":{"line":6,"column":9}}},
        {"id":"2","status":"Ignored","mutatorName":"ObjectLiteral","location":{"start":{"line":7,"column":1},"end":{"line":7,"column":9}}}
      ]
    }
  }
}`

// Um arquivo sem mutante SOBREVIVENTE tem score 100 — inclusive quando o motivo é que
// nenhum chegou a rodar. É a definição do número: 100% quer dizer "nada sobreviveu".
//
// Antes disto o score ficava ausente e o arquivo caía num limbo: o sinal era gravado
// vazio e o gate pedia "rode a ferramenta de mutação" para um arquivo em que ela já
// tinha rodado e ignorado tudo corretamente — um pedido impossível de atender, porque
// rodar de novo daria o mesmo resultado.
func TestParseMutationTudoIgnoradoDaCem(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "mutation.json")
	if err := os.WriteFile(p, []byte(fixtureTudoIgnorado), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["src/constants/categories.ts"]

	if fm.Score != 100 {
		t.Errorf("score = %.1f, queria 100 (nenhum mutante sobreviveu)", fm.Score)
	}
	if fm.Killed != 0 || fm.Survived != 0 {
		t.Errorf("killed/survived = %d/%d, queria 0/0 — ignorado não é morto nem vivo",
			fm.Killed, fm.Survived)
	}
	// O contador é o que impede o 100 de ser ambíguo numa lista ordenada: sem ele, este
	// arquivo e um com todos os mutantes MORTOS ficariam indistinguíveis.
	if fm.Ignored != 2 {
		t.Errorf("ignored = %d, queria 2", fm.Ignored)
	}
}

// E o ignorado não entra no score de quem TEM medida — ele fica fora do experimento nos
// dois sentidos. Sem esta garantia, um `disable` numa linha inflaria o número da unidade.
func TestParseMutationIgnoradoNaoInflaScore(t *testing.T) {
	misto := `{"schemaVersion":"1.0","files":{"a.ts":{"mutants":[
	  {"id":"1","status":"Killed","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}}},
	  {"id":"2","status":"Survived","location":{"start":{"line":2,"column":1},"end":{"line":2,"column":2}}},
	  {"id":"3","status":"Ignored","location":{"start":{"line":3,"column":1},"end":{"line":3,"column":2}}},
	  {"id":"4","status":"Ignored","location":{"start":{"line":4,"column":1},"end":{"line":4,"column":2}}}
	]}}}`
	dir := t.TempDir()
	p := filepath.Join(dir, "m.json")
	if err := os.WriteFile(p, []byte(misto), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["a.ts"]
	if fm.Score != 50 {
		t.Errorf("score = %.1f, queria 50 (1 morto de 2 que rodaram) — os 2 ignorados "+
			"não podem entrar no denominador NEM no numerador", fm.Score)
	}
	if fm.Ignored != 2 {
		t.Errorf("ignored = %d, queria 2", fm.Ignored)
	}
}

// Arquivo que NENHUM teste executa: todo mutante é NoCoverage. O gate de mutação não é
// dono disso — quem cobra "não há teste que execute esta linha" é o de cobertura —, então
// o score é 100 e o número fica gravado à parte.
//
// MEDIDO no app de referência em 25/08, e é o caso comum, não a exceção: `repositories/transactions.ts`
// dá 175 NoCoverage e zero sobrevivente, porque a prova dele é de INTEGRAÇÃO, que a config
// de mutação exclui de propósito (teste que fala com serviço real não distingue "quebrei a
// regra" de "a rede oscilou").
func TestParseMutationSemCoberturaNaoEhSobrevivente(t *testing.T) {
	semCobertura := `{"schemaVersion":"1.0","files":{"repo.ts":{"mutants":[
	  {"id":"1","status":"NoCoverage","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}}},
	  {"id":"2","status":"NoCoverage","location":{"start":{"line":2,"column":1},"end":{"line":2,"column":2}}},
	  {"id":"3","status":"NoCoverage","location":{"start":{"line":3,"column":1},"end":{"line":3,"column":2}}}
	]}}}`
	dir := t.TempDir()
	p := filepath.Join(dir, "m.json")
	if err := os.WriteFile(p, []byte(semCobertura), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["repo.ts"]

	if fm.Survived != 0 {
		t.Errorf("survived = %d, queria 0 — sem cobertura não é sobrevivente", fm.Survived)
	}
	if fm.NoCoverage != 3 {
		t.Errorf("noCoverage = %d, queria 3", fm.NoCoverage)
	}
	if fm.Score != 100 {
		t.Errorf("score = %.1f, queria 100 — nenhum mutante executado, nada sobreviveu",
			fm.Score)
	}
}

// E o contraste que dá sentido à separação: quando parte do arquivo É executada, os
// mutantes que sobrevivem ALI continuam sendo achado deste gate. O caso real que motivou
// isto foi o `models/holidays.ts`: 32 sem cobertura e 15 sobreviventes de verdade — hoje
// ele se perdia no meio de 187 zeros, e o número passa a dizer "o teste roda e não
// verifica", que é acionável.
func TestParseMutationParcialContaSoOExecutado(t *testing.T) {
	parcial := `{"schemaVersion":"1.0","files":{"h.ts":{"mutants":[
	  {"id":"1","status":"NoCoverage","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}}},
	  {"id":"2","status":"NoCoverage","location":{"start":{"line":2,"column":1},"end":{"line":2,"column":2}}},
	  {"id":"3","status":"Survived","location":{"start":{"line":3,"column":1},"end":{"line":3,"column":2}}},
	  {"id":"4","status":"Killed","location":{"start":{"line":4,"column":1},"end":{"line":4,"column":2}}}
	]}}}`
	dir := t.TempDir()
	p := filepath.Join(dir, "m.json")
	if err := os.WriteFile(p, []byte(parcial), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["h.ts"]

	if fm.Score != 50 {
		t.Errorf("score = %.1f, queria 50 (1 morto de 2 EXECUTADOS) — os 2 sem cobertura "+
			"não podem entrar no denominador", fm.Score)
	}
	if fm.NoCoverage != 2 || fm.Survived != 1 || fm.Killed != 1 {
		t.Errorf("noCov/surv/killed = %d/%d/%d, queria 2/1/1",
			fm.NoCoverage, fm.Survived, fm.Killed)
	}
}
