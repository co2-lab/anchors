// @anchors
//   ref: GVGDG

package governance

// reportBugGuide teaches an agent to tell a bug in Anchors from a problem of the project,
// to report the first upstream, and what to do while the fix does not come.
const reportBugGuide = `# Report-bug guide (when Anchors itself is wrong)

You are using Anchors in a project, and sometimes what is wrong is not the project: it is
Anchors. Reported where Anchors is fixed — github.com/co2-lab/anchors —, the fix reaches
every project that uses it, yours included. Kept to yourself, it is worked around here and
met again in the next project. Reporting it is part of the work, not a favour.

## Is it Anchors, or the project?

Signs it is ANCHORS:

- a command crashes, panics, or exits with an error that names no file of yours;
- a gate's verdict contradicts its own documentation (` + "`anchors doctor`" + ` says what each gate
  measures; the site has a page per gate) or what the file plainly holds;
- the same input gives different verdicts on two runs, or a command says it did something
  it did not do (a ` + "`map build`" + ` that changes nothing it should have);
- a message names a command, a flag or a file that does not exist;
- a file Anchors seeds (a workflow, a guide, the header) is wrong as seeded.

Signs it is the PROJECT — not a bug, do not report it:

- a gate fails and the file really does not meet what the gate measures: fix the file;
- the ` + "`anchors.yaml`" + ` declares something that does not match the tree (a layer pattern,
  a gate's ` + "`on:`" + `): ` + "`anchors doctor`" + ` names most of these;
- a gate is not what this project wants: that is configuration (` + "`blocking`" + `, ` + "`severity`" + `,
  ` + "`skip_on`" + `, removing it), or a feature request, not a bug.

In doubt, reproduce it in a scratch directory with a made-up minimal ` + "`anchors.yaml`" + `: if it
reproduces there, it is Anchors.

## How to report

1. Reduce it to a MINIMAL CASE, made up: the smallest anchors.yaml, files and commands that
   show it. Never copy the project's code, names, paths or data — the repository is PUBLIC.
2. Write it in Anchors' terms: the command or gate, what it did, what it should do.
3. Look at it before sending:

       anchors report-bug "<what happened>" --expected "<what should happen>" --repro <case> --dry-run

4. Send it — the same command without ` + "`--dry-run`" + `. The version and the platform are filled
   in. An open issue with the same title gets a "seen again" comment instead of a
   duplicate, so report even when you suspect someone already did: the count of sightings
   is how a bug is prioritised.

The command refuses a text that carries the project's path, your home folder or the
project's repository. That guard catches the obvious, not everything: the rest is yours.

If ` + "`gh`" + ` cannot create it (no login, no network), the command prints a prefilled link: hand
it to the user to open.

In a project whose queue is on GitHub, ` + "`anchors escalate --bug --upstream \"<reason>\"`" + `
does both at once — the bug card in the project, and the report to Anchors.

## While the fix does not come

- Do not edit Anchors' installed files, the map or a seeded workflow by hand to get past
  it: the next ` + "`map build`" + ` or upgrade undoes it, and the report loses its case.
- If the bug blocks your work, waive the one gate on the one target, with the issue as
  the reason (` + "`[skip-<gate>@<CODE>: Anchors issue <url>]`" + ` in the commit message), and
  tell the user. A waiver that names its reason is lifted when the fix ships.
- Tell the user what you reported and its address, in one line.
- When a new Anchors release fixes it, remove the waiver.
`
