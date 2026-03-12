package files

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.File, int64, error) {
	var out []entities.File
	allowedSorts := map[string]string{
		"createdAtUtc":        "created_at_utc",
		"uploadedAtUtc":       "uploaded_at_utc",
		"originalName":        "original_name",
		"storageRelativePath": "storage_relative_path",
		"storageKind":         "storage_kind",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.File{}, &out, opts, []string{"original_name", "storage_relative_path", "storage_kind"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.File, error) {
	var out entities.File
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
