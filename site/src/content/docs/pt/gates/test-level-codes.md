---
title: "Gate: test-level-codes"
description: "Garante que cada nível de teste referencia apenas códigos permitidos para seu escopo."
---

> **Identificador do Gate:** `test-level-codes` / `codigos-de-nivel-de-teste`  
> **Código Interno:** `TLVCD`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `Feature` |
| **Alvos Avaliados (`on`)** | `feature` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que cada nível de teste referencia apenas códigos permitidos para seu escopo.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede testes unitários de referenciarem regras visuais (-VR) e vice-versa.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Filtra os cenários de acordo com as regras de allow e exclude configuradas para cada nível no anchors.yaml.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Nenhum código proibido pelo nível foi referenciado.
- **`✗ Fail` (Reprovado):** Cenário de um nível (ex: @unit) referencia código proibido (ex: LOGIN-VR).
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Níveis sem restrição declarada.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: test-level-codes
    on: [feature]
    check: test-level-codes
    blocking: true
    levels:
      nivel-unit: { exclude: ['-VR$'] }
      nivel-vr:   { allow: ['-VR$'] } 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos com múltiplos regimes de teste.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Mova o cenário para a feature do nível correspondente ou ajuste o código de identidade da regra.
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
