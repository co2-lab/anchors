// @anchors
//   code: STGST
//   ref: VTRST

package gate

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// What changes a screen is proven through what the screen shows, and two gates tie each
// change of a visual unit to it:
//
//	validation-transitions   every validation leads from a state to a state — it is a
//	                         trigger of the State Flow —, or says `@no-state: <reason>`
//	error-message-declared   every error names the message it shows, or says
//	                         `@no-message: <reason>`
//
// With them, a validation is covered by the capture of the state it leads to, and an error
// by the capture of its message: the screen is one, what changes is the state or the text.
// They run, like the visual-regression gates, on a visual unit's main code file.

var (
	fromCols    = []string{"from", "de", "desde"}
	triggerCols = []string{"trigger", "gatilho", "disparador"}
	toCols      = []string{"to", "para", "a", "hacia"}
)

// visualSpec is the spec beside a visual unit's main code file, and its code.
func visualSpec(n mapx.Node, root string) (content, code string, why string, ok bool) {
	if n.Kind != mapx.KindCode {
		return "", "", "", false
	}
	b, err := readFile(root, strings.TrimSuffix(n.ID, path.Ext(n.ID))+".spec.md")
	if err != nil {
		return "", "", i18n.T("gate.vr_states.skip_no_spec"), false
	}
	m := specCodeRE().FindStringSubmatch(string(b))
	if m == nil {
		return "", "", i18n.T("gate.vr_states.skip_no_code"), false
	}
	return string(b), m[1], "", true
}

// optOut reads a `@no-<what>: <reason>` on a line: present, and whether it has a reason.
func optOut(line, what string) (present, reasoned bool) {
	m := regexp.MustCompile(`@no-` + what + `\b(?::\s*([^|]*))?`).FindStringSubmatch(line)
	if m == nil {
		return false, false
	}
	return true, strings.TrimSpace(strings.Trim(strings.TrimSpace(m[1]), "-—>")) != ""
}

// rowText is a table row's cells joined, to look for a marker in any of them.
func rowText(r map[string]string) string {
	var b strings.Builder
	for _, v := range r {
		b.WriteString(v + " | ")
	}
	return b.String()
}

// checkValidationTransitions: does every validation of a visual unit lead from a state to a
// state, or say why it changes none?
func checkValidationTransitions(_ string, n mapx.Node, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	content, code, why, ok := visualSpec(n, root)
	if !ok {
		return Skip, why
	}
	codeRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-[A-Z]\d{2}\b`)
	var validations []string
	exempt := map[string]bool{}
	var unreasoned []string
	for _, key := range []string{"section.title.validations", "section.title.presentation_validations"} {
		for _, r := range sectionRows(content, key) {
			first, _ := col(r, ruleCols...)
			v := codeRE.FindString(first)
			if v == "" {
				continue
			}
			validations = append(validations, v)
			if present, reasoned := optOut(rowText(r), "state"); present {
				if reasoned {
					exempt[v] = true
				} else {
					unreasoned = append(unreasoned, v)
				}
			}
		}
	}
	if len(validations) == 0 {
		return Skip, i18n.T("gate.validation_transitions.skip_no_validations")
	}
	states := map[string]bool{}
	for _, s := range registeredStates(content, code, stateLetter(cfg)) {
		states[code+"-"+s] = true
	}
	stateRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-[A-Z]\d{2}\b`)
	linked := map[string]bool{}
	var unknown []string
	for _, r := range sectionRows(content, "section.title.state_flow") {
		trigger, _ := col(r, triggerCols...)
		from, _ := col(r, fromCols...)
		to, _ := col(r, toCols...)
		f, t := stateRE.FindString(from), stateRE.FindString(to)
		for _, v := range codeRE.FindAllString(trigger, -1) {
			if !states[f] || !states[t] {
				unknown = append(unknown, fmt.Sprintf("%s (%s → %s)", v, orDash(f), orDash(t)))
				continue
			}
			linked[v] = true
		}
	}
	var unlinked []string
	for _, v := range dedupe(validations) {
		if !linked[v] && !exempt[v] {
			unlinked = append(unlinked, v)
		}
	}
	var gaps []string
	if len(unlinked) > 0 {
		gaps = append(gaps, i18n.T("gate.validation_transitions.no_transition", len(unlinked), strings.Join(unlinked, ", ")))
	}
	if len(unknown) > 0 {
		gaps = append(gaps, i18n.T("gate.validation_transitions.unknown_state", strings.Join(dedupe(unknown), ", ")))
	}
	if len(unreasoned) > 0 {
		gaps = append(gaps, i18n.T("gate.validation_transitions.no_reason", len(unreasoned), strings.Join(dedupe(unreasoned), ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// checkErrorMessageDeclared: does every error of a visual unit name the message it shows, or
// say why it shows none?
func checkErrorMessageDeclared(_ string, n mapx.Node, root string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	content, code, why, ok := visualSpec(n, root)
	if !ok {
		return Skip, why
	}
	codeRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-[A-Z]\d{2}\b`)
	messages := map[string]bool{}
	for _, m := range sectionCodes(content, "section.title.messages", code) {
		messages[code+"-"+m] = true
	}
	rows := sectionRows(content, "section.title.errors")
	if len(rows) == 0 {
		return Skip, i18n.T("gate.error_message.skip_no_errors")
	}
	var silent, unknown []string
	for _, r := range rows {
		first, _ := col(r, ruleCols...)
		e := codeRE.FindString(first)
		if e == "" {
			continue
		}
		if present, reasoned := optOut(rowText(r), "message"); present && reasoned {
			continue
		}
		var cited []string
		for _, c := range codeRE.FindAllString(rowText(r), -1) {
			if c != e {
				cited = append(cited, c)
			}
		}
		shown := false
		for _, c := range cited {
			if messages[c] {
				shown = true
			}
		}
		switch {
		case shown:
		case len(cited) > 0 && len(messages) > 0:
			unknown = append(unknown, fmt.Sprintf("%s (%s)", e, strings.Join(cited, ", ")))
		default:
			silent = append(silent, e)
		}
	}
	var gaps []string
	if len(silent) > 0 {
		gaps = append(gaps, i18n.T("gate.error_message.no_message", len(silent), strings.Join(silent, ", "), code))
	}
	if len(unknown) > 0 {
		gaps = append(gaps, i18n.T("gate.error_message.unknown_message", strings.Join(unknown, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
