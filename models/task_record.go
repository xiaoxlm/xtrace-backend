package models

import (
	"fmt"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/poster"
	"golang.org/x/sync/errgroup"
	"sync"
)

type TaskRecord struct {
	Id           int64  `json:"id" gorm:"primaryKey"`
	EventId      int64  `json:"event_id"`
	GroupId      int64  `json:"group_id"`
	IbexAddress  string `json:"ibex_address"`
	IbexAuthUser string `json:"ibex_auth_user"`
	IbexAuthPass string `json:"ibex_auth_pass"`
	Title        string `json:"title"`
	Account      string `json:"account"`
	Batch        int    `json:"batch"`
	Tolerance    int    `json:"tolerance"`
	Timeout      int    `json:"timeout"`
	Pause        string `json:"pause"`
	Script       string `json:"script"`
	Args         string `json:"args"`
	CreateAt     int64  `json:"create_at"`
	CreateBy     string `json:"create_by"`
	Status       string `json:"status" gorm:"-"`
}

func (r *TaskRecord) TableName() string {
	return "task_record"
}

// create task
func (r *TaskRecord) Add(ctx *ctx.Context) error {
	if !ctx.IsCenter {
		err := poster.PostByUrls(ctx, "/v1/n9e/task-record-add", r)
		return err
	}

	return Insert(ctx, r)
}

// list task, filter by group_id, create_by
func TaskRecordTotal(ctx *ctx.Context, bgids []int64, beginTime int64, createBy, query string) (int64, error) {
	session := DB(ctx).Model(&TaskRecord{}).Where("create_at > ?", beginTime)
	if len(bgids) > 0 {
		session = session.Where("group_id in (?)", bgids)
	}

	if createBy != "" {
		session = session.Where("create_by = ?", createBy)
	}

	if query != "" {
		session = session.Where("title like ?", "%"+query+"%")
	}

	return Count(session)
}

func TaskRecordGets(ctx *ctx.Context, bgids []int64, beginTime int64, createBy, query string, limit, offset int) ([]*TaskRecord, error) {
	session := DB(ctx).Where("create_at > ?", beginTime).Order("create_at desc").Limit(limit).Offset(offset)
	if len(bgids) > 0 {
		session = session.Where("group_id in (?)", bgids)
	}

	if createBy != "" {
		session = session.Where("create_by = ?", createBy)
	}

	if query != "" {
		session = session.Where("title like ?", "%"+query+"%")
	}

	var lst []*TaskRecord
	err := session.Find(&lst).Error

	if len(lst) > 0 {
		var taskRecordIds []int64
		for _, r := range lst {
			taskRecordIds = append(taskRecordIds, r.Id)
		}

		statusMap, err := TaskRecordStatus(ctx, taskRecordIds)
		if err != nil {
			return nil, err
		}

		for _, r := range lst {
			r.Status = statusMap[r.Id]
		}
	}

	return lst, err
}

type TaskHost struct {
	Id     string `json:"id"`
	Host   string `json:"host"`
	Status string `json:"status"`
}

func TaskRecordStatus(ctx *ctx.Context, taskRecordIds []int64) (map[int64] /*id*/ string /*状态*/, error) {
	wg := errgroup.Group{}
	var lock = sync.Mutex{}
	var ret = make(map[int64]string)
	for _, id := range taskRecordIds {
		taskRecordId := id

		wg.Go(func() error {
			statusList, err := getTaskRecordById(ctx, taskRecordId)
			if err != nil {
				return err
			}
			status := "success"
			for _, sta := range statusList {
				if sta.Status != status {
					status = sta.Status
					break
				}
			}
			lock.Lock()
			ret[taskRecordId] = status
			lock.Unlock()

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return ret, nil
}

func getTaskRecordById(ctx *ctx.Context, taskRecordId int64) ([]TaskHost, error) {
	tableSuffix := taskRecordId % 100
	var results []TaskHost

	err := DB(ctx).Table(fmt.Sprintf("task_host_%d", tableSuffix)).Where("id = ?", taskRecordId).Find(&results).Error

	return results, err
}

// update is_done field
func (r *TaskRecord) UpdateIsDone(ctx *ctx.Context, isDone int) error {
	return DB(ctx).Model(r).Update("is_done", isDone).Error
}
