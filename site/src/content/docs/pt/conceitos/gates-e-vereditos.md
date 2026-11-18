---
title: Gates e Vereditos
description: "O que é um gate, como funcionam as verificações automáticas, a diferença entre bloqueante e informativo, e os quatro vereditos possíveis."
---

Se as [Âncoras](/pt/docs/concepts/ancora/) são os pontos fixos que seguram a corda do projeto, os **Gates (Portais de Qualidade)** são os verificadores que inspecionam se a corda está bem amarrada.

Um **gate** é uma pergunta objetiva que o projeto faz a si mesmo em momentos-chave (como antes de commitar ou no CI), e cuja resposta o Anchors confronta de forma rigorosa.

---

## 1. A Estrutura de um Gate

Um gate é declarado no arquivo [`anchors.yaml`](/pt/docs/anchors-yaml//). Veja um exemplo típico:

```yaml
gates:
  - name: unit-complete          # O nome do gate como aparece no terminal
    on: [spec]                    # Em quais tipos de arquivo ele roda
    check: unit-complete         # O verificador interno executado
    blocking: true                # Reprovar BARRA o commit ou apenas avisa?
    when: [pre-commit, ci]        # Quando ele deve ser avaliado
    measures: "a spec tem código, feature e teste que a realizam"
```

### Os Campos Principais:
- **`on`**: Os tipos de arquivo alvo (`spec`, `feature`, `test`, `code`, `doc`, `plan`, `flag`).
- **`check`**: O nome de um verificador **interno determinístico** embutido no binário do Anchors.
- **`run`**: Um comando de ferramenta **externa** do seu ecossistema (ex: `gitleaks`, `osv-scanner`, `npm audit`).
- **`ask`**: A pergunta para um gate de **[Julgamento por IA](/pt/docs/concepts/julgamento-ia/)**.
- **`blocking`**: Define se a falha é um erro fatal que cancela a operação (`true`) ou apenas um aviso educativo (`false`).
- **`when`**: As fases do ciclo em que o gate roda (`pre-commit`, `pre-push`, `ci`).

---

## 2. Informativo vs Bloqueante: A Regra de Ouro

Uma das maiores causas de fracasso ao adotar ferramentas de qualidade em equipes é ligar 50 regras rigorosas de uma vez. No primeiro dia, tudo quebra, os desenvolvedores ficam frustrados e desligam a ferramenta.

O Anchors adota uma doutrina clara:

> [!IMPORTANT]
> **Todo gate novo deve nascer INFORMATIVO (`blocking: false`).**

O ciclo de adoção saudável funciona em três passos:
1. **Ligue como informativo (`blocking: false`)**: O gate roda, gera relatórios e você descobre quantos arquivos violam a regra sem travar o trabalho de ninguém.
2. **Estabilize e corrija**: A equipe ou a IA ajusta os pontos apontados.
3. **Promova para bloqueante (`blocking: true`)**: Quando a métrica estiver limpa, você torna o gate bloqueante para garantir que ninguém mais cometa aquele erro no futuro.

---

## 3. Os Quatro Desfechos (Vereditos)

Quando um gate é executado sobre um arquivo, o resultado não é apenas "passou" ou "falhou". Existem **quatro desfechos possíveis**:

| Símbolo | Veredito | Significado | Exemplo Prático |
| :---: | :--- | --- | --- |
| `✓` | **Pass (Aprovado)** | O arquivo cumpriu todos os requisitos da verificação. | A spec possui código, feature e teste associados. |
| `✗` | **Fail (Reprovado)** | O arquivo violou a regra. Se o gate for `blocking: true`, o commit ou PR é barrado. | A spec cita 5 regras, mas a feature só cobre 2. |
| `~` | **Indeterminado / Pending** | O gate **não teve o que confrontar** ou os dados de teste ainda não foram ingeridos. | Um gate de mutação em um arquivo que ainda não rodou a ferramenta de mutação. |
| `Skip` | **Dispensado** | A regra não se aplica honestamente a este tipo de arquivo ou camada. | Um gate de especificação rodando sobre uma [Camada Reconhecida](/pt/docs/layers/) (declarativa). |

### Por que o `~` (Indeterminado) NÃO é uma falha?
Este é um detalhe crucial para não se confundir: **Indeterminado não significa erro**.

Se um gate de cobertura de teste encontra um arquivo que não declara teste, ele não pode afirmar que a cobertura é zero nem que é cem: ele simplesmente não tem o que medir.

Confundir *falta de medição* com *reprovação* faz as pessoas "consertarem" o que não está quebrado e polui o histórico com hacks desnecessários.

---

## 4. Onde Ver os Gates

Você pode rodar todos os gates configurados no seu repositório a qualquer momento com:

```sh
anchors check
```

E para verificar a saúde geral do ecossistema e se todas as ferramentas necessárias estão instaladas:

```sh
anchors doctor
```

---

## 5. Próximos Passos

- [Catálogo Completo de Gates](/pt/docs/gates//): Conheça todos os gates disponíveis no Anchors.
- [Julgamento por IA](/pt/docs/concepts/julgamento-ia/): Entenda os gates que fazem perguntas semânticas.
- [Maturidade e Saúde](/pt/docs/concepts/maturidade-e-saude/): Como os gates compõem a qualidade geral do projeto.
