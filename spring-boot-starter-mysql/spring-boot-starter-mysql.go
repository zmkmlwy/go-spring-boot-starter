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

package SpringBootStarterMysql

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/didi/go-spring/spring-core"
	"github.com/didi/go-spring/spring-mysql"
	"github.com/go-spring/go-spring-boot/spring-boot"
)

func init() {
	SpringBoot.RegisterModule(func(ctx SpringCore.SpringContext) {

		host := ctx.GetProperties("spring.datasource.host")
		port, _ := ctx.GetDefaultProperties("spring.datasource.port", "3306")
		username := ctx.GetProperties("spring.datasource.username")
		password := ctx.GetProperties("spring.datasource.password")
		dbName := ctx.GetProperties("spring.datasource.db")

		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)

		db, err := sql.Open("mysql", url)
		if err != nil {
			panic(err)
		}

		ctx.RegisterBean(&SpringMysql.MysqlTemplateImpl{
			DB: db,
		})
	})
}
