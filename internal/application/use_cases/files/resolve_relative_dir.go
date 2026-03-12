package files

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	filesreq "photogallery/api_go/internal/application/dtos/request/files"
	"photogallery/api_go/internal/application/mappers"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) resolveRelativeDir(ctx context.Context, link filesreq.FileUploadLinkDTO) (string, *uuid.UUID, enums.FileStorageKind, error) {
	now := time.Now().UTC()
	repos := u.uow.Repositories()
	entityID, err := uuid.Parse(link.EntityID)
	if err != nil {
		return "", nil, "", err
	}
	entityType := strings.ToLower(strings.TrimSpace(link.EntityType))
	var customer *entities.Customer
	switch entityType {
	case "customer":
		customer, err = repos.Customers().GetByID(ctx, entityID)
	case "sale":
		sale, e := repos.Sales().GetByID(ctx, entityID)
		if e != nil {
			err = e
			break
		}
		if sale.CustomerID != nil {
			customer, err = repos.Customers().GetByID(ctx, *sale.CustomerID)
		}
	case "workorder":
		wo, e := repos.WorkOrders().GetByID(ctx, entityID)
		if e != nil {
			err = e
			break
		}
		sale, e := repos.Sales().GetByID(ctx, wo.SaleID)
		if e != nil {
			err = e
			break
		}
		if sale.CustomerID != nil {
			customer, err = repos.Customers().GetByID(ctx, *sale.CustomerID)
		}
	case "saleitem":
		si, e := repos.SaleItems().GetByID(ctx, entityID)
		if e != nil {
			err = e
			break
		}
		sale, e := repos.Sales().GetByID(ctx, si.SaleID)
		if e != nil {
			err = e
			break
		}
		if sale.CustomerID != nil {
			customer, err = repos.Customers().GetByID(ctx, *sale.CustomerID)
		}
	case "appointment":
		appt, e := repos.Appointments().GetByID(ctx, entityID)
		if e != nil {
			err = e
			break
		}
		customer, err = repos.Customers().GetByID(ctx, appt.CustomerID)
	case "companyinternal":
		return fmt.Sprintf("empresa/interno/%04d/%02d", now.Year(), now.Month()), nil, enums.StorageCompanyInternal, nil
	default:
		return "", nil, "", errors.New("unsupported entityType")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, "", err
	}
	if customer == nil {
		return fmt.Sprintf("empresa/interno/%04d/%02d", now.Year(), now.Month()), nil, enums.StorageCompanyInternal, nil
	}
	key := mappers.CustomerKey(customer)
	base := "clientes/" + key
	cid := customer.ID
	switch entityType {
	case "sale":
		return filepath.ToSlash(base + "/ventas/" + entityID.String()), &cid, enums.StorageSale, nil
	case "saleitem":
		si, err := repos.SaleItems().GetByID(ctx, entityID)
		if err != nil {
			return "", nil, "", err
		}
		return filepath.ToSlash(base + "/ventas/" + si.SaleID.String()), &cid, enums.StorageSaleItem, nil
	case "workorder":
		return filepath.ToSlash(base + "/ordenes/" + entityID.String()), &cid, enums.StorageWorkOrder, nil
	case "appointment":
		return filepath.ToSlash(base + "/citas/" + entityID.String()), &cid, enums.StorageAppointment, nil
	case "customer":
		return filepath.ToSlash(base + "/general"), &cid, enums.StorageCustomerGeneral, nil
	default:
		return filepath.ToSlash(base + "/general"), &cid, enums.StorageCustomerGeneral, nil
	}
}
