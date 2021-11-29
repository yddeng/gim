package im

import (
	"errors"
	"github.com/yddeng/utils/task"
	"math/rand"
)

var (
	taskPools []*task.TaskPool
)

func InitTaskPool(num int) {
	if num < 1 {
		num = 1
	}

	taskPools = make([]*task.TaskPool, 0, num)
	for i := 0; i < num; i++ {
		pool := task.NewTaskPool(1, 2048)
		taskPools = append(taskPools, pool)
	}

}

// random
func PostTask(f func()) error {
	return PostTaskHash(rand.Int(), f)
}

// hash
func PostTaskHash(n int, f func()) error {
	if len(taskPools) == 0 {
		return errors.New("task not init. ")
	}
	idx := n % len(taskPools)
	return taskPools[idx].Submit(f)
}
