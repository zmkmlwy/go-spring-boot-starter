/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package WebStarter

import (
	"fmt"

	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring/spring-boot"
	"github.com/go-spring/go-spring/spring-core"
)

func init() {
	SpringBoot.RegisterModule(func(ctx SpringCore.SpringContext) {
		ctx.RegisterBean(new(WebContainerConfig))
		ctx.RegisterBean(new(WebContainerStarter))
	})
}

//
// 定义 Web Bean 初始化接口
//
type WebBeanInitialization interface {
	InitWebBean(c SpringWeb.WebContainer, ctx SpringCore.SpringContext)
}

//
// Web 容器配置
//
type WebContainerConfig struct {
	EnableHTTP  bool   `value:"${web.server.enable:=true}"`      // 是否启用 HTTP
	Port        int32  `value:"${web.server.port:=8080}"`        // HTTP 端口
	EnableHTTPS bool   `value:"${web.server.ssl.enable:=false}"` // 是否启用 HTTPS
	SSLPort     int32  `value:"${web.server.ssl.port:=8443}"`    // SSL 端口
	SSLCert     string `value:"${web.server.ssl.cert:=}"`        // SSL 证书
	SSLKey      string `value:"${web.server.ssl.key:=}"`         // SSL 秘钥
}

//
// Web 容器启动器
//
type WebContainerStarter struct {
	Config       *WebContainerConfig    `autowire:""`
	Container    SpringWeb.WebContainer `autowire:"WebContainer?"`
	SSLContainer SpringWeb.WebContainer `autowire:"WebSSLContainer?"`
}

//
// 启动 Web 容器
//
func (starter *WebContainerStarter) runContainer(ctx SpringBoot.ApplicationContext,
	ssl bool, address string, certFile string, keyFile string) {

	var c SpringWeb.WebContainer

	if ssl {
		c = starter.SSLContainer
	} else {
		c = starter.Container
	}

	if c == nil {
		// 如果用户没有创建则使用默认的 Web 容器
		c = SpringWeb.WebContainerFactory()
	}

	var beans []WebBeanInitialization
	ctx.CollectBeans(&beans)

	// 初始化 Web Beans
	for _, bean := range beans {
		bean.InitWebBean(c, ctx)
	}

	// 启动 Web 容器
	ctx.SafeGoroutine(func() {

		var err error
		if ssl {
			starter.SSLContainer = c
			err = c.StartTLS(address, certFile, keyFile)
		} else {
			starter.Container = c
			err = c.Start(address)
		}

		fmt.Println("exit web server goroutine", err)
	})
}

func (starter *WebContainerStarter) OnStartApplication(ctx SpringBoot.ApplicationContext) {

	// 启动 HTTP 容器
	if starter.Config.EnableHTTP {
		address := fmt.Sprintf(":%d", starter.Config.Port)
		fmt.Printf("listen on %s%s\n", SpringUtils.LocalIPv4(), address)
		starter.runContainer(ctx, false, address, "", "")
	}

	// 启动 HTTPS 容器
	if starter.Config.EnableHTTPS {
		address := fmt.Sprintf(":%d", starter.Config.SSLPort)
		fmt.Printf("listen on %s%s\n", SpringUtils.LocalIPv4(), address)
		starter.runContainer(ctx, false, address, starter.Config.SSLCert, starter.Config.SSLKey)
	}
}

func (starter *WebContainerStarter) OnStopApplication(ctx SpringBoot.ApplicationContext) {

	// 停止 HTTP 容器
	if starter.Container != nil {
		starter.Container.Stop()
	}

	// 停止 HTTPS 容器
	if starter.SSLContainer != nil {
		starter.SSLContainer.Stop()
	}
}
