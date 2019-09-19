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

package SpringBootWeb

import (
	"fmt"
	"github.com/didi/go-spring/spring-web"
	"github.com/didi/go-spring/spring-utils"
	Logger "github.com/didi/go-spring/spring-logger"
	"github.com/go-spring/go-spring-boot/spring-boot"
)

//
// 容器
//
type WebContainerWrapper struct {
	container       SpringWeb.WebContainer
	ServerPort      int32  `value:"${server.port:=8080}"`
	ServerSSLEnable bool   `value:"${server.ssl.enable:=false}"`
	ServerSSLCert   string `value:"${server.ssl.cert:=}"`
	ServerSSLKey    string `value:"${server.ssl.key:=}"`
}

func Wrapper(container SpringWeb.WebContainer) *WebContainerWrapper {
	return &WebContainerWrapper{
		container: container,
	}
}

func (wrapper *WebContainerWrapper) OnStartApplication(ctx SpringBoot.SpringApplicationContext) {

	var controllerBeans []SpringWeb.SpringControllerInitialization
	ctx.FindBeansByType(&controllerBeans)

	for _, bean := range controllerBeans {
		bean.InitController(wrapper.container)
	}

	ctx.SafeGoroutine(func() {
		Logger.Infoln("run server goroutine")

		address := fmt.Sprintf(":%d", wrapper.ServerPort)
		Logger.Debugf("listening on %s%s\n", SpringUtils.LocalIPv4(), address)

		var err error
		if wrapper.ServerSSLEnable {
			err = wrapper.container.StartTLS(address, wrapper.ServerSSLCert, wrapper.ServerSSLKey)
		} else {
			err = wrapper.container.Start(address)
		}
		Logger.Infoln("exit server goroutine", err)
	})
}

func (wrapper *WebContainerWrapper) OnStopApplication(ctx SpringBoot.SpringApplicationContext) {
	wrapper.container.Stop()
}
