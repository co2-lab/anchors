// @anchors
//   code: GDMNG
//   ref: GVGDG

package governance

// monitorGuide is the habit of waiting on a long process: launch it in the background, watch
// `anchors monitor`, react to each line (DESIGN-process-monitor.md).
const monitorGuide = `# Monitor guide (waiting on a long process without waiting for nothing)

## The problem

A test suite, a mutation run, a build or the commit hook can take many minutes. A watch that
speaks only when the process ends misses the costly cases: a process that died, or hung,
never ends — and the agent waits for nothing while the user stares at a screen that says
nothing. Measured in a project: agents sat idle from 40 minutes to 4 hours before anyone
noticed.

## The habit

1. Launch the long command in the background.
2. Watch 'anchors monitor' with your watch tool — one notification per line.
3. React to each line as it comes; do not wait for the end to look.
4. When the watch expires, start it again: the monitor resumes where it stopped, says what
   happened meanwhile, and says nothing twice.

'anchors monitor --until-done' exits when the runs it watches have all ended, with their
outcome as the last line — for a watch that should end with the work.

## The lines, and what each asks of you

- '▶ <runner> — <who>: <command>' — a run started. Note it.
- '⏸ <runner>: no output or CPU for <time>' — it hangs. Look at its output; kill it and run it
  again, or find what it waits for.
- '▶ <runner> moving again' — it was only slow.
- '☠ <runner> gone without a summary' — it died. Read its last lines, then run it again; do
  not report a result it never produced.
- '✓/✗ <runner> ended' — with its exit code or its summary. Go on, or fix what failed.
- '✓/✗ report <path>' — a test report landed, with its failures: triage them while the rest
  runs.
- '<time> running=[...] load=<n>' — the heartbeat. 'running=[NOTHING]' while work is pending
  means a step never started: start it. A load marked HIGH explains timeouts — do not chase
  them as bugs before it falls.

## What the monitor sees

- Anchors' own long commands (test, mutation, check, verify — the commit hook —, map build,
  docs build, ingest) record their runs, with their exit.
- The agent's long and background commands, when the project installed the agent hook
  ('anchors install-hooks --agent'): their command line and the file their output goes to.
- 'anchors monitor run -- <command>' records anything else, with its output and its exit.
- Anything else a runner recognizes under the project — jest, vitest, go test, pytest,
  maestro, playwright, stryker —, from the process table. How such a process ended is known
  only from its output's summary.

The user sees the same in 'anchors status', under "Running".

## Timing

The defaults — read every 30s, heartbeat every 5m, a run stalled after its runner's quiet
time — come from the code; a 'monitor:' block in anchors.yaml changes them, and the flags
('--every', '--heartbeat', '--progress', '--stall') override both. A project declares its own
runners there too: the pattern of the command line, its quiet time, and the patterns of its
output that say it passed, failed or advanced.
`
