package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// obligation-honored: confronta as OBRIGAÇÕES TRANSVERSAIS declaradas (`obligations:`
// no anchors.yaml) contra a realidade. Um nó cujo header carrega o atributo-gatilho
// DEVE aparecer nos arquivos que a obrigação exige.
//
// É a classe de erro que nenhuma spec-por-unidade pega, porque o dever mora FORA da
// unidade. O caso que motivou (real): um modelo de dado novo, com anotações de texto
// livre do usuário, ficou de fora do script de exclusão de conta — dado pessoal que
// nunca seria apagado. O projeto já tinha sido mordido por isso antes, e o próprio
// script carrega um comentário chamando de "violação LGPD silenciosa". Repetiu, e
// nenhum dos 9 gates existentes viu, porque todos olham para dentro da unidade.
//
// EXCEÇÃO HONESTA: um nó pode se eximir com `obligation_waived: <nome> — <motivo>` no
// header. O motivo é OBRIGATÓRIO (waiver sem justificativa é tratado como ausente).
// É a mesma doutrina de @noPropagation: o silêncio não é permitido, mas a exceção
// legítima fica registrada e localizada em vez de o time desligar o gate.
func checkObligationHonored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	obrigacoes := allObligations(root, cfg)
	if len(obrigacoes) == 0 {
		return Skip, "" // projeto sem obrigações declaradas — nada a confrontar
	}

	var violated, pending []string
	for _, ob := range obrigacoes {
		if ob.When == "" || !headerHasAttr(content, ob.When) {
			continue // o nó não dispara esta obrigação
		}
		if reason := waiverFor(content, ob.Name); reason != "" {
			continue // exceção declarada COM motivo
		}
		name := nodeIdentifier(n.ID)
		// O nó pode declarar COMO é referenciado no destino, quando a derivação não dá
		// conta. É comum: o nome do artefato e o da referência seguem convenções
		// diferentes, e às vezes INCONSISTENTES entre si (num projeto real,
		// `Transaction`→`TRANSACTIONS_TABLE_NAME` mas `AuditItem`→`AUDIT_ITEM_TABLE_NAME`
		// — plural em um, singular no outro, sem regra derivável). Adivinhar pluralização
		// no engine seria pôr a bagunça de um projeto dentro do framework; declarar é
		// honesto e local.
		// PRECEDÊNCIA do token — a ordem importa e foi medida:
		//
		//  1. `identified_as_form` da OBRIGAÇÃO, quando declarado. Existe porque deveres
		//     diferentes procuram FORMAS diferentes do mesmo artefato: o de LGPD procura a
		//     env var (`METADATA_ENTRY_TABLE_NAME`) nos handlers; o de provisionamento
		//     procura o nome do modelo (`tables['MetadataEntry']`) na infra.
		//  2. `identified_as` do NÓ. É a única fonte que conhece a irregularidade real do
		//     projeto — `BalancePoint` vira `BALANCE_POINTS_TABLE_NAME` (plural) e
		//     `AuditItem` vira `AUDIT_ITEM_TABLE_NAME` (singular), sem regra derivável.
		//     Inverter esta ordem acusa 28 modelos que estão corretos (medido).
		//  3. `identified_by` da obrigação (forma automática), e por fim o nome cru.
		var token string
		switch {
		case ob.IdentifiedAsForm != "":
			token = applyIdentifierForm(name, ob.IdentifiedAsForm)
		case headerAttr(content, "identified_as") != "":
			token = headerAttr(content, "identified_as")
		case ob.IdentifiedBy != "":
			token = applyIdentifierForm(name, ob.IdentifiedBy)
		default:
			token = name
		}
		var faltando []string
		for _, glob := range ob.MustAppearIn {
			ok, checked := tokenAppearsIn(root, glob, token)
			if !checked {
				continue // glob não casou arquivo nenhum — não inventa violação
			}
			if !ok {
				faltando = append(faltando, glob)
			}
		}
		if len(faltando) > 0 {
			sort.Strings(faltando)
			msg := fmt.Sprintf("`%s` (como `%s`) não aparece em: %s",
				name, token, strings.Join(faltando, ", "))
			if ob.Because != "" {
				msg += " — " + ob.Because
			}
			// DEVER RECONHECIDO: o terceiro estado, entre cumprir e dispensar.
			//
			// Sem ele o gate só oferecia duas saídas, e nenhuma servia ao caso mais comum:
			// a obrigação é REAL e será cumprida noutra fase (o handler que deve consumir
			// a tabela ainda não existe). Dispensar seria mentira — o dever não deixou de
			// existir. Deixar vermelho confunde a dívida assumida com o esquecimento, que
			// é exatamente a distinção que o pilar existe para preservar.
			//
			// Quem declara `obligation_pending` afirma três coisas: que conhece o dever,
			// que ele continua valendo, e QUANDO será pago. É registro, não dispensa —
			// por isso o veredito é Pendente (aparece no relatório), nunca Pass.
			if quando := pendingFor(content, ob.Name); quando != "" {
				pending = append(pending, fmt.Sprintf("[%s] %s — DÍVIDA ASSUMIDA: %s",
					ob.Name, msg, quando))
				continue
			}
			violated = append(violated, fmt.Sprintf("[%s] %s", ob.Name, msg))
		}
	}
	if len(violated) == 0 && len(pending) == 0 {
		return Pass, ""
	}
	if len(violated) == 0 {
		sort.Strings(pending)
		return Pending, "obrigação transversal com DÍVIDA ASSUMIDA (não é falha — é registro): " +
			strings.Join(pending, "; ")
	}
	sort.Strings(violated)
	msg := "obrigação transversal não cumprida: " + strings.Join(violated, "; ") +
		". Três saídas, e só uma é silêncio:\n" +
		"  1. CUMPRA o dever;\n" +
		"  2. se ele será pago depois, assuma a dívida com " +
		"`obligation_pending: <nome> — <quando/onde>` (fica Pendente no relatório, visível);\n" +
		"  3. se o dever NÃO se aplica a este nó, dispense com " +
		"`obligation_waived: <nome> — <motivo>`.\n" +
		"Dispensar o que ainda se deve é mentira; deixar vermelho o que já foi reconhecido " +
		"confunde dívida com esquecimento."
	if len(pending) > 0 {
		msg += fmt.Sprintf(" (há também %d dívida(s) assumida(s) neste nó)", len(pending))
	}
	return Fail, msg
}

