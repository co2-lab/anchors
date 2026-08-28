// Package checklog espelha a saída do `anchors check` num arquivo, para que ela
// possa ser RELIDA sem re-executar.
//
// O `check --all` custa minutos. Sem o espelho, toda pergunta sobre o resultado
// ("o que reprovou?", "quais gates passaram?") obriga a rodar tudo de novo só
// para reler o que já havia saído na tela — e a tela, num terminal com scroll
// limitado ou numa sessão de agente, some.
//
// O arquivo é separado por ESCOPO. O `--all` é a foto completa e cara; o
// `--changed` roda a cada commit sobre um ou dois arquivos. Num arquivo só, o
// pre-commit apagaria a foto completa no commit seguinte, que é justamente
// quando ela seria mais útil.
//
// O cabeçalho carrega o instante, o HEAD e o estado da árvore. Sem isso o
// arquivo mente por omissão: uma leitura de ontem parece a de agora, e a
// diferença entre "o gate passou" e "o gate passava antes daquela mudança" é
// exatamente o que se quer saber.
package checklog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Dir é a pasta de estado, a mesma do daemon (efêmero, não versionável).
const Dir = ".anchors"

// NomeDoEscopo devolve o arquivo onde a saída deste escopo é espelhada.
func NomeDoEscopo(all bool) string {
	if all {
		return "check-all.txt"
	}
	return "check-changed.txt"
}

// Espelho duplica tudo que for escrito para o stdout no arquivo.
type Espelho struct {
	arq      *os.File
	original *os.File
	pipeR    *os.File
	pipeW    *os.File
	pronto   chan struct{}
}

// Abrir começa a espelhar o stdout no arquivo do escopo. O `cabecalho` é escrito
// antes de qualquer saída do comando.
//
// Falha ao abrir o arquivo NÃO é erro do check: o espelho é conveniência, e
// impedir a varredura porque o disco está cheio seria trocar o essencial pelo
// acessório. Devolve nil e o check segue escrevendo só na tela.
func Abrir(root string, all bool, cabecalho string) *Espelho {
	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil
	}
	arq, err := os.Create(filepath.Join(dir, NomeDoEscopo(all)))
	if err != nil {
		return nil
	}
	if _, err := io.WriteString(arq, cabecalho); err != nil {
		arq.Close()
		return nil
	}

	r, w, err := os.Pipe()
	if err != nil {
		arq.Close()
		return nil
	}
	e := &Espelho{arq: arq, original: os.Stdout, pipeR: r, pipeW: w, pronto: make(chan struct{})}
	os.Stdout = w

	// A cópia roda em goroutine: o pipe tem buffer finito, e sem alguém drenando
	// do outro lado uma saída longa (o `--all` passa de 60KB) travaria o próprio
	// comando que a produz.
	go func() {
		defer close(e.pronto)
		io.Copy(io.MultiWriter(e.original, arq), r)
	}()
	return e
}

// Fechar devolve o stdout ao lugar e espera a cópia terminar. Seguro com nil,
// para que o chamador possa usar `defer` sem checar.
func (e *Espelho) Fechar() {
	if e == nil {
		return
	}
	os.Stdout = e.original
	e.pipeW.Close()
	<-e.pronto
	e.pipeR.Close()
	e.arq.Close()
}

// Caminho é o arquivo escrito, para o comando poder dizer onde ficou.
func (e *Espelho) Caminho() string {
	if e == nil {
		return ""
	}
	return e.arq.Name()
}

// Cabecalho monta as linhas de contexto do topo do arquivo: quando rodou, sobre
// qual commit, com quantas mudanças pendentes e com que argumentos.
//
// O estado da árvore importa tanto quanto o commit: um check rodado com 40
// arquivos modificados descreve algo que não está em lugar nenhum do histórico,
// e quem relê precisa saber disso.
func Cabecalho(comando, head, assunto string, sujos int, quando time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n", comando)
	fmt.Fprintf(&b, "# quando: %s\n", quando.Format("2006-01-02 15:04:05 -0700"))
	if head != "" {
		fmt.Fprintf(&b, "# HEAD:   %s %s\n", head, assunto)
	}
	switch {
	case sujos < 0:
		// `sujos` negativo é "não deu para contar" (sem git/sem repo), NÃO "limpa".
		// Um relatório que afirma limpeza sem ter conseguido olhar carimba uma foto que
		// nunca existiu — e quem relê não tem como saber disso.
		b.WriteString("# árvore: desconhecida (sem repositório git — não deu para conferir)\n")
	case sujos == 1:
		b.WriteString("# árvore: 1 arquivo modificado (não commitado)\n")
	case sujos > 1:
		fmt.Fprintf(&b, "# árvore: %d arquivos modificados (não commitados)\n", sujos)
	default:
		b.WriteString("# árvore: limpa\n")
	}
	b.WriteString("#\n")
	b.WriteString("# Espelho da saída do check — releia daqui em vez de re-executar.\n")
	b.WriteString("# Se o HEAD ou a árvore mudaram desde `quando`, esta foto envelheceu.\n\n")
	return b.String()
}
