package config

import (
	"fmt"
	"strconv"
	"strings"
)

// A COMPARAÇÃO de versão, e por que ela não é um `<` ingênuo.
//
// O `staleBinaryWarning` do `check` compara por IGUALDADE, de propósito: ele confronta
// quem gravou o mapa com quem está rodando, e ali qualquer diferença importa. Aqui a
// pergunta é outra — "este binário atende ao mínimo?" —, e ela é ORDINAL.
//
// O caso que obriga a tratar `dev` à parte: um build local se identifica como "dev", e
// ordenar "dev" contra "0.1.84" não tem resposta certa. Ele pode ser mais novo que
// qualquer release (a árvore de quem desenvolve o Anchors) ou mais velho que todas.
// Recusá-lo tornaria impossível desenvolver o próprio Anchors contra um projeto que
// declara mínimo; aceitá-lo em silêncio esconde o caso em que ele é velho.
//
// A saída é AVISAR e não barrar: `dev` responde "não sei dizer".

// VersionOrder compara duas versões `MAJOR.MINOR.PATCH`.
//
// Devolve -1, 0 ou 1, e um erro quando alguma das duas não é ordenável — o que inclui
// "dev", string vazia, e qualquer coisa que não seja numérica.
func VersionOrder(a, b string) (int, error) {
	pa, err := parseVersion(a)
	if err != nil {
		return 0, err
	}
	pb, err := parseVersion(b)
	if err != nil {
		return 0, err
	}
	for i := range pa {
		switch {
		case pa[i] < pb[i]:
			return -1, nil
		case pa[i] > pb[i]:
			return 1, nil
		}
	}
	return 0, nil
}

// parseVersion aceita `0.1.84` e `v0.1.84`, e nada mais.
//
// Não aceita sufixo (`0.1.84-rc1`): uma pré-release comparada como se fosse a final
// responderia "atende" a quem tem menos do que o mínimo pede, e o silêncio seria pior que
// a recusa de comparar.
func parseVersion(v string) ([3]int, error) {
	var out [3]int
	v = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(v), "v"))
	if v == "" {
		return out, fmt.Errorf("versão vazia")
	}
	partes := strings.Split(v, ".")
	if len(partes) != 3 {
		return out, fmt.Errorf("versão %q não é MAJOR.MINOR.PATCH", v)
	}
	for i, p := range partes {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return out, fmt.Errorf("versão %q não é ordenável", v)
		}
		out[i] = n
	}
	return out, nil
}

// AtendeMinVersion responde se `rodando` satisfaz o mínimo declarado.
//
// Os três casos, e cada um tem uma razão:
//
//   - mínimo não declarado → atende. A maioria dos projetos não declara, e exigir a
//     declaração para poder trabalhar inverteria o padrão.
//   - alguma das duas não é ordenável (`dev`, formato estranho) → NÃO atende, com erro.
//     Quem chama decide se avisa ou barra; o que não se faz é fingir que sabe.
//   - ambas ordenáveis → a comparação.
func AtendeMinVersion(minimo, rodando string) (bool, error) {
	if strings.TrimSpace(minimo) == "" {
		return true, nil
	}
	ord, err := VersionOrder(rodando, minimo)
	if err != nil {
		return false, err
	}
	return ord >= 0, nil
}

// validarMinVersion recusa a CARGA quando o campo não é `MAJOR.MINOR.PATCH`.
//
// Aqui é o lugar, e não na comparação: `min_version: 0.1` ou `min_version: latest` só
// falharia ao comparar, e o efeito de falhar lá é o AVISO SUMIR — o campo existe para
// forçar uma atualização, e um valor mal escrito produz exatamente o silêncio que ele
// deveria quebrar. Quem escreve descobre no primeiro comando, não quando alguém não for
// avisado.
//
// `dev` é recusado NESTE campo, e isso não contradiz o tratamento do binário `dev` em
// `AtendeMinVersion`: um binário de desenvolvimento é um fato que se encontra, e o mínimo
// é uma decisão que alguém escreve. Declarar "no mínimo dev" não afirma nada.
func (c *Config) validarMinVersion() error {
	v := strings.TrimSpace(c.MinVersion)
	if v == "" {
		return nil
	}
	if _, err := parseVersion(v); err != nil {
		return fmt.Errorf("min_version: %q não é MAJOR.MINOR.PATCH (ex.: 0.1.84) — "+
			"o campo declara a versão mínima do binário que este projeto aceita, e um "+
			"valor que não se compara silencia o aviso que ele existe para dar", c.MinVersion)
	}
	return nil
}
