# Definición del proyecto (VERSIÓN COMPLETA + Storage por cliente)

## 1) Resumen
Plataforma para una foto-tienda que permite gestionar **catálogo**, **clientes**, **ventas con abonos**, **órdenes de producción/entrega**, **agenda de estudios**, **notificaciones**, **archivos (fotos)**, **caja (entradas/salidas)**, **trabajadores (sueldo + comisiones)**, **anulaciones/devoluciones/cancelaciones**, y **seguridad avanzada** por **Usuarios → Roles → Menús → Submenús → Permisos**, con **overrides por usuario**.

---

## 2) Alcance funcional (requerimientos)

### 2.1 Catálogo
- CRUD de **Categorías**.
- CRUD de **Productos/Servicios** con:
  - Precio base, costo, comisión (% o fijo), tipo (producto/servicio/estudio),
  - requiere entrega, tiempo estimado/fecha manual, activo/inactivo, notas.
- **Snapshot en ventas**: al vender, el ítem guarda precio/costo/comisión aplicados (para histórico).

### 2.2 Clientes
- CRUD de clientes, búsqueda rápida, datos de contacto.
- Reglas:
  - **Email puede ser unique** (si se decide usarlo como carpeta principal).
  - Recomendación: agregar `CustomerCode` unique para clientes sin email (fallback).

### 2.3 Ventas y pagos
- Crear venta con ítems:
  - precio editable, descuento, motivo, observación.
- Pagos múltiples por venta (abonos).
- Estados financieros: Pendiente / Abonada / Pagada / Anulada / Cancelada.
- Generación de **Orden de trabajo** cuando haya ítems con entrega.

### 2.4 Órdenes de producción / entregas
- Estados: Creado → En desarrollo → Listo para entrega → Entregado.
- Fecha estimada de entrega por orden (o por ítem si se usa detalle).
- Indicadores automáticos:
  - **Atrasado** si pasa la fecha estimada y no está entregado.
  - **Listo sin reclamar** con contador de días desde que quedó listo.

### 2.5 Estudios (agenda)
- Citas vinculadas a venta (ideal).
- Estados: Programada / Confirmada / Realizada / Cancelada / No asistió.
- Pagos y abonos aplicables.

### 2.6 Notificaciones
- Opt-in por venta/orden: “Enviar notificaciones Sí/No”.
- Reglas configurables:
  - 1 día antes de entrega
  - día de entrega
  - atrasos cada X días
  - listo sin reclamar (2/5/10 días, configurable)
- Historial de notificaciones:
  - Pendiente / Enviada / Fallida + reintentos.

### 2.7 Archivos (fotos) — **requerimiento ampliado (storage por carpetas)**
- Subir/ver/descargar/eliminar (según permisos).
- Asociar a Cliente/Venta/Orden/Ítem/Cita.
- Miniaturas para imágenes + control de tamaño y auditoría.
- **Path inicial configurable**: `FILES_BASE_PATH`.
- **Regla de almacenamiento físico**:

#### 2.7.1 Estructura base propuesta
- Carpeta raíz configurada:
  - `{FILES_BASE_PATH}/`

#### 2.7.2 Archivos que involucran cliente
> “Involucra cliente” = asociado a **Sale, WorkOrder, SaleItem, Appointment** o al mismo Customer.

- Se guardan bajo:
  - `{FILES_BASE_PATH}/clientes/{customerKey}/...`

- `customerKey`:
  - Preferido: **email del cliente** (unique) sanitizado.
  - Fallback: `CustomerCode` unique (si no hay email o no se usa email como key).

- Subcarpetas sugeridas:
  - Ventas: `ventas/{saleId}/`
  - Órdenes: `ordenes/{workOrderId}/`
  - Citas: `citas/{appointmentId}/`
  - General cliente: `general/`

- Ejemplos:
  - `clientes/juan_perez_gmail_com/ventas/9f3a.../img001.jpg`
  - `clientes/CUST-000245/ordenes/1a2b.../referencia.png`

#### 2.7.3 Archivos internos (empresa / no cliente)
- Si el archivo no está asociado a cliente (ej: manual interno, plantilla, archivos de operación):
  - `{FILES_BASE_PATH}/empresa/interno/{YYYY}/{MM}/...`

