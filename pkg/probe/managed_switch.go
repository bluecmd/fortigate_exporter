package probe

import (
	"log"
	"strconv"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeManagedSwitch(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		managedSwitchInfo = prometheus.NewDesc(
			"fortigate_managed_switch_info",
			"Infos about a managed switch",
			[]string{"vdom", "switch_name", "os_version", "serial", "state", "status"}, nil,
		)
		managedSwitchMaxPoeBudget = prometheus.NewDesc(
			"fortigate_managed_switch_max_poe_budget_watt",
			"Max poe budget watt",
			[]string{"vdom", "switch_name"}, nil,
		)
		portInfo = prometheus.NewDesc(
			"fortigate_managed_switch_port_info",
			"Infos about a switch port",
			[]string{"vdom", "switch_name", "port", "vlan", "duplex", "status", "poe_status", "poe_capable"}, nil,
		)

		portStatus = prometheus.NewDesc(
			"fortigate_managed_switch_port_status",
			"Port status up=1 down=0",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portPower = prometheus.NewDesc(
			"fortigate_managed_switch_port_power_watt",
			"Port power in watt",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portPowerStatus = prometheus.NewDesc(
			"fortigate_managed_switch_port_power_status",
			"Port power status",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		// Port stats
		portRxBytes = prometheus.NewDesc(
			"fortigate_managed_switch_rx_bytes_total",
			"Total number of received bytes",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portRxPackets = prometheus.NewDesc(
			"fortigate_managed_switch_rx_packets_total",
			"Total number of received packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portRxUcast = prometheus.NewDesc(
			"fortigate_managed_switch_rx_ucast_packets_total",
			"Total number of received unicast packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portRxMcast = prometheus.NewDesc(
			"fortigate_managed_switch_rx_mcast_packets_total",
			"Total number of received multicast packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portRxBcast = prometheus.NewDesc(
			"fortigate_managed_switch_rx_bcast_packets_total",
			"Total number of received broadcast packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portRxErrors = prometheus.NewDesc(
			"fortigate_managed_switch_rx_errors_total",
			"Total number of received errors",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portRxDrops = prometheus.NewDesc(
			"fortigate_managed_switch_rx_drops_total",
			"Total number of received drops",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portRxOversize = prometheus.NewDesc(
			"fortigate_managed_switch_rx_oversize_total",
			"Total number of received oversize",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portTxBytes = prometheus.NewDesc(
			"fortigate_managed_switch_tx_bytes_total",
			"Total number of transmitted bytes",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portTxPackets = prometheus.NewDesc(
			"fortigate_managed_switch_tx_packets_total",
			"Total number of transmitted packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portTxUcast = prometheus.NewDesc(
			"fortigate_managed_switch_tx_ucast_packets_total",
			"Total number of transmitted unicast packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portTxMcast = prometheus.NewDesc(
			"fortigate_managed_switch_tx_mcast_packets_total",
			"Total number of transmitted multicast packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portTxBcast = prometheus.NewDesc(
			"fortigate_managed_switch_tx_bcast_packets_total",
			"Total number of transmitted broadcast packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portTxErrors = prometheus.NewDesc(
			"fortigate_managed_switch_tx_errors_total",
			"Total number of transmitted errors",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portTxDrops = prometheus.NewDesc(
			"fortigate_managed_switch_tx_drops_total",
			"Total number of transmitted drops",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portTxOversize = prometheus.NewDesc(
			"fortigate_managed_switch_tx_oversize_total",
			"Total number of transmitted oversize",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		portUndersize = prometheus.NewDesc(
			"fortigate_managed_switch_under_size_total",
			"Total number of under size",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portFragments = prometheus.NewDesc(
			"fortigate_managed_switch_fragments_total",
			"Total number of fragments",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portJabbers = prometheus.NewDesc(
			"fortigate_managed_switch_jabbers_total",
			"Total number of jabbers",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portCollisions = prometheus.NewDesc(
			"fortigate_managed_switch_collisions_total",
			"Total number of collisions",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portCrcAlignments = prometheus.NewDesc(
			"fortigate_managed_switch_crc_alignments_total",
			"Total number of crc alignments",
			[]string{"vdom", "switch_name", "port"}, nil,
		)
		portL3Packets = prometheus.NewDesc(
			"fortigate_managed_switch_l3_packets_total",
			"Total number of l3 packets",
			[]string{"vdom", "switch_name", "port"}, nil,
		)

		// Port

	)

	type Port struct {
		Interface   string  `json:"interface"`
		Status      string  `json:"status"`
		Duplex      string  `json:"duplex"`
		Speed       float64 `json:"speed"`
		PortPower   float64 `json:"port_power"`
		PowerStatus float64 `json:"power_status"`
		Vlan        string  `json:"vlan"`
		PoeCapable  bool    `json:"poe_capable"`
		PoeStatus   string  `json:"poe_status"`
	}
	type PortStat struct {
		BytesRx       float64 `json:"rx-bytes"`
		BytesTx       float64 `json:"tx-bytes"`
		PacketsRx     float64 `json:"rx-packets"`
		PacketsTx     float64 `json:"tx-packets"`
		ErrorsRx      float64 `json:"rx-errors"`
		ErrorsTx      float64 `json:"tx-errors"`
		DroppedRx     float64 `json:"rx-drops"`
		DroppedTx     float64 `json:"tx-drops"`
		UcastRx       float64 `json:"rx-ucast"`
		UcastTx       float64 `json:"tx-ucast"`
		McastRx       float64 `json:"rx-mcast"`
		McastTx       float64 `json:"tx-mcast"`
		BcastRx       float64 `json:"rx-bcast"`
		BcastTx       float64 `json:"tx-bcast"`
		OversizeRx    float64 `json:"rx-oversize"`
		OversizeTx    float64 `json:"tx-oversize"`
		Collisions    float64 `json:"collisions"`
		CrcAlignments float64 `json:"crc-alignments"`
		L3Packets     float64 `json:"l3packets"`
		Fragments     float64 `json:"fragments"`
		Undersize     float64 `json:"undersize"`
		Jabbers       float64 `json:"jabbers"`
	}

	type Results struct {
		Name           string              `json:"name"`
		VDOM           string              `json:"vdom"`
		Serial         string              `json:"serial"`
		OSVersion      string              `json:"os_version"`
		State          string              `json:"state"`
		Status         string              `json:"status"`
		ConnectingFrom string              `json:"connecting_from"`
		JoinTimeRaw    float64             `json:"join_time_raw"`
		MaxPoeBudget   float64             `json:"max_poe_budget"`
		Ports          []Port              `json:"ports"`
		PortStats      map[string]PortStat `json:"port_stats"`
	}

	type managedResponse []struct {
		Results []Results `json:"results"`
	}

	// Consider implementing pagination to remove this limit of 1000 entries
	var response managedResponse
	if err := c.Get("api/v2/monitor/switch-controller/managed-switch/status", "vdom=*&start=0&poe=true&port_stats=true&transceiver=true&count=1000", &response); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	var m []prometheus.Metric
	for _, rs := range response {
		for _, result := range rs.Results {
			m = append(m, prometheus.MustNewConstMetric(managedSwitchInfo, prometheus.CounterValue, 1, result.VDOM, result.Name, result.OSVersion, result.Serial, result.State, result.Status))
			m = append(m, prometheus.MustNewConstMetric(managedSwitchMaxPoeBudget, prometheus.CounterValue, result.MaxPoeBudget, result.VDOM, result.Name))
			for _, port := range result.Ports {
				if port.Status == "up" {
					m = append(m, prometheus.MustNewConstMetric(portStatus, prometheus.GaugeValue, 1, result.VDOM, result.Name, port.Interface))
				} else {
					m = append(m, prometheus.MustNewConstMetric(portStatus, prometheus.GaugeValue, 0, result.VDOM, result.Name, port.Interface))
				}
				m = append(m, prometheus.MustNewConstMetric(portInfo, prometheus.GaugeValue, 1, result.VDOM, result.Name, port.Interface, port.Vlan, port.Duplex, port.Status, port.PoeStatus, strconv.FormatBool(port.PoeCapable)))
				m = append(m, prometheus.MustNewConstMetric(portPower, prometheus.GaugeValue, port.PortPower, result.VDOM, result.Name, port.Interface))
				m = append(m, prometheus.MustNewConstMetric(portPowerStatus, prometheus.GaugeValue, port.PowerStatus, result.VDOM, result.Name, port.Interface))

			}
			for portName, port := range result.PortStats {
				m = append(m, prometheus.MustNewConstMetric(portRxBytes, prometheus.CounterValue, port.BytesRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxBytes, prometheus.CounterValue, port.BytesTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxPackets, prometheus.CounterValue, port.PacketsRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxPackets, prometheus.CounterValue, port.PacketsTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxUcast, prometheus.CounterValue, port.UcastRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxUcast, prometheus.CounterValue, port.UcastTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxMcast, prometheus.CounterValue, port.McastRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxMcast, prometheus.CounterValue, port.McastTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxBcast, prometheus.CounterValue, port.BcastRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxBcast, prometheus.CounterValue, port.BcastTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxErrors, prometheus.CounterValue, port.ErrorsRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxErrors, prometheus.CounterValue, port.ErrorsTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxDrops, prometheus.CounterValue, port.DroppedRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxDrops, prometheus.CounterValue, port.DroppedTx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portRxOversize, prometheus.CounterValue, port.OversizeRx, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portTxOversize, prometheus.CounterValue, port.OversizeTx, result.VDOM, result.Name, portName))

				m = append(m, prometheus.MustNewConstMetric(portUndersize, prometheus.CounterValue, port.Undersize, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portFragments, prometheus.CounterValue, port.Fragments, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portJabbers, prometheus.CounterValue, port.Jabbers, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portCollisions, prometheus.CounterValue, port.Collisions, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portCrcAlignments, prometheus.CounterValue, port.CrcAlignments, result.VDOM, result.Name, portName))
				m = append(m, prometheus.MustNewConstMetric(portL3Packets, prometheus.CounterValue, port.L3Packets, result.VDOM, result.Name, portName))
			}
		}
	}

	return m, true
}
