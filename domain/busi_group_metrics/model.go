package busi_group_metrics

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/lie-flat-planet/httputil"
)

type MetricsWithThresholds struct {
	MetricUniqueID models.MetricUniqueID    `json:"metricUniqueID"`
	MetricDesc     string                   `json:"metricDesc"`
	HostIP         string                   `json:"hostIP"`
	Metrics        httputil.MetricsValues   `json:"metrics"`
	Child          []*MetricsWithThresholds `json:"child,omitempty"` // TODO
}
