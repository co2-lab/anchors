---
title: Doutrina de Produto
description: "O que é a Doutrina de Produto, como declarar regras transversais sem duplicar código e como usar a anotação @realizes."
---

Nem toda regra de negócio pertence a uma única tela ou função isolada.

Pense em regras como estas:
- *"Qualquer valor monetário exibido na aplicação deve seguir o padrão brasileiro (`R$ 1.234,56`)."*
- *"Se o usuário ficar inativo por 15 minutos, a sessão deve expirar em qualquer tela."*
- *"Dados sensíveis (como CPF ou cartão) nunca podem ser gravados em logs abertos."*

Se você copiar e colar essas regras na spec de cada tela do seu sistema, acontecerá o desastre inevitável: no dia em que a regra mudar de 15 para 20 minutos, alguém atualizará cinco specs e esquecerá as outras dez. O produto diverge e ninguém percebe.

Para resolver isso, o Anchors introduz a **Doutrina de Produto (Product Doctrine)**.

---

## 1. O que é uma Doutrina de Produto?

Uma **Doutrina de Produto** é uma especificação centralizada de regras **transversais** que cortam múltiplas unidades do projeto.

Elas ficam localizadas em uma camada dedicada (geralmente `product/*.doctrine.md`), declarada no [`anchors.yaml`](/pt/docs/anchors-yaml//).

```
product/
├── formatacao_moeda.doctrine.md
├── seguranca_sessao.doctrine.md
└── protecao_dados.doctrine.md
```

Em vez de cada tela inventar a sua própria regra, o arquivo de doutrina é a **única fonte da verdade**.

---

## 2. A Anotação `@realizes`

Como o Anchors sabe que uma tela ou serviço do seu sistema está cumprindo a Doutrina de Produto? Através da anotação **`@realizes`**.

Quando uma spec local implementa ou respeita uma regra da doutrina, ela aponta para o código da regra central:

### 1. No arquivo de Doutrina (`product/seguranca_sessao.doctrine.md`):
```markdown
### SESS-D01 — Expiração por inatividade
Toda tela autenticada deve invalidar a sessão após 15 minutos sem interação do usuário.
```

### 2. Na Spec da Tela (`screens/PerfilUsuario.spec.md`):
```markdown
### PERFIL-B04 — Encerramento de sessão
Encerra a sessão e redireciona para a tela de login.
@realizes SESS-D01
```

O Anchors conecta as duas pontas no [Grafo de Dependências](/pt/docs/concepts/grafo-e-mapa/).

---

## 3. Os Gates que Protegem a Doutrina

O Anchors fornece quatro [gates](/pt/docs/concepts/gates-e-vereditos/) especializados para governar a doutrina:

- [`spec-doctrine-exists`](/pt/docs/gates//spec-doctrine-exists/): Garante que toda anotação `@realizes` aponte para um código de doutrina que realmente existe no repositório.
- [`doctrine-realized`](/pt/docs/gates//doctrine-realized/): Acusa se uma regra de doutrina foi declarada como obrigatória no projeto, mas nenhuma unidade do sistema a realiza.
- [`doctrine-not-duplicated`](/pt/docs/gates//doctrine-not-duplicated/): Impede que o texto de uma regra de doutrina seja copiado e colado dentro de specs locais (forçando o uso de `@realizes`).
- [`spec-realizes-doctrine`](/pt/docs/gates//spec-realizes-doctrine/): Garante que a spec comprove como ela cumpre a doutrina.

---

## 4. Próximos Passos

- [Camadas do Projeto](/pt/docs/layers/): Veja como configurar a camada `product` no `anchors.yaml`.
- [Feature Flags](/pt/docs/concepts/feature-flags/): Como governar variações de código sem poluir a spec base.
- [Catálogo de Gates](/pt/docs/gates//): Conheça os gates de doutrina em detalhes.
