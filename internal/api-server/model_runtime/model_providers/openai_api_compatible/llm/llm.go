// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_runner"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/shopspring/decimal"
)

type IOpenApiCompactLargeLanguage interface {
	Invoke(ctx context.Context)
}

type openApiCompactLargeLanguageModel struct {
	*app_chat_runner.AppBaseChatRunner
	*biz_entity_chat.StreamGenerateQueue
	biz_entity.IAIModelRuntime
	FullAssistantContent string
	Usage                interface{}
	ChunkIndex           int
	Delimiter            string
	Model                string
	Stream               bool
	User                 string
	Stop                 []string
	Credentials          map[string]interface{}
	PromptMessages       []*po_entity_chat.PromptMessage
	ModelParameters      map[string]interface{}
}

func NewOpenApiCompactLargeLanguageModel(promptMessages []*po_entity_chat.PromptMessage, modelParameters map[string]interface{}, credentials map[string]interface{}, queueManager *biz_entity_chat.StreamGenerateQueue, model string, stream bool, modelRuntime biz_entity.IAIModelRuntime) *openApiCompactLargeLanguageModel {
	return &openApiCompactLargeLanguageModel{
		PromptMessages:      promptMessages,
		Credentials:         credentials,
		ModelParameters:     modelParameters,
		StreamGenerateQueue: queueManager,
		Model:               model,
		IAIModelRuntime:     modelRuntime,
		Stream:              stream,
		AppBaseChatRunner:   app_chat_runner.NewAppBaseChatRunner(),
	}
}

func (m *openApiCompactLargeLanguageModel) Invoke(ctx context.Context) {
	m.generate(ctx)
}

func (m *openApiCompactLargeLanguageModel) generate(ctx context.Context) {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Accept-Charset": "utf-8",
	}

	if extraHeaders, ok := m.Credentials["extra_headers"]; ok {
		if extraHeadersMap, ok := extraHeaders.(map[string]string); ok {
			for k, v := range extraHeadersMap {
				if _, ok := headers[k]; !ok {
					headers[k] = v
				}
			}
		}
	}

	if apiKey, ok := m.Credentials["api_key"]; ok {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	}

	endpointUrl, ok := m.Credentials["endpoint_url"]

	if !ok || endpointUrl == "" {
		m.PushErr(errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("Model %s not have endpoint url", m.Model)))
		return
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {
		m.PushErr(errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("Model %s not have endpoint url", m.Model)))
		return
	}

	if !strings.HasSuffix(endpointUrlStr, "/") {
		endpointUrlStr = fmt.Sprintf("%s/", endpointUrlStr)
	}

	requestData := map[string]interface{}{
		"model":  m.Model,
		"stream": m.Stream,
	}

	for k, v := range m.ModelParameters {
		requestData[k] = v
	}
	messageItems := make([]map[string]interface{}, 0)

	completionType := m.Credentials["mode"]
	if completionType == string(po_entity.CHAT) {
		endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "chat/completions")

		if err != nil {
			m.PushErr(errors.WithCode(code.ErrRunTimeCaller, err.Error()))
			return
		}
		endpointUrlStr = endpointJoinUrl

		for _, promptMessage := range m.PromptMessages {
			messageItem, err := promptMessage.ConvertToRequestData()

			if err != nil {
				m.PushErr(err)
				return
			}
			messageItems = append(messageItems, messageItem)
		}
	}

	requestData["messages"] = messageItems

	if len(m.Stop) > 1 {
		requestData["stop"] = m.Stop
	}

	if m.User != "" {
		requestData["user"] = m.User
	}

	client := http.Client{
		Timeout: time.Duration(300) * time.Second,
	}

	log.Infof("Invoke llm request body %+v", requestData)
	requestBodyData, err := json.Marshal(requestData)

	if err != nil {
		m.PushErr(errors.WithCode(code.ErrEncodingJSON, err.Error()))
		return
	}

	req, err := http.NewRequest("POST", endpointUrlStr, bytes.NewReader(requestBodyData))

	if err != nil {
		m.PushErr(errors.WithCode(code.ErrRunTimeCaller, err.Error()))
		return
	}

	if len(headers) > 0 {
		for headerKey, headerValue := range headers {
			req.Header.Set(headerKey, headerValue)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		m.PushErr(errors.WithCode(code.ErrCallLargeLanguageModel, err.Error()))
		return
	}

	defer response.Body.Close()

	if m.Stream {
		m.handleStreamResponse(ctx, response)
	}
}

