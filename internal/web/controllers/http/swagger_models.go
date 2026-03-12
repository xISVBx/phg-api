package http

import (
	"photogallery/api_go/internal/application/dtos"
	appointmentreq "photogallery/api_go/internal/application/dtos/request/appointment"
	authreq "photogallery/api_go/internal/application/dtos/request/auth"
	cashreq "photogallery/api_go/internal/application/dtos/request/cash"
	catalogreq "photogallery/api_go/internal/application/dtos/request/catalog"
	customerreq "photogallery/api_go/internal/application/dtos/request/customer"
	filesreq "photogallery/api_go/internal/application/dtos/request/files"
	salesreq "photogallery/api_go/internal/application/dtos/request/sales"
	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	systemreq "photogallery/api_go/internal/application/dtos/request/system"
	workerreq "photogallery/api_go/internal/application/dtos/request/worker"
	workorderreq "photogallery/api_go/internal/application/dtos/request/workorder"
	authresp "photogallery/api_go/internal/application/dtos/responses/auth"
	securityresp "photogallery/api_go/internal/application/dtos/responses/security"
	"photogallery/api_go/internal/domain/entities"
	webutils "photogallery/api_go/internal/web/utils"
)

type SwaggerLoginRequest = authreq.LoginRequestDTO
type SwaggerRefreshRequest = authreq.RefreshRequestDTO
type SwaggerChangePasswordRequest = authreq.ChangePasswordRequestDTO

type SwaggerCreateUserRequest = securityreq.CreateUserRequestDTO
type SwaggerUpdateUserRequest = securityreq.UpdateUserRequestDTO
type SwaggerSetPasswordRequest = securityreq.SetPasswordRequestDTO
type SwaggerSetPrimaryRoleRequest = securityreq.SetPrimaryRoleRequestDTO
type SwaggerReplaceOverridesRequest = securityreq.ReplaceOverridesRequestDTO
type SwaggerOverrideItemRequest = securityreq.OverrideItemDTO
type SwaggerReplaceRolePermissionsRequest = securityreq.ReplaceRolePermissionsRequestDTO
type SwaggerCreateRoleRequest = securityreq.CreateRoleRequestDTO
type SwaggerUpdateRoleRequest = securityreq.UpdateRoleRequestDTO

type SwaggerCreateCategoryRequest = catalogreq.CreateCategoryRequestDTO
type SwaggerUpdateCategoryRequest = catalogreq.UpdateCategoryRequestDTO
type SwaggerCreateProductRequest = catalogreq.CreateProductRequestDTO
type SwaggerUpdateProductRequest = catalogreq.UpdateProductRequestDTO

type SwaggerCreateCustomerRequest = customerreq.CreateCustomerRequestDTO
type SwaggerUpdateCustomerRequest = customerreq.UpdateCustomerRequestDTO

type SwaggerCreateSaleRequest = salesreq.CreateSaleRequestDTO
type SwaggerRegisterSalePaymentRequest = salesreq.RegisterSalePaymentRequestDTO

type SwaggerCreateWorkOrderRequest = workorderreq.CreateWorkOrderRequestDTO
type SwaggerUpdateWorkOrderRequest = workorderreq.UpdateWorkOrderRequestDTO

type SwaggerCreateAppointmentRequest = appointmentreq.CreateAppointmentRequestDTO
type SwaggerUpdateAppointmentRequest = appointmentreq.UpdateAppointmentRequestDTO

type SwaggerFileUploadLinkRequest = filesreq.FileUploadLinkDTO

type SwaggerCreateCashCategoryRequest = cashreq.CreateCashCategoryRequestDTO
type SwaggerUpdateCashCategoryRequest = cashreq.UpdateCashCategoryRequestDTO
type SwaggerCreateCashMovementRequest = cashreq.CreateCashMovementRequestDTO
type SwaggerUpdateCashMovementRequest = cashreq.UpdateCashMovementRequestDTO

type SwaggerCreateWorkerRequest = workerreq.CreateWorkerRequestDTO
type SwaggerUpdateWorkerRequest = workerreq.UpdateWorkerRequestDTO
type SwaggerPayCommissionRequest = workerreq.PayCommissionRequestDTO
type SwaggerPaySalaryRequest = workerreq.PaySalaryRequestDTO

type SwaggerSetSettingRequest = systemreq.SetAppSettingRequestDTO
type SwaggerRoleRequest = entities.Role
type SwaggerMenuRequest = entities.Menu
type SwaggerSubMenuRequest = entities.SubMenu
type SwaggerPermissionRequest = entities.Permission

type SwaggerMeta struct {
	Total int64 `json:"total"`
}

type SwaggerErrorResponse = webutils.ProblemDetails

type SwaggerOKFlag struct {
	OK bool `json:"ok"`
}

type SwaggerOKFlagResponse struct {
	Data SwaggerOKFlag `json:"data"`
}

type SwaggerHealthData struct {
	Status  string `json:"status"`
	TimeUtc string `json:"timeUtc"`
	Version string `json:"version"`
}

type SwaggerHealthResponse struct {
	Data SwaggerHealthData `json:"data"`
}

type SwaggerSystemInfoData struct {
	Name    string `json:"name"`
	TimeUtc string `json:"timeUtc"`
}

type SwaggerSystemInfoResponse struct {
	Data SwaggerSystemInfoData `json:"data"`
}

type SwaggerBackupData struct {
	Status string `json:"status"`
}

type SwaggerBackupResponse struct {
	Data SwaggerBackupData `json:"data"`
}

