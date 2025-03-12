package busi_group_metrics

import (
	"fmt"
	"sort"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/lie-flat-planet/httputil"
	"github.com/prometheus/common/model"
)

type metricsMappingEntity struct {
	metricUniqueID  models.MetricUniqueID
	labels          map[string]string
	desc            string
	category        string
	panel           *models.Panel
	metricsFromProm model.Value

	// need to convert
	multiMetricsData bool // 一个表达式是否会有多个值
	metricsData      httputil.MetricsFromExpr
}

func newMetricsMappingEntity(metricUniqueID models.MetricUniqueID, labels map[string]string, desc string, category string, panel *models.Panel, metricsFromProm model.Value) (*metricsMappingEntity, error) {
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
	metricsData, err := httputil.ParseModelValue2MetricsData(m.metricsFromProm)
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

func (m *metricsMappingEntity) setColorByMetricsValues(mValue *httputil.MetricsValues) error {
	var (
		color = "#6C53B1" // 绿色
		steps = m.sortThresholdsStepsByValue()
	)

	for _, step := range steps {
		if step.Value == nil {
			continue
		}

		v := *step.Value

		if mValue.Value >= v {
			color = step.Color
			break
		}
	}

	mValue.Color = color
	return nil
}

func (m *metricsMappingEntity) sortThresholdsStepsByValue() []models.Step {
	thresholdsSteps := m.panel.Options.Thresholds.Steps

	sort.Slice(thresholdsSteps, func(i, j int) bool {
		if thresholdsSteps[i].Value == nil || thresholdsSteps[j].Value == nil {
			return false
		}
		return *thresholdsSteps[i].Value > *thresholdsSteps[j].Value
	})

	return thresholdsSteps
}
