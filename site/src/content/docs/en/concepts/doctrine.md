---
title: "Product Doctrine"
description: "Understand how Anchors enforces cross-cutting business rules and non-functional requirements across all units."
---

Every company has core business principles that apply everywhere. For example:
- *"All monetary amounts must be stored in cents as integers, never floats."*
- *"Personal identifying information (PII) must never appear in application logs."*
- *"Every user action that modifies account balance must generate an immutable audit record."*

When developers build new features, they often forget these cross-cutting rules. In Anchors, we govern these invariants through **Product Doctrine**.

---

## 1. What is Product Doctrine?

**Product Doctrine** consists of non-negotiable, overarching policies stored in dedicated doctrine specifications:

```
docs/doctrine/
├── financial_integrity.doctrine.md
├── privacy_and_pii.doctrine.md
└── audit_logging.doctrine.md
```

Unlike regular unit specs that govern a single function or file, doctrine specs govern **entire domains or the whole application**.

---

## 2. Realizing Doctrine in Feature Specs

When a feature specification operates within a governed domain, it must explicitly **realize** relevant doctrine rules:

```markdown
# Checkout Service Spec

@realizes: FIN-DOC-01 (Store currency in cents)
@realizes: AUDIT-DOC-03 (Emit audit event on state transition)

### CHECKOUT-B01 — Calculate Final Total
Calculates cart subtotal, applies discount coupon, and emits AUDIT_CHECKOUT_COMPLETED.
```

If a spec modifies financial values but does not declare how it realizes `FIN-DOC-01`, the doctrine gates catch the omission.

---

## 3. Key Doctrine Gates

- [`spec-doctrine-exists`](/docs/gates/spec-doctrine-exists/): Verifies that all declared doctrine files exist and are valid.
- [`spec-realizes-doctrine`](/docs/gates/spec-realizes-doctrine/): Ensures feature specs implement required cross-cutting doctrine.
- [`doctrine-realized`](/docs/gates/doctrine-realized/): Verifies that doctrine rules are realized in actual source code and tests.
- [`doctrine-not-duplicated`](/docs/gates/doctrine-not-duplicated/): Prevents teams from writing conflicting local policies.
