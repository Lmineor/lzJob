package store

import (
	"github.com/Lmineor/lzJob/context"
	"gorm.io/gorm"
	"time"
)

func GetTasksInternal(ctx context.LZContext, after time.Duration) (ts []Task, err error) {
	err = ctx.DB.Find(&ts, "exec_time BETWEEN ? AND ?", time.Now(), time.Now().Add(time.Minute*after) ).Error
	return ts, err
}
func GetTasks(ctx context.LZContext, pageSize int, page int) (ts []Task, total int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	err = ctx.DB.Model(&Task{}).Count(&total).Error
	if err != nil {
		return ts, total, err
	}
	err = ctx.DB.Limit(limit).Offset(offset).Find(&ts).Error
	return ts, total, err
}

func GetTask(ctx context.LZContext, taskID string) (ts Task, err error) {
	err = ctx.DB.First(&ts, "id = ?", taskID).Error
	return ts, err
}

func GetTasksResult(ctx context.LZContext, pageSize int, page int) (
	tsResult []TaskResult, total int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	err = ctx.DB.Model(&TaskResult{}).Count(&total).Error
	if err != nil {
		return tsResult, total, err
	}
	err = ctx.DB.Limit(limit).Offset(offset).Find(&tsResult).Error
	return tsResult, total, err
}

func GetTaskResult(ctx context.LZContext, taskId string, pageSize int, page int) (
	tsResult []TaskResult, total int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	whereCond := ctx.DB.Model(&TaskResult{}).Where("task_id = ?", taskId)
	err = whereCond.Count(&total).Error
	if err != nil {
		return tsResult, total, err
	}
	err = whereCond.Limit(limit).Offset(offset).Find(&tsResult).Error
	return tsResult, total, err
}

func AddTask(ctx context.LZContext, t Task) error {
	return ctx.DB.Create(&t).Error
}

func AddTaskResult(ctx context.LZContext, taskResult TaskResult) error {
	return ctx.DB.Create(taskResult).Error
}

func UpdateTask(ctx context.LZContext, ID string, t Task) error {
	var exist Task
	err := ctx.DB.Where("id = ?", ID).Find(&exist).Error
	if err != nil {
		return err
	}
	exist.Status = t.Status
	exist.Cron = t.Cron
	exist.Type = t.Type
	exist.ExecTime = t.ExecTime
	exist.ParamBody = t.ParamBody
	exist.ExtInfo = t.ExtInfo
	exist.TriggerMethod = t.TriggerMethod
	exist.TriggerSpi = t.TriggerSpi
	return ctx.DB.Save(&exist).Error
}

func DeleteTask(ctx context.LZContext, ID string) error {
	db := ctx.DB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_id = ?", ID).Delete(&TaskResult{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", ID).Delete(&Task{}).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

func DeleteTaskResult(ctx context.LZContext, ID string) error {
	return ctx.DB.Delete(&TaskResult{}, "id = ?", ID).Error
}
