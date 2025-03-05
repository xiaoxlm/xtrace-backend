package router

import (
	"github.com/ccfos/nightingale/v6/domain/busi_group_metrics"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type MyParams struct {
	BusiGroupID    string `url:"id" binding:"required"`
	IBN            string `form:"ibn" binding:"required"`            // query params
	MetricUniqueID string `form:"metricUniqueID" binding:"required"` // query params
}

func (rt *Router) listBusiGroupMetrics(ctx *gin.Context) {
	var params = MyParams{}
	err := ctx.ShouldBindUri(&params)
	ginx.Dangerous(err)

	agg, err := busi_group_metrics.FactoryAggBusiGroupMetrics(rt.Ctx, params.BusiGroupID, params.IBN, params.MetricUniqueID)
	ginx.Dangerous(err)

	data, err := agg.FormData()
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(gin.H{
		"data": data,
	}, nil)

}
