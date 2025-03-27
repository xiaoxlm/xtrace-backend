package controller

import (
	"github.com/ccfos/nightingale/v6/pkg/util"
	"testing"
)

func TestListTrafficData(t *testing.T) {
	ibn := "算网A"

	dataList, err := ListTrafficData(tmpCtx, ibn)
	if err != nil {
		t.Fatal(err)
	}

	util.LogJSON(dataList)
}
