package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// --- o PROGRESSO mora fora do plano ---
//
// Um plano é DECISÃO: o que vai ser feito, em que ordem, e por quê. Alterá-lo tem de
// significar que a decisão mudou — é sobre isso que o `plano-alterado-justificado` cobra
// uma revisão (`{CODIGO}-R0001`), e é a única defesa contra o projeto derivar em silêncio.
//
// Enquanto os checkboxes de fase viviam no plano, marcar `- [x]` era ALTERAR o plano. E
// aí três coisas quebravam de uma vez:
//
//  1. o gate não distinguia "mudei a direção" de "terminei uma fase", e cobrava revisão
//     das duas. Uma revisão que diz "concluí a fase 1" é ruído — e ruído em gate
//     bloqueante é o que faz alguém desligá-lo.
//
//  2. o carimbo de julgamento da aresta `plano → spec` guarda a rev das DUAS pontas.
//     Marcar a fase mudava a rev do plano, e o julgamento da SPEC caía — um julgamento
//     sobre spec e código, derrubado por um checkbox. O efeito é circular: concluir a
//     fase invalida a verificação da spec que a fase entregou.
//
//  3. medido no blue-eyes: `plans/0017-mutacao.md` tinha DOIS commits — o que o criou (83
//     linhas) e um que mudou 1 linha, `- [ ]` para `- [x]`. Cem por cento das alterações
//     pós-criação eram progresso.
//
// O arquivo de progresso fica FORA DO MAPA de propósito. Se entrasse como nó com camada,
// os gates voltariam a confrontá-lo e o problema renasceria com outro nome — inclusive o
// `plano-alterado-justificado` cobrando revisão de um arquivo cuja única função é mudar.
// Ele é estado, não decisão: ninguém precisa justificar por que o estado avançou.

// sufixoProgresso liga o arquivo de estado ao plano: mesmo nome, sufixo fixo.
//
// A definição CANÔNICA é a do `scan` — é ele que precisa manter o arquivo fora do mapa,
// e uma segunda constante aqui poderia divergir dela em silêncio. Esta é a mesma string,
// e o teste `TestProgresso_sufixoBateComOScan` confronta as duas.
//
// A colocação ao lado (em vez de uma pasta `progress/`) é o que faz os dois serem lidos
// juntos: quem abre o plano vê o companheiro na mesma listagem, e um plano cujo progresso
// não existe fica visível pela ausência.
const sufixoProgresso = "-progress.md"

// caminhoDeProgresso devolve o arquivo de progresso de um plano.
func caminhoDeProgresso(plano string) string {
	ext := filepath.Ext(plano)
	return strings.TrimSuffix(plano, ext) + sufixoProgresso
}

// fasesDoPlano lê os códigos de fase declarados nos cabeçalhos do plano.
//
// A fonte é o CABEÇALHO (`### PLTFR-F01 — ...`), a mesma que os gates `fase-existe` e
// `fase-ordenada` já usam. Ler daqui em vez de manter uma segunda lista é o que garante
// que o progresso fale das fases que existem: uma fase renomeada aparece, uma inventada
// não.
var faseNoCabecalhoRE = regexp.MustCompile(`(?m)^#{2,4}\s+([A-Z0-9]{4,5}-F\d{2})\b[^\S\n]*—?[^\S\n]*(.*)$`)

type faseDoPlano struct {
	Codigo string
	Titulo string
}

func fasesDoPlano(conteudo string) []faseDoPlano {
	var out []faseDoPlano
	for _, m := range faseNoCabecalhoRE.FindAllStringSubmatch(conteudo, -1) {
		out = append(out, faseDoPlano{Codigo: m[1], Titulo: strings.TrimSpace(m[2])})
	}
	return out
}

// escreveProgressoInicial cria o `-progress.md` de um plano, com uma linha por fase.
//
// Não sobrescreve: o arquivo guarda o estado do trabalho, e regravá-lo apagaria o que já
// foi registrado. Um plano que ganha fase nova tem a linha acrescentada à mão — o comando
// não reescreve estado que não é dele.
func escreveProgressoInicial(planoPath, conteudoPlano, codigo string) (string, error) {
	destino := caminhoDeProgresso(planoPath)
	if _, err := os.Stat(destino); err == nil {
		return "", fmt.Errorf("%s já existe — o progresso é estado, e não sobrescrevo", destino)
	}

	var b strings.Builder
	b.WriteString("<!-- anchors:progress — o ESTADO do plano `" +
		filepath.Base(planoPath) + "`.\n\n")
	b.WriteString("Este arquivo fica FORA DO MAPA de propósito: ele existe para MUDAR, e um\n")
	b.WriteString("arquivo que muda por natureza não pode ser confrontado pelos gates que\n")
	b.WriteString("cobram justificativa de mudança. O plano ao lado é a DECISÃO; alterá-lo tem\n")
	b.WriteString("de significar que a decisão mudou.\n\n")
	b.WriteString("Marque `[x]` aqui, nunca no plano.\n-->\n\n")
	b.WriteString("# Progresso — " + codigo + "\n\n")

	fases := fasesDoPlano(conteudoPlano)
	if len(fases) == 0 {
		b.WriteString("TODO: o plano ainda não declara fases. Quando declarar, acrescente uma\n")
		b.WriteString("seção por fase aqui, com um item por spec semeada.\n")
	}
	for _, f := range fases {
		b.WriteString("## " + f.Codigo)
		if f.Titulo != "" {
			b.WriteString(" — " + f.Titulo)
		}
		b.WriteString("\n\n- [ ] TODO: um item por spec que esta fase semeia\n\n")
	}

	if err := os.WriteFile(destino, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return destino, nil
}
