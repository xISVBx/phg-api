package files

import (
	"context"

	"github.com/google/uuid"
	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
)

func (u *UseCase) Delete(ctx context.Context, actor, fileID uuid.UUID) error {
	return u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		file, err := repos.Files().GetByID(ctx, fileID)
		if err != nil {
			return err
		}
		if err := repos.FileLinks().DeleteByFile(ctx, fileID); err != nil {
			return err
		}
		if err := repos.Files().Delete(ctx, fileID); err != nil {
			return err
		}
		if err := u.storage.Delete(ctx, file.StorageRelativePath); err != nil {
			return err
		}
		common.CreateAudit(ctx, repos, &actor, "File", fileID.String(), "DELETE", file)
		return nil
	})
}
