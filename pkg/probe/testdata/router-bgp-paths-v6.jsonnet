# api/v2/monitor/router/bgp/paths6?vdom=*&count=10000
[
  {
    "http_method":"GET",
    "results":[
      {
        "nlri_prefix":"::",
        "nlri_prefix_len":0,
        "learned_from":"fd00::1",
        "next_hop_local":"ffd00::2",
        "next_hop_global":"fd00::1",
        "origin":"igp",
        "is_best":true
      },{
        "nlri_prefix":"fd02::",
        "nlri_prefix_len":0,
        "learned_from":"fd00::1",
        "next_hop_local":"ffd00::2",
        "next_hop_global":"fd00::1",
        "origin":"igp",
        "is_best":false
      },{
        "nlri_prefix":"fd03::",
        "nlri_prefix_len":0,
        "learned_from":"fd00::1",
        "next_hop_local":"ffd00::2",
        "next_hop_global":"fd00::1",
        "origin":"igp",
        "is_best":true
      },
      {
        "nlri_prefix":"2001:678:f40::",
        "nlri_prefix_len":48,
        "learned_from":"::",
        "next_hop_local":"::",
        "next_hop_global":"::",
        "origin":"igp",
        "is_best":true
      }
    ],
    "vdom":"root",
    "path":"router",
    "name":"bgp",
    "action":"paths",
    "status":"success",
    "serial":"FGT61FT000000000",
    "version":"v7.0.0",
    "build":66
  }
]