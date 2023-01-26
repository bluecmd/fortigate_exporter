package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"

)

func probeSwitchPortStats(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mSwitchStatus = prometheus.NewDesc(
			"fortiswitch_status",
			"Whether the switch is connected or not",
			[]string{"vdom", "name", "fgt_peer_intf_name", "connection_from", "state"}, nil,
		)
		mPortStatus = prometheus.NewDesc(
			"fortiswitch_port_status",
			"Whether the switch port is up or not",
			[]string{"vdom", "name", "interface", "vlan", "duplex"}, nil,
		)
		mPortSpeed = prometheus.NewDesc(
			"fortiswitch_port_speed_bps",
			"Speed negotiated on the interface in bits/s",
			[]string{"vdom", "name", "interface", "vlan", "duplex"}, nil,
		)
		mTxPkts = prometheus.NewDesc(
			"fortiswitch_port_transmit_packets_total",
			"Number of packets transmitted on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxB = prometheus.NewDesc(
			"fortiswitch_port_transmit_bytes_total",
			"Number of bytes transmitted on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxUcast = prometheus.NewDesc(
			"fortiswitch_port_transmit_unicast_packets_total",
			"Number of unicast packets transmitted on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxMcast = prometheus.NewDesc(
			"fortiswitch_port_transmit_multicast_packets_total",
			"Number of multicast packets transmitted on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxBcast = prometheus.NewDesc(
			"fortiswitch_port_transmit_broadcast_packets_total",
			"Number of broadcast packets transmitted on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxErr = prometheus.NewDesc(
			"fortiswitch_port_transmit_errors_total",
			"Number of transmission errors detected on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxDrops = prometheus.NewDesc(
			"fortiswitch_port_transmit_drops_total",
			"Number of dropped packets detected during transmission on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mTxOverS = prometheus.NewDesc(
			"fortiswitch_port_transmit_oversized_packets_total",
			"Number of oversized packets transmitted on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
		mRxPkts = prometheus.NewDesc(
			"fortiswitch_port_receive_packets_total",
			"Number of packets received on the interface",
			[]string{"vdom", "name", "interface"}, nil,
		)
                mRxB = prometheus.NewDesc(
                        "fortiswitch_port_receive_bytes_total",
                        "Number of bytes received on the interface",
                        []string{"vdom", "name", "interface"}, nil,
                )
                mRxUcast = prometheus.NewDesc(
                        "fortiswitch_port_receive_unicast_packets_total",
                        "Number of unicast packets received on the interface",
                        []string{"vdom", "name", "interface"}, nil,
                )
                mRxMcast = prometheus.NewDesc(
                        "fortiswitch_port_receive_multicast_packets_total",
                        "Number of multicast packets received on the interface",
                        []string{"vdom", "name", "interface"}, nil,
                )
                mRxBcast = prometheus.NewDesc(
                        "fortiswitch_port_receive_broadcast_packets_total",
                        "Number of broadcast packets received on the interface",
                        []string{"vdom", "name", "interface"}, nil,
                )
                mRxErr = prometheus.NewDesc(
                        "fortiswitch_port_receive_errors_total",
                        "Number of transmission errors detected on the interface",
                        []string{"vdom", "name", "interface"}, nil,
                )
                mRxDrops = prometheus.NewDesc(
                        "fortiswitch_port_receive_drops_total",
                        "Number of dropped packets detected during transmission on the interface",
                        []string{"vdom", "name", "interface"}, nil,
                )
                mRxOverS = prometheus.NewDesc(
                        "fortiswitch_port_receive_oversized_packets_total",
                        "Number of oversized packets received on the interface",
                        []string{"vdom", "name", "interface"}, nil,		
		)
	)

	type portStats struct {
		TxPackets float64 `json:"tx-packets"`
		TxBytes   float64 `json:"tx-bytes"`
		TxErrors  float64 `json:"tx-errors"`
		TxMcast   float64 `json:"tx-mcast"`
		TxUcast   float64 `json:"tx-ucast"`
		TxBcast  float64 `json:"tx-bcast"`
		TxDrops   float64 `json:"tx-drops"`
		TxOversize   float64 `json:"tx-oversize"`
		RxPackets float64 `json:"rx-packets"`
		RxBytes   float64 `json:"rx-bytes"`
		RxErrors  float64 `json:"rx-errors"`
                RxMcast   float64 `json:"rx-mcast"`
                RxUcast   float64 `json:"rx-ucast"`
                RxBcast  float64 `json:"rx-bcast"`
                RxDrops   float64 `json:"rx-drops"`
                RxOversize   float64 `json:"rx-oversize"`		
	}
	type portsInfo struct {
		Interface	string
		Status		string
		Duplex		string
		Speed		float64
		Vlan		string
	}
	type swResult struct {
		Name		string
		FgPeerIntfName	string `json:"fgt_peer_intf_name"`
		Status		string
		State		string
		Connection	string `json:"connecting_from"`
		VDOM		string
		Ports		[]portsInfo
		PortStats	map[string]portStats `json:"port_stats"`
	}
	
	type swResponse struct {
		//Results map[string]swResult
		Results []swResult `json:"results"`
		//VDOM    string
	}
	var r swResponse
	//var r map[string]swResponse
	//var r []swResponse

	if err := c.Get("api/v2/monitor/switch-controller/managed-switch", "port_stats=true", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
//	for _, v := range r {
		for _, swr := range r.Results {
			swStatus := 0.0
			if swr.Status == "Connected" {
				swStatus = 1.0
			}
			m = append(m, prometheus.MustNewConstMetric(mSwitchStatus, prometheus.GaugeValue, swStatus, swr.VDOM, swr.Name, swr.FgPeerIntfName, swr.Connection, swr.State))
			
			for _, pi := range swr.Ports {
				pStatus := 0.0
				if pi.Status == "up" {
					pStatus = 1.0
				}
				m = append(m, prometheus.MustNewConstMetric(mPortStatus, prometheus.GaugeValue, pStatus, swr.VDOM, swr.Name, pi.Interface, pi.Vlan, pi.Duplex))
				m = append(m, prometheus.MustNewConstMetric(mPortSpeed, prometheus.GaugeValue, pi.Speed*1000*1000, swr.VDOM, swr.Name, pi.Interface, pi.Vlan, pi.Duplex))

			}

			for port, ps := range swr.PortStats {
				m = append(m, prometheus.MustNewConstMetric(mTxPkts, prometheus.CounterValue, ps.TxPackets, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxB, prometheus.CounterValue, ps.TxBytes, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxUcast, prometheus.CounterValue, ps.TxUcast, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxBcast, prometheus.CounterValue, ps.TxBcast, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxMcast, prometheus.CounterValue, ps.TxMcast, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxErr, prometheus.CounterValue, ps.TxErrors, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxDrops, prometheus.CounterValue, ps.TxDrops, swr.VDOM, swr.Name, port))
				m = append(m, prometheus.MustNewConstMetric(mTxOverS, prometheus.CounterValue, ps.TxOversize, swr.VDOM, swr.Name, port))
                        	m = append(m, prometheus.MustNewConstMetric(mRxPkts, prometheus.CounterValue, ps.RxPackets, swr.VDOM, swr.Name, port))
                        	m = append(m, prometheus.MustNewConstMetric(mRxB, prometheus.CounterValue, ps.RxBytes, swr.VDOM, swr.Name, port))
                        	m = append(m, prometheus.MustNewConstMetric(mRxUcast, prometheus.CounterValue, ps.RxUcast, swr.VDOM, swr.Name, port))
                        	m = append(m, prometheus.MustNewConstMetric(mRxBcast, prometheus.CounterValue, ps.RxBcast, swr.VDOM, swr.Name, port))
                        	m = append(m, prometheus.MustNewConstMetric(mRxMcast, prometheus.CounterValue, ps.RxMcast, swr.VDOM, swr.Name, port))
                        	m = append(m, prometheus.MustNewConstMetric(mRxErr, prometheus.CounterValue, ps.RxErrors, swr.VDOM, swr.Name, port))
                       		m = append(m, prometheus.MustNewConstMetric(mRxDrops, prometheus.CounterValue, ps.RxDrops, swr.VDOM, swr.Name, port))
                       		m = append(m, prometheus.MustNewConstMetric(mRxOverS, prometheus.CounterValue, ps.RxOversize, swr.VDOM, swr.Name, port))				
                        }
		}
//	}
	return m, true
}
