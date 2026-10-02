---
title: "Gate: placeholder-filled"
description: "Verifica se os placeholders deixados por geradores ou templates foram preenchidos."
---

> **Identificador do Gate:** `placeholder-filled` / `placeholder-preenchido`  
> **Código Interno:** `PLCFL`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` |
| **Alvos Avaliados (`on`)** | `spec` `plan` `feature` `doc` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Verifica se os placeholders deixados por geradores ou templates foram preenchidos.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede commitar arquivos com textos como '[Descreva aqui]' ou 'TODO: preencher'.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Busca marcadores universais de template não preenchidos no corpo do documento.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Nenhum placeholder de template esquecido.
- **`✗ Fail` (Reprovado):** Encontrado texto de placeholder não substituído.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: placeholder-filled
    on: [spec, feature]
    check: placeholder-filled
    blocking: true
    measures: "o esqueleto emitido pelo gerador foi preenchido" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para evitar entrega de documentação pela metade.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Substitua os textos de placeholder pelo conteúdo real da regra.
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
