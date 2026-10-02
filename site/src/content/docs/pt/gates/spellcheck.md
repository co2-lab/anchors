---
title: "Gate: spellcheck"
description: "Elimina erros de digitação e ortografia em identificadores, comentários e textos."
---

> **Identificador do Gate:** `spellcheck` / `verificador-ortografico`  
> **Código Interno:** `EXCMX`  
> **Categoria:** [Segurança e Higiene Externa](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` |
| **Alvos Avaliados (`on`)** | `code` `doc` `spec` `feature` |
| **Tipo de Verificação** | `Externo (run)` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Elimina erros de digitação e ortografia em identificadores, comentários e textos.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Erros de digitação geram nomes de variáveis bizarros, quebram buscas textuais e passam impressão de descuido.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Executa o `typos` (binário nativo ultra-rápido) ou `cspell` sobre o repositório.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Sem erros ortográficos detectados.
- **`✗ Fail` (Reprovado):** Palavra com erro de digitação encontrada.
- **`~ Indeterminado/Pending`:** Ferramenta não instalada.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: spellcheck
    on: [code, doc, spec]
    scope: project
    run: "typos"
    needs_tool: typos
    install_hint: "brew install typos"
    blocking: false
    when: [pre-commit, ci] 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo no pre-commit; rode com `typos -w` para correção automática.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Corrija a palavra no texto ou cadastre-a no dicionário do projeto (`.typos.toml`).
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
