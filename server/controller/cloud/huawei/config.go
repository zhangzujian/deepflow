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

package huawei

import (
	"strings"

	"github.com/bitly/go-simplejson"

	cloudcommon "github.com/deepflowio/deepflow/server/controller/cloud/common"
	"github.com/deepflowio/deepflow/server/controller/common"
	"github.com/deepflowio/deepflow/server/libs/logger"
)

var DEFAULT_DOMAIN = "myhuaweicloud.com"
var DEFAULT_IAM_HOST_PREFIX = "iam"
var DEFAULT_IAM_HOST = strings.Join([]string{DEFAULT_IAM_HOST_PREFIX, DEFAULT_DOMAIN}, ".")

type Config struct {
	RegionLcuuid   string
	AccountName    string
	IAMName        string
	Password       string
	IAMHost        string
	IAMHostPrefix  string
	ProjectID      string
	ProjectName    string
	Domain         string // 用于构造访问华为云的endpoint，需与DeepFlow自身domain做区分
	IncludeRegions map[string]bool
}

func (c *Config) LoadFromString(orgID int, sConf string) (err error) {
	jConf, err := simplejson.NewJson([]byte(sConf))
	if err != nil {
		log.Error("convert config string: %s to json failed: %v", sConf, err, logger.NewORGPrefix(orgID))
		return
	}
	c.AccountName, err = jConf.Get("account_name").String()
	if err != nil {
		log.Error("account_name must be specified", logger.NewORGPrefix(orgID))
		return
	}
	c.IAMName, err = jConf.Get("iam_name").String()
	if err != nil {
		log.Error("iam_name must be specified", logger.NewORGPrefix(orgID))
		return
	}
	pswd, err := jConf.Get("password").String()
	if err != nil {
		log.Error("password must be specified", logger.NewORGPrefix(orgID))
		return
	}
	dpswd, err := common.DecryptSecretKey(pswd)
	if err != nil {
		log.Errorf("decrypt password failed (%s)", logger.NewORGPrefix(orgID))
		return
	}
	c.Password = dpswd

	c.IAMHost = jConf.Get("iam_host").MustString()
	if c.IAMHost == "" {
		c.IAMHost = DEFAULT_IAM_HOST
	}
	c.IAMHostPrefix = strings.Split(c.IAMHost, ".")[0]
	c.ProjectID, err = jConf.Get("project_id").String()
	if err != nil {
		log.Error("project_id must be specified", logger.NewORGPrefix(orgID))
		return
	}
	c.ProjectName, err = jConf.Get("region_name").String()
	if err != nil {
		log.Error("region_name must be specified", logger.NewORGPrefix(orgID))
		return
	}
	c.Domain = jConf.Get("domain").MustString()
	if c.Domain == "" {
		c.Domain = DEFAULT_DOMAIN
	}
	c.RegionLcuuid, err = jConf.Get("region_uuid").String()
	if err != nil {
		log.Error("region_uuid must be specified", logger.NewORGPrefix(orgID))
		return
	}
	c.IncludeRegions = cloudcommon.UniqRegions(jConf.Get("include_regions").MustString())
	return
}
