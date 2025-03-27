package traffic_topo

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

func FactoryAgg(ctx *ctx.Context, ibn string) (*Agg, error) {
	nodesList, err := models.TrafficTopoListAll(ctx)
	if err != nil {
		return nil, err
	}

	var nodes []*EntityNode
	for _, node := range nodesList {
		nodes = append(nodes, newEntityNode(ibn, node))
	}

	return newAgg(nodes), nil
}
