<!-- @anchors
  code: TIMNG
  updated_at: 2026-09-26
-->
# Flag: timing-metrics

> **Code**: `TIMNG`

Whether `anchors check` measures and prints how long each gate took.

The flag is real and this file is the axis confronting itself: `--timing` was added to
find what makes a scan expensive, and it found it — `docs-fresh` was 97% of a 6m49s run.

## Scenarios

| Scenario | When the value | Then |
| --- | --- | --- |
| `TIMNG-G01` | `= "off"` | the check prints only the verdicts, and measures no time <!-- @no-govern: honored by reportTiming in cmd/anchors/quality/check.go:1110 and proven in check_timing_test.go; the CLI layer has no spec citing it with @gated-by yet --> |
| `TIMNG-G02` | `= "on"` | the check also prints time per gate, and the slowest targets <!-- @no-govern: honored by reportTiming in cmd/anchors/quality/check.go:1110 and proven in check_timing_test.go; the CLI layer has no spec citing it with @gated-by yet --> |
| `TIMNG-G03` | `absent` | the same as `off` — measuring is opt-in, never a default cost <!-- @no-govern: default defined in cobra flag registration and proven in check_timing_test.go; the CLI layer has no spec citing it with @gated-by yet --> |
