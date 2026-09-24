<!-- anchors:generated from doct/arquitetura.md.tmpl — inputs:0341e111430cd376 — DO NOT EDIT: run `anchors docs build` -->


# Arquitetura



## Nível 1 — Contexto

O sistema como uma caixa só: quem o usa, e com que sistemas externos ele fala.

```mermaid
graph TB
    user["Pessoa<br/><small>quem usa o sistema</small>"]
    sys["O SISTEMA<br/><small>o que este repositório constrói</small>"]
    ext["Sistema externo<br/><small>de onde vêm os dados</small>"]

    user -->|"usa"| sys
    sys -->|"consulta"| ext

    classDef pessoa fill:#08427b,stroke:#052e56,color:#fff
    classDef sistema fill:#1168bd,stroke:#0b4884,color:#fff
    classDef externo fill:#999,stroke:#6b6b6b,color:#fff
    class user pessoa
    class sys sistema
    class ext externo
```

<!-- FORA deste nível: como o sistema é montado por dentro. Isso é o nível 2. -->

## Nível 2 — Contêineres

O zoom da caixa "O SISTEMA": o que roda ou armazena separado, e **com que protocolo cada
par conversa**. O protocolo é o que diz o que acontece quando a conversa falha.


```mermaid
graph TB
    user["Pessoa"]

    subgraph sistema["O SISTEMA"]
    end



    classDef pessoa fill:#08427b,stroke:#052e56,color:#fff
    classDef conteiner fill:#438dd5,stroke:#2e6295,color:#fff
    classDef externo fill:#999,stroke:#6b6b6b,color:#fff
    class user pessoa

```

> **Nenhum contêiner declarado.** O nível 2 e os de nível 3 saem vazios até a Estrutura
> declarar o que roda separado. Ver o bloco de contêineres no `anchors.yaml`.


<!-- FORA deste nível: as peças DENTRO de cada contêiner. Isso é o nível 3 — e há um
     diagrama por contêiner, porque cada nível amplia UMA caixa do anterior. -->

## Nível 3 — Componentes

Um diagrama **por contêiner**: cada um amplia uma caixa do nível 2. Os externos não
aparecem aqui — não temos componentes dentro deles, e desenhá-los afirmaria um
conhecimento que não temos.


### Camadas fora de todo contêiner

Estas camadas existem no projeto e nenhum contêiner as declara — não aparecem em diagrama
de nível 3 nenhum. Ou falta declará-las, ou elas não rodam em lugar nenhum:

- `gate`


## Nível 4 — Código

<!-- Só onde houver algo não óbvio. O mapa que o Anchors mantém já é a versão de
     máquina, e repeti-lo aqui seria duplicação sem leitor. -->

_Nada a destacar por enquanto._
