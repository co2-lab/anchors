// Package migra leva um projeto do formato ANTIGO do mapa para o atual.
//
// POR QUE EXISTE. As chaves de YAML do produto estão em inglês — `lang` traduz o que se
// LÊ, e uma chave é identificador, não prosa —, e quatro nasceram em português:
// `gerado_por`, `code_declarado`, `julgamentos` e `trinca_opcional`.
//
// Renomeá-las sem migração quebraria todo projeto existente, e o modo de quebrar seria o
// pior possível: SILENCIOSO. Um mapa com `julgamentos` perderia os carimbos de julgamento
// de IA, e o `check` refaria todos — cobrando de novo o que alguém já respondeu. Um
// `anchors.yaml` com `trinca_opcional` perderia a dispensa declarada, e o gate
// `trinca-completa` voltaria a cobrar o que o time decidiu não ter.
//
// A ALTERNATIVA A MIGRAR seria ler os dois nomes para sempre. Ela não fecha: o arquivo
// nunca se conserta, cada leitor precisa saber dos dois, e a terceira renomeação (quando
// houver) herda o problema das duas anteriores.
//
// A migração roda uma vez e deixa o projeto no formato novo.
package migra

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// versionRE lê o `version:` de topo sem carregar o arquivo pelo parser de YAML.
//
// Não usa o parser de propósito: o arquivo pode estar num formato que o parser atual
// recusa — que é exatamente o caso que se quer detectar e consertar.
var versionRE = regexp.MustCompile(`(?m)^version:[[:space:]]*([0-9]+)[[:space:]]*$`)

// FormatOf lê a versão de formato declarada no arquivo.
//
// Devolve 1 quando não há `version:`: os primeiros arquivos do produto não o declaravam, e
// tratá-los como formato 1 é o que permite migrá-los em vez de recusá-los.
func FormatOf(path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	m := versionRE.FindSubmatch(b)
	if m == nil {
		return 1, nil
	}
	var v int
	if _, err := fmt.Sscanf(string(m[1]), "%d", &v); err != nil {
		return 0, fmt.Errorf("%s: `version:` não é um número", path)
	}
	return v, nil
}

// keyRE casa uma chave YAML no começo da linha, em qualquer indentação, item de lista
// incluído — e só a chave.
//
// A ancoragem é o que torna a troca segura: `julgamentos` aparece em comentário e em prosa
// dentro dos arquivos do Anchors, e uma substituição por substring reescreveria texto que
// não é chave nenhuma.
func keyRE(nome string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^([[:space:]]*(?:-[[:space:]]+)?)` + regexp.QuoteMeta(nome) + `:`)
}

// valueRE casa `chave: valor` e devolve o valor, ancorado no começo da linha.
//
// A ancoragem é o que separa migração de reescrita cega: `gate` e `id` aparecem em
// comentário e em prosa, e uma substituição por substring corromperia texto que não é
// chave nenhuma.
func valueRE(chave string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^([[:space:]]*(?:-[[:space:]]+)?` +
		regexp.QuoteMeta(chave) + `:[[:space:]]*)([^[:space:]#]+)[[:space:]]*$`)
}

// Result descreve o que a migração fez — ou faria — num arquivo.
type Result struct {
	Path     string
	De, To   int
	Replaced map[string]int // chave antiga → quantas ocorrências
	Changed  bool
}

// MigrateFile leva um arquivo do formato antigo ao `destino`.
//
// `dryRun` mede sem escrever: é o que o `doctor` usa para AVISAR sem mexer no repositório
// de alguém que não pediu.
func MigrateFile(path string, destino int, dryRun bool) (*Result, error) {
	de, err := FormatOf(path)
	if err != nil {
		return nil, err
	}
	r := &Result{Path: path, De: de, To: destino, Replaced: map[string]int{}}
	if de >= destino {
		return r, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	texto := string(b)

	base := filepath.Base(path)
	passos, err := StepsFrom(de, destino)
	if err != nil {
		return nil, err
	}

	// OS PASSOS EM ORDEM, cada um sobre o resultado do anterior. Um projeto no formato 1
	// com o binário no 4 atravessa 1→2, 2→3, 3→4 — e nunca 1→4 direto, que exigiria um
	// passo que conhecesse os três saltos.
	for _, passo := range passos {
		for de, para := range passo.RenameKeys[base] {
			re := keyRE(de)
			if n := len(re.FindAllString(texto, -1)); n > 0 {
				r.Replaced[de] = n
				texto = re.ReplaceAllString(texto, "${1}"+para+":")
			}
		}
		for chave, valores := range passo.RenameValues[base] {
			re := valueRE(chave)
			texto = re.ReplaceAllStringFunc(texto, func(linha string) string {
				m := re.FindStringSubmatch(linha)
				if m == nil {
					return linha
				}
				novo, ok := valores[m[2]]
				if !ok {
					return linha
				}
				r.Replaced[chave+": "+m[2]]++
				return m[1] + novo
			})
		}
	}

	// O `version:` sobe MESMO SEM nada trocado: o formato é do arquivo, não do que ele por
	// acaso contém. Um mapa sem julgamento nenhum ainda é formato 2 depois de migrado, e
	// deixá-lo em 1 faria a migração rodar de novo a cada comando.
	if versionRE.MatchString(texto) {
		texto = versionRE.ReplaceAllString(texto, fmt.Sprintf("version: %d", destino))
	} else if base == "anchors.graph.yaml" {
		// O mapa SEM `version:` é de antes de o campo existir. Ele entra logo depois do
		// cabeçalho de comentário, que é onde o `Save` o escreve — acima dele, sairia do
		// lugar na próxima gravação.
		linhas := strings.SplitAfter(texto, "\n")
		i := 0
		for i < len(linhas) && strings.HasPrefix(linhas[i], "#") {
			i++
		}
		linhas = append(linhas[:i], append([]string{fmt.Sprintf("version: %d\n", destino)}, linhas[i:]...)...)
		texto = strings.Join(linhas, "")
	}

	r.Changed = texto != string(b)
	if !r.Changed || dryRun {
		return r, nil
	}

	modo := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		modo = info.Mode()
	}
	return r, os.WriteFile(path, []byte(texto), modo)
}
