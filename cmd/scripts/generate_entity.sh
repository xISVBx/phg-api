#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "Uso: $0 <entity_snake_case>"
  exit 1
fi

ENTITY="$1"
ENTITY_CAMEL="$(echo "$ENTITY" | awk -F'_' '{for(i=1;i<=NF;i++)$i=toupper(substr($i,1,1)) substr($i,2)}1' OFS='')"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

write_if_missing() {
  local file="$1"
  local content="$2"
  if [[ ! -f "$file" ]]; then
    printf "%s" "$content" > "$file"
    echo "[create] $file"
  fi
}

mkdir -p "$ROOT/internal/domain/entities"
mkdir -p "$ROOT/internal/application/dtos/request/$ENTITY"
mkdir -p "$ROOT/internal/application/use_cases/$ENTITY"
mkdir -p "$ROOT/internal/web/controllers/http"
mkdir -p "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY"

write_if_missing "$ROOT/internal/domain/entities/${ENTITY}.go" "package entities

type ${ENTITY_CAMEL} struct {
	BaseEntity
	Name string
}
"

write_if_missing "$ROOT/internal/application/dtos/request/$ENTITY/Create${ENTITY_CAMEL}RequestDTO.go" "package $ENTITY

type Create${ENTITY_CAMEL}RequestDTO struct {
	Name string \\`json:\"name\" binding:\"required\"\\`
}
"

write_if_missing "$ROOT/internal/application/dtos/request/$ENTITY/Update${ENTITY_CAMEL}RequestDTO.go" "package $ENTITY

type Update${ENTITY_CAMEL}RequestDTO = Create${ENTITY_CAMEL}RequestDTO
"

write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/base.go" "package $ENTITY

import appif \"photogallery/api_go/internal/application/interfaces\"

type UseCase struct{ uow appif.IUnitOfWork }

func New(uow appif.IUnitOfWork) *UseCase { return &UseCase{uow: uow} }
"

write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/list.go" "package $ENTITY

// TODO: implementar list use case para $ENTITY.
"
write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/get.go" "package $ENTITY

// TODO: implementar get use case para $ENTITY.
"
write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/create.go" "package $ENTITY

// TODO: implementar create use case para $ENTITY.
"
write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/update.go" "package $ENTITY

// TODO: implementar update use case para $ENTITY.
"
write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/set_active.go" "package $ENTITY

// TODO: implementar activate/deactivate use case para $ENTITY.
"
write_if_missing "$ROOT/internal/application/use_cases/$ENTITY/${ENTITY}_test.go" "package ${ENTITY}_test

import \"testing\"

func Test${ENTITY_CAMEL}_TODO(t *testing.T) {
	t.Skip(\"TODO: implementar tests de integración de casos de uso\")
}
"

write_if_missing "$ROOT/internal/web/controllers/http/${ENTITY}.go" "package http

// TODO: implementar controller HTTP para entidad $ENTITY.
"
write_if_missing "$ROOT/internal/web/controllers/http/${ENTITY}_test.go" "package http_test

import \"testing\"

func Test${ENTITY_CAMEL}Controller_TODO(t *testing.T) {
	t.Skip(\"TODO: implementar tests de controller\")
}
"

write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/repository.go" "package $ENTITY

import \"gorm.io/gorm\"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
"

write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/read.go" "package $ENTITY

// TODO: consultas de lectura de $ENTITY.
"
write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/write.go" "package $ENTITY

// TODO: comandos de escritura de $ENTITY.
"
write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/statistics.go" "package $ENTITY

// TODO: métricas de $ENTITY.
"
write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/reporting.go" "package $ENTITY

// TODO: reporting de $ENTITY.
"
write_if_missing "$ROOT/internal/infrastructure/persistence/repositories/$ENTITY/repository_test.go" "package ${ENTITY}_test

import \"testing\"

func Test${ENTITY_CAMEL}Repository_TODO(t *testing.T) {
	t.Skip(\"TODO: implementar tests de integración de repositorio\")
}
"

echo "Estructura creada/actualizada para entidad: $ENTITY"
echo "Recuerda registrar contrato de repositorio, UoW, agregador de use cases y rutas."
