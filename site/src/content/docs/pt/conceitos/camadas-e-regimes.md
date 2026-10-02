---
title: Camadas Arquiteturais e Regimes
description: "A planta da casa do seu projeto, a divisão entre camadas regidas e reconhecidas, e os regimes de teste."
---

Todo projeto de software precisa de uma **planta** para não virar um labirinto caótico. No Anchors, chamamos essa planta de **Camadas (Layers)**, e os diferentes modos de confrontar a realidade de **Regimes de Teste**.

> [!NOTE]
> Para a documentação exaustiva com exemplos de configuração e tabelas de compatibilidade, consulte a página [Camadas do Projeto](/pt/docs/layers/).

---

## 1. As Camadas: A Planta da Casa

Uma **camada** é um agrupamento de arquivos que possuem a mesma responsabilidade arquitetural no sistema.

No Anchors, você declara suas camadas no arquivo [`anchors.yaml`](/pt/docs/anchors-yaml//). A distinção mais fundamental é entre:

### A. Camadas Regidas (Regime de Regra)
- Onde moram as **regras de negócio** e comportamentos do sistema.
- Exigem a [Unidade](/pt/docs/concepts/unidade/) completa: Spec + Feature + Teste + Código.
- Exemplos: `usecase`, `service`, `domain`, `comando`.

### B. Camadas Reconhecidas (Regime Declarativo)
- Onde moram os **adaptadores técnicos e infraestrutura** que não inventam regras de negócio próprias.
- Dispensam specs completas e cenários em Gherkin. A sua identificação no cabeçalho `@anchors` é simplesmente o nome da sua camada (`layer: infra` ou `layer: dao`).
- Exemplos: `infra`, `dao`, `types`, `doc`.

---

## 2. Os Regimes de Teste

Nem todo teste deve ser testado da mesma forma. Você não testa uma regra de cálculo financeiro abrindo um navegador com Selenium, e não testa a cor de um botão com um teste unitário de CPU.

O Anchors reconhece **quatro regimes canônicos de teste**:

```
┌────────────────────────────────────────────────────────┐
│                   REGIMES DE TESTE                     │
└──────────────────────────────────┬─────────────────────┘
                                   │
      ┌───────────────┬────────────┴───┬───────────────┐
      ▼               ▼                ▼               ▼
┌───────────┐   ┌───────────┐    ┌───────────┐   ┌───────────┐
│   UNIT    │   │INTEGRATION│    │    E2E    │   │    VR     │
├───────────┤   ├───────────┤    ├───────────┤   ├───────────┤
│ Regras    │   │ Banco,    │    │ Fluxo     │   │ Regressão │
│ puras em  │   │ filas,    │    │ completo  │   │ visual,   │
│ memória   │   │ rede      │    │ da ponta  │   │ baseline  │
│ e funções │   │ externa   │    │ ao fim    │   │ de telas  │
└───────────┘   └───────────┘    └───────────┘   └───────────┘
```

| Regime | O que valida | Onde a prova mora (Superfície) |
| --- | --- | --- |
| **`unit`** | Funções puras, cálculos, validações de domínio isoladas de I/O. | Arquivo de teste unitário ao lado do código (`Login_test.go`). |
| **`integration`** | A integração do código com banco de dados, filas ou adaptadores reais. | Pasta de testes de integração (`tests/integration/`). |
| **`e2e`** | O fluxo de ponta a ponta navegando pela aplicação completa. | Roteiros executáveis (`.yaml`, Playwright, Cypress). |
| **`vr`** | A fidelidade visual e ausência de distorção de layout. | Imagens de baseline (`.png`). |

---

## 3. O Roteamento de Cenários por Regime

Uma mesma [Feature](/pt/docs/concepts/unidade/) pode conter cenários em regimes diferentes. Por exemplo:
- Cenários com tag `@unit` são confrontados contra o arquivo de teste unitário.
- Cenários com tag `@vr` são confrontados contra as imagens de baseline do gate [`vr-baseline`](/pt/docs/gates//vr-baseline/).

O Anchors sabe rotear cada cenário para a **superfície correta**, impedindo que um teste unitário falhe por falta de assert visual, ou que uma tela seja aprovada sem teste de regressão.

---

## 4. Próximos Passos

- [Guia Completo de Camadas](/pt/docs/layers/): Configuração detalhada no `anchors.yaml`.
- [Gate layer-boundary](/pt/docs/gates//layer-boundary/): Como barrar importações proibidas entre camadas.
- [Gate test-level-codes](/pt/docs/gates//test-level-codes/): Como garantir que testes unitários e testes visuais não misturem códigos indevidos.
