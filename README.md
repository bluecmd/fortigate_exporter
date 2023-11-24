# fortigate_exporter

![Go](https://github.com/bluecmd/fortigate_exporter/workflows/Go/badge.svg)
![Docker](https://github.com/bluecmd/fortigate_exporter/workflows/Docker/badge.svg)
[![Docker Repository on Quay](https://quay.io/repository/bluecmd/fortigate_exporter/status "Docker Repository on Quay")](https://quay.io/repository/bluecmd/fortigate_exporter)
[![Matrix](https://img.shields.io/matrix/fortigate_exporter:matrix.org)](https://matrix.to/#/#fortigate_exporter:matrix.org)

Prometheus exporter for FortiGate® firewalls.

**NOTE:** This is not an official Fortinet product, it is developed fully independently by professionals and hobbyists alike.

  * [Supported Metrics](#supported-metrics)
  * [Usage](#usage)
    + [Dynamic configuration](#dynamic-configuration)
    + [Available CLI parameters](#available-cli-parameters)
    + [Fortigate Configuration](#fortigate-configuration)
    + [Prometheus Configuration](#prometheus-configuration)
    + [Docker](#docker)
      - [docker-compose](#docker-compose)
  * [Known Issues](#known-issues)
  * [Missing Metrics?](#missing-metrics)

## Supported Metrics

Right now the exporter supports a quite limited set of metrics, but it is very easy to add!
Open an issue if your favorite metric is missing.

For example PromQL usage, see [EXAMPLES](EXAMPLES.md).

Supported metrics right now as follows.

Global:

 * _System/SensorInfo_
   * `fortigate_sensor_fan_rpm`
   * `fortigate_sensor_temperature_celsius`
   * `fortigate_sensor_voltage_volts`
 * _System/Status_
   * `fortigate_version_info`
 * _System/Time/Clock_
   * `fortigate_time_seconds`
 * _System/Resource/Usage_
   * `fortigate_cpu_usage_ratio`
   * `fortigate_memory_usage_ratio`
   * `fortigate_current_sessions`
 * _System/HAChecksums_
   * `fortigate_ha_member_has_role`
 * _License/Status_
   * `fortigate_license_vdom_usage`
   * `fortigate_license_vdom_max`
 * _WebUI/State_
   * `fortigate_last_reboot_seconds`
   * `fortigate_last_snapshot_seconds`

Per-VDOM:

 * _System/VDOMResources_
   * `fortigate_vdom_cpu_usage_ratio`
   * `fortigate_vdom_memory_usage_ratio`
   * `fortigate_vdom_current_sessions`
 * _Firewall/Policies_
   * `fortigate_policy_active_sessions`
   * `fortigate_policy_bytes_total`
   * `fortigate_policy_hit_count_total`
   * `fortigate_policy_packets_total`
 * _Firewall/IpPool_
   * `fortigate_ippool_available_ratio`
   * `fortigate_ippool_used_ips`
   * `fortigate_ippool_total_ips`
   * `fortigate_ippool_clients`
   * `fortigate_ippool_used_items`
   * `fortigate_ippool_total_items`
 * _System/Fortimanager/Status_
   * `fortigate_fortimanager_connection_status`
   * `fortigate_fortimanager_registration_status`
 * _System/Interface_
   * `fortigate_interface_link_up`
   * `fortigate_interface_speed_bps`
   * `fortigate_interface_transmit_packets_total`
   * `fortigate_interface_receive_packets_total`
   * `fortigate_interface_transmit_bytes_total`
   * `fortigate_interface_receive_bytes_total`
   * `fortigate_interface_transmit_errors_total`
   * `fortigate_interface_receive_errors_total`
 * _System/SDNConnector_
   * `fortigate_system_sdn_connector_status`
   * `fortigate_system_sdn_connector_last_update_seconds`
 * _User/Fsso_
   * `fortigate_user_fsso_info`
 * _VPN/Ssl/Connections_
   * `fortigate_vpn_connections`
   * `fortigate_vpn_users`
 * _VPN/Ssl/Stats_
   * `fortigate_vpn_ssl_users`
   * `fortigate_vpn_ssl_tunnels`
   * `fortigate_vpn_ssl_connections`
 * _VPN/IPSec_
   * `fortigate_ipsec_tunnel_receive_bytes_total`
   * `fortigate_ipsec_tunnel_transmit_bytes_total`
   * `fortigate_ipsec_tunnel_up`
 * _Wifi/APStatus_
   * `fortigate_wifi_access_points`
   * `fortigate_wifi_fabric_clients`
   * `fortigate_wifi_fabric_max_allowed_clients`
 * _Log/Fortianalyzer/Status_
   * `fortigate_log_fortianalyzer_registration_info`
   * `fortigate_log_fortianalyzer_logs_received`
 * _Log/Fortianalyzer/Queue_
   * `fortigate_log_fortianalyzer_queue_connections`
   * `fortigate_log_fortianalyzer_queue_logs`
 * _Log/DiskUsage_
   * `fortigate_log_disk_used_bytes`
   * `fortigate_log_disk_total_bytes`

 Per-HA-Member and VDOM:
 * _System/HAStatistics_
   * `fortigate_ha_member_info`
   * `fortigate_ha_member_cpu_usage_ratio`
   * `fortigate_ha_member_memory_usage_ratio`
   * `fortigate_ha_member_network_usage_ratio`
   * `fortigate_ha_member_sessions`
   * `fortigate_ha_member_packets_total`
   * `fortigate_ha_member_virus_events_total`
   * `fortigate_ha_member_bytes_total`
   * `fortigate_ha_member_ips_events_total`

 Per-Link and VDOM:
 * _System/LinkMonitor_
   * `fortigate_link_status`
   * `fortigate_link_latency_seconds`
   * `fortigate_link_latency_jitter_seconds`
   * `fortigate_link_packet_loss_ratio`
   * `fortigate_link_packet_sent_total`
   * `fortigate_link_packet_received_total`
   * `fortigate_link_active_sessions`
   * `fortigate_link_bandwidth_tx_byte_per_second`
   * `fortigate_link_bandwidth_rx_byte_per_second`
   * `fortigate_link_status_change_time_seconds`

 Per-SDWAN and VDOM:
 * _VirtualWAN/HealthCheck_
   * `fortigate_virtual_wan_status`
   * `fortigate_virtual_wan_latency_seconds`
   * `fortigate_virtual_wan_latency_jitter_seconds`
   * `fortigate_virtual_wan_packet_loss_ratio`
   * `fortigate_virtual_wan_packet_sent_total`
   * `fortigate_virtual_wan_packet_received_total`
   * `fortigate_virtual_wan_active_sessions`
   * `fortigate_virtual_wan_bandwidth_tx_byte_per_second`
   * `fortigate_virtual_wan_bandwidth_rx_byte_per_second`
   * `fortigate_virtual_wan_status_change_time_seconds`

 Per-BGP-Neighbor and VDOM:
 * _BGP/Neighbors/IPv4_
   * `fortigate_bgp_neighbor_ipv4_info`
 * _BGP/Neighbors/IPv6_
   * `fortigate_bgp_neighbor_ipv6_info`
 * _BGP/NeighborPaths/IPv4_
   * `fortigate_bgp_neighbor_ipv4_paths`
   * `fortigate_bgp_neighbor_ipv4_best_paths`
 * _BGP/NeighborPaths/IPv6_
   * `fortigate_bgp_neighbor_ipv6_paths`
   * `fortigate_bgp_neighbor_ipv6_best_paths`

 Per-OSPF-Neighbor and VDOM:
 * _OSPF/Neighbors_
   * `fortigate_ospf_neighbor_info`

 Per-VirtualServer and VDOM:
 * _Firewall/LoadBalance_
   * `fortigate_lb_virtual_server_info`

 Per-RealServer for each VirtualServer and VDOM:
 * _Firewall/LoadBalance_
   * `fortigate_lb_real_server_info`
   * `fortigate_lb_real_server_mode`
   * `fortigate_lb_real_server_status`
   * `fortigate_lb_real_server_active_sessions`
   * `fortigate_lb_real_server_rtt_seconds`
   * `fortigate_lb_real_server_processed_bytes_total`

 Per-Certificate:
 * _System/AvailableCertificates_
   * `fortigate_certificate_info`
   * `fortigate_certificate_valid_from_seconds`
   * `fortigate_certificate_valid_to_seconds`
   * `fortigate_certificate_cmdb_references`

Per-VDOM and Wifi-Client:
 * _Wifi/Clients_
   * `fortigate_wifi_client_info`
   * `fortigate_wifi_client_data_rate_bps`
   * `fortigate_wifi_client_bandwidth_rx_bps`
   * `fortigate_wifi_client_bandwidth_tx_bps`
   * `fortigate_wifi_client_signal_strength_dBm`
   * `fortigate_wifi_client_signal_noise_dBm`
   * `fortigate_wifi_client_tx_discard_ratio`
   * `fortigate_wifi_client_tx_retries_ratio`

Per-VDOM and managed access point:
 * _Wifi/ManagedAP_
   * `fortigate_wifi_managed_ap_info`
   * `fortigate_wifi_managed_ap_join_time_seconds`
   * `fortigate_wifi_managed_ap_cpu_usage_ratio`
   * `fortigate_wifi_managed_ap_memory_free_bytes`
   * `fortigate_wifi_managed_ap_memory_bytes_total`

Per-VDOM, managed access point and radio:
 * _Wifi/ManagedAP_
   * `fortigate_wifi_managed_ap_radio_info`
   * `fortigate_wifi_managed_ap_radio_client_count`
   * `fortigate_wifi_managed_ap_radio_operating_tx_power_ratio`
   * `fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio`
   * `fortigate_wifi_managed_ap_radio_bandwidth_rx_bps`
   * `fortigate_wifi_managed_ap_radio_rx_bytes_total`
   * `fortigate_wifi_managed_ap_radio_tx_bytes_total`
   * `fortigate_wifi_managed_ap_radio_interfering_aps`
   * `fortigate_wifi_managed_ap_radio_tx_power_ratio`
   * `fortigate_wifi_managed_ap_radio_tx_discard_ratio`
   * `fortigate_wifi_managed_ap_radio_tx_retries_ratio`

Per-VDOM, managed access point and interface:
 * _Wifi/ManagedAP_
   * `fortigate_wifi_managed_ap_interface_rx_bytes_total`
   * `fortigate_wifi_managed_ap_interface_tx_bytes_total`
   * `fortigate_wifi_managed_ap_interface_rx_packets_total`
   * `fortigate_wifi_managed_ap_interface_tx_packets_total`
   * `fortigate_wifi_managed_ap_interface_rx_errors_total`
   * `fortigate_wifi_managed_ap_interface_tx_errors_total`
   * `fortigate_wifi_managed_ap_interface_rx_dropped_packets_total`
   * `fortigate_wifi_managed_ap_interface_tx_dropped_packets_total`

Per-VDOM, managed switch and interface:
* _Switch/ManagedSwitch_
  * `fortigate_managed_switch_collisions_total`
  * `fortigate_managed_switch_crc_alignments_total`
  * `fortigate_managed_switch_fragments_total`
  * `fortigate_managed_switch_info`
  * `fortigate_managed_switch_jabbers_total`
  * `fortigate_managed_switch_l3_packets_total`
  * `fortigate_managed_switch_max_poe_budget_watt`
  * `fortigate_managed_switch_port_info`
  * `fortigate_managed_switch_port_power_status`
  * `fortigate_managed_switch_port_power_watt`
  * `fortigate_managed_switch_port_status`
  * `fortigate_managed_switch_rx_bcast_packets_total`
  * `fortigate_managed_switch_rx_bytes_total`
  * `fortigate_managed_switch_rx_drops_total`
  * `fortigate_managed_switch_rx_errors_total`
  * `fortigate_managed_switch_rx_mcast_packets_total`
  * `fortigate_managed_switch_rx_oversize_total`
  * `fortigate_managed_switch_rx_packets_total`
  * `fortigate_managed_switch_rx_ucast_packets_total`
  * `fortigate_managed_switch_tx_bcast_packets_total`
  * `fortigate_managed_switch_tx_bytes_total`
  * `fortigate_managed_switch_tx_drops_total`
  * `fortigate_managed_switch_tx_errors_total`
  * `fortigate_managed_switch_tx_mcast_packets_total`
  * `fortigate_managed_switch_tx_oversize_total`
  * `fortigate_managed_switch_tx_packets_total`
  * `fortigate_managed_switch_tx_ucast_packets_total`
  * `fortigate_managed_switch_under_size_total`
    
## Usage

Example:

```
$ ./fortigate_exporter -auth-file ~/fortigate-key.yaml
# or
$ docker run -d -p 9710:9710 -v /path/to/fortigate-key.yaml:/config/fortigate-key.yaml quay.io/bluecmd/fortigate_exporter:master
```

Where `fortigate-key.yaml` contains pairs of FortiGate targets and API keys in the following format:

```
"https://my-fortigate":
  token: api-key-goes-here
  # If you have a smaller fortigate unit you might want
  # to exclude sensors as they do not have any
  probes:
    exclude:
      - System/SensorInfo

"https://my-other-fortigate:8443":
  token: api-key-goes-here
```

NOTE: Currently only token authentication is supported. FortiGate does not allow usage of tokens on non-HTTPS connections,
which means that currently you need HTTPS to be configured properly.

You can select which probes you want to run on a per target basis.

- Probes can be included or excluded under the optional `probes` section by defining `include` and/or `exclude` lists.
- Each probe name, that can be run by the fortigate exporter, is compared to the `include`/`exclude` lists.
- Inclusion/exclusion of a probe is based on a prefix match, therefore lists must contains entries starting with a probe name to be included/excluded.
- Prefix match is case sensitive.
- `include` list is evaluated before `exclude` list, therefore `exclude` list can exclude a previously included probe.

Example:

```
"https://my-fortigate":
  token: api-key-goes-here
  probes:
    include:
      - System
      - VPN
      - Firewall/Policies
      # Include only probes with name starting with: System or VPN + probe: Firewall/Policies
      # Other probes are excluded because there were not explictly included
"https://my-other-fortigate:8443":
  token: api-key-goes-here
  probes:
    exclude:
      - Wifi
      - Firewall/LoadBalance
      # Exclude probes with name starting with: Wifi + probe: Firewall/LoadBalance
      # All other probes are included by default because include list is empty
"https://my-other-orther-fortigate:8443":
  token: api-key-goes-here
  probes:
    include:
      - System
      - Firewall
    exclude:
      - System/LinkMonitor
      # Inlcude probes with name starting with: System and Firewall
      # Then exclude probe: System/LinkMonitor
```

Special cases:

- If `probes` isn't set or is empty, all probes will be run against the target.
- If `include` list is empty, by default, all probes will be selected to be run against the target.
- If `include` contains an entry `- ''`, then all probes are included (equivalent to not defining `include`)
- If `exclude` contains an entry `- ''`, then all probes are excluded (equivalent to not defining the target)


To probe a FortiGate, do something like `curl 'localhost:9710/probe?target=https://my-fortigate'`

### Dynamic configuration
In use cases where the Fortigates that is to be scraped through the fortigate-exporter is configured in 
Prometheus using some discovery method it becomes problematic that the `fortigate-key.yaml` configuration also
has to be updated for each fortigate, and that the fortigate-exporter needs to be restarted on each change. 
For that scenario the token can be passed as a query parameter, `token`, to the fortigate. 

Example:
```bash
curl 'localhost:9710/probe?target=https://192.168.2.31&token=ghi6eItWzWewgbrFMsazvBVwDjZzzb'
```
It is also possible to pass a `profile` query parameter. The value will match an entry in the `fortigate-key.yaml` 
file, but only to use the `probes` section for include/exclude directives.

Example:
```bash
curl 'localhost:9710/probe?target=https://192.168.2.31&token=ghi6eItWzWewgbrFMsazvBVwDjZzzb&profile=fs124e'
```
The `profile=fs124e` would match the following entry in `fortigate-key.yaml`.

Example:
```yaml
fs124e:
  # token: not used 
  probes:
    include:
      - System
      - Firewall
    exclude:
      - System/LinkMonitor
```



### Available CLI parameters

| flag  | default value  |  description  |
|---|---|---|
| -auth-file      | fortigate-key.yaml  | path to the location of the key file |
| -listen         | :9710  | address to listen for incoming requests  |
| -scrape-timeout | 30     | timeout in seconds  |
| -https-timeout  | 10     | timeout in seconds for establishment of HTTPS connections  |
| -insecure       | _not set_  | allows to turn off security validation of TLS certificates  |
| -extra-ca-certs | (none) | comma-separated files containing extra PEMs to trust for TLS connections in addition to the system trust store |
| -max-bgp-paths  | 10000  | Sets maximum amount of BGP paths to fetch, value is per IP stack version (IPv4 & IPv6) |
| -max-vpn-users  | 0      | Sets maximum amount of VPN users to fetch (0 eq. none by default) |

### FortiGate Configuration

Read permission is enough for Fortigate exporter purpose.
To improve security, limit permissions to required ones only (least privilege principle).

| probe name | permission | API URL |
|---|---|---|
| *Default Global*            | *any*              |api/v2/monitor/system/status |
|BGP/NeighborPaths/IPv4       | netgrp.route-cfg   |api/v2/monitor/router/bgp/paths |
|BGP/NeighborPaths/IPv6       | netgrp.route-cfg   |api/v2/monitor/router/bgp/paths6 |
|BGP/Neighbors/IPv4           | netgrp.route-cfg   |api/v2/monitor/router/bgp/neighbors |
|BGP/Neighbors/IPv6           | netgrp.route-cfg   |api/v2/monitor/router/bgp/neighbors6 |
|Firewall/IpPool              | fwgrp.policy       |api/v2/monitor/firewall/ippool |
|Firewall/LoadBalance         | fwgrp.others       |api/v2/monitor/firewall/load-balance |
|Firewall/Policies            | fwgrp.policy       |api/v2/monitor/firewall/policy/select<br>api/v2/monitor/firewall/policy6/select<br>api/v2/cmdb/firewall/policy<br>api/v2/cmdb/firewall/policy6 |
|License/Status               | *any*              |api/v2/monitor/license/status/select |
|Log/Fortianalyzer/Status     | loggrp.config      |api/v2/monitor/log/fortianalyzer |
|Log/Fortianalyzer/Queue      | loggrp.config      |api/v2/monitor/log/fortianalyzer-queue |
|Log/DiskUsage                | loggrp.config      |api/v2/monitor/log/current-disk-usage |
|System/AvailableCertificates | *any*              |api/v2/monitor/system/available-certificates |
|System/Fortimanager/Status   | sysgrp.cfg         |api/v2/monitor/system/fortimanager/status |
|System/HAStatistics          | sysgrp.cfg         |api/v2/monitor/system/ha-statistics<br>api/v2/cmdb/system/ha |
|System/Interface             | netgrp.cfg         |api/v2/monitor/system/interface/select |
|System/LinkMonitor           | sysgrp.cfg         |api/v2/monitor/system/link-monitor |
|System/Resource/Usage        | sysgrp.cfg         |api/v2/monitor/system/resource/usage |
|System/SensorInfo            | sysgrp.cfg         |api/v2/monitor/system/sensor-info |
|System/Status                | *any*              |api/v2/monitor/system/status |
|System/Time/Clock            | sysgrp.cfg         |api/v2/monitor/system/time |
|System/VDOMResources         | sysgrp.cfg         |api/v2/monitor/system/resource/usage |
|User/Fsso                    | authgrp            |api/v2/monitor/user/fsso |
|VPN/IPSec                    | vpngrp             |api/v2/monitor/vpn/ipsec |
|VPN/Ssl/Connections          | vpngrp             |api/v2/monitor/vpn/ssl |
|VPN/Ssl/Stats                | vpngrp             |api/v2/monitor/vpn/ssl/stats |
|VirtualWAN/HealthCheck       | netgrp.cfg         |api/v2/monitor/virtual-wan/health-check |
|Wifi/APStatus                | wifi               |api/v2/monitor/wifi/ap_status |
|Wifi/Clients                 | wifi               |api/v2/monitor/wifi/client |
|Wifi/ManagedAP               | wifi               |api/v2/monitor/wifi/managed_ap |
|Switch/ManagedSwitch         | switch	           |api/v2/monitor/switch-controller/managed-switch|
If you omit to grant some of these permissions you will receive log messages warning about
403 errors and relevant metrics will be unavailable, but other metrics will still work.
If you do not need some probes to be run, do not grant permission for them and use `include/exclude` feature (see `Usage` section).

The following example Admin Profile describes the permissions that needs to be granted
to the monitor user in order for all metrics to be available.

```
config system accprofile
    edit "monitor"
        # global scope will fail on non multi-VDOM firewall
        set scope global
        set authgrp read
        # As of FortiOS 6.2.1 it seems `fwgrp-permissions.other` is removed,
        # use 'fwgrp read' to get load balance servers metrics
        set fwgrp custom
        set loggrp custom
        set netgrp custom
        set sysgrp custom
        set vpngrp read
        set wifi read
        # will fail for most recent FortiOS
        set system-diagnostics disable
        config fwgrp-permission
            set policy read
            set others read
        end
        config netgrp-permission
            set cfg read
            set route-cfg read
        end
        config loggrp-permission
            set config read
        end
        config sysgrp-permission
            set cfg read
        end
    next
end
```

### Prometheus Configuration

An example configuration for Prometheus looks something like this:

```yaml
  - job_name: 'fortigate_exporter'
    metrics_path: /probe
    static_configs:
      - targets:
        - https://my-fortigate
        - https://my-other-fortigate:8443
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
        # Drop the https:// and port (if specified) for the 'instance=' label
        regex: '(?:.+)(?::\/\/)([^:]*).*'
      - target_label: __address__
        replacement: '[::1]:9710'
```
In above configuration only the targets and the replacement values needs to be changed as per your environment. Where target is URL of the Fortigate firewall and the replacement (at the bottom) will be the FQDN of system where the node exporter is running, example ```replacement: 'YourSystem.public.corp.com:9710'```

If using [Dynamic configuration](#dynamic-configuration):
```yaml
  - job_name: 'fortigate_exporter'
    metrics_path: /probe
    file_sd_configs:
      - files:
          - /etc/prometheus/file_sd/fws/*.yml
    params:
      profile:
      - fs124e
    relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [token]
      target_label: __param_token
    - source_labels: [__param_target]
      regex: '(?:.+)(?::\/\/)([^:]*).*'
      target_label: instance
    - target_label: __address__
      replacement: '[::1]:9710'
    - action: labeldrop
      regex: token
```
> Make sure to use the last labeldrop on the `token` label so that the tokens is not be part of your time series.

> Since `token` is a label it will be shown in the Prometheus webgui at `http://<your prometheus>:9090/targets`.
> 
> **Make sure you protect your Prometheus if you add the token part of your prometheus config** 
> 
> Some options to protect Prometheus:
> - Only expose UI to localhost --web.listen-address="127.0.0.1:9090"
> - Basic authentication access - https://prometheus.io/docs/guides/basic-auth/
> - **It is your responsibility!**

### Docker

You can either use the automatic builds on
[quay.io](https://quay.io/repository/bluecmd/fortigate_exporter) or build yourself
like this:

```bash
docker build -t fortigate_exporter .
docker run -d -p 9710:9710 -v /path/to/fortigate-key.yaml:/config/fortigate-key.yaml fortigate_exporter
```

#### docker-compose

```yaml
prometheus_fortigate_exporter:
  build: ./
  ports:
    - 9710:9710
  volumes:
    - /path/to/fortigate-key.yaml:/config/fortigate-key.yaml
  # Applying multiple parameters
  command: ["-auth-file", "/config/fortigate-key.yaml", "-insecure"]
  restart: unless-stopped
```

## Known Issues

This is a collection of known issues that for some reason cannot be fixed,
but might be possible to work around.

 * Probing causing [httpsd memory leak in FortiOS 6.2.x](https://github.com/bluecmd/fortigate_exporter/issues/62) ([Workaround](https://github.com/bluecmd/fortigate_exporter/issues/62#issuecomment-798602061))

## Missing Metrics?

Please [file an issue](https://github.com/bluecmd/fortigate_exporter/issues/new) describing what metrics you'd like to see.
Include as much details as possible please, e.g. how the perfect Prometheus metric would look for your use-case.

An alternative to using this exporter is to use generic SNMP polling, e.g. using a Prometheus SNMP exporter
([official](https://github.com/prometheus/snmp_exporter), [alternative](https://github.com/dhtech/snmpexporter)).
Note that there are limitations (e.g. [1](https://kb.fortinet.com/kb/documentLink.do?externalID=FD47703))
in what FortiGate supports querying via SNMP.

## Legal

Fortinet®, and FortiGate® are registered trademarks of Fortinet, Inc.

This is not an official Fortinet product.
