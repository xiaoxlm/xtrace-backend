package busi_group_metrics

type AggBusiGroupMetrics struct {
	root *busiGroupMetricsTransformer
}

// TODO
func (agg *AggBusiGroupMetrics) FormData() {
	agg.root.transform()
	agg.root.getData()
}
