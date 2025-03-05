package busi_group_metrics

import (
	"github.com/ccfos/nightingale/v6/models"
	"strconv"
	"strings"
)

// 将查询出来的指标和 panel 数据结合
type busiGroupMetricsTransformer struct {
	metricUniqueID string
	metricsData    []metricsData
	panel          *models.Panel

	outputData *MetricsWithThresholds
}

func (trans *busiGroupMetricsTransformer) combine() error {
	thresholdsLen := len(trans.panel.Options.Thresholds.Steps)
	var çolor string
	for i := thresholdsLen - 1; i >= 0; i-- {
		step := trans.panel.Options.Thresholds.Steps[i]
		if step.Value == nil {
			çolor = step.Color
			break
		}

		metricValueSTR := trans.metricsData[0].values[0].value
		metricValue, err := strconv.ParseFloat(strings.TrimSpace(metricValueSTR), 64)
		if err != nil {
			return err
		}

		v := *step.Value

		if metricValue >= v {
			çolor = step.Color
			break
		}
	}

	trans.outputData = &MetricsWithThresholds{
		MetricUniqueID: trans.metricUniqueID,
		Metrics:        trans.metricsData[0].values[0],
		Color:          çolor,
	}

	return nil
}

func (trans *busiGroupMetricsTransformer) getData() (*MetricsWithThresholds, error) {
	if err := trans.combine(); err != nil {
		return nil, err
	}
	return trans.outputData, nil
}

type MetricsWithThresholds struct {
	MetricUniqueID string
	Metrics        metricsValues
	Color          string
}
