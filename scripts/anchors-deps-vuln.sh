#!/usr/bin/env bash
# Gate `dependencia-vulneravel` — osv-scanner (Apache-2.0, Google).
#
# Lê o `go.sum` direto: determinístico, e a única rede é a base OSV. Adaptado do gate
# equivalente de um projeto Node — lá o alvo era `yarn.lock`, aqui é o módulo Go.
#
# INFORMATIVO por decisão, e o motivo é o mesmo do original: a maioria das CVEs de um
# projeto pequeno é TRANSITIVA de toolchain, e bloquear commit por CVE de dependência de
# teste trava o trabalho sem reduzir risco. A métrica evolui; quando estiver baixa e
# estável, vira bloqueante.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)" || exit 1

command -v osv-scanner >/dev/null 2>&1 || {
  echo "osv-scanner não instalado — \`brew install osv-scanner\`. Gate PULADO (não aprovado)."
  exit 0
}

TMP=$(mktemp)
osv-scanner scan source --lockfile cli/go.mod --format json > "$TMP" 2>/dev/null
python3 - "$TMP" <<'PY'
import json, sys, collections
try:
    d = json.load(open(sys.argv[1]))
except Exception:
    print("osv-scanner não produziu saída legível — gate PULADO")
    raise SystemExit(0)

sev = collections.Counter()
pacotes = collections.defaultdict(set)
for r in d.get('results', []):
    for p in r.get('packages', []):
        nome = p['package']['name']
        for v in p.get('vulnerabilities', []):
            s = (v.get('database_specific') or {}).get('severity', 'UNKNOWN')
            sev[s] += 1
            pacotes[s].add(nome)

tot = sum(sev.values())
if tot == 0:
    print("nenhuma vulnerabilidade conhecida em cli/go.mod")
    raise SystemExit(0)

ordem = ['CRITICAL', 'HIGH', 'MODERATE', 'LOW', 'UNKNOWN']
print(f"{tot} vulnerabilidades em {len({n for s in pacotes.values() for n in s})} pacotes:")
for s in ordem:
    if sev[s]:
        print(f"  {s:9} {sev[s]:4}  ({', '.join(sorted(pacotes[s])[:6])}{'…' if len(pacotes[s]) > 6 else ''})")
print("")
print("Informativo: confira se a CVE atinge caminho de EXECUÇÃO ou só toolchain —")
print("`go mod why <pacote>` mostra quem a puxa.")
PY
rm -f "$TMP"
exit 0
