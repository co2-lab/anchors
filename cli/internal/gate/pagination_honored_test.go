package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cfgDialeto monta a config de um projeto. Repare no que é de LINGUAGEM (vem de `family`)
// e no que é de PROVEDOR DE DADOS (`collection_query` e `cursor`, declarados aqui):
// `QueryCommand` e `LastEvaluatedKey` são nomes do SDK da AWS, não de TypeScript. Se
// estivessem embutidos no Anchors, seria lock-in de vendor disfarçado de suporte.
func cfgDialeto(family, collectionQuery, cursor string) *config.Config {
	return &config.Config{Dialect: &config.Dialect{
		Family: family, CollectionQuery: collectionQuery, Cursor: cursor,
	}}
}

// O projeto de referência dos testes: TypeScript sobre DynamoDB — o mesmo do app de referência.
func cfgTSDynamo() *config.Config {
	return cfgDialeto("ts", `(?i)(QueryCommand|ScanCommand)`, `(?i)(LastEvaluatedKey|ExclusiveStartKey)`)
}

func rodaPaginacao(t *testing.T, content string) (Verdict, string) {
	t.Helper()
	return checkPaginationHonored(content, mapx.Node{Kind: mapx.KindCode}, "", nil, cfgTSDynamo())
}

// Os quatro casos abaixo saíram do MESMO arquivo real
// (packages/backend/models/commissionLedger.ts + familyMembers.ts). O valor do gate está
// em separá-los: se acusar todos, é ruído; se acusar nenhum, é decoração.
func TestPaginationHonored(t *testing.T) {
	casos := []struct {
		nome     string
		código   string
		esperado Verdict
		contem   string
	}{
		{
			// O padrão correto: drena o cursor até o fim.
			nome: "loop de cursor drena tudo → passa",
			código: `
export async function listMembersByUser(userId: string): Promise<FamilyMember[]> {
  const items: FamilyMember[] = []
  let lastKey: Record<string, unknown> | undefined
  do {
    const result = await docClient.send(new QueryCommand({ TableName: T, ExclusiveStartKey: lastKey }))
    items.push(...(result.Items ?? []))
    lastKey = result.LastEvaluatedKey
  } while (lastKey)
  return items
}`,
			esperado: Pass,
		},
		{
			// O chamador PASSOU o limite: ele sabe que há mais. Não é promessa quebrada.
			nome: "limite escolhido pelo chamador → passa",
			código: `
export async function listCommissionsByReseller(resellerId: string, limit: number): Promise<CommissionEntry[]> {
  const res = await docClient.send(new QueryCommand({ TableName: T, Limit: limit }))
  return (res.Items ?? []) as CommissionEntry[]
}`,
			esperado: Pass,
		},
		{
			// O pior caso, e o mais comum: default que o chamador não vê. O nome promete
			// "as liberáveis"; a 101ª nunca é liberada e nada avisa.
			nome:   "default escondido → reprova nomeando o limite",
			contem: "limit = 100",
			código: `
export async function listReleasableCommissions(nowIso: string, limit = 100): Promise<CommissionEntry[]> {
  const res = await docClient.send(new QueryCommand({ TableName: T, Limit: limit }))
  return (res.Items ?? []) as CommissionEntry[]
}`,
			esperado: Fail,
		},
		{
			// Sem limite e sem paginação: trunca na página do provedor (1 MB no DynamoDB).
			nome:   "sem limite e sem cursor → reprova",
			contem: "trunca na página do provedor",
			código: `
export async function listCommissionsByOrganization(organizationId: string): Promise<CommissionEntry[]> {
  const res = await docClient.send(new QueryCommand({ TableName: T, IndexName: 'byOrganization' }))
  return (res.Items ?? []) as CommissionEntry[]
}`,
			esperado: Fail,
		},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			v, d := rodaPaginacao(t, c.código)
			if v != c.esperado {
				t.Fatalf("veredito = %s, queria %s (detalhe: %s)", v, c.esperado, d)
			}
			if c.contem != "" && !strings.Contains(d, c.contem) {
				t.Fatalf("detalhe não menciona %q: %s", c.contem, d)
			}
		})
	}
}

