# Anchors — Levantamento visual e de conteúdo (Design Audit)

> Documento de leitura, gerado para alimentar a criação de um design system.
> Nenhum arquivo do repositório foi alterado para produzir este levantamento.

---

## 1. Stack e estrutura

**Versões (de `site/node_modules/*/package.json`, resolvidas a partir de `site/package.json`):**

- Astro: `7.2.1` (declarado como `^7.1.6` em `site/package.json`)
- `@astrojs/starlight`: `0.41.7` (declarado como `^7.1.6`... não, declarado como `^0.41.7`)
- `sharp`: `^0.35.3` (não resolvido em node_modules na leitura; versão declarada)
- Tailwind: não existe. Nenhuma referência a `tailwind` em `package.json`, `astro.config.mjs`, ou no repo (fora de `node_modules`).
- Lib de UI/componentes: não existe. Nenhum design-system framework (nenhum shadcn, nenhum Radix, nenhum UnoCSS). CSS puro.

**Árvore de `site/` (até 3 níveis, só diretórios + arquivos de config/estilo, excluindo `node_modules/`, `dist/`, `.astro/`):**

```
site/
├── .gitignore
├── astro.config.mjs
├── package.json
├── package-lock.json
├── tsconfig.json
├── public/
│   ├── anchors-icon-black.svg
│   ├── anchors-icon-red.svg
│   ├── anchors-lockup-black.svg
│   ├── anchors-lockup-white.svg
│   ├── anchors-mark-black.svg
│   ├── anchors-mark-red.svg
│   ├── anchors-mark-white.svg
│   ├── co2lab.png
│   └── favicon.svg
├── src/
│   ├── content.config.ts
│   ├── components/
│   │   └── Landing.astro
│   ├── content/
│   │   └── docs/
│   │       ├── en/        (10 arquivos .md)
│   │       └── pt/        (10 arquivos .md)
│   ├── i18n/
│   │   └── ui.ts
│   ├── pages/
│   │   ├── index.astro
│   │   └── [lang]/
│   │       └── index.astro
│   └── styles/
│       └── landing.css
```

**Customização do tema Starlight:**

- Não existe `customCss` no `astro.config.mjs` — nenhum CSS é injetado no Starlight.
- Não existe chave `components:` no config — nenhum componente do Starlight (Header, Sidebar, PageTitle, etc.) foi sobrescrito.
- Não existe `expressiveCode`, `tableOfContents` customizado, nem `pagination` customizada — tudo default.
- A única customização declarada é: `title`, `favicon`, `logo` (`light`/`dark`, dois SVGs diferentes), `defaultLocale`, `locales`, `social` (link do GitHub), e `sidebar` (grupos e itens nomeados manualmente, sem `autogenerate`).
- Conclusão: **as docs (`/docs/**`) usam 100% do tema visual default do Starlight** — cor, tipografia, componentes (Card, Aside, Tabs etc.) são os do pacote `@astrojs/starlight@0.41.7`, não personalizados neste repo. Só a landing (`/`, `/en/`, `/es/`, `/ja/`, `/zh/`) tem CSS próprio.

---

## 2. Tokens de cor

### 2.1 Conteúdo integral dos arquivos CSS globais

Só existe **um** arquivo CSS global no repo: `site/src/styles/landing.css` (715 linhas, carregado só nas páginas de landing via `import '../styles/landing.css'` em `src/pages/index.astro` e `src/pages/[lang]/index.astro`).

**Não existe** nenhum `:root` do Starlight sobrescrito, nenhum `[data-theme=dark]` custom, nenhum arquivo CSS a mais.

Bloco `:root` completo (`site/src/styles/landing.css:2-21`):

```css
:root {
	--bg: #060a0c;
	--bg-2: #0a1116;
	--panel: #10181e;
	--panel-2: #152029;
	--border: #223038;
	--text: #e8f0f2;
	--muted: #9ab0b8;
	--faint: #6b8088;
	--acc: #2dd4bf;
	--acc-2: #38bdf8;
	--acc-3: #fb923c;
	--structure: #38bdf8;
	--planning: #a78bfa;
	--spec: #2dd4bf;
	--traceability: #fbbf24;
	--propagation: #fb923c;
	--quality: #f472b6;
	--maxw: 1120px;
}
```

Não há bloco `@media (prefers-color-scheme: dark)` nem `[data-theme]` neste arquivo — **a landing é single-theme (dark only)**, não responde a light mode.

### 2.2 Tabela de variáveis

| Variável | Valor | Onde é usada (exemplos) |
|---|---|---|
| `--bg` | `#060a0c` | `body { background }` — `landing.css:33` |
| `--bg-2` | `#0a1116` | gradiente do `.card` — `landing.css:514` |
| `--panel` | `#10181e` | fundo de `.problem-card`, `.code-window`, `.stat`, `.nav-cta`, `.btn-ghost` |
| `--panel-2` | `#152029` | fundo de `.code-bar`, `.lang-menu`, chip de código inline (`.section-lead code`) |
| `--border` | `#223038` | borda de quase todo componente com `border: 1px solid var(--border)` |
| `--text` | `#e8f0f2` | `body { color }` — `landing.css:34` |
| `--muted` | `#9ab0b8` | corpo de texto secundário: `.nav-links`, `.hero-sub`, `.problem-card p`, `.card p` |
| `--faint` | `#6b8088` | texto terciário: `.hero-note`, `.code-file`, `.footer`, `.stats-note` |
| `--acc` | `#2dd4bf` (teal) | `.problem-num`, `.grad` (início do gradiente), `.btn-primary` (início do gradiente), `.nav-cta:hover` |
| `--acc-2` | `#38bdf8` (azul) | `.eyebrow`, `.brand-name` (gradiente), `.grad` (fim), `.section-lead code` |
| `--acc-3` | `#fb923c` (laranja) | `.grad` (meio), `.problem-punch strong` |
| `--structure` | `#38bdf8` | `.chip.structure` — `landing.css:463-466` |
| `--planning` | `#a78bfa` (roxo) | `.chip.planning` — declarada mas **não usada em nenhum HTML** (ver §11) |
| `--spec` | `#2dd4bf` | `.chip.spec` |
| `--traceability` | `#fbbf24` (amarelo) | `.chip.traceability` — declarada mas **não usada em nenhum HTML** (ver §11) |
| `--propagation` | `#fb923c` | `.chip.propagation` — declarada mas **não usada em nenhum HTML** (ver §11) |
| `--quality` | `#f472b6` (rosa) | `.chip.quality` |
| `--maxw` | `1120px` | `max-width` de `.nav`, `.section`, `.footer` |

**Cores hard-coded fora das variáveis** (não centralizadas em token, usadas diretamente):

