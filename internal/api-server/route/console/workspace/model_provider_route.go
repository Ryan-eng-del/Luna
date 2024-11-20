// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"

	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/model-provider/provider"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	"github.com/lunarianss/Luna/internal/api-server/domain/account"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type ModelProviderRoutes struct{}

func (r *ModelProviderRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	modelProviderDao := dao.NewModelProvider(gormIns)
	modelDao := dao.NewModelDao(gormIns)
	accountDao := dao.NewAccountDao(gormIns)
	tenantDao := dao.NewTenantDao(gormIns)

	// domain
	modelProviderDomain := domain.NewModelProviderDomain(modelProviderDao, modelDao)
	accountDomain := account.NewAccountDomain(accountDao, nil, nil, nil, tenantDao)

	// service
	modelProviderService := service.NewModelProviderService(modelProviderDomain, accountDomain)
	modelProviderController := controller.NewModelProviderController(modelProviderService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api/workspaces/current")

	modelProviderV1.Use(middleware.TokenAuthMiddleware())

	modelProviderV1.GET("/model-providers", modelProviderController.List)
	modelProviderV1.GET("/model-providers/:provider/:iconType/:lang", modelProviderController.ListIcons)
	modelProviderV1.POST("/model-providers/:provider", modelProviderController.SaveProviderCredential)

	return nil
}

func (r *ModelProviderRoutes) GetModule() string {
	return "providers"
}
