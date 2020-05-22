# fortigate_exporter

Prometheus exporter for Fortigate firewalls.

## Usage

Example:

```
./fortigate_exporter -api-key-file ~/fortigate.prom.key
```

Where `~/fortigate.prom.key` contains pairs of Fortigate targets and API keys in the following format:

```
https://my-fortigate             api-key-goes-here
http://my-unsafe-fortigate:8080  api-key-goes-here
```

To probe a Fortigate, do something like `curl 'localhost:9710/probe?target=https://my-fortigate'`

