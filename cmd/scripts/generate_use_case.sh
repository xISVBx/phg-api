#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
  echo "Uso: $0 <entity_snake_case> <use_case_snake_case>"
  exit 1
fi

ENTITY="$1"
CASE="$2"
ENTITY_CAMEL="$(echo "$ENTITY" | awk -F'_' '{for(i=1;i<=NF;i++)$i=toupper(substr($i,1,1)) substr($i,2)}1' OFS='')"
CASE_CAMEL="$(echo "$CASE" | awk -F'_' '{for(i=1;i<=NF;i++)$i=toupper(substr($i,1,1)) substr($i,2)}1' OFS='')"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

write_if_missing() {
  local file="$1"
  local content="$2"
  if [[ ! -f "$file" ]]; then
    printf "%s" "$content" > "$file"
    echo "[create] $file"
  fi
}

mkdir -p "$ROOT/internal/application/use_cases/$ENTITY"
mkdir -p "$ROOT/internal/application/dtos/request/$ENTITY"
mkdir -p "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY"
mkdir -p "$ROOT/internal/web/controllers/http"

write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/base.go" "package $ENTITY

import appif \"photogallery/api_go/internal/application/interfaces\"

type UseCase struct{ uow appif.IUnitOfWork }

func New(uow appif.IUnitOfWork) *UseCase { return &UseCase{uow: uow} }
"

write_if_missing "$ROOT/internal/application/dtos/request/$ENTITY/${CASE_CAMEL}RequestDTO.go" "package $ENTITY

type ${CASE_CAMEL}RequestDTO struct {
	// TODO: campos de entrada para caso de uso $CASE
}
"

write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/${CASE}.go" "package $ENTITY

// TODO: implementar caso de uso $CASE para entidad $ENTITY.
"

write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/${CASE}_test.go" "package ${ENTITY}_test

import \"testing\"

func Test${ENTITY_CAMEL}_${CASE_CAMEL}_TODO(t *testing.T) {
	t.Skip(\"TODO: implementar tests de integración para $ENTITY/$CASE\")
}
"

write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/repository_test.go" "package ${ENTITY}_test

import \"testing\"

func Test${ENTITY_CAMEL}Repository_TODO(t *testing.T) {
	t.Skip(\"TODO: implementar tests de integración de repositorio\")
}
"

write_if_missing "$ROOT/internal/web/controllers/http/${ENTITY}_test.go" "package http_test

import \"testing\"

func Test${ENTITY_CAMEL}Controller_TODO(t *testing.T) {
	t.Skip(\"TODO: implementar tests de controller\")
}
"

echo "Caso de uso creado: $ENTITY/$CASE"
