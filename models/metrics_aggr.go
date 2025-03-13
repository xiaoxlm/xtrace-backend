package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"

	"gorm.io/datatypes"

	"gorm.io/gorm"
)

// 用于 xtrace
type MetricsAggr struct {
	ID              uint                                `json:"id" gorm:"primarykey"`
	UniqueName      string                              `json:"uniqueName" gorm:"unique;size:255"`
	Desc            string                              `json:"desc" gorm:"size:255"`
	MetricUniqueIDs datatypes.JSONSlice[MetricUniqueID] `json:"metricUniqueIDs"`
	Category        MetricsCategory                     `json:"category"` // 类别
	CreatedAt       DateTime                            `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       DateTime                            `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
	gorm.DeletedAt
}

func (m *MetricsAggr) BeforeSave(tx *gorm.DB) error {
	if len(m.MetricUniqueIDs) > 2 {
		return fmt.Errorf("metric unique ids length cannot be greater than 2")
	}
	return nil
}

func MetricsAggrGetByUniqueName(ctx *ctx.Context, uniqueName MetricsUniqueName) (*MetricsAggr, error) {
	var ma = MetricsAggr{}
	err := DB(ctx).Where("unique_name = ?", uniqueName).First(&ma).Error
	return &ma, err
}

func MetricsAggrList(ctx *ctx.Context, search MetricsAggr) ([]*MetricsAggr, error) {
	db := DB(ctx).Where("id > ?", 0)

	if search.Category != "" {
		db = db.Where("category = ?", search.Category)
	}

	if search.Desc != "" {
		db = db.Where("`desc` like ?", "%"+search.Desc+"%")
	}

	var list []*MetricsAggr
	err := db.Find(&list).Error
	return list, err
}

type MetricsUniqueName string

const (
	MetricsUniqueName_Gpu_Util     MetricsUniqueName = "gpu_util"
	MetricsUniqueName_Gpu_Mem_Util MetricsUniqueName = "gpu_mem_util"
	MetricsUniqueName_Gpu_Temp     MetricsUniqueName = "gpu_temp"
	MetricsUniqueName_Gpu_Power    MetricsUniqueName = "gpu_power"
)

type DateTime struct {
	time.Time
}

func (t DateTime) GormDataType() string {
	return "datetime"
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t *DateTime) Scan(value interface{}) error {
	if v, ok := value.(time.Time); ok {
		*t = DateTime{Time: v}
		return nil
	}
	return fmt.Errorf("failed to scan time value: %v", value)
}

func (t DateTime) Value() (driver.Value, error) {
	return t.Time, nil
}

type MetricsCategory string

const (
	MetricsCategory_Gpu MetricsCategory = "gpu"
)
