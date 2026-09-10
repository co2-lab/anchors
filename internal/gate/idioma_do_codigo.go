package gate

import (
	"regexp"
	"strings"
)

// --- o gate que impede o código de voltar a misturar idiomas ---
//
// O Anchors nasceu com o código em português e migrou para o inglês, porque o projeto vai
// ser aberto a outros devs: um contribuidor externo não deveria precisar de português para
// ler `carimbo()` ou `DonoDoArquivo`.
//
// Esse trabalho se desfaz sozinho se nada o defender. Basta um PR de quem não sabe da
// regra — e ela não está em lugar nenhum que o compilador leia.
//
// O QUE ESTE GATE OLHA, e o que ele deliberadamente ignora:
//
//	identificadores   func, type, var, const — SIM
//	comentários       NÃO: eles carregam as medições e os porquês, no idioma do time
//	strings           NÃO: mensagem ao usuário vai pelo i18n, e traduzi-la aqui seria
//	                  o erro oposto
//
// A régua é a mesma do resto: o que se ESCREVE no código é inglês; o que se LÊ é
// traduzido.
//
// COMO A DETECÇÃO FUNCIONA, e por que não é um dicionário.
//
// A primeira tentativa comparou cada palavra contra o dicionário do sistema
// (`/usr/share/dict/words`). Ela acusava `AbsRoot`, `FileOwner` e `Classify` — palavras
// COMPOSTAS em inglês, que nenhum dicionário comum tem. Medido: 491 acusações para 47
// casos reais, 90% de falso positivo.
//
// Um gate com essa taxa é pior que gate nenhum: a saída barata vira desligá-lo.
//
// O que funciona é cruzar dois sinais que o inglês não produz:
//
//  1. MORFOLOGIA — terminações que só existem em português (-ção, -ável, -mento).
//     Nenhuma palavra inglesa termina em `-cao` ou `-ivel`.
//  2. RADICAIS — o vocabulário do próprio domínio, que se repete no código
//     (alvo, carimbo, peça, trinca).
//
// Nenhum dos dois é embutido de fora: são 60 linhas de tabela, versionadas, que crescem
// quando alguém introduz um termo novo.

// ptEndings são sufixos que o inglês não tem.
//
// Cada um foi confirmado contra o código real antes de entrar: `-ura` casa `leitura` e
// `estrutura`, mas também casaria `future` — por isso o casamento exige que a palavra
// INTEIRA não seja inglesa comum, e a lista de exceções abaixo cobre o resto.
var ptEndings = regexp.MustCompile(
	`(cao|coes|ada|ado|ados|adas|ida|ido|idos|idas|agem|ancia|encia|avel|ivel|` +
		`dade|mento|eiro|eira|oso|osa|eza|inho|inha|ismo|ista)$`)

// radicaisPT é o vocabulário do domínio que aparece em identificador.
//
// Substantivos e verbos que este projeto usava, e que qualquer projeto em português
// usaria. A lista é curta de propósito: ela não tenta ser um dicionário, só cobrir o que
// de fato vira nome de função.
var radicaisPT = map[string]bool{
	// substantivos do domínio
	"alvo": true, "arquivo": true, "camada": true, "carimbo": true, "cenario": true,
	"codigo": true, "conteudo": true, "dono": true, "duble": true, "escopo": true,
	"fase": true, "fonte": true, "guia": true, "irmao": true, "linha": true,
	"mapa": true, "marca": true, "nivel": true, "peca": true, "plano": true,
	"ponta": true, "prazo": true, "prosa": true, "raiz": true, "regra": true,
	"rota": true, "saida": true, "secao": true, "sigla": true, "trecho": true,
	"trinca": true, "dispensa": true, "entrega": true, "achado": true, "aresta": true,
	"chave": true, "destino": true, "divida": true, "etapa": true, "faixa": true,
	"fila": true, "juizo": true, "motivo": true, "mudanca": true, "ordem": true,
	"passo": true, "pergunta": true, "placar": true, "rastro": true, "sinal": true,
	"titulo": true, "unidade": true, "usuario": true, "vazio": true,
	// verbos, na forma que aparece em nome de função
	"abre": true, "acha": true, "aplica": true, "avisa": true, "busca": true,
	"carrega": true, "cita": true, "cobra": true, "cobre": true, "conta": true,
	"corrige": true, "cria": true, "deixa": true, "escreve": true, "exige": true,
	"falta": true, "fecha": true, "ignora": true, "impede": true, "junta": true,
	"lista": true, "marca_": true, "mede": true, "mostra": true, "move": true,
	"nasce": true, "olha": true, "passa": true, "pega": true, "prova": true,
	"puxa": true, "quebra": true, "recusa": true, "registra": true, "roda": true,
	"segue": true, "tira": true, "toca": true, "troca": true, "usa": true,
	"vale": true, "varre": true, "vira": true, "volta": true, "normaliza": true,
	"imprime": true, "verifica": true, "confronta": true, "declara": true,
}

