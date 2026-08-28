#!/usr/bin/env bash
# Gate `padrao-inseguro` — semgrep (LGPL-2.1).
#
# Usa o ruleset `p/gosec` (a régua consagrada de segurança em Go) em vez de regras
# próprias. O gate equivalente num projeto Node carrega regras escritas à mão para
# padrões que AQUELA base provou recorrentes; aqui ainda não há esse corpus, e inventar
# regra sem caso medido produz falso positivo que ninguém confia.
#
# INFORMATIVO: um achado do gosec é uma pergunta ("isto é seguro no seu contexto?"), não
# um veredito. Bloquear no primeiro dia treina a equipe a ignorar — o oposto do que o
# gate quer.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)" || exit 1

command -v semgrep >/dev/null 2>&1 || {
  echo "semgrep não instalado — \`brew install semgrep\`. Gate PULADO (não aprovado)."
  exit 0
}

# Sem argumento, varre o projeto inteiro (é o modo `scope_full`); com, só os arquivos
# passados pelo motor.
ALVO=("$@")
[ ${#ALVO[@]} -eq 0 ] && ALVO=("cli/")

TMP=$(mktemp)
semgrep --config=p/gosec --json --quiet "${ALVO[@]}" > "$TMP" 2>/dev/null

python3 - "$TMP" <<'PY'
import json, sys, collections
try:
    d = json.load(open(sys.argv[1]))
except Exception:
    print("semgrep não produziu saída legível — gate PULADO")
    raise SystemExit(0)

achados = d.get('results', [])
if not achados:
    print("nenhum padrão inseguro conhecido (p/gosec)")
    raise SystemExit(0)

por_regra = collections.Counter()
exemplos = {}
for a in achados:
    regra = a['check_id'].split('.')[-1]
    por_regra[regra] += 1
    exemplos.setdefault(regra, f"{a['path']}:{a['start']['line']}")

print(f"{len(achados)} achado(s) do gosec em {len({a['path'] for a in achados})} arquivo(s):")
for regra, n in por_regra.most_common(10):
    print(f"  {n:3}x {regra}  (ex.: {exemplos[regra]})")
print("")
print("Informativo: cada achado é uma PERGUNTA sobre o contexto, não um veredito.")
print("O que for deliberado, marque com `// nosemgrep: <regra>` e o motivo.")
PY
rm -f "$TMP"
exit 0
