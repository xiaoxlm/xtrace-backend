package busi_group_metrics

type AggBusiGroupMetrics struct {
	root *busiGroupMetricsTransformer
}

func (agg *AggBusiGroupMetrics) FormData() ([]*BusiGroupMetrics, error) {
	thresholds, err := agg.root.listData()
	if err != nil {
		return nil, err
	}

	var ret []*BusiGroupMetrics
	for _, threshold := range thresholds {
		ret = append(ret, metricsWithThresholdsToBusiGroupMetrics(threshold))
	}

	return ret, nil
}

type BusiGroupMetrics struct {
	MetricUniqueID string                 `json:"metricUniqueID"`
	HostIP         string                 `json:"hostIP"`
	Metrics        BusiGroupMetricsValues `json:"metrics"`
	Color          string                 `json:"color"`
}

type BusiGroupMetricsValues struct {
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

func metricsWithThresholdsToBusiGroupMetrics(thresholds *metricsWithThresholds) *BusiGroupMetrics {
	return &BusiGroupMetrics{
		MetricUniqueID: thresholds.metricUniqueID,
		HostIP:         thresholds.hostIP,
		Metrics: BusiGroupMetricsValues{
			Value:     thresholds.metrics.value,
			Timestamp: thresholds.metrics.timestamp,
		},
		Color: thresholds.color,
	}
}
