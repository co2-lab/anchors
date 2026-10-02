---
title: "anchors init & new"
description: "Configuração do projeto a partir das próprias pastas e da linguagem, o guia de contribuição, git hooks e geração de artefatos."
---

Iniciar um projeto novo ou plugar o Anchors numa base de código existente começa pelo `anchors init`. Ele lê o projeto como ele é e pergunta só as decisões que são suas.

---

## 1. Inicialização do Projeto (`anchors init`)

```bash
# Configuração interativa
anchors init

# Para um agente ou um script: primeiro as perguntas, em JSON…
anchors init --non-interactive
# …depois as respostas, em flags
anchors init --non-interactive --artifacts=spec,feature,test,code --colocation
# ou aceitar todos os defaults lidos do disco
anchors init --non-interactive --defaults
```

O `anchors init`:

1. **Não propõe estrutura.** O Anchors não se importa com a organização das pastas — só que as camadas fiquem separadas. Toda pasta com código é candidata a camada, com o nome da pasta, e você mantém as que são camadas. Tipos de camada como *pontos de entrada, casos de uso, domínio, repositórios, infraestrutura, apresentação* são ilustração, nunca um layout para mover arquivos.
2. **Lê a linguagem como dialeto.** A família vem do manifesto na raiz (`go.mod`, `package.json`, `pyproject.toml`…) ou da extensão mais comum, e o nome dos testes vem dos seus próprios arquivos de teste (`*_test.go`, `*.spec.ts`, `test_*.py`…). Fica gravada em `dialect.family`, e os gates leem seus testes e comentários desde o primeiro check.
3. **Semeia os gates que têm relação com as suas camadas** — um gate é semeado quando uma camada declarada é de um tipo que ele mede e os campos que ele pressupõe estão declarados. O resto do catálogo que cobre suas camadas é listado pelo [`anchors doctor`](/pt/docs/cli/commands/doctor/).
4. **Semeia o `CONTRIBUTING.md`** a partir da configuração que grava: a ordem do trabalho (spec → feature → teste → código), as camadas declaradas, os comandos do dia a dia e quais gates barram um commit. Um `CONTRIBUTING.md` que você já tem nunca é tocado; o trecho que seria acrescentado é mostrado.
5. **Semeia o guia de cabeçalho** e escreve o [`anchors.yaml`](/pt/docs/anchors-yaml/).

Ao semear gates com `run:`, ele também diz como escrever um passo de gate que faz mais que chamar uma ferramenta: na língua do projeto (`go run ./tools/gates <gate>`, `node tools/gates.mjs <gate>`, `python -m tools.gates <gate>`), não em shell — um script de shell quebra na plataforma.

| Flag | O que responde |
| --- | --- |
| `--artifacts` | os tipos de âncora que o projeto usa (`spec,feature,test,code,guide,plan`) |
| `--layers` | as pastas de código a manter como camadas |
| `--colocation` | spec, feature e teste ficam ao lado do código |
| `--gates` | semear os gates padrão (`true`) |
| `--header` | semear o guia de cabeçalho (`true`) |
| `--contributing` | semear o `CONTRIBUTING.md` quando o projeto não tem um (`true`) |
| `--workflow`, `--repo`, `--labels` | onde mora a fila de trabalho: `local`, `manual` ou `github` |
| `--governs GUIA=tag1,tag2` | quais tags um guia rege |

---

## 2. Geração de Artefatos (`anchors new`)

Gera um artefato com o cabeçalho `@anchors` e a identidade já resolvidos. `--out` é obrigatório: o artefato nasce ao lado da unidade que descreve.

```bash
# Uma spec ganha um código único novo
anchors new spec Fatura --out src/billing/fatura.spec.md

# Uma feature ou um teste ganha um ref para a spec
anchors new feature Fatura --out src/billing/fatura.feature --code FATUR

# As seções que um tipo oferece, e os presets dele
anchors new spec --list-sections
```

---

## 3. Git Hooks (`anchors install-hooks`)

```bash
anchors install-hooks
```

Instala o pre-commit que roda os gates sobre o que está em stage, para que um commit que reprova um gate bloqueante não entre.
