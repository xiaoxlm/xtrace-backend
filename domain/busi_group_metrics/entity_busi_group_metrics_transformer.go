package busi_group_metrics

import "github.com/ccfos/nightingale/v6/models"

// 将查询出来的指标和 panel 数据结合
type busiGroupMetricsTransformer struct {
	//expr        *busiGroupMetricsExpr
	metricsData interface{}
	panel       *models.Panel

	outputData interface{}
}

// TODO 将查询值与阈值组合
func (trans *busiGroupMetricsTransformer) transform() {}

// TODO
func (trans *busiGroupMetricsTransformer) getData() {}
