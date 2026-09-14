package initx

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A FAIXA DE ESCALONADOS é o único pedaço da página do board que decide algo: as colunas
// listam o que existe, e esta faixa julga o que está travado, quem trava quanto, e em que
// ordem responder.
//
// Ela nasceu porque `anchors:needs-user` não tem coluna: o card não estava sem destaque, ele
// DESAPARECIA da página. Ficava aberto no GitHub bloqueando o que dependia dele, e o board
// mostrava um projeto fluindo.
//
// O JS de uma página embutida não tem suíte, e a consequência é conhecida — quatro mutações
// deste comportamento passariam sem nada acusar. Este teste roda o JS de verdade, num `node`,
// e confronta o que ele produz. Sem `node` na máquina, pula: a régua é do CI e de quem tem a
// ferramenta, e travar o `go test` de quem não tem seria pior que a falta do teste.

//go:embed board
var boardParaTeste embed.FS

// dubleDeDOM é o mínimo para a página carregar fora do navegador. Um PROXY, e não um objeto
// enumerado: a página toca dezenas de métodos de DOM ao inicializar, e remendar um por um a
// cada erro custa mais que responder a todos. O que o teste observa é `hidden` e `innerHTML`
// do elemento da faixa.
const dubleDeDOM = `
const els = {};
function mk(id) {
  if (els[id]) return els[id];
  const alvo = { id, hidden: false, innerHTML: '', textContent: '', dataset: {}, style: {} };
  return els[id] = new Proxy(alvo, {
    get(o, k) { return k in o ? o[k] : (typeof k === 'string' ? () => mk(id + ':' + k) : undefined); },
    set(o, k, v) { o[k] = v; return true; },
  });
}
const qualquer = new Proxy(function(){}, {
  get(_, k) { return k === 'matches' ? false : qualquer; },
  set() { return true; }, apply() { return qualquer; }, construct() { return qualquer; },
});
globalThis.document = new Proxy({ getElementById: mk }, { get(o, k) { return k in o ? o[k] : qualquer; } });
globalThis.window = qualquer;
globalThis.localStorage = { getItem: () => null, setItem() {} };
globalThis.fetch = () => new Promise(() => {});
globalThis.setInterval = () => 0;
globalThis.setTimeout = () => 0;
globalThis.matchMedia = () => ({ matches: false, addEventListener() {} });
`

var scriptRE = regexp.MustCompile(`(?s)<script>(.*?)</script>`)

// rodaDesenhaBloqueio executa `desenhaBloqueio` com os itens dados e devolve o estado do
// elemento da faixa.
func rodaDesenhaBloqueio(t *testing.T, itens []map[string]any) (hidden bool, html string) {
	t.Helper()

	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node ausente: o JS da página não pode ser confrontado nesta máquina")
	}

	pagina, err := fs.ReadFile(boardParaTeste, "board/anchors-board.html")
	if err != nil {
		t.Fatalf("página do board: %v", err)
	}

	// A página tem MAIS DE UM bloco `<script>`, e um regex guloso pegaria o `</script>` do
	// último — o primeiro erro que este teste cometeu.
	var js strings.Builder
	for _, m := range scriptRE.FindAllStringSubmatch(string(pagina), -1) {
		js.WriteString(m[1])
		js.WriteString("\n")
	}
	if js.Len() == 0 {
		t.Fatal("nenhum bloco <script> na página: o teste não tem o que confrontar")
	}

	dados, err := json.Marshal(map[string]any{"itens": itens})
	if err != nil {
		t.Fatalf("dados: %v", err)
	}

	prog := dubleDeDOM + js.String() + `
const bloq = mk('bloqueio');
desenhaBloqueio(` + string(dados) + `);
console.log(JSON.stringify({hidden: bloq.hidden, html: bloq.innerHTML}));
`
	arq := filepath.Join(t.TempDir(), "board.mjs")
	if err := os.WriteFile(arq, []byte(prog), 0o600); err != nil {
		t.Fatalf("escrever: %v", err)
	}

	out, err := exec.Command("node", arq).Output()
	if err != nil {
		t.Fatalf("node: %v\n%s", err, out)
	}
	var r struct {
		Hidden bool   `json:"hidden"`
		HTML   string `json:"html"`
	}
	// A última linha: a página pode logar antes.
	linhas := strings.Split(strings.TrimSpace(string(out)), "\n")
	if err := json.Unmarshal([]byte(linhas[len(linhas)-1]), &r); err != nil {
		t.Fatalf("saída do node não é o JSON esperado: %v\n%s", err, out)
	}
	return r.Hidden, r.HTML
}

