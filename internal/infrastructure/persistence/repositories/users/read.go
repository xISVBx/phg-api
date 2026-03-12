package users

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.User, int64, error) {
	var out []entities.User
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"username":     "username",
		"fullName":     "full_name",
		"email":        "email",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.User{}, &out, opts, []string{"username", "full_name", "email"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var out entities.User
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var out entities.User
	if err := r.db.WithContext(ctx).First(&out, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entities.User, error) {
	var out entities.User
	value := strings.TrimSpace(usernameOrEmail)
	if err := r.db.WithContext(ctx).
		Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?)", value, value).
		First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
