# AGENTS.md — Reglas Operativas para Agentes de Código (Proyecto API Go)

## 1) Objetivo de este archivo
Este archivo define cómo debe trabajar un agente (Codex) en este repositorio para implementar cambios consistentes con la arquitectura Dilogum, el alcance MVP y la convención técnica del proyecto.

## 2) Arquitectura obligatoria
Usar arquitectura por capas con separación estricta:
- `domain`: entidades, enums y contratos de repositorio.
- `application`: casos de uso, DTOs, puertos/interfaces, mappers.
- `infrastructure`: persistencia (GORM/DB), UoW, providers, servicios externos.
- `web`: rutas, controllers, middlewares y utilidades HTTP.

Regla general:
- No mezclar responsabilidades entre capas.
- Mantener dependencias dirigidas hacia adentro (`web -> application -> domain`; `infrastructure` implementa puertos).

## 3) Estructura de carpetas y archivos
Mantener esta granularidad mínima:
- `internal/domain/entities/`: 1 archivo por entidad.
- `internal/domain/enums/`: 1 archivo por enum (o grupo pequeño coherente).
- `internal/domain/repositories/`: 1 archivo por interfaz de repositorio.
- `internal/application/dtos/request/<entity>/`: 1 archivo por DTO.
- `internal/application/use_cases/<entity>/`:
  - `base.go` (UseCase + constructor)
  - 1 archivo por caso de uso (`create.go`, `update.go`, `delete.go`, `list.go`, `get.go`, etc.).
- `internal/infrastructure/persistence/repositories/<entity>/`:
  - `repository.go`
  - `read.go`
  - `write.go`
  - `statistics.go`
  - `reporting.go`
- `internal/web/controllers/http/`: controller por entidad/recurso (no por caso de uso).

## 4) Reglas por capa
### Domain
- Entidades y reglas de dominio puras.
- Sin dependencias de framework HTTP o transporte.
- Contratos de repositorio por interfaz, sin implementación.

### Application
- Los casos de uso orquestan negocio y transacciones.
- No acceder DB directa ni usar GORM aquí.
- Usar DTOs de request/response y mappers.
- Si hay múltiples repos/entidades en la misma operación: usar UoW + transacción.

### Infrastructure
- Implementaciones reales de repositorios, UoW, servicios externos.
- División obligatoria `read/write/statistics/reporting` por entidad.
- Notificaciones en MVP: interfaz + implementación Noop + TODO para integración real.

### Web
- Controllers delgados: parseo/validación, llamada a caso de uso, respuesta.
- Sin lógica de negocio en controllers.
- Registrar rutas en archivos de rutas por módulo.

## 5) Convenciones de naming
- Entidades: `PascalCase` (`User`, `WorkOrder`, `CashMovement`).
- DTOs: `CreateXRequestDTO`, `UpdateXRequestDTO`, `ListXQueryDTO`.
- Casos de uso: métodos explícitos (`CreateCustomer`, `RegisterPayment`, etc.).
- Archivos de use case por acción: `create_*.go`, `update_*.go`, `get_*.go`, `list_*.go`, `delete_*.go`.
- Errores de negocio: códigos claros (`BAD_REQUEST`, `NOT_FOUND`, `FORBIDDEN`, etc.).

## 6) Reglas de endpoints
- Base URL: `/api/v1`.
- Auth: JWT Bearer (`Authorization: Bearer <token>`).
- Formato de respuesta:
  - éxito: `{ "data": ..., "meta": ... }`
  - error: `application/problem+json` (RFC 7807) con:
    - estándar: `type`, `title`, `status`, `detail`, `instance`
    - extensiones del proyecto: `code`, `errors`
- No inventar endpoints fuera del documento funcional del proyecto.
- Mantener Swagger sincronizado:
  - todo endpoint HTTP debe quedar documentado
  - todo endpoint debe tener comentarios `swag` completos: `@Summary`, `@Description`, `@Tags`, `@Accept` (si aplica), `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`
  - si no se agregan esos comentarios, el endpoint no se ve/describe correctamente en Swagger y Scalar
  - las anotaciones pueden vivir en el handler o en stubs dedicados (ej. `swagger_routes.go`), pero siempre deben mapear 1:1 con rutas reales
  - regenerar con `make swag`
  - exponer `/swagger/index.html` (y mantener `/scalar` si está activo).

