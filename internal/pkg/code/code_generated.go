// Copyright 2024 Benjamin <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Code generated by "codegen -type=int /Users/max/Documents/coding/Backend/Golang/Personal/Luna/internal/pkg/code"; DO NOT EDIT.

package code

import "github.com/lunarianss/Luna/pkg/errors" // init register error codes defines in this source code to `github.com/lunarianss/Luna/pkg/errors
func init() {
	errors.Enroll(ErrAppMapMode, 500, "Error occurred while attempt to index from appTemplate using mode")
	errors.Enroll(ErrAppNotFoundRelatedConfig, 500, "Error occurred while attempt to find app related config")
	errors.Enroll(ErrEmailCode, 500, "Error occurred when email code is incorrect")
	errors.Enroll(ErrTokenEmail, 500, "Error occurred when email is incorrect")
	errors.Enroll(ErrTenantAlreadyExist, 500, "Error occurred when tenant is already exist")
	errors.Enroll(ErrSuccess, 200, "OK")
	errors.Enroll(ErrUnknown, 500, "Internal server error")
	errors.Enroll(ErrBind, 400, "Error occurred while binding the request body to the struct")
	errors.Enroll(ErrValidation, 400, "Validation failed")
	errors.Enroll(ErrPageNotFound, 404, "Page not found")
	errors.Enroll(ErrRestfulId, 400, "Error occurred while parse restful id from url")
	errors.Enroll(ErrRunTimeCaller, 500, "Error occurred while call go inner function")
	errors.Enroll(ErrRunTimeConfig, 500, "Error occurred while runtime config is nil")
	errors.Enroll(ErrDatabase, 500, "Database error")
	errors.Enroll(ErrRecordNotFound, 500, "Database record not found")
	errors.Enroll(ErrScanToField, 400, "Database scan error to field")
	errors.Enroll(ErrEncrypt, 401, "Error occurred while encrypting the user password")
	errors.Enroll(ErrSignatureInvalid, 401, "Signature is invalid")
	errors.Enroll(ErrExpired, 401, "Token expired")
	errors.Enroll(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	errors.Enroll(ErrMissingHeader, 401, "The `Authorization` header was empty")
	errors.Enroll(ErrPasswordIncorrect, 401, "Password was incorrect")
	errors.Enroll(ErrPermissionDenied, 403, "Permission denied")
	errors.Enroll(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	errors.Enroll(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	errors.Enroll(ErrInvalidJSON, 500, "Data is not valid JSON")
	errors.Enroll(ErrEncodingJSON, 500, "JSON data could not be encoded")
	errors.Enroll(ErrDecodingJSON, 500, "JSON data could not be decoded")
	errors.Enroll(ErrInvalidYaml, 500, "Data is not valid Yaml")
	errors.Enroll(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	errors.Enroll(ErrDecodingYaml, 500, "Yaml data could not be decoded")
	errors.Enroll(ErrTokenGenerate, 500, "Error occurred when Token generate")
	errors.Enroll(ErrTokenExpired, 500, "Error occurred when Token expired")
	errors.Enroll(ErrTokenInvalid, 401, "Token invalid")
	errors.Enroll(ErrTokenMethodErr, 500, "Unexpected signing method")
	errors.Enroll(ErrTokenInsNotFound, 500, "Jwt instance is not found")
	errors.Enroll(ErrRedisSetKey, 500, "Error occurred when set key, value to redis")
	errors.Enroll(ErrRedisSetExpire, 500, "Error occurred when set expire  to redis")
	errors.Enroll(ErrRedisRuntime, 500, "Error occurred when invoke redis api")
	errors.Enroll(ErrRedisDataExpire, 500, "Error occurred when redis data is expired")
	errors.Enroll(ErrOnlyOverrideConfigInDebugger, 500, "Error occurred while attempt to override config in non-debug mode")
	errors.Enroll(ErrModelEmptyInConfig, 500, "Error occurred while attempt to index model from config")
	errors.Enroll(ErrRequiredCorrectProvider, 500, "Error occurred when provider is not found or provider isn't include in the provider list")
	errors.Enroll(ErrRequiredModelName, 500, "Error occurred when model name is not found in model config")
	errors.Enroll(ErrRequiredCorrectModel, 500, "Error occurred when model is not found or model isn't include in the model list")
	errors.Enroll(ErrRequiredOverrideConfig, 500, "Config_from is ARGS that override_config_dict is required")
	errors.Enroll(ErrNotFoundModelRegistry, 500, "Model registry is not found in the model registry list")
	errors.Enroll(ErrRequiredPromptType, 500, "Prompt type is required when convert to prompt template")
	errors.Enroll(ErrProviderMapModel, 500, "Error occurred while attempt to index from providerMpa using provider")
	errors.Enroll(ErrProviderNotHaveIcon, 500, "Error occurred while provider entity doesn't have icon property")
	errors.Enroll(ErrToOriginModelType, 500, "Error occurred while convert to origin model type")
	errors.Enroll(ErrDefaultModelNotFound, 500, "Error occurred while trying to convert default model to unknown")
	errors.Enroll(ErrModelSchemaNotFound, 500, "Error occurred while attempt to index from predefined models using model name")
	errors.Enroll(ErrAllModelsEmpty, 500, "Error occurred when all models are empty")
	errors.Enroll(ErrModelNotHavePosition, 500, "Error occurred when models haven't position definition")
	errors.Enroll(ErrModelNotHaveEndPoint, 500, "Error occurred when models haven't url endpoint")
	errors.Enroll(ErrModelUrlNotConvertUrl, 500, "Error occurred when models url interface{} convert ot string ")
	errors.Enroll(ErrTypeOfPromptMessage, 500, "When prompt type is user, the type of message is neither string or []*promptMessageContent")
	errors.Enroll(ErrCallLargeLanguageModel, 500, "Error occurred when call large language model post api")
	errors.Enroll(ErrConvertDelimiterString, 500, "Error occurred when convert delimiter to string")
}
