# PromQL usage examples for Fortigate Exporter

PromQL is an advanced language and proper usage of it requires firstly to understand
the language itself and secondly how the metrics are organized.

To help the users of this exporter, here are some examples to illustrate
how one could use the power of PromQL to query interesting facts about one's
environment.

For an introduction to the PromQL basics, see the official Prometheus
[documentation](https://prometheus.io/docs/prometheus/latest/querying/examples/).

## Most Active Firewall Policies

Using the `topk` function it is easy to get the most active firewall policies
based on whatever ranking function you want.

Example: `topk(3, rate(fortigate_policy_bytes_total[15m])) * 8`

This will return the top 3 most active policies based upon the bytes transfered on
average the last 15m. The output will be bits/s.

| Element | Value| 
|---------|------|
`{id="9",instance="fgt-a",job="fortigate",name="SSH",protocol="ipv6",uuid="5cd4b62e-4904-51eb-b4a9-f52e75461e52",vdom="bluecmd"}` |	102805.48826815643
`{id="28",instance="fgt-a",job="fortigate",name="fortigate exporter",protocol="ipv6",uuid="8753dcd2-4a07-51eb-bf78-c61aa31a8e1e",vdom="bluecmd"}` |	34443.977653631286
`{id="11",instance="fgt-a",job="fortigate",name="tera cluster mgmt",protocol="ipv4",uuid="19a2c192-4905-51eb-9d14-5d1249566588",vdom="bluecmd"}` | 15239.401117318435

## Adding Version Information

Using `group_left` one can add data from other metadata metrics like `fortigate_version_info`.

Example: `fortigate_memory_usage_ratio * on (instance) group_left (version) fortigate_version_info`

| Element | Value| 
|---------|------|
`{instance="fgt-test",job="fortigate",version="v6.4.5"}`	| 0.12
`{instance="fgt-a",job="fortigate",version="v6.4.4"}`	| 0.23
`{instance="fgt-b",job="fortigate",version="v6.4.4"}`	| 0.16

## Expiring Certificates

Dealing with certificates can unfortunately be quite complex.
However, by using PromQL it is possible to construct a query that returns
all certificates expiring within the coming 90 days.

Example:

```
floor( # return whole days
  (
    (
      fortigate_certificate_valid_to_seconds and 
      fortigate_certificate_cmdb_references > 0 and  # only include certificates that are used for something
      on (instance,name,vdom) fortigate_certificate_info{status="valid"}  # we do not care about things like CSRs
    )
    - time()
  ) / 86400 # convert seconds to days
)
< 90 # number of days to filter for
```

| Element | Value| 
|---------|------|
`{instance="fgt-test",job="fortigate",name="LetsEncrypt-2021-03-13",scope="global",source="user",vdom="root"}` |	89
`{instance="fgt-a",job="fortigate",name="LetsEncrypt-2021-01-05",scope="global",source="user",vdom="root"}` |	21
`{instance="fgt-b",job="fortigate",name="LetsEncrypt-2021-01-05",scope="global",source="user",vdom="root"}` |	21