## 7) Módulos MVP obligatorios
El agente debe priorizar cobertura funcional MVP para:
- Auth (login/refresh/me/change-password/my-permissions)
- Seguridad RBAC (users/roles/menus/submenus/permissions/overrides)
- Catálogo (categories/products)
- Clientes
- Ventas + pagos/abonos + creación automática de work order
- Órdenes de trabajo
- Agenda (appointments)
- Archivos (upload/download/delete + links + auditoría)
- Caja (categorías/movimientos)
- Trabajadores (CRUD + pago comisiones FIFO + pago sueldo)
- Auditoría y AppSettings

WebSockets:
- En MVP no se exponen eventos WS. Mantener estructura preparada sin implementación funcional obligatoria.

## 8) Reglas críticas de archivos (storage)
Config relevante:
- `FILES_BASE_PATH`
- `FILES_MAX_SIZE_MB`
- `FILES_ALLOWED_MIME`
- `CUSTOMER_KEY_MODE` (`email|code|both`)

Reglas de path físico:
- Si involucra cliente (Customer/Sale/WorkOrder/SaleItem/Appointment):
  - `clientes/{customerKey}/ventas/{saleId}/...`
  - `clientes/{customerKey}/ordenes/{workOrderId}/...`
  - `clientes/{customerKey}/citas/{appointmentId}/...`
  - `clientes/{customerKey}/general/...`
- Si es interno empresa:
  - `empresa/interno/{YYYY}/{MM}/...`

`customerKey`:
- Preferir email sanitizado (si existe y es usable).
- Fallback: `CustomerCode` único.

Persistencia mínima en BD:
- `StorageRelativePath`
- `StorageKind`
- relación en `FileLinks` (+ `CustomerId` cuando aplique)
- auditoría de `UPLOAD`, `DELETE`, `LINK` (y `MOVE` si existe)

Validaciones de seguridad:
- tamaño máximo
- MIME permitido
- sanitización de filename
- prevenir path traversal

## 9) Seguridad y auditoría
- Implementar RBAC con permisos efectivos (rol base + overrides Grant/Revoke).
- Auditar operaciones críticas (descuentos, anulaciones, pagos, caja, cambios de estado, archivos, seguridad).
- No hardcodear secretos; usar configuración/env.

## 10) Reglas de testing
Estrategia recomendada en este proyecto:
- Priorizar tests de integración para casos de uso y repositorios (DB real de pruebas).
- Colocar tests junto al código del módulo:
  - repos: `internal/infrastructure/persistence/repositories/<entity>/*_test.go`
  - use cases: `internal/application/use_cases/<entity>/*_test.go`
  - controllers: `internal/web/controllers/http/*_test.go`
- Organizar tests por feature/use case (no fragmentados sin criterio).
- Usar `DATABASE_DSN_TEST` para integración.
- Ejecutar suite y verificar resultados antes de cerrar cambios.

## 11) Checklist operativo para cambios
Antes de cerrar cualquier feature endpoint:
1. DTOs por entidad y por archivo creados/ajustados.
2. Caso(s) de uso en archivos separados por acción.
3. Contrato(s) de repositorio por interfaz.
4. Implementación repos en `read/write/statistics/reporting`.
5. Controller delgado actualizado.
6. Ruta registrada con middleware adecuado.
7. Swagger actualizado (`make swag`).
   - verificar que `docs/swagger/swagger.json` y `swagger.yaml` incluyan request/response/error correctos por endpoint.
8. Tests ejecutados y consistentes.
9. Documentación funcional mínima actualizada si aplica.

## 12) Restricciones de implementación
- No mover lógica de negocio a controller.
- No romper separación de capas.
- No introducir endpoints fuera de alcance funcional definido.
- No implementar notificaciones reales en MVP (solo contrato + Noop + TODO).