func TestFaixaDeBloqueio_semEscalonadoNaoAparece(t *testing.T) {
	hidden, html := rodaDesenhaBloqueio(t, []map[string]any{
		{"estado": "anchors:to-do", "numero": 1, "codigo": "A", "titulo": "x", "url": "u"},
	})
	if !hidden {
		t.Error("sem card `needs-user` a faixa deveria estar escondida")
	}
	// Esconder e não LIMPAR deixaria o conteúdo velho pronto para reaparecer na próxima
	// vez que a faixa fosse mostrada — com decisões já resolvidas.
	if html != "" {
		t.Errorf("a faixa escondida deveria estar vazia, e tem %q", html)
	}
}

func TestFaixaDeBloqueio_ordenaPorQuantosCardsATravaTrava(t *testing.T) {
	// O número BAIXO trava pouco, o ALTO trava muito: as duas ordenações possíveis dão
	// resultados opostos aqui. Com dados em que elas concordam, a mutação que ordena por
	// número sobrevive — foi o que aconteceu na primeira versão deste teste.
	//
	// E o `DEC2` depende do `DEC1`: serve para provar que a contagem NÃO conta como
	// bloqueado um card que também está parado esperando uma pessoa.
	hidden, html := rodaDesenhaBloqueio(t, []map[string]any{
		{"estado": "anchors:needs-user", "numero": 12, "codigo": "DEC1", "titulo": "[DEC1] qual vocabulário?", "url": "https://x/12"},
		{"estado": "anchors:needs-user", "numero": 99, "codigo": "DEC2", "titulo": "[DEC2] onde fica o limite?", "needs": []string{"DEC1"}, "url": "https://x/99"},
		{"estado": "anchors:to-do", "numero": 5, "codigo": "Z1", "titulo": "z", "needs": []string{"DEC2"}, "url": "u"},
		{"estado": "anchors:to-do", "numero": 6, "codigo": "Z2", "titulo": "z", "needs": []string{"DEC2"}, "url": "u"},
		{"estado": "anchors:to-do", "numero": 7, "codigo": "Z3", "titulo": "z", "needs": []string{"DEC1"}, "url": "u"},
	})
	if hidden {
		t.Fatal("com dois cards `needs-user` a faixa deveria aparecer")
	}

	i99, i12 := strings.Index(html, "#99"), strings.Index(html, "#12")
	if i99 < 0 || i12 < 0 {
		t.Fatalf("a faixa deveria mostrar os DOIS números; html=%q", html)
	}
	if i99 > i12 {
		t.Error("quem trava mais cards deveria vir primeiro: o #99 trava 2 e o #12 trava 1")
	}
	if !strings.Contains(html, "trava 2 cards") {
		t.Error("o #99 trava dois cards, e a faixa deveria dizer quantos")
	}
	// `trava 1 card<` e não `trava 1 card`: sem o delimitador, o "trava 1 card" casaria
	// dentro de "trava 1 cards" — que a regra do plural não deveria produzir.
	if !strings.Contains(html, "trava 1 card<") {
		t.Error("o #12 trava só o Z3 (o DEC2 também está parado, e não conta como bloqueado)")
	}
	if !strings.Contains(html, "Esperando você · 2") {
		t.Error("a faixa deveria dizer QUANTAS decisões estão paradas")
	}
	// O `[CODIGO]` do título é ruído aqui: o código já aparece na coluna do board, e na
	// faixa o que importa é a pergunta.
	if strings.Contains(html, "[DEC2]") {
		t.Error("o prefixo `[CODIGO]` deveria sair do título na faixa")
	}
}