// pendingFor devolve o QUANDO/ONDE de uma dívida assumida, ou "" se não há.
// Como o waiver, exige a justificativa escrita: `obligation_pending: <nome> — <quando>`.
// Um marcador nu não assume dívida nenhuma — só esconde melhor.
func pendingFor(content, obligation string) string {
	return declaracaoComMotivo(content, "obligation_pending", obligation)
}

// headerHasAttr diz se o header `@anchors` declara o atributo-gatilho (ex.:
// "carries: pii"). Compara normalizando espaço, para tolerar a escrita à mão.
func headerHasAttr(content, attr string) bool {
	k, v, ok := strings.Cut(attr, ":")
	if !ok {
		return false
	}
	re := regexp.MustCompile(`(?mi)^\s*(?://|#|<!--|\*)?\s*` +
		regexp.QuoteMeta(strings.TrimSpace(k)) + `:\s*` +
		regexp.QuoteMeta(strings.TrimSpace(v)) + `\s*$`)
	return re.MatchString(headerOf(content))
}

// waiverFor devolve o MOTIVO do waiver desta obrigação, ou "" se não há waiver válido.
// Waiver sem motivo não vale — é o que separa a exceção honesta do silêncio.
func waiverFor(content, obligation string) string {
	// O separador é o TRAVESSÃO (—) ou ` - ` COM espaços. Aceitar `-` colado permitiria
	// que o próprio nome da obrigação fornecesse o separador: em `pii-purgavel`, o `-`
	// interno faria o regex ler nome=`pii` e motivo=`purgavel`, auto-eximindo qualquer
	// obrigação com hífen no nome.
	return declaracaoComMotivo(content, "obligation_waived", obligation)
}

