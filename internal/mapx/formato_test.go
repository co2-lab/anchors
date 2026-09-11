package mapx

import (
	"strings"
	"testing"
)

// O CONTRATO: o `version:` do mapa diz em que FORMATO o arquivo está, e o binário recusa o
// que não sabe ler.
//
// Sem isso, um binário que encontra formato desconhecido interpreta o que reconhece, ignora
// o resto, e a próxima gravação escreve só o que sobrou — perdendo dado sem nada acusar. O
// exemplo concreto é o `julgamentos` renomeado: os carimbos de julgamento de IA evaporam e
// o `check` recobra o que alguém já respondeu.

func TestConfereFormato_aceitaSoAFaixaQueSabeLer(t *testing.T) {
	if err := ConfereFormato("m.yaml", FormatoAtual); err != nil {
		t.Errorf("o formato que este binário escreve tem de ser legível: %v", err)
	}
	if err := ConfereFormato("m.yaml", FormatoMinimoLegivel); err != nil {
		t.Errorf("o mínimo legível tem de ser aceito: %v", err)
	}
}

// Um mapa do FUTURO pede atualização do binário, e a mensagem tem de dizer o RISCO —
// senão a saída barata é apagar o arquivo e reconstruir, que perde os julgamentos.
func TestConfereFormato_formatoDoFuturoPedeAtualizacao(t *testing.T) {
	err := ConfereFormato("m.yaml", FormatoAtual+1)
	if err == nil {
		t.Fatal("um formato acima do que o binário escreve tem de ser recusado")
	}
	msg := err.Error()
	for _, quer := range []string{"mais NOVO", "silêncio", "julgamento", "upgrade"} {
		if !strings.Contains(msg, quer) {
			t.Errorf("a mensagem deveria conter %q; veio:\n%s", quer, msg)
		}
	}
}

// Um mapa do PASSADO pede MIGRAÇÃO, e a mensagem nomeia o comando. Um erro que diz "rode X"
// sem que X exista é pior que nenhum: quem lê tenta, falha, e desconfia da próxima
// mensagem.
func TestConfereFormato_formatoAntigoPedeMigracao(t *testing.T) {
	err := ConfereFormato("m.yaml", 1)
	if err == nil {
		t.Fatal("o formato 1 tem de ser recusado — ele é migrado, não lido")
	}
	if !strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("a mensagem deveria nomear o comando que conserta; veio:\n%s", err)
	}
}

// Mapa SEM `version:` é de antes de o campo existir — formato 1, e migrável. Tratá-lo como
// "0" e recusar com a mensagem do futuro mandaria a pessoa atualizar um binário que já é o
// mais novo.
func TestConfereFormato_semVersionEhFormatoUm(t *testing.T) {
	err := ConfereFormato("m.yaml", 0)
	if err == nil {
		t.Fatal("sem `version:` o mapa é formato 1, e precisa migrar")
	}
	if !strings.Contains(err.Error(), "migrate") {
		t.Errorf("deveria pedir MIGRAÇÃO, não atualização; veio:\n%s", err)
	}
	// A mensagem tem de dizer FORMATO 1, e não "formato 0".
	//
	// Sem a normalização, o 0 ainda cai na faixa "abaixo do mínimo" e produz a mensagem
	// certa pelo motivo errado — e quem lê "está no formato 0" procura um arquivo
	// corrompido, não um arquivo antigo. A mutação que remove a normalização passava com
	// a primeira versão deste teste.
	if !strings.Contains(err.Error(), "formato 1") {
		t.Errorf("a mensagem deveria dizer FORMATO 1 (o arquivo é antigo, não corrompido); veio:\n%s", err)
	}
	if strings.Contains(err.Error(), "formato 0") {
		t.Error("`formato 0` não existe: um mapa sem `version:` é o formato 1")
	}
}

// As duas mensagens são DISTINTAS, e a distinção é o que manda a pessoa para o lado certo:
// o futuro pede binário novo, o passado pede migração. Um "erro ao carregar o mapa" para os
// dois faria procurar corrupção onde há só versão.
func TestConfereFormato_asDuasMensagensNaoSeConfundem(t *testing.T) {
	futuro := ConfereFormato("m.yaml", FormatoAtual+1).Error()
	passado := ConfereFormato("m.yaml", 1).Error()

	if strings.Contains(futuro, "anchors migrate") {
		t.Error("o mapa do futuro não se conserta migrando — ele pede binário novo")
	}
	if strings.Contains(passado, "mais NOVO") {
		t.Error("o mapa antigo não foi gravado por binário mais novo")
	}
}