// O nome é o único lugar onde a promessa está escrita — `Promise<T[]>` é idêntico para
// uma página e para o conjunto. Quem nomeia o recorte não prometeu o total.
func TestPaginationNomeDelimitaAPromessa(t *testing.T) {
	corpo := `(): Promise<X[]> {
  const res = await docClient.send(new QueryCommand({ TableName: T }))
  return res.Items
}`
	casos := map[string]Verdict{
		"listRecentTransactions": Pass, // "recentes" = recorte declarado
		"listTopResellers":       Pass,
		"getFirstPageOfItems":    Pass,
		"findNextBatch":          Pass,
		"getById":                Pass, // não promete conjunto nenhum
		"listAllTransactions":    Fail, // promete o conjunto
		"getAllPending":          Fail,
		"findAllByOwner":         Fail,
	}
	for nome, esperado := range casos {
		t.Run(nome, func(t *testing.T) {
			v, d := rodaPaginacao(t, "export async function "+nome+corpo)
			if v != esperado {
				t.Fatalf("%s: veredito = %s, queria %s (%s)", nome, v, esperado, d)
			}
		})
	}
}

// A assimetria fortalece o veredito: se as irmãs paginam, o padrão é conhecido no módulo
// e a omissão é esquecimento. É o mesmo raciocínio do sibling-guard.
func TestPaginationAssimetriaEntreIrmas(t *testing.T) {
	código := `
export async function listMembersByGroup(groupId: string): Promise<M[]> {
  let lastKey
  do {
    const r = await docClient.send(new QueryCommand({ ExclusiveStartKey: lastKey }))
    lastKey = r.LastEvaluatedKey
  } while (lastKey)
  return items
}

export async function listMembersByUser(userId: string): Promise<M[]> {
  let lastKey
  do {
    const r = await docClient.send(new QueryCommand({ ExclusiveStartKey: lastKey }))
    lastKey = r.LastEvaluatedKey
  } while (lastKey)
  return items
}

export async function listActiveMemberships(userId: string): Promise<M[]> {
  const res = await docClient.send(new QueryCommand({ TableName: T }))
  return res.Items
}`
	v, d := rodaPaginacao(t, código)
	if v != Fail {
		t.Fatalf("veredito = %s, queria Fail (%s)", v, d)
	}
	if !strings.Contains(d, "listActiveMemberships") {
		t.Errorf("não nomeou a função assimétrica: %s", d)
	}
	if !strings.Contains(d, "2 funções") || !strings.Contains(d, "padrão é conhecido") {
		t.Errorf("não usou a assimetria como prova: %s", d)
	}
}

// Opt-out honesto (CONCEPT §5.1): a dispensa vale com razão escrita, e NÃO vale nua.
// Marcador sem razão é o buraco com nome bonito.
func TestPaginationDispensaExigeRazao(t *testing.T) {
	base := `
export async function listSystemFlags(): Promise<F[]> {
  // %s
  const res = await docClient.send(new QueryCommand({ TableName: T }))
  return res.Items
}`
	comRazão := strings.Replace(base, "%s", "@no-paginate: tabela de flags tem no máximo 12 linhas, fixadas em migration", 1)
	if v, d := rodaPaginacao(t, comRazão); v != Pass {
		t.Errorf("dispensa COM razão deveria passar, foi %s (%s)", v, d)
	}
	semRazão := strings.Replace(base, "%s", "@no-paginate:", 1)
	if v, _ := rodaPaginacao(t, semRazão); v != Fail {
		t.Errorf("marcador NU deveria continuar reprovando, foi %s", v)
	}
	só := strings.Replace(base, "%s", "@no-paginate", 1)
	if v, _ := rodaPaginacao(t, só); v != Fail {
		t.Errorf("marcador sem dois-pontos deveria continuar reprovando, foi %s", v)
	}
}

