package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// plan-source-declared: um plano que NOMEIA uma fonte na prosa tem de declarar, no
// `needs:`, o plano que a constrói.
//
// O CASO REAL, medido no blue-eyes. O plano 0008 (Frontend/Web) dizia:
//
//	Fonte: **GA4**, e ela é a exceção arquitetural do projeto.
//
// e declarava `needs: plans/0005-home-e-indice.md` — só isso. O adaptador do GA4 vinha do
// plano 0002, e essa dependência existia SÓ NA PROSA.
//
// O que aconteceu: o 0002 foi revisado (`PLTFR-R0002`) e o `Ga4Adapter` foi REMOVIDO, com
// um argumento correto — nenhum documento do projeto o sustentava. O 0008 continuou
// intacto, dependendo de uma fonte que ninguém mais ia construir. Os dois planos ficaram
// internamente coerentes, e a contradição só apareceu meses depois, ao começar o 0008.
//
// POR QUE NENHUM GATE PEGAVA. O `dependency-honored` confronta o `needs:` DECLARADO; aqui
// o defeito é o `needs:` que FALTA. E o `plan-seeds-valid` olha o que o plano promete
// criar, não o que ele promete consumir.
//
// É a versão entre PLANOS do passo 5 do guia de revisão ("a prosa envelhece"): ninguém
// alterou o 0008, e ele ficou errado assim mesmo — porque outro documento mudou.
//
// O QUE ELE MEDE. Uma linha de fonte (`Fonte:`/`Fontes:`) nomeia uma ou mais fontes em
// negrito. Para cada uma, o gate procura o adaptador correspondente entre as sementes de
// TODOS os planos. Se o adaptador existe noutro plano e este não o declara no `needs:`,
// reprova.
//
// O que ele NÃO cobra: source cujo adaptador ninguém semeia (pode ser source de um plano
// futuro), e a ordem das fases — quem faz isso é o `fase-ordenada`.

// reSourceLine casa a declaração de fonte na prosa de um plano.
var reSourceLine = regexp.MustCompile(`(?m)^Fontes?:\s*(.+)$`)

// reBoldSource extrai cada fonte nomeada da linha. O negrito é a convenção que os
// planos já usam (`**Prometheus**`, `**ELK**`, `**GA4**`) — e ela é o que separa a source
// da prosa que a explica.
var reBoldSource = regexp.MustCompile(`\*\*([^*]+)\*\*`)

func checkPlanSourceDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, "não é um plano — o gate confronta a FONTE que um plano nomeia"
	}
	if g == nil {
		return Pending, "sem mapa — o gate precisa das sementes dos outros planos"
	}

	sources := namedSources(content)
	if len(sources) == 0 {
		return Skip, "o plano não nomeia fonte (`Fonte:` / `Fontes:`) — nada a confrontar"
	}

	// De onde vem cada adaptador: nome da source -> plano que o semeia.
	owners := adapterOwners(g)
	if len(owners) == 0 {
		return Pending, "nenhum plano semeia adaptador — o gate não tem contra o que confrontar"
	}

	declared := map[string]bool{}
	for _, d := range n.Needs {
		declared[d] = true
	}

	var missing []string
	for _, f := range sources {
		owner, ok := owners[normalizeSource(f)]
		if !ok || owner == n.ID || declared[owner] {
			continue
		}
		missing = append(missing, fmt.Sprintf("%s (o adaptador está em %s)", f, owner))
	}
	if len(missing) == 0 {
		return Pass, fmt.Sprintf("as %d fonte(s) nomeada(s) têm o plano do adaptador no `needs:`", len(sources))
	}
	sort.Strings(missing)

	return Fail, "o plano NOMEIA fonte cujo adaptador vem de outro plano, e não o declara no `needs:`:\n" +
		"      " + strings.Join(missing, "\n      ") + "\n" +
		"      A dependência existe só na PROSA. Medido: o 0008 nomeava GA4 e declarava apenas o 0005;\n" +
		"      o 0002 removeu o `Ga4Adapter` numa revisão, e o 0008 ficou dependendo de uma fonte que\n" +
		"      ninguém ia construir — sem que nada acusasse, porque os dois planos seguiram coerentes."
}

// namedSources extrai as fontes em negrito das linhas `Fonte:`/`Fontes:`.
func namedSources(content string) []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range reSourceLine.FindAllStringSubmatch(content, -1) {
		for _, f := range reBoldSource.FindAllStringSubmatch(m[1], -1) {
			nome := strings.TrimSpace(f[1])
			if nome != "" && !seen[nome] {
				seen[nome] = true
				out = append(out, nome)
			}
		}
	}
	return out
}

// adapterOwners mapeia o nome da fonte para o plano que semeia o adaptador dela.
//
// A convenção é o nome do arquivo: `Ga4Adapter.spec.md` -> `ga4`, `ElkAdapter.spec.md` ->
// `elk`. Ela vale porque é a que os três adaptadores existentes já seguem, e porque o gate
// falha para o lado seguro: um nome fora do padrão não casa e o gate não acusa nada.
func adapterOwners(g *mapx.Graph) map[string]string {
	owners := map[string]string{}
	// As sementes são ARESTAS (`EdgeSeeds`), não campo do nó: é o que liga a decisão de
	// fazer ao artefato que nasce, e o que sobrevive quando a spec ainda não existe.
	for _, e := range g.Edges {
		if e.Type != mapx.EdgeSeeds {
			continue
		}
		base := filepathBase(e.To)
		if !strings.HasSuffix(base, "Adapter.spec.md") {
			continue
		}
		source := strings.TrimSuffix(base, "Adapter.spec.md")
		owners[normalizeSource(source)] = e.From
	}
	return owners
}

// normalizeSource reduz o nome ao que compara: minúsculas, sem espaço e sem pontuação.
// `GA4` casa com `Ga4`, e `CloudWatch` com `cloudwatch`.
func normalizeSource(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