| Valor | Onde | Propósito |
|---|---|---|
| `#fff` (branco puro) | `.brand-name` (gradiente), `.btn-primary` (`color: #fff`) | texto sobre gradiente |
| `#6d28d9` (roxo escuro) | `.btn-primary { background: linear-gradient(100deg, var(--acc), #6d28d9) }` | fim do gradiente do botão primário — não é nenhum token declarado |
| `rgba(124, 92, 255, 0.25)`, `rgba(124, 92, 255, 0.18)`, `rgba(124, 92, 255, 0.1)`, `rgba(124, 92, 255, 0.08)` | `.bg-glow`, `.badge`, `.card-icon`, `.pill`, `.cta-final` | roxo translúcido (`#7c5cff` aprox.) usado em glows e fundos — **não corresponde a nenhuma variável de `:root`**; é resíduo do tema original do Bavel (ver §11) |
| `rgba(34, 211, 238, 0.14)` | `.bg-glow` | glow ciano |
| `rgba(244, 114, 182, 0.1)` | `.bg-glow` | glow rosa |
| `#4ade80` | `.badge .dot` (bolinha "em desenvolvimento ativo") | verde, único uso, não tokenizado |
| `#ff5f56`, `#ffbd2e`, `#27c93f` | `.d.red`, `.d.yellow`, `.d.green` (bolinhas estilo macOS da code window) | cores fixas de referência, não tokenizadas |
| `#cdd6f4` | `.code` (cor base do texto do bloco de código) | não tokenizado |
| `#6b7089`, `#c792ea`, `#c3e88d`, `#f78c6c` | `.code .cm`, `.code .kw`, `.code .st`, `.code .nu` (syntax highlight do bloco de código fake) | paleta estilo "Tokyo Night" / Catppuccin, não tokenizada, não relacionada às cores dos pilares |
| `#c9bfff` | `.pill` (cor do texto) | roxo claro, não tokenizado, órfão do tema Bavel |

**Nenhuma dessas cores é do Starlight** — o Starlight não é tocado neste arquivo. As cores do Starlight (usadas em `/docs/**`) são inteiramente as default do pacote `@astrojs/starlight@0.41.7` e não estão declaradas em nenhum arquivo deste repo.

**Tailwind:** não existe — não há `theme.extend` para colar.

---

## 3. Tipografia

**Fontes usadas — landing (`landing.css:33-40`):**

```css
font-family:
	ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto,
	Helvetica, Arial, sans-serif;
```

- 100% system font stack. Não há `@font-face`, não há `<link>` para Google Fonts, não há fonte self-hosted. Não existe nenhum arquivo de fonte (`.woff`, `.woff2`, `.ttf`) em nenhum lugar do repo.
- Fonte mono (`landing.css:48-51`, aplicada a `code`):
  ```css
  font-family:
  	ui-monospace, "SF Mono", "JetBrains Mono", Menlo, Consolas,
  	monospace;
  ```
  Também system stack — nenhuma fonte mono é carregada externamente, apesar de "JetBrains Mono" aparecer na lista (só é usada se já estiver instalada localmente no SO do visitante).

**Fontes usadas — docs (Starlight):** não customizadas neste repo — o que o Starlight aplica por padrão (Atkinson Hyperlegible para o corpo, ou o que a versão 0.41.7 definir; não verificável sem inspecionar o pacote, e este levantamento não abre `node_modules` do Starlight para citar como "nosso").

**Fonte usada nos logos** (`site/public/anchors-lockup-*.svg`, atributo `font-family` do elemento `<text>`):
```
font-family="Archivo, 'Helvetica Neue', Arial, sans-serif"
```
"Archivo" não é carregada em lugar nenhum do site — só aparece hard-coded dentro do SVG do lockup. Se o SVG for renderizado num ambiente sem a fonte Archivo instalada, cai para Helvetica Neue/Arial.

**Escala real (landing, valores literais do CSS):**

| Elemento | font-size | line-height | font-weight | letter-spacing |
|---|---|---|---|---|
| `.hero-title` (H1) | `clamp(2.6rem, 7vw, 4.6rem)` | `1.05` | `800` | `-0.03em` |
| `.section-head h2` (H2) | `clamp(1.8rem, 4vw, 2.7rem)` | `1.15` | `800` | `-0.02em` |
| `.cta-final h2` (H2 variante) | `clamp(1.9rem, 4vw, 2.6rem)` | não declarado (herda `body`: `1.6`) | `800` | `-0.02em` |
| `.problem-card h3` / `.card h3` (H3) | `1.2rem` | herda | não declarado (herda default do navegador, `bold` por tag `<h3>`) | não declarado |
| `.hero-sub` (subhead) | `clamp(1.05rem, 2.2vw, 1.3rem)` | herda `1.6` | não declarado | não declarado |
| `.section-lead` (lead de seção) | `1.08rem` | herda | não declarado | não declarado |
| body (default) | tamanho não declarado (herda `1rem` do navegador) | `1.6` (`landing.css:38`) | não declarado | não declarado |
| `.problem-card p` / `.card p` (corpo pequeno) | `0.96rem` / `0.98rem` | herda | não declarado | não declarado |
| `.stats-note`, `.hero-note` (small/caption) | `0.95rem` | herda | não declarado | não declarado |
| `.eyebrow` (label pequena, uppercase) | `0.8rem` | não declarado | `700` | `0.18em` |
| `.stat-num` (número grande de estatística) | `clamp(2.2rem, 5vw, 3rem)` | `1` | `800` | não declarado |
| `.code` (bloco de código) | `0.86rem` (`0.76rem` abaixo de 480px) | `1.75` | não declarado | não declarado (`tab-size: 2`) |

Não existe nenhum outro nível de heading tokenizado (H4, H5, H6) — o CSS não define nada além de h2/h3 explícitos.

---

## 4. Espaçamento, borda, elevação, movimento

**Espaçamento (valores mais frequentes, extraídos literalmente do CSS — não há escala nomeada, tudo é valor direto em px/rem):**

- `24px` — padding horizontal padrão de containers (`.nav`, `.section`, `.footer`)
- `20px` — gap mais comum entre cards em grid (`.problem-grid`, `.cards`, `.stats`)
- `16px`, `14px`, `12px`, `10px`, `8px`, `6px` — gaps/paddings menores, usados sem padrão aparente de progressão (não é uma escala geométrica nem um múltiplo fixo — ex.: `.nav` tem `padding: 22px 24px`, `.badge` tem `padding: 6px 14px`, `.btn` tem `padding: 13px 26px`)
- `70px` — padding vertical de `.section` (`padding: 70px 24px`)
- `40px` — margem vertical recorrente entre grandes blocos (`.hero`, `.problem-punch`)

**border-radius (todos os valores encontrados no arquivo, em ordem de aparição):**

| Valor | Onde |
|---|---|
| `999px` (pill/full) | `.badge`, `.nav-cta`, `.chip`, `.pill` |
| `24px` | `.cta-final` |
| `16px` | `.problem-card`, `.code-window`, `.card`, `.stat` |
| `12px` | `.btn`, `.lang-menu` |
| `11px` | `.card-icon` |
| `8px` | `.lang-menu a` |
| `6px` | `.section-lead code` / `.problem-punch code` / `.flow-caption code` |
| `50%` (círculo) | `.badge .dot`, `.code-bar .d` |

Não há um único valor reutilizado consistentemente para "cards" — `.problem-card`, `.code-window`, `.card`, `.stat` compartilham `16px`, mas `.cta-final` (que também é um "card" grande) usa `24px`, e o botão usa `12px`. Não há uma escala nomeada (ex. `--radius-sm/md/lg`) declarada em variável.

**box-shadow (todos os valores encontrados, literais):**

