# fortigate_exporter

![Go](https://github.com/bluecmd/fortigate_exporter/workflows/Go/badge.svg)
![Docker](https://github.com/bluecmd/fortigate_exporter/workflows/Docker/badge.svg)
[![Docker Repository on Quay](https://quay.io/repository/bluecmd/fortigate_exporter/status "Docker Repository on Quay")](https://quay.io/repository/bluecmd/fortigate_exporter)

Prometheus exporter for FortiGate® firewalls.

**NOTE:** This is not an official Fortinet product, it is developed fully independently by professionals and hobbyists alike.

  * [Supported Metrics](#supported-metrics)
  * [Usage](#usage)
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

 * `fortigate_version_info`
 * `fortigate_cpu_usage_ratio`
 * `fortigate_memory_usage_ratio`
 * `fortigate_current_sessions`
 * `fortigate_license_vdom_usage`
 * `fortigate_license_vdom_max`

Per-VDOM:

 * `fortigate_vdom_cpu_usage_ratio`
 * `fortigate_vdom_memory_usage_ratio`
 * `fortigate_vdom_current_sessions`
 * `fortigate_policy_active_sessions`
 * `fortigate_policy_bytes_total`
 * `fortigate_policy_hit_count_total`
 * `fortigate_policy_packets_total`
 * `fortigate_interface_link_up`
 * `fortigate_interface_speed_bps`
 * `fortigate_interface_transmit_packets_total`
 * `fortigate_interface_receive_packets_total`
 * `fortigate_interface_transmit_bytes_total`
 * `fortigate_interface_receive_bytes_total`
 * `fortigate_interface_transmit_errors_total`
 * `fortigate_interface_receive_errors_total`
 * `fortigate_vpn_connections`
 * `fortigate_ipsec_tunnel_receive_bytes_total`
 * `fortigate_ipsec_tunnel_transmit_bytes_total`
 * `fortigate_ipsec_tunnel_up`
 * `fortigate_wifi_access_points`
 * `fortigate_wifi_fabric_clients`
 * `fortigate_wifi_fabric_max_allowed_clients`

 Per-HA-Member and VDOM:
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
 * `fortigate_bgp_neighbor_ipv4_info`
 * `fortigate_bgp_neighbor_ipv6_info`
 * `fortigate_bgp_neighbor_ipv4_paths`
 * `fortigate_bgp_neighbor_ipv6_paths`
 * `fortigate_bgp_neighbor_ipv4_best_paths`
 * `fortigate_bgp_neighbor_ipv6_best_paths`

 Per-VirtualServer and VDOM:
 * `fortigate_lb_virtual_server_info`

 Per-RealServer for each VirtualServer and VDOM:
 * `fortigate_lb_real_server_info`
 * `fortigate_lb_real_server_mode`
 * `fortigate_lb_real_server_status`
 * `fortigate_lb_real_server_active_sessions`
 * `fortigate_lb_real_server_rtt_seconds`
 * `fortigate_lb_real_server_processed_bytes_total`

 Per-Certificate
 * `fortigate_certificate_info`
 * `fortigate_certificate_valid_from_seconds`
 * `fortigate_certificate_valid_to_seconds`
 * `fortigate_certificate_cmdb_references`

Per-VDOM and Wifi-Client
 * `fortigate_wifi_client_info`
 * `fortigate_wifi_client_data_rate_bps`
 * `fortigate_wifi_client_bandwidth_rx_bps`
 * `fortigate_wifi_client_bandwidth_tx_bps`
 * `fortigate_wifi_client_signal_strength_dBm`
 * `fortigate_wifi_client_signal_noise_dBm`
 * `fortigate_wifi_client_tx_discard_ratio`
 * `fortigate_wifi_client_tx_retries_ratio`

Per-VDOM and managed access point:
 * `fortigate_wifi_managed_ap_info`
 * `fortigate_wifi_managed_ap_join_time_seconds`
 * `fortigate_wifi_managed_ap_cpu_usage_ratio`
 * `fortigate_wifi_managed_ap_memory_free_bytes`
 * `fortigate_wifi_managed_ap_memory_bytes_total`

Per-VDOM, managed access point and radio:
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
 * `fortigate_wifi_managed_ap_interface_rx_bytes_total`
 * `fortigate_wifi_managed_ap_interface_tx_bytes_total`
 * `fortigate_wifi_managed_ap_interface_rx_packets_total`
 * `fortigate_wifi_managed_ap_interface_tx_packets_total`
 * `fortigate_wifi_managed_ap_interface_rx_errors_total`
 * `fortigate_wifi_managed_ap_interface_tx_errors_total`
 * `fortigate_wifi_managed_ap_interface_rx_dropped_packets_total`
 * `fortigate_wifi_managed_ap_interface_tx_dropped_packets_total`

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
"https://my-other-fortigate:8443":
  token: api-key-goes-here
```

NOTE: Currently only token authentication is supported. FortiGate does not allow usage of tokens on non-HTTPS connections,
which means that currently you need HTTPS to be configured properly.

You can select which probes or probe categories you want to run per target, for example:

```
"https://my-fortigate":
  token: api-key-goes-here
  probes:
    - System
"https://my-other-fortigate:8443":
  token: api-key-goes-here
  probes:
    - BGPNeighborsIPv4
    - Wifi
```

If `probes` isn't set or is empty, all probes will be run against the target.

To probe a FortiGate, do something like `curl 'localhost:9710/probe?target=https://my-fortigate'`

### Available CLI parameters

| flag  | default value  |  description  |
|---|---|---|
| -auth-file      | fortigate-key.yaml  | path to the location of the key file |
| -listen         | :9710  | address to listen for incoming requests  |
| -scrape-timeout | 30     | timeout in seconds  |
| -https-timeout  | 10     | timeout in seconds for establishment of HTTPS connections  |
| -insecure       | false  | allows to turn off security validation of TLS certificates  |
| -extra-ca-certs | (none) | comma-separated files containing extra PEMs to trust for TLS connections in addition to the system trust store |
| -max-bgp-paths  | 10000  | Sets maximum amount of BGP paths to fetch, value is per IP stack version (IPv4 a& IPv6) |
### FortiGate Configuration

The following example Admin Profile describes the permissions that needs to be granted
to the monitor user in order for all metrics to be available.
If you omit to grant some of these permissions you will receive log messages warning about
403 errors and relevant metrics will be unavailable, but other metrics will still work.

```
config system accprofile
    edit "monitor"
        set scope global
        set secfabgrp read
        set netgrp custom
        # As of FortiOS 6.2.1 it seems `fwgrp-permissions.other` is removed,
        # use 'fwgrp read' to get load balance servers metrics
        set fwgrp custom
        set vpngrp read
        set system-diagnostics disable
        config netgrp-permission
            set cfg read
        end
        config sysgrp-permission
            # Sysgrp.cfg is needed for the following optional functions:
            # - HA group name
            # If you do not wish to grant this permission, the relevant
            # labels/metrics will be absent.
            set cfg read
            # If you wish to collect ipv6 bgp neighbours, add this:
            set route-cfg read
        end
        config fwgrp-permission
            set policy read
            # fwgrp.other is need for load balance servers
            # set other read
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
  command: ["-auth-file", "/config/fortigate-key.yaml", "-insecure", "true"]
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
