# v6.5 Changelog

### DeepFlow release v6.5

#### New Feature
* feat: Add Changelog [#7095](https://github.com/deepflowio/deepflow/pull/7095)
* feat: CK’s username and password support the use of special characters [#7229](https://github.com/deepflowio/deepflow/pull/7119)

#### Bug Fix
* fix: Show metrics use query cache can be configured [#7271](https://github.com/deepflowio/deepflow/pull/7271) by [xiaochaoren1](https://github.com/xiaochaoren1)
* fix: server recorder prints unnecessary error logs [#7264](https://github.com/deepflowio/deepflow/pull/7264) by [ZhengYa-0110](https://github.com/ZhengYa-0110)
* fix: Network interface card list supports duplicate mac [#7253](https://github.com/deepflowio/deepflow/pull/7253) by [xiaochaoren1](https://github.com/xiaochaoren1)
* fix: uses long connection to connect to CK for datasources manager [#7242](https://github.com/deepflowio/deepflow/pull/7242) by [lzf575](https://github.com/lzf575)
* fix: server recorder prints unnecessary error logs [#7192](https://github.com/deepflowio/deepflow/pull/7192) by [ZhengYa-0110](https://github.com/ZhengYa-0110)
* fix: when modifying the TTL of the CK table fails, the connection to CK needs to be closed [#7202](https://github.com/deepflowio/deepflow/pull/7202) by [lzf575](https://github.com/lzf575)
* fix: getting wrong org-id and team-id when the Agent version is less than v6.5.9 [#7202](https://github.com/deepflowio/deepflow/pull/7202) by [lzf575](https://github.com/lzf575)
* fix: agent Packet fanout can only be set in TapMode::Local mode [#7229](https://github.com/deepflowio/deepflow/pull/7229) by [TomatoMr](https://github.com/TomatoMr)
* fix: Attribute supports Chinese [#7228](https://github.com/deepflowio/deepflow/pull/7228) by [xiaochaoren1](https://github.com/xiaochaoren1)
* fix: agent - incorrect grpc log collected by uprobe [#7201](https://github.com/deepflowio/deepflow/pull/7201) by [yuanchaoa](https://github.com/yuanchaoa)
* fix: Modify the value of team's short lcuuid corresponding to org_id [#7176](https://github.com/deepflowio/deepflow/pull/7176) by [jin-xiaofeng](https://github.com/jin-xiaofeng)
* fix: prometheus data cannot be labeled with universal tags, if slow-decoder is used. [#7100](https://github.com/deepflowio/deepflow/pull/7100)
* fix: agent - eBPF strengthening protocol inference for SOFARPC and MySQL [#7110](https://github.com/deepflowio/deepflow/pull/7110)
* fix: server controller changes prometheus label version when data does not change actually. [#7115](https://github.com/deepflowio/deepflow/pull/7115)
* fix: k8s refresh api http client close keep alive. [#7117](https://github.com/deepflowio/deepflow/pull/7117)
* fix: Ingester always update prometheus labels even if labels version has not changed.  [#7128](https://github.com/deepflowio/deepflow/pull/7128)

**[Changelog for v6.5](https://www.deepflow.io/docs/release-notes/release-6.5-ce)**<br/>

#### NEW FEATURE
* feat: agent - support vhost user [#7269](https://github.com/deepflowio/deepflow/pull/7269) by [yuanchaoa](https://github.com/yuanchaoa)
* feat: add debug ctl to rebalance agent by traffic [#7261](https://github.com/deepflowio/deepflow/pull/7261) by [roryye](https://github.com/roryye)
* feat: agent - eBPF Add JAVA symbol file generation log [#7257](https://github.com/deepflowio/deepflow/pull/7257) by [yinjiping](https://github.com/yinjiping)
* feat: Add volcengine icon const [#7232](https://github.com/deepflowio/deepflow/pull/7232) by [xiaochaoren1](https://github.com/xiaochaoren1)
* feat: add volcengine icon const [#7180](https://github.com/deepflowio/deepflow/pull/7180) by [askyrie](https://github.com/askyrie)

#### Documentation
* docs: rename opentemetry to opentelemetry [#7246](https://github.com/deepflowio/deepflow/pull/7246) by [lzf575](https://github.com/lzf575)


#### Chore
* chore: update cli dependencies [#7250](https://github.com/deepflowio/deepflow/pull/7250) by [lzf575](https://github.com/lzf575)


#### OTHER
* bump golang.org/x/net to v0.26.0 [#7234](https://github.com/deepflowio/deepflow/pull/7234) by [zhangzujian](https://github.com/zhangzujian)


#### Performance
* perf: add setting ttl_only_drop_parts to the CK table to make TTL more efficient [#7266](https://github.com/deepflowio/deepflow/pull/7266) by [lzf575](https://github.com/lzf575)


#### Testing
* chore: use the latest go version to build server/cli [#7235](https://github.com/deepflowio/deepflow/pull/7235) by [zhangzujian](https://github.com/zhangzujian)