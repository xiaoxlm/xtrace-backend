package router

import (
	"strconv"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/ccfos/nightingale/v6/center/controller"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type MyParams struct {
	//BusiGroupID    uint   `uri:"id" binding:"required"`
	IBN        string                   `form:"ibn" binding:"required"`        // query params
	MetricType models.MetricsUniqueName `form:"metricType" binding:"required"` // query params
}

func (rt *Router) listBusiGroupMetrics(ctx *gin.Context) {
	id := ctx.Param("id")
	busiGroupID, err := strconv.Atoi(id)
	ginx.Dangerous(err)
	var params = MyParams{}

	err = ctx.ShouldBind(&params)
	ginx.Dangerous(err)

	data, err := controller.ListBusiGroupMetrics(rt.Ctx, uint(busiGroupID), params.IBN, params.MetricType)
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(data, nil)

}
