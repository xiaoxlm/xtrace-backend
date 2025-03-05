package busi_group_metrics

import (
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type metricsData struct {
	values []metricsValues // 时序数值
}

type metricsValues struct {
	value     string
	timestamp int64
}

func promCommonModelValue2MetricsData(promValues []model.Value) ([]metricsData, error) {
	var ret []metricsData
	for _, result := range promValues {
		switch result.Type() {
		case model.ValScalar:
			logrus.Warnf("need to parse 'Scalar' type value")
		case model.ValVector:
			logrus.Warnf("need to parse 'Vector' type value")
		case model.ValMatrix:
			var values []metricsValues
			matrix := result.(model.Matrix)
			for _, stream := range matrix {
				for _, value := range stream.Values {
					values = append(values, metricsValues{
						value:     value.Value.String(),
						timestamp: value.Timestamp.Unix(),
					})
				}
			}

			ret = append(ret, metricsData{
				values: values,
			})
		case model.ValString:
			logrus.Warnf("need to parse 'String' type value")
		default:
			return nil, fmt.Errorf("unknown metric type: %s", result.Type())
		}
	}

	return ret, nil
}
