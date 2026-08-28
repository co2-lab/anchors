package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
)

// confrontarEntrega confronta o que o autor DECLARA contra o que está no disco, na hora
// de registrar — antes de o registro virar material para o revisor.
//
// As duas verificações nasceram de divergências reais, medidas na primeira rodada em que
// o fluxo `deliver → review` funcionou:
//
//  1. O registro afirmava "atualiza as contagens de DTAX (50→51 modelos…)" e o arquivo
//     NÃO fora tocado. Afirmação factual, falsificável, falsa — e trivialmente checável.
//  2. O registro dizia "as funções são PAGINADAS (loop em LastEvaluatedKey)" e citava o
//     gate `pagination-honored` pelo nome. O gate estava VERMELHO naquele arquivo. O
//     autor declarou a virtude e entregou o defeito, contra o gate que ele mesmo nomeou.
//
// Nenhuma das duas é falha de atenção que se resolva pedindo mais cuidado: a primeira é
// memória, a segunda é o custo de reler a saída de um gate informativo. O registro é o
// lugar certo para confrontá-las, porque é onde o autor AFIRMA.
//
// O confronto não BLOQUEIA o registro: ele imprime. Bloquear empurraria o autor a
// declarar menos — e o valor do registro está em ele declarar mais.
func confrontarEntrega(root string, files []string, unit string) {
	avisos, confrontou := arquivosNaoTocados(root, files)
	if !confrontou {
		// Sem git este confronto NÃO ACONTECEU. Calar aqui era o pior silêncio do
		// comando: a ausência de aviso é lida como "os arquivos declarados conferem",
		// que é uma afirmação que ninguém verificou.
		fmt.Println("\n⚠ não deu para confrontar os arquivos declarados contra o que mudou:")
		fmt.Println("  " + gitmeta.Explica(gitmeta.Verifica(root), "ler o estado do working tree"))
		fmt.Println("  a entrega segue, mas ninguém verificou se você fez o que declarou.")
	}
	if len(avisos) > 0 {
		fmt.Println("\n⚠ arquivos declarados que NÃO aparecem no diff:")
		for _, a := range avisos {
			fmt.Printf("    %s\n", a)
		}
		fmt.Println("  Ou você não fez o que declarou, ou já commitou — em qualquer caso, o")
		fmt.Println("  registro está afirmando algo que o disco não confirma.")
	}
	if aviso := mutacaoNaoMedida(root, unit); aviso != "" {
		fmt.Println("\n⚠ " + aviso)
	}
	if linhas := gatesVermelhos(root, files, unit); len(linhas) > 0 {
		fmt.Println("\n⚠ gates INFORMATIVOS reprovando nos arquivos desta entrega:")
		for _, l := range linhas {
			fmt.Printf("    %s\n", l)
		}
		fmt.Println("  Informativo não barra a promoção — e é exatamente por isso que passa")
		fmt.Println("  despercebido. Se você declarou o contrário do que o gate diz, confira")
		fmt.Println("  antes que o revisor confira por você.")
	}
}

// arquivosNaoTocados devolve os arquivos declarados que o git não vê como modificados
// nem como novos. Usa o git porque é a única fonte que sabe o que MUDOU — não basta o
// arquivo existir.
// O segundo retorno diz se o confronto ACONTECEU. Sem ele, "nenhum arquivo suspeito"
// e "não tive como olhar" seriam o mesmo `nil` — e quem lê a saída concluiria que está
// tudo certo por não ver aviso nenhum.
func arquivosNaoTocados(root string, files []string) (avisos []string, confrontou bool) {
	out, err := exec.Command("git", "-C", root, "status", "--porcelain").Output()
	if err != nil {
		return nil, false // sem git: não há como confrontar, e inventar seria pior
	}
	// `git status --porcelain` COLAPSA um diretório inteiramente novo numa linha só
	// (`?? pasta/`), sem listar o que há dentro. Comparar caminho exato cegava o confronto
	// justamente no caso mais comum — toda entrega que cria um diretório novo (um handler
	// novo, uma feature nova). Medido: o worker declarou 3 arquivos, 2 foram acusados de
	// não existir, e os 3 estavam no disco.
	//
	// Então uma linha terminada em `/` marca um PREFIXO tocado, não um caminho.
	tocado := map[string]bool{}
	var prefixos []string
	for _, l := range strings.Split(string(out), "\n") {
		if len(l) < 4 {
			continue
		}
		p := strings.TrimSpace(l[3:])
		if strings.HasSuffix(p, "/") {
			prefixos = append(prefixos, p)
			continue
		}
		tocado[p] = true
	}
	var faltando []string
	for _, f := range files {
		rel := relTo(root, f)
		if tocado[rel] || sobPrefixo(rel, prefixos) {
			continue
		}
		faltando = append(faltando, rel)
	}
	return faltando, true
}

