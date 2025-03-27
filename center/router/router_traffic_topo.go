package router

import (
	"github.com/ccfos/nightingale/v6/center/controller"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) listTrafficData(ctx *gin.Context) {
	ibn := ctx.Query("ibn")

	data, err := controller.ListTrafficData(rt.Ctx, ibn)
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(data, nil)
}
