package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemSensorInfo(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/sensor-info", "testdata/system-sensor-info.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemSensorInfo, c, r) {
		t.Errorf("probeSystemSensorInfo() returned non-success")
	}

	em := `
	# HELP fortigate_sensor_fan_rpm Sensor fan rotation speed in RPM
	# TYPE fortigate_sensor_fan_rpm gauge
	fortigate_sensor_fan_rpm{name="FAN1"} 2900
	fortigate_sensor_fan_rpm{name="FAN2"} 2400
	fortigate_sensor_fan_rpm{name="FAN3"} 3000
	fortigate_sensor_fan_rpm{name="FAN4"} 2500
	fortigate_sensor_fan_rpm{name="FAN5"} 2900
	fortigate_sensor_fan_rpm{name="FAN6"} 2600
	fortigate_sensor_fan_rpm{name="PS1 Fan 1"} 4096
	fortigate_sensor_fan_rpm{name="PS2 Fan 1"} 4224
	# HELP fortigate_sensor_temperature_celsius Sensor temperature in degree celsius
	# TYPE fortigate_sensor_temperature_celsius gauge
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 0"} 40
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 1"} 42
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 2"} 42
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 3"} 41
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 4"} 43
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 5"} 41
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 6"} 44
	fortigate_sensor_temperature_celsius{name="CPU 0 Core 7"} 43
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 0"} 41
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 1"} 42
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 2"} 42
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 3"} 41
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 4"} 43
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 5"} 41
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 6"} 44
	fortigate_sensor_temperature_celsius{name="CPU 1 Core 7"} 43
	fortigate_sensor_temperature_celsius{name="DTS CPU0"} 47
	fortigate_sensor_temperature_celsius{name="DTS CPU1"} 49
	fortigate_sensor_temperature_celsius{name="PS1 Temp"} 25
	fortigate_sensor_temperature_celsius{name="PS2 Temp"} 25
	fortigate_sensor_temperature_celsius{name="TD1"} 31
	fortigate_sensor_temperature_celsius{name="TD2"} 38
	fortigate_sensor_temperature_celsius{name="TD3"} 27
	fortigate_sensor_temperature_celsius{name="TD4"} 30
	fortigate_sensor_temperature_celsius{name="TS1"} 31
	fortigate_sensor_temperature_celsius{name="TS2"} 31
	fortigate_sensor_temperature_celsius{name="TS3"} 32
	fortigate_sensor_temperature_celsius{name="TS4"} 32
	fortigate_sensor_temperature_celsius{name="TS5"} 31
	# HELP fortigate_sensor_voltage_volts Sensor voltage in volts
	# TYPE fortigate_sensor_voltage_volts gauge
	fortigate_sensor_voltage_volts{name="+12V"} 12.077
	fortigate_sensor_voltage_volts{name="+3.3VSB"} 3.264
	fortigate_sensor_voltage_volts{name="+3.3VSB_SMC"} 3.264
	fortigate_sensor_voltage_volts{name="3VDD"} 3.264
	fortigate_sensor_voltage_volts{name="CPU0 PVCCIN"} 1.792
	fortigate_sensor_voltage_volts{name="CPU1 PVCCIN"} 1.792
	fortigate_sensor_voltage_volts{name="MAC_1.025V"} 1.027
	fortigate_sensor_voltage_volts{name="MAC_AVS 1V"} 0.99
	fortigate_sensor_voltage_volts{name="P1V05_PCH"} 1.008
	fortigate_sensor_voltage_volts{name="P3V3_AUX"} 3.3126
	fortigate_sensor_voltage_volts{name="PS1 VIN"} 224
	fortigate_sensor_voltage_volts{name="PS1 VOUT_12V"} 12.032
	fortigate_sensor_voltage_volts{name="PS2 VIN"} 226
	fortigate_sensor_voltage_volts{name="PS2 VOUT_12V"} 12.032
	fortigate_sensor_voltage_volts{name="PVCCIO"} 1.04
	fortigate_sensor_voltage_volts{name="PVDDQ AB"} 1.2
	fortigate_sensor_voltage_volts{name="PVDDQ EF"} 1.2
	fortigate_sensor_voltage_volts{name="PVTT AB"} 0.592
	fortigate_sensor_voltage_volts{name="PVTT CD"} 0.592
	fortigate_sensor_voltage_volts{name="PVTT GH"} 0.592
	fortigate_sensor_voltage_volts{name="VCC1.15V"} 1.1581
	fortigate_sensor_voltage_volts{name="VCC2.5V"} 2.5169
	fortigate_sensor_voltage_volts{name="VCC3V3"} 3.3126
	fortigate_sensor_voltage_volts{name="VCC5V"} 4.999
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
