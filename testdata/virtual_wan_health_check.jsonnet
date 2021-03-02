# api/v2/monitor/virtual-wan/health-check?vdom=*
[
    {
        "http_method": "GET",
        "results": {
            "Internet Check": {
                "WAN1_VL300": {
                    "status": "up",
                    "latency": 5.611332893371582,
                    "jitter": 0.031166711822152138,
                    "packet_loss": 0,
                    "packet_sent": 306958,
                    "packet_received": 306895,
                    "sla_targets_met": [
                        1
                    ],
                    "session": 710,
                    "tx_bandwidth": 117296,
                    "rx_bandwidth": 257003,
                    "state_changed": 1614107800
                },
                "wan2": {
                    "status": "disable"
                }
            }
        },
        "vdom": "root",
        "path": "virtual-wan",
        "name": "health-check",
        "status": "success",
        "serial": "FGT60EXXXXXXXXXX",
        "version": "v6.4.5",
        "build": 1828
    }
]
