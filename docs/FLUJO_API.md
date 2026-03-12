# Flujo de la API (Guía práctica)

Este documento explica el flujo real de uso de la API para frontend/Postman/QA, con foco en autenticación y operaciones comunes.

Base URL: `/api/v1`

## 1) Flujo de autenticación

### 1.1 Login
Endpoint: `POST /api/v1/auth/login`

Body:
```json
{
  "username": "admin",
  "password": "Admin123*"
}
```

Respuesta (200):
- `accessToken`: token corto para autorizar requests.
- `refreshToken`: token para renovar sesión.
- `expiresIn`: segundos de vida del access token.
- `user`: datos básicos del usuario.

Errores:
- `400` cuando el body es inválido (faltan campos o formato incorrecto).
- `401` cuando las credenciales no son válidas.

Uso:
- Guardar `accessToken` y `refreshToken` en cliente.
- En cada endpoint protegido enviar:
  - `Authorization: Bearer <accessToken>`

### 1.2 ¿Para qué sirve el refresh?
Endpoint: `POST /api/v1/auth/refresh`

Se usa cuando:
- el `accessToken` expiró,
- o estás cerca de expirar y quieres renovar sin pedir login otra vez.

Body:
```json
{
  "refreshToken": "<refresh-token>"
}
```

Resultado:
- entrega un nuevo `accessToken` (y nuevo refresh según implementación).

Errores:
- `400` cuando el body está mal formado.
- `401` cuando el refresh token es inválido o expiró.

Regla práctica en frontend:
1. Haces request con access token.
2. Si responde 401 por token expirado, llamas `/auth/refresh`.
3. Reintentas el request original con el nuevo access token.
4. Si refresh falla, envías al usuario a login.

### 1.3 Consultar usuario actual
Endpoint: `GET /api/v1/auth/me`

Uso:
- al iniciar la app después de restaurar tokens,
- para mostrar perfil actual.

### 1.4 Permisos efectivos (RBAC)
Endpoint: `GET /api/v1/auth/my-permissions`

Uso:
- construir menú lateral dinámico,
- habilitar/ocultar botones por permiso.

### 1.5 Cambio de contraseña (usuario autenticado)
Endpoint: `POST /api/v1/auth/change-password`

Body:
```json
{
  "oldPassword": "Actual123*",
  "newPassword": "NuevaSegura123*"
}
```

Uso recomendado:
- pedir logout de otras sesiones si en futuro se implementa.

---

## 2) Flujo para crear un producto nuevo

Normalmente requiere que ya exista una categoría.

### Paso 1: crear categoría (si no existe)
Endpoint: `POST /api/v1/categories`

Body ejemplo:
```json
{
  "name": "Impresiones",
  "description": "Productos impresos"
}
```

### Paso 2: crear producto
Endpoint: `POST /api/v1/products`

Body ejemplo:
```json
{
  "categoryId": "<uuid-categoria>",
  "name": "Album Premium",
  "type": "Product",
  "basePrice": 250,
  "cost": 120,
  "commissionType": "Fixed",
  "commissionValue": 25,
  "requiresDelivery": true,
  "defaultLeadDays": 3,
  "notes": "Producto premium"
}
```

### Paso 3: verificar creación
- `GET /api/v1/products`
- `GET /api/v1/products/{id}`

### Paso 4 (opcional): activar/desactivar
- `PATCH /api/v1/products/{id}/activate`
- `PATCH /api/v1/products/{id}/deactivate`

---

## 3) Flujo para cambiar contraseña (admin a otro usuario)

Endpoint: `PATCH /api/v1/users/{id}/password`

Body:
```json
{
  "newPassword": "Temporal123*"
}
```

Cuándo usarlo:
- soporte/admin resetea contraseña de un usuario.

Diferencia con `auth/change-password`:
- `auth/change-password` lo usa el mismo usuario con su clave actual.
- `users/{id}/password` lo usa admin sobre otra cuenta.

