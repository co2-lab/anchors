package config

// A TABELA DE ALIAS que este arquivo testava FOI REMOVIDA.
//
// Ela aceitava os nomes de gate em português e os convertia na carga, para sempre. Os
// testes daqui provavam que a conversão funcionava — e ela funcionava; o problema era o
// desenho: o arquivo nunca se consertava, e o mapa acumulava carimbos nos dois formatos.
// Medido no blue-eyes: 40 julgamentos gravados como `regra-cumprida` convivendo com 2 como
// `rule-fulfilled`.
//
// Havia também uma assimetria: a LEITURA normalizava (`mapx.mesmoGate`), a ESCRITA não —
// um projeto que renomeasse o gate ganhava um SEGUNDO carimbo em vez de atualizar o
// primeiro.
//
// A conversão virou o passo de migração `1→2` (`internal/migra/formato_2.go`), e as réguas
// que valiam a pena foram com ela:
//
//   · `internal/migra` testa a conversão em si (renomeia chave e valor, não toca em
//     comentário, é idempotente, respeita o arquivo de destino);
//   · `internal/initx/vocabulario_test.go` confere que todo destino do de-para corresponde
//     a um gate que EXISTE, e que nenhum gate default nasceu com nome em português.
//
// Este arquivo fica como registro: quem procurar pelos testes do alias encontra por que
// eles não estão mais aqui.
