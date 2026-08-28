#!/usr/bin/env bash
# Gate `secret-nao-vazado` — gitleaks (MIT).
#
# Segredo vazado é o único achado desta lista que NÃO tem volta: rotacionar chave é caro e
# o histórico do git guarda a original. Por isso é o único que nasce BLOQUEANTE.
#
# No pre-commit varre só o STAGED (0,7s medido); no CI varre o histórico (29s medido) —
# daí o gate declarar `cost: fast` e o script escolher o modo pelo contexto.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)" || exit 1

command -v gitleaks >/dev/null 2>&1 || {
  echo "gitleaks não instalado — \`brew install gitleaks\`. Gate PULADO (não aprovado)."
  exit 0
}

if [ "${ANCHORS_PHASE:-pre-commit}" = "ci" ]; then
  gitleaks detect --no-banner --redact --exit-code 1 && { echo "sem segredos no histórico"; exit 0; }
else
  gitleaks protect --staged --no-banner --redact --exit-code 1 && { echo "sem segredos no staged"; exit 0; }
fi
echo ""
echo "Segredo detectado. Se for falso positivo, declare em .gitleaks.toml (allowlist) com o"
echo "MOTIVO — a allowlist é a declaração de que alguém olhou e aprovou."
exit 1
