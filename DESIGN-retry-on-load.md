<!-- @anchors
  code: RTLOD
  layer: desenho
-->

# Retry on load — a test that timed out under load runs once more, and says so

> APPROVED on 2026-10-13 (the user, on MIF's report). Decisions D01–D03 below.
>
> IMPLEMENTED on 2026-10-13. The fixtures are MIF's three failures:
> - the hook's `Exceeded timeout` is told by its message;
> - the `waitFor` that rethrows `Expected: true / Received: false` is told by its line;
> - an assertion failing outside a wait is not told.
>
> Under a played load past 2x, the failed file alone ran again, the run passed, and the test was
> recorded as passed on retry, whole and with `--changed`. Under the threshold, or for an
> assertion, the failure stood. `tests-retried` lists a file's retried tests of the last 7 days
> and calls a test retried 3 times unstable. Its window is days, not runs: the record keeps
> retries, not every run.

## The problem

Measured at MIF on 2026-10-13. A Maestro round, the agents' chain and tests in worktrees ran at
once, with a load of 20 to 77. Three suites of the mobile app failed on time alone:

- two `waitFor` of 1 s ran out (useIRSnapshot UIRNU-B01#01, useBalancePoints USBPU-B01#01);
- an `afterEach` of the testing library ran past 20 s (ImportStatusBanner IPSBI-V06#01).

Run alone, they passed: 14/14 and 25/25. Each time the check failed `tests-pass` and
`scenario-coverage`, and the agent ran the file again by hand and redid the check.

A runner's own retry (`jest.retryTimes`) repeats any failure, a wrong assertion included, and
hides that it repeated. Only Anchors knows the load at the moment of the failure and writes
the evidence.

## The design

- **One retry, only of a test that failed on time, only under load.** A failed test is run
  again once when both hold:
  - its failure reads as a timeout. The pattern is matched against the failure's message **and**
    the line of the test where it failed, which Anchors reads at the stack's file:line. Measured
    at MIF: a `waitFor` that runs out does not say "timeout". It rethrows its last assertion
    (`Expected: true / Received: false`), and only the line, `await waitFor(...)`, tells the
    wait. An assertion that fails inside a wait is, by definition, a wait that ran out.
    - The built-in patterns cover the runners Anchors knows, e.g. jest's `Exceeded timeout`,
      `waitFor(` and `findBy`.
    - A suite can set its own `timeout_pattern`.
    - An assertion that fails anywhere else is never repeated;
  - the load during the run passed the threshold. `anchors test` samples the load while the
    run lasts and keeps the highest, against the CPUs.
- **The retry runs only those tests, with less contention.** It uses the suite's
  `retry_run:`, or its `run_changed:` when it has none, with the failed files in `{{files}}`.
  For example: `run_changed: "jest --runInBand {{files}}"`.
- **The evidence is honest.**
  - A test that passes on the retry is passed, recorded as `passed on retry (load 31/10,
    timeout)`.
  - The run says how many tests passed that way.
  - An informative gate, `tests-retried`, lists the tests that needed a retry, and how often
    in their recent runs. A test that times out every time stays visible as unstable instead
    of disappearing behind the retry.
  - A test that fails again is a real failure.
- **Configured per suite, in anchors.yaml:**

      tests:
        - layer: unit
          run: "pnpm jest"
          run_changed: "pnpm jest --findRelatedTests {{files}}"
          retry_on_load:
            load: 2x          # the run's highest load against the CPUs
            max: 1
            run: "pnpm jest --runInBand {{files}}"   # optional

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `RTLOD-W01` | **The retry and its evidence.** The highest load of a run, the timeout patterns over the message and the failing line, the retry of the failed files, the evidence `passed on retry`, the count. | With MIF's three failures as fixtures — the hook's `Exceeded timeout`, and the two `waitFor` that rethrow their assertion — under load they are run again and pass with the mark; the same failures without load, and an assertion failing outside a wait under load, are not repeated. |
| `RTLOD-W02` | **`tests-retried`.** The informative gate over the retried tests and how often. | A test retried 3 times in 7 days is listed as unstable; a retry older, or of another file, does not count. |

## What NOT to do

- **Repeat a failed assertion.** Only a timeout, and only under load.
- **Hide the retry.** A pass on retry is a pass with a mark, never a plain pass.
- **Retry without limit.** One retry by default; a test that needs more is unstable, and says so.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `RTLOD-D01` | On or off by default? | **Decided: off** (the user, 2026-10-13). A suite turns it on with `retry_on_load:`. A project that never runs under load keeps every failure as it is. |
| `RTLOD-D02` | The load threshold. | **Decided: 2× the CPUs** (the user, 2026-10-13), over the highest load the run saw, sampled by `anchors test` while it runs. The monitor's heartbeat read about 13 on 10 CPUs near MIF's `waitFor` failures. That is one sample every 5 minutes, and the run's own peak may be higher. |
| `RTLOD-D03` | Does the check accept a test that passed on retry? | **Decided: yes** (the user, 2026-10-13), as passed, with the mark; `tests-retried` informs. |
