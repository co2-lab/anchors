package mapx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// O CASO REAL, medido no blue-eyes.
//
// O plano 0010 (Redis) semeia `packages/lambdas/redis/InstanceList.spec.md`. O plano 0009
// (Database) já tinha entregue `packages/lambdas/database/InstanceList.spec.md`.
//
// O fallback "resolve por nome quando ele é único" — feito para CITAÇÃO EM PROSA, onde o
// autor escreve só o nome do arquivo — foi aplicado a uma semente que traz o CAMINHO
// INTEIRO. O alvo do 0010 não existia, o nome base era único no repositório, e o mapa
// ganhou uma aresta do plano 0010 para a spec do plano 0009.
//
// O efeito medido: o `anchors next` respondeu "2 de 4", depois "3 de 3", para o mesmo
// plano cujas QUATRO sementes não existem — e nenhum dos dois números estava certo. Pior
// que o número errado: o mapa afirma que o 0010 semeia specs que são do 0009, e todo gate
// relacional passa a confrontar o par errado.
//
// A régua: caminho DECLARADO é caminho. Só o nome nu (sem separador) pode ser resolvido
// por busca — e o `extractSeeds` já garante que nome nu não vira semente.
func TestBuild_semeanteComCaminhoNaoResolvePorNomeDeOutroDiretorio(t *testing.T) {
	files := []scan.File{
		{
			Path:  "plans/0010-redis.md",
			Kind:  "plan",
			Seeds: []string{"packages/lambdas/redis/InstanceList.spec.md"},
		},
		{Path: "packages/lambdas/database/InstanceList.spec.md", Kind: "spec"},
	}

	g := Build(files, &config.Config{}, nil)

	for _, e := range g.Edges {
		if e.Type != EdgeSeeds || e.From != "plans/0010-redis.md" {
			continue
		}
		if e.To == "packages/lambdas/database/InstanceList.spec.md" {
			t.Fatalf("o plano do REDIS ficou semeando a spec do DATABASE:\n  %s → %s\n"+
				"  o caminho declarado era packages/lambdas/redis/InstanceList.spec.md", e.From, e.To)
		}
	}
}

// A semente cujo caminho declarado EXISTE continua virando aresta — é o caso normal.
func TestBuild_semeanteComCaminhoQueExisteViraAresta(t *testing.T) {
	files := []scan.File{
		{
			Path:  "plans/0010-redis.md",
			Kind:  "plan",
			Seeds: []string{"packages/lambdas/redis/InstanceList.spec.md"},
		},
		{Path: "packages/lambdas/redis/InstanceList.spec.md", Kind: "spec"},
	}

	g := Build(files, &config.Config{}, nil)

	achou := false
	for _, e := range g.Edges {
		if e.Type == EdgeSeeds && e.From == "plans/0010-redis.md" &&
			e.To == "packages/lambdas/redis/InstanceList.spec.md" {
			achou = true
		}
	}
	if !achou {
		t.Error("a semente cujo caminho existe não virou aresta")
	}
}

// E o caminho que NÃO existe em lugar nenhum também não vira aresta por adivinhação: o
// plano promete criar um arquivo, e o mapa não deve inventar que ele já está lá.
func TestBuild_semeanteInexistenteNaoViraArestaParaOutroArquivo(t *testing.T) {
	files := []scan.File{
		{
			Path:  "plans/0010-redis.md",
			Kind:  "plan",
			Seeds: []string{"packages/lambdas/redis/Coisa.spec.md"},
		},
		{Path: "packages/outro/lugar/Coisa.spec.md", Kind: "spec"},
	}

	g := Build(files, &config.Config{}, nil)

	for _, e := range g.Edges {
		if e.Type == EdgeSeeds && e.To == "packages/outro/lugar/Coisa.spec.md" {
			t.Errorf("resolveu por nome um caminho declarado que não existe: %s → %s", e.From, e.To)
		}
	}
}
