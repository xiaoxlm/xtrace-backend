package models

import (
	"testing"

	"gorm.io/datatypes"
)

func TestDateTime_CreateMetricsAggr(t *testing.T) {

	metricsAggr := []*MetricsAggr{
		{
			UniqueName: string(MetricsUniqueName_Gpu_Util),
			Desc:       "gpu利用率",
			MetricUniqueIDs: datatypes.JSONSlice[MetricUniqueID]{
				MetricUniqueID_Gpu_Avg_Util,
				MetricUniqueID_Gpu_All_Util,
			},
			Category: "gpu",
		},
		{
			UniqueName: string(MetricsUniqueName_Gpu_Mem_Util),
			Desc:       "gpu内存利用率", 
			MetricUniqueIDs: datatypes.JSONSlice[MetricUniqueID]{
				MetricUniqueID_Gpu_Mem_Avg_Util,
			},
			Category: "gpu",
		},
		{
			UniqueName: string(MetricsUniqueName_Gpu_Temp),
			Desc:       "gpu温度",
			MetricUniqueIDs: datatypes.JSONSlice[MetricUniqueID]{
				MetricUniqueID_Gpu_Avg_Temp,
			},
			Category: "gpu",
		},
	}

	if err := gormDB.Create(metricsAggr).Error; err != nil {
		t.Fatal(err)
	}
}
