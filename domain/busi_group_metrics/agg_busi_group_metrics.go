package busi_group_metrics

type AggBusiGroupMetrics struct {
	root *entityMetricTreeEntity
}

func (agg *AggBusiGroupMetrics) ListMetrics(metricsAggrDesc string) ([]*MetricsWithThresholds, error) {
	return agg.root.listAvgData(metricsAggrDesc)
}