func (m *openApiCompactLargeLanguageModel) sendStreamChunkToQueue(ctx context.Context, messageId string, assistantPromptMessage *biz_entity_chat.AssistantPromptMessage) {
	streamResultChunk := &biz_entity_chat.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &biz_entity_chat.LLMResultChunkDelta{
			Index:   m.ChunkIndex,
			Message: assistantPromptMessage,
		},
	}
	m.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, false, nil)
}

func (m *openApiCompactLargeLanguageModel) sendErrorChunkToQueue(ctx context.Context, code error) {
	defer m.Close()
	err := errors.WithMessage(code, fmt.Sprintf("Error ocurred when handle stream: %#+v", code))
	m.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err)
}

func handleTokenCount(count any) (int64, bool) {
	var countInt int64
	ok := true

	switch v := count.(type) {
	case float64:
		countInt = int64(v)
	case int:
		countInt = int64(v)
	case string:
		countFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			ok = false
			log.Infof("token count %v can't convert to float", count)
		}
		countInt = int64(countFloat)
	default:
		ok = false
	}
	return countInt, ok
}

func (m *openApiCompactLargeLanguageModel) sendStreamFinalChunkToQueue(ctx context.Context, messageId string, finalReason string, fullAssistant string, usage map[string]interface{}) {
	defer m.Close()

	var (
		err               error
		ok                bool
		completeTokensInt int64
		promptTokensInt   int64
	)

	promptTokens, found := usage["prompt_tokens"]

	if found {
		promptTokensInt, ok = handleTokenCount(promptTokens)
	}

	if !ok {
		// 使用 openai tokenizlier
		promptTokensInt = 10
	}

	completeTokens, found := usage["completion_tokens"]

	if found {
		completeTokensInt, ok = handleTokenCount(completeTokens)
	}

	if !ok {
		// 使用 openai tokenizlier
		completeTokensInt = 20
	}

	promptPriceInfo, err := m.GetPrice(m.Model, m.Credentials, biz_entity.INPUT, promptTokensInt)

	if err != nil {
		m.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err)
	}

	completePriceInfo, err := m.GetPrice(m.Model, m.Credentials, biz_entity.OUTPUT, completeTokensInt)

	promptTotal := decimal.NewFromFloat(promptPriceInfo.TotalAmount)
	completeTotal := decimal.NewFromFloat(completePriceInfo.TotalAmount)
	totalAmount := promptTotal.Add(completeTotal).InexactFloat64()

	if err != nil {
		m.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err)
	}

	llmUsage := &biz_entity_chat.LLMUsage{
		PromptTokens:        promptTokensInt,
		PromptUnitPrice:     promptPriceInfo.UnitPrice,
		PromptPriceUnit:     promptPriceInfo.Unit,
		PromptPrice:         promptPriceInfo.TotalAmount,
		CompletionTokens:    completeTokensInt,
		CompletionUnitPrice: completePriceInfo.UnitPrice,
		CompletionPriceUnit: completePriceInfo.Unit,
		CompletionPrice:     completePriceInfo.TotalAmount,
		Currency:            promptPriceInfo.Currency,
		Latency:             1.0,
		TotalTokens:         promptTokensInt + completeTokensInt,
		TotalPrice:          totalAmount,
	}

	streamResultChunk := &biz_entity_chat.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &biz_entity_chat.LLMResultChunkDelta{
			Index: m.ChunkIndex,
			Message: &biz_entity_chat.AssistantPromptMessage{
				PromptMessage: &po_entity_chat.PromptMessage{
					Content: fullAssistant,
				},
			},
			FinishReason: finalReason,
			Usage:        llmUsage,
		},
	}
	m.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, true, nil)
}

