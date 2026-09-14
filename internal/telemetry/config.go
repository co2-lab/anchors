package telemetry

import (
	"os"
	"strings"
)

// --- O OPT-OUT, e por que ele é fácil de encontrar ---
//
// A telemetry é LIGADA por padrão, e a decisão é do mantenedor. O custo dessa escolha é
// conhecido: quem descobre depois que um binário mandava dados sem ter perguntado perde a
// confiança na ferramenta inteira — e o Anchors não fica ao lado do trabalho, ele barra
// commit e move card.
//
// O que compensa esse custo é o desligamento ser TRIVIAL e ANUNCIADO. Três caminhos, e o
// primeiro funciona sem editar arquivo nenhum:
//
//	ANCHORS_TELEMETRY=off          variável de ambiente (vale para tudo, inclusive CI)
//	telemetry: off                 no `anchors.yaml` (vale para o projeto, versionado)
//	anchors telemetry off          grava em `.anchors/settings.yaml` (vale para a máquina)
//
// E o AVISO aparece duas vezes: no `anchors init`, e na PRIMEIRA execução de qualquer
// comando numa máquina que ainda não viu o aviso. A segunda é o que alcança quem instalou
// com `go install` num projeto que outra pessoa configurou — que é o caso comum num time.

// EnvVar desliga sem editar arquivo. `off`, `0`, `false` e `no` funcionam: quem
// procura como desligar tenta uma dessas, e exigir a palavra exata seria uma pegadinha.
const EnvVar = "ANCHORS_TELEMETRY"

// Config é o que decide se e para onde emitir.
type Config struct {
	// Enabled: o padrão é true. Ver o bloco acima sobre a escolha.
	Enabled bool
	// Endpoint OTLP/HTTP. Vazio usa o padrão do produto.
	Endpoint string
	// Headers de autenticação, do ambiente — nunca do arquivo versionado, que iria para
	// o repositório de quem usa.
	Headers map[string]string
	// NoCodes remove o código de identidade das unidades dos eventos.
	//
	// O código (`RLSGR`, `ALFDL`) é vocabulário do PROJETO de quem usa, e a sequência
	// deles desenha a arquitetura dele. Para achar caso anômalo os números bastam — o
	// código só ajuda a investigar um caso específico, e quem quer esse nível send o
	// relatório local.
	NoCodes bool
}

// Disabled responde se o ambiente ou a configuração desligou a emissão.
//
// A ORDEM importa: o ambiente vence o arquivo. Quem roda um comando numa máquina de CI
// precisa conseguir desligar sem commitar — e commitar para desligar telemetry seria
// pedir que a decisão de uma pessoa vire mudança no repositório do time.
func Disabled(doArquivo string) bool {
	if v := strings.ToLower(strings.TrimSpace(os.Getenv(EnvVar))); v != "" {
		return v == "off" || v == "0" || v == "false" || v == "no"
	}
	v := strings.ToLower(strings.TrimSpace(doArquivo))
	return v == "off" || v == "0" || v == "false" || v == "no"
}
