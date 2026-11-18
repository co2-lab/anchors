---
title: "Feature Flags"
description: "How Anchors governs feature toggles as first-class architectural gates with dual-state scenario validation."
---

Feature flags are an essential tool for continuous delivery and dark launching. However, in most codebases, feature flags become tech debt:
- Developers forget to remove flags after a feature is fully launched.
- Tests only verify the flag enabled state, leaving the disabled state broken.
- Dead branches linger in the codebase forever.

Anchors treats **Feature Flags** as first-class architectural entities.

---

## 1. Declaring Flags in Anchors

Feature flags are declared in dedicated `.flag.md` files:

```markdown
---
flag: NEW_CHECKOUT_FLOW
status: active
default: false
sunset_milestone: v2.4.0
---

# NEW_CHECKOUT_FLOW

Enables the redesigned one-step checkout wizard with instant payment processing.
```

---

## 2. Dual-State Scenario Enforcement

When a feature is governed by a flag, Anchors requires tests for **both states**:
1. Behavior when the flag is **ON**.
2. Fallback behavior when the flag is **OFF**.

In your Gherkin `.feature` file:

```gherkin
@flag:NEW_CHECKOUT_FLOW=enabled
Scenario: User checks out with one-step wizard
  Given the feature flag "NEW_CHECKOUT_FLOW" is enabled
  When the user clicks "Instant Checkout"
  Then the one-step modal is displayed

@flag:NEW_CHECKOUT_FLOW=disabled
Scenario: User checks out with legacy flow
  Given the feature flag "NEW_CHECKOUT_FLOW" is disabled
  When the user clicks "Checkout"
  Then the standard multi-step form is displayed
```

---

## 3. Key Feature Flag Gates

- [`flag-covered`](/docs/gates/flag-covered/): Ensures all flagged code paths have corresponding tests for both ON and OFF states.
- [`flag-scenario-exists`](/docs/gates/flag-scenario-exists/): Verifies that flag scenarios exist in feature files.
- [`flag-scenarios-complete`](/docs/gates/flag-scenarios-complete/): Prevents partial test coverage of flagged features.
- [`flag-scenario-grammar`](/docs/gates/flag-scenario-grammar/): Enforces standard `@flag:` tag syntax.
