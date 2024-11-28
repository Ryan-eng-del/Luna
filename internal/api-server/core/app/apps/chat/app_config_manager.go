package chat

import (
	"context"
	"encoding/json"

	"github.com/lunarianss/Luna/internal/api-server/_domain/app/entity/po_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/_domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/model_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/prompt_template"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ChatAppConfigManager struct {
	ProviderDomain *domain_service.ProviderDomain
}

func NewChatAppConfigManager(providerDomain *domain_service.ProviderDomain) *ChatAppConfigManager {
	return &ChatAppConfigManager{
		ProviderDomain: providerDomain,
	}
}

func (m *ChatAppConfigManager) ConfigValidate(ctx context.Context, tenantID string, config map[string]any) (map[string]any, error) {
	var (
		relatedConfigKeys        []string
		currentRelatedConfigKeys []string
	)

	// model
	modelConfigManager := model_config.NewModelConfigManager(m.ProviderDomain)
	config, currentRelatedConfigKeys, err := modelConfigManager.ValidateAndSetDefaults(ctx, tenantID, config)

	if err != nil {
		return nil, err
	}

	relatedConfigKeys = append(relatedConfigKeys, currentRelatedConfigKeys...)

	// todo Filter out extra parameters
	return config, nil
}

func (m *ChatAppConfigManager) getAppConfig(ctx context.Context, appModel *po_entity.App, appModelConfig *po_entity.AppModelConfig, conversation *po_entity_chat.Conversation, overrideConfigDict map[string]any) (*app_config.ChatAppConfig, error) {

	var (
		configFrom app_config.EasyUIBasedAppModelConfigFrom
		configDict map[string]interface{}
	)

	if overrideConfigDict != nil {
		configFrom = app_config.Args
	} else if conversation != nil {
		configFrom = app_config.ConversationSpecificConfig
	} else {
		configFrom = app_config.AppLatestConfig
	}

	if configFrom != app_config.Args {
		appModelByte, err := json.Marshal(appModelConfig)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(appModelByte, &configDict); err != nil {
			return nil, err
		}
	} else {
		if overrideConfigDict == nil {
			return nil, errors.WithCode(code.ErrRequiredOverrideConfig, "override_config_dict is required when config_from is ARGS")
		}

		configDict = overrideConfigDict
	}

	modelConfigManager := model_config.NewModelConfigManager(m.ProviderDomain)
	modelConfigEntity, err := modelConfigManager.Convert(ctx, configDict)

	if err != nil {
		return nil, err
	}

	promptTemplateConfigManager := prompt_template.PromptTemplateConfigManager{}

	promptTemplate, err := promptTemplateConfigManager.Convert(configDict)

	if err != nil {
		return nil, err
	}

	return &app_config.ChatAppConfig{
		EasyUIBasedAppConfig: &app_config.EasyUIBasedAppConfig{
			AppConfig: &app_config.AppConfig{
				TenantID:               appModel.TenantID,
				AppID:                  appModel.ID,
				AppMode:                appModel.Mode,
				SensitiveWordAvoidance: nil,
				AdditionalFeatures:     nil,
			},
			AppModelConfigDict: configDict,
			AppModelConfigFrom: configFrom,
			AppModelConfigID:   appModelConfig.ID,
			Model:              modelConfigEntity,
			PromptTemplate:     promptTemplate,
			Dataset:            nil,
		},
	}, nil
}