---

## 4) Flujo de cliente + venta + pago

### 4.1 Crear cliente
Endpoint: `POST /api/v1/customers`

Body ejemplo:
```json
{
  "fullName": "Juan Perez",
  "phone": "3001234567",
  "email": "juan@mail.com",
  "customerCode": "CUST-0001",
  "document": "123456",
  "notes": "Cliente frecuente"
}
```

### 4.2 Crear venta
Endpoint: `POST /api/v1/sales`

Body ejemplo:
```json
{
  "customerId": "<uuid-cliente>",
  "notifyOptIn": true,
  "items": [
    {
      "productId": "<uuid-producto>",
      "quantity": 1,
      "unitPrice": 250,
      "discount": 0,
      "discountReason": "",
      "notes": ""
    }
  ]
}
```

Notas:
- Si hay ítems que requieren entrega, el sistema puede crear orden de trabajo automáticamente.

### 4.3 Registrar abono/pago
Endpoint: `POST /api/v1/sales/{id}/payments`

Body ejemplo:
```json
{
  "method": "Cash",
  "amount": 100,
  "reference": "abono inicial"
}
```

### 4.4 Consultar detalle de venta
Endpoint: `GET /api/v1/sales/{id}`

Devuelve:
- venta,
- items,
- pagos.

---

## 5) Flujo de subida de archivos

Endpoint: `POST /api/v1/files/upload`

Tipo request: `multipart/form-data`

Campos:
- `file` (archivo)
- `entityType` (ej: `Sale`, `WorkOrder`, `Appointment`, etc.)
- `entityId` (UUID entidad)
- `customerId` (opcional)
- `notes` (opcional)

### Reglas de guardado
Si involucra cliente:
- `clientes/{customerKey}/ventas/{saleId}/...`
- `clientes/{customerKey}/ordenes/{workOrderId}/...`
- `clientes/{customerKey}/citas/{appointmentId}/...`
- `clientes/{customerKey}/general/...`

Si es interno:
- `empresa/interno/{YYYY}/{MM}/...`

`customerKey`:
- email sanitizado,
- fallback: `CustomerCode`.

### Consultar/descargar/eliminar
- `GET /api/v1/files`
- `GET /api/v1/files/{id}/download`
- `DELETE /api/v1/files/{id}`

---

## 6) Flujo de workers y pagos

### Crear trabajador
- `POST /api/v1/workers`

### Pagar comisiones (FIFO)
- `POST /api/v1/workers/pay-commissions`

Body:
```json
{
  "workerId": "<uuid>",
  "method": "Cash",
  "amount": 100000,
  "notes": "pago parcial"
}
```

### Pagar salario
- `POST /api/v1/workers/pay-salary`

---

## 7) Flujo de configuración

### Obtener setting
- `GET /api/v1/settings/{key}`

### Actualizar setting
- `PUT /api/v1/settings/{key}`

Body:
```json
{
  "value": "nuevo valor"
}
```

Ejemplos típicos:
- `FILES_BASE_PATH`
- flags del sistema

---

## 8) Formato de errores

La API usa Problem Details (`application/problem+json`):

```json
{
  "type": "https://api.photogallery.local/problems/validation-error",
  "title": "Validation Error",
  "status": 400,
  "detail": "invalid uuid",
  "instance": "/api/v1/products/abc",
  "code": "VALIDATION_ERROR",
  "errors": ["detalle opcional"]
}
```

Campos:
- estándar RFC7807: `type`, `title`, `status`, `detail`, `instance`
- extensión de proyecto: `code`, `errors`

---

## 9) Secuencia recomendada para un cliente nuevo (resumen)

1. Login.
2. Consultar `me` y `my-permissions`.
3. Cargar catálogos base (`categories`, `products`, `cash/categories`, etc.).
4. Operar módulos protegidos con bearer token.
5. Ante 401, ejecutar refresh y reintentar.
6. Usar Problem Details para manejo uniforme de errores.
