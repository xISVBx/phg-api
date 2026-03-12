# Normas Para Agente IA: API Go (Arquitectura Dilogum)

## 1. Objetivo

Este documento define el estándar obligatorio para construir APIs Go en este proyecto:

- Arquitectura por capas (`domain`, `application`, `infrastructure`, `web`).
- Unit of Work (UoW) para transacciones multi-repositorio.
- Casos de uso como orquestadores de negocio.
- Repositorios por entidad con separación interna por tipo de consulta/comando.
- DTOs y mappers separados del dominio.
- Swagger/OpenAPI obligatorio.

## 2. Estructura base obligatoria

```text
cmd/
  api/main.go

docs/
  swagger/
  uses_cases/
  websockets/

internal/
  domain/
    entities/
      <entity>.go                      # 1 archivo por entidad
    enums/
      <enum_or_group>.go               # 1 archivo por enum/grupo
    repositories/
      repositories.go                  # contratos de repositorio

  application/
    dtos/
      request/
        <entity>/
          <DTO>.go                     # 1 archivo por DTO
      responses/
      queries/
    interfaces/
      services/
      websockets/
      uow.go
    mappers/
    use_cases/
      common/
      <entity>/
        base.go                        # type UseCase + constructor
        <case>.go                      # 1 archivo por caso de uso
      use_cases.go                     # agregador central

  infrastructure/
    config/
    persistence/
      database.go
      repositories/
        repositories.go                # compat/adaptadores generales (si aplica)
        <entity>/
          repository.go                # struct repository + constructor
          read.go                      # consultas de lectura
          write.go                     # comandos de escritura
          statistics.go                # métricas/agregados
          reporting.go                 # reportes
      scripts/
      migration/
    providers/
    services/
    uow/
    websocket/
    cron/

  web/
    routes.go
    web.go
    controllers/http/
      <entity>.go                      # controller por entidad/recurso
    controllers/ws/
    handlers/
    middlewares/
    utils/
```

## 3. Reglas de arquitectura por capa

### 3.1 Domain

- `entities`: solo dominio, sin detalles HTTP/DB/framework.
- `repositories`: contratos (interfaces), no implementación.
- Prohibido importar `gin`, `gorm`, `jwt`, `http`.

### 3.2 Application

- `use_cases/<entity>/<case>.go`: lógica de negocio por caso de uso.
- Un caso de uso por archivo.
- Solo depende de `domain`, `interfaces`, `dtos`, `mappers`.
- No usar acceso directo a DB.

### 3.3 Infrastructure

- Implementa contratos de repositorio y puertos externos.
- Repos por entidad con separación obligatoria:
  - `read.go`
  - `write.go`
  - `statistics.go`
  - `reporting.go`
- `uow` crea repos con `db` o `tx`.

### 3.4 Web (Delivery)

- Controllers delgados.
- Controller por entidad/recurso (`controllers/http/<entity>.go`).
- Sin lógica de negocio en controller.

## 4. Patrón Unit of Work (obligatorio)

Contrato esperado:

- `Repositories()`: acceso a set de repositorios.
- `Transaction(ctx, fn)`: operaciones atómicas multi-entidad.

Reglas:

- Si un caso de uso toca más de una entidad, usar `Transaction`.
- Efectos externos (archivos, notificaciones) preferiblemente post-commit o con patrón robusto.

## 5. Flujo estándar request -> response

1. `routes.go` define endpoint y middlewares.
2. Controller parsea request.
3. Controller llama caso de uso.
4. Use case valida reglas y orquesta repos/servicios.
5. Repositorio persiste/consulta.
6. Use case retorna resultado.
7. Controller responde en formato estándar.

## 6. Convenciones de código

- Un archivo por entidad en `domain/entities`.
- Un archivo por enum en `domain/enums`.
- Un archivo por caso de uso en `application/use_cases/<entity>/`.
- Un archivo por DTO en `application/dtos/request/<entity>/`.
- Un archivo de controller por entidad en `web/controllers/http/`.
- Nombres sugeridos:
  - `CreateXRequestDTO`, `UpdateXRequestDTO`, `ListXQueryDTO`.
  - `create.go`, `update.go`, `delete.go`, `list.go`, `get.go`.
- Errores de negocio: `BadRequest`, `NotFound`, `Unauthorized`, `Forbidden`.
- Respuesta HTTP de error obligatoria: `application/problem+json` (RFC 7807) con:
  - `type`, `title`, `status`, `detail`, `instance`
  - extensiones del proyecto: `code`, `errors`

## 7. Plantilla mínima para una entidad nueva

Para una entidad `Product`:

1. `internal/domain/entities/product.go`
2. `internal/domain/repositories/repositories.go` (extender contrato)
3. `internal/infrastructure/persistence/repositories/products/repository.go`
4. `internal/infrastructure/persistence/repositories/products/read.go`
5. `internal/infrastructure/persistence/repositories/products/write.go`
6. `internal/infrastructure/persistence/repositories/products/statistics.go`
7. `internal/infrastructure/persistence/repositories/products/reporting.go`
8. `internal/application/dtos/request/product/CreateProductRequestDTO.go` (y los demás DTOs)
9. `internal/application/use_cases/product/base.go`
10. `internal/application/use_cases/product/create.go` (y demás casos)
11. `internal/web/controllers/http/product.go`
12. Registrar en:
   - `application/use_cases/use_cases.go`
   - `infrastructure/uow/uow.go`
   - `web/routes.go`

## 8. Swagger/OpenAPI (obligatorio)

- Todo endpoint HTTP debe tener anotaciones `swag`.
- Regenerar docs con `make swag`.
- Exponer UI en `/swagger/index.html`.

## 9. Scripts y scaffolding

Si se usan scripts de scaffolding, deben generar la estructura nueva por entidad/caso:

- `use_cases/<entity>/base.go` + archivos por caso.
- `dtos/request/<entity>/` con 1 archivo por DTO.
- `repositories/<entity>/{repository,read,write,statistics,reporting}.go`.
- `controllers/http/<entity>.go`.

## 10. Checklist obligatorio al crear endpoint nuevo

1. Crear/ajustar DTO(s) en carpeta de entidad.
2. Crear/ajustar caso de uso (archivo propio).
3. Crear/ajustar contrato de repositorio.
4. Implementar read/write/statistics/reporting según corresponda.
5. Controller por entidad, método delgado.
6. Registrar ruta y middlewares.
7. Agregar anotaciones Swagger.
8. Ejecutar `make swag`.
9. Ejecutar `go test ./...`.
10. Documentar en `docs/uses_cases/...`.

## 11. Reglas para agente IA

- No mezclar lógica de negocio en controllers.
- No acceder DB desde use case sin repositorio.
- Toda operación multi-entidad va en `uow.Transaction`.
- Respetar separación por archivos y carpetas de esta norma.
- No exponer secretos en código.
- Mantener Swagger sincronizado con rutas reales.
- Si se implementa WS, documentar eventos en `docs/websockets`.

## 12. Recomendaciones de mejora

- Tests por capa (use case/repository/controller).
- Migraciones versionadas para producción.
- CI con lint + tests.
- Outbox para efectos externos robustos.
