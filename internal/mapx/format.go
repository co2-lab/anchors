package mapx

import "fmt"

// A VERSÃO DE FORMATO do `anchors.graph.yaml`, e por que ela não é a versão do binário.
//
// O campo `version:` do mapa sempre existiu e sempre valeu 1 — ninguém o usava para nada.
// Agora ele é o contrato: diz em que FORMATO o arquivo está escrito, e o binário recusa o
// que não sabe ler.
//
// POR QUE ISSO, E NÃO UM TETO DECLARADO NO `anchors.yaml`.
//
// A alternativa era o time declarar `max_version: 0.1.90` à mão. Ela funciona e tem um
// defeito que aparece com o tempo: depende de alguém lembrar de subir o teto a cada
// migração, e quem esquece passa a barrar atualização legítima — o freio vira atrito, e
// atrito se contorna.
//
// O formato não depende de ninguém lembrar. Ele sobe quando o Anchors muda o arquivo de um
// jeito que a versão anterior não entende, e é o próprio produto que sabe quando isso
// aconteceu.
//
// QUANDO SUBIR. Só quando a mudança torna o arquivo ILEGÍVEL para a versão anterior:
// renomear uma chave, mudar a forma de um valor, remover um campo que era obrigatório.
//
// NÃO sobe por campo novo opcional — um binário velho o ignora, e ignorar não é quebrar.
// Subir por isso tornaria toda adição uma barreira, e o número perderia o significado.
const (
	// FormatoAtual é o que este binário ESCREVE.
	FormatoAtual = 4

	// FormatoMinimoLegivel é o mais antigo que ele lê sem migrar.
	//
	// Hoje é igual ao atual: o formato 1 não é lido, ele é MIGRADO — as chaves em
	// português (`gerado_por`, `code_declarado`, `julgamentos`) foram renomeadas, e ler os
	// dois nomes para sempre faria o arquivo nunca se consertar.
	FormatoMinimoLegivel = 2
)

// ErroDeFormato distingue "não sei ler isto" de "o arquivo está corrompido".
//
// A diferença importa para quem lê a mensagem: um arquivo do FUTURO pede atualização do
// binário, e um do PASSADO pede migração. Tratar os dois como "erro ao carregar o mapa"
// mandaria a pessoa procurar corrupção onde há só versão.
type ErroDeFormato struct {
	Path       string
	Encontrado int
	Escrevo    int
	MinLegivel int
}

func (e *ErroDeFormato) Error() string {
	if e.Encontrado > e.Escrevo {
		return fmt.Sprintf(
			"%s is in format %d, and this binary reads up to %d.\n\n"+
				"  The map was written by an Anchors NEWER than yours. Continuing\n"+
				"  would rewrite the file in the old format, and what the newer version\n"+
				"  wrote would be lost in silence — the judgment stamps included,\n"+
				"  which `check` would have to redo.\n\n"+
				"  Update: `brew upgrade anchors` or\n"+
				"  `go install github.com/co2-lab/anchors/cmd/anchors@latest`",
			e.Path, e.Encontrado, e.Escrevo)
	}
	return fmt.Sprintf(
		"%s is in format %d, and this binary reads from %d onwards.\n\n"+
			"  The map is from an earlier version of Anchors and needs to be MIGRATED.\n"+
			"  The migration is automatic and runs once:\n\n"+
			"      anchors migrate\n\n"+
			"  It rewrites the file in the current format and leaves it ready to commit.",
		e.Path, e.Encontrado, e.MinLegivel)
}

// ConfereFormato responde se este binário pode operar sobre o mapa encontrado.
//
// Chamado na CARGA, antes de qualquer decisão: um mapa que não se entende não deve ser
// parcialmente interpretado — é assim que se perde dado sem nada acusar.
func ConfereFormato(path string, encontrado int) error {
	// Formato 0 é arquivo sem `version:` — os primeiros mapas do produto. Trata como 1,
	// que é o que a migração espera.
	if encontrado == 0 {
		encontrado = 1
	}
	if encontrado >= FormatoMinimoLegivel && encontrado <= FormatoAtual {
		return nil
	}
	return &ErroDeFormato{
		Path: path, Encontrado: encontrado,
		Escrevo: FormatoAtual, MinLegivel: FormatoMinimoLegivel,
	}
}
