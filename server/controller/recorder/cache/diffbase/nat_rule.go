/**
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

package diffbase

import (
	ctrlrcommon "github.com/deepflowio/deepflow/server/controller/common"
	metadbmodel "github.com/deepflowio/deepflow/server/controller/db/metadb/model"
)

func (b *DataSet) AddNATRule(dbItem *metadbmodel.NATRule, seq int) {
	b.NATRules[dbItem.Lcuuid] = &NATRule{
		DiffBase: DiffBase{
			Sequence: seq,
			Lcuuid:   dbItem.Lcuuid,
		},
	}
	b.GetLogFunc()(addDiffBase(ctrlrcommon.RESOURCE_TYPE_NAT_RULE_EN, b.NATRules[dbItem.Lcuuid]), b.metadata.LogPrefixes)
}

func (b *DataSet) DeleteNATRule(lcuuid string) {
	delete(b.NATRules, lcuuid)
	log.Info(deleteDiffBase(ctrlrcommon.RESOURCE_TYPE_NAT_RULE_EN, lcuuid), b.metadata.LogPrefixes)
}

type NATRule struct {
	DiffBase
}
