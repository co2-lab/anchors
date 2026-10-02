---
title: "Gate: spec-sections"
description: "Verifica se a spec cataloga regras estruturadas (em cabeçalho, tabela ou lista) e usa o idioma correto."
---

> **Identificador do Gate:** `spec-sections` / `spec-completa`  
> **Código Interno:** `SFMSP`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `Spec` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Verifica se a spec cataloga regras estruturadas (em cabeçalho, tabela ou lista) e usa o idioma correto.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Uma spec com texto solto em prosa não permite que a IA ou o CLI indexem as regras. Além disso, evita misturar idiomas nas seções.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Busca regras com código em cabeçalhos (###), linhas de tabela (|) ou bullets em negrito. Verifica também se os títulos batem com o lang configurado.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Pelo menos uma regra está catalogada de forma estruturada e os títulos seguem o idioma configurado.
- **`✗ Fail` (Reprovado):** A spec só contém prosa corrida sem regras catalogadas, ou possui seções em idioma diferente do lang declarado.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Camadas declarativas ou specs de coordenação.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: spec-sections
    on: [spec]
    check: spec-sections
    blocking: true
    measures: "a spec tem ao menos uma regra catalogada estruturada" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante desde o início. Garante que os templates de spec sejam preenchidos de verdade.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Estruture suas regras usando títulos com código (### CODE-B01) ou tabelas, e use o idioma padrão do projeto nas seções.
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
