package router

import (
	"github.com/ccfos/nightingale/v6/center/controller"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"strconv"
)

type MyParams struct {
	//BusiGroupID    uint   `uri:"id" binding:"required"`
	IBN            string `form:"ibn" binding:"required"`            // query params
	MetricUniqueID string `form:"metricUniqueID" binding:"required"` // query params
}

func (rt *Router) listBusiGroupMetrics(ctx *gin.Context) {
	id := ctx.Param("id")
	busiGroupID, err := strconv.Atoi(id)
	ginx.Dangerous(err)
	var params = MyParams{}

	err = ctx.ShouldBind(&params)
	ginx.Dangerous(err)

	data, err := controller.ListBusiGroupMetrics(rt.Ctx, uint(busiGroupID), params.IBN, params.MetricUniqueID)
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(data, nil)

}
