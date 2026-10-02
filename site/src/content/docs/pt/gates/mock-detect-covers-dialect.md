---
title: "Gate: mock-detect-cobre-o-dialeto"
description: "Avalia se a regex declarada para detectar dublês alcança todas as formas que o projeto usa."
---

> **Identificador do Gate:** `mock-detect-cobre-o-dialeto` / `mock-detect-cobre-dialeto`  
> **Código Interno:** `MCSTM`  
> **Categoria:** [Dublês e Mocks](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `test` |
| **Alvos Avaliados (`on`)** | `test` |
| **Tipo de Verificação** | `Julgamento por IA (ask)` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Avalia se a regex declarada para detectar dublês alcança todas as formas que o projeto usa.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Se a regex de mock-detect falhar, dublês passam sem validação de contrato.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Pergunta ao modelo se existem padrões de mock no repositório que escaparam da regex configurada.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Veredito `pass` emitido por julgamento.
- **`✗ Fail` (Reprovado):** Veredito `fail` acusando padrões de mock não cobertos.
- **`~ Indeterminado/Pending`:** Julgamento não registrado.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: mock-detect-cobre-o-dialeto
    on: [test]
    ask: "o padrão de detecção de mocks alcança todas as formas usadas no repositório?"
    blocking: false 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo; avalie periodicamente com `anchors judge`.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Ajuste o regex `mock_detect` no anchors.yaml e emita novo veredito.
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
