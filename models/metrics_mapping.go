package models

import (
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

func (MetricsMapping) TableName() string {
	return "metrics_mapping"
}

func MetricsMappingGetByMetricUniqueID(ctx *ctx.Context, metricUniqueID string) (*MetricsMapping, error) {
	var mm = MetricsMapping{}
	err := DB(ctx).Where("metric_unique_id = ?", metricUniqueID).First(&mm).Error
	return &mm, err
}
