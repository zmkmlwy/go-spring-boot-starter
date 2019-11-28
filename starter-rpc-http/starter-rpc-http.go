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

package HttpRpcStarter

import (
	"fmt"

	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/go-spring/go-spring-rpc/spring-rpc"
	"github.com/go-spring/go-spring-rpc/spring-rpc-http"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring/spring-boot"
)

func init() {
	SpringBoot.RegisterBean(new(RpcContainerConfig))
	SpringBoot.RegisterBean(new(RpcContainerStarter))
}

//
// RPC 容器配置
//
type RpcContainerConfig struct {
	EnableHTTP  bool   `value:"${rpc.server.enable:=true}"`      // 是否启用 HTTP
	Port        int32  `value:"${rpc.server.port:=9090}"`        // HTTP 端口
	EnableHTTPS bool   `value:"${rpc.server.ssl.enable:=false}"` // 是否启用 HTTPS
	SSLPort     int32  `value:"${rpc.server.ssl.port:=9443}"`    // SSL 端口
	SSLCert     string `value:"${rpc.server.ssl.cert:=}"`        // SSL 证书
	SSLKey      string `value:"${rpc.server.ssl.key:=}"`         // SSL 秘钥
}

//
// RPC 容器启动器
//
type RpcContainerStarter struct {
	Config       *RpcContainerConfig `autowire:""`
	Container    SpringRpc.RpcContainer
	SSLContainer SpringRpc.RpcContainer
}

//
// 启动 RPC 容器
//
func (starter *RpcContainerStarter) runContainer(ctx SpringBoot.ApplicationContext,
	ssl bool, address string, certFile string, keyFile string) {

	// 创建 RPC 容器对象
	c := SpringHttpRpc.NewContainer(SpringWeb.WebContainerFactory())

	var beans []SpringRpc.RpcBeanInitialization
	ctx.CollectBeans(&beans)

	// 初始化 RPC Beans
	for _, bean := range beans {
		bean.InitRpcBean(c)
	}

	// 启动 RPC 容器
	ctx.SafeGoroutine(func() {

		var err error
		if ssl {
			starter.SSLContainer = c
			err = c.StartTLS(address, certFile, keyFile)
		} else {
			starter.Container = c
			err = c.Start(address)
		}

		fmt.Println("exit rpc server goroutine", err)
	})
}

func (starter *RpcContainerStarter) OnStartApplication(ctx SpringBoot.ApplicationContext) {

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

func (starter *RpcContainerStarter) OnStopApplication(ctx SpringBoot.ApplicationContext) {

	// 停止 HTTP 容器
	if starter.Container != nil {
		starter.Container.Stop()
	}

	// 停止 HTTPS 容器
	if starter.SSLContainer != nil {
		starter.SSLContainer.Stop()
	}
}
