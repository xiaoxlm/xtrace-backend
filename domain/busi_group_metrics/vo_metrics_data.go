package busi_group_metrics

import (
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type metricsFromExpr []metricsInfo // 一个表达式得到的数据

type metricsInfo struct {
	metric map[string]string
	values []metricsValues // 时序数值
}

type metricsValues struct {
	value     string
	timestamp int64
}

func promCommonModelValue2MetricsData(promValues []model.Value) ([]metricsFromExpr, error) {
	var ret []metricsFromExpr

	for _, result := range promValues {
		mData, err := parseModelValue2metricsData(result)
		if err != nil {
			return nil, err
		}

		ret = append(ret, mData)
	}

	return ret, nil
}

func parseModelValue2metricsData(commonModelValue model.Value) (metricsFromExpr, error) {
	var ret metricsFromExpr
	switch commonModelValue.Type() {
	case model.ValScalar:
		logrus.Warnf("need to parse 'Scalar' type value")
	case model.ValVector:
		logrus.Warnf("need to parse 'Vector' type value")
	case model.ValMatrix:
		matrix := commonModelValue.(model.Matrix)
		for _, sample := range matrix {
			var values []metricsValues
			for _, value := range sample.Values {
				values = append(values, metricsValues{
					value:     value.Value.String(),
					timestamp: value.Timestamp.Unix(),
				})
			}

			var m = make(map[string]string)
			for k, v := range sample.Metric {
				m[string(k)] = string(v)
			}

			ret = append(ret, metricsInfo{
				metric: m,
				values: values,
			})
		}

	case model.ValString:
		logrus.Warnf("need to parse 'String' type value")
	default:
		return nil, fmt.Errorf("unknown metric type: %s", commonModelValue.Type())
	}

	return ret, nil
}
