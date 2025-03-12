package controller

import (
	"github.com/ccfos/nightingale/v6/domain/busi_group_metrics"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

func ListBusiGroupMetrics(ctx *ctx.Context, busiGroupID uint, ibn string, uniqueName models.MetricsUniqueName) ([]*busi_group_metrics.MetricsWithThresholds, error) {

	aggr, err := models.MetricsAggrGetByUniqueName(ctx, uniqueName)
	if err != nil {
		return nil, err
	}

	// TODO child
	agg, err := busi_group_metrics.FactoryAggBusiGroupMetrics(ctx, busiGroupID, ibn, aggr.MetricUniqueIDs[0])
	if err != nil {
		return nil, err
	}

	data, err := agg.ListMetrics()
	if err != nil {
		return nil, err
	}

	return data, nil
}
