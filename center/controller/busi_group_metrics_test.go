package controller

import (
	"context"
	"github.com/ccfos/nightingale/v6/conf"
	n9eCtx "github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/util"
	"github.com/ccfos/nightingale/v6/storage"
	"path/filepath"
	"testing"
)

var (
	tmpCtx *n9eCtx.Context
)

func init() {
	// 加载配置
	projectPath, err := util.GetProjectPath()
	if err != nil {
		panic(err)
	}

	configDir := filepath.Join(projectPath, "etc")
	config, err := conf.InitConfig(configDir, "")
	if err != nil {
		panic(err)
	}

	db, err := storage.New(config.DB)
	if err != nil {
		panic(err)
	}

	tmpCtx = n9eCtx.NewContext(context.Background(), db, true)
}

func TestListBusiGroupMetrics(t *testing.T) {
	t.Run("cpu_util", func(t *testing.T) {
		data, err := ListBusiGroupMetrics(tmpCtx, 1, "算网A", "cpu_util")
		if err != nil {
			t.Fatal(err)
		}
		util.LogJSON(data)
	})

	t.Run("mem_util", func(t *testing.T) {
		data, err := ListBusiGroupMetrics(tmpCtx, 1, "算网A", "mem_util")
		if err != nil {
			t.Fatal(err)
		}
		util.LogJSON(data)
	})

	t.Run("gpu_util", func(t *testing.T) {
		data, err := ListBusiGroupMetrics(tmpCtx, 1, "算网A", "gpu_util")
		if err != nil {
			t.Fatal(err)
		}
		util.LogJSON(data)
	})

}
