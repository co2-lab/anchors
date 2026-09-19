package mapcmd

import (
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// staleMapNodes devolve os caminhos cujo conteúdo mudou desde que o mapa foi construído.
//
// PRESENÇA e IDADE são perguntas diferentes, e o hook só fazia a primeira. O pre-commit já
// barra o arquivo REGIDO que está FORA do mapa — o arquivo novo, de quem nunca rodou
// `anchors map build`. O que ele não via é o arquivo que ESTÁ no mapa carregando a `rev` de
// uma versão anterior: o `map build` rodou, e o trabalho continuou depois dele.
//
// MEDIDO no projeto de referência: das 12 reprovações do pipeline `gates`, SETE foram o mapa
// desatualizado — a maior causa isolada. E nas sete o agente tinha commitado o mapa. Ninguém
// esqueceu de gerá-lo; todos geraram cedo demais e seguiram editando.
//
// O custo é o pior tipo: o trabalho está CERTO, e o CI reprova por causa da foto. São ~6
// minutos de pipeline para descobrir o que um hash responde na hora — e o agente que recebe
// a reprovação vai procurar defeito no trabalho, porque é o que uma reprovação significa.
//
// POR QUE COMPARAR `rev`. É o sha256 curto do conteúdo (`scan.File.Rev`, `mapx.Node.Rev`):
// hash contra hash, sem olhar texto, linguagem ou idioma. Não usa data — `updated_at` muda
// num `git checkout` sem o conteúdo mudar, e não muda numa edição que preserve o mtime.
func StaleMapNodes(root string, g *mapx.Graph, cfg *config.Config) []string {
	if g == nil {
		return nil
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		// Sem varredura não há o que comparar, e ERRAR AQUI SERIA PIOR: esta é uma
		// conferência auxiliar, e derrubar o `check` por causa dela esconderia o
		// relatório que o usuário veio buscar.
		return nil
	}
	current := make(map[string]string, len(files))
	for _, f := range files {
		current[f.Path] = f.Rev
	}

	var stale []string
	for _, n := range g.Nodes {
		r, ok := current[n.ID]
		// AUSENTE não é desta conferência: um nó cujo arquivo sumiu é outro caso (o
		// `map build` o remove), e acusá-lo aqui misturaria dois problemas numa
		// mensagem só. `rev` vazia dos dois lados também não diz nada.
		if !ok || r == "" || n.Rev == "" {
			continue
		}
		if r != n.Rev {
			stale = append(stale, n.ID)
		}
	}
	return stale
}