| Valor | Onde |
|---|---|
| `0 8px 30px rgba(124, 92, 255, 0.4)` | `.btn-primary` |
| `0 12px 40px rgba(0, 0, 0, 0.5)` | `.lang-menu` |
| `0 30px 80px rgba(0, 0, 0, 0.5)` | `.code-window` |
| `0 0 10px #4ade80` | `.badge .dot` (glow do ponto verde) |

Note que o shadow de `.btn-primary` usa a mesma cor roxa órfã (`rgba(124, 92, 255, ...)`) do §2 — não usa `--acc` (teal), apesar do botão em si usar `--acc` no gradiente de fundo.

**Transições (todas as declarações, literais):**

| Propriedade | Duração | Easing | Onde |
|---|---|---|---|
| `transform` | `0.15s` | `ease` | `.btn` |
| `box-shadow` | `0.2s` | `ease` | `.btn` |
| `border-color` | `0.2s` | `ease` | `.btn`, `.card` |
| `transform` | `0.2s` | `ease` | `.card` |
| `opacity` | `0.2s` | `ease` | `.footer-org` |

Não há nenhuma declaração `@keyframes` no repo — nenhuma animação, só transitions em hover.

**Estados (hover/focus/active) — valores exatos por elemento:**

| Elemento | Hover | Focus-visible | Active/pressed |
|---|---|---|---|
| `.btn` | `transform: translateY(-2px)` | não declarado | não declarado |
| `.btn-ghost` | `border-color: var(--acc-2)` | não declarado | não declarado |
| `.nav-links a` | `color: var(--text)` | não declarado | não declarado |
| `.nav-cta` | `border-color: var(--acc)` | não declarado | não declarado |
| `.card` | `border-color: var(--acc)`, `transform: translateY(-3px)` | não declarado | não declarado |
| `.footer-org` | `opacity: 1` (de `0.85`) | não declarado | não declarado |
| `.lang-switch` (hover no pai revela menu) | `.lang-current { color: var(--text) }`, `.lang-menu { display: flex }` | não declarado | não declarado |
| `.lang-menu a` | `background: var(--panel)`, `color: var(--text)` | não declarado | não declarado |
| `.footer-copy a` | `color: var(--text)` | não declarado | não declarado |

**Nenhum seletor `:focus-visible` existe no arquivo inteiro** — zero estilos de foco de teclado customizados na landing. Acessibilidade de foco depende inteiramente do default do navegador.

**Breakpoints (todos, literais):**

```css
@media (max-width: 820px) { … }
@media (max-width: 480px) { … }
```

Só dois breakpoints, ambos `max-width`, mobile-first não é o padrão (é desktop-first com overrides). Não há breakpoint para tablet intermediário nem para telas grandes (>1120px o conteúdo só centraliza via `--maxw`).

**Larguras máximas de container (todas, literais):**

| Valor | Onde |
|---|---|
| `1120px` | `--maxw`, usado em `.nav`, `.section`, `.footer` |
| `900px` | `.hero` |
| `860px` | `.code-window` |
| `760px` | `.cta-final` |
| `720px` | `.section-head` |
| `640px` | `.hero-sub`, `.flow-caption` |
| `620px` | `.problem-punch` |
| `560px` | `.hero-note`, `.stats-note` |
| `520px` | `.cta-final p` |
| `460px` | `.footer p` |

Dez valores de max-width diferentes, todos hard-coded, nenhum reaproveitado como variável.

---

## 5. Componentes existentes

### `site/src/components/Landing.astro`

- **Propósito:** único componente Astro do repo. Renderiza a landing de marketing inteira (nav, hero, 4 seções de conteúdo, CTA final, footer) para um idioma.
- **Props:** `{ lang: Lang }` — um único prop obrigatório, tipo importado de `../i18n/ui`.
- **Slots:** nenhum — o componente não aceita children/slots, todo o conteúdo vem de `t(key)` (função de tradução) e de constantes internas do frontmatter (`repoUrl`, `orgUrl`, `flowSample`, ícones SVG inline como strings).
- **Onde é usado:** `site/src/pages/index.astro` (com `lang="pt"`) e `site/src/pages/[lang]/index.astro` (com `lang` dinâmico via `getStaticPaths`, para `en`/`es`/`ja`/`zh`).
- **Não tem `<style>` próprio** — todo o CSS do componente vem de `landing.css`, importado nas páginas que o envolvem, não no componente em si.
- **Estrutura interna (blocos, na ordem em que aparecem no arquivo):** `<header class="nav">`, `<section class="hero">`, `<section id="problema">`, `<section id="solucao">` (com "code window" fake), `<section id="pilares">` (grid de 4 cards com ícone SVG inline), `<section id="cli">` (stats + pills), `<section class="cta-final">`, `<footer class="footer">`.
- Ícones são todos SVG inline via `set:html` (strings de SVG cru no frontmatter, ex. `githubIcon`, `arrowIcon`) ou hardcoded diretamente no template (os 4 ícones dos cards de pilares) — **não vêm de uma lib de ícones**, são colados um a um.

### `site/src/pages/index.astro`

- **Propósito:** entrypoint da landing em português (rota `/`, locale raiz/default).
- Monta o HTML shell (`<!doctype html>`, `<head>` com meta tags e favicon, `<body>` com dois `<div>` decorativos de fundo — `.bg-stars`, `.bg-glow` — e `<Landing lang="pt" />`).
- Não é um componente reutilizável, é uma page.

### `site/src/pages/[lang]/index.astro`

- **Propósito:** mesma coisa que o `index.astro` acima, mas para os 4 locales secundários (`en`, `es`, `ja`, `zh`), usando `getStaticPaths()` para gerar uma rota estática por idioma.
- Idêntico em estrutura ao `index.astro` raiz, só muda o import relativo (`../../` em vez de `../`) e o `lang` vem de `Astro.params`.

### Componentes do Starlight usados sem customização

Nenhum componente Starlight (`Card`, `CardGrid`, `Tabs`, `TabItem`, `Aside`, `Steps`, `Code`, `LinkCard`, `FileTree`, `Badge`) é usado atualmente nos arquivos `.md` de `src/content/docs/**` — uma varredura por `<Card`, `<Tabs`, `<Aside`, `<Steps` nesses arquivos não encontra nenhuma ocorrência. As páginas de docs são Markdown puro (texto, headings, listas, tabelas, blocos de código com highlighting via Shiki default do Starlight) — nenhum componente MDX interativo foi adotado ainda, apesar de `.mdx` não estar em uso (todos os arquivos de doc são `.md`).

A página `pt/index.md` usa o `template: splash` do Starlight com `hero:` (tagline + actions) — esse é o único uso de uma feature "componentizada" do Starlight (o Hero splash), e não é customizado, é a config default do template.

---

## 6. Landing page — estrutura real

Idioma de referência: **português** (`ui.pt` em `site/src/i18n/ui.ts`), que é o locale raiz. Seção por seção, de cima para baixo, com texto EXATO.

### Header / Nav
Layout: barra horizontal, `justify-content: space-between`, logo+nome à esquerda, links + seletor de idioma + CTA à direita.

- Brand: logo (`anchors-mark-white.svg`) + texto `Anchors`
- Nav links: `O problema` · `A solução` · `Os pilares` · `O CLI`
- Seletor de idioma: mostra o idioma atual (`Português`), dropdown com os 5 idiomas
- CTA da nav: `Docs`

