// @anchors
//   code: GNCGD
//   ref: GVGDG

package governance

// navigationGuide is the ruler for an app's navigation: the screen spec's In and Out
// tables, the flag on each navigation call, and the gates that hold the three together.
// Its examples splice the at sign (flagAt), or this file would carry them as real flags.
const navigationGuide = `# Navigation guide (the screens, where each leads, and by which rule)

## The problem

An app's navigation lives in three places that drift apart without a word: the screen's
spec says where it leads, the code calls navigate, and the screen it arrives at says where it
comes from. A screen the code stopped reaching, a back navigation nobody wrote down, a spec
promising a route the code never takes — each is invisible until someone walks the app.
Anchors ties the three with the same rule the dependency chain follows: the comment flag is
the declaration; the code is the proof.

## The screen's spec: In and Out

A screen is a spec of the 'screen' layer. Its Navigation section has two tables:

    ## Navigation

    ### In

    | Origin | Action |
    | --- | --- |
    | HomeScreen | tap a goal |

    ### Out

    | Rule | Destination | Action |
    | --- | --- | --- |
    | 'GLDTG-A01' | GoalEditScreen | edit |

- Out: one row per navigation the screen makes, with the rule that triggers it — a code of
  this spec — and the screen it leads to, by name or route.
- In: one row per screen that leads here.
- The headings and columns are read in English, Portuguese or Spanish (Entrada/Saída,
  Origem/Destino, Regra).

## The code: a flag on every call

Every navigation call carries its flag on its line, or the line before — navigate, push,
replace, goBack, reset, popToTop: a back navigation is a navigation too.

    navigation.navigate('GoalEdit') // ` + flagAt + `navigates: GLETG [GLDTG-A01]
    // ` + flagAt + `navigates: HOMEH, GOALG
    navigation.goBack()
    close() // ` + flagAt + `no-nav: closes a modal of this same screen

- '@navigates: <CODE>[, <CODE>...] [<RULE>]' — the screen (or screens, for a back navigation
  that returns to more than one) this call leads to, by its code, and the Out rule it
  answers.
- '` + flagAt + `no-nav: <reason>' — this call is no screen edge.
- 'anchors check --fix' writes the flag on every call whose route names one screen; a back
  navigation, whose target the code does not say, is the author's to write.

## The gates

- 'nav-annotated' — every navigation call carries a flag or a waiver, and a flag on a call
  whose route names a screen names that screen;
- 'nav-matches-spec' — every Out row is answered by a flag in the screen's code, and every
  flag by an Out row;
- 'nav-symmetric' — every Out of one screen is an In of the other, and every In an Out;
- 'nav-reachable' — every screen is reached from an entry route along the edges.

## The project declares

    dialect:
      navigation_call: '<pattern with a (?P<route>…) group>'   # the family's when omitted
    navigation:
      entry: [Home]                                            # the roots of reachability

With no entry declared, 'nav-reachable' is pending: there is nothing to reach from.

## Seeing it

- 'anchors map nav [screen]' prints every edge, or one screen's.
- '{{ with navigation }}' in a template under 'doct/' draws the map as a Mermaid flowchart —
  an entry as a stadium, a screen no entry reaches dashed, each edge labelled with its rule;
  'anchors docs init' seeds 'navigation.md.tmpl' when the project has screens.
`
