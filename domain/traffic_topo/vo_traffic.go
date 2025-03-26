package traffic_topo

import (
	"fmt"
	"github.com/ccfos/nightingale/v6/center/service/prometheus"
	"github.com/ccfos/nightingale/v6/models"
	prometheus2 "github.com/ccfos/nightingale/v6/models/prometheus"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/lie-flat-planet/httputil"
	"github.com/prometheus/common/model"
	"time"
)

type VoTraffic struct {
	ibn       string
	hostIP    string
	labels    models.LabelExpressionSlice
	connect   *models.Connect
	inPromql  string
	outPromql string

	inCompletePromql  string
	outCompletePromql string

	inMetricsData  httputil.MetricsFromExpr
	outMetricsData httputil.MetricsFromExpr
}

func (vo *VoTraffic) getMetricsData(ctx *ctx.Context) error {
	promAddr, err := prometheus.GetPrometheusSource(ctx)
	if err != nil {
		return err
	}

	// 获取指标数据
	startTime := time.Now().Unix()
	endTime := time.Now().Unix()

	inData, err := prometheus2.NewPrometheus(promAddr).QueryRange(ctx.Ctx, prometheus2.QueryFormItem{
		Start: startTime,
		End:   endTime,
		Step:  10,
		Query: vo.inCompletePromql, // in
	})
	if err != nil {
		return err
	}

	outData, err := prometheus2.NewPrometheus(promAddr).QueryRange(ctx.Ctx, prometheus2.QueryFormItem{
		Start: startTime,
		End:   endTime,
		Step:  10,
		Query: vo.outCompletePromql, // in
	})
	if err != nil {
		return err
	}

	metricsData, err := httputil.PromCommonModelValue([]model.Value{inData, outData})
	if err != nil {
		return err
	}

	vo.inMetricsData = metricsData[0]
	vo.outMetricsData = metricsData[1]

	return nil
}

func (vo *VoTraffic) completePromql() error {
	var labelExpressions = ""

	for _, l := range vo.labels {
		switch PromLabelName(l.Label) {
		case PromLabelName_IBN:
			l.Value = vo.ibn
		case PromLabelName_HostIP:
			l.Value = vo.hostIP
		case PromLabelName_Device:
			l.Value = vo.connect.SelfPortDevice
		case PromLabelName_IfName:
			l.Value = vo.connect.SelfPortDevice
		}

		str, err := l.ToString()
		if err != nil {
			return err
		}

		labelExpressions += str + ","
	}

	vo.inCompletePromql = fmt.Sprintf(vo.inPromql, labelExpressions)
	vo.outCompletePromql = fmt.Sprintf(vo.outPromql, labelExpressions)

	return nil
}
