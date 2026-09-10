package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// --- um derivado com VÁRIOS padrões ---
//
// `files: {code: "{{dir}}/{{name}}.ts"}` liga uma spec a UM arquivo, e isso cobre a
// unidade normal: uma tela, um handler, um componente.
//
// Não cobre a spec de CONFIGURAÇÃO. `TypeScriptConfig` descreve seis `tsconfig.json`
// espalhados pelos pacotes; `Workspace` descreve o `pnpm-workspace.yaml` e quatro
// `package.json`. Com um padrão só, a trinca dessas specs nunca fecha — e o gate reprova
// para sempre um trabalho que ESTÁ feito, ensinando a dispensar por hábito.
//
// Daí `Padroes`: o mesmo campo aceita string (a forma de sempre) ou lista.
//
//	files:
//	  code:
//	    - "tsconfig.base.json"
//	    - "packages/*/tsconfig.json"
//	    - "apps/*/tsconfig.json"
type Padroes []string

// UnmarshalYAML aceita as duas formas. A string continua válida porque a maioria das
// specs governa um arquivo só, e exigir lista de todas seria ruído em troca de nada.
func (p *Padroes) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		var s string
		if err := n.Decode(&s); err != nil {
			return err
		}
		*p = Padroes{s}
		return nil
	case yaml.SequenceNode:
		var l []string
		if err := n.Decode(&l); err != nil {
			return err
		}
		if len(l) == 0 {
			return fmt.Errorf("lista de padrões vazia: declare ao menos um, ou remova a camada")
		}
		*p = l
		return nil
	}
	return fmt.Errorf("padrão de derivado deve ser texto ou lista, veio %v", n.Kind)
}

// MarshalYAML escreve de volta na forma mais simples que couber: uma lista de um item
// vira string. Sem isto, todo `anchors init` reescreveria a config num formato mais
// verboso que o do usuário, e o diff mostraria mudança onde não houve.
func (p Padroes) MarshalYAML() (any, error) {
	if len(p) == 1 {
		return p[0], nil
	}
	return []string(p), nil
}