- Si quieres separar:
  - `{FILES_BASE_PATH}/empresa/finanzas/{YYYY}/{MM}/...`
  - `{FILES_BASE_PATH}/empresa/otros/{YYYY}/{MM}/...`

#### 2.7.4 Reglas de BD para archivos
- La BD debe guardar:
  - `StorageRelativePath` (ruta relativa desde base path)
  - `StorageKind` (clasificación: Sale/WorkOrder/Appointment/CustomerGeneral/CompanyInternal/etc.)
  - `CustomerId` (si aplica) para consultas rápidas.
  - asociación (FileLinks).

#### 2.7.5 Configuración de storage
- En configuración debe existir:
  - `FILES_BASE_PATH`
  - `FILES_MAX_SIZE_MB`
  - `FILES_ALLOWED_MIME` (lista)
  - `CUSTOMER_KEY_MODE = email|code|both` (recomendado: `both` con fallback)

### 2.8 Anulaciones, cancelaciones y devoluciones
- **Anulación**: motivo obligatorio; si hubo pagos, definir tratamiento:
  - reembolso / saldo a favor / retención (según política).
- **Cancelación por no reclamar**:
  - cuando abonaron y no reclaman; definir tratamiento (configurable).
- **Devolución** total/parcial por ítem/cantidad:
  - genera ajuste financiero y afecta comisiones.

### 2.9 Caja (entradas/salidas)
- Movimientos de caja:
  - entradas (pagos de ventas, ingresos varios)
  - salidas (egresos, reembolsos, pagos a trabajadores)
- Categorías configurables.
- (Opcional recomendado) Apertura y cierre de caja con arqueo/diferencias.

### 2.10 Trabajadores, comisiones y pagos
- CRUD trabajadores:
  - sueldo fijo (0 o >0), solo comisiones o ambos.
- Cartera de comisiones (ledger):
  - Ganada / Parcialmente pagada / Pagada / Ajustada/Anulada.
- Pago masivo de comisiones por monto (FIFO):
  - Ingreso de monto → recorre comisiones pendientes y marca Pagadas/Parcial.
- Pagos de sueldo como salida de caja.
- Historial de pagos (comisiones y sueldos).

### 2.11 Seguridad (Usuarios/Roles/Menús/Submenús/Permisos) + overrides
- CRUD de usuarios/roles/menús/submenús/permisos.
- Overrides por usuario (grant/revoke) sobre rol base.
- Auditoría obligatoria:
  - descuentos, anulaciones/devoluciones/cancelaciones
  - pagos/reembolsos
  - caja
  - pagos a trabajadores
  - cambios de estado de órdenes
  - **archivos: upload/delete/link/move**

---

## 3) Menús, submenús y permisos (lista exacta propuesta)
> Convención: `MENU_CODE / SUBMENU_CODE / PERM_CODE`

### 3.1 Menú: Dashboard (`DASH`)
**Submenús**
- `DASH_HOME` — Inicio
- `DASH_ALERTS` — Alertas (atrasados, listo sin reclamar, cuentas por cobrar)

**Permisos**
- `VIEW`
- `VIEW_ALERTS`

---

### 3.2 Menú: Ventas (`SALES`)
**Submenús**
- `SALES_NEW` — Realizar venta
- `SALES_LIST` — Listado de ventas
- `SALES_DAY` — Ventas del día
- `SALES_PAYMENTS` — Pagos/Abonos
- `SALES_DOCS` — Documentos/Comprobantes

**Permisos**
- `VIEW`, `CREATE`, `EDIT`
- `APPLY_DISCOUNT`, `EDIT_PRICE`
- `REGISTER_PAYMENT`, `REFUND_PAYMENT`
- `EXPORT`, `PRINT`

---

### 3.3 Menú: Producción y Entregas (`OPS`)
**Submenús**
- `OPS_ORDERS` — Órdenes
- `OPS_DELIVERIES` — Entregas
- `OPS_STATUS_BOARD` — Tablero

**Permisos**
- `VIEW`, `EDIT`
- `CHANGE_STATUS`, `SET_DUE_DATE`, `MARK_DELIVERED`
- `EXPORT`

