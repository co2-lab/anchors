package common

import (
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
