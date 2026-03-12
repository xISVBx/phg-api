# Endpoints del Backend (Go) — API REST v1

> Backend propuesto: **Go** (Gin/Echo/Fiber) + PostgreSQL/MySQL/SQLite (según despliegue) + JWT.
>
> Convenciones
- Base URL: `/api/v1`
- Auth: `Authorization: Bearer <token>`
- Fechas: ISO8601 UTC (`2026-03-04T19:30:00Z`)
- Paginación: `?page=1&pageSize=25`
- Búsqueda: `?q=texto`
- Orden: `?sort=createdAtUtc&dir=desc`
- Respuestas:
  - OK: `{ "data": ..., "meta": ... }`
  - Error: `application/problem+json` (RFC 7807)
    ```json
    {
      "type": "https://api.photogallery.local/problems/validation-error",
      "title": "Validation Error",
      "status": 400,
      "detail": "mensaje de error",
      "instance": "/api/v1/recurso",
      "code": "VALIDATION_ERROR",
      "errors": ["detalle opcional"]
    }
    ```

---

## 0) Health / System
### Health
- `GET /api/v1/health`
  - 200 OK: `{ data: { status: "ok", timeUtc: "...", version: "..." } }`

### System Info
- `GET /api/v1/system/info`
- `POST /api/v1/system/backup/run` *(si aplica)*
- `GET  /api/v1/system/backup/status` *(si aplica)*

---

## 1) Auth & Sesión
### Login / Refresh / Me
- `POST /api/v1/auth/login`
  - Body:
    ```json
    { "username": "admin", "password": "1234" }
    ```
  - 200:
    ```json
    { "data": { "accessToken": "...", "refreshToken": "...", "expiresIn": 3600, "user": { "id":"...", "fullName":"..." } } }
    ```

- `POST /api/v1/auth/refresh`
  - Body: `{ "refreshToken": "..." }`

- `GET /api/v1/auth/me`
  - 200: usuario actual + rol base

- `POST /api/v1/auth/logout` *(opcional: blacklist refresh token)*
- `POST /api/v1/auth/change-password`
  - Body: `{ "oldPassword":"...", "newPassword":"..." }`

### Permisos efectivos
- `GET /api/v1/auth/my-permissions`
  - Devuelve árbol: Menú → Submenú → Permisos efectivos (rol + overrides)

---

## 2) Seguridad: Usuarios / Roles / Menús / Submenús / Permisos / Overrides
> Recomendación: endpoints **bulk** para evitar 100 llamadas.

### 2.1 Users
- `GET    /api/v1/users`
  - Filtros: `?q=juan&isActive=true`
- `POST   /api/v1/users`
- `GET    /api/v1/users/{id}`
- `PUT    /api/v1/users/{id}`
- `PATCH  /api/v1/users/{id}/activate`
- `PATCH  /api/v1/users/{id}/deactivate`
- `PATCH  /api/v1/users/{id}/password` *(admin set password)*

#### Rol base / roles
- `GET /api/v1/users/{id}/roles`
- `PUT /api/v1/users/{id}/roles`
  - Body ejemplo (1 rol):
    ```json
    { "primaryRoleId": "role_uuid" }
    ```

#### Overrides
- `GET  /api/v1/users/{id}/overrides`
- `PUT  /api/v1/users/{id}/overrides` *(bulk replace)*
  - Body:
    ```json
    {
      "items": [
        { "subMenuId":"...", "permissionId":"...", "mode":"Grant" },
        { "subMenuId":"...", "permissionId":"...", "mode":"Revoke" }
      ]
    }
    ```
- `POST /api/v1/users/{id}/overrides`
- `DELETE /api/v1/users/{id}/overrides/{overrideId}`

---

### 2.2 Roles
- `GET    /api/v1/roles`
- `POST   /api/v1/roles`
- `GET    /api/v1/roles/{id}`
- `PUT    /api/v1/roles/{id}`
- `PATCH  /api/v1/roles/{id}/activate`
- `PATCH  /api/v1/roles/{id}/deactivate`

#### Permisos del rol (RoleSubMenuPermissions)
- `GET /api/v1/roles/{id}/permissions`
- `PUT /api/v1/roles/{id}/permissions` *(bulk set)*
  - Body:
    ```json
    {
      "items": [
        { "subMenuId":"...", "permissionIds":["...","..."] }
      ]
    }
    ```

---

### 2.3 Menus / SubMenus / Permissions
- `GET  /api/v1/menus`
- `POST /api/v1/menus`
- `GET  /api/v1/menus/{id}`
- `PUT  /api/v1/menus/{id}`
- `PATCH /api/v1/menus/{id}/activate`
- `PATCH /api/v1/menus/{id}/deactivate`

- `GET  /api/v1/submenus`
- `POST /api/v1/submenus`
- `GET  /api/v1/submenus/{id}`
- `PUT  /api/v1/submenus/{id}`
- `PATCH /api/v1/submenus/{id}/activate`
- `PATCH /api/v1/submenus/{id}/deactivate`

- `GET  /api/v1/permissions`
- `POST /api/v1/permissions`
- `GET  /api/v1/permissions/{id}`
- `PUT  /api/v1/permissions/{id}`

---

### 2.4 Auditoría
- `GET /api/v1/audit-logs`
  - Filtros: `?from=...&to=...&actorUserId=...&entityType=Sale&action=VOID`
- `GET /api/v1/audit-logs/{id}`

---

## 3) Catálogo
### 3.1 Cat
