package busi_group_metrics

// TODO
func factoryBusiGroupMetrics(metricUniqueID, busiGroupID string) (*busiGroupMetricsExpr, error) {
	// 0. 准备表的内容
	// 1. 根据 busi_group_id 获取节点列表
	// 2. 根据metricUniqueID，获取监控表达式
	// 2.1 解析表达式
	// 2.2 metrics_mapping表要增加映射 panel 的字段
	// 3. 拿到表达式涉及到的阈值
	// 4. 查询值与阈值组合
	return newBusiGroupMetricsExpr("", "", nil)
}

// TODO
func factoryBusiGroupMetricsTransformer(expr *busiGroupMetricsExpr) (*busiGroupMetricsTransformer, error) {

}

// TODO AggBusiGroupMetrics
func FactoryAggBusiGroupMetrics(metricUniqueID, busiGroupID string) (*AggBusiGroupMetrics, error) {
	exprVO, err := factoryBusiGroupMetrics(metricUniqueID, busiGroupID)
	// 然后用表达式查询数据
}
