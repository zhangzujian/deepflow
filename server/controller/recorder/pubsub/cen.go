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

package pubsub

import (
	"github.com/deepflowio/deepflow/server/controller/common"
	"github.com/deepflowio/deepflow/server/controller/recorder/pubsub/message"
)

type CEN struct {
	ResourcePubSubComponent[
		*message.CENAdd,
		message.CENAdd,
		message.AddNoneAddition,
		*message.CENUpdate,
		message.CENUpdate,
		*message.CENFieldsUpdate,
		message.CENFieldsUpdate,
		*message.CENDelete,
		message.CENDelete,
		message.DeleteNoneAddition]
}

func NewCEN() *CEN {
	return &CEN{
		ResourcePubSubComponent[
			*message.CENAdd,
			message.CENAdd,
			message.AddNoneAddition,
			*message.CENUpdate,
			message.CENUpdate,
			*message.CENFieldsUpdate,
			message.CENFieldsUpdate,
			*message.CENDelete,
			message.CENDelete,
			message.DeleteNoneAddition,
		]{
			PubSubComponent: newPubSubComponent(PubSubTypeCEN),
			resourceType:    common.RESOURCE_TYPE_CEN_EN,
		},
	}
}
