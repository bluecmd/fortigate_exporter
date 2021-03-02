# fortigate_exporter

![Go](https://github.com/bluecmd/fortigate_exporter/workflows/Go/badge.svg)
![Docker](https://github.com/bluecmd/fortigate_exporter/workflows/Docker/badge.svg)
[![Docker Repository on Quay](https://quay.io/repository/bluecmd/fortigate_exporter/status "Docker Repository on Quay")](https://quay.io/repository/bluecmd/fortigate_exporter)

Prometheus exporter for Fortigate firewalls.

## Supported Metrics

Right now the exporter supports a quite limited set of metrics, but it is very easy to add!
Open an issue if your favorite metric is missing.

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

  Per-Certificate
 * `fortigate_certificate_info`
 * `fortigate_certificate_valid_from_seconds`
 * `fortigate_certificate_valid_to_seconds`
 * `fortigate_certificate_cmdb_references`

## Usage

Example:

```
$ ./fortigate_exporter -auth-file ~/fortigate-key.yaml
# or
$ docker run -d -p 9710:9710 -v /path/to/fortigate-key.yaml:/config/fortigate-key.yaml quay.io/bluecmd/fortigate_exporter
```

Where `fortigate-key.yaml` contains pairs of Fortigate targets and API keys in the following format:

```
"https://my-fortigate":
  token: api-key-goes-here
"https://my-other-fortigate:8443":
  token: api-key-goes-here
```

NOTE: Currently only token authentication is supported. Fortigate does not allow usage of tokens on non-HTTPS connections,
which means that currently you need HTTPS to be configured properly.

To probe a Fortigate, do something like `curl 'localhost:9710/probe?target=https://my-fortigate'`

## Available CLI parameters
| flag  | default value  |  description  |
|---|---|---|
| -auth-file      | /config/fortigate-key.yaml  | path to the location of the key file |
| -listen         | :9710  | address to listen for incoming requests  |
| -scrape-timeout | 30     | timeout in seconds  |
| -https-timeout  | 10     | timeout in seconds for establishment of HTTPS connections  |
| -insecure       | false  | allows to turn off security validation of TLS certificates  |
| -extra-ca-certs | (none) | comma-separated files containing extra PEMs to trust for TLS connections in addition to the system trust store |

## Fortigate Configuration

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
        end
        config fwgrp-permission
            set policy read
        end
    next
end
```

## Prometheus Configuration

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

## Docker

You can either use the automatic builds on
[quay.io](https://quay.io/repository/bluecmd/fortigate_exporter) or build yourself
like this:

```bash
docker build -t fortigate_exporter .
docker run -d -p 9710:9710 -v /path/to/fortigate-key.yaml:/config/fortigate-key.yaml fortigate_exporter
```

### docker-compose

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

## Missing Metrics?

Please [file an issue](https://github.com/bluecmd/fortigate_exporter/issues/new) describing what metrics you'd like to see.
Include as much details as possible please, e.g. how the perfect Prometheus metric would look for your use-case.

An alternative to using this exporter is to use generic SNMP polling, e.g. using a Prometheus SNMP exporter
([official](https://github.com/prometheus/snmp_exporter), [alternative](https://github.com/dhtech/snmpexporter)).
Note that there are limitations (e.g. [1](https://kb.fortinet.com/kb/documentLink.do?externalID=FD47703))
in what Fortigate supports querying via SNMP.
