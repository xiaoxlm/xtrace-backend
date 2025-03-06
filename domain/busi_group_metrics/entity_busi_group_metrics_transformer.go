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
	metricsData    MetricsFromExpr
	panel          *models.Panel

	outputData []*MetricsWithThresholds
}

func (trans *busiGroupMetricsTransformer) check() error {
	if len(trans.metricsData) < 1 {
		return fmt.Errorf("busiGroupMetricsTransformer no metrics data")
	}

	return nil
}

func (trans *busiGroupMetricsTransformer) getColor(mValue MetricsValues) (string, error) {
	var (
		color         string
		thresholdsLen = len(trans.panel.Options.Thresholds.Steps)
	)

	for i := thresholdsLen - 1; i >= 0; i-- {
		step := trans.panel.Options.Thresholds.Steps[i]
		if step.Value == nil {
			color = step.Color
			break
		}

		metricValueSTR := mValue.Value
		metricValue, err := strconv.ParseFloat(strings.TrimSpace(metricValueSTR), 64)
		if err != nil {
			return "", err
		}

		v := *step.Value

		if metricValue >= v {
			color = step.Color
			break
		}
	}

	return color, nil
}

func (trans *busiGroupMetricsTransformer) combine() error {
	if err := trans.check(); err != nil {
		return err
	}

	for _, m := range trans.metricsData {
		if len(m.Values) < 1 {
			return fmt.Errorf("trans.metricsData's values is empty")
		}

		color, err := trans.getColor(m.Values[0])
		if err != nil {
			return err
		}

		trans.outputData = append(trans.outputData, &MetricsWithThresholds{
			MetricUniqueID: trans.metricUniqueID,
			HostIP:         m.Metric["host_ip"],
			Metrics:        m.Values[0],
			Color:          color,
		})
	}

	return nil
}

func (trans *busiGroupMetricsTransformer) listData() ([]*MetricsWithThresholds, error) {
	if err := trans.combine(); err != nil {
		return nil, err
	}
	return trans.outputData, nil
}

type MetricsWithThresholds struct {
	MetricUniqueID string `json:"metricUniqueID"`
	HostIP         string `json:"hostIP"`
	Metrics        MetricsValues `json:"metrics"`
	Color          string        `json:"color"`
}

