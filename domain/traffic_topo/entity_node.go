package traffic_topo

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

type EntityNode struct {
	hostIP   string
	typ      string
	traffics []*VoTraffic // 流向父亲节点的流量
}

func newEntityNode(ibn string, node *models.TrafficTopo) *EntityNode {
	if node.Type == models.TrafficTopoType_Spine {
		return &EntityNode{
			hostIP:   node.IP,
			typ:      node.Type.String(),
			traffics: []*VoTraffic{newVoTraffic(ibn, node.IP, node.Labels, nil, node.InPromql, node.OutPromql)},
		}
	}

	var traffics []*VoTraffic
	for _, conn := range node.Connects {
		traffics = append(traffics, newVoTraffic(ibn, node.IP, node.Labels, &conn, node.InPromql, node.OutPromql))
	}

	return &EntityNode{
		hostIP:   node.IP,
		typ:      node.Type.String(),
		traffics: traffics,
	}
}

func (entity *EntityNode) export(ctx *ctx.Context, startTime, endTime int64) (*NodeTrafficData, error) {
	err := entity.formToParentTraffics(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return convert2NodeTrafficData(entity), nil
}

func (entity *EntityNode) formToParentTraffics(ctx *ctx.Context, startTime, endTime int64) error {
	for _, tra := range entity.traffics {
		if err := tra.completePromql(); err != nil {
			return err
		}

		if err := tra.getMetricsData(ctx, startTime, endTime); err != nil {
			return err
		}
	}

	return nil
}

func convert2NodeTrafficData(node *EntityNode) *NodeTrafficData {
	var traffics []*trafficDataModel
	for _, traffic := range node.traffics {
		traffics = append(traffics, convert2TrafficDataModel(traffic))
	}

	return &NodeTrafficData{
		HostIP:   node.hostIP,
		Type:     node.typ,
		Traffics: traffics,
	}
}
