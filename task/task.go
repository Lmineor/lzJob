package task

import (
	"github.com/Lmineor/lzJob/context"
	"github.com/Lmineor/lzJob/pqueue"
	"github.com/Lmineor/lzJob/store"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
	"time"
)

func Trigger(ctx context.LZContext) {
	go func() {
		queryTasksIn5min(ctx)
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ticker.C:
				queryTasksIn5min(ctx)
			}
		}
	}()

	for {
		doTask(ctx)
	}
}

func queryTasksIn5min(ctx context.LZContext) {
	parsedTasks, err := GetTasksInternal(ctx, time.Duration(5))
	if err != nil {
		klog.Errorf("failed to read tasks from GetTasksInternal: %s", err.Error())
		return
	}
	for _, task := range parsedTasks {
		timeTask := pqueue.TimeTask{Task: &task}
		ctx.PQ.Push(&timeTask)

	}
}

// 这样，定时器就不需要每隔 1 秒就扫描一遍任务列表了。
// 它拿队首任务的执行时间点，与当前时间点相减，得到一个时间间隔 T。
// 这个时间间隔 T 就是，从当前时间开始，需要等待多久，才会有第一个任务需要被执行。
// 这样，定时器就可以设定在 T 秒之后，再来执行任务。从当前时间点到（T-1）秒这段时间里，
// 定时器都不需要做任何事情。当 T 秒时间过去之后，定时器取优先级队列中队首的任务执行。
// 然后再计算新的队首任务的执行时间点与当前时间点的差值，把这个值作为定时器执行下一个任务需要等待的时间。
// 这样，定时器既不用间隔 1 秒就轮询一次，也不用遍历整个任务列表，性能也就提高了。
func doTask(ctx context.LZContext) {
	currentTask := ctx.PQ.Pop()
	if currentTask == nil {
		return
	}
	timeInterval := currentTask.ExecTime.Unix() - time.Now().Unix()
	timer := time.NewTimer(time.Duration(timeInterval))
	select {
	case <-timer.C:
		go currentTask.CallBack()
	}
}

func GetTasksInternal(ctx context.LZContext, after time.Duration) (ts []store.Task, err error) {
	err = ctx.DB.Find(&ts, "exec_time BETWEEN ? AND ?", time.Now(), time.Now().Add(time.Minute*after)).Error
	return ts, err
}
func GetTasks(ctx context.LZContext, pageSize int, page int) (ts []store.Task, total int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	err = ctx.DB.Model(&store.Task{}).Count(&total).Error
	if err != nil {
		return ts, total, err
	}
	err = ctx.DB.Limit(limit).Offset(offset).Find(&ts).Error
	return ts, total, err
}

func GetTask(ctx context.LZContext, taskID string) (ts store.Task, err error) {
	err = ctx.DB.First(&ts, "id = ?", taskID).Error
	return ts, err
}

func GetTasksResult(ctx context.LZContext, pageSize int, page int) (
	tsResult []store.TaskResult, total int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	err = ctx.DB.Model(&store.TaskResult{}).Count(&total).Error
	if err != nil {
		return tsResult, total, err
	}
	err = ctx.DB.Limit(limit).Offset(offset).Find(&tsResult).Error
	return tsResult, total, err
}

func GetTaskResult(ctx context.LZContext, taskId string, pageSize int, page int) (
	tsResult []store.TaskResult, total int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	whereCond := ctx.DB.Model(&store.TaskResult{}).Where("task_id = ?", taskId)
	err = whereCond.Count(&total).Error
	if err != nil {
		return tsResult, total, err
	}
	err = whereCond.Limit(limit).Offset(offset).Find(&tsResult).Error
	return tsResult, total, err
}

func AddTask(ctx context.LZContext, t store.Task) error {
	return ctx.DB.Create(&t).Error
}

func AddTaskResult(ctx context.LZContext, taskResult store.TaskResult) error {
	return ctx.DB.Create(taskResult).Error
}

func UpdateTask(ctx context.LZContext, ID string, t store.Task) error {
	var exist store.Task
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
		if err := tx.Where("task_id = ?", ID).Delete(&store.TaskResult{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", ID).Delete(&store.Task{}).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

func DeleteTaskResult(ctx context.LZContext, ID string) error {
	return ctx.DB.Delete(&store.TaskResult{}, "id = ?", ID).Error
}
