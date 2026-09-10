package scan

import "testing"

func TestRegiaoAninhadaFechaNaOrdem(t *testing.T) {
	// O caso que motivou a região: o MLETX-B05 vive DENTRO do MLETX-A03 (é o catch do
	// try dele). Sem delimitação declarada, inferir o fim por indentação dava ao A03
	// um intervalo que ia até o próximo marcador de indentação menor — e em JSX isso
	// engolia centenas de linhas (medido: 489).
	src := `func persist() {
  // #region [MLETX-A03]: confirmar persiste o lançamento.
  put(item)
    // #region [MLETX-B05]: falha ao gravar avisa o usuário.
    catch(err)
    // #endregion [MLETX-B05]
  toast()
  // #endregion [MLETX-A03]
}`
	regs, erros := Regioes(src)
	if len(erros) != 0 {
		t.Fatalf("não devia haver erro de pareamento, veio %+v", erros)
	}
	if len(regs) != 2 {
		t.Fatalf("esperava 2 regiões, veio %d", len(regs))
	}
	// a INTERNA fecha primeiro, então sai primeiro
	if regs[0].Code != "MLETX-B05" || regs[0].Start != 4 || regs[0].End != 6 {
		t.Errorf("interna errada: %+v", regs[0])
	}
	if regs[1].Code != "MLETX-A03" || regs[1].Start != 2 || regs[1].End != 8 {
		t.Errorf("externa errada: %+v", regs[1])
	}
}

func TestRegiaoJSXUmaLinhaNaoEngoleOArquivo(t *testing.T) {
	// A regressão que a heurística de indentação produzia: um marcador de UMA linha em
	// JSX, com indentação profunda, "engolia" tudo até o próximo marcador de indentação
	// menor. Com região declarada, o intervalo é o que o autor escreveu.
	src := `<View>
            {/* #region [MLETX-A01]: alternar tipo define saída/entrada. */}
            <TypeSegment />
            {/* #endregion [MLETX-A01] */}
</View>
` + longo(500)
	regs, erros := Regioes(src)
	if len(erros) != 0 || len(regs) != 1 {
		t.Fatalf("esperava 1 região sem erro, veio %d regiões %+v", len(regs), erros)
	}
	if got := regs[0].End - regs[0].Start + 1; got != 3 {
		t.Errorf("a região devia cobrir 3 linhas, cobriu %d", got)
	}
}

func TestRevDaRegiaoIgnoraMudancaForaDela(t *testing.T) {
	// O ganho central: mexer FORA da região não vence o carimbo dela. É o que separa
	// "qual arquivo mudou" de "qual requisito mudou".
	a := "// #region [CODEX-A01]: x\nfaz()\n// #endregion [CODEX-A01]\noutraCoisa()"
	b := "// #region [CODEX-A01]: x\nfaz()\n// #endregion [CODEX-A01]\noutraCoisaMUDADA()"
	ra, _ := Regioes(a)
	rb, _ := Regioes(b)
	if ra[0].Rev != rb[0].Rev {
		t.Errorf("mudança fora da região não devia mudar o rev dela: %s != %s", ra[0].Rev, rb[0].Rev)
	}
	c := "// #region [CODEX-A01]: x\nfazDIFERENTE()\n// #endregion [CODEX-A01]\noutraCoisa()"
	rc, _ := Regioes(c)
	if ra[0].Rev == rc[0].Rev {
		t.Error("mudança DENTRO da região tinha de mudar o rev")
	}
}

func TestErrosDePareamento(t *testing.T) {
	casos := []struct {
		nome, src, kind string
	}{
		{"sem fecho", "// #region [CODEX-A01]: x\nfaz()", "sem-fecho"},
		{"fecho órfão", "faz()\n// #endregion [CODEX-A01]", "fecho-orfao"},
		{"fecho trocado", "// #region [CODEX-A01]: x\nfaz()\n// #endregion [CODEX-B02]", "fecho-trocado"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			_, erros := Regioes(c.src)
			if len(erros) != 1 || erros[0].Kind != c.kind {
				t.Fatalf("esperava 1 erro %q, veio %+v", c.kind, erros)
			}
		})
	}
}

func TestFechoTrocadoNaoCascateiaNasSeguintes(t *testing.T) {
	// Um fecho trocado é UM defeito; se ele desbalanceasse a pilha, toda região seguinte
	// pareceria errada e o relatório culparia código correto.
	src := `// #region [CODEX-A01]: x
faz()
// #endregion [CODEX-B02]
// #region [CODEX-A02]: y
faz2()
// #endregion [CODEX-A02]`
	regs, erros := Regioes(src)
	if len(erros) != 1 || erros[0].Kind != "fecho-trocado" {
		t.Fatalf("esperava exatamente 1 erro fecho-trocado, veio %+v", erros)
	}
	if len(regs) != 2 {
		t.Fatalf("as 2 regiões ainda têm de sair, veio %d", len(regs))
	}
}

func TestArquivoSemRegiaoNaoEhErro(t *testing.T) {
	// A delimitação é OPCIONAL: onde não há região, vale o rev do arquivo (comportamento
	// anterior). Um projeto não precisa converter nada para continuar válido.
	regs, erros := Regioes("// CODEX-A01: comentário à moda antiga\nfaz()")
	if len(regs) != 0 || len(erros) != 0 {
		t.Errorf("sem região = sem erro e sem intervalo, veio %d/%d", len(regs), len(erros))
	}
}

func longo(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "linha\n"
	}
	return s
}

func TestComposeRefsResolveRelativoAoRoteiro(t *testing.T) {
	// Os caminhos num flow são relativos ao DIRETÓRIO dele (`../../utils/login.yaml`).
	// Resolver contra a raiz produziria alvos inexistentes, e a aresta seria descartada
	// em silêncio — o pior resultado, porque parece "sem dependência".
	src := `- runFlow: ../../utils/login.yaml
- runFlow:
    file: ../../utils/dismissOsDialogs.yaml
    env:
      SKIP_SETUP: 'true'
- runScript: ../../scripts/dataLoader.js`
	got := ComposeRefs(src, "apps/mobile/.maestro/suites/smoke/SS-03.yaml")
	esperado := []string{
		"apps/mobile/.maestro/utils/login.yaml",
		"apps/mobile/.maestro/utils/dismissOsDialogs.yaml",
	}
	if len(got) != len(esperado) {
		t.Fatalf("esperava %v, veio %v", esperado, got)
	}
	for i := range esperado {
		if got[i] != esperado[i] {
			t.Errorf("[%d] esperava %q, veio %q", i, esperado[i], got[i])
		}
	}
}

func TestComposeRefsIgnoraRunScriptEDuplicata(t *testing.T) {
	// `runScript:` aponta um .js que é DADO de entrada, não composição de roteiro.
	// E o mesmo util citado duas vezes é uma dependência, não duas.
	src := `- runScript: ../scripts/dataLoader.js
- runFlow: ../utils/login.yaml
- runFlow: ../utils/login.yaml`
	got := ComposeRefs(src, "a/b/flow.yaml")
	if len(got) != 1 || got[0] != "a/utils/login.yaml" {
		t.Fatalf("esperava só o login.yaml uma vez, veio %v", got)
	}
}

func TestComposeRefsSemComposicao(t *testing.T) {
	if got := ComposeRefs("- tapOn:\n    id: ':x'", "a/f.yaml"); got != nil {
		t.Errorf("flow sem runFlow não tem dependência de composição, veio %v", got)
	}
}
