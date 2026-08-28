package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// layer-boundary: uma camada não alcança o que não é dela.
//
// O Anchors já declarava camadas — mas não confrontava se elas se RESPEITAM. Declarar
// `screens/`, `hooks/`, `repositories/` e não verificar que a tela não fala direto com o
// repositório é desenhar a arquitetura e não defendê-la: o `layers:` vira documentação.
//
// A lacuna aparecia porque cada projeto resolvia sozinho. Um projeto real mantinha um
// script de 366 linhas com 15 regras arquiteturais escritas à mão em shell — todas da
// mesma forma: "arquivos que casam ESTE padrão não podem conter AQUELE outro". Quinze
// instâncias de um único mecanismo, reimplementado porque o framework não o oferecia.
//
// Declaração (anchors.yaml), agnóstica de linguagem — quem escreve o padrão é o projeto:
//
//	boundaries:
//	  - layer: screens                  # a quem a regra se aplica (nome da camada)
//	    forbid: "from '@/repositories"   # o que não pode aparecer
//	    because: "tela fala com hook, não com dado"
//	    severity: error                  # error (default) | warn — a maturação do gate
//
// O opt-out honesto (CONCEPT §5.1): `@allow-boundary: <razão>` na linha dispensa AQUELA
// linha, com a razão escrita. Dívida reconhecida fica visível e datada no código, em vez
// de virar uma exceção numa lista distante que ninguém revisita.
func checkLayerBoundary(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, "a fronteira se lê no código — é ele que importa ou não importa"
	}
	if cfg == nil || len(cfg.Boundaries) == 0 {
		// Pendente, não Pass: o Anchors não sabe a arquitetura do projeto, e fingir que
		// verificou seria pior do que dizer o que falta (QUALITY §7, o terceiro estado).
		return Pending, "o projeto não declarou fronteiras de camada. Declare `boundaries:` " +
			"no anchors.yaml (`layer` + `forbid` + `because`) — sem isso, `layers:` é " +
			"documentação: descreve a arquitetura sem defendê-la"
	}

	// A camada do nó vem da mesma classificação que o scan usa — a regra se declara pelo
	// nome da camada (`layers:`), que é o vocabulário que o projeto já tem. Uma regra sem
	// `layer` vale para TODO código: é assim que se declara uma proibição global (ex.:
	// "sem relógio cru em lugar nenhum").
	minhaCamada, _ := scan.Classify(n.ID, cfg)
	var aplicaveis []config.Boundary
	for _, b := range cfg.Boundaries {
		if b.Layer == "" || b.Layer == minhaCamada {
			aplicaveis = append(aplicaveis, b)
		}
	}
	if len(aplicaveis) == 0 {
		return Skip, fmt.Sprintf("nenhuma fronteira declarada para a camada %q", minhaCamada)
	}

	linhas := strings.Split(content, "\n")
	var erros, avisos []string
	for _, b := range aplicaveis {
		// `(?s)` para o `.` cruzar a quebra de linha: o padrão é casado contra o arquivo
		// INTEIRO, não linha a linha.
		//
		// Era linha a linha, e isso furava toda regra cujo alvo se espalha por várias
		// linhas — um import formatado pelo Prettier, por exemplo. Medido no app de referência: das 7
		// telas que importam `Modal` do react-native, o gate acusava UMA — a única com o
		// import numa linha só. As outras 6 quebram em várias porque passam de 100
		// colunas, e escapavam. A tela era acusada por FORMATAÇÃO, não por ser diferente
		// das outras, e o projeto não tinha como saber disso: o gate ficava verde.
		//
		// Um projeto que queira ancorar o padrão numa linha só continua podendo — é o que
		// `^` e `$` fazem, e com `(?s)` eles seguem valendo por linha (Go não liga
		// `(?m)` junto).
		re, err := regexp.Compile("(?s)" + b.Forbid)
		if err != nil {
			// Padrão inválido é erro de CONFIG, e precisa aparecer — um regex quebrado
			// que fosse ignorado em silêncio desligaria a regra sem ninguém saber.
			erros = append(erros, fmt.Sprintf("a fronteira %q tem `forbid` inválido (%v)", b.Layer, err))
			continue
		}
		for _, loc := range re.FindAllStringIndex(content, -1) {
			ini, fim := linhaDoOffset(content, loc[0]), linhaDoOffset(content, loc[1]-1)
			// A dispensa vale se estiver em QUALQUER linha do trecho casado (num import
			// multilinha ela fica na linha do `from`, não na do `import`) ou na linha
			// imediatamente acima.
			if trechoDispensa(linhas, ini, fim) || linhaAnteriorDispensa(linhas, ini) {
				continue
			}
			achado := fmt.Sprintf("linha %d: %s", ini+1, descreveFronteira(b))
			if b.Severity == "warn" {
				avisos = append(avisos, achado)
			} else {
				erros = append(erros, achado)
			}
		}
	}

	if len(erros) == 0 && len(avisos) == 0 {
		return Pass, ""
	}
	// Avisos sozinhos não reprovam — é a maturação por regra (QUALITY §7) dentro do gate:
	// o projeto marca `severity: warn` no que ainda está migrando, sem desligar o gate
	// inteiro nem perder o registro.
	if len(erros) == 0 {
		return Pending, "fronteira em migração (`severity: warn`): " + juntaAte(avisos, 5)
	}
	msg := "fronteira de camada violada: " + juntaAte(erros, 5)
	if len(avisos) > 0 {
		msg += fmt.Sprintf(" — e mais %d aviso(s) de regra em migração", len(avisos))
	}
	return Fail, msg + ". Se a violação é dívida reconhecida, marque a linha com " +
		"`@allow-boundary: <razão>`: a exceção fica visível no código, não numa lista distante"
}

