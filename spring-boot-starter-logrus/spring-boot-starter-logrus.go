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

package SpringBootStarterLogrus

import (
	"os"
	"strings"
	"github.com/sirupsen/logrus"
	"github.com/didi/go-spring/spring-core"
	"github.com/didi/go-spring/spring-logrus"
	"github.com/didi/go-spring/spring-logger"
	"github.com/go-spring/go-spring-boot/spring-boot"
)

func init() {
	SpringLogger.SetLogger(logrus.StandardLogger())

	logrus.SetFormatter(new(SpringLogrus.TextFormatter))
	logrus.SetLevel(logrus.TraceLevel)

	output := SpringLogrus.NewNullOutput()
	logrus.SetOutput(output)

	SpringBoot.RegisterModule(func(context SpringCore.SpringContext) {
		properties := context.GetPrefixProperties("logger.appender.")

		appenderMap := make(map[string]SpringLogger.LoggerAppender)
		for key := range properties {
			ss := strings.Split(key, ".")
			appenderMap[ss[2]] = nil
		}

		for key := range appenderMap {
			if t, ok := properties["logger.appender."+key+".type"]; ok {
				switch t {
				case "ConsoleAppender":

					level := logrus.DebugLevel
					if l, ok := properties["logger.appender."+key+".level"]; ok {
						level, _ = logrus.ParseLevel(l)
					}

					appender := SpringLogger.NewConsoleAppender()
					appenderMap[key] = appender
					logrus.AddHook(SpringLogrus.NewSpringLogrusHook(appender, level))

				case "FileAppender":

					level := logrus.DebugLevel
					if l, ok := properties["logger.appender."+key+".level"]; ok {
						level, _ = logrus.ParseLevel(l)
					}

					workDir, _ := os.Getwd()
					app := context.GetProperties("spring.application.name")

					filePath := workDir + "/log/" + app + ".log"
					if pattern, ok := properties["logger.appender."+key+".pattern"]; ok {
						filePath = workDir + "/log/" + pattern
					}

					logFile, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					appender := SpringLogger.NewFileAppender(logFile)
					appenderMap[key] = appender
					logrus.AddHook(SpringLogrus.NewSpringLogrusHook(appender, level))
				}
			}
		}

		output.Output(appenderMap)
	})
}