type SwaggerAuthUser struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
}

type SwaggerAuthLoginData struct {
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	ExpiresIn    int64           `json:"expiresIn"`
	User         SwaggerAuthUser `json:"user"`
}

type SwaggerAuthLoginResponse struct {
	Data SwaggerAuthLoginData `json:"data"`
}

type SwaggerAuthMeResponse struct {
	Data authresp.MeResponseDTO `json:"data"`
}

type SwaggerPermissionTreeResponse struct {
	Data []dtos.EffectivePermissionNode `json:"data"`
}

type SwaggerUserListResponse struct {
	Data []securityresp.UserResponseDTO `json:"data"`
	Meta SwaggerMeta                    `json:"meta"`
}

type SwaggerUserResponse struct {
	Data securityresp.UserResponseDTO `json:"data"`
}

type SwaggerRoleListResponse struct {
	Data []securityresp.RoleResponseDTO `json:"data"`
	Meta SwaggerMeta                    `json:"meta"`
}

type SwaggerRoleResponse struct {
	Data entities.Role `json:"data"`
}

type SwaggerRoleDetailResponse struct {
	Data securityresp.RoleResponseDTO `json:"data"`
}

type SwaggerRolePermissionsResponse struct {
	Data []entities.RoleSubMenuPermission `json:"data"`
}

type SwaggerMenuListResponse struct {
	Data []entities.Menu `json:"data"`
	Meta SwaggerMeta     `json:"meta"`
}

type SwaggerMenuResponse struct {
	Data entities.Menu `json:"data"`
}

type SwaggerSubMenuListResponse struct {
	Data []entities.SubMenu `json:"data"`
	Meta SwaggerMeta        `json:"meta"`
}

type SwaggerSubMenuResponse struct {
	Data entities.SubMenu `json:"data"`
}

type SwaggerPermissionListResponse struct {
	Data []entities.Permission `json:"data"`
	Meta SwaggerMeta           `json:"meta"`
}

type SwaggerPermissionResponse struct {
	Data entities.Permission `json:"data"`
}

type SwaggerOverrideListResponse struct {
	Data []entities.UserPermissionOverride `json:"data"`
}

type SwaggerOverrideResponse struct {
	Data entities.UserPermissionOverride `json:"data"`
}

type SwaggerUserRolesResponse struct {
	Data []entities.UserRole `json:"data"`
}

type SwaggerAuditListResponse struct {
	Data []entities.AuditLog `json:"data"`
	Meta SwaggerMeta         `json:"meta"`
}

type SwaggerAuditResponse struct {
	Data entities.AuditLog `json:"data"`
}

type SwaggerCategoryListResponse struct {
	Data []entities.Category `json:"data"`
	Meta SwaggerMeta         `json:"meta"`
}

type SwaggerCategoryResponse struct {
	Data entities.Category `json:"data"`
}

type SwaggerProductListResponse struct {
	Data []entities.Product `json:"data"`
	Meta SwaggerMeta        `json:"meta"`
}

type SwaggerProductResponse struct {
	Data entities.Product `json:"data"`
}

type SwaggerCustomerListResponse struct {
	Data []entities.Customer `json:"data"`
	Meta SwaggerMeta         `json:"meta"`
}

type SwaggerCustomerResponse struct {
	Data entities.Customer `json:"data"`
}

type SwaggerSaleListResponse struct {
	Data []entities.Sale `json:"data"`
	Meta SwaggerMeta     `json:"meta"`
}

type SwaggerSaleResponse struct {
	Data entities.Sale `json:"data"`
}

type SwaggerSalePaymentResponse struct {
	Data entities.SalePayment `json:"data"`
}

type SwaggerSaleDetailData struct {
	Sale     entities.Sale          `json:"sale"`
	Items    []entities.SaleItem    `json:"items"`
	Payments []entities.SalePayment `json:"payments"`
}

type SwaggerSaleDetailResponse struct {
	Data SwaggerSaleDetailData `json:"data"`
}

type SwaggerWorkOrderListResponse struct {
	Data []entities.WorkOrder `json:"data"`
	Meta SwaggerMeta          `json:"meta"`
}

type SwaggerWorkOrderResponse struct {
	Data entities.WorkOrder `json:"data"`
}

type SwaggerAppointmentListResponse struct {
	Data []entities.Appointment `json:"data"`
	Meta SwaggerMeta            `json:"meta"`
}

type SwaggerAppointmentResponse struct {
	Data entities.Appointment `json:"data"`
}

type SwaggerFileListResponse struct {
	Data []entities.File `json:"data"`
	Meta SwaggerMeta     `json:"meta"`
}

type SwaggerFileResponse struct {
	Data entities.File `json:"data"`
}

type SwaggerCashCategoryListResponse struct {
	Data []entities.CashCategory `json:"data"`
	Meta SwaggerMeta             `json:"meta"`
}

type SwaggerCashCategoryResponse struct {
	Data entities.CashCategory `json:"data"`
}

type SwaggerCashMovementListResponse struct {
	Data []entities.CashMovement `json:"data"`
	Meta SwaggerMeta             `json:"meta"`
}

type SwaggerCashMovementResponse struct {
	Data entities.CashMovement `json:"data"`
}

type SwaggerWorkerListResponse struct {
	Data []entities.Worker `json:"data"`
	Meta SwaggerMeta       `json:"meta"`
}

type SwaggerWorkerResponse struct {
	Data entities.Worker `json:"data"`
}

type SwaggerSettingResponse struct {
	Data entities.AppSetting `json:"data"`
}