---

### 3.4 Menú: Estudios (Agenda) (`APPT`)
**Submenús**
- `APPT_CALENDAR`
- `APPT_CREATE`
- `APPT_LIST`

**Permisos**
- `VIEW`, `CREATE`, `EDIT`, `CANCEL`, `CHANGE_STATUS`
- `EXPORT`

---

### 3.5 Menú: Clientes (`CRM`)
**Submenús**
- `CRM_CLIENTS`

**Permisos**
- `VIEW`, `CREATE`, `EDIT`, `DEACTIVATE`, `EXPORT`

---

### 3.6 Menú: Catálogo (`CAT`)
**Submenús**
- `CAT_CATEGORIES`
- `CAT_PRODUCTS`
- `CAT_PRICING`

**Permisos**
- `VIEW`, `CREATE`, `EDIT`, `DEACTIVATE`, `EXPORT`

---

### 3.7 Menú: Archivos (Fotos) (`FILES`)
**Submenús**
- `FILES_LIBRARY`
- `FILES_UPLOAD`

**Permisos**
- `VIEW`
- `UPLOAD`
- `DOWNLOAD`
- `DELETE`
- `LINK_TO_ENTITY`
- `MOVE_FILE` *(si se permite reubicar al corregir asociación)*
- `EXPORT`

---

### 3.8 Menú: Caja y Movimientos (`CASH`)
**Submenús**
- `CASH_MOVES`
- `CASH_CATEGORIES`
- `CASH_OPEN_CLOSE`
- `CASH_RECONCILE`

**Permisos**
- `VIEW`
- `CREATE_MOVE`, `EDIT_MOVE`
- `OPEN_CASH`, `CLOSE_CASH`
- `EXPORT`

---

### 3.9 Menú: Anulaciones/Devoluciones/Ajustes (`ADJ`)
**Submenús**
- `ADJ_VOID`
- `ADJ_REFUNDS`
- `ADJ_RETURNS`
- `ADJ_CANCEL_NOCLAIM`
- `ADJ_CREDITNOTES`

**Permisos**
- `VIEW`
- `VOID_SALE`, `RETURN_ITEM`, `CANCEL_SALE`
- `CREATE_REFUND`, `CREATE_CREDIT_NOTE`
- `APPROVE`
- `EXPORT`

---

### 3.10 Menú: Trabajadores y Comisiones (`HR`)
**Submenús**
- `HR_WORKERS`
- `HR_COMMISSIONS`
- `HR_PAY_COMM`
- `HR_PAYROLL`
- `HR_PAY_HISTORY`

**Permisos**
- `VIEW`, `CREATE`, `EDIT`, `DEACTIVATE`
- `PAY_COMMISSIONS`, `PAY_SALARY`
- `ADJUST_COMMISSION`
- `EXPORT`

---

### 3.11 Menú: Reportes (`RPT`)
**Submenús**
- `RPT_SALES`
- `RPT_CASH`
- `RPT_ORDERS`
- `RPT_COMMISSIONS`
- `RPT_ACCOUNTS_RECEIVABLE`
- `RPT_ADJUSTMENTS`

**Permisos**
- `VIEW`, `EXPORT`

---

### 3.12 Menú: Configuración y Seguridad (`SET`)
**Submenús**
- `SET_USERS`
- `SET_ROLES`
- `SET_MENUS`
- `SET_SUBMENUS`
- `SET_PERMS`
- `SET_USER_OVERRIDES`
- `SET_NOTIF_RULES`
- `SET_BACKUP`
- `SET_STORAGE` *(para FILES_BASE_PATH y reglas de storage)*

**Permisos**
- `VIEW`, `CREATE`, `EDIT`, `DEACTIVATE`
- `ASSIGN_ROLE`, `ASSIGN_MENU`, `ASSIGN_PERMISSIONS`, `MANAGE_OVERRIDES`
- `BACKUP_RUN`
- `CONFIGURE_STORAGE`

---

## 4) Historias de usuario (con criterios de aceptación y pruebas)

