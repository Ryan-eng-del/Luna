// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/openai_api_compatible/llm"
	provider_register "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	biz_entity_chat_prompt_message "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/chat_prompt_message"
	biz_entity_base_stream_generator "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/stream_base_generator"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
)

type zhipuLargeLanguageModel struct {
	llm.IOpenApiCompactLargeLanguage
}

func init() {
	NewZhipuLargeLanguageModel().Register()
}

func NewZhipuLargeLanguageModel() *zhipuLargeLanguageModel {
	return &zhipuLargeLanguageModel{}
}

var _ provider_register.IModelRegistry = (*zhipuLargeLanguageModel)(nil)

func (m *zhipuLargeLanguageModel) Invoke(ctx context.Context, queueManager biz_entity_base_stream_generator.IStreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []biz_entity_chat_prompt_message.IPromptMessage, modelRuntime biz_entity.IAIModelRuntime, tools []*biz_entity_chat_prompt_message.PromptMessageTool) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, model, modelRuntime, tools)
	m.IOpenApiCompactLargeLanguage.Invoke(ctx, queueManager)
}

func (m *zhipuLargeLanguageModel) InvokeNonStream(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []biz_entity_chat_prompt_message.IPromptMessage, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_base_stream_generator.LLMResult, error) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, model, modelRuntime, nil)
	return m.IOpenApiCompactLargeLanguage.InvokeNonStream(ctx)
}
func (m *zhipuLargeLanguageModel) Register() {
	provider_register.ModelRuntimeRegistry.RegisterLargeModelInstance(m)
}

func (m *zhipuLargeLanguageModel) RegisterName() string {
	return "zhipuai/llm"
}

func (m *zhipuLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://open.bigmodel.cn/api/paas/v4"
	return credentials
}