// gatesVermelhos roda os gates do projeto sobre os arquivos entregues e devolve os que
// reprovam — inclusive os INFORMATIVOS, que são justamente os que somem do relatório.
func gatesVermelhos(root string, files []string, unit string) []string {
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil || len(cfg.Gates) == 0 {
		return nil
	}
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		return nil
	}
	alvo := map[string]bool{}
	for _, f := range append(append([]string{}, files...), unit) {
		alvo[relTo(root, f)] = true
	}
	var nós []mapx.Node
	for _, n := range g.Nodes {
		if alvo[n.ID] {
			nós = append(nós, n)
		}
	}
	if len(nós) == 0 {
		return nil
	}
	var out []string
	for _, r := range gate.RunWithConfig(cfg.Gates, nós, root, g, cfg) {
		if r.Verdict != gate.Fail {
			continue
		}
		detalhe := r.Detail
		if i := strings.Index(detalhe, ". "); i > 0 && i < 160 {
			detalhe = detalhe[:i+1] // primeira frase basta para reconhecer o problema
		}
		out = append(out, fmt.Sprintf("%s @ %s — %s", r.Gate, r.Target, detalhe))
	}
	return out
}

// mutacaoNaoMedida avisa quando a unidade tem teste e nenhum sinal de mutação ingerido.
//
// A mutação MANUAL — que o `anchors work test` manda fazer — pega o que o autor lembra de
// mutar. Medido numa rodada real: o autor mutou 6 pontos e matou os 6, e mesmo assim uma
// regra escapou; um revisor independente inverteu o desempate dela e os 22 testes
// continuaram verdes. Não foi descuido: ninguém pensa em inverter o próprio critério.
//
// A ferramenta muta o que ninguém pensou. Este aviso não exige que exista uma (nem todo
// stack tem), mas impede que a ausência passe como se o teste estivesse provado — que é o
// que acontece hoje, com `mutation-score` em `~` no meio de outros indeterminados.
func mutacaoNaoMedida(root, unit string) string {
	if unit == "" {
		return ""
	}
	base := strings.TrimSuffix(unit, filepath.Ext(unit))
	temTeste := false
	for _, suf := range []string{".test.ts", ".test.tsx", ".spec.ts", "_test.go", "_test.py", "_spec.rb"} {
		if _, err := os.Stat(filepath.Join(root, base+suf)); err == nil {
			temTeste = true
			break
		}
	}
	if !temTeste {
		return "" // sem teste, o assunto é outro (e o trinca-completa já cobra)
	}
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		return ""
	}
	for _, n := range g.Nodes {
		if n.ID != unit || n.Signal == nil {
			continue
		}
		if n.Signal.MutantsKilled > 0 || n.Signal.MutantsSurvived > 0 {
			return "" // já medido
		}
	}
	return "esta unidade tem teste e NENHUM sinal de mutação ingerido.\n" +
		"  Mutar à mão pega o que você lembra de mutar — e o que escapa é justamente o que\n" +
		"  você não pensou em atacar (medido: 6 mutações manuais, 6 mortas, e ainda assim um\n" +
		"  revisor inverteu um critério e os testes ficaram verdes). Se o stack tiver\n" +
		"  ferramenta, rode-a e `anchors ingest --mutation <relatório>`."
}

// sobPrefixo diz se o caminho está sob algum diretório que o git reportou como novo.
func sobPrefixo(rel string, prefixos []string) bool {
	for _, p := range prefixos {
		if strings.HasPrefix(rel, p) {
			return true
		}
	}
	return false
}
