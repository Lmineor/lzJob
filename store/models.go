package store

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Task struct {
	ID            uuid.UUID `json:"id" gorm:"comment:任务id"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Cron          string    `json:"cron"`
	ExecTime      time.Time `json:"exec_time"`
	ParamBody     string    `json:"param_body"`
	ExtInfo       string    `json:"ext_info"`
	TriggerMethod string    `json:"trigger_method"`
	TriggerSpi    string    `json:"trigger_spi"`
}

type TaskResult struct {
	ID         uuid.UUID `json:"id" gorm:"comment:任务结果id"`
	TaskID     uuid.UUID `json:"task_id" gorm:"comment:任务id"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	FinishTime time.Time `json:"finish_time"`
	Desc       string    `json:"desc"`
}