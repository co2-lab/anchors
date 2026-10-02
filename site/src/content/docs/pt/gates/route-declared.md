---
title: "Gate: route-declared"
description: "Garante que uma tela ou endpoint declara como se chega nela e nomeia seus vizinhos de navegação."
---

> **Identificador do Gate:** `route-declared` / `rota-declarada`  
> **Código Interno:** `RTDCL`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `UI` `Telas` `API` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Interno Determinístico` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que uma tela ou endpoint declara como se chega nela e nomeia seus vizinhos de navegação.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede o surgimento de telas órfãs que existem no repositório mas ninguém sabe como acessar.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica a presença da seção de Rotas ou Cabeçalho de Navegação na spec.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A rota de entrada e vizinhos estão formalmente declarados.
- **`✗ Fail` (Reprovado):** Spec de tela ou rota sem declaração de caminho de acesso.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Camadas que não representam pontos de entrada de usuário ou API.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: route-declared
    on: [spec]
    check: route-declared
    blocking: true
    measures: "uma tela declara como se chega nela e seus vizinhos" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos de frontend (React, Flutter, Vue) e APIs.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Declare a rota na seção de Identidade ou Rotas da spec.
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
