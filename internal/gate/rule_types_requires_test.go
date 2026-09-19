package gate

import (
	"github.com/co2-lab/anchors/internal/config"
)

func cfgComRequires() *config.Config {
	return &config.Config{RuleTypes: []config.RuleType{
		{Letter: "B", Term: "Behavior",
			Sections:            []string{"Eventos / Callbacks", "Comportamentos"},
			SectionsRequireCode: []string{"Eventos / Callbacks"}},
		{Letter: "S", Term: "State",
			Sections: []string{"Variantes", "Estados Visuais"}},
	}}
}
