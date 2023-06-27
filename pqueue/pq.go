package pqueue

import (
	"fmt"
	"github.com/Lmineor/lzJob/store"
	"sync"
	"time"
)

// 针对这些问题，我们就可以用优先级队列来解决。
// 我们按照任务设定的执行时间，将这些任务存储在优先级队列中，队
// 列首部（也就是小顶堆的堆顶）存储的是最先执行的任务。
// 这样，定时器就不需要每隔 1 秒就扫描一遍任务列表了。
// 它拿队首任务的执行时间点，与当前时间点相减，得到一个时间间隔 T。
// 这个时间间隔 T 就是，从当前时间开始，需要等待多久，才会有第一个任务需要被执行。
// 这样，定时器就可以设定在 T 秒之后，再来执行任务。从当前时间点到（T-1）秒这段时间里，
// 定时器都不需要做任何事情。当 T 秒时间过去之后，定时器取优先级队列中队首的任务执行。
// 然后再计算新的队首任务的执行时间点与当前时间点的差值，把这个值作为定时器执行下一个任务需要等待的时间。
// 这样，定时器既不用间隔 1 秒就轮询一次，也不用遍历整个任务列表，性能也就提高了。
type PriorityQueue struct {
	mLock      sync.Mutex                   // 互斥锁，queues和priorities并发操作时使用
	queues     map[time.Time]chan *TimeTask // 优先级队列map
	pushChan   chan *TimeTask               // 推送任务管道
	priorities []time.Time                  // 记录时间的切片（时间从小到大排列）
}

type TimeTask struct {
	*store.Task
}

func (t TimeTask) CallBack() {
	fmt.Println("callback called")
	fmt.Printf("exec time is %d", t.ExecTime.Nanosecond())

}
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		queues:   make(map[time.Time]chan *TimeTask), // 初始化优先级队列map
		pushChan: make(chan *TimeTask, 100),
	}
	// 监听pushChan
	go pq.listenPushChan()
	return pq
}

func (pq *PriorityQueue) listenPushChan() {
	for {
		select {
		case taskEle := <-pq.pushChan:
			priority := taskEle.ExecTime
			pq.mLock.Lock()
			if v, ok := pq.queues[priority]; ok {
				pq.mLock.Unlock()
				// 之前推送过相同优先级的任务
				// 将推送的任务塞到对应优先级的队列中
				v <- taskEle
				continue
			}

			// 如果这是一个新的优先级，则需要插入优先级切片，并且新建一个优先级的queue
			// 通过二分法寻找新优先级的切片插入位置
			index := pq.getNewPriorityInsertIndex(priority, 0, len(pq.priorities)-1)

			// index右侧元素均需要向后移动一个单位
			pq.moveNextPriorities(index, priority)

			// 创建一个新优先级队列
			pq.queues[priority] = make(chan *TimeTask, 10000)

			// 将任务塞到新的优先级队列中
			pq.queues[priority] <- taskEle
			pq.mLock.Unlock()
		}
	}
}

func (pq *PriorityQueue) Push(t *TimeTask) {
	pq.pushChan <- t
}

// 通过二分法寻找新优先级的切片插入位置
func (pq *PriorityQueue) getNewPriorityInsertIndex(priority time.Time, leftIndex, rightIndex int) (index int) {
	if len(pq.priorities) == 0 {
		// 如果当前优先级切片没有元素，则插入的index就是0
		return 0
	}

	length := rightIndex - leftIndex
	if pq.priorities[leftIndex].After(priority) {
		// 如果当前切片中最小的时间都超过了插入的时间，则插入位置应该是最左边
		return leftIndex
	}

	if pq.priorities[rightIndex].Before(priority) {
		// 如果当前切片中最大的时间都没超过插入的时间，则插入位置应该是最右边
		return rightIndex + 1
	}

	if length == 1 && pq.priorities[leftIndex].Before(priority) && pq.priorities[rightIndex].After(priority) {
		// 如果插入的优先级刚好在仅有的两个优先级之间，则中间的位置就是插入位置
		return leftIndex + 1
	}

	middleVal := pq.priorities[leftIndex+length/2]

	// 这里用二分法递归的方式，一直寻找正确的插入位置
	if priority.Before(middleVal) {
		return pq.getNewPriorityInsertIndex(priority, leftIndex, leftIndex+length/2)
	} else {
		return pq.getNewPriorityInsertIndex(priority, leftIndex+length/2, rightIndex)
	}
}

// index右侧元素均需要向后移动一个单位
func (pq *PriorityQueue) moveNextPriorities(index int, priority time.Time) {
	pq.priorities = append(pq.priorities, time.Now())
	copy(pq.priorities[index+1:], pq.priorities[index:])

	pq.priorities[index] = priority
}

// 消费者轮询获取最高优先级的任务
func (pq *PriorityQueue) Consume() {
	for {
		task := pq.Pop()
		if task == nil {
			// 未获取到任务，则继续轮询
			continue
		}

		// 获取到了任务，就执行任务
		task.CallBack()
	}
}

// 取出最高优先级队列中的一个任务
func (pq *PriorityQueue) Pop() *TimeTask {
	pq.mLock.Lock()
	defer pq.mLock.Unlock()
	for i := 0; i <= len(pq.priorities)-1; i++ {
		//for i := len(pq.priorities) - 1; i >= 0; i-- {
		if len(pq.queues[pq.priorities[i]]) == 0 {
			// 如果当前优先级的队列没有任务，则看低一级优先级的队列中有没有任务
			continue
		}

		// 如果当前优先级的队列里有任务，则取出一个任务。
		return <-pq.queues[pq.priorities[i]]
	}

	// 如果所有队列都没有任务，则返回null
	return nil
}