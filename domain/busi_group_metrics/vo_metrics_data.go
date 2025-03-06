package busi_group_metrics

import (
	"fmt"

	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type MetricsFromExpr []MetricsInfo // 一个表达式得到的数据

type MetricsInfo struct {
	Metric map[string]string `json:"metric"`
	Values []MetricsValues   `json:"values"` // 时序数值
}

type MetricsValues struct {
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

func promCommonModelValue2MetricsData(promValues []model.Value) ([]MetricsFromExpr, error) {
	var ret []MetricsFromExpr

	for _, result := range promValues {
		mData, err := parseModelValue2metricsData(result)
		if err != nil {
			return nil, err
		}

		ret = append(ret, mData)
	}

	return ret, nil
}

func parseModelValue2metricsData(commonModelValue model.Value) (MetricsFromExpr, error) {
	var ret MetricsFromExpr
	switch commonModelValue.Type() {
	case model.ValScalar:
		logrus.Warnf("need to parse 'Scalar' type value")
	case model.ValVector:
		logrus.Warnf("need to parse 'Vector' type value")
	case model.ValMatrix:
		matrix := commonModelValue.(model.Matrix)
		for _, sample := range matrix {
			var values []MetricsValues
			for _, value := range sample.Values {
				values = append(values, MetricsValues{
					Value:     value.Value.String(),
					Timestamp: value.Timestamp.Unix(),
				})
			}

			var m = make(map[string]string)
			for k, v := range sample.Metric {
				m[string(k)] = string(v)
			}

			ret = append(ret, MetricsInfo{
				Metric: m,
				Values: values,
			})
		}

	case model.ValString:
		logrus.Warnf("need to parse 'String' type value")
	default:
		return nil, fmt.Errorf("unknown metric type: %s", commonModelValue.Type())
	}

	return ret, nil
}
