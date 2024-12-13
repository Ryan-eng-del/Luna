// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_chat_runner

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type appChatRunner struct {
	*AppBaseChatRunner
	AppDomain  *domain_service.AppDomain
	ChatDomain *chatDomain.ChatDomain
}

func NewAppChatRunner(appBaseChatRunner *AppBaseChatRunner, appDomain *domain_service.AppDomain, chatDomain *chatDomain.ChatDomain) *appChatRunner {
	return &appChatRunner{
		AppBaseChatRunner: appBaseChatRunner,
		AppDomain:         appDomain,
		ChatDomain:        chatDomain,
	}
}

func (r *appChatRunner) Run(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager *biz_entity.StreamGenerateQueue) {

	var (
		memory token_buffer_memory.ITokenBufferMemory
	)

	appRecord, err := r.AppDomain.AppRepo.GetAppByID(ctx, applicationGenerateEntity.AppConfig.AppID)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	credentials, err := applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.GetCurrentCredentials(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType, applicationGenerateEntity.AppConfig.Model.Model)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	if applicationGenerateEntity.ConversationID != "" {
		modelCaller := model_registry.NewModelRegisterCaller(applicationGenerateEntity.AppConfig.Model.Model, string(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType), applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider, credentials, applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance)

		memory = token_buffer_memory.NewTokenBufferMemory(conversation, modelCaller, r.ChatDomain)
	}

	promptMessages, stop, err := r.OrganizePromptMessage(ctx, appRecord, applicationGenerateEntity.ModelConf, applicationGenerateEntity.AppConfig.PromptTemplate, applicationGenerateEntity.Inputs, nil, applicationGenerateEntity.Query, "", memory)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	modelInstance := model_registry.NewModelRegisterCaller(applicationGenerateEntity.AppConfig.Model.Model, string(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType), applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider, credentials, applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance)

	if applicationGenerateEntity.Stream {
		modelInstance.InvokeLLM(ctx, promptMessages, queueManager, applicationGenerateEntity.ModelConf.Parameters, nil, stop, applicationGenerateEntity.UserID, nil)
	}
}
