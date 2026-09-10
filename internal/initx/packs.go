package initx

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// packsFS carrega os packs de conformidade que o `anchors init` semeia no projeto.
//
// Eles VIAJAM no binário e são COPIADOS para `packs/` — não são lidos de dentro dele. A
// diferença importa: no projeto, o pack é material, versionado e AUDITÁVEL. Quem precisa
// justificar conformidade abre o arquivo e lê o dever ao lado do artigo que o origina; se
// vivesse só no binário, a resposta seria "confie na versão do CLI".
//
// E permite o que a lei exige na prática: adaptar. Uma norma muda, um projeto tem uma
// interpretação jurídica própria, um mercado tem exigência adicional — o pack no projeto
// se edita. O pack no binário obrigaria a esperar um release.
//
//go:embed packs
var packsFS embed.FS

// SeedPacks copia os packs embutidos para `packs/` na raiz do projeto.
//
// Copiar TODOS (e não só os da jurisdição declarada) é deliberado: o projeto que amanhã
// abrir um mercado novo encontra o pack pronto e declara uma linha, em vez de descobrir
// que precisa escrever LGPD do zero. O que filtra é a adoção (`packs:` no anchors.yaml),
// não a presença no disco.
//
// NÃO sobrescreve pack que já existe: o projeto pode tê-lo adaptado, e sobrescrever a
// adaptação de alguém é o tipo de perda silenciosa que este framework existe para evitar.
func SeedPacks(root string) (criados, preservados []string, err error) {
	err = fs.WalkDir(packsFS, "packs", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, rerr := packsFS.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		dest := filepath.Join(root, p) // `packs/privacy/lgpd.yaml`
		if _, serr := os.Stat(dest); serr == nil {
			preservados = append(preservados, p)
			return nil
		}
		if merr := os.MkdirAll(filepath.Dir(dest), 0o755); merr != nil {
			return merr
		}
		if werr := os.WriteFile(dest, b, 0o644); werr != nil {
			return werr
		}
		criados = append(criados, p)
		return nil
	})
	sort.Strings(criados)
	sort.Strings(preservados)
	return criados, preservados, err
}

// AvailablePacks lista os packs embutidos, agrupados por domínio — para o `init`
// perguntar e para o `anchors compliance` dizer o que existe mas não foi adotado.
func AvailablePacks() map[string][]string {
	out := map[string][]string{}
	_ = fs.WalkDir(packsFS, "packs", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".yaml") {
			return err
		}
		rel := strings.TrimPrefix(p, "packs/")
		dom := filepath.Dir(rel)
		name := strings.TrimSuffix(filepath.Base(rel), ".yaml")
		out[dom] = append(out[dom], dom+"/"+name)
		_ = fmt.Sprint(name)
		return nil
	})
	for k := range out {
		sort.Strings(out[k])
	}
	return out
}
