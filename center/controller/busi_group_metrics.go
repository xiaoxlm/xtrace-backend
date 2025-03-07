package controller

import (
	"github.com/ccfos/nightingale/v6/domain/busi_group_metrics"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

func ListBusiGroupMetrics(ctx *ctx.Context, busiGroupID uint, ibn, metricUniqueID string) ([]*busi_group_metrics.MetricsWithThresholds, error) {
	agg, err := busi_group_metrics.FactoryAggBusiGroupMetrics(ctx, busiGroupID, ibn, metricUniqueID)
	if err != nil {
		return nil, err
	}

	data, err := agg.ListMetrics()
	if err != nil {
		return nil, err
	}

	return data, nil
}
