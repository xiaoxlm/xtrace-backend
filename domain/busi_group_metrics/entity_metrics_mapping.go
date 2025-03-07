package busi_group_metrics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/prometheus/common/model"
)

type metricsMappingEntity struct {
	metricUniqueID  string
	labels          map[string]string
	desc            string
	category        string
	panel           *models.Panel
	metricsFromProm model.Value

	// need to convert
	multiMetricsData bool // 一个表达式是否会有多个值
	metricsData      prom.MetricsFromExpr
}

func newMetricsMappingEntity(metricUniqueID string, labels map[string]string, desc string, category string, panel *models.Panel, metricsFromProm model.Value) (*metricsMappingEntity, error) {
	entity := &metricsMappingEntity{
		metricUniqueID:  metricUniqueID,
		labels:          labels,
		desc:            desc,
		category:        category,
		panel:           panel,
		metricsFromProm: metricsFromProm,
	}

	if entity.metricUniqueID == "cpu_avg_util" {
		entity.multiMetricsData = false
	}

	if err := entity.check(); err != nil {
		return nil, err
	}

	return entity, nil
}

func (m *metricsMappingEntity) check() error {

	if m.metricsFromProm == nil {
		return fmt.Errorf("metricsFromProm is nil in metricsMappingEntity")
	}

	return nil
}

func (m *metricsMappingEntity) entry() error {
	if err := m.parseToMetricsData(); err != nil {
		return err
	}

	if err := m.setMetricsDataColor(); err != nil {
		return err
	}

	return nil
}

// 将prom的model.Value转换为metricsData
func (m *metricsMappingEntity) parseToMetricsData() error {
	metricsData, err := prom.ParseModelValue2MetricsData(m.metricsFromProm)
	if err != nil {
		return err
	}

	m.metricsData = metricsData

	return nil
}

func (m *metricsMappingEntity) setMetricsDataColor() error {
	for i, data := range m.metricsData {

		for j := range data.Values {
			if err := m.setColorByMetricsValues(&data.Values[j]); err != nil {
				return err
			}
		}

		m.metricsData[i] = data
	}

	return nil
}

func (m *metricsMappingEntity) setColorByMetricsValues(mValue *prom.MetricsValues) error {
	var (
		color         string
		thresholdsLen = len(m.panel.Options.Thresholds.Steps)
	)

	for i := thresholdsLen - 1; i >= 0; i-- {
		step := m.panel.Options.Thresholds.Steps[i]
		if step.Value == nil {
			color = step.Color
			break
		}

		metricValueSTR := mValue.Value
		metricValue, err := strconv.ParseFloat(strings.TrimSpace(metricValueSTR), 64)
		if err != nil {
			return err
		}

		v := *step.Value

		if metricValue >= v {
			color = step.Color
			break
		}
	}

	mValue.Color = color
	return nil
}
