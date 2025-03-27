package controller

import (
	"github.com/ccfos/nightingale/v6/domain/traffic_topo"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

func ListTrafficData(ctx *ctx.Context, ibn string) ([]*traffic_topo.NodeTrafficData, error) {
	agg, err := traffic_topo.FactoryAgg(ctx, ibn)
	if err != nil {
		return nil, err
	}

	return agg.ListTrafficData(ctx)
}
