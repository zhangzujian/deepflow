START TRANSACTION;

-- modify start, add upgrade sql

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"\",\"PROM_SQL\":\"delta(min(deepflow_system__deepflow_agent_monitor__create_time)by(host)[1m:10s])\",\"interval\":60,\"metric\":\"process_start_time_delta\",\"time_tag\":\"toi\"}", target_field="{\"displayName\":\"process_start_time_delta\", \"unit\": \"毫秒\"}", name="采集器重启" WHERE name="进程启动";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_dispatcher\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.kernel_drops`) AS `dispatcher.metrics.kernel_drops`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.kernel_drops`) AS `dispatcher.metrics.kernel_drops`\"]}]}", target_field="{\"displayName\":\"dispatcher.metrics.kernel_drops\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (dispatcher.metrics.kernel_drops)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_queue\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.overwritten`) AS `queue.metrics.overwritten`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.overwritten`) AS `queue.metrics.overwritten`\"]}]}", target_field="{\"displayName\":\"queue.metrics.overwritten\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (queue.metrics.overwritten)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_l7_session_aggr\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.throttle-drop`) AS `l7_session_aggr.metrics.throttle-drop`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.throttle-drop`) AS `l7_session_aggr.metrics.throttle-drop`\"]}]}", target_field="{\"displayName\":\"l7_session_aggr.metrics.throttle-drop\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (l7_session_aggr.metrics.throttle-drop)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_flow_aggr\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop-in-throttle`) AS `flow_aggr.metrics.drop-in-throttle`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop-in-throttle`) AS `flow_aggr.metrics.drop-in-throttle`\"]}]}", target_field="{\"displayName\":\"flow_aggr.metrics.drop-in-throttle\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (flow_aggr.metrics.drop-in-throttle)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_ebpf_collector\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.kern_lost`) AS `ebpf_collector.metrics.kern_lost`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.kern_lost`) AS `ebpf_collector.metrics.kern_lost`\"]}]}", target_field="{\"displayName\":\"ebpf_collector.metrics.kern_lost\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (ebpf_collector.metrics.kern_lost)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_ebpf_collector\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.user_enqueue_lost`) AS `ebpf_collector.metrics.user_enqueue_lost`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.kernuser_enqueue_lost_lost`) AS `ebpf_collector.metrics.user_enqueue_lost`\"]}]}", target_field="{\"displayName\":\"ebpf_collector.metrics.user_enqueue_lost\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (ebpf_collector.metrics.user_enqueue_lost)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_dispatcher\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.invalid_packets`) AS `dispatcher.metrics.invalid_packets`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.invalid_packets`) AS `dispatcher.metrics.invalid_packets`\"]}]}", target_field="{\"displayName\":\"dispatcher.metrics.invalid_packets\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (dispatcher.metrics.invalid_packets)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_dispatcher\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.err`) AS `dispatcher.metrics.err`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.err`) AS `dispatcher.metrics.err`\"]}]}", target_field="{\"displayName\":\"dispatcher.metrics.err\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (dispatcher.metrics.err)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_flow_map\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop_by_window`) AS `flow_map.metrics.drop_by_window`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop_by_window`) AS `flow_map.metrics.drop_by_window`\"]}]}", target_field="{\"displayName\":\"flow_map.metrics.drop_before_window\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (flow_map.metrics.drop_before_window)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_flow_aggr\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop-before-window`) AS `flow_aggr.metrics.drop-before-window`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop-before-window`) AS `flow_aggr.metrics.drop-before-window`\"]}]}", target_field="{\"displayName\":\"flow_aggr.metrics.drop-before-window\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (flow_aggr.metrics.drop-before-window)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_quadruple_generator\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop-before-window`) AS `quadruple_generator.metrics.drop-before-window`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop-before-window`) AS `quadruple_generator.metrics.drop-before-window`\"]}]}", target_field="{\"displayName\":\"quadruple_generator.metrics.drop-before-window\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (quadruple_generator.metrics.drop-before-window)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_collector\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop-before-window`) AS `collector.metrics.drop-before-window`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop-before-window`) AS `collector.metrics.drop-before-window`\"]}]}", target_field="{\"displayName\":\"collector.metrics.drop-before-window\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (collector.metrics.drop-before-window)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_collector\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop-inactive`) AS `collector.metrics.drop-inactive`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop-inactive`) AS `collector.metrics.drop-inactive`\"]}]}", target_field="{\"displayName\":\"collector.metrics.drop-inactive\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (collector.metrics.drop-inactive)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_agent_collect_sender\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.dropped`) AS `collect_sender.metrics.dropped`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.dropped`) AS `collect_sender.metrics.dropped`\"]}]}", target_field="{\"displayName\":\"collect_sender.metrics.dropped\", \"unit\": \"\"}" WHERE name="采集器数据丢失 (collect_sender.metrics.dropped)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_server_ingester_recviver\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.invalid`) AS `ingester.recviver.metrics.invalid`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.invalid`) AS `ingester.recviver.metrics.invalid`\"]}]}", target_field="{\"displayName\":\"ingester.recviver.metrics.invalid\", \"unit\": \"\"}" WHERE name="数据节点数据丢失 (ingester.recviver.metrics.invalid)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_server_ingester_queue\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.overwritten`) AS `ingester.queue.metrics.overwritten`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.overwritten`) AS `ingester.queue.metrics.overwritten`\"]}]}", target_field="{\"displayName\":\"ingester.queue.metrics.overwritten\", \"unit\": \"\"}" WHERE name="数据节点数据丢失 (ingester.queue.metrics.overwritten)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_server_ingester_decoder\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.drop_count`) AS `ingester.decoder.metrics.drop_count`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.drop_count`) AS `ingester.decoder.metrics.drop_count`\"]}]}", target_field="{\"displayName\":\"ingester.decoder.metrics.drop_count\", \"unit\": \"\"}" WHERE name="数据节点数据丢失 (ingester.decoder.metrics.drop_count)";

UPDATE alarm_policy SET query_params="{\"DATABASE\":\"deepflow_system\",\"TABLE\":\"deepflow_server_ingester_ckwriter\",\"interval\":60,\"fill\": \"none\",\"window_size\":1,\"QUERIES\":[{\"QUERY_ID\":\"R1\",\"SELECT\":\"Sum(`metrics.write_failed_count`) AS `ingester.ckwriter.metrics.write_failed_count`\",\"WHERE\":\"1=1\",\"GROUP_BY\":\"`tag.host`\",\"METRICS\":[\"Sum(`metrics.write_failed_count`) AS `ingester.ckwriter.metrics.write_failed_count`\"]}]}", target_field="{\"displayName\":\"ingester.ckwriter.metrics.write_failed_count\", \"unit\": \"\"}" WHERE name="数据节点数据丢失 (ingester.ckwriter.metrics.write_failed_count)";


-- update db_version to latest, remeber update DB_VERSION_EXPECT in migrate/init.go
UPDATE db_version SET version='6.4.1.30';
-- modify end

COMMIT;