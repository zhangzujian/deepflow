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

package config

import (
	configs "github.com/deepflowio/deepflow/server/controller/config/common"
)

type MonitorConfig struct {
	HealthCheckInterval         int                           `default:"60" yaml:"health_check_interval"`
	HealthCheckPort             int                           `default:"30417" yaml:"health_check_port"`
	HealthCheckHandleChannelLen int                           `default:"1000" yaml:"health_check_handle_channel_len"`
	LicenseCheckInterval        int                           `default:"60" yaml:"license_check_interval"`
	VTapCheckInterval           int                           `default:"60" yaml:"vtap_check_interval"`
	ExceptionTimeFrame          int                           `default:"3600" yaml:"exception_time_frame"`
	AutoRebalanceVTap           bool                          `default:"true" yaml:"auto_rebalance_vtap"`
	RebalanceCheckInterval      int                           `default:"300" yaml:"rebalance_check_interval"` // unit: second
	VTapAutoDelete              VTapAutoDelete                `yaml:"vtap_auto_delete"`
	Warrant                     configs.Warrant               `yaml:"warrant"`
	IngesterLoadBalancingConfig IngesterLoadBalancingStrategy `yaml:"ingester-load-balancing-strategy"`
	SyncDefaultORGDataInterval  int                           `default:"10" yaml:"sync_default_org_data_interval"`
}

type IngesterLoadBalancingStrategy struct {
	Algorithm         string `default:"by-ingested-data" yaml:"algorithm"` // options: by-ingested-data, by-agent-count
	DataDuration      int    `default:"86400" yaml:"data-duration"`        // default: 1d
	RebalanceInterval int    `default:"3600" yaml:"rebalance-interval"`    // default: 1h
}

type VTapAutoDelete struct {
	Enabled     bool `default:"true" yaml:"enabled"`
	LostTimeMax int  `default:"3600" yaml:"lost_time_max"` // unit: second
}
