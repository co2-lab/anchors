---
title: "Traceability & Identity Codes"
description: "How identity codes like AUTH-B01 connect specs, features, tests, and code into an unbroken chain of custody."
---

How do you prove that line 142 of `invoice_service.go` exists to satisfy Requirement 3.2 of the billing specification?

In traditional software development, this connection exists only in git commit messages, PR descriptions, or developers' heads. Over time, as code is refactored, the link vanishes.

Anchors establishes **deterministic traceability** through **Identity Codes**.

---

## 1. The Anatomy of an Identity Code

Every business rule, scenario, and contract in Anchors carries a standardized identity code:

```
                  AUTH - B 0 1
                   │    │  │
                   │    │  └── Sequence number (01, 02, ...)
                   │    └───── Rule type letter:
                   │           B = Business Rule
                   │           V = Validation Rule
                   │           S = Security Rule
                   │           F = Failure / Error Handling
                   │           P = Performance / Scale
                   └────────── Module / Domain prefix (AUTH, BILL, USER)
```

Examples:
- `AUTH-B01`: Authentication Business Rule #1 (e.g. Issue JWT on valid credentials).
- `BILL-V02`: Billing Validation Rule #2 (e.g. Invoice amount must be greater than zero).
- `PAY-F01`: Payment Failure Rule #1 (e.g. Handle payment gateway timeout with exponential backoff).

---

## 2. Connecting the Four Pieces with Identity Codes

The identity code appears across all four artifacts of [The Unit](/docs/concepts/unit/):

### 1. In the Spec (`*.spec.md`):
```markdown
### AUTH-B01 — Issue Session JWT
Upon successful credential validation, the system MUST issue a cryptographically
signed JWT with an expiration of 15 minutes.
```

### 2. In the Feature File (`*.feature`):
```gherkin
@rule:AUTH-B01
Scenario: User logs in with valid credentials
  Given a registered user with email "alice@example.com"
  When they submit valid credentials
  Then the system responds with a signed JWT token valid for 15 minutes
```

### 3. In the Test Code (`*_test.go`):
```go
// @test:AUTH-B01
func TestIssueSessionJWT_ValidCredentials(t *testing.T) {
    token, err := service.Authenticate("alice@example.com", "secret")
    assert.NoError(t, err)
    assert.True(t, token.ExpiresInMinutes(15))
}
```

### 4. In the Source Code (`*.go`):
```go
// @rule:AUTH-B01
func (s *AuthService) IssueToken(user *User) (string, error) {
    return s.jwtSigner.Sign(user.Claims(), 15 * time.Minute)
}
```

---

## 3. How Gates Verify Traceability

When you run Anchors, traceability gates inspect this chain:
- [`spec-feature-match`](/docs/gates/spec-feature-match/): Ensures code `AUTH-B01` defined in the spec appears in the `.feature` file.
- [`feature-test-match`](/docs/gates/feature-test-match/): Ensures scenario `@rule:AUTH-B01` is implemented in test code.
- [`code-cataloged`](/docs/gates/code-cataloged/): Ensures the function tagged `@rule:AUTH-B01` exists in the codebase.

If someone deletes the test or changes the rule without updating the other artifacts, the build fails immediately.

---

## 4. Related Concepts

- [The Unit](/docs/concepts/unit/): The bundle united by identity codes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): How gates validate codes automatically.
