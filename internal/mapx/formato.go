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
	FormatoAtual = 2

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
			"%s está no formato %d, e este binário lê até o %d.\n\n"+
				"  O mapa foi gravado por um Anchors mais NOVO que o seu. Continuar\n"+
				"  reescreveria o arquivo no formato antigo, e o que a versão nova\n"+
				"  gravou se perderia em silêncio — inclusive os carimbos de julgamento,\n"+
				"  que o `check` teria de refazer.\n\n"+
				"  Atualize: `brew upgrade anchors` ou\n"+
				"  `go install github.com/co2-lab/anchors/cmd/anchors@latest`",
			e.Path, e.Encontrado, e.Escrevo)
	}
	return fmt.Sprintf(
		"%s está no formato %d, e este binário lê a partir do %d.\n\n"+
			"  O mapa é de uma versão anterior do Anchors e precisa ser MIGRADO.\n"+
			"  A migração é automática e roda uma vez:\n\n"+
			"      anchors migrate\n\n"+
			"  Ela reescreve o arquivo no formato atual e o deixa pronto para commit.",
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
