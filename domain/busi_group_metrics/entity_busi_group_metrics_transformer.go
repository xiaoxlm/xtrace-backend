package busi_group_metrics

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/ccfos/nightingale/v6/models"
)

// 将查询出来的指标和 panel 数据结合
type busiGroupMetricsTransformer struct {
	metricUniqueID string
	metricsData    []metricsData
	panel          *models.Panel

	outputData []*metricsWithThresholds
}

func (trans *busiGroupMetricsTransformer) check() error {
	if len(trans.metricsData) < 1 {
		return fmt.Errorf("busiGroupMetricsTransformer no metrics data")
	}

	return nil
}

func (trans *busiGroupMetricsTransformer) combine() error {
	if err := trans.check(); err != nil {
		return err
	}

	thresholdsLen := len(trans.panel.Options.Thresholds.Steps)
	for _, m := range trans.metricsData {
		if len(m.values) < 1 {
			return fmt.Errorf("trans.metricsData's values is empty")
		}

		var çolor string
		for i := thresholdsLen - 1; i >= 0; i-- {
			step := trans.panel.Options.Thresholds.Steps[i]
			if step.Value == nil {
				çolor = step.Color
				break
			}

			metricValueSTR := m.values[0].value
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

		trans.outputData = append(trans.outputData, &metricsWithThresholds{
			metricUniqueID: trans.metricUniqueID,
			hostIP:         m.metric["host_ip"],
			metrics:        m.values[0],
			color:          çolor,
		})
	}

	return nil
}

func (trans *busiGroupMetricsTransformer) listData() ([]*metricsWithThresholds, error) {
	if err := trans.combine(); err != nil {
		return nil, err
	}
	return trans.outputData, nil
}

type metricsWithThresholds struct {
	metricUniqueID string
	hostIP         string
	metrics        metricsValues
	color          string
}
