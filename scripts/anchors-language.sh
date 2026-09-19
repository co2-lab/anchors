#!/usr/bin/env bash
# Confronta o IDIOMA dos identificadores do código.
#
# O Anchors nasceu com o código em português e migrou para o inglês, porque o projeto vai
# ser aberto a outros devs. Esse trabalho se desfaz sozinho se nada o defender — basta um
# PR de quem não sabe da regra, e ela não está em lugar nenhum que o compilador leia.
#
# THE POLICY: everything in this repository is written in English. The single exception is
# the translation catalog (`internal/i18n/locales/*.json`), which exists to hold every
# language the product speaks.
#
# What this gate reads: DECLARATIONS (func, type, var, const).
# What it deliberately does not read, and why that is not an exemption:
#   - comments: reading prose for language is guesswork, and a false positive here costs
#     the gate its credibility. The policy still covers them — by review, not by regex.
#   - strings: the translation catalog already resolves them by the project `lang:`.
#     Confronting them here would duplicate the ruler in two places that would diverge.
set -euo pipefail

cd "$(dirname "$0")/.."
saida=$(go test ./internal/gate/ -run TestNenhumIdentificadorEmPortugues 2>&1) || {
  echo "$saida" | grep -E "identificador|TOTAL" | head -20
  exit 1
}
echo "✓ nenhum identificador em português"