// O gate se cala onde não reconhece — não chuta. Um Fail falso custa mais que um Skip.
func TestPaginationSilencioOndeNaoReconhece(t *testing.T) {
	casos := map[string]string{
		"sem consulta de coleção": `
export async function listNames(): Promise<string[]> {
  return ['a', 'b', 'c']
}`,
		"dialeto não reconhecido": `
export async function listThings(): Promise<T[]> {
  return await someUnknownOrm.retrieveEverything()
}`,
	}
	for nome, código := range casos {
		t.Run(nome, func(t *testing.T) {
			if v, d := rodaPaginacao(t, código); v != Skip {
				t.Fatalf("veredito = %s, queria Skip (%s)", v, d)
			}
		})
	}
}

// Cursor presente mas FORA de laço não drena nada: ou é repassado ao chamador (e aí a
// responsabilidade é dele), ou está lá sem uso. A prova de drenagem é cursor + laço.
func TestPaginationCursorSemLacoNaoConta(t *testing.T) {
	código := `
export async function listAllItems(): Promise<I[]> {
  const res = await docClient.send(new QueryCommand({ TableName: T }))
  const nextToken = res.LastEvaluatedKey
  return res.Items
}`
	if v, d := rodaPaginacao(t, código); v != Fail {
		t.Fatalf("cursor sem laço deveria reprovar, foi %s (%s)", v, d)
	}
}

// A PROVA do agnosticismo: o mesmo defeito conceitual — nome promete o conjunto, código
// devolve a primeira página — escrito em 4 linguagens sem nada de TS/JS no gate. Se algum
// dia alguém cravar sintaxe de novo, este teste cai.
func TestPaginationAgnosticoEntreLinguagens(t *testing.T) {
	casos := []struct {
		family, query, cursor string
		defeituoso, correto   string
	}{
		{
			family: "python", query: `(?i)\.(query|scan)\(`, cursor: `(?i)next_token`,
			defeituoso: `
def list_all_orders(customer_id):
    resp = table.query(KeyConditionExpression=Key('customer').eq(customer_id))
    return resp['Items']
`,
			correto: `
def list_all_orders(customer_id):
    items, next_token = [], None
    while True:
        resp = table.query(KeyConditionExpression=Key('customer').eq(customer_id), next_token=next_token)
        items.extend(resp['Items'])
        next_token = resp.get('next_token')
        if not next_token:
            break
    return items
`,
		},
		{
			family: "go", query: `(?i)\.Query\(`, cursor: `(?i)nextToken`,
			defeituoso: `
func ListAllOrders(ctx context.Context, customer string) ([]Order, error) {
	out, err := db.Query(ctx, &QueryInput{Customer: customer})
	return out.Items, err
}
`,
			correto: `
func ListAllOrders(ctx context.Context, customer string) ([]Order, error) {
	var all []Order
	var nextToken string
	for {
		out, err := db.Query(ctx, &QueryInput{Customer: customer, NextToken: nextToken})
		if err != nil {
			return nil, err
		}
		all = append(all, out.Items...)
		nextToken = out.NextToken
		if nextToken == "" {
			break
		}
	}
	return all, nil
}
`,
		},
		{
			family: "java", query: `(?i)\.(query|createQuery)\(`, cursor: `(?i)nextToken`,
			defeituoso: `
public List<Order> listAllOrders(String customerId) {
    QueryResponse resp = client.query(QueryRequest.builder().customer(customerId).build());
    return resp.items();
}
`,
			correto: `
public List<Order> listAllOrders(String customerId) {
    List<Order> all = new ArrayList<>();
    String nextToken = null;
    do {
        QueryResponse resp = client.query(QueryRequest.builder().customer(customerId).nextToken(nextToken).build());
        all.addAll(resp.items());
        nextToken = resp.nextToken();
    } while (nextToken != null);
    return all;
}
`,
		},
		{
			family: "rust", query: `(?i)\.query\(`, cursor: `(?i)next_token`,
			defeituoso: `
pub async fn list_all_orders(customer: &str) -> Vec<Order> {
    let resp = client.query().customer(customer).send().await.unwrap();
    resp.items
}
`,
			correto: `
pub async fn list_all_orders(customer: &str) -> Vec<Order> {
    let mut all = Vec::new();
    let mut next_token = None;
    loop {
        let resp = client.query().customer(customer).next_token(next_token).send().await.unwrap();
        all.extend(resp.items);
        next_token = resp.next_token;
        if next_token.is_none() { break }
    }
    all
}
`,
		},
	}

	for _, c := range casos {
		t.Run(c.family+" — sem paginar reprova", func(t *testing.T) {
			v, d := checkPaginationHonored(c.defeituoso, mapx.Node{Kind: mapx.KindCode}, "", nil,
				cfgDialeto(c.family, c.query, c.cursor))
			if v != Fail {
				t.Fatalf("veredito = %s, queria Fail (%s)", v, d)
			}
		})
		t.Run(c.family+" — paginando passa", func(t *testing.T) {
			v, d := checkPaginationHonored(c.correto, mapx.Node{Kind: mapx.KindCode}, "", nil,
				cfgDialeto(c.family, c.query, c.cursor))
			if v != Pass {
				t.Fatalf("veredito = %s, queria Pass (%s)", v, d)
			}
		})
	}
}

