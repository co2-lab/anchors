package common

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AliasDeFlag declara que antiga é o nome anterior de nova.
func AliasDeFlag(cmd *cobra.Command, nova, antiga string) {
	f := cmd.Flags().Lookup(nova)
	if f == nil {
		panic(fmt.Sprintf("aliasDeFlag: %q does not exist in %q", nova, cmd.Name()))
	}
	switch f.Value.Type() {
	case "bool":
		cmd.Flags().Bool(antiga, false, f.Usage)
	default:
		cmd.Flags().String(antiga, "", f.Usage)
	}
	_ = cmd.Flags().MarkDeprecated(antiga, fmt.Sprintf("use --%s", nova))
	_ = cmd.Flags().MarkHidden(antiga)
}

// ResolveAliases copia o valor da flag ANTIGA para a nova, quando só a antiga foi passada.
func ResolveAliases(cmd *cobra.Command, pares map[string]string) error {
	for nova, antiga := range pares {
		fa := cmd.Flags().Lookup(antiga)
		fn := cmd.Flags().Lookup(nova)
		if fa == nil || fn == nil || !fa.Changed || fn.Changed {
			continue
		}
		if err := cmd.Flags().Set(nova, fa.Value.String()); err != nil {
			return fmt.Errorf("--%s: %w", antiga, err)
		}
	}
	return nil
}
