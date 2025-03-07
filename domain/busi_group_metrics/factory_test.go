package busi_group_metrics

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

func Test_factoryMetricsMappingEntity(t *testing.T) {
	got, err := factoryMetricsMappingEntity(tmpCtx, 1, "算网A", "cpu_util")
	if err != nil {
		t.Fatal(err)
	}

	if err = got.parseToMetricsData(); err != nil {
		t.Fatal(err)
	}

	if err = got.setMetricsDataColor(); err != nil {
		t.Fatal(err)
	}

}
