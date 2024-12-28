// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_chat_runner

import (
	"context"
	"time"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	"github.com/lunarianss/Luna/internal/api-server/core/app_feature"

	"github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	"github.com/redis/go-redis/v9"
)

type appChatRunner struct {
	*AppBaseChatRunner
	AppDomain      *domain_service.AppDomain
	ChatDomain     *chatDomain.ChatDomain
	ProviderDomain *providerDomain.ProviderDomain
	DatasetDomain  *datasetDomain.DatasetDomain
	redis          *redis.Client
}

func NewAppChatRunner(appBaseChatRunner *AppBaseChatRunner, appDomain *domain_service.AppDomain, chatDomain *chatDomain.ChatDomain, providerDomain *providerDomain.ProviderDomain, datasetDomain *datasetDomain.DatasetDomain, redis *redis.Client) *appChatRunner {
	return &appChatRunner{
		AppBaseChatRunner: appBaseChatRunner,
		AppDomain:         appDomain,
		ChatDomain:        chatDomain,
		ProviderDomain:    providerDomain,
		DatasetDomain:     datasetDomain,
		redis:             redis,
	}
}

func (r *appChatRunner) baseRun(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, conversation *po_entity_chat.Conversation) (model_registry.IModelRegistryCall, []*po_entity_chat.PromptMessage, []string, *po_entity.App, error) {

	var (
		memory token_buffer_memory.ITokenBufferMemory
	)

	appRecord, err := r.AppDomain.AppRepo.GetAppByID(ctx, applicationGenerateEntity.AppConfig.AppID)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	credentials, err := applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.GetCurrentCredentials(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType, applicationGenerateEntity.AppConfig.Model.Model)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	if applicationGenerateEntity.ConversationID != "" {
		modelCaller := model_registry.NewModelRegisterCaller(applicationGenerateEntity.AppConfig.Model.Model, string(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType), applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider, credentials, applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance)

		memory = token_buffer_memory.NewTokenBufferMemory(conversation, modelCaller, r.ChatDomain)
	}

	promptMessages, stop, err := r.OrganizePromptMessage(ctx, appRecord, applicationGenerateEntity.ModelConf, applicationGenerateEntity.AppConfig.PromptTemplate, applicationGenerateEntity.Inputs, nil, applicationGenerateEntity.Query, "", memory)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	modelInstance := model_registry.NewModelRegisterCaller(applicationGenerateEntity.AppConfig.Model.Model, string(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType), applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider, credentials, applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance)

	return modelInstance, promptMessages, stop, appRecord, nil
}

func (r *appChatRunner) DirectOutStream(applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager *biz_entity.StreamGenerateQueue, text string, promptMessages []*po_entity_chat.PromptMessage) {

	index := 0
	for i, token := range text {
		tokenStr := string(token)
		llmResultChunk := &biz_entity.LLMResultChunk{
			Model:         applicationGenerateEntity.ModelConf.Model,
			PromptMessage: promptMessages,
			Delta: &biz_entity.LLMResultChunkDelta{
				Index:   index,
				Message: biz_entity.NewAssistantPromptMessage(tokenStr),
			},
		}
		event := biz_entity.NewAppQueueEvent(biz_entity.LLMChunk)
		queueManager.Push(&biz_entity.QueueLLMChunkEvent{
			AppQueueEvent: event,
			Chunk:         llmResultChunk})
		i++
		time.Sleep(10 * time.Millisecond)
	}

	queueManager.FinalManual(&biz_entity.QueueMessageEndEvent{
		AppQueueEvent: biz_entity.NewAppQueueEvent(biz_entity.MessageEnd),
		LLMResult: &biz_entity.LLMResult{
			Model:         applicationGenerateEntity.ModelConf.Model,
			PromptMessage: promptMessages,
			Message:       biz_entity.NewAssistantPromptMessage(text),
			Usage:         biz_entity.NewEmptyLLMUsage(),
		},
	})
}

func (r *appChatRunner) Run(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager *biz_entity.StreamGenerateQueue) {

	modelInstance, promptMessages, stop, app, err := r.baseRun(ctx, applicationGenerateEntity, conversation)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	if applicationGenerateEntity.Query != "" {
		annotation, err := r.QueryAppAnnotationToReply(ctx, app, message, applicationGenerateEntity.Query, applicationGenerateEntity.UserID, string(applicationGenerateEntity.InvokeFrom))

		if err != nil {
			queueManager.PushErr(err)
			return
		}

		if annotation != nil {
			queueEvent := biz_entity.NewAppQueueEvent(biz_entity.AnnotationReply)
			queueManager.Push(&biz_entity.QueueAnnotationReplyEvent{
				AppQueueEvent:       queueEvent,
				MessageAnnotationID: annotation.ID,
			})

			r.DirectOutStream(applicationGenerateEntity, message, conversation, queueManager, annotation.Content, promptMessages)
			return
		}
	}

	modelInstance.InvokeLLM(ctx, promptMessages, queueManager, applicationGenerateEntity.ModelConf.Parameters, nil, stop, applicationGenerateEntity.UserID, nil)
}

func (r *appChatRunner) RunNonStream(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation) (*biz_entity.LLMResult, error) {

	modelInstance, promptMessages, stop, app, err := r.baseRun(ctx, applicationGenerateEntity, conversation)

	if err != nil {
		return nil, err
	}

	if applicationGenerateEntity.Query != "" {
		annotation, err := r.QueryAppAnnotationToReply(ctx, app, message, applicationGenerateEntity.Query, applicationGenerateEntity.UserID, string(applicationGenerateEntity.InvokeFrom))

		if err != nil {
			return nil, err
		}

		if annotation != nil {
			return &biz_entity.LLMResult{
				Model:         applicationGenerateEntity.ModelConf.Model,
				PromptMessage: promptMessages,
				Message:       biz_entity.NewAssistantPromptMessage(annotation.Content),
				Usage:         biz_entity.NewEmptyLLMUsage(),
			}, nil
		}
	}

	return modelInstance.InvokeLLMNonStream(ctx, promptMessages, applicationGenerateEntity.ModelConf.Parameters, nil, stop, applicationGenerateEntity.UserID, nil)
}

func (r *appChatRunner) QueryAppAnnotationToReply(ctx context.Context, appRecord *po_entity.App, message *po_entity_chat.Message, query, accountID, invokeFrom string) (*po_entity_chat.MessageAnnotation, error) {
	return app_feature.NewAnnotationReplyFeature(r.ChatDomain, r.DatasetDomain, r.ProviderDomain, r.redis).Query(ctx, appRecord, message, query, accountID, invokeFrom)
}
