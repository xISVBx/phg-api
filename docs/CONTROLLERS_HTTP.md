# Controllers HTTP (Mapa y responsabilidades)

Este documento lista los controllers HTTP de la API y el endpoint que atiende cada método.

Base path: `/api/v1`

## Reglas generales
- Los controllers son delgados: validan/parsing de request, llaman use cases y responden.
- La lógica de negocio vive en `application/use_cases`.
- Errores: `application/problem+json` (Problem Details RFC7807).

## system_auth.go
- `Health` -> `GET /health`
- `SystemInfo` -> `GET /system/info`
- `BackupRun` -> `POST /system/backup/run`
- `BackupStatus` -> `GET /system/backup/status`
- `Login` -> `POST /auth/login`
- `Refresh` -> `POST /auth/refresh`
- `Me` -> `GET /auth/me`
- `Logout` -> `POST /auth/logout`
- `ChangePassword` -> `POST /auth/change-password`
- `MyPermissions` -> `GET /auth/my-permissions`

## security_users.go
- `ListUsers` -> `GET /users`
- `CreateUser` -> `POST /users`
- `GetUser` -> `GET /users/:id`
- `UpdateUser` -> `PUT /users/:id`
- `ActivateUser` -> `PATCH /users/:id/activate`
- `DeactivateUser` -> `PATCH /users/:id/deactivate`
- `SetUserPassword` -> `PATCH /users/:id/password`
- `GetUserRoles` -> `GET /users/:id/roles`
- `SetUserRoles` -> `PUT /users/:id/roles`
- `ListUserOverrides` -> `GET /users/:id/overrides`
- `ReplaceUserOverrides` -> `PUT /users/:id/overrides`
- `CreateUserOverride` -> `POST /users/:id/overrides`
- `DeleteUserOverride` -> `DELETE /users/:id/overrides/:overrideId`

## security_acl.go
- `ListRoles` -> `GET /roles`
- `CreateRole` -> `POST /roles`
- `GetRole` -> `GET /roles/:id`
- `UpdateRole` -> `PUT /roles/:id`
- `ActivateRole` -> `PATCH /roles/:id/activate`
- `DeactivateRole` -> `PATCH /roles/:id/deactivate`
- `GetRolePermissions` -> `GET /roles/:id/permissions`
- `SetRolePermissions` -> `PUT /roles/:id/permissions`
- `ListMenus` -> `GET /menus`
- `CreateMenu` -> `POST /menus`
- `GetMenu` -> `GET /menus/:id`
- `UpdateMenu` -> `PUT /menus/:id`
- `ActivateMenu` -> `PATCH /menus/:id/activate`
- `DeactivateMenu` -> `PATCH /menus/:id/deactivate`
- `ListSubMenus` -> `GET /submenus`
- `CreateSubMenu` -> `POST /submenus`
- `GetSubMenu` -> `GET /submenus/:id`
- `UpdateSubMenu` -> `PUT /submenus/:id`
- `ActivateSubMenu` -> `PATCH /submenus/:id/activate`
- `DeactivateSubMenu` -> `PATCH /submenus/:id/deactivate`
- `ListPermissions` -> `GET /permissions`
- `CreatePermission` -> `POST /permissions`
- `GetPermission` -> `GET /permissions/:id`
- `UpdatePermission` -> `PUT /permissions/:id`

## audit.go
- `ListAudit` -> `GET /audit-logs`
- `GetAudit` -> `GET /audit-logs/:id`

## catalog.go
- `ListCategories` -> `GET /categories`
- `CreateCategory` -> `POST /categories`
- `GetCategory` -> `GET /categories/:id`
- `UpdateCategory` -> `PUT /categories/:id`
- `ActivateCategory` -> `PATCH /categories/:id/activate`
- `DeactivateCategory` -> `PATCH /categories/:id/deactivate`
- `ListProducts` -> `GET /products`
- `CreateProduct` -> `POST /products`
- `GetProduct` -> `GET /products/:id`
- `UpdateProduct` -> `PUT /products/:id`
- `ActivateProduct` -> `PATCH /products/:id/activate`
- `DeactivateProduct` -> `PATCH /products/:id/deactivate`

## customer.go
- `ListCustomers` -> `GET /customers`
- `CreateCustomer` -> `POST /customers`
- `GetCustomer` -> `GET /customers/:id`
- `UpdateCustomer` -> `PUT /customers/:id`
- `ActivateCustomer` -> `PATCH /customers/:id/activate`
- `DeactivateCustomer` -> `PATCH /customers/:id/deactivate`

## sales.go
- `ListSales` -> `GET /sales`
- `CreateSale` -> `POST /sales`
- `GetSale` -> `GET /sales/:id`
- `RegisterSalePayment` -> `POST /sales/:id/payments`

## workorders.go
- `ListWorkOrders` -> `GET /work-orders`
- `CreateWorkOrder` -> `POST /work-orders`
- `GetWorkOrder` -> `GET /work-orders/:id`
- `UpdateWorkOrder` -> `PUT /work-orders/:id`

## appointments.go
- `ListAppointments` -> `GET /appointments`
- `CreateAppointment` -> `POST /appointments`
- `GetAppointment` -> `GET /appointments/:id`
- `UpdateAppointment` -> `PUT /appointments/:id`

## files.go
- `ListFiles` -> `GET /files`
- `UploadFile` -> `POST /files/upload`
- `DownloadFile` -> `GET /files/:id/download`
- `DeleteFile` -> `DELETE /files/:id`

## cash.go
- `ListCashCategories` -> `GET /cash/categories`
- `CreateCashCategory` -> `POST /cash/categories`
- `GetCashCategory` -> `GET /cash/categories/:id`
- `UpdateCashCategory` -> `PUT /cash/categories/:id`
- `ListCashMovements` -> `GET /cash/movements`
- `CreateCashMovement` -> `POST /cash/movements`
- `GetCashMovement` -> `GET /cash/movements/:id`
- `UpdateCashMovement` -> `PUT /cash/movements/:id`

## workers.go
- `ListWorkers` -> `GET /workers`
- `CreateWorker` -> `POST /workers`
- `GetWorker` -> `GET /workers/:id`
- `UpdateWorker` -> `PUT /workers/:id`
- `ActivateWorker` -> `PATCH /workers/:id/activate`
- `DeactivateWorker` -> `PATCH /workers/:id/deactivate`
- `PayCommissions` -> `POST /workers/pay-commissions`
- `PaySalary` -> `POST /workers/pay-salary`

## settings.go
- `GetSystemSetting` -> `GET /settings/:key`
- `SetSystemSetting` -> `PUT /settings/:key`
