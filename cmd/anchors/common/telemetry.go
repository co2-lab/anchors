package common

import (
	"os"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/telemetry"
	"github.com/spf13/cobra"
)

// Emitter é o destino dos eventos deste processo.
var Emitter *telemetry.Emitter

// NoticeTelemetry mostra o aviso na primeira execução desta máquina e liga o emitter.
func NoticeTelemetry(cmd *cobra.Command) {
	root := ProjectRoot(cmd)

	var declared string
	if root != "" {
		if cfg, err := config.Load(config.DefaultFile); err == nil && cfg != nil {
			declared = cfg.Telemetry
		}
	}
	if telemetry.Disabled(declared) {
		return
	}

	telemetry.Notice(os.Stderr, root)

	Emitter = telemetry.NewEmitter(telemetry.Config{
		Enabled:  true,
		Endpoint: os.Getenv("ANCHORS_TELEMETRY_ENDPOINT"),
		Headers:  telemetryHeaders(),
	}, Version)
}

func telemetryHeaders() map[string]string {
	if k := os.Getenv("ANCHORS_TELEMETRY_KEY"); k != "" {
		return map[string]string{"x-honeycomb-team": k}
	}
	return nil
}

// ProjectRoot devolve a raiz, ou "" fora de um projeto.
func ProjectRoot(cmd *cobra.Command) string {
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

// FlushTelemetry espera os envios em voo antes de o processo sair.
func FlushTelemetry() {
	if Emitter != nil {
		Emitter.Flush()
	}
}
