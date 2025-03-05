package busi_group_metrics

type AggBusiGroupMetrics struct {
	root *busiGroupMetricsTransformer
}

func (agg *AggBusiGroupMetrics) FormData() (*MetricsWithThresholds, error) {
	return agg.root.getData()
}
