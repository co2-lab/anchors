#!/usr/bin/env bash
# Confronta o IDIOMA dos identificadores do código.
#
# O Anchors nasceu com o código em português e migrou para o inglês, porque o projeto vai
# ser aberto a outros devs. Esse trabalho se desfaz sozinho se nada o defender — basta um
# PR de quem não sabe da regra, e ela não está em lugar nenhum que o compilador leia.
#
# O que este gate olha: DECLARAÇÕES (func, type, var, const).
# O que ele deliberadamente ignora:
#   - comentários, que carregam as medições e os porquês, no idioma do time
#   - strings, que vão pelo i18n — traduzi-las aqui seria o erro oposto
set -euo pipefail

cd "$(dirname "$0")/.."
saida=$(go test ./internal/gate/ -run TestNenhumIdentificadorEmPortugues 2>&1) || {
  echo "$saida" | grep -E "identificador|TOTAL" | head -20
  exit 1
}
echo "✓ nenhum identificador em português"
