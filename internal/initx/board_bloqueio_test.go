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
    get(o, k) {
      if (k in o) return o[k];
      // 'querySelectorAll' devolve LISTA. Um método genérico que devolve outro nó faz o
      // spread '[...]' da página estourar, e o erro aponta para o código dela em vez do
      // dublê — o sintoma é "exit status 1" apontando uma linha que está correta.
      if (k === 'querySelectorAll') return () => [];
      return typeof k === 'string' ? () => mk(id + ':' + k) : undefined;
    },
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
	return rodaNaPagina(t, "bloqueio", "desenhaBloqueio", itens)
}

// rodaNaPagina executa uma função da página sobre os itens dados e devolve o elemento que
// ela escreveu. Generalizado do `desenhaBloqueio` quando a marca do card na COLUNA passou a
// precisar do mesmo aparato — o dublê de DOM e a extração do script são os mesmos.
func rodaNaPagina(t *testing.T, elemento, fn string, itens []map[string]any) (hidden bool, html string) {
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

	dados, err := json.Marshal(map[string]any{"items": itens})
	if err != nil {
		t.Fatalf("dados: %v", err)
	}

	prog := dubleDeDOM + js.String() + `
const alvo = mk('` + elemento + `');
` + fn + `(` + string(dados) + `);
console.log(JSON.stringify({hidden: alvo.hidden, html: alvo.innerHTML}));
`
	arq := filepath.Join(t.TempDir(), "board.mjs")
	if err := os.WriteFile(arq, []byte(prog), 0o600); err != nil {
		t.Fatalf("escrever: %v", err)
	}

	// `CombinedOutput` e não `Output`: o `Output` engole o stderr, e um erro dentro do JS
	// vira "exit status 1" sem dizer onde — dez minutos procurando o que a mensagem já
	// teria dito.
	out, err := exec.Command("node", arq).CombinedOutput()
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
		{"state": "anchors:to-do", "number": 1, "code": "A", "title": "x", "url": "u"},
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
		{"state": "anchors:needs-user", "number": 12, "code": "DEC1", "title": "[DEC1] qual vocabulário?", "url": "https://x/12"},
		{"state": "anchors:needs-user", "number": 99, "code": "DEC2", "title": "[DEC2] onde fica o limite?", "needs": []string{"DEC1"}, "url": "https://x/99"},
		{"state": "anchors:to-do", "number": 5, "code": "Z1", "title": "z", "needs": []string{"DEC2"}, "url": "u"},
		{"state": "anchors:to-do", "number": 6, "code": "Z2", "title": "z", "needs": []string{"DEC2"}, "url": "u"},
		{"state": "anchors:to-do", "number": 7, "code": "Z3", "title": "z", "needs": []string{"DEC1"}, "url": "u"},
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
		{"state": "anchors:in-review", "escalated": true, "number": 311, "code": "NTDSN",
			"title": "[NTDSN] o que sai da VPC", "url": "https://x/311"},
		{"state": "anchors:to-do", "escalated": true, "number": 443, "code": "DEC",
			"title": "[DEC] a revisão achou defeito crítico", "url": "https://x/443"},
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
		{"state": "anchors:in-review", "escalated": true, "number": 311, "code": "NTDSN",
			"title": "[NTDSN] o que sai da VPC", "url": "https://x/311"},
		{"state": "anchors:to-do", "escalated": true, "number": 500, "code": "DEC",
			"title": "[DEC] qual vocabulário?", "url": "https://x/500"},
		// O card da mudança: não está escalado, e DESTRAVA o 311.
		{"state": "anchors:to-do", "number": 444, "code": "FIX", "unblocks": "311",
			"title": "[destrava #311] try/catch por token", "url": "https://x/444"},
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
		{"state": "anchors:needs-user", "number": 7, "code": "X",
			"title": "<img src=x onerror=alert(1)>", "url": "https://x/7"},
	})
	if strings.Contains(html, "<img src=x") {
		t.Errorf("o título deveria ser escapado; html=%q", html)
	}
}

// O CARD ESCALADO precisa de marca NA COLUNA, não só na faixa.
//
// A faixa lista o que espera uma pessoa; a coluna é onde se olha para saber o que está
// acontecendo com um card específico. Sem marca ali, um card parado esperando decisão é
// visualmente IDÊNTICO ao card ao lado que está sendo trabalhado.
//
// Medido: o usuário procurou o #311 na coluna, encontrou, e perguntou "e o destaque no
// card?" — a faixa mostrava, e o card não dizia nada.
func TestCardEscalado_temMarcaNaColuna(t *testing.T) {
	_, html := rodaNaPagina(t, "colunas", "desenha", []map[string]any{
		{"state": "anchors:in-review", "escalated": true, "number": 311, "code": "NTDSN",
			"title": "[NTDSN] o que sai da VPC", "url": "u"},
		{"state": "anchors:in-review", "number": 312, "code": "OUTRO",
			"title": "[OUTRO] revisão normal", "url": "u"},
	})

	if !strings.Contains(html, `class="card escalado"`) {
		t.Error("o card escalado precisa de classe própria — a borda é o que o distingue " +
			"na varredura da coluna")
	}
	if !strings.Contains(html, "esperando você") {
		t.Error("o selo precisa DIZER o que o card espera: uma borda colorida sozinha " +
			"exige que quem vê já saiba o que a cor significa")
	}
	// E o card normal NÃO pode ganhar nenhum dos dois — uma marca que aparece em tudo não
	// marca nada.
	if n := strings.Count(html, `class="card escalado"`); n != 1 {
		t.Errorf("só o card escalado deveria ter a classe; apareceu %d vezes", n)
	}
	if n := strings.Count(html, "esperando você"); n != 1 {
		t.Errorf("só o card escalado deveria ter o selo; apareceu %d vezes", n)
	}
}

