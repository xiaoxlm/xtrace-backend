package tests

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TaskHost struct {
	Id     string `json:"id"`
	Host   string `json:"host"`
	Status string `json:"status"`
}

func ListStatus(taskID string) error {

	dsn := "root:uWXf87plmQGz8zMM@tcp(10.10.1.84:3306)/n9e_v6?charset=utf8mb4&parseTime=True&loc=Local"
	dbd, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	var results []TaskHost
	// Table("task_host_20").
	err = dbd.Table("task_host_20").Where("id IN ?", []int64{1920, 1820}).Find(&results).Error
	fmt.Println(results)

	return err
}
