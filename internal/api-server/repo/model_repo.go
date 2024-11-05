package repo

import (
	"context"

	model "github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ModelRepo interface {
	// GetTenantModels get all models by searchModel
	GetTenantModel(ctx context.Context, tenantId, providerName, modelName, modelType string) (*model.ProviderModel, error)
	// UpdateModel updates model
	UpdateModel(ctx context.Context, model *model.ProviderModel) error
	// CreateModel create model
	CreateModel(ctx context.Context, model *model.ProviderModel) error
}
