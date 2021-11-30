package im

import (
	"errors"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/queue"
	"math/rand"
)

var task_pool *taskPool

type taskPool struct {
	cqs []*queue.ChannelQueue
}

func newTaskPool(num int) *taskPool {
	if num < 1 {
		num = 1
	}

	cqs := make([]*queue.ChannelQueue, 0, num)
	for i := 0; i < num; i++ {
		cq := queue.NewChannelQueue(1024)
		cqs = append(cqs, cq)
	}
	return &taskPool{cqs: cqs}
}

func (this *taskPool) postTaskHash(n int, f func()) error {
	if this == nil || len(this.cqs) == 0 {
		return errors.New("task not init. ")
	}

	idx := n % len(this.cqs)
	return this.cqs[idx].PushB(f)
}

func (this *taskPool) postTask(f func()) error {
	return this.postTaskHash(rand.Int(), f)
}

func (this *taskPool) run() {
	for i, cq := range this.cqs {
		go func(i int) {
			log.Debugf("task pool queue %d run. ", i)
			for {
				e, open := cq.Pop()
				if !open {
					log.Debugf("task pool queue %d stopped. ", i)
					return
				}

				e.(func())()
			}
		}(i)
	}
}
