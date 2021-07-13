# api/v2/monitor/wifi/client

[
    {
      "http_method":"GET",
      "results":[
        {
          "sta_ip":"0.0.0.0",
          "sta_ip6":[
            "::"
          ],
          "sta_rate":"130.0 mbps",
          "sta_snr":"36 dB",
          "sta_idle_time":"61 sec",
          "sta_assoc_time":"04\/20\/21 19:24",
          "sta_mac":"00:00:00:00:00:00",
          "sta_auth":"pass",
          "sta_manuf":"Nintendo Co., Ltd.",
          "ip":"0.0.0.0",
          "wtp_name":"3rd Floor",
          "wtp_id":"FP221E0000000000",
          "wtp_radio":2,
          "wtp_ip":"0.0.0.0",
          "vap_name":"example-SSID-wifi",
          "ssid":"example-SSID",
          "mac":"00:00:00:AA:00:00",
          "authentication":"pass",
          "captive_portal_authenticated":0,
          "manufacturer":"Nintendo Co., Ltd.",
          "data_rate":1300,
          "data_rate_bps":130000000,
          "snr":36,
          "idle_time":61,
          "association_time":1618939468,
          "bandwidth_tx":0,
          "bandwidth_rx":0,
          "lan_authenticated":false,
          "channel":44,
          "signal":-59,
          "vci":"",
          "host":"",
          "security":10,
          "security_str":"wpa2_only_personal",
          "encrypt":1,
          "noise":-95,
          "radio_type":"802.11ac",
          "mimo":"2x2",
          "vlan_id":21,
          "tx_discard_percentage":0,
          "tx_retry_percentage":0,
          "mpsk_name":"",
          "health":{
            "signal_strength":{
              "value":-59,
              "severity":"good"
            },
            "snr":{
              "value":36,
              "severity":"good"
            },
            "band":{
              "value":"5ghz",
              "severity":"good"
            },
            "transmission_retry":{
              "value":0,
              "severity":"good"
            },
            "transmission_discard":{
              "value":0,
              "severity":"good"
            }
          }
        },
        {
          "sta_ip":"0.0.0.0",
          "sta_ip6":[
            "::"
          ],
          "sta_rate":"1.0 mbps",
          "sta_snr":"36 dB",
          "sta_idle_time":"21 sec",
          "sta_assoc_time":"04\/19\/21 21:44",
          "sta_mac":"00:00:00:00:00:00",
          "sta_auth":"pass",
          "sta_manuf":"Espressif Inc.",
          "ip":"0.0.0.0",
          "wtp_name":"2nd Floor",
          "wtp_id":"FP221E0000000000",
          "wtp_radio":1,
          "wtp_ip":"0.0.0.0",
          "vap_name":"example-SSID-wifi",
          "ssid":"example-SSID",
          "mac":"00:00:00:00:00:00",
          "hostname":"wled-WLED",
          "authentication":"pass",
          "captive_portal_authenticated":0,
          "manufacturer":"Espressif Inc.",
          "data_rate":10,
          "data_rate_bps":1000000,
          "snr":36,
          "idle_time":21,
          "association_time":1618861445,
          "bandwidth_tx":0,
          "bandwidth_rx":0,
          "lan_authenticated":false,
          "channel":6,
          "signal":-59,
          "vci":"",
          "host":"wled-WLED",
          "security":10,
          "security_str":"wpa2_only_personal",
          "encrypt":1,
          "noise":-95,
          "radio_type":"802.11n",
          "mimo":"1x1",
          "vlan_id":21,
          "tx_discard_percentage":0,
          "tx_retry_percentage":0,
          "mpsk_name":"",
          "health":{
            "signal_strength":{
              "value":-59,
              "severity":"good"
            },
            "snr":{
              "value":36,
              "severity":"good"
            },
            "band":{
              "value":"24ghz",
              "severity":"fair"
            },
            "transmission_retry":{
              "value":0,
              "severity":"good"
            },
            "transmission_discard":{
              "value":0,
              "severity":"good"
            }
          }
        }
      ],
      "vdom":"root",
      "path":"wifi",
      "name":"client",
      "status":"success",
      "serial":"FGT61FT000000000",
      "version":"v6.4.5",
      "build":1828
    }
]
