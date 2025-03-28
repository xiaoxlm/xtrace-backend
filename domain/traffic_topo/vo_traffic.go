package traffic_topo

import (
	"fmt"
	"github.com/ccfos/nightingale/v6/center/service/prometheus"
	"github.com/ccfos/nightingale/v6/models"
	prometheus2 "github.com/ccfos/nightingale/v6/models/prometheus"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/lie-flat-planet/httputil"
	"github.com/prometheus/common/model"
	"strings"
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

func newVoTraffic(ibn, hostIP string, labels models.LabelExpressionSlice, connect *models.Connect, inPromql, outPromql string) *VoTraffic {
	vo := &VoTraffic{
		ibn:       ibn,
		hostIP:    hostIP,
		labels:    labels,
		connect:   connect,
		inPromql:  inPromql,
		outPromql: outPromql,
	}

	return vo
}

func (vo *VoTraffic) getMetricsData(ctx *ctx.Context, startTime, endTime int64) error {
	if vo.isRoot() {
		return nil
	}

	promAddr, err := prometheus.GetPrometheusSource(ctx)
	if err != nil {
		return err
	}

	// 获取指标数据
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
	if vo.isRoot() {
		return nil
	}

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
		case PromLabelName_ExportedInstance:
			l.Value = vo.hostIP

		}

		str, err := l.ToString()
		if err != nil {
			return err
		}

		labelExpressions += str + ","
	}

	labelExpressions = strings.TrimRight(labelExpressions, ",")

	vo.inCompletePromql = fmt.Sprintf(vo.inPromql, labelExpressions)
	vo.outCompletePromql = fmt.Sprintf(vo.outPromql, labelExpressions)

	return nil
}

func (vo *VoTraffic) isRoot() bool {
	return vo.connect == nil
}

func convert2TrafficDataModel(vo *VoTraffic) *trafficDataModel {
	if vo.isRoot() {
		return &trafficDataModel{
			HostIP:         vo.hostIP,
			InMetricsData:  make(httputil.MetricsFromExpr, 0),
			OutMetricsData: make(httputil.MetricsFromExpr, 0),
		}
	}

	return &trafficDataModel{
		HostIP:         vo.hostIP,
		PortDevice:     vo.connect.SelfPortDevice,
		ParentIP:       vo.connect.ParentIP,
		InMetricsData:  vo.inMetricsData,
		OutMetricsData: vo.outMetricsData,
		InExpr:         vo.inCompletePromql,
		OutExpr:        vo.outCompletePromql,
	}
}
