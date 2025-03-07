package busi_group_metrics

import "github.com/ccfos/nightingale/v6/pkg/prom"

type MetricsData struct {
	MetricUniqueID string `json:"metricUniqueID"`
	Labels map[string]string `json:"labels"`
	Desc string `json:"desc"`
	Category string `json:"category"`
	BoardPayloadID uint `json:"boardPayloadID"`
	PanelID uint `json:"panelID"`
	Data prom.MetricsFromExpr `json:"data"`
}
