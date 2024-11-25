package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/app"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat"
	modelDomain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type AppRoutes struct{}

func (a *AppRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// dao
	appDao := dao.NewAppDao(gormIns)
	modelDao := dao.NewModelDao(gormIns)
	providerDao := dao.NewModelProvider(gormIns)
	accountDao := dao.NewAccountDao(gormIns)
	tenantDao := dao.NewTenantDao(gormIns)
	appRunningDao := dao.NewAppRunningDao(gormIns)
	messageDao := dao.NewMessageDao(gormIns)

	// domain
	appDomain := domain.NewAppDomain(appDao, appRunningDao, messageDao)
	modelDomain := modelDomain.NewModelDomain(modelDao)
	providerDomain := providerDomain.NewModelProviderDomain(providerDao, modelDao)
	accountDomain := accountDomain.NewAccountDomain(accountDao, nil, nil, nil, tenantDao)
	chatDomain := chatDomain.NewChatDomain(messageDao)

	// service
	appService := service.NewAppService(appDomain, modelDomain, providerDomain, accountDomain, gormIns, config)
	chatService := service.NewChatService(appDomain, providerDomain, accountDomain, chatDomain)

	appController := controller.NewAppController(appService, chatService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	modelProviderV1.Use(middleware.TokenAuthMiddleware())
	modelProviderV1.POST("/apps", appController.Create)
	modelProviderV1.GET("/apps", appController.List)
	modelProviderV1.GET("/apps/:appID", appController.Detail)
	modelProviderV1.POST("/apps/:appID/model-config", appController.UpdateModelConfig)
	modelProviderV1.POST("/apps/:appID/chat-messages", appController.ChatMessage)
	return nil
}

func (r *AppRoutes) GetModule() string {
	return "app"
}
