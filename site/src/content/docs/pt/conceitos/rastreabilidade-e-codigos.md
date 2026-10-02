---
title: Rastreabilidade e Códigos de Identidade
description: "Por que o Anchors usa códigos de identidade (ex: AUTH-B01) para ligar a spec ao código, em vez de depender de caminhos de arquivo."
---

Se você já tentou rastrear requisitos em projetos de software antigos, provavelmente viu planilhas que diziam coisas como: *"O requisito 12 está implementado no arquivo `src/services/auth_helper.js` na linha 45"*.

O que acontece quando alguém renomeia esse arquivo para `authenticationService.ts` ou move funções para outro módulo? **Toda a rastreabilidade é destruída.** O link quebra, a planilha fica desatualizada e ninguém mais sabe se o requisito continua existindo.

No Anchors, aprendemos uma lição fundamental:

> [!IMPORTANT]
> **O código de identidade é a chave de conexão, NUNCA o caminho do arquivo.**

---

## 1. A Anatomia de um Código de Identidade

Para que uma regra seja rastreável de ponta a ponta sem ambiguidades, ela recebe um identificador curto e único, como `AUTH-B01` ou `PAY-V02`.

Veja a estrutura desse código:

```
        AUTH   ─   B   01
         │         │   │
         │         │   └── Número sequencial da regra (01, 02, 03...)
         │         │
         │         └────── Letra de natureza da regra (B = Business)
         │
         └──────────────── Prefixo da Unidade (ex: AUTH = Autenticação)
```

### O Vocabulário de Letras de Regras

As letras indicam a natureza do que está sendo verificado. O vocabulário padrão do Anchors inclui:

| Letra | Natureza da Regra | O que descreve | Exemplo |
| :---: | :--- | --- | --- |
| **`B`** | **Business Rule** | Regra de negócio pura, lógica de cálculo, decisão. | `CART-B01`: Desconto de 10% para compras acima de R$ 200. |
| **`V`** | **Validation** | Validação de formato, campos obrigatórios, tipos. | `USER-V01`: O e-mail precisa conter um `@` e domínio válido. |
| **`E`** | **Error Condition** | Tratamento explícito de falhas e exceções. | `AUTH-E01`: Retorna status 401 quando o token expirar. |
| **`S`** | **State Transition** | Máquina de estados e transições de ciclo de vida. | `PED-S01`: Pedido no estado "pago" transiciona para "em preparo". |
| **`DS`** | **Data State** | Estados de dados em tela (vazio, carregando, erro). | `HOME-DS-empty`: Mensagem exibida quando não há itens. |
| **`VR`** | **Visual Regression**| Requisito de aparência estética e fidelidade visual. | `BTN-VR`: Cor e espaçamento do botão primário. |

O gate [`rule-types`](/pt/docs/gates//rule-types/) confere se todas as letras usadas nas suas specs pertencem ao vocabulário declarado no seu [`anchors.yaml`](/pt/docs/anchors-yaml//).

---

## 2. Como o Código Conecta as Quatro Peças da Unidade

O mesmo código de identidade acompanha a regra em **todas as suas representações**:

```
1. NA SPEC (Login.spec.md)
   ### AUTH-B01 — Bloqueio de conta após três erros
   Se o usuário errar a senha por 3 vezes consecutivas, a conta é bloqueada.

2. NA FEATURE (Login.feature)
   @AUTH-B01
   Cenário: Bloqueio de conta após 3 erros consecutivos
     Dado que o usuário errou a senha 2 vezes...

3. NO TESTE (Login_test.go)
   func TestLogin_AUTH_B01_BloqueioConta(t *testing.T) {
       // Prova executável do cenário AUTH-B01
   }

4. NO CÓDIGO FONTE (Login.go)
   // #region AUTH-B01
   if user.FailedAttempts >= 3 {
       user.Locked = true
   }
   // #endregion AUTH-B01
```

Se o desenvolvedor renomear o arquivo `Login.go` para `AuthService.go`, ou mover o teste para outra pasta, **a rastreabilidade continua 100% intacta**, porque os gates procuram pelo identificador `AUTH-B01`, e não pelo nome do arquivo no disco.

---

## 3. As Vantagens Práticas para o seu Dia a Dia

1. **Busca instantânea no repositório**: Quer saber onde uma regra está implementada? Basta dar um `grep` ou busca global no VSCode por `AUTH-B01`. Você achará a spec, o cenário em Gherkin, o teste e a linha exata do código fonte em 1 segundo.
2. **Sem testes órfãos**: O gate [`test-traceable`](/pt/docs/gates//test-traceable/) acusa testes que não citam o código do cenário que estão provando.
3. **Sem regras esquecidas**: O gate [`spec-feature-match`](/pt/docs/gates//spec-feature-match/) garante que nenhuma regra da spec fique sem cenário na feature.
4. **Alinhamento entre cenário e teste**: O gate [`feature-test-match`](/pt/docs/gates//feature-test-match/) verifica se o teste ligado realmente exercita o que o cenário promete.

---

## 4. Próximos Passos

- [A Unidade](/pt/docs/concepts/unidade/): Veja como o código de identidade une a spec ao teste.
- [Propagação e Impacto](/pt/docs/concepts/propagacao-e-impacto/): O que acontece quando você altera o código de uma regra.
- [Gate code-cataloged](/pt/docs/gates//code-cataloged/): Como garantir que todo export do código esteja catalogado na spec.
