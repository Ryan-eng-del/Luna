// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/lunarianss/Luna/internal/api-server/service"
)

type AppController struct {
	appService  *service.AppService
	chatService *service.ChatService
}

func NewAppController(appSrv *service.AppService, chatSrv *service.ChatService) *AppController {
	return &AppController{appService: appSrv, chatService: chatSrv}
}
