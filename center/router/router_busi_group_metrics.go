package router

import (
	"fmt"
	"strconv"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/ccfos/nightingale/v6/center/controller"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type MetricType string

const (
	MetricType_Gpu_Util MetricType = "gpu_util"
)

func (mt MetricType) ToMetricUniqueID() (models.MetricUniqueID, error) {
	switch mt {
	case MetricType_Gpu_Util:
		return models.MetricUniqueID_Avg_Gpu_Util, nil
	}

	return "", fmt.Errorf("invalid metric type: %s", mt)

}

type MyParams struct {
	//BusiGroupID    uint   `uri:"id" binding:"required"`
	IBN        string     `form:"ibn" binding:"required"`        // query params
	MetricType MetricType `form:"metricType" binding:"required"` // query params
}

func (rt *Router) listBusiGroupMetrics(ctx *gin.Context) {
	id := ctx.Param("id")
	busiGroupID, err := strconv.Atoi(id)
	ginx.Dangerous(err)
	var params = MyParams{}

	err = ctx.ShouldBind(&params)
	ginx.Dangerous(err)

	metricUniqueID, err := params.MetricType.ToMetricUniqueID()
	ginx.Dangerous(err)

	data, err := controller.ListBusiGroupMetrics(rt.Ctx, uint(busiGroupID), params.IBN, metricUniqueID)
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(data, nil)

}