### Hero
Layout: centralizado, coluna única, texto centralizado.

- Logo grande (mark)
- Badge: `Em desenvolvimento ativo` (com bolinha verde pulsante)
- H1 (duas linhas, segunda com gradiente):
  > Sessões de IA esquecem.
  > **O projeto não pode.**
- Subhead:
  > Anchors é um framework de continuidade para desenvolvimento assistido por IA: âncoras que guiam o trabalho na ida, seguram a corda, e confrontam o código na volta — e não podem mentir.
- CTAs: `Ler o conceito →` (botão primário) · `⌥ Ver no GitHub` (botão ghost, com ícone GitHub)
- Nota de rodapé do hero:
  > O nome vem da **escalada**: pontos fixos que dizem para onde ir e demarcam por onde você passou.

### Seção "O problema" (`id="problema"`)
Layout: eyebrow + H2 centralizados; grid de 3 colunas (cards numerados); punchline centralizada abaixo.

- Eyebrow: `O PROBLEMA`
- H2: `Toda sessão de IA sofre de amnésia estrutural.`
- Card 01 — título: `A sessão acaba, o contexto some` — corpo: `Outro agente, outro dia, outro humano — a próxima sessão recomeça sem saber as regras, sem saber o que já foi decidido e por quê.`
- Card 02 — título: `Código cresce, direção não` — corpo: `O projeto ganha linhas, mas perde rumo. Ferramentas de agente já resolvem a amnésia dentro de uma sessão — não entre elas.`
- Card 03 — título: `Testes verdes, requisitos sem prova` — corpo: `A suíte passa enquanto pedaços inteiros da spec seguem sem cobertura nenhuma — e ninguém percebe até doer.`
- Punchline:
  > O que falta não é memória. É **artefato ancorado no repositório** — algo que a sessão futura é obrigada a respeitar.

### Seção "A solução" (`id="solucao"`)
Layout: eyebrow + H2 + lead centralizados; "code window" fake (mockup de terminal/editor) centralizada abaixo; legenda de fluxo (chips + setas) abaixo da code window; caption final.

- Eyebrow: `A SOLUÇÃO`
- H2: `Âncoras que apontam, seguram e demarcam.`
- Lead: `Cada âncora cumpre uma ou mais das três funções da escalada: diz para onde ir, segura a corda para você não cair, e marca por onde você passou. Juntas formam um grafo material, versionado — e cada aresta sabe se está em dia.`
- Code window: barra de título estilo macOS (3 bolinhas) com nome de arquivo `ciclo.sh`, conteúdo (bloco de código, texto exato incluindo comentários):
  ```
  -- planejar: a IA lê o guide de plano, decide o norte
  anchors guide plan

  @@@ spec  -- a origem da verdade: o requisito ganha um código de cenário
  SPCRX-V01: usuário autentica com e-mail e senha

  @@@ map  -- rastreabilidade: liga spec → feature → teste → código
  anchors map build

  @@@ check  -- qualidade: gates determinísticos + julgamento por IA
  anchors check --all  // diverge? vira issue, nunca se resolve sozinho
  ```
- Legenda de fluxo (chips coloridos + setas entre eles): `Planejamento semeia` → `Spec define a verdade` → `Rastreabilidade conecta` → `Propagação e Qualidade confrontam`
- Caption: `Quando algo diverge, o Anchors detecta e apresenta — nunca arbitra sozinho. Ou o trabalho está errado, ou a âncora ficou para trás.`

### Seção "Os pilares" (`id="pilares"`)
Layout: eyebrow + H2 + lead centralizados; grid 2×2 de cards com ícone.

- Eyebrow: `A TESE`
- H2: `Maturidade não é um número. É o vigor dos pilares.`
- Lead: `Um projeto é maduro quando tem *todos os seus pilares* **implementados e vigorosos** — não quando um dashboard isolado diz que sim.`
- Card 1 — título: `Duas vias, nunca se contradizem` — corpo: `O plano vivo (âncoras + grafo, a verdade atual) e o histórico (issues resolvidas, laudos, datados e imutáveis).`
- Card 2 — título: `Opt-out honesto, não silencioso` — corpo: `Dispensar uma exigência é explícito, registrado e datado. Dispensa a exigência, nunca o registro de que foi dispensada.`
- Card 3 — título: `Sincronia incremental` — corpo: `Cada aresta do grafo sabe se está em dia. Uma mudança repropaga só o que ela tocou — nunca o projeto inteiro.`
- Card 4 — título: `Issue nunca reabre` — corpo: `Quando um confronto falha e o Anchors não resolve sozinho, registra uma issue imutável. Recorrência é issue nova — o histórico nunca é reescrito.`

### Seção "O CLI" (`id="cli"`)
Layout: eyebrow + H2 centralizados; grid de 4 estatísticas; linha de pills (comandos); nota abaixo.

- Eyebrow: `O CLI`
- H2: `A ferramenta que uma IA opera para exercitar o ciclo.`
- Stats: `~2000` `nós no grafo do POC (app de referência)` · `~2800` `arestas de dependência mapeadas` · `6` `pilares implementados` · `40` `issues reais encontradas em \`check --all\``
- Pills (comandos): `init` `map build` `watch` `check` `doctor` `guide`
- Nota: `O Anchors não embute IA: é a ferramenta que a IA usa, em qualquer cliente — Claude Code, GPT, Gemini. Ela lê \`anchors guide\`, aprende o fluxo, e opera com os comandos.`

### CTA final
Layout: card centralizado com glow radial de fundo.

- H2: `Explore a documentação.`
- Corpo: `Do mecanismo comum aos seis pilares: entenda como o Anchors mantém um projeto coerente ao longo do tempo.`
- CTAs: `Ler os docs →` (primário) · `⌥ GitHub` (ghost)

### Footer
Layout: centralizado, borda superior separando do resto.

- Brand: logo pequeno + `Anchors`
- Tagline: `Um framework de continuidade para desenvolvimento assistido por IA — âncoras que não podem mentir sobre o código.`
- Linha de organização: `um projeto` + logo `co2lab`
- Linha de copyright: `© 2026 · co2lab · fonte disponível (BSL) · GitHub` (o ano é gerado via `new Date().getFullYear()`, não hard-coded)

---

## 7. Docs — estrutura

### Sidebar completa (de `site/astro.config.mjs`, na ordem declarada)

```
Conceito                      (en: Concept)
  └─ Conceito                 → /docs/conceito/

Os pilares                    (en: The pillars)
  ├─ Estrutura de Projeto     (en: Project Structure)   → /docs/estrutura/
  ├─ Planejamento             (en: Planning)            → /docs/planejamento/
  ├─ Spec                     (en: Spec)                → /docs/spec/
  ├─ Tipos de Spec            (en: Spec Types)          → /docs/tipos-de-spec/
  ├─ Rastreabilidade          (en: Traceability)        → /docs/rastreabilidade/
  ├─ Propagação               (en: Propagation)         → /docs/propagacao/
  └─ Qualidade                (en: Quality)             → /docs/qualidade/

O CLI                         (en: The CLI)
  └─ Visão geral              (en: Overview)            → /docs/cli/
```

Cada grupo do sidebar tem `items` declarados manualmente (via `slug:`), não usa `autogenerate: { directory: ... }` em nenhum grupo.

