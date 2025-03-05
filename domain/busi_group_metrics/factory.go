package busi_group_metrics

import (
	"github.com/ccfos/nightingale/v6/center/service/prometheus"
	"github.com/ccfos/nightingale/v6/models"
	prometheus2 "github.com/ccfos/nightingale/v6/models/prometheus"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"time"
)

func FactoryAggBusiGroupMetrics(ctx *ctx.Context, busiGroupID, ibn, metricUniqueID string) (*AggBusiGroupMetrics, error) {
	var (
		err            error
		targetIdents   []string                   // 获取节点列表
		metricsMapping = &models.MetricsMapping{} // 根据metricUniqueID，获取监控表达式
	)
	{
		targetIdents, err = models.TargetGroupIdsGetByGroupID(ctx, busiGroupID)
		if err != nil {
			return nil, err
		}

		metricsMapping, err = models.MetricsMappingGetByMetricUniqueID(ctx, metricUniqueID)
		if err != nil {
			return nil, err
		}

	}

	// 根据表达式获取指标数据
	var metricsValue []metricsData
	{
		// 解析表达式
		exprVO, err := newBusiGroupMetricsExpr(metricsMapping.Expression, ibn, targetIdents)
		if err != nil {
			return nil, err
		}

		promAddr, err := prometheus.GetPrometheusSource(ctx)
		if err != nil {
			return nil, err
		}

		// 获取指标数据
		commonValues, err := prometheus2.NewPrometheus(promAddr).BatchQueryRange(ctx.Ctx, []prometheus2.QueryFormItem{
			{
				Start: time.Now().Unix(),
				End:   time.Now().Unix(),
				Step:  15,
				Query: exprVO.getParsedExpr(),
			},
		})
		if err != nil {
			return nil, err
		}

		// 转换数据结构
		metricsValue, err = promCommonModelValue2MetricsData(commonValues)
		if err != nil {
			return nil, err
		}
	}

	// 获取metrics_mapping表要增加映射 panel 的字段
	panelContent, err := models.GetPanelContent(ctx, metricsMapping.BoardPayloadID, metricsMapping.PanelID)
	if err != nil {
		return nil, err
	}

	// 查询值与阈值组合
	return &AggBusiGroupMetrics{
		root: &busiGroupMetricsTransformer{
			metricUniqueID: metricUniqueID,
			metricsData:    metricsValue,
			panel:          panelContent,
		},
	}, nil
}
