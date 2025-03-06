package router

import (
	"github.com/ccfos/nightingale/v6/center/controller"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type MyParams struct {
	BusiGroupID    uint   `url:"id" binding:"required"`
	IBN            string `form:"ibn" binding:"required"`            // query params
	MetricUniqueID string `form:"metricUniqueID" binding:"required"` // query params
}

func (rt *Router) listBusiGroupMetrics(ctx *gin.Context) {
	var params = MyParams{}
	err := ctx.ShouldBindUri(&params)
	ginx.Dangerous(err)

	data, err := controller.ListBusiGroupMetrics(rt.Ctx, params.BusiGroupID, params.IBN, params.MetricUniqueID)
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(gin.H{
		"data": data,
	}, nil)

}