// Sem dialeto declarado, o gate NÃO pode passar: ele não olhou o código. Pendente com o
// nome do campo que falta — a correção é uma linha de YAML, não uma investigação.
func TestPaginationSemDialetoEhPendente(t *testing.T) {
	código := `
export async function listAllThings(): Promise<T[]> {
  return (await docClient.send(new QueryCommand({}))).Items
}`
	casos := map[string]*config.Config{
		"config nula":                  nil,
		"config sem dialect":           {},
		"dialect sem exported_func":    {Dialect: &config.Dialect{CollectionQuery: `QueryCommand`}},
		"dialect sem collection_query": {Dialect: &config.Dialect{Family: "ts"}},
	}
	for nome, cfg := range casos {
		t.Run(nome, func(t *testing.T) {
			v, d := checkPaginationHonored(código, mapx.Node{Kind: mapx.KindCode}, "", nil, cfg)
			if v != Pending {
				t.Fatalf("veredito = %s, queria Pending — um ✓ aqui seria mentira (%s)", v, d)
			}
			if !strings.Contains(d, "dialect") {
				t.Errorf("Pendente não diz o que declarar: %s", d)
			}
		})
	}
}

func TestPaginationOptOutSaiDoRelatorio(t *testing.T) {
	// "Não declarei ainda" e "não se aplica a mim" são estados diferentes com o mesmo
	// sintoma. Um projeto sem banco não tem `collection_query` para declarar, e cobrar dele
	// eternamente transforma o relatório em ruído que se aprende a ignorar — levando embora
	// os avisos que importam.
	//
	// Sem opt-out: Pending (o gate não verificou, e diz).
	semOptOut := &config.Config{Dialect: &config.Dialect{ExportedFunc: `(?m)^export function (\w+)`}}
	n := mapx.Node{ID: "x.ts", Kind: mapx.KindCode}
	if v, _ := checkPaginationHonored("export function listAll() {}", n, "", nil, semOptOut); v != Pending {
		t.Fatalf("sem collection_query nem opt-out, esperava Pending, veio %v", v)
	}
	// Com opt-out: Skip — alguém decidiu, e a decisão está escrita no anchors.yaml.
	comOptOut := &config.Config{Dialect: &config.Dialect{
		ExportedFunc: `(?m)^export function (\w+)`,
		OptOut:       []string{"collection_query"},
	}}
	v, msg := checkPaginationHonored("export function listAll() {}", n, "", nil, comOptOut)
	if v != Skip {
		t.Fatalf("com opt-out, esperava Skip, veio %v — %s", v, msg)
	}
	if !strings.Contains(msg, "opt_out") {
		t.Errorf("a mensagem deve dizer que foi dispensado por opt_out, veio: %s", msg)
	}
}

