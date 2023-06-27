package lzJob

import (
	"github.com/Lmineor/lzJob/context"
	"github.com/Lmineor/lzJob/pqueue"
	"github.com/Lmineor/lzJob/store"
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
	parsedTasks, err := store.GetTasksInternal(ctx, time.Duration(5))
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
	timeInterval := currentTask.ExecTime.Unix() - time.Now().Unix()
	timer := time.NewTimer(time.Duration(timeInterval))
	select {
	case <-timer.C:
		go currentTask.CallBack()
	}
}
