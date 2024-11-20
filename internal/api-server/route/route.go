// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	authRoute "github.com/lunarianss/Luna/internal/api-server/route/console/auth"
	setupRoute "github.com/lunarianss/Luna/internal/api-server/route/console/setup"
	consoleWorkSpaceRoute "github.com/lunarianss/Luna/internal/api-server/route/console/workspace"
	"github.com/lunarianss/Luna/internal/pkg/server"
)

// Route unified registration portal
func init() {
	server.RegisterRoute(&blogRoutes{})

	// console/workspace
	server.RegisterRoute(&consoleWorkSpaceRoute.ModelProviderRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.ModelRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.AppRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.AccountRoute{})
	server.RegisterRoute(&consoleWorkSpaceRoute.WorkspaceRoutes{})

	// console/auth
	server.RegisterRoute(&authRoute.AuthRoutes{})

	// console/setup
	server.RegisterRoute(&setupRoute.SetupRoutes{})

	server.RegisterRoute(&staticRoute{})

}