### 4.1 Catálogo
**HU-01: Crear categoría**
- Como Admin, quiero crear categorías para organizar los productos.
- Criterios de aceptación:
  - Given estoy en Categorías, When creo una con nombre válido, Then queda activa y visible en listados.
  - Given nombre duplicado, When intento guardar, Then se impide y muestra error.
- Pruebas:
  - Crear “Portarretratos”.
  - Duplicar “Portarretratos” → falla.

**HU-02: Crear producto con comisión y costo**
- Como Admin, quiero crear productos con precio, costo y comisión.
- Criterios de aceptación:
  - Given datos completos, When guardo, Then aparece en ventas.
  - Given comisión inválida, When guardo, Then error.
- Pruebas:
  - Comisión 10% ok.
  - Comisión 150% → error.

---

### 4.2 Clientes
**HU-03: Registrar cliente (email/code unique)**
- Como Vendedor, quiero registrar clientes para asociarles ventas y archivos.
- Criterios de aceptación:
  - Given se usa email como key, When guardo email duplicado, Then error por unique.
  - Given cliente sin email, When guardo, Then se exige/genera `CustomerCode` unique.
- Pruebas:
  - Email duplicado falla.
  - Sin email genera code.

---

### 4.3 Ventas y pagos
**HU-04: Realizar una venta**
- Como Vendedor, quiero crear una venta con varios ítems.
- Criterios de aceptación:
  - Given productos activos, When agrego ítems, Then total se calcula bien.
- Pruebas:
  - 2 ítems, verificar subtotal/total.

**HU-05: Editar precio con motivo y permiso**
- Como Vendedor, quiero cambiar precio en casos excepcionales.
- Criterios de aceptación:
  - Given sin `EDIT_PRICE`, When intento cambiar precio, Then se bloquea.
  - Given con permiso, When cambio precio, Then motivo obligatorio y auditado.
- Pruebas:
  - Sin permiso bloquea.
  - Con permiso exige motivo.

**HU-06: Registrar abonos**
- Como Recepción, quiero registrar abonos a una venta.
- Criterios de aceptación:
  - Given saldo pendiente, When registro pago, Then saldo baja y estado cambia.
- Pruebas:
  - Abono parcial → Abonada.
  - Pago total → Pagada.

---

### 4.4 Órdenes/Entregas
**HU-07: Orden automática al confirmar venta**
- Como Sistema, quiero crear orden si hay ítems con entrega.
- Criterios de aceptación:
  - Given venta con entrega, When confirmo, Then se crea WorkOrder con due date.
- Pruebas:
  - Con entrega crea orden; sin entrega no.

**HU-08: Cambiar estados de una orden**
- Como Producción, quiero cambiar estados de la orden.
- Criterios de aceptación:
  - Given orden existe, When cambio estado, Then se registra usuario/fecha y se ve en tablero.
- Pruebas:
  - Creado → En desarrollo.
  - En desarrollo → Listo.

**HU-09: Alertas de atrasos y listo sin reclamar**
- Como Admin, quiero ver alertas.
- Criterios de aceptación:
  - Given hoy > dueDate y no entregado, Then aparece como Atrasado.
  - Given Ready y pasan X días, Then aparece como Listo sin reclamar con contador.
- Pruebas:
  - Simular fechas.

---

### 4.5 Agenda
**HU-10: Agendar estudio**
- Como Recepción, quiero agendar una cita.
- Criterios de aceptación:
  - Given producto tipo estudio, When guardo venta, Then permite/solicita agendar.
- Pruebas:
  - Venta con estudio muestra agenda.

---

### 4.6 Archivos (CRÍTICO: carpeta por cliente + tipos)
**HU-11: Subir imagen asociada a una venta y guardarla en carpeta del cliente**
- Como Usuario autorizado, quiero subir imágenes para el proceso de una venta y que otro usuario pueda verlas/descargarlas.
- Criterios de aceptación:
  - Given `FILES_BASE_PATH` está configurado, When subo una imagen **asociada a Sale** con un cliente, Then:
    - El sistema obtiene `customerKey` (email sanitizado si existe; si no, `CustomerCode`).
    - Guarda físicamente en:  
      `clientes/{customerKey}/ventas/{saleId}/`
    - Registra en BD:
      - `Files.StorageRelativePath = clientes/{customerKey}/ventas/{saleId}/{filename}`
      - `Files.StorageKind = Sale`
      - `FileLinks` con `EntityType=Sale` y `EntityId=saleId`
      - `FileLinks.CustomerId = customerId`
    - Registra auditoría (UPLOAD + LINK).
  - Given la venta/orden/cita se puede consultar, When otro usuario con permiso `VIEW` y `DOWNLOAD` entra, Then puede ver y descargar.
  - Given `FILES_BASE_PATH` no está configurado, When subo, Then error de configuración y no crea registros.
