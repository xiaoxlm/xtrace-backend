package models

import (
	"fmt"
	"testing"
)

func TestTrafficTopo(t *testing.T) {
	t.Run("#find", func(t *testing.T) {
		topo := TrafficTopo{}

		if err := gormDB.Debug().Where("ip=?", "10.10.1.89").Find(&topo).Error; err != nil {
			t.Fatal(err)
		}

		if topo.Type == TrafficTopoType_Spine {
			return
		}

		mm := topo.Labels.KeyToLabelExpression()
		mm["IBN"].Value = "算网A"
		mm["host_ip"].Value = topo.IP
		mm["device"].Value = "mlx5_0"

		labelExpr, err := ConvertKeyToLabelExpressionToString(mm)
		if err != nil {
			t.Fatal(err)
		}

		inPromql := fmt.Sprintf(topo.InPromql, labelExpr)
		outPromql := fmt.Sprintf(topo.OutPromql, labelExpr)

		fmt.Println("inPromql:", inPromql)
		fmt.Println("outPromql:", outPromql)
	})

	t.Run("#create", func(t *testing.T) {
		topoes := []TrafficTopo{
			{
				IP:   "10.10.1.84",
				Type: TrafficTopoType_Node,
				Connects: ConnectSlice{
					{
						ParentIP:       "10.10.1.90",
						SelfPortDevice: "mlx5_0",
					},
				},
				InPromql:  "rate(node_infiniband_port_data_received_bytes_total{%s}[15s])",
				OutPromql: "-rate(node_infiniband_port_data_transmitted_bytes_total{%s}[15s])",
				Labels: LabelExpressionSlice{
					{
						Label:    "IBN",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "host_ip",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "device",
						Operator: "=",
						Value:    "",
					},
				},
			},
			{
				IP:   "10.10.1.85",
				Type: TrafficTopoType_Node,
				Connects: ConnectSlice{
					{
						ParentIP:       "10.10.1.90",
						SelfPortDevice: "mlx5_0",
					},
				},
				InPromql:  "rate(node_infiniband_port_data_received_bytes_total{%s}[15s])",
				OutPromql: "-rate(node_infiniband_port_data_transmitted_bytes_total{%s}[15s])",
				Labels: LabelExpressionSlice{
					{
						Label:    "IBN",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "host_ip",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "device",
						Operator: "=",
						Value:    "",
					},
				},
			},
			{
				IP:   "10.10.1.90",
				Type: TrafficTopoType_Leaf,
				Connects: ConnectSlice{
					{
						ParentIP:       "10.10.1.88",
						SelfPortDevice: "IB1/1",
					},
					{
						ParentIP:       "10.10.1.89",
						SelfPortDevice: "IB1/2",
					},
				},
				InPromql:  "rate(ifHCInOctets_total{%s}[15s])",
				OutPromql: "-rate(ifHCOutOctets_total{%s}[15s])",
				Labels: LabelExpressionSlice{
					{
						Label:    "IBN",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "exported_instance",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "ifName",
						Operator: "=",
						Value:    "",
					},
				},
			},
			{
				IP:        "10.10.1.88",
				Type:      TrafficTopoType_Spine,
				Connects:  nil,
				InPromql:  "rate(ifHCInOctets_total{%s}[15s])",
				OutPromql: "-rate(ifHCOutOctets_total{%s}[15s])",
				Labels: LabelExpressionSlice{
					{
						Label:    "IBN",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "exported_instance",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "ifName",
						Operator: "=",
						Value:    "",
					},
				},
			},
			{
				IP:        "10.10.1.89",
				Type:      TrafficTopoType_Spine,
				Connects:  nil,
				InPromql:  "rate(ifHCInOctets_total{%s}[15s])",
				OutPromql: "-rate(ifHCOutOctets_total{%s}[15s])",
				Labels: LabelExpressionSlice{
					{
						Label:    "IBN",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "exported_instance",
						Operator: "=",
						Value:    "",
					},
					{
						Label:    "ifName",
						Operator: "=",
						Value:    "",
					},
				},
			},
		}

		if err := gormDB.Create(&topoes).Error; err != nil {
			t.Fatal(err)
		}
	})
}

func TestMigrateTrafficTopo(t *testing.T) {
	err := gormDB.AutoMigrate(&TrafficTopo{})
	if err != nil {
		t.Fatal(err)
	}
}
