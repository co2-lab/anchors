package main

import (
	"os"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/telemetry"
	"github.com/spf13/cobra"
)

// emitter é o destino dos eventos deste processo.
//
// Global e não passado por parâmetro: a alternativa seria acrescentar um argumento a cada
// um dos 43 comandos, e o primeiro que fosse esquecido perderia os eventos em silêncio.
// Nil quando a telemetry está desligada — e `Emit` num nil não faz nada, por desenho.
var emitter *telemetry.Emitter

// noticeTelemetry mostra o aviso na primeira execução desta máquina e liga o emitter.
//
// AVISA ANTES DE LIGAR, e a ordem é a razão de tudo isto estar no `PersistentPreRunE`: o
// aviso precisa sair antes do primeiro evento. Invertido, a coleta começaria antes de a
// pessoa saber que existe.
func noticeTelemetry(cmd *cobra.Command) {
	root := projectRoot(cmd)

	// O ARQUIVO só é lido para saber se o projeto declarou `telemetry: off`. Se ele não
	// carregar (projeto sem `anchors.yaml`, comando rodando fora de projeto), o padrão do
	// produto vale — e o aviso sai igual, porque quem roda `anchors --help` numa pasta
	// qualquer também merece saber.
	var declared string
	if root != "" {
		if cfg, err := config.Load(config.DefaultFile); err == nil && cfg != nil {
			declared = cfg.Telemetry
		}
	}
	if telemetry.Disabled(declared) {
		return
	}

	// O aviso vai para o STDERR: o stdout de muitos comandos é consumido por script
	// (`anchors code list --json`), e um aviso no meio do JSON quebraria quem o lê.
	telemetry.Notice(os.Stderr, root)

	emitter = telemetry.NewEmitter(telemetry.Config{
		Enabled:  true,
		Endpoint: os.Getenv("ANCHORS_TELEMETRY_ENDPOINT"),
		Headers:  telemetryHeaders(),
	}, version)
}

// telemetryHeaders lê a credencial do AMBIENTE, nunca do arquivo.
//
// O `anchors.yaml` é versionado: uma chave ali iria para o repositório de quem usa, e
// `secret-nao-vazado` — gate do próprio produto — a acusaria. O ambiente é onde credencial
// mora.
func telemetryHeaders() map[string]string {
	if k := os.Getenv("ANCHORS_TELEMETRY_KEY"); k != "" {
		return map[string]string{"x-honeycomb-team": k}
	}
	return nil
}

// projectRoot devolve a raiz, ou "" fora de um projeto.
//
// Sem erro: o chamador não tem o que fazer com ele. Fora de projeto o aviso ainda sai (para
// o diretório atual), e é o comportamento certo — quem roda `anchors guide` numa pasta
// qualquer também precisa saber da coleta.
func projectRoot(cmd *cobra.Command) string {
	if f := cmd.Flags().Lookup("root"); f != nil && f.Value.String() != "" {
		if abs, err := config.AbsRoot(f.Value.String()); err == nil {
			return abs
		}
	}
	abs, err := config.AbsRoot(".")
	if err != nil {
		return ""
	}
	return abs
}

// flushTelemetry espera os envios em voo antes de o processo sair.
//
// Nil-safe: sem telemetria o emitter é nil, e `Flush` num nil não faz nada.
func flushTelemetry() { emitter.Flush() }
