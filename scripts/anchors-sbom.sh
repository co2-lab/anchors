#!/usr/bin/env bash
# Gate `sbom-gerado` — syft (Apache-2.0).
#
# O SBOM é ARTEFATO de auditoria, não uma regra que possa reprovar código — daí ser
# informativo. Ganha importância com a abertura do projeto: quem distribui binário passa a
# receber a pergunta "o que tem aqui dentro?", e um inventário gerado responde melhor do
# que uma lista mantida à mão.
#
# Só em CI: varrer a árvore inteira custa tempo que não pertence a um pre-commit.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)" || exit 1

command -v syft >/dev/null 2>&1 || {
  echo "syft não instalado — \`brew install syft\`. Gate PULADO (não aprovado)."
  exit 0
}

mkdir -p reports
OUT=reports/sbom-cyclonedx.json
syft scan dir:cli --output "cyclonedx-json=$OUT" -q 2>/dev/null || {
  echo "syft falhou. Gate PULADO (não aprovado)."
  exit 0
}

python3 - "$OUT" <<'PY'
import json, sys, collections
d = json.load(open(sys.argv[1]))
comps = d.get('components', [])
tipos = collections.Counter(c.get('type', '?') for c in comps)
print(f"SBOM CycloneDX gerado: {len(comps)} componentes")
for t, n in tipos.most_common(6):
    print(f"  {t:14} {n}")
print(f"  → {sys.argv[1]}")
PY
exit 0
