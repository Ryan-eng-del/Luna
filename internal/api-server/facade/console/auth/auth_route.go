// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/auth"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/email"
	"github.com/lunarianss/Luna/internal/infrastructure/mq"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
)

type AuthRoutes struct{}

func (a *AuthRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	redisIns, err := redis.GetRedisIns(nil)

	if err != nil {
		return err
	}

	mqProducerIns, err := mq.GetMQProducerIns(nil)

	if err != nil {
		return err
	}

	email, err := email.GetEmailSMTPIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// repo
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)

	// domain
	accountDomain := domain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	tenantDomain := domain.NewTenantDomain(tenantRepo)

	// service
	accountService := service.NewAccountService(accountDomain, tenantDomain, gormIns, mqProducerIns)
	accountController := controller.NewAuthController(accountService)

	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.POST("/email-code-login", accountController.SendEmailCode)
	authV1.POST("/email-code-login/validity", accountController.EmailValidity)
	authV1.POST("/refresh-token", accountController.RefreshToken)
	return nil
}

func (r *AuthRoutes) GetModule() string {
	return "auth"
}