### Padrões de página

- **Frontmatter usado:** todas as páginas de conteúdo (fora do splash) usam só `title:` (ex.: `pt/conceito.md:1-3` → `---\ntitle: Conceito\n---`). A página `pt/index.md` (e presumivelmente `en/index.md`) usa o conjunto maior: `title`, `description`, `template: splash`, `hero: { tagline, actions: [...] }`.
- **Componentes recorrentes dentro do corpo Markdown:** nenhum componente Starlight — as páginas usam apenas Markdown puro: headings (`#`, `##`, `###`), blockquotes (`>`) para os avisos de abertura de cada pilar, listas, tabelas, e blocos de código com fences triplos.
- Todas as páginas de pilar abrem com um blockquote de contexto (ex.: `pt/estrutura.md` abre com `> Este documento define o **pilar de Estrutura de Projeto**...`), seguindo o mesmo padrão em todos os 7 documentos de pilar/conceito.

---

## 8. Voz e copy

### 10 trechos literais

1. (Landing, hero) — `Sessões de IA esquecem. O projeto não pode.`
2. (Landing, hero sub) — `Anchors é um framework de continuidade para desenvolvimento assistido por IA: âncoras que guiam o trabalho na ida, seguram a corda, e confrontam o código na volta — e não podem mentir.`
3. (README.md:145) — `Âncoras não podem mentir. Cada âncora é confrontada contra a realidade.`
4. (CONCEPT.md:16-17) — `Todo trabalho assistido por IA sofre de amnésia estrutural. Uma sessão de IA é sem estado entre invocações...`
5. (docs/cli.md, também `cli/cmd/anchors/root.go:10-15`) — `Anchors mantém um projeto coerente ao longo do tempo através de âncoras: documentos que guiam o desenvolvimento e confrontam o que foi feito.`
6. (CLI, `cmd/anchors/doctor.go:70`) — `✓ nenhuma ponta sistêmica encontrada — ecossistema íntegro`
7. (CLI, `cmd/anchors/check.go:307,309`) — `✓ pode promover — todos os gates bloqueantes passaram` / `✗ barrado — %d gate(s) bloqueante(s) reprovaram`
8. (CLI, `cmd/anchors/doctor.go:97`) — `(diagnóstico — nada foi bloqueado; decida o que conciliar)`
9. (CLI, `cmd/anchors/init.go:20-24`, help text) — `Escaneia o projeto de forma determinística (sem IA), propõe uma Estrutura, e confirma/ajusta com você via perguntas — chegando a um anchors.yaml correto. O grosso é inferido; as perguntas cobrem só as decisões humanas.`
10. (Landing, seção pilares) — `Opt-out honesto, não silencioso. Dispensar uma exigência é explícito, registrado e datado. Dispensa a exigência, nunca o registro de que foi dispensada.`

### Convenções

- **Pessoa:** predominantemente impessoal/terceira pessoa para descrever o sistema ("o Anchors detecta", "cada âncora cumpre"), mudando para segunda pessoa direta ("você") quando se dirige ao usuário/IA operando a ferramenta (ex.: "decida o que conciliar", "para você não cair"). Não usa "nós" para a equipe/produto — quando "nós" aparece, é dentro da metáfora de escalada ("para você não cair", genérico).
- **Caixa dos títulos:** sentence case em todos os headings (H1/H2/H3) da landing e das docs — nunca Title Case. Ex.: `Toda sessão de IA sofre de amnésia estrutural.` (não `Toda Sessão de IA Sofre de Amnésia Estrutural`). Eyebrows são a única exceção: sempre `UPPERCASE` com `letter-spacing` largo (ex.: `O PROBLEMA`, `A SOLUÇÃO`, `A TESE`).
- **Imperativo:** usado em CTAs e em nomes de comando/subcomando (`Ler o conceito`, `Ver no GitHub`, `map build`, `check --all`), não usado em títulos de seção ou corpo de texto.
- **Tamanho de frase:** frases curtas a médias predominam; parágrafos de 2-3 frases no máximo tanto na landing quanto nas docs. Uso frequente de travessão (—) para aposição/pausa retórica em vez de vírgula ou dois-pontos — é a marca de pontuação mais característica do texto (aparece em praticamente todo parágrafo de peso).
- **Emoji:** não usado em nenhum lugar do texto de produto (landing, docs, README, CLI). O único uso de símbolo "decorativo" é a âncora `⚓` no README (`# ⚓ Anchors`) e a bolinha de status (`● Em desenvolvimento ativo`), que é CSS (`.badge .dot`), não emoji.
- **Jargão:** alto e deliberado — o texto não simplifica termos técnicos para leigos ("grafo", "aresta", "stale", "gate", "carimbo" aparecem sem gloss inline na landing). A metáfora de escalada (âncora, corda, rota) é usada como o principal recurso de tornar o jargão técnico compreensível, não a simplificação lexical.

### Glossário do produto (definições como aparecem nos documentos)

| Termo | Definição, como usada nos docs |
|---|---|
| **Âncora** | "Qualquer documento que redigimos e que nos acompanha na escalada do projeto — apontando o próximo passo, segurando a corda para não cairmos, ou marcando por onde passamos." (CONCEPT.md, definição-raiz) |
| **Grafo** | O conjunto de âncoras e suas dependências, "material, versionado, de arestas tipadas (muitos-para-muitos)" — materializado no arquivo `anchors.graph.yaml` |
| **Sincronia** | Propriedade de cada aresta do grafo saber "se está em dia (carimbo por relação)" — usada pela Propagação para invalidar de forma incremental |
| **Issue** | "Quando um confronto falha e o Anchors não pode resolver sozinho, registra uma issue: arquivo em pasta-estado (`todo`/`doing`/`done`), imutável, aberta só por confronto, que nunca reabre (recorrência é issue nova)" |
| **Opt-out honesto** | "Dispensar uma exigência de forma explícita, registrada e datada (`@no-test`, `--no-block`, maturação). Dispensa a exigência, nunca o registro." |
| **Maturidade** | "Um projeto é maduro quando tem todos os seus pilares implementados e vigorosos — não é um número medido num artefato isolado, é a presença e o vigor dos pilares no projeto como um todo." (CONCEPT.md §1.1) |
| **Pilar** | Cada aplicação especializada do mecanismo comum (âncora, grafo, sincronia, issues) a um propósito nomeado: Estrutura de Projeto, Planejamento, Spec, Rastreabilidade, Propagação, Qualidade |
| **Spec** | "A origem da verdade: a âncora-base, o safepoint do qual tudo pende" (CONCEPT.md §1.1, resumo do pilar Spec) |
| **Rastreabilidade** | "A cola, em duas metades: dá a cada requisito uma identidade contínua através de suas formas (spec → feature → teste → código) e mantém o mapa de dependências entre os arquivos" |
| **Propagação** | "O motor: faz uma alteração num ponto percorrer o organismo pela Rastreabilidade, marcando o que ficou stale, até tudo voltar a ser coerente" |
| **Stale** | Estado de uma aresta do grafo que mudou e "ainda não foi reconfrontado" (`cmd/anchors/stale.go:19`) — termo em inglês usado sem tradução mesmo em texto português |
| **Amnésia estrutural** | O problema-raiz que o Anchors ataca: "Uma sessão de IA é sem estado entre invocações: a sessão acaba, o contexto some, e a próxima sessão... recomeça sem as regras" (CONCEPT.md §1) |
| **Anti-drift** | Princípio nomeado em CONCEPT.md §"Anti-drift é a lei" — não citado literalmente na landing, mas é o nome interno do princípio de "as âncoras não podem mentir" |
| **Carimbo** | O registro de timestamp/validação numa aresta que marca "está em dia" (não tem seção própria de definição, aparece operacionalmente em CONCEPT.md e QUALITY.md) |
| **Fila** (queue) | Mecanismo do watcher: "o watcher **enfileira** o trabalho, a IA **puxa** da fila (a conversa nunca fica presa)" — README.md |

