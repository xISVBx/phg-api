package files

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	filesreq "photogallery/api_go/internal/application/dtos/request/files"
	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Upload(ctx context.Context, actor uuid.UUID, fh *multipart.FileHeader, link filesreq.FileUploadLinkDTO) (*entities.File, error) {
	relDir, customerID, storageKind, err := u.resolveRelativeDir(ctx, link)
	if err != nil {
		return nil, err
	}
	stored, err := u.storage.Store(ctx, fh, relDir)
	if err != nil {
		return nil, err
	}
	var created *entities.File
	err = u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		file := &entities.File{OriginalName: stored.OriginalName, ContentType: stored.ContentType, SizeBytes: stored.SizeBytes, UploadedByUserID: actor, UploadedAtUtc: time.Now().UTC(), StorageRelativePath: filepath.ToSlash(stored.StorageRelativePath), StorageKind: string(storageKind), StorageBaseKey: stored.StorageBaseKey}
		if err := repos.Files().Create(ctx, file); err != nil {
			return err
		}
		entityID, err := uuid.Parse(link.EntityID)
		if err != nil {
			return err
		}
		fl := &entities.FileLink{FileID: file.ID, EntityType: link.EntityType, EntityID: entityID, CustomerID: customerID, Notes: link.Notes}
		if err := repos.FileLinks().Create(ctx, fl); err != nil {
			return err
		}
		common.CreateAudit(ctx, repos, &actor, "File", file.ID.String(), "UPLOAD", file)
		common.CreateAudit(ctx, repos, &actor, "FileLink", fl.ID.String(), "LINK", fl)
		created = file
		return nil
	})
	if err != nil {
		_ = u.storage.Delete(ctx, stored.StorageRelativePath)
		return nil, err
	}
	return created, nil
}
