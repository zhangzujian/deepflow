/*
 * Copyright (c) 2024 Yunshan PodIngressRuleBackends
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

package listener

import (
	cloudmodel "github.com/deepflowio/deepflow/server/controller/cloud/model"
	metadbmodel "github.com/deepflowio/deepflow/server/controller/db/metadb/model"
	"github.com/deepflowio/deepflow/server/controller/recorder/cache"
	"github.com/deepflowio/deepflow/server/controller/recorder/cache/diffbase"
)

type PodIngressRuleBackend struct {
	cache *cache.Cache
}

func NewPodIngressRuleBackend(c *cache.Cache) *PodIngressRuleBackend {
	listener := &PodIngressRuleBackend{
		cache: c,
	}
	return listener
}

func (b *PodIngressRuleBackend) OnUpdaterAdded(addedDBItems []*metadbmodel.PodIngressRuleBackend) {
	b.cache.AddPodIngressRuleBackends(addedDBItems)
}

func (b *PodIngressRuleBackend) OnUpdaterUpdated(cloudItem *cloudmodel.PodIngressRuleBackend, diffBase *diffbase.PodIngressRuleBackend) {
}

func (b *PodIngressRuleBackend) OnUpdaterDeleted(lcuuids []string, deletedDBItems []*metadbmodel.PodIngressRuleBackend) {
	b.cache.DeletePodIngressRuleBackends(lcuuids)
}
