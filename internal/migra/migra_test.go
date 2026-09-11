package migra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func escreve(t *testing.T, dir, nome, conteudo string) string {
	t.Helper()
	p := filepath.Join(dir, nome)
	if err := os.WriteFile(p, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// O CASO QUE MOTIVA TUDO: `julgamentos` virou `judgments`, e perder a chave não daria erro
// — daria SILÊNCIO. Os carimbos de julgamento de IA sumiriam e o `check` refaria todos,
// cobrando de novo o que alguém já respondeu.
func TestMigraArquivo_renomeiaAsChavesDoMapa(t *testing.T) {
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", `# cabeçalho
version: 1
gerado_por: 0.1.83
nodes:
    - id: a.spec.md
      kind: spec
      code_declarado: true
edges:
    - from: a.spec.md
      to: a.ts
      julgamentos:
        - gate: doc-self-contained
          verdict: pass
`)
	r, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Changed {
		t.Fatal("o arquivo deveria ter sido alterado")
	}

	got, _ := os.ReadFile(p)
	texto := string(got)
	for _, novo := range []string{"generated_by:", "code_declared:", "judgments:", "version: 2"} {
		if !strings.Contains(texto, novo) {
			t.Errorf("esperava %q no arquivo migrado", novo)
		}
	}
	for _, velho := range []string{"gerado_por:", "code_declarado:", "julgamentos:"} {
		if strings.Contains(texto, velho) {
			t.Errorf("a chave antiga %q continua no arquivo", velho)
		}
	}
	// O VALOR tem de sobreviver: uma migração que renomeia e perde o conteúdo é pior que
	// nenhuma, porque parece ter funcionado.
	if !strings.Contains(texto, "0.1.83") || !strings.Contains(texto, "doc-self-contained") {
		t.Error("a migração perdeu valor ao renomear a chave")
	}
}

// A ancoragem da chave é o que separa migração de reescrita cega: `julgamentos` aparece em
// COMENTÁRIO e em prosa dentro dos arquivos do Anchors, e trocar por substring
// corromperia texto que não é chave nenhuma.
func TestMigraArquivo_naoTocaOQueNaoEhChave(t *testing.T) {
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", `# os julgamentos ficam aqui; gerado_por diz quem gravou
version: 1
nodes:
    - id: a.spec.md
      # este nó tem julgamentos pendentes
      title: "sobre julgamentos e code_declarado"
`)
	if _, err := MigrateFile(p, 2, false); err != nil {
		t.Fatal(err)
	}
	texto := string(mustRead(t, p))
	if !strings.Contains(texto, "# os julgamentos ficam aqui; gerado_por diz quem gravou") {
		t.Error("o comentário foi reescrito — a troca não está ancorada na chave")
	}
	if !strings.Contains(texto, `"sobre julgamentos e code_declarado"`) {
		t.Error("um VALOR com o texto da chave foi reescrito")
	}
}

// `trinca_opcional` vive no `anchors.yaml`, e `julgamentos` no mapa. Trocar a chave no
// arquivo errado reescreveria algo que por acaso tem o mesmo nome.
func TestMigraArquivo_cadaChaveNoSeuArquivo(t *testing.T) {
	dir := t.TempDir()
	cfg := escreve(t, dir, "anchors.yaml", `version: 1
layers:
    spec:
        trinca_opcional: [covered-by, tested-by]
`)
	if _, err := MigrateFile(cfg, 2, false); err != nil {
		t.Fatal(err)
	}
	texto := string(mustRead(t, cfg))
	if !strings.Contains(texto, "triad_optional:") {
		t.Error("o `trinca_opcional` do anchors.yaml não foi migrado")
	}
	// E o VALOR da dispensa, que é o que não pode sumir: sem ele o gate `trinca-completa`
	// volta a cobrar o que o time decidiu não ter.
	if !strings.Contains(texto, "[covered-by, tested-by]") {
		t.Error("a migração perdeu a dispensa declarada")
	}

	// A OUTRA METADE da regra: uma chave DO MAPA não pode ser trocada no `anchors.yaml`.
	//
	// A primeira versão deste teste só verificava que a chave certa migrou, e sobreviveu à
	// mutação que remove o filtro por arquivo — porque um `anchors.yaml` comum não tem
	// `julgamentos` para ser trocado por engano. Aqui ele tem, de propósito: é um nome de
	// GATE, e o gate se chama assim em qualquer arquivo.
	cfg2 := escreve(t, dir, "anchors.yaml", `version: 1
gates:
    - id: julgamentos
      when: [pre-commit]
layers:
    spec:
        gerado_por: quem-sabe
`)
	if _, err := MigrateFile(cfg2, 2, false); err != nil {
		t.Fatal(err)
	}
	t2 := string(mustRead(t, cfg2))
	if strings.Contains(t2, "judgments") {
		t.Error("uma chave DO MAPA foi trocada no anchors.yaml")
	}
	if strings.Contains(t2, "generated_by") {
		t.Error("o `gerado_por` é chave do MAPA e não deveria ser tocado aqui")
	}
}

// Migrar duas vezes não muda nada na segunda: o `version:` sobe MESMO sem chave trocada, e
// sem isso a migração rodaria de novo a cada comando.
func TestMigraArquivo_ehIdempotente(t *testing.T) {
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", "version: 1\nnodes: []\nedges: []\n")

	r1, err := MigrateFile(p, 2, false)
	if err != nil || !r1.Changed {
		t.Fatalf("a primeira migração deveria alterar; err=%v", err)
	}
	primeira := mustRead(t, p)

	r2, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if r2.Changed {
		t.Error("a segunda migração não deveria alterar nada")
	}
	if string(mustRead(t, p)) != string(primeira) {
		t.Error("migrar de novo mudou o arquivo")
	}
}

// `dryRun` MEDE sem escrever — é o que permite o `doctor` avisar sem mexer no repositório
// de alguém que não pediu.
func TestMigraArquivo_dryRunNaoEscreve(t *testing.T) {
	dir := t.TempDir()
	original := "version: 1\ngerado_por: dev\nnodes: []\n"
	p := escreve(t, dir, "anchors.graph.yaml", original)

	r, err := MigrateFile(p, 2, true)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Changed {
		t.Error("o dry-run deveria RELATAR que há o que mudar")
	}
	if r.Replaced["gerado_por"] != 1 {
		t.Errorf("o dry-run deveria contar a chave a trocar; veio %v", r.Replaced)
	}
	if string(mustRead(t, p)) != original {
		t.Error("o dry-run ESCREVEU no arquivo")
	}
}

// Um mapa de antes de o campo `version:` existir: ele é formato 1, e migrá-lo exige
// inserir a linha — não só trocá-la.
func TestMigraArquivo_mapaSemVersionGanhaOCampo(t *testing.T) {
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", "# anchors.graph.yaml — o mapa\n# derivado\nnodes: []\n")

	if de, _ := FormatOf(p); de != 1 {
		t.Fatalf("sem `version:` o formato deveria ser 1, e veio %d", de)
	}
	if _, err := MigrateFile(p, 2, false); err != nil {
		t.Fatal(err)
	}
	texto := string(mustRead(t, p))
	if !strings.Contains(texto, "version: 2") {
		t.Error("o campo `version:` não foi inserido")
	}
	// Depois do cabeçalho, não antes: o `Save` escreve os comentários primeiro, e um
	// `version:` acima deles sairia do lugar na próxima gravação.
	if strings.HasPrefix(texto, "version:") {
		t.Error("o `version:` entrou ANTES do cabeçalho de comentário")
	}
}

func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
