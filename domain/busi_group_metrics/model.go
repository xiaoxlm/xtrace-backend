package busi_group_metrics

import "github.com/ccfos/nightingale/v6/pkg/prom"

type MetricsWithThresholds struct {
	MetricUniqueID string                   `json:"metricUniqueID"`
	MetricDesc     string                   `json:"metricDesc"`
	HostIP         string                   `json:"hostIP"`
	Metrics        prom.MetricsValues       `json:"metrics"`
	Child          []*MetricsWithThresholds `json:"-"` // TODO
}
