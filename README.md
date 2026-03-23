# PhotoGallery API Go

Backend MVP en Go + Gin + PostgreSQL + GORM con arquitectura por capas tipo Dilogum:
- `domain` entidades + interfaces de repositorio.
- `application` DTOs, mappers, puertos y use cases.
- `infrastructure` GORM repos, UoW, JWT/bcrypt, servicios externos.
- `web` controllers HTTP, middlewares y rutas.

## Cómo leer la estructura del proyecto
La estructura está pensada para que cada capa tenga una responsabilidad clara y para que los casos de uso no dependan de GORM ni de Gin directamente.

### `internal/domain`
Contiene el núcleo del negocio:
- `entities/`: estructuras persistidas y conceptos del dominio.
- `repositories/`: contratos que la aplicación usa para leer/escribir datos.
- `enums/` y `value_objects/`: tipos de apoyo del dominio.

Regla práctica:
- aquí no debería entrar lógica HTTP, GORM, JWT, storage ni detalles de framework.

### `internal/application`
Contiene la lógica de aplicación:
- `dtos/`: payloads de entrada y salida.
- `interfaces/`: puertos que necesita la aplicación, por ejemplo `IUnitOfWork`.
- `use_cases/`: operaciones del sistema agrupadas por módulo.

Acá vive la orquestación:
- validar/parsing de IDs y fechas
- pedir datos a repositorios
- ejecutar reglas de negocio
- decidir si una operación requiere transacción
- registrar auditoría

### `internal/infrastructure`
Contiene implementaciones concretas:
- `persistence/repositories/`: implementación GORM de cada repositorio.
- `uow/`: `UnitOfWork` que construye el `RepositorySet` y abre transacciones.
- `services/`: JWT, hash de passwords, file storage, notifier, etc.
- `bootstrap/` y `seed/`: inicialización y datos de desarrollo.

### `internal/web`
Contiene la capa HTTP:
- `controllers/http/`: handlers Gin.
- `middlewares/`: auth, permisos, request id, logger, etc.
- `routes*.go`: registro de rutas.
- `utils/`: helpers de respuesta y extracción de datos del request.

Flujo típico:
`route -> controller -> use case -> repository/UoW -> DB`

## Guía rápida para trabajar en el código
Si vas a agregar o modificar una funcionalidad, normalmente el recorrido es este:

1. Definir o ajustar DTOs en `internal/application/dtos`.
2. Implementar la lógica en `internal/application/use_cases/<modulo>`.
3. Usar contratos de `internal/domain/repositories`.
4. Si hace falta persistencia nueva, implementar el método en `internal/infrastructure/persistence/repositories/<modulo>`.
5. Exponerlo desde controller y ruta en `internal/web`.

Ejemplo real:
- el controller obtiene el usuario autenticado con `actorID(c)` y llama al caso de uso.
- el caso de uso decide si usa `uow.Repositories()` o `uow.Transaction(...)`.
- los repositorios concretos reciben `context.Context` y operan con la instancia GORM correcta.

## Actor: qué es y cómo se usa
En este proyecto `actor` es el `uuid.UUID` del usuario autenticado que está ejecutando la acción.

Normalmente aparece así en los use cases:
- `Create(ctx context.Context, actor uuid.UUID, in ...)`
- `Update(ctx context.Context, actor, id uuid.UUID, in ...)`

Su propósito principal es:
- guardar `CreatedByUserID` o `UpdatedByUserID` cuando la entidad lo necesita
- registrar auditoría en `AuditLogs`
- saber quién disparó operaciones sensibles como pagos, cambios de permisos o settings

En la capa HTTP sale de la sesión autenticada:
- `actorID(c)` usa `webutils.MustUserID(c)` y devuelve el UUID del usuario actual

Regla práctica:
- si una acción depende del usuario logueado o debe auditarse, el use case debería recibir `actor`
- si es una lectura pública o una operación interna sin usuario, no siempre hace falta

## Unit of Work (`uow`) y `RepositorySet`
La aplicación no crea repositorios uno por uno. Todo entra por `IUnitOfWork`:

```go
type IUnitOfWork interface {
	Repositories() RepositorySet
	Transaction(ctx context.Context, fn func(repos RepositorySet) error) error
}
```

`RepositorySet` es un agrupador de repositorios del dominio:
- `Users()`
- `Roles()`
- `Customers()`
- `Sales()`
- `WorkOrders()`
- etc.

La implementación concreta está en `internal/infrastructure/uow/uow.go`.

### Cuándo usar `uow.Repositories()`
Usa `uow.Repositories()` cuando la operación es simple y no necesita atomicidad entre varias escrituras.

