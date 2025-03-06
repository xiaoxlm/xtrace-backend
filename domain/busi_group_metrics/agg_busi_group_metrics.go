package busi_group_metrics

type AggBusiGroupMetrics struct {
	root *busiGroupMetricsTransformer
}

func (agg *AggBusiGroupMetrics) FormData() ([]*metricsWithThresholds, error) {
	return agg.root.listData()
}