func descreveFronteira(b config.Boundary) string {
	quem := "esta camada"
	if b.Layer != "" {
		quem = "a camada `" + b.Layer + "`"
	}
	s := fmt.Sprintf("%s não pode conter `%s`", quem, b.Forbid)
	if b.Because != "" {
		s += " — " + b.Because
	}
	return s
}

func juntaAte(xs []string, n int) string {
	sort.Strings(xs)
	if len(xs) > n {
		return strings.Join(xs[:n], "; ") + fmt.Sprintf(" (e mais %d)", len(xs)-n)
	}
	return strings.Join(xs, "; ")
}

// allowBoundaryRE — a dispensa exige RAZÃO escrita. `[^\S\n]` = espaço/tab e não quebra
// de linha: sem isso a razão seria "achada" na linha seguinte e um marcador nu passaria.
var allowBoundaryRE = regexp.MustCompile(`@allow-boundary[^\S\n]*:[^\S\n]*\S+`)

// linhaAnteriorDispensa aceita a marcação no comentário ACIMA da linha, além de na
// própria linha: em várias linguagens o import não tem onde receber um comentário de
// fim de linha legível, e obrigar a marcação inline empurraria o autor a não marcar.
func linhaAnteriorDispensa(linhas []string, i int) bool {
	return i > 0 && allowBoundaryRE.MatchString(linhas[i-1])
}

// linhaDoOffset converte um offset em bytes no conteúdo para o índice da linha (0-based).
func linhaDoOffset(content string, off int) int {
	if off > len(content) {
		off = len(content)
	}
	return strings.Count(content[:off], "\n")
}

// trechoDispensa aceita `@allow-boundary:` em qualquer linha do trecho casado. Num import
// que o Prettier quebrou, a marcação natural fica na linha do `from` — cobrar que ela
// esteja na primeira linha do casamento seria exigir que o autor soubesse onde o regex
// começou a casar.
func trechoDispensa(linhas []string, ini, fim int) bool {
	for i := ini; i <= fim && i < len(linhas); i++ {
		if allowBoundaryRE.MatchString(linhas[i]) {
			return true
		}
	}
	return false
}