Ejemplo típico:
- crear una sola entidad
- actualizar una sola tabla
- listar o consultar datos

Patrón real del proyecto:

```go
err := u.uow.Repositories().Customers().Create(ctx, item)
if err == nil {
	common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Customer", item.ID.String(), "CREATE", item)
}
```

Esto funciona bien para operaciones pequeñas, pero hay que usarlo con criterio:
- cada llamada obtiene repositorios montados sobre `u.db`
- si necesitas que varias escrituras vivan o fallen juntas, esto no alcanza

### Cuándo usar `uow.Transaction(...)`
Usa `uow.Transaction(...)` cuando una operación:
- escribe en varias tablas
- necesita consistencia total
- crea auditoría que debe quedar en la misma transacción
- hace lecturas y escrituras que deben ver el mismo contexto transaccional

Patrón real del proyecto:

```go
err := u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
	if err := repos.Sales().Create(ctx, sale); err != nil {
		return err
	}
	if err := repos.SaleItems().CreateMany(ctx, saleItems); err != nil {
		return err
	}
	common.CreateAudit(ctx, repos, &actor, "Sale", sale.ID.String(), "CREATE", sale)
	return nil
})
```

Eso hace que todo use el mismo `tx *gorm.DB`:
- venta
- items
- work order auto-creada
- auditoría

Si algo falla, todo hace rollback.

### Regla importante dentro de una transacción
Si ya entraste a `uow.Transaction(...)`, usa el `repos` que te entrega el callback.

Correcto:
- `repos.Users().Create(...)`
- `common.CreateAudit(ctx, repos, ...)`

Incorrecto dentro de la transacción:
- `u.uow.Repositories().Users().Create(...)`
- `common.CreateAudit(ctx, u.uow.Repositories(), ...)`

La razón es simple:
- `repos` está ligado al `tx` transaccional
- `u.uow.Repositories()` crea repositorios sobre la conexión base, fuera de ese `tx`

## Cómo están organizados los repositorios
Cada repositorio concreto vive en:
- `internal/infrastructure/persistence/repositories/<modulo>/`

La estructura usual es:
- `repository.go`: struct base con `db *gorm.DB`
- `read.go`: consultas
- `write.go`: escrituras
- `reporting.go`: reportes o queries especiales
- `statistics.go`: agregaciones o métricas

Ejemplo con `customers`:
- `repository.go` define `type Repository struct { db *gorm.DB }`
- `read.go` implementa `List`, `GetByID`
- `write.go` implementa `Create`, `Update`, `SetActive`

Esta separación ayuda a que:
- el archivo no crezca sin control
- sea claro dónde vive cada método
- sea más fácil mantener queries complejas

## Cómo decidir dónde poner la lógica
Regla práctica por capa:

### Controller
Debe hacer solo trabajo HTTP:
- bind de JSON o query params
- parseo de path params
- devolver status codes y responses
- obtener `actorID(c)`

No debería:
- hablar con GORM directamente
- meter reglas de negocio
- coordinar varias tablas

### Use case
Es el lugar correcto para:
- validar reglas del negocio
- coordinar varios repositorios
- abrir transacciones
- disparar auditoría
- llamar servicios externos vía interfaces

### Repository
Debe encargarse de:
- queries a DB
- filtros, paginación, `Preload`, agregaciones
- persistencia de entidades

No debería:
- conocer Gin
- decidir respuestas HTTP
- contener orquestación de negocio de alto nivel

## Auditoría: patrón esperado
La auditoría se centraliza con `common.CreateAudit(...)`.

Se usa para registrar:
- actor
- entidad afectada
- acción realizada
- payload serializado en JSON

Recomendación:
- si la operación es transaccional, registra la auditoría con el `repos` recibido en la transacción
- si la operación es simple, se puede registrar después de la escritura exitosa

## Convenciones útiles al crear un nuevo use case
- El constructor del módulo normalmente recibe solo las dependencias necesarias.
- La dependencia principal casi siempre es `appif.IUnitOfWork`.
- Si el caso de uso requiere password hashing, JWT, files o notificaciones, eso entra como interfaz adicional.

Patrón común:

```go
type UseCase struct{
	uow appif.IUnitOfWork
}

func New(uow appif.IUnitOfWork) *UseCase {
	return &UseCase{uow: uow}
}
```

Si el caso de uso hace varias escrituras:
- usar `u.uow.Transaction(...)`

Si solo consulta:
- usar `u.uow.Repositories()`

