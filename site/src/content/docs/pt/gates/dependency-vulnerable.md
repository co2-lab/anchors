---
title: "Gate: dependency-vulnerable"
description: "Audita dependências e lockfiles contra bases públicas de vulnerabilidades conhecidas (CVEs)."
---

> **Identificador do Gate:** `dependency-vulnerable` / `dependencia-vulneravel`  
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

Audita dependências e lockfiles contra bases públicas de vulnerabilidades conhecidas (CVEs).

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Usar pacotes vulneráveis abre portas para invasões e violações de segurança.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Invoca ferramentas como `osv-scanner`, `govulncheck`, `pnpm audit` ou `pip-audit`.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Nenhuma CVE conhecida encontrada nas dependências.
- **`✗ Fail` (Reprovado):** Vulnerabilidade conhecida detectada.
- **`~ Indeterminado/Pending`:** Ferramenta de auditoria não instalada.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: dependency-vulnerable
    on: [code]
    scope: project
    run: "osv-scanner scan source -r ."
    needs_tool: osv-scanner
    install_hint: "brew install osv-scanner"
    blocking: true
    when: [pre-push, ci] 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante no pre-push e CI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Atualize a versão da biblioteca vulnerável no seu package.json, go.mod ou requirements.txt.
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