- Pruebas:
  - Subir imagen a venta con cliente con email → valida ruta.
  - Subir imagen a venta con cliente sin email → usa `CustomerCode`.
  - Descargar desde otra sesión → ok.
  - Sin permisos DOWNLOAD → bloquea.

**HU-11B: Subir imagen asociada a una orden y guardarla en carpeta del cliente**
- Criterios de aceptación:
  - When subo a `WorkOrder`, Then guarda en:
    - `clientes/{customerKey}/ordenes/{workOrderId}/`
  - Y BD guarda StorageRelativePath + link a WorkOrder.
- Pruebas:
  - Subir a orden y verificar ruta.

**HU-11C: Subir archivo interno de la empresa**
- Como Admin/Usuario autorizado, quiero subir archivos internos.
- Criterios de aceptación:
  - Given tipo `CompanyInternal`, When subo, Then guarda en:
    - `empresa/interno/{YYYY}/{MM}/`
  - BD registra `StorageKind=CompanyInternal` sin CustomerId.
- Pruebas:
  - Subir archivo interno y verificar ruta.

**HU-11D: Configurar path inicial de almacenamiento**
- Como Admin, quiero configurar `FILES_BASE_PATH`.
- Criterios de aceptación:
  - Given soy admin, When guardo un path válido, Then queda persistido y se valida que exista/sea accesible.
- Pruebas:
  - Path existente OK.
  - Path inválido → error.

---

### 4.7 Ajustes (anulación/devolución/cancelación)
**HU-12: Anular venta**
- Criterios de aceptación:
  - Sin pagos → anula directo.
  - Con pagos → exige tratamiento (refund/storeCredit/retain) y refleja caja.
- Pruebas:
  - Anular sin pagos.
  - Anular con pagos y reembolso.

**HU-13: Devolver ítem**
- Criterios de aceptación:
  - Devolución parcial crea ajuste, afecta total y comisiones.
- Pruebas:
  - Devolución y revisar comisión.

**HU-14: Cancelación por no reclamar**
- Criterios de aceptación:
  - Exige elegir tratamiento del abono y registra caja si aplica.
- Pruebas:
  - Retener vs devolver.

---

### 4.8 Caja
**HU-15: Registrar egreso**
- Criterios de aceptación:
  - Requiere categoría, método, valor; queda en reportes y auditoría.
- Pruebas:
  - Egreso insumos.

**HU-16: Apertura/cierre caja**
- Criterios de aceptación:
  - Apertura con monto inicial; cierre calcula diferencia.
- Pruebas:
  - Cerrar con diferencia.

---

### 4.9 Trabajadores y comisiones
**HU-17: Generar comisiones**
- Criterios de aceptación:
  - Confirmar venta crea entradas de comisión por ítem.
- Pruebas:
  - 2 ítems → 2 comisiones.

**HU-18: Pagar comisiones por monto (FIFO)**
- Criterios de aceptación:
  - Paga más antiguas primero y marca pagada/parcial; crea salida de caja.
- Pruebas:
  - 60k/60k/60k, pago 100k.

**HU-19: Pagar sueldo**
- Criterios de aceptación:
  - Crea salida de caja y queda historial.
- Pruebas:
  - Pago quincenal.

---

### 4.10 Seguridad
**HU-20: Overrides**
- Criterios de aceptación:
  - Grant/Revoke sobre rol base funciona y se ve reflejado en UI y API.
- Pruebas:
  - Rol sin permiso + grant → permite.
  - Rol con permiso + revoke → bloquea.

---

## 5) Modelo de base de datos (propuesto) — COMPLETO con storage por carpetas

