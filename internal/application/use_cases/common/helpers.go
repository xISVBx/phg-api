package common

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func QueryOpts(page, pageSize int, q, sort, dir string) drepo.QueryOptions {
	return drepo.QueryOptions{Page: page, PageSize: pageSize, Q: q, Sort: sort, Dir: dir}
}

func CreateAudit(ctx context.Context, repos appif.RepositorySet, actorID *uuid.UUID, entityType, entityID, action string, data any) {
	payload := "{}"
	if data != nil {
		if b, err := json.Marshal(data); err == nil {
			payload = string(b)
		}
	}
	_ = repos.AuditLogs().Create(ctx, &entities.AuditLog{ActorUserID: actorID, EntityType: entityType, EntityID: entityID, Action: action, DataJSON: payload})
}
