package im

import (
	"github.com/yddeng/utils"
	"github.com/yddeng/utils/queue"
	"github.com/yddeng/utils/task"
	"reflect"
)

var (
	taskPool      *task.TaskPool
	taskQueueList []*TaskQueue
	taskQueue     *TaskQueue
)

type TaskQueue struct {
	cq *queue.ChannelQueue
}

func NewTaskQueue(channelSize int) *TaskQueue {
	return &TaskQueue{cq: queue.NewChannelQueue(channelSize)}
}

func (this *TaskQueue) Push(task func()) error {
	return this.cq.PushB(task)
}

func (this *TaskQueue) Run() {
	for {
		e, open := this.cq.Pop()
		if !open {
			return
		}
		e.(func())()
	}
}

/*
 将耗时函数异步处理后，调回原线程
*/
type wrapFunc func(callback interface{}, args ...interface{}) error

func WrapFunc(oriFunc interface{}) wrapFunc {
	oriF := reflect.ValueOf(oriFunc)

	if oriF.Kind() != reflect.Func {
		panic("WrapFunc oriFunc is not a func")
	}

	return func(callback interface{}, args ...interface{}) error {
		f := func() {
			out, err := utils.CallFunc(oriFunc, args...)
			if err != nil {
				panic(err)
			}

			if len(out) > 0 {
				_ = taskQueue.Push(func() {
					utils.CallFunc(callback, out...)
				})
			} else {
				_ = taskQueue.Push(func() {
					utils.CallFunc(callback)
				})
			}
		}

		return taskPool.Submit(f)
	}
}
