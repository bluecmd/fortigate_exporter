# api/v2/monitor/system/available-certificates?vdom=*
[
  {
    "http_method":"GET",
    "results":[
      {
        "name":"Fortinet_CA_SSL",
        "source":"factory",
        "comments":"This is the default CA certificate the SSL Inspection will use when generating new server certificates.",
        "range":"global",
        "exists":true,
        "is_ssl_server_cert":true,
        "is_proxy_ssl_cert":true,
        "is_general_allowable_cert":true,
        "is_default_local":false,
        "is_built_in":true,
        "is_wifi_cert":false,
        "is_deep_inspection_cert":true,
        "has_valid_cert_key":true,
        "key_type":"RSA",
        "key_size":2048,
        "is_local_ca_cert_strict":true,
        "is_local_ca_cert":true,
        "type":"local-ca",
        "status":"valid",
        "valid_from":1472285182,
        "valid_from_raw":"2016-08-27 08:06:22  GMT",
        "valid_to":1787904382,
        "valid_to_raw":"2026-08-28 08:06:22  GMT",
        "signature_algorithm":"SHA256",
        "subject":{
          "C":"US",
          "ST":"California",
          "L":"Sunnyvale",
          "O":"Fortinet",
          "OU":"Certificate Authority",
          "CN":"FGT61E4QXXXXXXXX",
          "emailAddress":"support@fortinet.com"
        },
        "subject_raw":"C = US, ST = California, L = Sunnyvale, O = Fortinet, OU = Certificate Authority, CN = FGT61E4QXXXXXXXX, emailAddress = support@fortinet.com",
        "issuer":{
          "C":"US",
          "ST":"California",
          "L":"Sunnyvale",
          "O":"Fortinet",
          "OU":"Certificate Authority",
          "CN":"FGT61E4QXXXXXXXX",
          "emailAddress":"support@fortinet.com"
        },
        "issuer_raw":"C = US, ST = California, L = Sunnyvale, O = Fortinet, OU = Certificate Authority, CN = FGT61E4QXXXXXXXX, emailAddress = support@fortinet.com",
        "fingerprint":"F9:96:25:CC:E0:E9:08:F7:1C",
        "version":3,
        "is_ca":true,
        "serial_number":"28:D7:XX:XX",
        "ext":[
          {
            "name":"X509v3 Basic Constraints",
            "data":"CA:TRUE",
            "critical":false
          }
        ],
        "q_path":"vpn.certificate",
        "q_name":"local",
        "q_ref":5,
        "q_static":true,
        "q_type":155
      }
    ],
    "vdom":"root",
    "path":"system",
    "name":"available-certificates",
    "status":"success",
    "serial":"FGT61E4QXXXXXXXX",
    "version":"v6.2.3",
    "build":1066
  }
]