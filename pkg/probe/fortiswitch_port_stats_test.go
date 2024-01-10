package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSwitchPortStats(t *testing.T) {
	c := newFakeClient()
	//c.prepare("api/v2/monitor/switch-controller/managed-switch?port_stats=true", "testdata/fsw-interface-chained.jsonnet")
	c.prepare("api/v2/monitor/switch-controller/managed-switch", "testdata/fsw-interface.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSwitchPortStats, c, r) {
		t.Errorf("probeSwitchPortStats() returned non-success")
	}

	em := `
	# HELP fortiswitch_status Whether the switch is connected or not
	# TYPE fortiswitch_status gauge
	fortiswitch_status{connection_from="169.254.2.2",fgt_peer_intf_name="fortilink",name="S108EN0000000000",state="Authorized",vdom="root"} 1
	fortiswitch_status{connection_from="169.254.2.3",fgt_peer_intf_name="fortilink",name="S108EP4N00000000",state="Authorized",vdom="root"} 1
	# HELP fortiswitch_port_status Whether the switch port is up or not
	# TYPE fortiswitch_port_status gauge
	fortiswitch_port_status{duplex="",interface="port10",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port10",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port2",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port3",name="S108EP4N00000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port5",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port6",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port6",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port7",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port9",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_status{duplex="",interface="port9",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_status{duplex="full",interface="port1",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1
	fortiswitch_port_status{duplex="full",interface="port1",name="S108EP4N00000000",vdom="root",vlan="fortiap"} 1
	fortiswitch_port_status{duplex="full",interface="port2",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1
	fortiswitch_port_status{duplex="full",interface="port3",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1
	fortiswitch_port_status{duplex="full",interface="port4",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1
	fortiswitch_port_status{duplex="full",interface="port4",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 1
	fortiswitch_port_status{duplex="full",interface="port5",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1
	fortiswitch_port_status{duplex="full",interface="port7",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 1
	fortiswitch_port_status{duplex="full",interface="port8",name="S108EN0000000000",vdom="root",vlan="vsw.fortilink"} 1
	fortiswitch_port_status{duplex="full",interface="port8",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 1
	# HELP fortiswitch_port_speed_bps Speed negotiated on the interface in bits/s
	# TYPE fortiswitch_port_speed_bps gauge
	fortiswitch_port_speed_bps{duplex="",interface="port10",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port10",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port2",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port3",name="S108EP4N00000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port5",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port6",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port6",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port7",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port9",name="S108EN0000000000",vdom="root",vlan="snf.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="",interface="port9",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 0
	fortiswitch_port_speed_bps{duplex="full",interface="port1",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port1",name="S108EP4N00000000",vdom="root",vlan="fortiap"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port2",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port3",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port4",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port4",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port5",name="S108EN0000000000",vdom="root",vlan="proxmox"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port7",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port8",name="S108EN0000000000",vdom="root",vlan="vsw.fortilink"} 1e+09
	fortiswitch_port_speed_bps{duplex="full",interface="port8",name="S108EP4N00000000",vdom="root",vlan="vsw.fortilink"} 1e+09
	# HELP fortiswitch_port_transmit_packets_total Number of packets transmitted on the interface
	# TYPE fortiswitch_port_transmit_packets_total counter
	fortiswitch_port_transmit_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 1.784539e+07
	fortiswitch_port_transmit_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 1.9953237e+07
	fortiswitch_port_transmit_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.585012238e+09
	fortiswitch_port_transmit_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 2.8127054e+08
	fortiswitch_port_transmit_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 2.538240956e+09
	fortiswitch_port_transmit_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 4.04470704e+09
	fortiswitch_port_transmit_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 1.315834415e+09
	fortiswitch_port_transmit_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.935104868e+09
	fortiswitch_port_transmit_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 3.771233202e+09
	fortiswitch_port_transmit_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 1.4567727e+07
	fortiswitch_port_transmit_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.7394893e+09
	fortiswitch_port_transmit_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 1.870034438e+09
	fortiswitch_port_transmit_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_transmit_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
	# HELP fortiswitch_port_receive_packets_total Number of packets received on the interface
	# TYPE fortiswitch_port_receive_packets_total counter
	fortiswitch_port_receive_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 3.2926722e+07
	fortiswitch_port_receive_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 1.7078313e+07
	fortiswitch_port_receive_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.36109511e+09
	fortiswitch_port_receive_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 1.28335158e+08
	fortiswitch_port_receive_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 3.943744942e+09
	fortiswitch_port_receive_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 3.674920338e+09
	fortiswitch_port_receive_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 7.67499331e+08
	fortiswitch_port_receive_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.739489487e+09
	fortiswitch_port_receive_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 3.00568994e+09
	fortiswitch_port_receive_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 420799
	fortiswitch_port_receive_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.935104683e+09
	fortiswitch_port_receive_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 2.207656699e+09
	fortiswitch_port_receive_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_receive_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
	#HELP fortiswitch_port_receive_broadcast_packets_total Number of broadcast packets received on the interface
        # TYPE fortiswitch_port_receive_broadcast_packets_total counter
        fortiswitch_port_receive_broadcast_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 4.713459e06
        fortiswitch_port_receive_broadcast_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 1.461722e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.430059e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 464
        fortiswitch_port_receive_broadcast_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 390496
        fortiswitch_port_receive_broadcast_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 838650
        fortiswitch_port_receive_broadcast_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 1.116203e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 3.088381e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 1.172753e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 207264
        fortiswitch_port_receive_broadcast_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 5.890675e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 5.685072e06
        fortiswitch_port_receive_broadcast_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_broadcast_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_receive_bytes_total Number of bytes received on the interface
        # TYPE fortiswitch_port_receive_bytes_total counter
        fortiswitch_port_receive_bytes_total{interface="internal",name="S108EN0000000000",vdom="root"} 1.1352224784e10
        fortiswitch_port_receive_bytes_total{interface="internal",name="S108EP4N00000000",vdom="root"} 6.407616215e09
        fortiswitch_port_receive_bytes_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.1003829024042e13
        fortiswitch_port_receive_bytes_total{interface="port1",name="S108EP4N00000000",vdom="root"} 6.348630629e10
        fortiswitch_port_receive_bytes_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port2",name="S108EN0000000000",vdom="root"} 8.919415523606e12
        fortiswitch_port_receive_bytes_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port3",name="S108EN0000000000",vdom="root"} 1.2770689222182e13
        fortiswitch_port_receive_bytes_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port4",name="S108EN0000000000",vdom="root"} 9.319869478488e12
        fortiswitch_port_receive_bytes_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.159387519622e12
        fortiswitch_port_receive_bytes_total{interface="port5",name="S108EN0000000000",vdom="root"} 1.0496453460264e13
        fortiswitch_port_receive_bytes_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port7",name="S108EP4N00000000",vdom="root"} 4.2523424e07
        fortiswitch_port_receive_bytes_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.324840613348e12
        fortiswitch_port_receive_bytes_total{interface="port8",name="S108EP4N00000000",vdom="root"} 1.612957216684e12
        fortiswitch_port_receive_bytes_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_bytes_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_receive_drops_total Number of dropped packets detected during transmission on the interface
        # TYPE fortiswitch_port_receive_drops_total counter
        fortiswitch_port_receive_drops_total{interface="internal",name="S108EN0000000000",vdom="root"} 3400
        fortiswitch_port_receive_drops_total{interface="internal",name="S108EP4N00000000",vdom="root"} 25
        fortiswitch_port_receive_drops_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.2092694e07
        fortiswitch_port_receive_drops_total{interface="port1",name="S108EP4N00000000",vdom="root"} 11
        fortiswitch_port_receive_drops_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port2",name="S108EN0000000000",vdom="root"} 1.1260376e07
        fortiswitch_port_receive_drops_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port3",name="S108EN0000000000",vdom="root"} 1.3489257e07
        fortiswitch_port_receive_drops_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port4",name="S108EN0000000000",vdom="root"} 5.530408e06
        fortiswitch_port_receive_drops_total{interface="port4",name="S108EP4N00000000",vdom="root"} 65652
        fortiswitch_port_receive_drops_total{interface="port5",name="S108EN0000000000",vdom="root"} 9.331414e06
        fortiswitch_port_receive_drops_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port7",name="S108EP4N00000000",vdom="root"} 36
        fortiswitch_port_receive_drops_total{interface="port8",name="S108EN0000000000",vdom="root"} 5.339649e06
        fortiswitch_port_receive_drops_total{interface="port8",name="S108EP4N00000000",vdom="root"} 5.849172e06
        fortiswitch_port_receive_drops_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_drops_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_receive_errors_total Number of transmission errors detected on the interface
        # TYPE fortiswitch_port_receive_errors_total counter
        fortiswitch_port_receive_errors_total{interface="internal",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="internal",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port1",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port1",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port2",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port3",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port4",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port4",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port5",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port7",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port8",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port8",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_errors_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_receive_multicast_packets_total Number of multicast packets received on the interface
        # TYPE fortiswitch_port_receive_multicast_packets_total counter
        fortiswitch_port_receive_multicast_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 2.2194933e07
        fortiswitch_port_receive_multicast_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 9.672064e06
        fortiswitch_port_receive_multicast_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 4715
        fortiswitch_port_receive_multicast_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 33113
        fortiswitch_port_receive_multicast_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 5287
        fortiswitch_port_receive_multicast_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 3235
        fortiswitch_port_receive_multicast_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 6.625242e06
        fortiswitch_port_receive_multicast_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.4073094e07
        fortiswitch_port_receive_multicast_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 2.386241e06
        fortiswitch_port_receive_multicast_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 3861
        fortiswitch_port_receive_multicast_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 8.5459295e07
        fortiswitch_port_receive_multicast_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 8.7964741e07
        fortiswitch_port_receive_multicast_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_multicast_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_receive_oversized_packets_total Number of oversized packets received on the interface
        # TYPE fortiswitch_port_receive_oversized_packets_total counter
        fortiswitch_port_receive_oversized_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_oversized_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_receive_unicast_packets_total Number of unicast packets received on the interface
        # TYPE fortiswitch_port_receive_unicast_packets_total counter
        fortiswitch_port_receive_unicast_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 6.01833e06
        fortiswitch_port_receive_unicast_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 5.944527e06
        fortiswitch_port_receive_unicast_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.359660336e09
        fortiswitch_port_receive_unicast_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 1.28301581e08
        fortiswitch_port_receive_unicast_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 3.943349159e09
        fortiswitch_port_receive_unicast_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 3.674078453e09
        fortiswitch_port_receive_unicast_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 7.59757886e08
        fortiswitch_port_receive_unicast_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.722328012e09
        fortiswitch_port_receive_unicast_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 3.002130946e09
        fortiswitch_port_receive_unicast_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 209674
        fortiswitch_port_receive_unicast_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.843754713e09
        fortiswitch_port_receive_unicast_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 2.114006886e09
        fortiswitch_port_receive_unicast_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_receive_unicast_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_broadcast_packets_total Number of broadcast packets transmitted on the interface
        # TYPE fortiswitch_port_transmit_broadcast_packets_total counter
        fortiswitch_port_transmit_broadcast_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 69913
        fortiswitch_port_transmit_broadcast_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 69879
        fortiswitch_port_transmit_broadcast_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 8.624954e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 1.296068e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 9.554441e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 8.991159e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 8.828099e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 5.890675e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 8.759829e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 850698
        fortiswitch_port_transmit_broadcast_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 3.088381e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 3.288593e06
        fortiswitch_port_transmit_broadcast_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_broadcast_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_bytes_total Number of bytes transmitted on the interface
        # TYPE fortiswitch_port_transmit_bytes_total counter
        fortiswitch_port_transmit_bytes_total{interface="internal",name="S108EN0000000000",vdom="root"} 4.844915593e09
        fortiswitch_port_transmit_bytes_total{interface="internal",name="S108EP4N00000000",vdom="root"} 5.181766502e09
        fortiswitch_port_transmit_bytes_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.0921799654249e13
        fortiswitch_port_transmit_bytes_total{interface="port1",name="S108EP4N00000000",vdom="root"} 2.89477628518e11
        fortiswitch_port_transmit_bytes_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port2",name="S108EN0000000000",vdom="root"} 5.84025450544e12
        fortiswitch_port_transmit_bytes_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port3",name="S108EN0000000000",vdom="root"} 1.2355499172422e13
        fortiswitch_port_transmit_bytes_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port4",name="S108EN0000000000",vdom="root"} 1.0343038183895e13
        fortiswitch_port_transmit_bytes_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.324840632702e12
        fortiswitch_port_transmit_bytes_total{interface="port5",name="S108EN0000000000",vdom="root"} 1.3253235010876e13
        fortiswitch_port_transmit_bytes_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port7",name="S108EP4N00000000",vdom="root"} 5.86307952e09
        fortiswitch_port_transmit_bytes_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.159387493899e12
        fortiswitch_port_transmit_bytes_total{interface="port8",name="S108EP4N00000000",vdom="root"} 1.225785748078e12
        fortiswitch_port_transmit_bytes_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_bytes_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_drops_total Number of dropped packets detected during transmission on the interface
        # TYPE fortiswitch_port_transmit_drops_total counter
        fortiswitch_port_transmit_drops_total{interface="internal",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="internal",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port1",name="S108EN0000000000",vdom="root"} 2
        fortiswitch_port_transmit_drops_total{interface="port1",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port2",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port3",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port4",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port4",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port5",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port7",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port8",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port8",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_drops_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_errors_total Number of transmission errors detected on the interface
        # TYPE fortiswitch_port_transmit_errors_total counter
        fortiswitch_port_transmit_errors_total{interface="internal",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="internal",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port1",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port1",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port2",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port3",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port4",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port4",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port5",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port7",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port8",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port8",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_errors_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_multicast_packets_total Number of multicast packets transmitted on the interface
        # TYPE fortiswitch_port_transmit_multicast_packets_total counter
        fortiswitch_port_transmit_multicast_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 1.1994229e07
        fortiswitch_port_transmit_multicast_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 1.4158987e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 8.1573697e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 1.48505e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 8.1931485e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 8.1930854e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 7.5610099e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 8.5459296e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 7.9655049e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 1.29357e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.4073093e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 1.0019549e07
        fortiswitch_port_transmit_multicast_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_multicast_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_oversized_packets_total Number of oversized packets transmitted on the interface
        # TYPE fortiswitch_port_transmit_oversized_packets_total counter
        fortiswitch_port_transmit_oversized_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_oversized_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
        # HELP fortiswitch_port_transmit_unicast_packets_total Number of unicast packets transmitted on the interface
        # TYPE fortiswitch_port_transmit_unicast_packets_total counter
        fortiswitch_port_transmit_unicast_packets_total{interface="internal",name="S108EN0000000000",vdom="root"} 5.781248e06
        fortiswitch_port_transmit_unicast_packets_total{interface="internal",name="S108EP4N00000000",vdom="root"} 5.724371e06
        fortiswitch_port_transmit_unicast_packets_total{interface="port1",name="S108EN0000000000",vdom="root"} 1.494813587e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port1",name="S108EP4N00000000",vdom="root"} 2.65123972e08
        fortiswitch_port_transmit_unicast_packets_total{interface="port10",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port10",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port2",name="S108EN0000000000",vdom="root"} 2.44675503e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port2",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port3",name="S108EN0000000000",vdom="root"} 3.953785027e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port3",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port4",name="S108EN0000000000",vdom="root"} 1.231396217e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port4",name="S108EP4N00000000",vdom="root"} 1.843754897e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port5",name="S108EN0000000000",vdom="root"} 3.682818324e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port5",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port6",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port6",name="S108EP4N00000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port7",name="S108EN0000000000",vdom="root"} 0
        fortiswitch_port_transmit_unicast_packets_total{interface="port7",name="S108EP4N00000000",vdom="root"} 781329
        fortiswitch_port_transmit_unicast_packets_total{interface="port8",name="S108EN0000000000",vdom="root"} 1.722327826e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port8",name="S108EP4N00000000",vdom="root"} 1.856726296e09
        fortiswitch_port_transmit_unicast_packets_total{interface="port9",name="S108EN0000000000",vdom="root"} 0
	fortiswitch_port_transmit_unicast_packets_total{interface="port9",name="S108EP4N00000000",vdom="root"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
