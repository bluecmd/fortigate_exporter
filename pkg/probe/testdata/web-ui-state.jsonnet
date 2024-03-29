# api/v2/monitor/web-ui/state
{
  "http_method":"GET",
  "results":{
    "model_name":"FortiGate",
    "model_number":"61F",
    "hostname":"fortigate",
    "model":"FGT61F",
    "model_subtype":"",
    "model_level":"low",
    "admin_using_default_password":false,
    "admin":{
      "name":"prometheus",
      "login_name":"prometheus",
      "profile":{
        "name":"prometheus",
        "q_origin_key":"prometheus",
        "scope":"global",
        "comments":"",
        "secfabgrp":"read",
        "ftviewgrp":"read",
        "authgrp":"read",
        "sysgrp":"read",
        "netgrp":"read",
        "loggrp":"read",
        "fwgrp":"read",
        "vpngrp":"read",
        "utmgrp":"read",
        "wanoptgrp":"read",
        "wifi":"read",
        "netgrp-permission":{
          "cfg":"none",
          "packet-capture":"none",
          "route-cfg":"none"
        },
        "sysgrp-permission":{
          "admin":"none",
          "upd":"none",
          "cfg":"none",
          "mnt":"none"
        },
        "fwgrp-permission":{
          "policy":"none",
          "address":"none",
          "service":"none",
          "schedule":"none",
          "others":"none"
        },
        "loggrp-permission":{
          "config":"none",
          "data-access":"none",
          "report-access":"none",
          "threat-weight":"none"
        },
        "utmgrp-permission":{
          "antivirus":"none",
          "ips":"none",
          "webfilter":"none",
          "emailfilter":"none",
          "data-loss-prevention":"none",
          "file-filter":"none",
          "application-control":"none",
          "icap":"none",
          "voip":"none",
          "waf":"none",
          "dnsfilter":"none",
          "endpoint-control":"none"
        },
        "admintimeout-override":"disable",
        "admintimeout":10,
        "system-diagnostics":"enable"
      },
      "global_admin":true,
      "super_admin":false,
      "ignore_release_overview":"",
      "ignore_invalid_signature_version":"",
      "dashboard_template":"",
      "guest_admin":false,
      "fmg_admin":false,
      "sso_login_type":"none",
      "remote_admin":false,
      "pki_admin":false,
      "vdoms":[
        "",
        "main",
        "root"
      ],
      "vdom_info":{
        "main":{
          "central_nat_enabled":true,
          "transparent_mode":false,
          "ngfw_mode":"profile-based",
          "features":{
            "gui-icap":false,
            "gui-implicit-policy":true,
            "gui-dns-database":true,
            "gui-load-balance":false,
            "gui-multicast-policy":false,
            "gui-dos-policy":false,
            "gui-object-colors":true,
            "gui-voip-profile":false,
            "gui-ap-profile":true,
            "gui-security-profile-group":false,
            "gui-local-in-policy":true,
            "gui-local-reports":true,
            "gui-explicit-proxy":true,
            "gui-dynamic-routing":true,
            "gui-sslvpn-personal-bookmarks":false,
            "gui-sslvpn-realms":false,
            "gui-policy-based-ipsec":false,
            "gui-threat-weight":true,
            "gui-spamfilter":false,
            "gui-file-filter":true,
            "gui-application-control":true,
            "gui-ips":true,
            "gui-endpoint-control":true,
            "gui-endpoint-control-advanced":false,
            "gui-dhcp-advanced":true,
            "gui-vpn":true,
            "gui-wireless-controller":false,
            "gui-switch-controller":false,
            "gui-fortiap-split-tunneling":false,
            "gui-traffic-shaping":true,
            "gui-wan-load-balancing":false,
            "gui-antivirus":false,
            "gui-webfilter":false,
            "gui-videofilter":true,
            "gui-dnsfilter":true,
            "gui-waf-profile":false,
            "gui-advanced-policy":true,
            "gui-allow-unnamed-policy":true,
            "gui-email-collection":false,
            "gui-multiple-interface-policy":true,
            "gui-policy-disclaimer":false,
            "gui-ztna":false
          },
          "virtual_wire_pair_count":0,
          "is_management_vdom":false,
          "log_device_state":{
            "memory":{
              "is_available":true,
              "is_enabled":false,
              "is_default":false,
              "is_ha_supported":true
            },
            "disk":{
              "is_available":true,
              "is_enabled":true,
              "is_loggable":true,
              "num_ssds_available":1,
              "disabled_by_default":false,
              "is_ha_supported":true,
              "is_fortiview_supported":true,
              "fortiview_weekly_data":false
            },
            "fortianalyzer":{
              "is_available":true,
              "is_enabled":false,
              "overrides_global_faz":false
            },
            "fortianalyzer_cloud":{
              "is_available":false,
              "is_enabled":false,
              "overrides_global_faz_cloud":false
            },
            "forticloud":{
              "is_available":true,
              "is_enabled":true,
              "is_faz_cloud":false
            }
          },
          "log_device_default":"forticloud",
          "resolve_hostnames":true
        },
        "root":{
          "central_nat_enabled":false,
          "transparent_mode":false,
          "ngfw_mode":"profile-based",
          "features":{
            "gui-icap":false,
            "gui-implicit-policy":true,
            "gui-dns-database":false,
            "gui-load-balance":false,
            "gui-multicast-policy":false,
            "gui-dos-policy":false,
            "gui-object-colors":true,
            "gui-voip-profile":false,
            "gui-ap-profile":true,
            "gui-security-profile-group":false,
            "gui-local-in-policy":true,
            "gui-local-reports":true,
            "gui-explicit-proxy":false,
            "gui-dynamic-routing":true,
            "gui-threat-weight":false,
            "gui-spamfilter":false,
            "gui-file-filter":true,
            "gui-application-control":true,
            "gui-ips":true,
            "gui-endpoint-control":true,
            "gui-endpoint-control-advanced":false,
            "gui-dhcp-advanced":true,
            "gui-vpn":false,
            "gui-wireless-controller":false,
            "gui-switch-controller":false,
            "gui-fortiap-split-tunneling":false,
            "gui-traffic-shaping":true,
            "gui-wan-load-balancing":false,
            "gui-antivirus":false,
            "gui-webfilter":false,
            "gui-videofilter":true,
            "gui-dnsfilter":false,
            "gui-waf-profile":false,
            "gui-fortiextender-controller":false,
            "gui-advanced-policy":true,
            "gui-allow-unnamed-policy":true,
            "gui-email-collection":false,
            "gui-multiple-interface-policy":true,
            "gui-policy-disclaimer":false,
            "gui-ztna":false
          },
          "virtual_wire_pair_count":0,
          "is_management_vdom":true,
          "log_device_state":{
            "memory":{
              "is_available":true,
              "is_enabled":false,
              "is_default":false,
              "is_ha_supported":true
            },
            "disk":{
              "is_available":true,
              "is_enabled":true,
              "is_loggable":true,
              "num_ssds_available":1,
              "disabled_by_default":false,
              "is_ha_supported":true,
              "is_fortiview_supported":true,
              "fortiview_weekly_data":false
            },
            "fortianalyzer":{
              "is_available":true,
              "is_enabled":false,
              "overrides_global_faz":false
            },
            "fortianalyzer_cloud":{
              "is_available":false,
              "is_enabled":false,
              "overrides_global_faz_cloud":false
            },
            "forticloud":{
              "is_available":true,
              "is_enabled":true,
              "is_faz_cloud":false
            }
          },
          "log_device_default":"forticloud",
          "resolve_hostnames":true
        }
      }
    },
    "fext_enabled":true,
    "fext_vlan_mode":false,
    "snapshot_utc_time":1659857566000,
    "utc_last_reboot":1657116965000,
    "time_zone_offset":0,
    "time_zone_text":"(GMT) Greenwich Mean Time",
    "time_zone_db_name":"Etc\/GMT",
    "centrally_managed":false,
    "fortimanager_backup_mode":false,
    "fips_cc_enabled":false,
    "fips_ciphers_enabled":false,
    "vdom_mode":"multi-vdom",
    "management_vdom":"root",
    "conserve_mode":false,
    "image_sign_status":"certified",
    "bios_security_level":1,
    "need_fs_check":false,
    "carrier_mode":false,
    "has_hyperscale_license":false,
    "csf_enabled":false,
    "csf_group_name":"",
    "csf_upstream_ip":"",
    "csf_sync_mode":"default",
    "csf_object_sync_mode":"local",
    "ha_mode":0,
    "is_ha_master":1,
    "ngfw_mode":"profile-based",
    "forced_low_crypto":false,
    "has_log_disk":true,
    "has_local_config_revisions":true,
    "lenc_mode":false,
    "usg_mode":false,
    "admin_https_redirection":true,
    "config_save_mode":"automatic",
    "debug_supported_daemons":[
      "node",
      "httpsd",
      "cmdb",
      "miglogd",
      "csfd",
      "sslvpnd"
    ],
    "is_vm":false,
    "theme":"jade",
    "language_code":"en",
    "cmgmt_override_cookie_name":"REDACTED",
    "ccsrf_token_cookie_name":"REDACTED",
    "file_downloading_cookie_name":"REDACTED",
    "autoscale_config_rec_override_cookie_name":"REDACTED",
    "initial_vdom":"",
    "timeout_minutes":60,
    "features":{
      "gui-ipv6":true,
      "gui-replacement-message-groups":false,
      "gui-local-out":false,
      "gui-certificates":true,
      "gui-custom-language":false,
      "gui-wireless-opensecurity":false,
      "gui-display-hostname":false,
      "gui-fortigate-cloud-sandbox":false,
      "gui-firmware-upgrade-warning":true,
      "gui-allow-default-hostname":false,
      "gui-forticare-registration-setup-warning":true,
      "gui-cdn-usage":true,
      "switch-controller":true,
      "wireless-controller":true,
      "fortiextender":true,
      "fortitoken-cloud-service":true
    },
    "date_format":'yyyy\/MM\/dd',
    "date_format_device":"system",
    "security_rating_result_submission":true,
    "security_rating_run_on_schedule":true
  },
  "vdom":"main",
  "path":"web-ui",
  "name":"state",
  "action":"",
  "status":"success",
  "serial":"FGT61FT000000000",
  "version":"v7.0.6",
  "build": 366
}