### 5.1 Seguridad
**Users**
- Id (PK)
- Username (unique)
- PasswordHash
- FullName
- Phone
- Email
- IsActive
- CreatedAtUtc

**Roles**
- Id (PK)
- Name (unique)
- Description
- IsActive

**Menus**
- Id (PK)
- Code (unique)
- Name
- DisplayOrder
- IsActive

**SubMenus**
- Id (PK)
- MenuId (FK -> Menus)
- Code (unique)
- Name
- Route (nullable)
- DisplayOrder
- IsActive

**Permissions**
- Id (PK)
- Code (unique)
- Name
- Description

**RoleSubMenuPermissions**
- RoleId (FK) (PK part)
- SubMenuId (FK) (PK part)
- PermissionId (FK) (PK part)

**UserRoles**
- UserId (FK) (PK part)
- RoleId (FK) (PK part)
- IsPrimary (bool)

**UserPermissionOverrides**
- Id (PK)
- UserId (FK)
- SubMenuId (FK)
- PermissionId (FK)
- Mode (enum: Grant, Revoke)
- CreatedAtUtc

**AuditLogs**
- Id (PK)
- ActorUserId (FK -> Users)
- EntityType (string)
- EntityId (uuid/string)
- Action (string)
- DataJson (text)
- CreatedAtUtc
- IpAddress (nullable)

---

### 5.2 Catálogo y clientes
**Categories**
- Id (PK)
- Name (unique)
- Description
- IsActive

**Products**
- Id (PK)
- CategoryId (FK)
- Name
- Type (enum: Physical, Service, Study)
- BasePrice (decimal)
- Cost (decimal)
- CommissionType (enum: None, Percent, Fixed)
- CommissionValue (decimal)
- RequiresDelivery (bool)
- DefaultLeadDays (nullable int)
- IsActive
- Notes

**Customers**
- Id (PK)
- FullName
- Phone
- Email (nullable) **(unique si se usa como key)**
- CustomerCode (nullable) **(unique; recomendado always set)**
- Document (nullable)
- Notes
- IsActive

---

### 5.3 Ventas
**Sales**
- Id (PK)
- CustomerId (FK -> Customers, nullable)
- SellerUserId (FK -> Users)
- Status (enum)
- NotifyOptIn (bool)
- Subtotal, DiscountTotal, Total
- TotalCostSnapshot
- TotalCommissionSnapshot
- CreatedAtUtc
- ConfirmedAtUtc (nullable)

**SaleItems**
- Id (PK)
- SaleId (FK -> Sales)
- ProductId (FK -> Products)
- Quantity
- UnitPriceSnapshot
- UnitCostSnapshot
- CommissionTypeSnapshot
- CommissionValueSnapshot
- DiscountSnapshot
- DiscountReason (nullable)
- Notes (nullable)
- RequiresDeliverySnapshot
- LeadDaysSnapshot (nullable)

**SalePayments**
- Id (PK)
- SaleId (FK)
- Method (enum)
- Amount
- Reference (nullable)
- PaidAtUtc
- CreatedByUserId (FK)

---

### 5.4 Órdenes
**WorkOrders**
- Id (PK)
- SaleId (FK -> Sales)
- Status (enum)
- DueDateUtc (nullable)
- ResponsibleUserId (nullable FK -> Users)
- Notes
- CreatedAtUtc
- DeliveredAtUtc (nullable)

**WorkOrderItems** *(opcional)*
- Id (PK)
- WorkOrderId (FK)
- SaleItemId (FK)
- Status (enum)
- DueDateUtc (nullable)
- Notes

---

### 5.5 Agenda
**Appointments**
- Id (PK)
- CustomerId (FK -> Customers)
- SaleId (nullable FK -> Sales)
- ProductId (FK -> Products)
- StartsAtUtc
- EndsAtUtc (nullable)
- Status (enum)
- Notes
- CreatedByUserId (FK)

---

### 5.6 Archivos (con storage relativo + customerKey)
**Files**
- Id (PK)
- OriginalName
- ContentType
- SizeBytes
- UploadedByUserId (FK -> Users)
- UploadedAtUtc
- StorageRelativePath **(NOT NULL)**  
  - Ej: `clientes/juan_gmail_com/ventas/<saleId>/img001.jpg`
