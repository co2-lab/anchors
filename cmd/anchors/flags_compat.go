package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// --- as FLAGS antigas em português, que ainda funcionam ---
//
// As flags são superfície pública tanto quanto os nomes de gate: elas estão em scripts,
// em pipelines de CI, em anotações de quem usa. Renomear sem alias quebraria tudo isso
// numa atualização de binário.
//
// O Cobra resolve isso de fábrica com `MarkDeprecated`: a flag antiga continua sendo
// aceita, e quem a usa recebe um aviso dizendo qual é a nova. Ela some numa versão maior,
// não numa correção.
//
// A régua do que traduzir é a mesma do vocabulário de gates: traduz-se o que se LÊ (a
// descrição da flag, que passa pelo i18n), não o que se ESCREVE na linha de comando.

// aliasDeFlag declara que `antiga` é o nome anterior de `nova`.
//
// Registra a antiga como flag de verdade — senão o Cobra a recusaria como desconhecida —,
// marca-a obsoleta, e a esconde do `--help`: quem está lendo a ajuda hoje deve ver só o
// nome novo, e quem tem a antiga num script recebe o aviso ao rodar.
func aliasDeFlag(cmd *cobra.Command, nova, antiga string) {
	f := cmd.Flags().Lookup(nova)
	if f == nil {
		// Programação defeituosa, não erro de usuário: um alias para flag inexistente
		// significa que alguém renomeou a flag e esqueceu de atualizar o alias.
		panic(fmt.Sprintf("aliasDeFlag: %q não existe em %q", nova, cmd.Name()))
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

// resolveAliases copia o valor da flag ANTIGA para a nova, quando só a antiga foi passada.
//
// Roda no `PreRunE` do comando, antes de o `RunE` ler as variáveis. Sem isto, quem usasse
// a flag antiga veria o aviso de obsolescência e o comando rodaria com o valor vazio — o
// pior desfecho possível, porque parece que funcionou.
func resolveAliases(cmd *cobra.Command, pares map[string]string) error {
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
