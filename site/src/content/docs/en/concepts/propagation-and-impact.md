---
title: "Propagation & Impact"
description: "How changes in specs and code propagate across the dependency graph and trigger the revision wave."
---

In complex software systems, no change exists in a vacuum. Changing the schema of a user profile might break an authentication token generator, which breaks an API response, which breaks three mobile app screens.

Traditional development relies on humans remembering what to test. Anchors computes the exact **Propagation Wave**.

---

## 1. What is the Propagation Wave?

When an anchor or specification is modified, Anchors treats that change as an **origin point** in the project graph and traverses downstream edges to identify all impacted dependents:

```
  [ Modified Spec ]
         │
         ├──► [ Dependent Spec A ] ──► [ Test Suite A ] (Stale)
         │
         └──► [ Dependent Spec B ] ──► [ Contract B ] (Stale)
                                              │
                                              └──► [ Client SDK ] (Impacted)
```

The set of all nodes affected by a change is called the **Wave**.

---

## 2. Stale vs. Fresh Evidence

Anchors tracks hash signatures for every artifact. When an upstream spec is modified:
- All downstream test results are marked **Stale**.
- The [`evidence-fresh`](/docs/gates/evidence-fresh/) gate blocks releases until fresh test runs confirm that dependents still function correctly.
- Downstream specs enter a **Revision Needed** state.

---

## 3. Computing the Wave with the CLI

You can calculate the exact impact of any proposed change before committing:

```bash
# Calculate downstream impact of modifying a spec
anchors wave src/services/auth/session.spec.md

# Output:
# Wave Origin: src/services/auth/session.spec.md
# Direct Dependents:
#   - src/services/billing/checkout.spec.md (via @uses: AUTH-B01)
#   - src/api/handlers/auth_handler.go
# Downstream Impact:
#   - tests/integration/checkout_flow_test.go [STALE]
# Recommended Action: Re-run integration suite and verify checkout contract.
```

---

## 4. Key Gates Governing Propagation

- [`evidence-fresh`](/docs/gates/evidence-fresh/): Ensures no test evidence is older than the spec it verifies.
- [`contract-impact`](/docs/gates/contract-impact/): Verifies that breaking changes to API contracts trigger explicit revision waves.
- [`updated-at-atual`](/docs/gates/updated-at-atual/): Ensures timestamps and revision numbers reflect current modifications.
