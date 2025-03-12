package models

import (
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 用于 ibn
type MetricsMapping struct {
	ID             uint              `gorm:"primarykey"`
	MetricUniqueID string            `json:"metricUniqueID" gorm:"unique"` // 告警唯一标识
	Labels         datatypes.JSONMap `json:"labels"`                       // 指标标签(key:标签名；value:标签描述)
	Expression     string            `json:"-"`                            // 表达式
	Desc           string            `json:"description"`                  // 描述
	Category       string            `json:"category"`                     // 类别
	BoardPayloadID uint              `json:"-"`                            // 监控面板id
	PanelID        string            `json:"-"`                            // 具体某个仪表图id
	CreatedAt      time.Time         `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time         `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
	gorm.DeletedAt
}

func (m MetricsMapping) LabelsToStringMap() map[string]string {
	ret := make(map[string]string)
	for k, v := range m.Labels {
		ret[k] = v.(string)
	}

	return ret
}

func (MetricsMapping) TableName() string {
	return "metrics_mapping"
}

func MetricsMappingGetByMetricUniqueID(ctx *ctx.Context, metricUniqueID MetricUniqueID) (*MetricsMapping, error) {
	var mm = MetricsMapping{}
	err := DB(ctx).Where("metric_unique_id = ?", metricUniqueID).First(&mm).Error
	return &mm, err
}

type MetricUniqueID string

const (
	MetricUniqueID_Cpu_Avg_Util     MetricUniqueID = "cpu_avg_util"
	MetricUniqueID_Mem_Util         MetricUniqueID = "mem_util"
	MetricUniqueID_Gpu_Mem_Avg_Util MetricUniqueID = "gpu_mem_avg_util"
	MetricUniqueID_Gpu_Avg_Util     MetricUniqueID = "gpu_avg_util"
	MetricUniqueID_Gpu_All_Util     MetricUniqueID = "gpu_all_util"
	MetricUniqueID_Gpu_Avg_Temp     MetricUniqueID = "gpu_avg_temp"
	MetricUniqueID_Disk_Util        MetricUniqueID = "disk_util"
	MetricUniqueID_Eth_Recv         MetricUniqueID = "eth_recv_bytes_rate"
	MetricUniqueID_Eth_Trans        MetricUniqueID = "eth_trans_bytes_rate"
	MetricUniqueID_IB_Recv          MetricUniqueID = "ib_recv_bytes_rate"
	MetricUniqueID_IB_Trans         MetricUniqueID = "ib_trans_bytes_rate"
)
