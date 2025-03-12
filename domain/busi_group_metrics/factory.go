package busi_group_metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/common/model"

	"github.com/ccfos/nightingale/v6/center/service/prometheus"
	"github.com/ccfos/nightingale/v6/models"
	prometheus2 "github.com/ccfos/nightingale/v6/models/prometheus"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

//func FactoryAggBusiGroupMetrics(ctx *ctx.Context, busiGroupID uint, ibn, metricUniqueID string) (*AggBusiGroupMetrics, error) {
//	var (
//		err            error
//		targetIdents   []string                   // 获取节点列表
//		metricsMapping = &models.MetricsMapping{} // 根据metricUniqueID，获取监控表达式
//	)
//	{
//		targetIdents, err = models.TargetGroupIdsGetByGroupID(ctx, busiGroupID)
//		if err != nil {
//			return nil, err
//		}
//
//		metricsMapping, err = models.MetricsMappingGetByMetricUniqueID(ctx, metricUniqueID)
//		if err != nil {
//			return nil, err
//		}
//
//	}
//
//	// 根据表达式获取指标数据
//	var metricsData prom.MetricsFromExpr
//	{
//		// 解析表达式
//		exprVO, err := newBusiGroupMetricsExpr(metricsMapping.Expression, ibn, targetIdents)
//		if err != nil {
//			return nil, err
//		}
//
//		promAddr, err := prometheus.GetPrometheusSource(ctx)
//		if err != nil {
//			return nil, err
//		}
//
//		// 获取指标数据
//		parsedExpr := exprVO.getParsedExpr()
//		commonValues, err := prometheus2.NewPrometheus(promAddr).BatchQueryRange(ctx.Ctx, []prometheus2.QueryFormItem{
//			{
//				Start: time.Now().Unix(),
//				End:   time.Now().Unix(),
//				Step:  15,
//				Query: parsedExpr,
//			},
//		})
//		if err != nil {
//			return nil, err
//		}
//		if len(commonValues) < 1 {
//			return nil, fmt.Errorf("empty metrics data from BatchQueryRange")
//		}
//
//		// 转换数据结构
//		tmpMetricsData, err := prom.PromCommonModelValue(commonValues)
//		if err != nil {
//			return nil, err
//		}
//		// 因为只有一条表达式
//		metricsData = tmpMetricsData[0]
//	}
//
//	// 获取metrics_mapping表要增加映射 panel 的字段
//	panelContent, err := models.GetPanelContent(ctx, metricsMapping.BoardPayloadID, metricsMapping.PanelID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 查询值与阈值组合
//	return &AggBusiGroupMetrics{
//		root: &busiGroupMetricsTransformer{
//			metricUniqueID: metricUniqueID,
//			metricsData:    metricsData,
//			panel:          panelContent,
//		},
//	}, nil
//}

func FactoryAggBusiGroupMetrics(ctx *ctx.Context, busiGroupID uint, ibn string, metricUniqueID models.MetricUniqueID) (*AggBusiGroupMetrics, error) {
	avg, err := factoryMetricsMappingEntity(ctx, busiGroupID, ibn, metricUniqueID)
	if err != nil {
		return nil, err
	}

	tree := factoryEntityMetricTree(ctx, avg, nil)

	return &AggBusiGroupMetrics{
		root: tree,
	}, nil
}

func factoryMetricsMappingEntity(ctx *ctx.Context, busiGroupID uint, ibn string, metricUniqueID models.MetricUniqueID) (*metricsMappingEntity, error) {
	var (
		err            error
		metricsMapping = &models.MetricsMapping{} // 根据metricUniqueID，获取监控
		exprVO         = &busiGroupMetricsExpr{}
	)
	{
		targetIdents, tErr := models.TargetGroupIdsGetByGroupID(ctx, busiGroupID)
		if tErr != nil {
			return nil, err
		}

		metricsMapping, err = models.MetricsMappingGetByMetricUniqueID(ctx, metricUniqueID)
		if err != nil {
			return nil, err
		}

		// 解析表达式
		exprVO, err = newBusiGroupMetricsExpr(metricsMapping.Expression, ibn, targetIdents)
		if err != nil {
			return nil, err
		}
	}

	var (
		panel           = &models.Panel{}
		metricsFromProm model.Value
	)
	{
		promAddr, err := prometheus.GetPrometheusSource(ctx)
		if err != nil {
			return nil, err
		}

		// 获取指标数据
		parsedExpr := exprVO.getParsedExpr()
		commonValues, err := prometheus2.NewPrometheus(promAddr).BatchQueryRange(ctx.Ctx, []prometheus2.QueryFormItem{
			{
				Start: time.Now().Unix(),
				End:   time.Now().Unix(),
				Step:  15,
				Query: parsedExpr,
			},
		})
		if len(commonValues) < 1 {
			return nil, fmt.Errorf("commonValues is empty")
		}
		metricsFromProm = commonValues[0]

		panel, err = models.GetPanelContent(ctx, metricsMapping.BoardPayloadID, metricsMapping.PanelID)
		if err != nil {
			return nil, err
		}
	}

	return newMetricsMappingEntity(metricUniqueID, metricsMapping.LabelsToStringMap(), metricsMapping.Desc, metricsMapping.Category, panel, metricsFromProm)
}

func factoryEntityMetricTree(ctx *ctx.Context, avg *metricsMappingEntity, child *metricsMappingEntity) *entityMetricTreeEntity {
	return &entityMetricTreeEntity{
		avg: avg,
		//children: child,
	}
}
