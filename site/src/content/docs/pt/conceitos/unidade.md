---
title: A Unidade (The Unit)
description: "Entenda como a Spec, a Feature, o Teste, o Código e a Documentação formam a unidade indivisível de governança do Anchors."
---

Se você perguntar a dez desenvolvedores como garantir que uma funcionalidade foi entregue com qualidade, a maioria responderá: *"Basta ter 100% de cobertura de testes"*.

No Anchors, sabemos que isso é uma armadilha perigosa. Você pode ter 100% de cobertura de testes rodando em um código que faz a coisa errada, ou testes vazios que não assustam nenhum bug.

Para resolver isso de forma definitiva, toda [Camada Regida](/pt/docs/layers/) no Anchors opera sobre o conceito de **A Unidade (The Unit)**.

---

## 1. O que é a Unidade?

A **Unidade** é o bloco fundamental e indivisível de desenvolvimento do Anchors. Em vez de espalhar regras no Notion, cenários na cabeça do QA, testes em pastas distantes e código solto, uma unidade é um **conjunto coeso de peças que realizam uma funcionalidade**:

```
                      ┌─────────────────────────┐
                      │          SPEC           │
                      │    (A Regra Escrita)    │
                      └────────────┬────────────┘
                                   │
                                   ▼
                      ┌─────────────────────────┐
                      │         FEATURE         │
                      │  (O Cenário em Gherkin) │
                      └────────────┬────────────┘
                                   │
              ┌────────────────────┼────────────────────┐
              ▼                    ▼                    ▼
 ┌─────────────────────────┐ ┌───────────┐ ┌─────────────────────────┐
 │          TESTE          │ │  CÓDIGO   │ │     DOCS / VISUAL       │
 │   (A Prova Executável)  │ │ (Lógica)  │ │ (Baselines / Manuais)   │
 └─────────────────────────┘ └───────────┘ └─────────────────────────┘
```

| Peça | O que ela faz | Formato | Pergunta que responde |
| --- | --- | --- | --- |
| **1. Spec** | Define as regras de negócio em formato estruturado. | `*.spec.md` | *"Qual é o comportamento esperado e as restrições?"* |
| **2. Feature** | Traduz as regras em cenários legíveis e executáveis. | `*.feature` | *"Quais passos o usuário/sistema executa para validar a regra?"* |
| **3. Teste** | Codifica os cenários em código de teste automatizado. | `*_test.go`, `*.test.ts` | *"O sistema realmente cumpre o cenário quando executado?"* |
| **4. Código** | Implementa as funções e tipos que executam a lógica. | `*.go`, `*.ts`, `*.py` | *"Como o sistema realiza o que foi pedido?"* |
| **5. Docs / Assets** | Documentação agregada ou baselines de tela (quando aplicável). | `*.doc.md`, `*.png` | *"Como isso é consumido externamente ou como a tela se parece?"* |

