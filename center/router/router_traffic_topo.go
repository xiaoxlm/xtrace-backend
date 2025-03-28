package router

import (
	"github.com/ccfos/nightingale/v6/center/controller"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"time"
)

func (rt *Router) listTrafficData(ctx *gin.Context) {
	ibn := ctx.Query("ibn")

	t := time.Now()
	start := t.Unix()
	end := t.Unix()

	data, err := controller.ListTrafficData(rt.Ctx, ibn, start, end)
	ginx.Dangerous(err)

	ginx.NewRender(ctx).Data(data, nil)
}
