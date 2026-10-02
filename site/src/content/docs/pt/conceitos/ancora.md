---
title: O que é uma Âncora?
description: "Entenda a metáfora da escalada, por que o Anchors existe e como as âncoras impedem que a inteligência artificial destrua seu projeto."
---

Se você já usou inteligência artificial (como ChatGPT, Claude ou Copilot) para programar em um projeto de médio ou grande porte, já deve ter presenciado este filme de terror:

1. No primeiro dia, a IA é rápida, cria arquivos funcionais e resolve o problema inicial.
2. Na semana seguinte, você pede para ela implementar uma nova funcionalidade. Ela **esquece** o padrão que adotou no dia anterior, inventa uma nova arquitetura, apaga validações antigas e quebra testes sem perceber.
3. No mês seguinte, o repositório virou um castelo de cartas: ninguém mais tem certeza do que está funcionando, quais regras existem ou por que certas decisões foram tomadas.

Esse fenômeno é chamado no Anchors de **amnésia estrutural**. A IA não tem memória persistente de longo prazo entre sessões.

A solução do Anchors para esse problema não é confiar na memória do modelo, e sim fincar **Âncoras** no próprio repositório Git.

---

## 1. A Metáfora da Escalada

O nome **Anchors** vem do montanhismo e da escalada em rocha.

Quando um alpinista sobe um paredão vertical de pedra, ele não escala no escuro sem proteção. Conforme sobe, ele crava **âncoras de aço** nas fendas da rocha e passa a corda por elas.

Uma âncora na escalada cumpre **três funções vitais**:

```
                 ▲
                 │  1. APONTA
                 │  (Mostra qual é o próximo ponto de apoio na rocha)
                 │
           [ ÂNCORA ] ─────── 2. SEGURA A CORDA (Safepoint)
                 │            (Se o alpinista escorregar, ele não despenca)
                 │
                 │  3. DEMARCA
                 ▼  (Registra o caminho percorrido para quem vem atrás)
```

No desenvolvimento de software assistido por IA, uma **âncora** funciona exatamente da mesma maneira:

1. **Aponta (Diz para onde ir)**: O documento serve de bússola. A IA lê a âncora antes de mexer no código para saber exatamente quais regras e padrões deve seguir.
2. **Segura a corda (Safepoint inegociável)**: Se a IA (ou você!) tentar commitar um código que viole as regras estabelecidas, a âncora **barra** o commit ou o merge. Ela impede a queda.
3. **Demarca (Rastro auditável)**: Qualquer desenvolvedor ou novo agente que entrar no repositório daqui a seis meses saberá exatamente quais decisões foram tomadas e por quê.

---

## 2. O que conta como uma Âncora?

A definição de âncora no Anchors é generosa: **qualquer artefato versionado no repositório que oriente o desenvolvimento e proteja o sistema contra regressões é uma âncora**.

Os principais tipos de âncoras são:

| Tipo de Âncora | Papel no Projeto | Exemplo |
| --- | --- | --- |
| **Especificações ([Specs](/pt/docs/spec//))** | A origem da verdade. Descreve em linguagem clara e estruturada as regras que o código deve cumprir. | `Login.spec.md` |
| **Cenários ([Features](/pt/docs/concepts/unidade/))** | Traduz as regras em cenários executáveis no formato Gherkin (Dado, Quando, Então). | `Login.feature` |
| **Testes Automatizados** | O código executável que comprova que o sistema se comporta como especificado. | `Login_test.go` |
| **Código Fonte** | A implementação concreta da funcionalidade. | `Login.go` |
| **Planos ([Planning](/pt/docs/planning/))** | A rota de subida: define em que ordem as tarefas e specs serão semeadas. | `PLANNING.md` |
| **Guias e Doutrina** | As regras inegociáveis de arquitetura, segurança e governança do repositório. | `QUALITY.md`, `product/*.doctrine.md` |

---

## 3. A Diferença entre Documentação Comum e uma Âncora

Muitas equipes escrevem páginas bonitas no Notion ou no Confluence que ninguém lê e que ficam obsoletas na primeira semana. **Isso NÃO é uma âncora.**

Uma documentação passiva aceita qualquer mentira: o código muda, a documentação continua dizendo o que fazia há um ano, e o projeto apodrece em silêncio.

Uma **Âncora**, por definição, **confronta a realidade através de [Gates](/pt/docs/concepts/gates-e-vereditos/) automatizados**:

```
        ÂNCORA (ex: Login.spec.md)
         ▲                     │
         │ Segura a corda      │ Aponta
         │ (Validação,         │ (A IA lê e gera
         │  check de gates)    │  o código seguindo a spec)
         │                     ▼
        CÓDIGO (ex: Login.go) ── Divergiu? ──┐
                                             │
                       ┌─────────────────────┴──────────────────┐
                       ▼                                        ▼
             Corrige o código                    Atualiza a spec e justifica
           (O gate barrou a falha)                  a mudança no histórico
```

Se a spec diz que uma senha precisa ter 8 caracteres e o código aceita 4, o Anchors acusa erro no commit. Você é forçado a escolher: ou corrige o código para respeitar a âncora, ou atualiza a âncora e justifica por que a regra de negócio mudou.

**Nenhuma âncora tem permissão para mentir sobre o código.**

---

## 4. Próximos Passos

Agora que você entende o que é uma âncora, veja como elas se conectam na prática:
- [A Unidade](/pt/docs/concepts/unidade/): A relação inseparável entre Spec, Feature, Teste e Código.
- [Rastreabilidade e Códigos de Identidade](/pt/docs/concepts/rastreabilidade-e-codigos/): Como o Anchors sabe qual teste prova qual regra.
- [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/): Os verificadores que seguram a corda.