> [!NOTE]
> **Nota histórica sobre o nome "Trinca":** Nas primeiras versões do Anchors, este conceito era informalmente chamado informalmente de *"A Trinca"* (código, feature e teste) que orbitavam ao redor da spec. O gate automatizado se chama [`unit-complete`](/pt/docs/gates//unit-complete/), mas o conceito oficial e natural é **A Unidade**, pois reúne todas as peças indispensáveis da entrega.

---

## 2. Por que todas as peças precisam existir?

Imagine o que acontece quando qualquer uma das peças é deixada de fora:

- **Código + Teste (sem Spec e sem Feature)**: Você tem código que funciona hoje, mas ninguém sabe *por que* ele funciona assim. Quando outro desenvolvedor ou uma IA for alterar o código daqui a três meses, não saberá quais comportamentos eram intencionais e quais eram bugs acidentais.
- **Spec + Código (sem Teste)**: É apenas uma promessa no papel. Ninguém prova que o código realmente cumpre o que a spec prometeu.
- **Spec sozinha (sem Código e sem Teste)**: É um plano fantasma. Pior ainda: sem o gate [`unit-complete`](/pt/docs/gates//unit-complete/), uma spec sozinha passaria em branco pelo CI como se estivesse tudo bem.

A Unidade transforma uma funcionalidade em um **fato auditável e inegociável**.

---

## 3. Como as peças conversam: O Código de Identidade

As peças da unidade não se conectam simplesmente por terem nomes de arquivos parecidos. Elas se conectam através de um [Código de Identidade de Cenário](/pt/docs/concepts/rastreabilidade-e-codigos/) compartilhado (ex: `AUTH-B01`).

Veja um exemplo real:

### 1. Na Spec (`Login.spec.md`)
```markdown
### AUTH-B01 — Bloqueio após três tentativas incorretas
Se o usuário errar a senha por 3 vezes consecutivas, a conta deve ser
bloqueada temporariamente por 15 minutos.
```

### 2. Na Feature (`Login.feature`)
```gherkin
@AUTH-B01
Cenário: Bloqueio de conta após três erros
  Dado que o usuário "maria@exemplo.com" errou a senha 2 vezes
  Quando ela tenta fazer login com a senha errada pela 3ª vez
  Então a conta deve ficar no status "bloqueada" por 15 minutos
```

### 3. No Teste (`Login_test.go`)
```go
func TestLogin_AUTH_B01_BloqueioAposTresErros(t *testing.T) {
    // Prova executável do cenário AUTH-B01
    // ...
}
```

### 4. No Código (`Login.go`)
```go
// #region AUTH-B01
if user.FailedAttempts >= 3 {
    user.LockUntil = time.Now().Add(15 * time.Minute)
    return ErrAccountLocked
}
// #endregion AUTH-B01
```

---

## 4. Os Gates que Protegem a Unidade

O Anchors possui um conjunto de [gates](/pt/docs/concepts/gates-e-vereditos/) especializados em garantir que a unidade esteja sempre íntegra:

- [`unit-complete`](/pt/docs/gates//unit-complete/): Impede que uma spec de camada regida exista sem suas peças correspondentes.
- [`spec-feature-match`](/pt/docs/gates//spec-feature-match/): Garante que **todas** as regras catalogadas na spec possuam pelo menos um cenário na feature.
- [`feature-test-match`](/pt/docs/gates//feature-test-match/): Garante que cada cenário da feature esteja realmente implementado no arquivo de teste pelo código e descrição.
- [`code-cataloged`](/pt/docs/gates//code-cataloged/): Garante que todas as funções públicas exportadas pelo código estejam documentadas na spec.

---

## 5. Onde as peças moram: Co-location vs Centralização

O Anchors suporta duas convenções para a organização física da unidade:

1. **Co-localizada (Recomendada)**: Todas as peças são irmãs no mesmo diretório:
   ```
   src/services/auth/
   ├── Login.go
   ├── Login.spec.md
   ├── Login.feature
   └── Login_test.go
   ```
2. **Centralizada / Por Região**: As peças ficam em árvores separadas (comum em alguns frameworks legados):
   ```
   src/services/auth/Login.go
   specs/auth/Login.spec.md
   features/auth/Login.feature
   tests/unit/auth/Login_test.go
   ```

Em ambos os casos, a declaração de superfícies no [`anchors.yaml`](/pt/docs/anchors-yaml//) orienta o Anchors sobre onde encontrar cada peça.

---

## 6. Próximos Passos

- [Camadas do Projeto](/pt/docs/layers/): Entenda quais camadas exigem a Unidade completa.
- [Rastreabilidade e Códigos](/pt/docs/concepts/rastreabilidade-e-codigos/): Como estruturar códigos como `AUTH-B01`.
- [Gate unit-complete](/pt/docs/gates//unit-complete/): A documentação completa do gate que vigia a integridade da unidade.
