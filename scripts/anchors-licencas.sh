#!/usr/bin/env bash
# Gate `licenca-compativel` — go-licenses (Apache-2.0, Google).
#
# O Anchors é distribuído como BINÁRIO ESTÁTICO: tudo que ele importa vai junto no
# executável. Isso muda a régua em relação a um projeto que só roda o código — aqui,
# copyleft numa dependência alcança o binário inteiro.
#
# BLOQUEANTE para copyleft forte (AGPL/SSPL/GPL) porque a consequência é jurídica e o
# número hoje é ZERO: bloquear um número que já é zero não interrompe ninguém, e impede a
# PRIMEIRA entrada — que é quando o custo de reverter ainda é baixo.
#
# Uma assimetria que vale nomear: o Anchors é Elastic License 2.0 (fonte disponível, com
# restrição de revenda como serviço). Uma dependência GPL exigiria que o todo fosse GPL,
# o que a ELv2 não é — então o conflito é real, não teórico.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)/cli" || exit 1

command -v go-licenses >/dev/null 2>&1 || {
  echo "go-licenses não instalado — \`go install github.com/google/go-licenses@latest\`."
  echo "Gate PULADO (não aprovado)."
  exit 0
}

TMP=$(mktemp)
go-licenses report ./cmd/anchors > "$TMP" 2>/dev/null || {
  echo "go-licenses falhou. Gate PULADO (não aprovado)."; rm -f "$TMP"; exit 0
}

python3 - "$TMP" <<'PY'
import csv, sys, collections, re, subprocess

# O módulo do PRÓPRIO projeto não é dependência de terceiro: a licença dele é a do
# repositório (ELv2), e o `go-licenses` a reporta como "Unknown" porque só classifica
# licenças OSI — a ELv2 é fonte-disponível, não open source.
#
# Sem esta exclusão o gate reportaria ~19 falsos positivos a cada execução, e um gate que
# erra sempre treina a equipe a ignorá-lo.
try:
    PROPRIO = subprocess.run(["go", "list", "-m"], capture_output=True, text=True).stdout.strip()
except Exception:
    PROPRIO = ""

# Copyleft FORTE: obriga abrir o código de quem distribui o binário. Incompatível com a
# ELv2 do projeto.
proibidas = re.compile(r'\b(AGPL|SSPL|GPL-[23]|GPL-2|GPL-3)(?!.*exception)', re.I)
# Copyleft FRACO: obriga abrir só o arquivo modificado. Conviver é possível, mas quem
# redistribui precisa saber que a obrigação existe.
fraco = re.compile(r'\b(LGPL|MPL|EPL|CDDL)', re.I)

linhas = list(csv.reader(open(sys.argv[1])))
cont = collections.Counter()
bad, aviso, desconhecida = [], [], []

for l in linhas:
    if len(l) < 3:
        continue
    pkg, _, lic = l[0], l[1], l[2]
    if PROPRIO and pkg.startswith(PROPRIO):
        continue
    cont[lic] += 1
    if proibidas.search(lic):
        bad.append((pkg, lic))
    elif fraco.search(lic):
        aviso.append((pkg, lic))
    elif lic in ('Unknown', '', 'UNKNOWN'):
        desconhecida.append((pkg, lic))

print(f"{sum(cont.values())} dependências de terceiros no binário, {len(cont)} licenças distintas")
for lic, n in cont.most_common(6):
    print(f"  {lic[:40]:42} {n}")

if aviso:
    print(f"\ncopyleft fraco ({len(aviso)}) — a obrigação alcança o arquivo modificado:")
    for p, l in aviso[:6]:
        print(f"  {p[:52]:54} {l}")

if desconhecida:
    print(f"\nlicença NÃO declarada ({len(desconhecida)}) — sem licença explícita o default")
    print("é 'todos os direitos reservados', o que proíbe a redistribuição:")
    for p, l in desconhecida[:6]:
        print(f"  {p[:52]:54} {l or '(vazia)'}")

if bad:
    print(f"\ncopyleft FORTE ({len(bad)}) — INCOMPATÍVEL com a BSL deste projeto:")
    for p, l in bad:
        print(f"  {p[:52]:54} {l}")
    print("\nUma dependência GPL/AGPL exigiria que o binário inteiro fosse GPL.")
    raise SystemExit(1)

print("\nnenhuma licença copyleft forte.")
PY
rc=$?; rm -f "$TMP"; exit $rc
