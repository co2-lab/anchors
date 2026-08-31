package main

// reviewGuide é a régua de quem REVISA um PR. Ele existe porque o `claim` entrega cards
// em `ready-to-review` e não dizia o que conferir — o revisor caía no PR sem saber o que
// é dele e o que já foi medido por script.
//
// A divisão é a que separa as duas coisas: o que se confronta por SCRIPT roda como check
// do PR (mapa em dia, gates, cards declarados); o que exige JULGAMENTO é do revisor. Um
// review que repete o que o check já mediu gasta o revisor no que a máquina faz melhor —
// e deixa passar justamente o que só ele veria.
const reviewGuide = `# Guia de review (a régua de quem revisa um PR)

## O que NÃO é seu

Os checks do PR já confrontaram, por script:

- o mapa commitado corresponde ao repositório;
- os gates bloqueantes passam (` + "`anchors check --all`" + `);
- os cards que o PR fecha estão declarados no corpo;
- os pipelines estão atualizados.

Se algum reprovou, o PR volta para quem o abriu — não é trabalho de review.
Reconferir isso à mão gasta você no que a máquina faz melhor.

## O que é seu

O que exige JULGAMENTO, e por isso nenhum script alcança:

### A spec decide o que precisava decidir?

Um gate confere que a spec tem as seções e que as regras têm código. Nenhum confere
se a regra RESPONDE a pergunta que o código vai fazer. Uma spec que diz "o sistema
deve ser rápido" passa em todos os gates e não decide nada — e quem implementa vai
inventar o limiar.

### O código realiza a regra, ou só a cita?

O gate ` + "`regra-cumprida`" + ` já pergunta isso a uma IA, e o veredito dela está no
mapa. Seu trabalho é diferente: confrontar o JULGAMENTO com o trecho. Um veredito
"pass" sobre marcação num lugar genérico (topo do arquivo, import) é o defeito que a
automação erra com mais frequência.

### O teste PROVA, ou só executa?

Cobertura diz que a linha rodou. Um teste que chama a função e não afirma nada sobre
o resultado dá 100% de cobertura e não prova coisa alguma. Leia as asserções, não o
número.

### O que o PR mudou sem dizer?

Diff que o autor não menciona na descrição é onde mora o que ele não percebeu que
mudou. Compare a descrição com o que de fato foi tocado.

### A mudança de spec muda a DIREÇÃO?

Se o PR altera uma spec ou plano, a revisão (` + "`{CODIGO}-R000N`" + `) diz o que
mudou e por quê. A pergunta é se aquilo era correção de FORMA — e se não for, se
alguém decidiu. O gate confere que a revisão EXISTE; se ela é honesta é você quem vê.

## Ao terminar

Aprovar não é "não achei nada": é afirmar que você OLHOU o que é seu. Se olhou e não
há o que dizer, aprove — a revisão que só passa o olho é pior que nenhuma, porque
cria a impressão de que alguém conferiu.

Achou algo que não é deste PR? Registre em vez de só comentar:

    anchors escalate "<o que está errado>" --sobre <arquivo> --card <o card do PR>
`