// O ROADMAP nao abria os detalhes, e a incoerencia era invisivel.
//
// Nas abas Board e Arvore o clique num card abre a modal — dono, historico de posse, o
// que o card pede. No Roadmap nao abria NADA: o rotulo da esquerda nao tinha
// `data-number` nem listener, e a barra era um `<a href>` que SAIA para o GitHub.
//
// O custo: quem queria ver de quem e' um card nao tinha como, justamente na aba que
// mostra o trabalho no tempo — onde a pergunta "quem esta com isso?" mais aparece.
//
// Relatado pelo usuario: "estou clicando nos itens do board e nao esta abrindo a modal,
// eu queria ver os owners dos cards que estao na coluna in-review".
func TestRoadmapAbreOsDetalhes(t *testing.T) {
	b, err := fs.ReadFile(boardFS, "board/anchors-board.html")
	if err != nil {
		t.Fatal(err)
	}
	html := string(b)

	// O container do roadmap entra na delegacao de clique.
	if !strings.Contains(html, `"colunas", "arvore", "v-roadmap"`) {
		t.Error("o `v-roadmap` nao esta na delegacao de clique — o clique no roadmap nao " +
			"chega ao `abre()`, e a aba fica sem detalhes")
	}

	// O seletor reconhece o rotulo do roadmap.
	if !strings.Contains(html, `.g-rot[data-number]`) {
		t.Error("o seletor nao casa o rotulo do roadmap (`.g-rot[data-number]`) — o " +
			"listener existe e nao acha alvo")
	}

	// E o rotulo carrega o numero, senao nao ha o que casar.
	if !strings.Contains(html, `data-number="${esc(i.number)}"`) {
		t.Error("o rotulo do roadmap nao emite `data-number` — o seletor nao tem o que ler")
	}
}

// O BOARD DIZ POR QUEM O CARD ESPERA, e não só que espera.
//
// O `needs-user` para o card — e não diz QUEM o segura. O board mostrava "esperando você"
// sem vínculo, e quem responde a decisão não sabe o que acabou de soltar: cada card tem de
// ser reencontrado à mão, e o que não for reencontrado segue parado depois de a decisão já
// ter saído.
//
// MEDIDO no projeto de referência: 23 decisões abertas eram SETE perguntas, e destravar as
// dependentes exigia reler card por card para descobrir quem esperava o quê.
func TestBoardMostraPorQuemOCardEspera(t *testing.T) {
	html, err := BoardHTML()
	if err != nil {
		t.Fatal(err)
	}

	// O CAMPO da coleta tem de chegar ao HTML: sem ele o selo nunca aparece.
	if !strings.Contains(html, "blockedBy") {
		t.Error("o board não lê `blockedBy` — o vínculo existe na label e não é desenhado")
	}

	// O NÚMERO no selo, e não só a palavra "bloqueado": é o número que responde "esperando
	// o quê?", e sem ele o selo repete o que o `needs-user` já dizia.
	if !strings.Contains(html, "bloqueado por #") {
		t.Error("o selo de bloqueio não mostra o NÚMERO de quem segura — sem ele o leitor " +
			"sabe que o card parou e não sabe onde a resposta vai")
	}

	// O CLIQUE vai para quem segura. Sem o `data-vai-para` o clique cai no card que o
	// contém e abre o próprio bloqueado, que é o que o leitor já está vendo.
	if !strings.Contains(html, "data-vai-para") {
		t.Error("o selo não leva ao card que segura — a pergunta de quem lê `bloqueado " +
			"por #123` é o que o #123 pede")
	}

	// A COR tem de existir nas TRÊS paletas: o board é lido em tema claro, escuro por
	// `prefers-color-scheme` e escuro por `data-theme`. Uma cor definida só num bloco
	// vira `inherit` nos outros, e o selo perde justamente o que o distingue do escalado.
	if n := strings.Count(html, "--bloqueio:"); n < 3 {
		t.Errorf("a cor de bloqueio está definida %d vez(es), e as paletas são 3 "+
			"(claro, prefers-color-scheme, data-theme) — nos temas que faltam o selo "+
			"perde a cor que o distingue do card apenas escalado", n)
	}
}
