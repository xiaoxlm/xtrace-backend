package busi_group_metrics

type AggBusiGroupMetrics struct {
	root *entityMetricTreeEntity
}

func (agg *AggBusiGroupMetrics) ListMetrics() ([]*MetricsWithThresholds, error) {
	return agg.root.listAvgData()
}
