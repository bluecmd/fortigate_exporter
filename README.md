# fortigate_exporter

![Go](https://github.com/bluecmd/fortigate_exporter/workflows/Go/badge.svg)

Prometheus exporter for Fortigate firewalls.

## Supported Metrics

Right now the exporter supports a quite limited set of metrics, but it is very easy to add!
Open an issue if your favorite metric is missing.

Supported metrics right now as follows.

Global:

 * `fortigate_cpu_usage_ratio`
 * `fortigate_memory_usage_ratio`
 * `fortigate_current_sessions`

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

## Usage

Example:

```
./fortigate_exporter -auth-file ~/fortigate-key.yaml
```

Where `~/fortigate-key.yaml` contains pairs of Fortigate targets and API keys in the following format:

```
"https://my-fortigate":
  token: api-key-goes-here
"https://my-other-fortigate:8443":
  token: api-key-goes-here
```

NOTE: Currently only token authentication is supported. Fortigate does not allow usage of tokens on non-HTTPS connections,
which means that currently you need HTTPS to be configured properly.

To probe a Fortigate, do something like `curl 'localhost:9710/probe?target=https://my-fortigate'`

## Prometheus Configuration

An example configuration for Prometheus looks something like this:

```yaml
  - job_name: 'fortigate_exporter'
    metrics_path: /probe
    static_configs:
      - targets:
        - https://fortigate.s6.network
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
        # Drop the https:// and port (if specified) for the 'instance=' label
        regex: '(?:.+)(?::\/\/)([^:]*)'
      - target_label: __address__
        replacement: '[::1]:9710'
```