func (m *openApiCompactLargeLanguageModel) handleStreamResponse(ctx context.Context, response *http.Response) {

	var (
		messageID    string
		finishReason string
	)

	delimiter, ok := m.Credentials["stream_mode_delimiter"]

	if !ok {
		delimiter = "\n\n"
	}

	m.Delimiter, ok = delimiter.(string)

	if !ok {
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrConvertDelimiterString, fmt.Sprintf("Can't convert delimiter %+v to string", delimiter)))
		return
	}

	scanner := bufio.NewScanner(response.Body)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), m.Delimiter); i >= 0 {
			return i + len(m.Delimiter), data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	var usage = make(map[string]interface{})

	for scanner.Scan() {
		var (
			assistantPromptMessage *biz_entity_chat.AssistantPromptMessage
		)
		chunk := strings.TrimSpace(scanner.Text())

		if chunk == "" || strings.HasPrefix(chunk, ":") {
			continue
		}

		chunk = strings.TrimPrefix(chunk, "data: ")
		chunk = strings.TrimSpace(chunk)

		if chunk == "[DONE]" {
			continue
		}

		var chunkJson map[string]interface{}

		err := json.Unmarshal([]byte(chunk), &chunkJson)

		if err != nil {
			m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrEncodingJSON, fmt.Sprintf("JSON data %+v could not be decoded, failed: %+v", chunk, err.Error())))
			return
		}

		// groq 返回 error
		if apiError, ok := chunkJson["error"]; ok {
			if apiMapErr, ok := apiError.(map[string]interface{}); ok {
				if ok {
					apiByteErr, err := json.Marshal(apiMapErr)

					if err != nil {
						m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrEncodingJSON, err.Error()))
						return
					}

					m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrCallLargeLanguageModel, string(apiByteErr)))
					return
				}
			}
		}

		var chunkChoice = make(map[string]interface{})

		// groq
		if usageChunk, ok := chunkJson["x_groq"]; ok {
			if v, ok := usageChunk.(map[string]interface{}); ok {
				usageAny, ok := v["usage"]
				if ok {
					usageMap, ok := usageAny.(map[string]interface{})
					if ok {
						usage = usageMap
					}
				}
			}
		}

		if chunkChoices, ok := chunkJson["choices"]; ok {
			if v, ok := chunkChoices.([]interface{}); ok {
				if vv, ok := v[0].(map[string]interface{}); ok {
					chunkChoice = vv
				}
			}
		}

		messageID, ok = chunkChoice["id"].(string)

		if !ok {
			messageID = ""
		}

		finishReason, ok = chunkChoice["finish_reason"].(string)

		if !ok {
			finishReason = "Finish_reason doesn't convert to string"
		}

		m.ChunkIndex += 1

		if delta, ok := chunkChoice["delta"]; ok {
			if deltaMap, ok := delta.(map[string]interface{}); ok {
				deltaContent := deltaMap["content"]
				if deltaContentStr, ok := deltaContent.(string); ok {
					m.FullAssistantContent += deltaContentStr
					assistantPromptMessage = biz_entity_chat.NewAssistantPromptMessage(po_entity_chat.ASSISTANT, deltaContentStr)
				}
			}
		} else {
			log.Warn("This chunk not property of delta and text")
			continue
		}

		if assistantPromptMessage != nil {
			m.sendStreamChunkToQueue(ctx, messageID, assistantPromptMessage)
		}
	}

	if err := scanner.Err(); err != nil {
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrRunTimeCaller, err.Error()))
		return
	}

	m.sendStreamFinalChunkToQueue(ctx, messageID, finishReason, m.FullAssistantContent, usage)
}
