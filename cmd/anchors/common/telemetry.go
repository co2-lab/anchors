// @anchors
//   ref: TLSTT

package common

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/telemetry"
	"github.com/spf13/cobra"
)

// Emitter é o destino dos eventos deste processo.
var Emitter *telemetry.Emitter

// NoticeTelemetry mostra o aviso na primeira execução neste projeto, nesta máquina, e liga
// o emitter.
func NoticeTelemetry(cmd *cobra.Command) {
	root := ProjectRoot(cmd)
	// OUTSIDE A PROJECT: no notice, no emitter, no mark. The root used to fall back to the
	// working directory, and the notice's marker was written there — a `.anchors/` left in
	// whatever directory the command ran from. With no project there is no `telemetry:`
	// to honour nor a project's unversioned directory to keep the mark in.
	if root == "" {
		return
	}
	cfgPath := filepath.Join(root, config.DefaultFile)
	if _, err := os.Stat(cfgPath); err != nil {
		return
	}

	// From the project ROOT, not the cwd: run from a subdirectory or with `--root`, the
	// cwd has no anchors.yaml, and a project's `telemetry: off` was ignored.
	declared := declaredTelemetry(cfgPath)
	if telemetry.Disabled(declared) {
		return
	}

	telemetry.Notice(os.Stderr, root)

	Emitter = telemetry.NewEmitter(telemetry.Config{
		Enabled:  true,
		Endpoint: os.Getenv("ANCHORS_TELEMETRY_ENDPOINT"),
		Headers:  telemetryHeaders(),
		// NoCodes stays false: no project setting declares it yet (see the spec).
	}, Version)
}

// telemetryLineRE reads the top-level `telemetry:` value as text.
var telemetryLineRE = regexp.MustCompile(`(?m)^telemetry:[ \t]*["']?([^\s"'#]*)`)

// declaredTelemetry returns the project's `telemetry:` value. When the configuration does
// not load (an unknown key, a format still to migrate) the line is read as text: skipping
// the file dropped a declared `telemetry: off`, and the notice and the emitter came up for a
// project that had opted out.
func declaredTelemetry(cfgPath string) string {
	if cfg, err := config.Load(cfgPath); err == nil && cfg != nil {
		return cfg.Telemetry
	}
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return ""
	}
	if m := telemetryLineRE.FindSubmatch(b); m != nil {
		return strings.TrimSpace(string(m[1]))
	}
	return ""
}

func telemetryHeaders() map[string]string {
	if k := os.Getenv("ANCHORS_TELEMETRY_KEY"); k != "" {
		return map[string]string{"x-honeycomb-team": k}
	}
	return nil
}

// ProjectRoot devolve a raiz, ou "" fora de um projeto.
//
// An explicit `--root` is taken as given. Without it, the root is the nearest directory
// above the working directory that holds `anchors.yaml`, and "" when none does — it used to
// return the working directory itself, contradicting this comment.
func ProjectRoot(cmd *cobra.Command) string {
	if f := cmd.Flags().Lookup("root"); f != nil && f.Value.String() != "" && f.Value.String() != "." {
		if abs, err := config.AbsRoot(f.Value.String()); err == nil {
			return abs
		}
	}
	abs, err := config.AbsRoot(".")
	if err != nil {
		return ""
	}
	if _, err := os.Stat(filepath.Join(abs, config.DefaultFile)); err != nil {
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
