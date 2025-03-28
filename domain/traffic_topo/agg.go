package traffic_topo

import "github.com/ccfos/nightingale/v6/pkg/ctx"

type Agg struct {
	nodes []*EntityNode
}

func newAgg(nodes []*EntityNode) *Agg {
	return &Agg{
		nodes: nodes,
	}
}

func (agg *Agg) ListTrafficData(ctx *ctx.Context, startTime, endTime int64) ([]*NodeTrafficData, error) {
	var ret []*NodeTrafficData

	for _, node := range agg.nodes {
		data, err := node.export(ctx, startTime, endTime)
		if err != nil {
			return nil, err
		}

		ret = append(ret, data)
	}

	return ret, nil
}