func TestPaginationMensagemOfereceOOptOut(t *testing.T) {
	// A mensagem tem de mostrar a saída. Sem ela, a única opção visível é conviver com o
	// aviso para sempre — e aviso permanente é aviso ignorado.
	_, msg := checkPaginationHonored("x", mapx.Node{ID: "x.ts", Kind: mapx.KindCode}, "", nil,
		&config.Config{})
	if !strings.Contains(msg, "opt_out") {
		t.Errorf("o Pending deve oferecer o opt-out, veio: %s", msg)
	}
}

func TestPaginationPegaOCasoRealDoSDK(t *testing.T) {
	// O código REALX que motivou a investigação (packages/backend/infra/cognito.ts, antes da
	// correção): `Limit` fixo em 20 na chamada e o `res.PaginationToken` descartado no
	// `return`. Um usuário com 21 aparelhos nunca via o 21º, e nenhum teste pegou — o teste
	// roda com 2 aparelhos, e com 2 o limite de 20 é irrelevante.
	//
	// O gate pega, mas pela regra do LIMITE, não pela do cursor. A distinção importa e está
	// documentada no gate: descartar um campo é NÃO mencioná-lo, e um gate que lê texto não
	// detecta a ausência de algo que nunca foi escrito. A mensagem é menos precisa (cita o
	// limite) e aponta a função certa, com "devolva o cursor" entre as saídas.
	src := `export async function cognitoListDevices(
  accessToken: string,
): Promise<DeviceType[]> {
  const res = await client.send(new ListDevicesCommand({ AccessToken: accessToken, Limit: 20 }))
  return res.Devices ?? []
}`
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts",
		CollectionQuery: `(?i)(ListDevicesCommand|QueryCommand|ScanCommand)`}}
	v, msg := checkPaginationHonored(src, mapx.Node{ID: "infra/cognito.ts", Kind: mapx.KindCode},
		"", nil, cfg)
	if v != Fail {
		t.Fatalf("teto fixo de 20 sem cursor deve reprovar, veio %v — %s", v, msg)
	}
	if !strings.Contains(msg, "cognitoListDevices") {
		t.Errorf("a mensagem deve nomear a função, veio: %s", msg)
	}
}

func TestPaginationPrefixoDeProvedorNaoEscondeAPromessa(t *testing.T) {
	// `cognitoListDevices` promete conjunto tanto quanto `listDevices`. O `setPromise` estava
	// ancorado em `^`, então o prefixo do provedor ESCONDIA a promessa e o gate nem olhava a
	// função — foi o que fez este caso passar despercebido.
	cfgTS := &config.Config{Dialect: &config.Dialect{Family: "ts"}}
	d := cfgTS.DialectFor()
	for _, nome := range []string{"cognitoListDevices", "s3ListObjects", "listDevices"} {
		if !prometeConjunto(nome, d) {
			t.Errorf("%s promete conjunto e não foi reconhecido", nome)
		}
	}
	// E o verbo tem de ser palavra: `all` interno não conta, senão vira falso positivo.
	for _, nome := range []string{"allocateSlot", "callbackUrl", "getUserById"} {
		if prometeConjunto(nome, d) {
			t.Errorf("%s NÃO promete conjunto — falso positivo", nome)
		}
	}
}

func TestPaginationSemCursorNoProvedorNaoInventa(t *testing.T) {
	// Se o provedor não expõe cursor, não há o que repassar. O gate não pode exigir o que
	// a API não oferece — inventar aqui produziria achado que ninguém consegue consertar.
	src := `export async function listTudo(): Promise<Item[]> {
  const res = await client.send(new QueryCommand({ TableName: 't' }))
  return res.Items ?? []
}`
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts",
		CollectionQuery: `(?i)QueryCommand`, Cursor: `(?i)PaginationToken`}}
	v, msg := checkPaginationHonored(src, mapx.Node{ID: "m.ts", Kind: mapx.KindCode}, "", nil, cfg)
	// Pode reprovar por OUTRO motivo (consulta sem paginar), mas não pela regra do cursor.
	if strings.Contains(msg, "não o devolve") {
		t.Errorf("sem cursor no provedor, a regra do cursor não se aplica; veredito %v — %s", v, msg)
	}
}
