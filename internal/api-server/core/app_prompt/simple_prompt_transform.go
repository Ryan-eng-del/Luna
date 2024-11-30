// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package prompt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type SimplePromptTransform struct {
}

func (s *SimplePromptTransform) GetPromptStrAndRules(appMode po_entity.AppMode, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity, prePrompt string, inputs map[string]interface{}, query string, context string, histories string) (string, *biz_entity_app_config.SimpleChatPromptConfig, error) {

	var (
		variables = make(map[string]interface{})
	)

	promptTemplateConfig, err := s.GetPromptTemplate(appMode, modelConfig.Provider, modelConfig.Model, prePrompt, context != "", false, histories != "")

	if err != nil {
		return "", nil, err
	}

	for _, key := range promptTemplateConfig.CustomVariableKeys {
		variables[key] = inputs[key]
	}

	for _, v := range promptTemplateConfig.SpecialVariableKeys {
		if v == "#context" {
			variables["#context#"] = context
		} else if v == "#query#" {
			variables["#query#"] = query
		} else if v == "#histories#" {
			variables["#histories#"] = histories
		}
	}

	prompt := promptTemplateConfig.PromptTemplate.Format(inputs, false)

	return prompt, promptTemplateConfig.PromptRules, nil

}
func (s *SimplePromptTransform) GetChatModelPromptMessage(appMode po_entity.AppMode, prePrompt string, inputs map[string]interface{}, query string, context string, files []string, memory any, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity) ([]*po_entity_chat.PromptMessage, []string, error) {

	var promptMessages []*po_entity_chat.PromptMessage

	prompt, _, err := s.GetPromptStrAndRules(appMode, modelConfig, prePrompt, inputs, query, context, "")

	if err != nil {
		return nil, nil, err
	}

	if prompt != "" && query != "" {
		promptMessages = append(promptMessages, po_entity_chat.NewSystemMessage(prompt))
	}

	if query != "" {
		promptMessages = append(promptMessages, s.GetLastUserMessage(query, nil))
	} else {
		promptMessages = append(promptMessages, s.GetLastUserMessage(query, nil))
	}

	return promptMessages, nil, nil
}

func (s *SimplePromptTransform) GetLastUserMessage(prompt string, files []string) *po_entity_chat.PromptMessage {
	return po_entity_chat.NewUserMessage(prompt)
}

func (s *SimplePromptTransform) GetPrompt(appMode po_entity.AppMode, promptTemplateEntity *biz_entity_app_config.PromptTemplateEntity, inputs map[string]interface{}, query string, files []string, context string, memory any, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity) ([]*po_entity_chat.PromptMessage, []string, error) {

	var (
		promptMessage []*po_entity_chat.PromptMessage
		stop          []string
		err           error
	)

	modelMode := modelConfig.Mode

	if modelMode == "chat" {
		promptMessage, stop, err = s.GetChatModelPromptMessage(appMode, promptTemplateEntity.SimplePromptTemplate, inputs, query, context, files, nil, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessage, stop, nil
}

func (s *SimplePromptTransform) GetPromptTemplate(appMode po_entity.AppMode, provider, model, prePrompt string, hasContext bool, queryInPrompt bool, withMemoryPrompt bool) (*biz_entity_app_config.SimpleChatPromptTransformConfig, error) {

	var (
		customVariableKeys  []string
		specialVariableKeys []string
		prompt              string
		templatePromptRules string
		templateParser      biz_entity_app_config.IPromptTemplateParser
	)

	promptRules, err := s.getPromptRole(appMode, provider, model)

	if err != nil {
		return nil, err
	}

	promptOrders := promptRules.SystemPromptOrders

	for _, promptOrder := range promptOrders {
		if promptOrder == "context_prompt" && hasContext {
			prompt += promptRules.ContextPrompt
			specialVariableKeys = append(specialVariableKeys, "#context#")
		} else if promptOrder == "pre_prompt" && prePrompt != "" {
			prompt += prePrompt + "\n"
			templateParser = biz_entity_app_config.NewPromptTemplateParse(prePrompt, false)
			customVariableKeys = templateParser.Extract()
		} else if promptOrder == "histories_prompt" && withMemoryPrompt {
			prompt += promptRules.HistoriesPrompt
			specialVariableKeys = append(specialVariableKeys, "#histories#")
		}
	}

	if queryInPrompt {
		templatePromptRules = promptRules.QueryPrompt
		prompt += templatePromptRules
		specialVariableKeys = append(specialVariableKeys, "#query#")
	}

	templateParser = biz_entity_app_config.NewPromptTemplateParse(prompt, false)

	return &biz_entity_app_config.SimpleChatPromptTransformConfig{
		PromptTemplate:      templateParser,
		CustomVariableKeys:  customVariableKeys,
		SpecialVariableKeys: specialVariableKeys,
		PromptRules:         promptRules,
	}, nil
}

func (s *SimplePromptTransform) getPromptRole(appMode po_entity.AppMode, provider, modelName string) (*biz_entity_app_config.SimpleChatPromptConfig, error) {

	var (
		promptRoleMap biz_entity_app_config.SimpleChatPromptConfig
	)

	promptFileName := s.promptFileName(appMode, provider, modelName)

	_, fullFilePath, _, ok := runtime.Caller(0)

	if !ok {
		return nil, errors.WithCode(code.ErrRunTimeCaller, "Fail to get runtime caller info")
	}

	fileDir := filepath.Dir(fullFilePath)

	roleJsonPath := fmt.Sprintf("%s/prompt_templates/%s", fileDir, promptFileName)

	roleJsonContent, err := os.ReadFile(roleJsonPath)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, "Read file %s failed, Error: %+v", roleJsonPath, err)
	}

	if err := json.Unmarshal(roleJsonContent, &promptRoleMap); err != nil {
		return nil, errors.WithCode(code.ErrDecodingJSON, err.Error())
	}

	return &promptRoleMap, nil
}

func (s *SimplePromptTransform) promptFileName(appMode po_entity.AppMode, _, _ string) string {
	if appMode == po_entity.COMPLETION {
		return "common_completion.json"
	} else {
		return "common_chat.json"
	}
}