## Recomendaciones prácticas para el equipo
- Antes de agregar lógica a un controller, preguntarse si realmente pertenece al use case.
- Antes de usar `u.uow.Repositories()` en una escritura, revisar si la operación toca más de una tabla o necesita auditoría atómica.
- Si un método necesita datos de varias entidades, preferir que el use case coordine varios repositorios en vez de esconder demasiada lógica cruzada dentro de un solo repo.
- Si agregas un método nuevo a un repo, primero declara el contrato en `internal/domain/repositories` y luego implementa la versión GORM en `infrastructure`.
- Mantén la firma `ctx context.Context` en toda la cadena: controller, use case y repository.

## Requisitos
- Go 1.25+
- PostgreSQL (extensión `pgcrypto`)
- Docker (opcional, para DB de tests)

## Configuración
1. Copiar `.env.example` a `.env` y ajustar variables.
2. Definir `FILES_BASE_PATH` (obligatorio para uploads).

### Bootstrap inicial
El bootstrap reutiliza la misma lógica en tres puntos:
- arranque automático de la app solo si `APP_ENV=prod|production` y `BOOTSTRAP_AUTO_ENABLED=true`
- endpoint manual excepcional `POST /api/v1/system/bootstrap/run`
- runner Docker one-shot `make docker-bootstrap`

Variables:
- `BOOTSTRAP_AUTO_ENABLED`
- `BOOTSTRAP_MANUAL_ENABLED`
- `BOOTSTRAP_TOKEN` solo para el endpoint manual
- `BOOTSTRAP_ADMIN_USERNAME`
- `BOOTSTRAP_ADMIN_PASSWORD`
- `BOOTSTRAP_ADMIN_EMAIL`
- `BOOTSTRAP_ADMIN_FULL_NAME`

El bootstrap crea y deja idempotentes:
- roles base (`Administrador`, `Vendedor`)
- permisos base (`READ`, `CREATE`, `UPDATE`, `DELETE`, `EXECUTE`)
- menús y submenús base
- grants RBAC base
- usuario admin inicial con rol primario `Administrador`

Propiedades:
- no duplica datos si ya existen
- no sobrescribe la contraseña del admin si el usuario ya existe
- activa el usuario admin si existe pero estaba inactivo
- usa una transacción con `pg_try_advisory_xact_lock(...)` para evitar carreras entre instancias

### Seed de desarrollo (opcional)
El seed solo corre cuando:
- `APP_ENV=dev` o `APP_ENV=development`
- `SEED_ENABLED=true`
- y no existe el usuario admin (`admin@gmail.com`)

Variables:
- `SEED_ENABLED`
- `SEED_ADMIN_USERNAME`
- `SEED_ADMIN_PASSWORD`
- `SEED_SELLER_USERNAME`
- `SEED_SELLER_PASSWORD`

El seed de desarrollo crea:
- roles base (`Administrador`, `Vendedor`)
- usuarios base (admin + vendedor)
- categorías y productos demo
- clientes demo
- trabajadores demo
- pedidos demo (ventas + items + ordenes de trabajo)
- citas demo
- movimientos de caja demo
- RBAC base (menús, submenús, permisos `READ/CREATE/UPDATE/DELETE/EXECUTE`)
- asignación de permisos al rol `Administrador`

Para reiniciar solo datos demo (sin borrar toda la DB) y volver a sembrar:
```bash
make seed-reset
```
Para intentar sembrar sin limpiar (respeta idempotencia):
```bash
make seed-start
```
Notas:
- Solo funciona en `APP_ENV=dev` o `APP_ENV=development`.
- El reset elimina únicamente los datos de seed identificados por marcadores/emails/códigos demo.

## Ejecutar
```bash
make run
```

## Docker (DB real + pgAdmin + API)
1. Copiar variables:
```bash
cp .env.example .env
```
2. Levantar todo:
```bash
make docker-up
```
Esto levanta:
- PostgreSQL con volumen persistente (`pg_data`)
- pgAdmin con volumen persistente (`pgadmin_data`)
- API Go en contenedor (`app_storage` para archivos)

Accesos:
- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- Scalar: `http://localhost:8080/scalar`
- pgAdmin: `http://localhost:5050`

pgAdmin se inicia con servidor preconfigurado (`docker/pgadmin/servers.json`) apuntando a `db:5432`.
Si ya habías levantado pgAdmin antes, ejecuta `make docker-reset` para recrear volumen e importar la configuración inicial.

Comandos útiles:
```bash
make docker-ps
make docker-bootstrap
make docker-logs
make docker-logs-tail
make restart
make restart-build
make docker-reset
make docker-down
```
El esquema se gestiona automáticamente con GORM `AutoMigrate` al iniciar la API.