---

## 9. CLI

### 9.1 Comandos e subcomandos (Short/help, literais, extraídos do código-fonte)

Comando raiz — `anchors`:
> Anchors — framework de continuidade para desenvolvimento assistido por IA

Subcomandos (nível 1, na ordem em que aparecem no `--help` real, alfabética por ser assim que Cobra lista):

| Comando | Short (descrição de uma linha) |
|---|---|
| `audit <arquivo>` | Dossiê de pendências de UM arquivo (gates + doctor), para correção em lote |
| `check` | Roda os gates de qualidade (o pipeline) |
| `code <nome>` | Gera um código de identidade único para uma nova unidade |
| `completion` | Generate the autocompletion script for the specified shell *(texto em inglês — gerado pelo Cobra, não pelo Anchors)* |
| `coverage [<spec>]` | Mostra a cobertura por cenário, por linha, do DIFF e o delta |
| `doctor` (alias `status`) | Raio-X do ecossistema: pontas sistêmicas, saúde e maturidade |
| `done <id>` | Fecha uma task reivindicada (move para o histórico `.anchors/done/`) |
| `drop <id>` | Descarta uma task da fila sem concluí-la (remove; não arquiva) |
| `governs [<guide>]` | Mostra quem cada guide rege (e quantos) — dimensiona a auditoria por guide |
| `guide` | Imprime os guias do Anchors para agentes de IA |
| `help` | Help about any command *(texto em inglês, gerado pelo Cobra)* |
| `impact <arquivo>` | Análise de impacto: o que uma alteração neste arquivo atinge |
| `ingest` | Ingere sinais de teste (JUnit/lcov) que o projeto gerou e os amarra ao mapa |
| `init` | Configura o projeto (anchors.yaml) por perguntas e respostas |
| `install-hooks` | Instala o git pre-commit que roda os gates sobre os arquivos staged |
| `judge <alvo>` | Registra o veredito de uma IA para um gate de julgamento |
| `map` | Opera o mapa de dependências (anchors.graph.yaml) |
| `new <kind> <nome>` | Emite o esqueleto de um artefato (spec/feature/test) conforme a régua |
| `next` | Puxa e reivindica o próximo item da fila (o worker chama) |
| `queue` | Lista as tasks vivas (o trabalho que o watcher enfileirou) |
| `reclaim` | Devolve à fila as tasks presas em claimed (worker morto) |
| `recode <antigo> <novo>` | Renomeia um código de identidade e propaga por todo o projeto |
| `report <perspectiva>` | Gera relatórios em docs/ por perspectiva (recortes do que o Anchors mede) |
| `stale` | Lista as arestas stale — o que mudou e ainda não foi reconfrontado |
| `watch` | O watcher em background: vê mudanças e ENFILEIRA trabalho |

Subcomandos de `map`: `show [arquivo]` ("Consulta o mapa: vizinhança de um nó, órfãos, estatísticas"), `build` ("Constrói o mapa de dependências a partir do projeto").

Subcomandos de `watch`: `start` ("Inicia o watcher em background"), `run` ("Roda o loop do watcher em foreground (uso interno do daemon)"), `status` ("Estado do watcher"), `stop` ("Encerra o watcher"), `pause` ("Pausa o watcher (sem encerrar)"), `resume` ("Retoma o watcher pausado"), `logs` ("Mostra o log do watcher").

Subcomandos de `guide`: `plan` ("Imprime o guia de plano (como um plano se estrutura e semeia specs)"), além de outros subguias referenciados no código (`spec`, `code`, `feature`, `test`, `guide`) conforme `cmd/anchors/guide_*.go`.

Subcomandos de `report`: um por perspectiva (nome dinâmico via `sp.name`/`sp.short`) mais `all` ("gera todas as perspectivas + um índice em docs/anchors/").

### 9.2 Output real do terminal

`anchors --help` (capturado rodando o binário compilado, sem cor, texto exato incluindo alinhamento por espaços):

```
Anchors mantém um projeto coerente ao longo do tempo através de âncoras:
documentos que guiam o desenvolvimento e confrontam o que foi feito.

Este CLI exercita o ciclo de vida do Anchors: constrói o mapa de dependências,
propaga alterações, roda os gates de qualidade e reporta a saúde do projeto.

Usage:
  anchors [command]

Available Commands:
  audit         Dossiê de pendências de UM arquivo (gates + doctor), para correção em lote
  check         Roda os gates de qualidade (o pipeline)
  code          Gera um código de identidade único para uma nova unidade
  completion    Generate the autocompletion script for the specified shell
  coverage      Mostra a cobertura por cenário, por linha, do DIFF e o delta
  doctor        Raio-X do ecossistema: pontas sistêmicas, saúde e maturidade
  done          Fecha uma task reivindicada (move para o histórico .anchors/done/)
  drop          Descarta uma task da fila sem concluí-la (remove; não arquiva)
  governs       Mostra quem cada guide rege (e quantos) — dimensiona a auditoria por guide
  guide         Imprime os guias do Anchors para agentes de IA
  help          Help about any command
  impact        Análise de impacto: o que uma alteração neste arquivo atinge
  ingest        Ingere sinais de teste (JUnit/lcov) que o projeto gerou e os amarra ao mapa
  init          Configura o projeto (anchors.yaml) por perguntas e respostas
  install-hooks Instala o git pre-commit que roda os gates sobre os arquivos staged
  judge         Registra o veredito de uma IA para um gate de julgamento
  map           Opera o mapa de dependências (anchors.graph.yaml)
  new           Emite o esqueleto de um artefato (spec/feature/test) conforme a régua
  next          Puxa e reivindica o próximo item da fila (o worker chama)
  queue         Lista as tasks vivas (o trabalho que o watcher enfileirou)
  reclaim       Devolve à fila as tasks presas em claimed (worker morto)
  recode        Renomeia um código de identidade e propaga por todo o projeto
  report        Gera relatórios em docs/ por perspectiva (recortes do que o Anchors mede)
  stale         Lista as arestas stale — o que mudou e ainda não foi reconfrontado
  watch         O watcher em background: vê mudanças e ENFILEIRA trabalho

Flags:
  -h, --help   help for anchors

Use "anchors [command] --help" for more information about a command.
```

Erro real (rodando `anchors map build` num diretório sem `anchors.yaml`):

```
Error: carregar anchors.yaml: open /tmp/anchors-test-proj/anchors.yaml: no such file or directory (rode `anchors init` para criar)
erro: carregar anchors.yaml: open /tmp/anchors-test-proj/anchors.yaml: no such file or directory (rode `anchors init` para criar)
```

