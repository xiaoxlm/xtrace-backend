package controller

import (
	"github.com/ccfos/nightingale/v6/domain/busi_group_metrics"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/util"
)

func ListBusiGroupMetrics(ctx *ctx.Context, busiGroupID uint, ibn, metricUniqueID string) error {
	agg, err := busi_group_metrics.FactoryAggBusiGroupMetrics(ctx, busiGroupID, ibn, metricUniqueID)
	if err != nil {
		return err
	}

	data, err := agg.FormData()
	if err != nil {
		return err
	}

	//data[0].Metrics

	util.LogJSON(data)
	return nil
}
