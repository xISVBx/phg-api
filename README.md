# PhotoGallery API Go

Backend MVP en Go + Gin + PostgreSQL + GORM con arquitectura por capas tipo Dilogum:
- `domain` entidades + interfaces de repositorio.
- `application` DTOs, mappers, puertos y use cases.
- `infrastructure` GORM repos, UoW, JWT/bcrypt, servicios externos.
- `web` controllers HTTP, middlewares y rutas.

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
