/*
 * Copyright (c) 2024 Yunshan PeerConnections
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

type PeerConnection struct {
	cache *cache.Cache
}

func NewPeerConnection(c *cache.Cache) *PeerConnection {
	listener := &PeerConnection{
		cache: c,
	}
	return listener
}

func (c *PeerConnection) OnUpdaterAdded(addedDBItems []*metadbmodel.PeerConnection) {
	c.cache.AddPeerConnections(addedDBItems)
}

func (c *PeerConnection) OnUpdaterUpdated(cloudItem *cloudmodel.PeerConnection, diffBase *diffbase.PeerConnection) {
	diffBase.Update(cloudItem)
}

func (c *PeerConnection) OnUpdaterDeleted(lcuuids []string, deletedDBItems []*metadbmodel.PeerConnection) {
	c.cache.DeletePeerConnections(lcuuids)
}
