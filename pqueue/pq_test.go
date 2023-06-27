package pqueue

import (
	"fmt"
	"github.com/Lmineor/lzJob/store"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	pq := NewPriorityQueue()
	rand.Seed(time.Now().Unix())

	// 我们在这里，随机生成一些优先级任务
	for i := 0; i < 100; i++ {
		time.Sleep(1)
		pq.Push(&TimeTask{&store.Task{
			ID:            uuid.UUID{},
			Type:          "cron",
			Status:        "init",
			Cron:          "dddd",
			ExecTime:      time.Now(),
			ParamBody:     "dd",
			ExtInfo:       "dd",
			TriggerMethod: "dd",
			TriggerSpi:    "dd",
		}})
	}

	// 这里会阻塞，消费者会轮询查询任务队列
	pq.Consume()
}