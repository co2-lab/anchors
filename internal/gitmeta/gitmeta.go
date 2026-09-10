// Package gitmeta lê metadados do git de um arquivo — hoje, a data do último commit
// que o tocou. É a fonte de verdade para o carimbo de alteração (updated_at): mais
// confiável que uma data mantida à mão, que sempre desatualiza.
package gitmeta

import (
	"os/exec"
	"strings"
	"time"
)

// Today devolve a data do sistema (AAAA-MM-DD). Usada como desempate quando o arquivo
// tem alterações não-commitadas: a "última alteração" é hoje, não o último commit. Um
// gate de data legitimamente precisa saber que dia é — é da natureza dele.
func Today() string { return time.Now().Format("2006-01-02") }

// LastCommitDate devolve a data (AAAA-MM-DD) do ÚLTIMO commit que modificou `rel`
// (caminho relativo à raiz). Só a data — nunca a hora: um arquivo pode ser editado
// num horário e commitado em outro do MESMO dia, e isso é irrelevante; só a mudança
// de DIA importa. Devolve ok=false se não é repo git, o arquivo nunca foi commitado,
// ou o git não está disponível.
func LastCommitDate(root, rel string) (date string, ok bool) {
	cmd := exec.Command("git", "-C", root, "log", "-1", "--format=%ad", "--date=short", "--", rel)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	d := strings.TrimSpace(string(out))
	if d == "" {
		return "", false // arquivo sem commit (novo/untracked)
	}
	return d, true
}

// AllCommitDates devolve, num só `git log`, a data (AAAA-MM-DD) do último commit que
// tocou CADA arquivo do repo — o carimbo de alteração de todos, sem 1 chamada git por
// arquivo. Percorre o log com --name-only e associa a data de cada commit ao primeiro
// aparecimento de cada arquivo (o mais recente, pois o log vem do mais novo ao antigo).
// Mapa vazio se não for repo git.
func AllCommitDates(root string) map[string]string {
	cmd := exec.Command("git", "-C", root, "log", "--format=D:%ad", "--date=short", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return map[string]string{}
	}
	dates := map[string]string{}
	var cur string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "D:") {
			cur = strings.TrimPrefix(line, "D:")
			continue
		}
		f := strings.TrimSpace(line)
		if f == "" || cur == "" {
			continue
		}
		if _, seen := dates[f]; !seen { // 1º aparecimento = commit mais recente
			dates[f] = cur
		}
	}
	return dates
}

// HasUncommittedChanges diz se `rel` tem alterações não-commitadas no working tree.
// Útil para não reclamar de um arquivo em edição ativa (a "última alteração" é hoje,
// não o último commit).
//
// ATENÇÃO ao `false`: ele significa "sem alteração pendente" E TAMBÉM "não deu para
// perguntar" (sem git, sem repo). Quem toma DECISÃO a partir disso deve usar
// `UncommittedChanges`, que separa os dois — confundi-los faz o Anchors afirmar sobre
// um estado que ele não conseguiu ler.
func HasUncommittedChanges(root, rel string) bool {
	mudou, _ := UncommittedChanges(root, rel)
	return mudou
}

// UncommittedChanges é o `HasUncommittedChanges` que não esconde a ignorância: o
// segundo retorno diz se a pergunta pôde ser FEITA. Sem ele, "o arquivo está limpo" e
// "não há repositório para consultar" chegam ao chamador como o mesmo `false`, e a
// diferença entre os dois é justamente o que decide se há um veredito a dar.
func UncommittedChanges(root, rel string) (mudou, sabido bool) {
	cmd := exec.Command("git", "-C", root, "status", "--porcelain", "--", rel)
	out, err := cmd.Output()
	if err != nil {
		return false, false
	}
	return strings.TrimSpace(string(out)) != "", true
}

// Head devolve o hash curto do HEAD e o assunto do commit. Serve para carimbar um
// relatório com a versão do código que ele descreve: sem isso, uma leitura salva
// não diz sobre QUE código ela fala, e envelhece sem avisar.
//
// ok=false quando não é repo git, o git não está disponível, ou o repo ainda não
// tem commit nenhum.
func Head(root string) (short, subject string, ok bool) {
	out, err := exec.Command("git", "-C", root, "log", "-1", "--format=%h%x00%s").Output()
	if err != nil {
		return "", "", false
	}
	partes := strings.SplitN(strings.TrimRight(string(out), "\n"), "\x00", 2)
	if len(partes) != 2 || partes[0] == "" {
		return "", "", false
	}
	return partes[0], partes[1], true
}

// DirtyCount conta os arquivos com alteração não-commitada na árvore inteira.
// Um relatório rodado sobre árvore suja descreve um estado que não está em commit
// nenhum — quem relê precisa saber que não vai reencontrá-lo pelo hash.
//
// Devolve -1 quando não deu para contar (sem git, sem repo). NÃO é 0: zero afirma
// "árvore limpa", e carimbar um relatório com essa afirmação sem tê-la verificado é
// exatamente o silêncio que tranquiliza — o pior tipo.
func DirtyCount(root string) int {
	out, err := exec.Command("git", "-C", root, "status", "--porcelain").Output()
	if err != nil {
		return -1
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}
