package busi_group_metrics

import (
	"github.com/lie-flat-planet/httputil"
)

type MetricsWithThresholds struct {
	MetricUniqueID string                   `json:"metricUniqueID"`
	MetricDesc     string                   `json:"metricDesc"`
	HostIP         string                   `json:"hostIP"`
	Metrics        httputil.MetricsValues   `json:"metrics"`
	Child          []*MetricsWithThresholds `json:"child"` // TODO
}