(Nota: a mensagem aparece duas vezes — uma vez prefixada por `Error:` (formato default do Cobra) e uma vez por `erro:` em minúsculo (o próprio código do comando também imprime, redundantemente, antes do Cobra formatar o retorno de `RunE`) — ver §11.)

Exemplo de saída de `doctor` quando tudo está limpo (`cmd/anchors/doctor.go:70`, texto literal do código-fonte, não executável neste levantamento sem um `anchors.yaml`/mapa reais):
```
anchors doctor — %d nós, %d arestas, %d camadas

✓ nenhuma ponta sistêmica encontrada — ecossistema íntegro
```

Exemplo de saída de `doctor` com achados (formato, do código-fonte):
```
anchors doctor — 2142 nós, 1437 arestas, 4 camadas

⚠ órfãos (3)
  ⚠ Foo.tsx — sem spec associada
    Bar.tsx — identidade ausente

resumo: 1 ponta(s) de atenção, 3 achado(s) no total
(diagnóstico — nada foi bloqueado; decida o que conciliar)
```

Exemplo de saída de `check` (formato, do código-fonte, `cmd/anchors/check.go`):
```
  spec-completeness   bloqueante  ✓12 ✗3 ~1
  test-coverage       julgamento  ⏳2 pendente(s) de IA

~ 1 indeterminado(s) — não é falha; o gate não teve o que confrontar:
  ~ spec-completeness @ Foo.spec.md
      o código ainda não existe

3 issue(s) — divergências registradas:
  ✗ [BLOQUEIA] spec-completeness @ Foo.spec.md
      faltam 2 cenários sem teste correspondente

✗ barrado — 1 gate(s) bloqueante(s) reprovaram
```

### 9.3 Uso de cor e símbolos no terminal

- **Nenhum código ANSI de cor é emitido pelo próprio CLI do Anchors.** Uma busca por `\x1b[`, `\033`, e por bibliotecas de cor de terminal (`lipgloss`, `termenv`, pacote `color`) nos comandos e lógica interna (`cmd/anchors/*.go`, `internal/*`) não encontra nenhum uso fora das dependências de TUI interativa (`huh`, usada só em `init` para o formulário de perguntas — essa sim colorida, mas é uma lib de terceiros, não estilo autoral do Anchors).
- **Todo o output "estático" (não interativo) do CLI é texto monocromático** — a única "cor" percebida é a do terminal do usuário aplicada aos símbolos Unicode abaixo.
- **Símbolos usados como indicador visual (todos, catalogados por grep no código-fonte):**

| Símbolo | Significado | Onde |
|---|---|---|
| `✓` | sucesso / passou / resolvido | `doctor.go`, `check.go`, `install_hooks.go` |
| `✗` | falha / reprovado / bloqueado | `audit.go`, `check.go`, `install_hooks.go` |
| `✔` | corrigido (variante de check, usada em um único lugar) | `check.go:83` |
| `~` | indeterminado / skip / pendente | `audit.go`, `check.go` |
| `⏳` | pendente de julgamento por IA | `audit.go`, `check.go` |
| `⚠` | atenção / warning | `audit.go`, `doctor.go`, `install_hooks.go` |
| `ℹ` | informativo | `doctor.go` (default de `mark` antes de checar se há warning) |
| `•` | item de lista (bullet) em texto de ajuda longo | `audit.go`, `guide.go` |
| `→` | consequência/implicação em texto de ajuda | `guide_header.go`, `guide.go` |
| `──▶` | aresta do grafo, direção da dependência | `stale.go:53` (`fmt.Printf("  %s ──%s──▶ %s  (%s)\n", ...)`) |
| `──────...` (linha de 66 traços) | separador visual em texto de shell script gerado (pre-commit hook) | `install_hooks.go:146,153` |

Não há uso de caracteres de box-drawing completos (┌─┐│└┘) em lugar nenhum — só o traço simples repetido como separador, e a seta composta `──▶` usada uma vez para desenhar arestas.

---

## 10. Assets

### Todos os arquivos de imagem no repo (excluindo `node_modules/`, `dist/`, e as capturas de tela geradas durante a sessão de verificação do site — `anchors-docs-*.png`, `anchors-landing-*.png`, na raiz do repo, que não são assets de produto)

| Caminho | Dimensões | Descrição |
|---|---|---|
| `site/public/favicon.svg` | `viewBox="0 0 100 104"`, `width="32" height="32"` | Ícone de âncora simplificada (traço), cor `#ec3013` (vermelho da marca), `stroke-width="16"`, sem fundo |
| `site/public/anchors-mark-white.svg` | `viewBox="0 0 100 104"`, `100×104` | Símbolo da âncora, traço branco (`#ffffff`), `stroke-width="12"`, sem fundo (transparente) |
| `site/public/anchors-mark-black.svg` | `viewBox="0 0 100 104"`, `100×104` | Mesmo símbolo, traço preto (`#201e1d`), sem fundo |
| `site/public/anchors-mark-red.svg` | `viewBox="0 0 100 104"`, `100×104` | Mesmo símbolo, traço vermelho (`#ec3013`), sem fundo |
| `site/public/anchors-icon-red.svg` | `viewBox="0 0 512 512"`, `512×512` | Símbolo da âncora centralizado sobre fundo quadrado sólido vermelho (`#ec3013`), traço branco — formato de app icon |
| `site/public/anchors-icon-black.svg` | `viewBox="0 0 512 512"`, `512×512` | Mesmo formato de app icon, fundo quadrado sólido `#201e1d` (quase-preto), traço branco |
| `site/public/anchors-lockup-white.svg` | `viewBox="0 0 420 104"`, `420×104` | Símbolo da âncora + texto `ANCHORS` (fonte Archivo, peso 800, `letter-spacing: -1.6`), tudo branco, sem fundo |
| `site/public/anchors-lockup-black.svg` | `viewBox="0 0 420 104"`, `420×104` | Símbolo em vermelho (`#ec3013`) + texto `ANCHORS` em preto (`#201e1d`), sem fundo |
| `site/public/co2lab.png` | `300×300` px (lido do header PNG) | Logo bitmap da organização co2lab, usado no footer da landing como "selo" do projeto |
| `simulation/larder/diagram.svg` | `viewBox="0 0 1320 1000"` | Diagrama próprio (não é asset de marca) ilustrando o ciclo de vida da simulação "Larder" usada como POC didático — usa uma paleta própria não relacionada à marca (`#5a45cc` nas setas, entre outras cores não catalogadas aqui por não pertencer ao design system de produto) |

Não existe nenhum arquivo `.ico` no repo (nenhum favicon bitmap legado) — só o SVG.

**Set de ícones:** não existe uma lib de ícones instalada (nenhuma dependência tipo `lucide`, `heroicons`, `@iconify` no `package.json`). Todos os ícones usados na landing (seta, GitHub, os 4 ícones dos cards de pilares) são SVGs colados manualmente inline no `Landing.astro`, com o estilo visual (stroke, sem fill, `stroke-width="2"`, `viewBox="0 0 24 24"`) idêntico ao da biblioteca **Lucide** — mas sem a dependência declarada, sugerindo que foram copiados manualmente do site/da lib Lucide sem instalá-la como pacote.

