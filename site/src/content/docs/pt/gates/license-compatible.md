---
title: "Gate: license-compatible"
description: "Impede a inclusão de dependências com licenças incompatíveis ou copyleft forte (como AGPL)."
---

> **Identificador do Gate:** `license-compatible` / `licenca-compativel`  
> **Código Interno:** `EXCMX`  
> **Categoria:** [Segurança e Higiene Externa](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` `code` |
| **Alvos Avaliados (`on`)** | `code` |
| **Tipo de Verificação** | `Externo (run)` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Impede a inclusão de dependências com licenças incompatíveis ou copyleft forte (como AGPL).

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita riscos jurídicos de contaminação de código proprietário por licenças não comerciais.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Executa `go-licenses`, `license-checker` ou `cargo-deny` comparando com a política de licenças.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as dependências possuem licenças permitidas.
- **`✗ Fail` (Reprovado):** Dependência com licença proibida ou desconhecida detectada.
- **`~ Indeterminado/Pending`:** Ferramenta de licenças não instalada.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: license-compatible
    on: [code]
    scope: project
    run: "go-licenses check ./... --disallowed_types=forbidden,restricted"
    needs_tool: go-licenses
    blocking: true
    when: [pre-push, ci] 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante no CI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Substitua a biblioteca com licença proibida por uma alternativa de licença permissiva (MIT, Apache-2.0, BSD).
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
