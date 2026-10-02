---
title: "Gate: marker-parity"
description: "Garante que a mesma regra aparece em ambas as pontas que a realizam (ex: frontend e backend)."
---

> **Identificador do Gate:** `marker-parity` / `paridade-de-marcador`  
> **Código Interno:** `MRPRM`  
> **Categoria:** [Falhas e Governança](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` |
| **Alvos Avaliados (`on`)** | `code` `spec` |
| **Tipo de Verificação** | `Composto / ComGate` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que a mesma regra aparece em ambas as pontas que a realizam (ex: frontend e backend).

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita assimetria regulatória: a tela diz que apaga o dado do usuário, mas o backend esquece de apagar.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se regras marcadas com prefixos simétricos possuem a contraparte correspondente na outra ponta.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A mesma regra está presente em ambas as pontas.
- **`✗ Fail` (Reprovado):** Regra simétrica presente em um lado mas ausente no outro.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: marker-parity
    blocking: true
    marker_prefix: "LGPD-"
    measures: "a mesma regra aparece nas DUAS pontas que a cumprem" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos com dados regulados (LGPD, saúde, financeiro).
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Implemente e marque a regra na ponta que faltava.
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