**Kit de logo em `/temp/`:** existe uma cópia idêntica (byte a byte, mesmo conteúdo) de todos os 7 SVGs de marca dentro de `temp/` na raiz do repo (`anchors-icon-black.svg`, `anchors-icon-red.svg`, `anchors-lockup-black.svg`, `anchors-lockup-white.svg`, `anchors-mark-black.svg`, `anchors-mark-red.svg`, `anchors-mark-red.svg`, mais `favicon.svg`) — parecem ser a origem/staging de onde os arquivos foram copiados para `site/public/`.

---

## 11. Débitos e incoerências

1. **A paleta de acento da landing (teal `#2dd4bf` / azul `#38bdf8` / laranja `#fb923c`) não deriva do vermelho da marca (`#ec3013`) presente em todos os logos.** O vermelho da marca não aparece em NENHUM lugar do CSS da landing (`landing.css`) nem do componente (`Landing.astro`) — a única cor com qualquer relação de matiz é `--acc-3: #fb923c` (laranja), que é adjacente mas não é o mesmo tom. Isso significa que hoje a landing e o logo comunicam duas identidades cromáticas diferentes lado a lado (ex.: no header, a mark branca do logo aparece ao lado de um brand-name com gradiente branco→azul, nenhum dos dois puxando para o vermelho do ícone/lockup vermelho que também existe no kit).

2. **Resíduo de cor roxa do tema original (Bavel), nunca migrada.** As seguintes declarações usam `rgba(124, 92, 255, ...)` (roxo, aproximadamente `#7c5cff`) que não corresponde a NENHUMA variável declarada em `:root` — é herança do projeto irmão de onde este CSS foi copiado (ver `.bg-glow`, `.badge`, `.card-icon`, `.pill`, `.cta-final`, e o `box-shadow` de `.btn-primary`). Também sobrevive `#6d28d9` (roxo escuro) hard-coded no gradiente de `.btn-primary`, e `#c9bfff` (roxo claro) em `.pill`. Nenhuma dessas 3 cores tem relação com a paleta atual de pilares (`--structure/--planning/--spec/--traceability/--propagation/--quality`) nem com o vermelho da marca.

3. **Três das seis cores de pilar declaradas nunca são usadas no HTML.** `--planning` (`#a78bfa`), `--traceability` (`#fbbf24`) e `--propagation` (`#fb923c`, que é idêntico a `--acc-3`) têm classes CSS correspondentes (`.chip.planning`, `.chip.traceability`, `.chip.propagation`) definidas em `landing.css:463-486`, mas o único lugar em que chips de fluxo são renderizados (`Landing.astro`, seção "solução") só usa as classes `structure`, `spec`, `traceability`, `quality` — e note que a segunda ocorrência de "traceability" no template (`chip traceability` associado ao texto "Rastreabilidade conecta") na verdade usa a variável `--traceability` corretamente, mas os outros dois pilares (Planejamento, Propagação) citados na legenda de fluxo usam as classes erradas: o texto "Planejamento semeia" está com classe `structure` (azul, cor de Estrutura) em vez de `planning` (roxo), e "Propagação e Qualidade confrontam" está com classe `quality` (rosa) em vez de `propagation` (laranja) — ou seja, **a legenda de cores dos pilares na seção "A solução" não bate com o nome dos pilares que ela ilustra.**

4. **Duplicidade `--acc-3` / `--propagation`:** ambas valem exatamente `#fb923c` — são o mesmo token com dois nomes, um "semântico por papel" (`--acc-3`, terceira cor de destaque genérica) e outro "semântico por domínio" (`--propagation`, cor do pilar Propagação). Não há definição de qual dos dois é a fonte da verdade.

5. **Erro duplicado no output do CLI.** Ao rodar um comando sem `anchors.yaml` presente, a mensagem de erro é impressa DUAS vezes — uma vez com prefixo `Error:` (formato automático do Cobra ao retornar `err` de `RunE`) e uma vez com prefixo `erro:` em minúsculo (aparentemente algum wrapper de main/root imprime de novo antes ou depois do Cobra formatar). Ver saída literal capturada em §9.2. Isso é uma inconsistência de comportamento, não uma inconsistência puramente visual, mas afeta diretamente qualquer diretriz de "voz do produto no terminal" que o design system queira estabelecer (hoje a mesma frase aparece com capitalização inconsistente, uma vez ao lado da outra).

6. **Nenhuma escala de espaçamento nomeada.** Todo valor de padding/margin/gap é literal (px ou rem), sem variável — o que torna qualquer ajuste de densidade (ex.: "aumentar todo o respiro em 20%") uma operação de find-and-replace manual em ~40 valores distintos, não uma mudança de token.

7. **border-radius sem sistema.** Ver tabela em §4 — quatro valores diferentes (`24px`, `16px`, `12px`, `11px`) usados para o que visualmente são todos "cards/painéis arredondados", sem hierarquia clara (o card maior, `.cta-final`, tem o raio MAIOR `24px`, mas `.card-icon`, um elemento pequeno, tem `11px`, que é quase do mesmo raio dos cards grandes `.problem-card`/`.card`/`.stat` em `16px` — não há uma relação proporcional óbvia entre tamanho do elemento e raio).

8. **Zero tratamento de foco de teclado (`:focus-visible`) na landing inteira.** Nenhum elemento interativo (links, botões, seletor de idioma) tem estilo de foco customizado — toda a acessibilidade de navegação por teclado depende do outline default do navegador, que pode inclusive ser cortado visualmente pelo `overflow-x: hidden` do `body` ou pelos `border-radius` grandes dos botões.

9. **Light mode inexistente na landing.** A landing é hard-coded para dark (não há bloco de light theme, nem `prefers-color-scheme`), enquanto as docs (Starlight) suportam os dois temas (o Starlight tem toggle de tema nativo, visível no screenshot capturado durante a verificação: seletor "Auto" no header das docs). Um visitante que force light mode no SO verá a seção de docs mudar de tema e a landing continuar sempre escura — inconsistência de comportamento entre as duas metades do site.

10. **Ícones sem fonte declarada como dependência.** Os ícones inline (estilo Lucide) são colados como string SVG crua no componente, sem que `lucide` (ou qualquer lib) conste em `package.json` — qualquer atualização futura de um ícone exige edição manual do SVG dentro do `.astro`, não há "trocar de lib e regenerar".

11. **A fonte "Archivo" citada no lockup nunca é carregada pelo site.** O texto `ANCHORS` dentro de `anchors-lockup-*.svg` declara `font-family="Archivo, ...` mas nenhuma página do site carrega essa fonte via `@font-face`/Google Fonts — o lockup, se usado num navegador sem Archivo instalada localmente, cairá silenciosamente para Helvetica/Arial, alterando a proporção/peso visual do logotipo.

12. **Nenhum componente é compartilhado entre landing e docs.** As duas metades do site (`/` e `/docs`) não compartilham nem cor, nem tipografia, nem espaçamento, nem componente algum — são dois sistemas visuais isolados coexistindo sob o mesmo domínio, unidos apenas pelo favicon e pelo logo (light/dark) configurado no Starlight. Isso não é necessariamente um "erro", mas é a maior lacuna estrutural para quem for desenhar um design system único: hoje não existe nenhuma variável, nem token, nem componente compartilhado entre as duas superfícies.
