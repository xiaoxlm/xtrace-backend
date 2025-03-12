package router

import "github.com/gin-gonic/gin"


func (rt *Router) listBusiGroupMetricsCategories(ctx *gin.Context) {
	
}


type MetricsCategoryEnum string

const (
	MetricsCategoryEnum_Gpu_Util MetricsCategoryEnum = "gpu_util"
	MetricsCategoryEnum_Gpu_Mem_Util MetricsCategoryEnum = "gpu_mem_util"
	MetricsCategoryEnum_Gpu_Temp MetricsCategoryEnum = "gpu_temp"
	MetricsCategoryEnum_Gpu_Power MetricsCategoryEnum = "gpu_power"
)


type MetricsCategoryDetail struct {
	UniqueID string `json:"uniqueID"`
	Desc string `json:"desc"`
}


