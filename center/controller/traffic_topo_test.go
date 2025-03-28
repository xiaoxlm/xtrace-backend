package controller

import (
	"github.com/ccfos/nightingale/v6/pkg/util"
	"testing"
	"time"
)

func TestListTrafficData(t *testing.T) {
	ibn := "算网A"

	tim := time.Now()
	start := tim.Unix()
	end := tim.Unix()
	dataList, err := ListTrafficData(tmpCtx, ibn, start, end)
	if err != nil {
		t.Fatal(err)
	}

	util.LogJSON(dataList)
}