### Ejemplo exacto: arranque automático en prod
```bash
export APP_ENV=prod
export BOOTSTRAP_AUTO_ENABLED=true
export BOOTSTRAP_ADMIN_USERNAME=admin
export BOOTSTRAP_ADMIN_PASSWORD='Admin123*'
export BOOTSTRAP_ADMIN_EMAIL=admin@photogallery.local
export BOOTSTRAP_ADMIN_FULL_NAME='Administrador General'
export BOOTSTRAP_MANUAL_ENABLED=true
export BOOTSTRAP_TOKEN='super-bootstrap-token'
make docker-up
```

### Ejemplo exacto: bootstrap manual por endpoint
```bash
curl -X POST http://localhost:8080/api/v1/system/bootstrap/run \
  -H 'X-Bootstrap-Token: super-bootstrap-token'
```

### Ejemplo exacto: bootstrap one-shot en Docker sin Go local
```bash
export BOOTSTRAP_MANUAL_ENABLED=true
export BOOTSTRAP_ADMIN_USERNAME=admin
export BOOTSTRAP_ADMIN_PASSWORD='Admin123*'
export BOOTSTRAP_ADMIN_EMAIL=admin@photogallery.local
make docker-bootstrap
```

### Prueba de idempotencia
```bash
curl -X POST http://localhost:8080/api/v1/system/bootstrap/run \
  -H 'X-Bootstrap-Token: super-bootstrap-token'

curl -X POST http://localhost:8080/api/v1/system/bootstrap/run \
  -H 'X-Bootstrap-Token: super-bootstrap-token'
```
La segunda ejecución debe registrar recursos existentes y no crear duplicados.

## Swagger
```bash
make swag
```
- UI: `http://localhost:8080/swagger/index.html`

## Formato de errores (Problem Details)
La API responde errores con `application/problem+json` siguiendo RFC 7807.

Estructura:
```json
{
  "type": "https://api.photogallery.local/problems/validation-error",
  "title": "Validation Error",
  "status": 400,
  "detail": "invalid uuid",
  "instance": "/api/v1/customers/abc",
  "code": "VALIDATION_ERROR",
  "errors": ["detalle opcional"]
}
```

Notas:
- `type`, `title`, `status`, `detail` siguen Problem Details.
- `code` y `errors` son extensiones del proyecto.

## Tests
### Suite general
```bash
make test
```

### Integración con Postgres
1. Levantar DB de test:
```bash
make test-db-up
```
2. Definir DSN de test (ver `.env.test.example`):
```bash
export DATABASE_DSN_TEST="host=localhost user=admin password=1234 dbname=photogallery_test port=5434 sslmode=disable"
```
3. Ejecutar:
```bash
make test-integration
```
4. Bajar DB de test:
```bash
make test-db-down
```

Los tests de integración hacen `skip` automático si `DATABASE_DSN_TEST` no está definido.

## Generadores
### Nueva entidad
```bash
make generate-entity name=reportes
```
Genera estructura base y también archivos `*_test.go` para:
- repositorio
- controller
- casos de uso

### Nuevo caso de uso
```bash
make generate-usecase entity=reportes name=reportes_por_mes
```
Genera:
- archivo del caso de uso
- DTO request
- test del caso de uso
- test base de repo/controller (si no existen)

## Storage de archivos (MVP implementado)
- Cliente (email sanitizado o `CustomerCode`):
  - `clientes/{customerKey}/ventas/{saleId}/...`
  - `clientes/{customerKey}/ordenes/{workOrderId}/...`
  - `clientes/{customerKey}/citas/{appointmentId}/...`
  - `clientes/{customerKey}/general/...`
- Interno empresa:
  - `empresa/interno/{YYYY}/{MM}/...`

Se persiste en BD:
- `Files.StorageRelativePath`
- `Files.StorageKind`
- `FileLinks` y `CustomerId` (cuando aplica)
- auditoría `UPLOAD`, `LINK`, `DELETE`.

## Notificaciones
Se dejó puerto e implementación `Noop` (sin envío real) con TODO para integrar WhatsApp/email sin romper arquitectura.

## Notas
- Los endpoints de seguridad del documento están implementados.
- Se incluyeron CRUDs MVP para catálogo, clientes, ventas, órdenes, agenda, caja, trabajadores, archivos y settings.
- Middleware RBAC efectivo quedó con TODO de enforcement fino por submenú/permiso.
- Mapa de controllers HTTP: `docs/CONTROLLERS_HTTP.md`.