// `needs-user` ACUMULA com o estado do trabalho, e a faixa tem de ver os dois.
//
// Um card em revisão que precisa de decisão fica com as DUAS labels, e o `estado` do JSON é
// a PRIMEIRA delas. A primeira versão filtrava por `estado === "anchors:needs-user"` e o
// card SUMIA da faixa — ficava só na coluna IN REVIEW, indistinguível de um card sendo
// revisado normalmente.
//
// Medido: o #311 do projeto de referência esperava uma decisão do usuário e não aparecia
// aqui. O usuário o viu no board sem nada indicando que ninguém ia tocá-lo.
func TestFaixaDeBloqueio_cardEscaladoComEstadoDeTrabalhoAparece(t *testing.T) {
	hidden, html := rodaDesenhaBloqueio(t, []map[string]any{
		{"estado": "anchors:in-review", "escalado": true, "numero": 311, "codigo": "NTDSN",
			"titulo": "[NTDSN] o que sai da VPC", "url": "https://x/311"},
		{"estado": "anchors:to-do", "escalado": true, "numero": 443, "codigo": "DEC",
			"titulo": "[DEC] a revisão achou defeito crítico", "url": "https://x/443"},
	})
	if hidden {
		t.Fatal("dois cards escalados e a faixa não apareceu")
	}
	if !strings.Contains(html, "#311") {
		t.Error("o card `in-review` + `needs-user` tem de aparecer — é o caso que sumia")
	}
	if !strings.Contains(html, "#443") {
		t.Error("o card `to-do` + `needs-user` também")
	}
	if !strings.Contains(html, "Esperando você · 2") {
		t.Error("a contagem tem de incluir os dois")
	}
	// E dizer ONDE o card parou: `in-review` esperando decisão é diferente de `to-do`
	// esperando decisão — no primeiro há trabalho meio feito, no segundo não.
	if !strings.Contains(html, "parado em in-review") {
		t.Error("a faixa deveria dizer em que estado o card parou")
	}
	// `to-do` não é "parado em": é o estado normal de quem ainda não foi pego.
	if strings.Contains(html, "parado em to-do") {
		t.Error("`to-do` não é um estado de trabalho interrompido")
	}
}

// O CARD QUE ESPERA VOCÊ é diferente do que espera TRABALHO JÁ DECIDIDO.
//
// Quando a decisão sai e gera trabalho, o card da mudança nasce com
// `anchors:desbloqueia-<n>` e o bloqueado continua com `needs-user` — porque remover a
// label antes da entrega faria o claim devolver o card, e o agente esbarraria no mesmo
// impasse.
//
// Os dois casos se parecem na tela se a faixa não os separar: um pede a atenção do usuário
// AGORA, o outro não pede nada — só espera uma fila andar.
func TestFaixaDeBloqueio_distingueQuemEsperaVoceDeQuemEsperaTrabalho(t *testing.T) {
	_, html := rodaDesenhaBloqueio(t, []map[string]any{
		{"estado": "anchors:in-review", "escalado": true, "numero": 311, "codigo": "NTDSN",
			"titulo": "[NTDSN] o que sai da VPC", "url": "https://x/311"},
		{"estado": "anchors:to-do", "escalado": true, "numero": 500, "codigo": "DEC",
			"titulo": "[DEC] qual vocabulário?", "url": "https://x/500"},
		// O card da mudança: não está escalado, e DESTRAVA o 311.
		{"estado": "anchors:to-do", "numero": 444, "codigo": "FIX", "destrava": "311",
			"titulo": "[destrava #311] try/catch por token", "url": "https://x/444"},
	})

	if !strings.Contains(html, "espera a entrega do #444") {
		t.Error("o card cuja decisão já saiu deveria dizer POR QUEM espera")
	}
	// O #500 ainda espera a PESSOA — não pode ganhar a mesma frase.
	i500 := strings.Index(html, "#500")
	i444 := strings.Index(html, "espera a entrega do #444")
	if i500 >= 0 && i444 >= 0 && i444 > i500 {
		t.Error("a frase de espera colou no card errado — ela é do #311, não do #500")
	}
	// E o card da MUDANÇA não entra na faixa: ele não espera ninguém, é trabalho normal.
	if strings.Contains(html, "#444 ") || strings.Contains(html, ">#444<") {
		t.Error("o card que destrava não é um card escalado — ele não deveria aparecer aqui")
	}
}

func TestFaixaDeBloqueio_escapaTituloHostil(t *testing.T) {
	// O título vem de uma issue, e qualquer pessoa com acesso ao repositório escreve
	// issue. A página é servida no Pages do projeto — um título com marcação executaria
	// no navegador de quem abre o board.
	_, html := rodaDesenhaBloqueio(t, []map[string]any{
		{"estado": "anchors:needs-user", "numero": 7, "codigo": "X",
			"titulo": "<img src=x onerror=alert(1)>", "url": "https://x/7"},
	})
	if strings.Contains(html, "<img src=x") {
		t.Errorf("o título deveria ser escapado; html=%q", html)
	}
}
