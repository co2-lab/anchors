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
