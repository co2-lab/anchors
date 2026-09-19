---
title: Freezing the project
description: The panic button — how to stop all work when a problem must be solved first.
---

Sometimes a problem appears and **nobody should work until it's solved**: the
plan points at a spec that doesn't exist, an architecture decision turned out
wrong, a credential leaked.

Without a brake, the team keeps producing against a base that will change — and
hours of work age before they can be delivered.

```sh
anchors freeze --motivo "plan 0002 points at a spec that does not exist — see #42"
```

To release:

```sh
anchors thaw
```

## Four layers, and none suffices alone

`freeze` engages all four in one operation. Separating them produces inconsistent
state — an active ruleset with no issue explaining it, or an `enabled: false`
nobody pushed.

| layer | what it stops | who can bypass |
| --- | --- | --- |
| `enabled: false` in `anchors.yaml` | commit and push on the machine of whoever already cloned | `--no-verify` |
| **the CLI** | every command that produces state | — |
| **the ruleset** on the remote | push and merge on the server | admin |
| **the claim** | assignment of new cards | — |

### Why the platform brake isn't enough

A ruleset blocking push and merge doesn't reach the machine of whoever **already
cloned**. There, `check`, `judge` and `ingest` keep running and writing to the
map.

Freezing means **stopping the production of state**, not just its delivery.

### Why the local hooks aren't enough

A hook is bypassable (`--no-verify`), and whoever clones after the freeze doesn't
pass through it until installing the hooks. That's why it's the **second** layer
— the first is the ruleset, which nobody bypasses from the inside.

The hooks exist so people **find out early**, not to be inviolable.

## The reason is mandatory

```sh
anchors freeze --motivo "..."     # without this, the command refuses
```

Not bureaucracy. It's the text **every refusal** will show — in the blocked
commit, the blocked push, the refused claim, the issue.

A freeze with no written reason is indistinguishable from broken configuration,
and whoever hits it tries to work around it instead of reading.

## Whoever fixes it gets through

This is deliberate, and it's written in every refusal message:

- `--no-verify` bypasses the hooks
- the **admin** bypasses the ruleset

A brake that prevents its own fix becomes the problem. The freeze exists to stop
work **by inertia**, not the correction that releases it.

## What keeps working while frozen

The rule: what runs frozen is what **doesn't produce project state**.

```sh
anchors status      # where the project is
anchors doctor      # the ecosystem X-ray
anchors guide ...   # the guides
anchors coverage    # what's already been measured
anchors impact      # what a change would reach
anchors thaw        # otherwise the freeze would be irreversible
```

Whoever is investigating needs those. What's refused is what **writes**: `check`,
`map build`, `judge`, `ingest`, `new`, `escalate`.

## What freeze does, step by step

1. writes `enabled: false` and the reason into `anchors.yaml`
2. **commits and pushes** — with `--no-verify`, because the hooks it just
   activated would refuse the freeze itself
3. creates the `anchors-freeze` ruleset on the remote, blocking push and merge on
   every branch
4. opens the `[congelado]` issue with the reason

Step 2 uses `--no-verify` for one more reason: `pre-commit` runs the gates, and
an urgent freeze **cannot depend on the suite being green** — the reason for the
freeze may be precisely that it isn't.

### The options

```sh
anchors freeze --motivo "..." --sem-ruleset   # local brake only
anchors freeze --motivo "..." --sem-push      # doesn't commit or push
```

## After the thaw

`thaw` returns `anchors.yaml` to its exact prior state — comments included.

One thing to know: the hooks cache the answer for up to **10 minutes**, to avoid
paying a `git fetch` on every commit. Anyone with a warm cache may take that long
to notice the release, or force it with a `git fetch`.

The cache is **asymmetric** on purpose: it caches "not frozen", but when it says
"frozen" it consults the network again. The reason is the cost of the error — a
cache that held the freeze would leave someone blocked with the project already
released, not understanding why and with nothing they could do.
