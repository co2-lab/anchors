package governance

// projectGuide é a régua da fase DESCOBRIR — a passada que faltava: o projeto que
// ainda NÃO existe. O `anchors init` resolve a Estrutura de um projeto que já tem
// arquivos no disco (ele INFERE do que está lá); num diretório vazio ele não tem o
// que inferir, e o agente parte para o plano sem saber em que linguagem escrever.
//
// O conceito vem do Polvo (a versão IDE do Anchors), que conduz o usuário por etapas
// de perguntas e destila um resumo no fim. Aqui ele é PORTADO, não copiado: o Polvo
// tem UI e modelo por trás (nove agentes devolvendo JSON que a interface renderiza em
// formulário); o Anchors não embute IA (ver DECISIONS) — quem tem o modelo é o agente
// que chama o CLI. Então a etapa não é código Go que pergunta: é esta régua, que o
// agente lê e executa na conversa que ele já está tendo com o usuário.
//
// A divisão em DOIS arquivos é o que separa dois públicos: o PROJECT.md é lido a cada
// vez que alguém vai escrever código (curto, técnico, decidido) e o INSIGHTS.md é lido
// quando alguém pergunta "por que isto?" (longo, com as perguntas, as respostas e as
// alternativas descartadas). Misturar os dois produz um documento que ninguém relê.
const projectGuide = `# Project guide (the ruler of the DISCOVER phase)

This is the FIRST pass, and it only exists when the project DOES NOT YET EXIST — an
empty directory, or nearly. Before any plan, spec or line of code, the
project needs to answer what it is technically: in which language, under which
paradigm, with which structure, with which formatting ruler.

Why before the plan: ` + "`anchors init`" + ` INFERS the Structure from what is on disk
(the extensions, the code directories, the co-location). In an empty project there is nothing
to infer — it asks "which code directories should be treated as layers?" and the
honest answer is "none yet". Without this pass, the agent writes the first spec
without knowing whether the project is Go or TypeScript, and the decision ends up taken by accident in
the first file someone creates.

## What this phase produces

Two files at the project ROOT, and the division between them is deliberate:

- ` + "`PROJECT.md`" + ` — the TECHNICAL SUMMARY. Short, decided, with no open alternative.
  It is what you read before writing each file, so every line in it is a rule that
  the code must obey. It carries no justification.
- ` + "`INSIGHTS.md`" + ` — the TRANSCRIPT. Every question you asked, the answer the
  user gave, and what was discarded along the way. It is what answers "why Postgres and
  not Mongo?" six months later, when nobody remembers the conversation.

The rule that keeps both useful: **the decision in PROJECT, the reason in INSIGHTS.** If you
are writing "because" in PROJECT.md, the text belongs in INSIGHTS.md. If you are
writing a new rule in INSIGHTS.md that is not in PROJECT.md, it will be
ignored — nobody rereads the transcript to write code.

## Who conducts it: YOU, in the conversation

Anchors asks nothing here — it embeds no model. **You are the interviewer**,
the agent operating Anchors, in the same conversation the user is in now.
The CLI only handed you this ruler.

That has three practical consequences:

- **The interview runs in the CONVERSATION, never in a background worker.** It is the user who
  answers; delegating to a subagent that does not talk to them produces no answer at all.
  This is the exception to the rule of not monopolizing the conversation — here the conversation IS the
  work.
- **One stage at a time, waiting for the answer.** Ask the question, STOP, read what the
  user answered, and only then formulate the next. Do not write the five stages in one
  reply and ask them to answer it all together: half the value is in each
  question being born narrowed by the previous answer.
- **The files are only written AT THE END**, after the inconsistency review.
  During the interview you carry the conclusions in the conversation, not on disk. Writing
  the PROJECT.md at stage 2 and correcting as you go produces a file that has already been read wrong.

## How to conduct it (the format of each question)

You are a senior architect having a conversation, not a form. The rules that make the
difference between a useful interview and a questionnaire the user abandons:

**ONE question at a time.** Never dump the five stages at once. The user answers
the first, you understand, and the second is born narrower because of the answer.

**Every question has THREE parts:**

1. **WHAT IS AT STAKE** (2–3 sentences). The CONCRETE consequence of the answer in the
   architecture — not about the process, about the system. Name things: which database,
   which deploy model, which test pattern, which scale ceiling. Never write
   "I am going to ask you some questions", "in this stage", "this is important" — the context
   is the architecture's, not the interview's.
2. **THE QUESTION** (1 sentence). Direct, no preamble.
3. **CONCRETE EXAMPLES** (3–4 lines). Start with "Examples to guide you:" and give
   one realistic scenario per line, at the level of detail you expect in the answer.
   A vague example produces a vague answer.

**BE OPINIONATED.** This is the rule that changes the outcome most. The user is rarely an
expert in everything you are going to ask — if you present a neutral menu of
eight options, they pick the one they recognize, not the one that fits. Use what you ALREADY KNOW from the
previous stages to eliminate what does not fit and present at most 2–3
viable options, **leading with your recommendation and its reason**. "Given that it is a team of
1 person with a 2-month deadline, I recommend a modular monolith because X — does that work, or
do you have a strong reason for something else?" is a better question than "monolith,
microservices or serverless?".

**If the user asks for help or says they do not know:** present the concrete tradeoffs
in THEIR context, and ask a narrower question that leads them to decide. Do not jump
to the next stage with the matter open.

## The stages

Five stages, in this order. The order matters: each one restricts the options of the next,
and inverting them produces a tool choice before knowing what for.

### Stage 1 — Purpose and form
What the system is, before any technology. Is it an API, a CLI, a web app, mobile,
a library, a worker, a monorepo with several of these? Who consumes it? That decides whether there is
an interface layer, whether there is a distribution build, whether there is session state.
Conclude with: what it is, who consumes it, which executable artifacts are born.

### Stage 2 — Language and runtime
Only now, and restricted by stage 1. Language, version, package manager,
execution runtime. Be opinionated based on the form and the team's experience —
a CLI distributed as a single binary pushes towards Go/Rust; a web app with a team
that already knows TypeScript rarely justifies switching.
Conclude with: language + version, manager, runtime.

### Stage 3 — Architecture and paradigm
The structural pattern and the paradigm: layered, clean/hexagonal, feature-sliced, MVC?
Object-oriented, functional, procedural? Vertical modules per feature or horizontal
layers? Where does the domain live?
This is what becomes the ` + "`layers:`" + ` of anchors.yaml — the answer here is not academic,
it decides the glob of each layer. Check the presets: ` + "`anchors init`" + ` offers
established structures per stack, and matching a preset saves work and error.
Conclude with: pattern, paradigm, organization (modular or layered), boundaries.

### Stage 4 — Macro structure and file conventions
The first-level directories and what lives in each. The extensions of each type of
file. The test naming convention (` + "`*_test.go`" + `, ` + "`*.test.ts`" + `, ` + "`*.spec.ts`" + `) and whether
test/spec sit BESIDE the code (co-location) or in a separate tree — ` + "`init`" + `
asks exactly that, and the answer here is the one it will use.
Conclude with: directory tree, extensions, test pattern, co-location yes/no.

### Stage 5 — Tooling and formatting
The mechanical ruler: indentation (tabs or spaces, how many), line width, formatter
(gofmt, prettier, black, rustfmt), linter and its configuration, naming convention
(camelCase, snake_case, PascalCase per kind of symbol). Editors and extensions the
team uses, and whether there is an ` + "`.editorconfig`" + `, ` + "`.vscode/`" + ` or equivalent to version.
This looks like a detail and is not: it is what makes the code YOU write indistinguishable
from what the team writes. Without this declared, every new file negotiates the style again.
Conclude with: indentation, formatter, linter, naming convention, editors/extensions.

### Closing — the inconsistency review
Before writing the files, reread the five conclusions looking for CONTRADICTION.
They are common and expensive:
  • a functional paradigm declared + a structure that only makes sense with classes
  • co-location "yes" + a macro structure with a separate ` + "`tests/`" + ` at the top
  • a language without static typing + a gate that requires types
  • a modular-by-feature pattern + directories organized by technical layer
Present each contradiction found with the TWO choices that collide, quoted
verbatim, and ask the user to resolve it. One at a time. Only write the files
when none is left.

## Writing the PROJECT.md

At the root. Only what was DECIDED — no "probably", no "to be defined". A field that
was not decided does not go in: an empty field in a document read before writing
code is worse than absence, because it looks like a decision.

` + "```md" + `
# Project: <name>

<one sentence: what the system is and for whom>

## Stack
| item | choice |
| --- | --- |
| language | Go 1.23 |
| manager | go mod |
| runtime | native binary |
| frameworks | cobra (CLI), huh (TUI) |

## Architecture
- **Pattern:** layered
- **Paradigm:** procedural with types; no inheritance
- **Organization:** horizontal layers (cmd/ → internal/)
- **Boundaries:** cmd/ only orchestrates I/O; the decision lives in internal/

## Macro structure
` + "```" + `
cmd/anchors/     commands (a thin shell per command)
internal/        the logic, one package per domain
guides/          the rulers of the artifacts
` + "```" + `

## File conventions
| item | choice |
| --- | --- |
| extensions | .go |
| test | ` + "`*_test.go`" + `, beside the code |
| co-location | yes |

## Formatting
| item | choice |
| --- | --- |
| indentation | tabs |
| formatter | gofmt |
| linter | go vet |
| names | PascalCase exported, camelCase internal |

## Tooling
- editors: VS Code, Neovim
- extensions: gopls
- versioned: .editorconfig
` + "```" + `

## Writing the INSIGHTS.md

Also at the root. One section per stage, and inside it each question with the answer
the user gave. What makes this file worth the disk is the third part: **what was
discarded and why.** A decision without the rejected alternatives cannot be revised —
whoever reopens it does not know whether the option they are proposing was already considered.

` + "```md" + `
# Project insights

> Transcript of the DISCOVER phase. The DECISIONS live in PROJECT.md; here are the
> questions, the answers and what was discarded.

## Stage 2 — Language and runtime

**Q:** <the question as you asked it>
**A:** <the user's answer, in their own words>

**Decided:** Go 1.23
**Discarded:**
- Rust — the learning curve does not pay off in a team that already ships in Go
- Node/TS — distribution as a single binary was a requirement from stage 1
` + "```" + `

## After writing

1. ` + "`anchors init`" + ` — now it has something to infer. The answers of stages 3, 4 and 5
   are exactly what it asks (preset, layers, co-location, test pattern);
   answer with what is in PROJECT.md, without renegotiating.
2. ` + "`anchors map build`" + ` — the two new files enter the map.
3. Move on to PLANNING (` + "`anchors guide plan`" + `): the first plan decides which specs
   are born, and now it is born knowing in which language they will be implemented.

## When the project ALREADY exists

Do not conduct the whole interview — the disk has already answered most of it. Read what
is there (extensions, directories, ` + "`.editorconfig`" + `, linter config, package
manifest), write the PROJECT.md with what you OBSERVED, and ask the user only
what the code does not reveal: the intended paradigm, the boundaries that should exist,
and what in today's code is debt rather than pattern. The INSIGHTS.md records what was
observed versus what was decided — the difference between the two is the list of debts.

## What to NEVER do

- **Writing the PROJECT.md without talking.** A technical summary you invented is a
  guess with the authority of a document. If there was no interview, there is no PROJECT.md.
- **Leaving a field open.** "To be defined" in PROJECT.md becomes a decision taken by
  accident in the first file someone writes.
- **Repeating the INSIGHTS inside the PROJECT.** The PROJECT.md stops being read the day
  it gets long, and it is the file that needs to be read always.
- **Letting both age in silence.** When a decision changes (a linter swap, a new module,
  a paradigm change), update the PROJECT.md and record in the
  INSIGHTS.md what changed and why. A PROJECT.md that lies about the project is worse
  than not existing — the agent obeys it.
`
