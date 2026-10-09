// @anchors
//   code: ERCMR
//   ref: CMCLC

package common

import (
	"fmt"

	"github.com/co2-lab/anchors/internal/i18n"
)

// ExitNotGoverned é o código de saída para "este caminho não é regido pela Estrutura".
const ExitNotGoverned = 3

// ErrNotGoverned indica que o caminho não é regido pela Estrutura.
type ErrNotGoverned struct {
	Path string
}

func (e ErrNotGoverned) Error() string {
	return i18n.T("not_governed", e.Path)
}

// ExitCode is the exit code a command ends with, returned instead of leaving the process on
// the spot: `main` closes what the command opened — the record of its run, the telemetry —
// and exits with the code, printing nothing more, for the command already said what it had
// to. A command that called os.Exit itself left its run recorded as died (reported from MIF:
// a blocked `anchors check`).
type ExitCode struct{ Code int }

func (e ExitCode) Error() string { return fmt.Sprintf("exit status %d", e.Code) }
