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

package listener

import (
	cloudmodel "github.com/deepflowio/deepflow/server/controller/cloud/model"
	metadbmodel "github.com/deepflowio/deepflow/server/controller/db/metadb/model"
	"github.com/deepflowio/deepflow/server/controller/recorder/cache"
	"github.com/deepflowio/deepflow/server/controller/recorder/cache/diffbase"
)

type LBTargetServer struct {
	cache *cache.Cache
}

func NewLBTargetServer(c *cache.Cache) *LBTargetServer {
	listener := &LBTargetServer{
		cache: c,
	}
	return listener
}

func (l *LBTargetServer) OnUpdaterAdded(addedDBItems []*metadbmodel.LBTargetServer) {
	l.cache.AddLBTargetServers(addedDBItems)
}

func (l *LBTargetServer) OnUpdaterUpdated(cloudItem *cloudmodel.LBTargetServer, diffBase *diffbase.LBTargetServer) {
	diffBase.Update(cloudItem)
}

func (l *LBTargetServer) OnUpdaterDeleted(lcuuids []string, deletedDBItems []*metadbmodel.LBTargetServer) {
	l.cache.DeleteLBTargetServers(lcuuids)
}