// exceptions são palavras que a morfologia acusaria e são inglês legítimo.
//
// `dado` casa `-ado`, mas `upgraded` também casaria se estivesse solto. A lista existe
// para que o gate não obrigue ninguém a renomear algo que já está certo.
var exceptions = map[string]bool{
	"data": true, "meta": true, "beta": true, "delta": true, "media": true,
	"schema": true, "comma": true, "gamma": true, "lambda": true, "alpha": true,
	"extra": true, "area": true, "idea": true, "arena": true, "camera": true,
	"formula": true, "quota": true, "vista": true, "list": true, "cascade": true,
	"grade": true, "made": true, "trade": true, "upgrade": true, "facade": true,
	"decade": true, "parade": true, "shade": true, "blade": true,
	// Inglês que a morfologia acusa, achado varrendo o próprio projeto:
	//   resolver  -er, mas é substantivo inglês (quem resolve)
	//   moves     plural de `move`, não de `mova`
	//   divider   -er, mesmo caso de `resolver`
	//   header    idem — e ele aparece em todo parser
	//   owner     idem
	"resolver": true, "moves": true, "divider": true, "header": true,
	"owner": true, "render": true, "order": true, "under": true,
	"other": true, "after": true, "over": true, "later": true,
	"filter": true, "counter": true, "parser": true, "writer": true,
	"reader": true, "buffer": true, "marker": true, "number": true,
	"member": true, "folder": true, "holder": true, "helper": true,
}

// goDeclaration casa uma DECLARAÇÃO de identificador.
//
// `^` exige coluna zero: um `func` dentro de string ou de comentário indentado não é
// declaração, e acusá-lo faria o gate reprovar exemplo de documentação.
var goDeclaration = regexp.MustCompile(
	`(?m)^(?:func\s+(?:\([^)]*\)\s*)?|type\s+|var\s+|const\s+)([A-Za-z][A-Za-z0-9]*)`)

// palavrasDoIdentificador quebra camelCase/PascalCase em palavras.
//
// `AbsRoot` → [Abs, Root]. Sem lookahead: o regexp do Go (RE2) não o suporta, e a
// alternativa — `[A-Z][a-z]*|[a-z]+` — produz o mesmo recorte para os casos que
// importam. Uma sigla como `ID` vira dois pedaços de uma letra, que a régua de tamanho
// descarta.
var palavrasDoIdentificador = regexp.MustCompile(`[A-Z][a-z]*|[a-z]+`)

// PalavraEhPT diz se uma palavra isolada é portuguesa.
//
// Exportada para o teste: é a decisão mais delicada do gate, e ela merece ser exercitada
// caso a caso em vez de só através do resultado final.
func PalavraEhPT(w string) bool {
	w = strings.ToLower(w)
	if len(w) < 3 || exceptions[w] {
		return false
	}
	if radicaisPT[w] || ptEndings.MatchString(w) {
		return true
	}
	// O PLURAL do radical: `peca` está na tabela, `pecas` não estaria. Listar as duas
	// formas de cada palavra dobraria a tabela sem acrescentar informação.
	if strings.HasSuffix(w, "s") && radicaisPT[strings.TrimSuffix(w, "s")] {
		return true
	}
	// O INFINITIVO: `desenvolver`, `carregar`, `imprimir`. O inglês também tem palavras
	// em `-er` (`owner`, `header`, `render`), então a régua exige que o RADICAL sem a
	// terminação seja português — `desenvolv` não é inglês, `own` é.
	//
	// Sem isso, `pecasPorDesenvolver` escapava: nenhuma das três palavras casava.
	for _, term := range []string{"ar", "er", "ir"} {
		if !strings.HasSuffix(w, term) {
			continue
		}
		raiz := strings.TrimSuffix(w, term)
		if len(raiz) < 4 {
			continue // curto demais: `par`, `ver`, `air` não dizem nada
		}
		if radicaisPT[raiz] || radicaisPT[raiz+"a"] || radicaisPT[raiz+"e"] {
			return true
		}
		// Radicais que só aparecem no infinitivo, e não como substantivo.
		if infinitivosPT[w] {
			return true
		}
	}
	return false
}

// infinitivosPT são verbos que aparecem em nome de função e cujo radical não é
// substantivo — `desenvolver` não tem `desenvolve` na tabela de radicais.
var infinitivosPT = map[string]bool{
	"desenvolver": true, "confrontar": true, "carregar": true, "imprimir": true,
	"escrever": true, "resolver": true, "concluir": true, "reabrir": true,
	"atribuir": true, "distribuir": true, "construir": true, "destruir": true,
	"promover": true, "remover": true, "resumir": true, "traduzir": true,
	"reduzir": true, "produzir": true, "conduzir": true, "corrigir": true,
	"decidir": true, "definir": true, "exibir": true, "inserir": true,
	"garantir": true, "permitir": true, "repetir": true, "seguir": true,
}

// IdentificadorEhPT diz se um identificador tem alguma palavra portuguesa.
func IdentificadorEhPT(ident string) (string, bool) {
	for _, w := range palavrasDoIdentificador.FindAllString(ident, -1) {
		if PalavraEhPT(w) {
			return strings.ToLower(w), true
		}
	}
	return "", false
}

// IdentificadoresPT devolve os identificadores em português declarados no conteúdo.
//
// Devolve o identificador E a palavra que o acusou: sem ela, quem lê a reprovação tem de
// adivinhar qual parte do nome está errada — e num `checkPlanoAlteradoJustificado` são
// quatro candidatas.
func IdentificadoresPT(conteudo string) map[string]string {
	out := map[string]string{}
	for _, m := range goDeclaration.FindAllStringSubmatch(conteudo, -1) {
		if w, ehPT := IdentificadorEhPT(m[1]); ehPT {
			out[m[1]] = w
		}
	}
	return out
}
