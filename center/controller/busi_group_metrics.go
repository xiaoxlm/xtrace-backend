package controller

import (
	"github.com/ccfos/nightingale/v6/domain/busi_group_metrics"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

func ListMetricsAggr(ctx *ctx.Context, search models.MetricsAggr) ([]*busi_group_metrics.MetricsAggr, error) {
	list, err := models.MetricsAggrList(ctx, search)
	if err != nil {
		return nil, err
	}

	aggrs := make([]*busi_group_metrics.MetricsAggr, 0, len(list))
	for _, item := range list {
		aggrs = append(aggrs, &busi_group_metrics.MetricsAggr{
			ID:         item.ID,
			UniqueName: item.UniqueName,
			Desc:       item.Desc,
			Category:   string(item.Category),
		})
	}

	return aggrs, nil
}

func ListBusiGroupMetrics(ctx *ctx.Context, busiGroupID uint, ibn string, uniqueName models.MetricsUniqueName) ([]*busi_group_metrics.MetricsWithThresholds, error) {

	metricsAggr, err := models.MetricsAggrGetByUniqueName(ctx, uniqueName)
	if err != nil {
		return nil, err
	}

	// TODO child
	agg, err := busi_group_metrics.FactoryAggBusiGroupMetrics(ctx, busiGroupID, ibn, metricsAggr.MetricUniqueIDs[0])
	if err != nil {
		return nil, err
	}

	data, err := agg.ListMetrics(metricsAggr.Desc)
	if err != nil {
		return nil, err
	}

	return data, nil
}
