package enums

type FileStorageKind string

const (
	StorageCustomerGeneral FileStorageKind = "CustomerGeneral"
	StorageSale            FileStorageKind = "Sale"
	StorageWorkOrder       FileStorageKind = "WorkOrder"
	StorageSaleItem        FileStorageKind = "SaleItem"
	StorageAppointment     FileStorageKind = "Appointment"
	StorageCompanyInternal FileStorageKind = "CompanyInternal"
	StorageCompanyFinance  FileStorageKind = "CompanyFinance"
	StorageOther           FileStorageKind = "Other"
)
