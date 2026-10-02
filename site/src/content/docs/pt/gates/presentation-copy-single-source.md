---
title: "Gate: presentation-copy-single-source"
description: "Garante que o texto exibido na apresentação vem de um código de mensagem central, e não de texto repetido."
---

> **Identificador do Gate:** `presentation-copy-single-source` / `texto-de-apresentacao-fonte-unica`  
> **Código Interno:** `PRSNT`  
> **Categoria:** [Apresentação e Telas](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `UI` `Telas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que o texto exibido na apresentação vem de um código de mensagem central, e não de texto repetido.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Facilita internacionalização (i18n) e impede que o mesmo texto mude em uma tela e continue antigo noutra.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se os textos de interface apontam para identificadores de strings ou chaves de i18n.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Textos referenciados por códigos de mensagens.
- **`✗ Fail` (Reprovado):** Textos literais repetidos na regra sem chave de catálogo.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: presentation-copy-single-source
    on: [spec]
    check: presentation-copy-single-source
    blocking: true
    measures: "o texto que a tela mostra vem de um código de mensagem" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Recomendado para apps multilíngues.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Extraia o texto para o catálogo de mensagens/i18n e use a chave correspondente na spec.
2. Reexecute a verificação no terminal:
   ```sh
   anchors check
   ```
3. Se o gate suportar correção automática, você pode tentar o comando:
   ```sh
   anchors check --fix
   ```

---

## 🔗 Conceitos Relacionados

- [Guia Completo de Camadas](/pt/docs/layers/): Entenda quais camadas exigem este gate.
- [A Unidade (The Unit)](/pt/docs/concepts/unidade/): A relação entre Spec, Feature, Teste e Código.
- [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/): A mecânica completa de avaliação do Anchors.
- [Catálogo Completo de Gates](/pt/docs/gates//): Retornar ao índice pesquisável de todos os gates.
