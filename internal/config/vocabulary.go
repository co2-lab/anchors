package config

// --- o VOCABULÁRIO em inglês ---
//
// Os nomes de gate são IDENTIFICADORES: eles vão para o `anchors.yaml` de cada projeto,
// para tutoriais, para respostas de fórum. São fixos em inglês e NÃO se traduzem — um
// `anchors.yaml` escrito por um time brasileiro tem de funcionar num time espanhol sem
// tradução nenhuma.
//
// A TABELA DE ALIAS QUE VIVIA AQUI FOI REMOVIDA.
//
// Ela aceitava os nomes em português — `regra-cumprida`, `trinca-completa`, `spec-completa`
// e outros dezoito — e os convertia na carga, para sempre. Isso resolve e não fecha: o
// arquivo nunca se conserta, e o mapa acumula carimbos nos DOIS formatos.
//
// E tinha um defeito de assimetria: a LEITURA normalizava (`mapx.mesmoGate`), a ESCRITA
// não. Um projeto que renomeasse o gate ganhava um SEGUNDO carimbo em vez de atualizar o
// primeiro. Medido no blue-eyes: 40 julgamentos gravados como `regra-cumprida` convivendo
// com 2 como `rule-fulfilled` — o mesmo gate contado duas vezes, e o `check` rejulgando o
// que já fora respondido.
//
// Com o contrato de formato, o alias deixou de ser necessário: a conversão virou o passo
// de migração `1→2` (ver `internal/migra/formato_2.go`), roda UMA VEZ, e o formato 2 só
// tem nome canônico. Quem está no formato 1 é recusado com a mensagem que manda migrar —
// não lido pela metade.
// gates default, que vive no `initx`. E `config` não pode importar `initx` — seria ciclo,
// já que o `initx` monta gates a partir de tipos daqui.
//
// A injeção resolve: o `initx` registra a lista no seu `init()`, e o teste a consulta. Sem
// isso o teste não teria contra o que confrontar, e um de-para apontando para o vazio
// passaria — o projeto carregaria, o gate viraria um nome que nenhum verificador conhece,
// e o `check` reportaria "gate sem nada para medir" sem dizer que a causa foi a conversão.

var defaultGateNames func() []string

// RegisterGateNames liga a lista de gates default ao pacote config.
func RegisterGateNames(f func() []string) { defaultGateNames = f }

// DefaultGateNamesForTest devolve os nomes registrados, ou nil se ninguém registrou.
func DefaultGateNamesForTest() []string {
	if defaultGateNames == nil {
		return nil
	}
	return defaultGateNames()
}
