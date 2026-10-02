---
title: "Gate: scenario-coverage"
description: "Garante que cada cenário declarado na spec tem um teste correspondente que rodou e passou."
---

> **Identificador do Gate:** `scenario-coverage` / `cobertura-de-cenario`  
> **Código Interno:** `SFMSP`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que cada cenário declarado na spec tem um teste correspondente que rodou e passou.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Fecha a rastreabilidade semântica: não basta o teste existir no arquivo, ele precisa ter rodado e sido aprovado.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Cruza os códigos de cenário da spec com a lista de códigos comprovados (ProvenCodes) no relatório ingerido.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os cenários declarados possuem teste verde no relatório.
- **`✗ Fail` (Reprovado):** Cenário declarado na spec não possui teste ou o teste falhou.
- **`~ Indeterminado/Pending`:** Relatório de testes ainda não ingerido.
- **`Skip` (Dispensado):** Camadas com dispensa de teste declarada.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: scenario-coverage
    blocking: true
    measures: "cada cenário da spec tem teste verde executado" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante no CI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Execute a suíte com `anchors test` e ingira o resultado com `anchors ingest`.
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
