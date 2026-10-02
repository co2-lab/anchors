---
title: "Gate: sbom-generated"
description: "Gera o inventário de software SBOM (CycloneDX ou SPDX) para compliance e auditorias."
---

> **Identificador do Gate:** `sbom-generated` / `sbom-gerado`  
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

Gera o inventário de software SBOM (CycloneDX ou SPDX) para compliance e auditorias.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Essencial para conformidade com padrões de segurança em distribuição de software corporativo.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Executa ferramentas como `syft` gerando o arquivo `sbom.json`.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Arquivo de SBOM gerado com sucesso.
- **`✗ Fail` (Reprovado):** Falha na geração do inventário.
- **`~ Indeterminado/Pending`:** Syft não instalado.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: sbom-generated
    on: [code]
    scope: project
    run: "syft scan dir:. -o cyclonedx-json=sbom.json -q"
    needs_tool: syft
    install_hint: "brew install syft"
    blocking: true
    when: [ci] 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante no CI de release.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Instale o `syft` e verifique as permissões de escrita do arquivo sbom.json.
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
