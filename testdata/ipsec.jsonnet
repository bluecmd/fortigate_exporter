# api/v2/monitor/vpn/ipsec?vdom=*
[
  {
    "http_method":"GET",
    "results":[
      {
        "proxyid":[
          {
            "proxy_src":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"up",
            "p2name":"tunnel_1-sub",
            "p2serial":1,
            "expire":11279,
            "incoming_bytes": 14298240,
            "outgoing_bytes": 14248560
          },
          {
            "proxy_src":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"down",
            "p2name":"tunnel_1-sub",
            "p2serial":12,
            "expire":11279,
            "incoming_bytes": 14298240,
            "outgoing_bytes": 14248560
          }          
        ],
        "name":"tunnel_1",
        "comments":"",
        "wizard-type":"custom",
        "creation_time":270801,
        "type":"automatic",
        "incoming_bytes": 14298240,
        "outgoing_bytes": 14248560,
        "rgwy":"1.2.3.4"
      }
    ],
    "vdom":"root",
    "path":"vpn",
    "name":"ipsec",
    "status":"success",
    "serial":"FGT61FT000000000",
    "version":"v6.0.10",
    "build":365
  },
]
