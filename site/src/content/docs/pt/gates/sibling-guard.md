---
title: "Gate: sibling-guard"
description: "Impede que módulos irmãos tratem parâmetros iguais de forma inconsistente ou se alcancem por caminhos proibidos."
---

> **Identificador do Gate:** `sibling-guard` / `guarda-de-irmaos`  
> **Código Interno:** `SBGRD`  
> **Categoria:** [Fronteiras Arquiteturais](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` `code` |
| **Alvos Avaliados (`on`)** | `code` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Impede que módulos irmãos tratem parâmetros iguais de forma inconsistente ou se alcancem por caminhos proibidos.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita acoplamento cruzado desordenado entre serviços que estão no mesmo nível hierárquico.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica os caminhos de chamada e tipos de parâmetros entre módulos do mesmo diretório ou nível.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Tratamento de parâmetros consistente e sem acoplamento proibido.
- **`✗ Fail` (Reprovado):** Inconsistência de parâmetros entre funções irmãs ou acesso indevido.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: sibling-guard
    blocking: true
    measures: "módulos irmãos tratam o mesmo parâmetro de forma consistente" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos médios e grandes.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Padronize os parâmetros das funções irmãs ou utilize um DTO comum.
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
