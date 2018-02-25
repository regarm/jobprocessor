package jobprocessor

import (
	"container/list"
)

type JobQueue struct {
	l           *list.List
	pushChannel chan Job
	popChannel  chan Job
}

func pop(l *list.List) Job {
	var job Job
	el := l.Front()
	if el == nil {
		return job
	}
	return l.Remove(el).(Job)
}

func (jq JobQueue) Push(j Job) {
	jq.pushChannel <- j
}

func (jq JobQueue) Pop() Job {
	return <-jq.popChannel
}

func New() JobQueue {
	jq := JobQueue{
		l:           list.New(),
		pushChannel: make(chan Job),
		popChannel:  make(chan Job),
	}
	jq.l.Init()

	go func() {
		defer close(jq.pushChannel)
		defer close(jq.popChannel)

		for {
			select {
			case toPush := <-jq.pushChannel:
				jq.l.PushBack(toPush)
			case jq.popChannel <- pop(jq.l):
			}
		}
	}()

	return jq
}
