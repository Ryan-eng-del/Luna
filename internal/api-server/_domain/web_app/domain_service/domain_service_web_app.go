package domain_service

import (
	"context"

	po_entity_web_app "github.com/lunarianss/Luna/internal/api-server/_domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/_domain/web_app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/_domain/web_app/repository"
)

type WebAppDomain struct {
	WebAppRepo repository.WebAppRepo
}

func NewWebAppDomain(webAppRepo repository.WebAppRepo) *WebAppDomain {
	return &WebAppDomain{
		WebAppRepo: webAppRepo,
	}
}

func (ard *WebAppDomain) CreateEndUser(ctx context.Context, appRecord *po_entity_web_app.App) (*po_entity.EndUser, error) {
	sessionID, err := ard.WebAppRepo.GenerateSessionForEndUser(ctx)

	if err != nil {
		return nil, err
	}

	endUser := &po_entity.EndUser{
		TenantID:  appRecord.TenantID,
		AppID:     appRecord.ID,
		Type:      "browser",
		SessionID: sessionID,
	}

	endUserRecord, err := ard.WebAppRepo.CreateEndUser(ctx, endUser, nil)

	if err != nil {
		return nil, err
	}

	return endUserRecord, nil
}
