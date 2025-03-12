package models

import (
	"github.com/ccfos/nightingale/v6/conf"
	"github.com/ccfos/nightingale/v6/pkg/util"
	"github.com/ccfos/nightingale/v6/storage"
	"gorm.io/gorm"
	"path/filepath"
	"testing"
)

var gormDB *gorm.DB

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

	gormDB, err = storage.New(config.DB)

	if err != nil {
		panic(err)
	}
}

func TestMigrate(t *testing.T) {
	err := gormDB.AutoMigrate(&MetricsAggr{})
	if err != nil {
		t.Fatal(err)
	}
}
