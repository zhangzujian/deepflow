/*
 * Copyright (c) 2024 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("db.metadb.migrator.common")

func CreateDatabase(dc *DBConfig) error {
	log.Infof(LogDBName(dc.Config.Database, "create database"))
	return dc.DB.Exec(dc.SqlFmt.CreateDatabase()).Error
}

func CreateDatabaseIfNotExists(dc *DBConfig) (bool, error) {
	var databaseName string
	dc.DB.Raw(dc.SqlFmt.SelectDatabase()).Scan(&databaseName)
	if databaseName == dc.Config.Database {
		return true, nil
	} else {
		err := CreateDatabase(dc)
		return false, err
	}
}
