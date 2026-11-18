---
title: Feature Flags
description: "Como o Anchors governa bifurcações de código e cenários alternativos através de feature flags sem duplicar especificações."
---

Quando uma equipe adota feature flags (sinalizadores de funcionalidade), um problema silencioso surge no código:

```typescript
if (isFeatureEnabled("novo-checkout")) {
    return processarCheckoutV2(carrinho);
} else {
    return processarCheckoutLegado(carrinho);
}
```

Um simples `if` acabou de criar **dois comportamentos completamente diferentes** no sistema.

Qual comportamento a especificação descreve? Se a spec descrever apenas o checkout novo, o checkout legado vira código órfão sem governança. Se você duplicar a spec inteira, terá duas specs quase idênticas divergindo a cada commit.

O Anchors resolve isso com uma camada dedicada de **Feature Flags** e a anotação **`@gated-by`**.

---

## 1. O que é um Arquivo de Flag?

No Anchors, cada feature flag importante é declarada como um artefato formal na pasta `flags/*.flag.md` (configurada como camada no [`anchors.yaml`](/pt/docs/anchors-yaml//)).

O arquivo da flag define **quais são os cenários e estados possíveis** da flag:

```markdown
<!-- @anchors
  code: FLG-CKOUT
  layer: flags
-->
# Flag: novo-checkout

## Cenários

### ON — Novo fluxo de checkout com pagamento Pix em um clique
O usuário visualiza a interface V2 e pode concluir com QR Code dinâmico.

### OFF — Fluxo tradicional legado
O usuário passa pelo formulário de 3 etapas com cartão de crédito.
```

O Anchors confronta os **cenários declarados**, e nunca o valor dinâmico em produção (que varia por usuário ou ambiente).

---

## 2. A Anotação `@gated-by`

Dentro da especificação da sua funcionalidade (ex: `Checkout.spec.md`), quando uma regra só é válida sob determinado estado da flag, você usa a anotação **`@gated-by`**:

```markdown
### CKOUT-B05 — Exibição de QR Code Pix imediato
Gera e apresenta o QR Code dinâmico na tela de confirmação.
@gated-by FLG-CKOUT.ON

### CKOUT-B06 — Redirecionamento para gateway legado
Abre a tela de checkout em três etapas.
@gated-by FLG-CKOUT.OFF
```

Dessa forma:
1. A spec continua sendo **uma só**, contendo tanto a regra nova quanto a percurso legado.
2. A IA ou o desenvolvedor sabem exatamente qual regra pertence a qual estado da flag.
3. Quando a flag for finalmente removida (deprecated), o Anchors aponta exatamente quais regras com `@gated-by FLG-CKOUT.OFF` devem ser limpas da spec e do código.

---

## 3. Os Gates que Protegem as Feature Flags

O Anchors fornece cinco [gates](/pt/docs/concepts/gates-e-vereditos/) dedicados:

- [`flag-scenario-grammar`](/pt/docs/gates//flag-scenario-grammar/): Valida se os arquivos de flag seguem a gramática esperada de estados (ex: `ON`, `OFF`).
- [`flag-scenarios-complete`](/pt/docs/gates//flag-scenarios-complete/): Garante que todos os estados declarados na flag sejam documentados.
- [`flag-scenario-exists`](/pt/docs/gates//flag-scenario-exists/): Impede que uma spec use `@gated-by` apontando para uma flag ou estado inexistente.
- [`flag-scenario-governs`](/pt/docs/gates//flag-scenario-governs/): Garante que a flag governe as regras adequadas.
- [`flag-covered`](/pt/docs/gates//flag-covered/): Comprova que ambos os ramos da flag (ON e OFF) possuem testes automatizados correspondentes.

---

## 4. Próximos Passos

- [Camadas do Projeto](/pt/docs/layers/): Como declarar a camada `flags` no `anchors.yaml`.
- [Doutrina de Produto](/pt/docs/concepts/doutrina-de-produto/): Regras transversais compartilhadas entre unidades.
- [Catálogo de Gates](/pt/docs/gates//): Conheça os gates de feature flags em detalhes.
