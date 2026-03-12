#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="${1:-.cache/test-all.log}"
mkdir -p "$(dirname "$LOG_FILE")"

: "${DATABASE_DSN_TEST:?DATABASE_DSN_TEST no definido}"

ESC=$'\033'
RESET="${ESC}[0m"
BOLD="${ESC}[1m"
RED="${ESC}[31m"
GREEN="${ESC}[32m"
BLUE="${ESC}[34m"
YELLOW="${ESC}[33m"

export GIN_MODE=release

echo -e ">> Ejecutando tests (pretty runner)..."
set +e
GOCACHE="${GOCACHE:-}" GOMODCACHE="${GOMODCACHE:-}" DATABASE_DSN_TEST="$DATABASE_DSN_TEST" go test ./... -p 1 -count=1 -v 2>&1 \
  | tee "$LOG_FILE" \
  | awk '
BEGIN {
  esc = sprintf("%c", 27)
  reset = esc "[0m"
  bold = esc "[1m"
  red = esc "[31m"
  green = esc "[32m"
  blue = esc "[34m"
  yellow = esc "[33m"
  total = 0
  passed = 0
  failed = 0
}
/^=== RUN   / {
  name = $0
  sub(/^=== RUN   /, "", name)
  printf("%sRUN%s      %s\n", blue, reset, name)
  total++
  next
}
/^--- PASS: / {
  name = $0
  sub(/^--- PASS: /, "", name)
  sub(/ .*/, "", name)
  printf("%sSUCCESS%s  %s\n", green, reset, name)
  passed++
  next
}
/^--- FAIL: / {
  name = $0
  sub(/^--- FAIL: /, "", name)
  sub(/ .*/, "", name)
  printf("%sFAIL%s     %s\n", red, reset, name)
  failed++
  next
}
/^ok[[:space:]]/ {
  printf("%sPKG OK%s   %s\n", green, reset, $2)
  next
}
/^FAIL[[:space:]]/ {
  printf("%sPKG FAIL%s %s\n", red, reset, $2)
  next
}
END {
  printf("\n%sTOTAL:%s  %s%d%s  %sPASSED:%s  %s%d%s  %sFAILED:%s  %s%d%s\n", \
    bold, reset, blue, total, reset, \
    bold, reset, green, passed, reset, \
    bold, reset, red, failed, reset)
}
'
go_status=${PIPESTATUS[0]}
set -e

if [[ $go_status -ne 0 ]]; then
  echo
  echo -e "${RED}${BOLD}Failure details:${RESET}"
  awk '
    { lines[NR] = $0 }
    /^--- FAIL: / { fail[++n] = NR }
    END {
      if (n == 0) {
        print "(No se detectaron líneas --- FAIL:, revisar log completo)"
      }
      for (i = 1; i <= n; i++) {
        s = fail[i] - 5
        if (s < 1) s = 1
        e = fail[i] + 6
        for (j = s; j <= e; j++) {
          if (j in lines) print lines[j]
        }
        print ""
      }
    }
  ' "$LOG_FILE"
  echo
  echo -e "${YELLOW}Log completo:${RESET} $LOG_FILE"
  exit 1
fi

echo -e "${GREEN}${BOLD}>> ✅ Todos los tests pasaron.${RESET}"
