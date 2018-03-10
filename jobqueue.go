package jobprocessor

import (
	"container/list"
)

type jobQueueChannels struct {
	pushChannel chan Job
	popChannel  chan chan Job
}
type jobQueueStatus struct {
	scheduled bool
}
type JobQueue struct {
	l              *list.List
	channels       *jobQueueChannels
	jobProcessor   *processor
	jobQueueStatus *jobQueueStatus
}

func (jq JobQueue) PushChannel() chan Job {
	return jq.channels.pushChannel
}

func (jq JobQueue) PopChannel() chan chan Job {
	return jq.channels.popChannel
}

func (jq JobQueue) Sched() {
	if !jq.jobQueueStatus.scheduled {
		prcs := NewProcessor(&jq)
		jq.jobProcessor = prcs
		prcs.init()
	}
	jq.jobQueueStatus.scheduled = true
}

func getFor(el *list.Element) Job {
	var job Job
	if el != nil {
		job = el.Value.(Job)
	}
	return job
}

func New() JobQueue {
	jq := JobQueue{
		l: list.New(),
		channels: &jobQueueChannels{
			pushChannel: make(chan Job),
			popChannel:  make(chan chan Job),
		},
		jobQueueStatus: &jobQueueStatus{},
	}
	jq.l.Init()

	pushPopListenerStarted := make(chan bool)
	go func() {
		defer close(jq.channels.pushChannel)
		defer close(jq.channels.popChannel)
		pushPopListenerStarted <- true
		for {
			select {
			case toPush := <-jq.channels.pushChannel:
				jq.l.PushBack(toPush)
			case toPop := <-jq.channels.popChannel:
				value := jq.l.Front()
				var job Job
				if value != nil {
					job = value.Value.(Job)
					jq.l.Remove(value)
				}
				toPop <- job
			}
		}
	}()

	<-pushPopListenerStarted
	close(pushPopListenerStarted)

	return jq
}
