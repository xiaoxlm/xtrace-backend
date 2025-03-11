package busi_group_metrics

type entityMetricTreeEntity struct {
	avg   *metricsMappingEntity
	child *metricsMappingEntity
}

func (entity *entityMetricTreeEntity) listAvgData() ([]*MetricsWithThresholds, error) {
	if err := entity.avg.entry(); err != nil {
		return nil, err 
	}

	var ret []*MetricsWithThresholds
	for _, data := range entity.avg.metricsData {
		ret = append(ret, &MetricsWithThresholds{
			MetricUniqueID: entity.avg.metricUniqueID,
			MetricDesc:     entity.avg.desc,
			HostIP:         data.Metric["host_ip"],
			Metrics:        data.Values[0],
			Child:          nil,
		})
	}

	return ret, nil
}