// declaracaoComMotivo lê `<campo>: <obrigação> — <motivo>` do header e devolve o motivo.
// Compartilhado por `obligation_waived` (dispensa) e `obligation_pending` (dívida
// assumida): as duas são declarações que só valem COM justificativa escrita.
func declaracaoComMotivo(content, campo, obligation string) string {
	re := regexp.MustCompile(`(?mi)^\s*(?://|#|<!--|\*)?\s*` + regexp.QuoteMeta(campo) + `:\s*` +
		regexp.QuoteMeta(obligation) + `\s*(?:—|\s-\s)\s*(\S.*?)\s*$`)
	if m := re.FindStringSubmatch(headerOf(content)); m != nil {
		return m[1]
	}
	return ""
}

// headerOf recorta o início do arquivo (onde vive o bloco @anchors). Limitar evita que
// uma menção no corpo do documento seja lida como declaração de header.
func headerOf(content string) string {
	if len(content) > 1200 {
		return content[:1200]
	}
	return content
}

// nodeIdentifier é o nome pelo qual o nó é conhecido: o basename sem as extensões
// compostas (`.spec.md`, `.test.ts`) nem a simples.
func nodeIdentifier(id string) string {
	base := filepath.Base(id)
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx"} {
		if strings.HasSuffix(base, suf) {
			return strings.TrimSuffix(base, suf)
		}
	}
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// applyIdentifierForm converte o nome do nó para a forma como ele é REFERENCIADO no
// destino. O nome do artefato e a referência a ele raramente coincidem — um modelo
// `MetadataEntry` aparece no script de purge como `METADATA_ENTRIES_TABLE_NAME`.
func applyIdentifierForm(name, form string) string {
	switch {
	case form == "" || form == "as-is":
		return name
	case form == "screaming-snake":
		return strings.ToUpper(splitCamel(name, "_"))
	case form == "snake":
		return strings.ToLower(splitCamel(name, "_"))
	case form == "kebab":
		return strings.ToLower(splitCamel(name, "-"))
	case strings.Contains(form, "{{NAME}}") || strings.Contains(form, "{{SCREAMING}}"):
		// template livre: {{NAME}} = nome, {{SCREAMING}} = SCREAMING_SNAKE
		r := strings.NewReplacer(
			"{{NAME}}", name,
			"{{SCREAMING}}", strings.ToUpper(splitCamel(name, "_")),
		)
		return r.Replace(form)
	}
	return name
}

var camelBoundaryRE = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func splitCamel(s, sep string) string {
	return camelBoundaryRE.ReplaceAllString(s, "${1}"+sep+"${2}")
}

// tokenAppearsIn procura o token nos arquivos que casam o glob. Devolve (achou,
// houve-arquivo-para-checar) — sem arquivo, não há violação a declarar.
func tokenAppearsIn(root, glob, token string) (bool, bool) {
	matches, err := doublestar.Glob(os.DirFS(root), glob)
	if err != nil || len(matches) == 0 {
		return false, false
	}
	// Fronteira própria: `\b` falha aqui porque `_` é word-char — `\bMETADATA_ENTRY\b`
	// NÃO casa dentro de `METADATA_ENTRY_TABLE_NAME`, que é justamente a forma como o
	// destino referencia o nó. Exigimos que o token não esteja colado a outro
	// identificador à esquerda; à direita, `_` é continuação legítima.
	re := regexp.MustCompile(`(^|[^A-Za-z0-9_$])` + regexp.QuoteMeta(token) + `([^A-Za-z0-9$]|$)`)
	for _, m := range matches {
		b, err := os.ReadFile(filepath.Join(root, m))
		if err != nil {
			continue
		}
		if re.Match(b) {
			return true, true
		}
	}
	return false, true
}

// headerAttr lê o valor de um atributo livre do header `@anchors` (ex.:
// `identified_as: TRANSACTIONS_TABLE_NAME`). Vazio se não declarado.
func headerAttr(content, key string) string {
	re := regexp.MustCompile(`(?mi)^\s*(?://|#|<!--|\*)?\s*` + regexp.QuoteMeta(key) + `:\s*(\S+)\s*$`)
	if m := re.FindStringSubmatch(headerOf(content)); m != nil {
		return m[1]
	}
	return ""
}