- StorageKind (enum):
  - `CustomerGeneral | Sale | WorkOrder | SaleItem | Appointment | CompanyInternal | CompanyFinance | Other`
- StorageBaseKey (nullable string):
  - Ej: `clientes/juan_gmail_com` o `empresa/interno/2026/03`

**FileLinks**
- Id (PK)
- FileId (FK -> Files)
- EntityType (enum/string: Customer, Sale, WorkOrder, SaleItem, Appointment)
- EntityId (uuid)
- CustomerId (nullable FK -> Customers)  *(si involucra cliente)*
- Notes (nullable)

---

### 5.7 Ajustes
**Adjustments**
- Id (PK)
- SaleId (FK -> Sales)
- Type (enum: Void, Return, CancelNoClaim, CreditNote, DebitNote)
- Reason
- PolicyResult (enum: Refund, StoreCredit, Retain, Penalty)
- AmountImpact
- CreatedByUserId (FK -> Users)
- CreatedAtUtc

**AdjustmentItems**
- Id (PK)
- AdjustmentId (FK)
- SaleItemId (FK)
- QuantityReturned
- AmountImpact

---

### 5.8 Caja
**CashCategories**
- Id (PK)
- Name (unique)
- Type (enum: In, Out, Both)
- IsActive

**CashSessions** *(opcional)*
- Id (PK)
- OpenedByUserId (FK)
- OpenedAtUtc
- OpeningAmount
- ClosedByUserId (nullable FK)
- ClosedAtUtc (nullable)
- ClosingCountedAmount (nullable)
- Difference (nullable)
- Status (enum: Open, Closed)

**CashMovements**
- Id (PK)
- SessionId (nullable FK)
- Type (enum: In, Out)
- CategoryId (FK)
- Method (enum)
- Amount
- Reference (nullable)
- RelatedEntityType (nullable)
- RelatedEntityId (nullable)
- Notes (nullable)
- CreatedByUserId (FK)
- CreatedAtUtc

---

### 5.9 Trabajadores / Comisiones
**Workers**
- Id (PK)
- FullName
- Phone
- Email (nullable)
- IsActive
- FixedSalary (decimal)
- SalaryPeriod (enum, nullable)
- Notes

**CommissionEntries**
- Id (PK)
- WorkerId (FK -> Workers)
- SaleItemId (FK -> SaleItems)
- EarnedAtUtc
- Amount
- PaidAmount
- Status (enum: Earned, PartiallyPaid, Paid, Voided, Adjusted)

**WorkerPayments**
- Id (PK)
- WorkerId (FK)
- Type (enum: Commission, Salary)
- Method (enum)
- Amount
- Notes (nullable)
- PaidAtUtc
- CreatedByUserId (FK)
- CashMovementId (nullable FK -> CashMovements)

**WorkerPaymentAllocations**
- Id (PK)
- WorkerPaymentId (FK)
- CommissionEntryId (FK)
- AmountApplied

---

### 5.10 Configuración (para FILES_BASE_PATH y demás)
**AppSettings**
- Key (PK)  // e.g. FILES_BASE_PATH
- Value (text)
- UpdatedAtUtc
- UpdatedByUserId (nullable FK -> Users)

---

## 6) Reglas clave (business rules)
1. **Snapshots** en venta.
2. **Override manda** (Grant/Revoke) sobre rol base.
3. **No borrar**: usar reversos/ajustes + auditoría.
4. **FIFO comisiones** en pagos por monto.
5. **Atrasos** y **listo sin reclamar** por fechas/estado.
6. **Storage**:
   - `FullPath = FILES_BASE_PATH + "/" + StorageRelativePath`
   - Si involucra cliente → `clientes/{customerKey}/...`
   - Si es interno → `empresa/interno/{YYYY}/{MM}/...`
7. Auditoría para: upload/delete/link/move de archivos.

---

## 7) Checklist de producción (mínimo)
- Permisos + auditoría en módulos críticos.
- Backups de BD + respaldo del `FILES_BASE_PATH`.
- Reportes mínimos listos.
- Validación de `FILES_BASE_PATH` en arranque (y desde configuración).
- Límites/MIME/seguridad en uploads.

---
