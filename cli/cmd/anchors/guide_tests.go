package main

// testGuide é a régua universal do artefato TEST: a prova executável dos cenários.
// Agnóstico de framework (Jest, pytest, go test, …).
const testGuide = `# Guia de teste (a régua executável)

O teste é a prova executável dos cenários da feature. Ele fecha a cadeia de
rastreabilidade: o mesmo código atravessa spec → feature → teste. Um teste que passa
por engano é pior que a ausência de teste — ele mente sobre a saúde do sistema.

## A pirâmide

Organize os testes por custo e escopo, do barato ao caro:
- UNIDADE — lógica pura, sem ambiente. Rápido, roda em qualquer lugar.
- INTEGRAÇÃO — um módulo com suas bordas (infraestrutura) mockadas; o domínio roda
  de verdade.
- PONTA-A-PONTA — o fluxo real atravessando as fronteiras do sistema.
Cada cenário da feature já veio classificado no seu nível natural — prove-o ali.

## As regras de ouro (o que separa teste real de teatro)

- COBRE O FELIZ + AO MENOS UM ERRO. Todo comportamento testado prova o caminho
  principal E pelo menos um caso de falha.
- TESTE O CÓDIGO DE PRODUÇÃO REAL. Proibido mock inline que REPLICA a lógica dentro
  do teste — isso não exercita o código real e gera zero cobertura. Proibido gravar e
  ler pela mesma via de baixo nível no mesmo teste — isso testa a biblioteca, não o
  seu código.
- MOCKE A INFRAESTRUTURA NA BORDA, não as unidades internas. Deixe o domínio executar
  de verdade; mocke só o que sai do processo (rede, disco, serviço externo). Mockar a
  unidade interna inteira zera a cobertura que importa.
- ISOLAMENTO E LIMPEZA. Dados de teste levam prefixo identificável; limpe no fim (e
  tenha uma rede de segurança global). Nunca deixe dado de teste permanente.
- CENTRALIZE O SETUP. Não crie dados dentro de cada caso quando há estado
  compartilhado — centralize o seed; evita corrida e contenção.
- ASSÍNCRONO/CONSISTÊNCIA EVENTUAL EXPLÍCITOS. Se o sistema é eventualmente
  consistente, espere por um predicado (retry/poll), não leia logo após escrever.

## Segurança de ambiente (quando o teste toca estado externo)

- DESCUBRA O AMBIENTE POR IDENTIDADE, NÃO POR NOME. Valide que os recursos pertencem
  ao ambiente de teste por uma marca inequívoca (tag/atributo), nunca inferindo pelo
  nome. Falhe ANTES de gravar se houver ambiguidade — para nunca rodar contra produção.
- NÃO EDITE O ARQUIVO DE AMBIENTE À MÃO. Gere-o por script, reproduzível.

## Honestidade sobre o verde (regra cultural)

Distinga as suítes que rodam sem credencial das que exigem ambiente externo. NUNCA
afirme "N/N verde" sem essa ressalva. "Passou o que dava para passar aqui" é uma
frase honesta; "tudo verde" quando metade nem rodou é uma mentira que custa caro
depois. (Vale a regra de ouro do projeto: ausência de prova não é prova de ausência.)

## Anti-padrões (recuse-os)

- Mock que replica a lógica do alvo → você testou o mock, não o código.
- Gravar e ler pela mesma via de baixo nível → você testou a biblioteca.
- Mockar a unidade interna inteira → cobertura zero no que importa.
- Ler logo após escrever num sistema eventual → teste intermitente (flaky).
- Inferir ambiente por nome → risco de rodar contra produção.
- "Tudo verde" sem dizer o que não rodou → relatório desonesto.

## Especialização do projeto

O framework concreto, os helpers, a estrutura de pastas e os gates de PR são do seu
projeto — veja o guide de teste na camada 'guide' do anchors.yaml. Este guia é a
doutrina universal; siga o dialeto do projeto quando existir, e avise se não existir.
`
