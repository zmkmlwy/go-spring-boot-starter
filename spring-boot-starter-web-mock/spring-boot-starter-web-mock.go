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

package SpringBootStarterWebMock

import (
	"os"
	"strings"
	"io/ioutil"
	"path/filepath"
	"github.com/didi/go-spring/spring-web"
	"github.com/didi/go-spring/spring-core"
	"github.com/go-spring/go-spring-boot/spring-boot"
)

func init() {
	SpringBoot.RegisterModule(func(ctx SpringCore.SpringContext) {
		location, _ := ctx.GetDefaultProperties("mock.config.location", "mock/")

		mc := new(SpringWeb.MockController)
		mc.Mapping = loadMockData(location)
		ctx.RegisterBean(mc)
	})
}

func loadMockData(mockDirPath string) map[string]*SpringWeb.MockData {
	mapping := make(map[string]*SpringWeb.MockData, 0)
	filepath.Walk(mockDirPath, func(path string, f os.FileInfo, err error) error {

		if f == nil || f.IsDir() {
			return nil
		}

		mockPath := strings.TrimPrefix(path, mockDirPath)

		uri := mockPath
		method := "get"

		if index := strings.LastIndex(mockPath, "."); index >= 0 {
			method = mockPath[index+1:]
			uri = mockPath[0:index]
		}

		uri = "/" + uri

		bytes, _ := ioutil.ReadFile(path)

		key := method + ":" + uri
		mapping[key] = &SpringWeb.MockData{
			Method: method,
			Path:   uri,
			Data:   string(bytes),
		}

		return nil
	})
	return mapping
}
