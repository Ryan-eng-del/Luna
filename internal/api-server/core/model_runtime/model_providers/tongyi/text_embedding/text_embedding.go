package text_embedding

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/openai_api_compatible/text_embedding"
	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"

	biz_entity_openai_standard_response "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/openai_standard_response"

	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
)

type tongyiTextEmbedding struct {
	text_embedding.IOpenApiCompactTextEmbeddingModel
}

func init() {
	NewTongyiTextEmbedding().Register()
}

func NewTongyiTextEmbedding() *tongyiTextEmbedding {
	return &tongyiTextEmbedding{}
}

func (m *tongyiTextEmbedding) RegisterName() string {
	return "tongyi/text-embedding"
}

func (m *tongyiTextEmbedding) Register() {
	model_registry.TextEmbeddingRegistry.RegisterLargeModelInstance(m)
}

func (m *tongyiTextEmbedding) Embedding(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user string, modelRuntime biz_entity.IAIModelRuntime, inputType string, texts []string) (*biz_entity_openai_standard_response.TextEmbeddingResult, error) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactTextEmbeddingModel = text_embedding.NewOpenApiCompactLargeLanguageModel(ctx, model, credentials, texts, modelRuntime)
	return m.IOpenApiCompactTextEmbeddingModel.Invoke(ctx)

}

func (m *tongyiTextEmbedding) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["endpoint_url"] = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	credentials["api_key"] = credentials["dashscope_api_key"]
	return credentials
}
