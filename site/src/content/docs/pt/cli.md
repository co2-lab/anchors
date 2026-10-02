---
title: "O CLI do Anchors"
description: "A ferramenta de linha de comando oficial para pair programming com IA e governança contínua spec-first."
---

O **CLI do Anchors** (`anchors`) é um binário rápido e estático escrito em Go que sustenta a governança contínua de software desenvolvido com inteligência artificial.

Em vez de embutir IA dentro da ferramenta, o Anchors foi desenhado como **a ferramenta que os agentes de IA operam**:
- A IA não precisa decorar o framework: ela consulta o binário (`anchors guide`), entende a rota e opera através dos comandos.
- Totalmente agnóstico ao cliente: funciona com Claude Code, Cursor, Windsurf, Copilot, Gemini CLI ou desenvolvedores humanos no terminal.
- Opera estritamente sobre texto e contratos de arquivos — sem servidores ocultos ou dependências pesadas.

---

## 🚀 Navegação Rápida

- [**Guia de Instalação**](/pt/docs/cli/installation/) — Instale via Homebrew, script shell, Go, Windows, Docker e CI/CD.
- [**Referência de Comandos**](/pt/docs/cli/commands/) — Índice completo e pesquisável dos 55 comandos do CLI.
- [**O fluxo de trabalho**](/pt/docs/workflow/) — Como operar o ciclo no dia a dia, da spec ao merge seguro.
- [**O anchors.yaml**](/pt/docs/anchors-yaml/) — O arquivo central de configuração do seu projeto.

---

## 🛠️ Comandos Essenciais

```bash
# 1. Inicializa um novo projeto ou configura uma base existente
anchors init

# 2. Diagnostica a saúde e maturidade estrutural do repositório
anchors doctor

# 3. Constrói e consulta o grafo de dependências
anchors map build
anchors impact src/services/auth/login.spec.md

# 4. Executa os gates de qualidade contra o stage ou projeto inteiro
anchors check
anchors check --all

# 5. Inicia o monitor de tarefas em background para IAs
anchors watch
```

---

## 🧭 Guias Detalhados de Comandos

- [**anchors check & verify**](/pt/docs/cli/commands/check/) — Avaliação de gates, modo estrito e filtros.
- [**anchors doctor & audit**](/pt/docs/cli/commands/doctor/) — Raio-X do ecossistema, artefatos órfãos e auditoria de arquivos.
- [**anchors map & impact**](/pt/docs/cli/commands/map/) — Grafo de dependências, raio de impacto e recodificação.
- [**anchors init & new**](/pt/docs/cli/commands/init/) — Assistente de setup, descoberta de camadas e geração de artefatos.
- [**anchors flow & queue**](/pt/docs/cli/commands/flow/) — Watcher em background, fila de tarefas e entregas autônomas.
- [**anchors freeze & thaw**](/pt/docs/cli/commands/freeze/) — Bloqueio de releases e congelamento de governança.
