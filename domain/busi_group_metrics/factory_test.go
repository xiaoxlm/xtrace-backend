package busi_group_metrics

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ccfos/nightingale/v6/conf"
	n9eCtx "github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/storage"
)

var (
	tmpCtx *n9eCtx.Context
)

func init() {
	// 加载配置
	projectPath, err := getProjectPath()
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

func TestFactoryAggBusiGroupMetrics(t *testing.T) {
	agg, err := FactoryAggBusiGroupMetrics(tmpCtx, 1, "算网A", "cpu_util")
	if err != nil {
		t.Fatal(err)
	}
	data, err := agg.FormData()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)

}

func getProjectPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find project root (where go.mod is located)
	projectRoot := currentDir
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return "", fmt.Errorf("Could not find project root - no go.mod file found")
		}
		projectRoot = parent
	}

	return projectRoot, nil
}
