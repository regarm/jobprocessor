package jobprocessor

type processorChannels struct {
	pauseChannel       chan bool
	workerCountChannel chan bool
}
type processorStatus struct {
	paused       bool
	scheduled    bool
	workersCount int
}
type processor struct {
	maxroutines       int
	jobQueue          *JobQueue
	processorStatus   *processorStatus
	processorChannels *processorChannels
}

func (prcs *processor) init() {
	if !prcs.processorStatus.scheduled {
		prcs.sched()
	}
}
func (prcs *processor) start() {
	prcs.processorChannels.pauseChannel <- false
}
func (prcs *processor) pause() {
	prcs.processorChannels.pauseChannel <- true
}

func (prcs *processor) sched() {
	invokedChannel := make(chan bool)
	go func() {
		invokedChannel <- true
		for {
			select {
			case <-prcs.processorChannels.workerCountChannel:
				prcs.processorStatus.workersCount--
				if !prcs.processorStatus.paused && prcs.processorStatus.workersCount < prcs.maxroutines {
					prcs.schedNewWoker()
				}
			case in := <-prcs.processorChannels.pauseChannel:
				prcs.processorStatus.paused = in
			default:
				if !prcs.processorStatus.paused && prcs.processorStatus.workersCount < prcs.maxroutines {
					prcs.schedNewWoker()
				}
			}
		}
	}()

	<-invokedChannel
	prcs.processorStatus.scheduled = true
}
func (prcs *processor) schedNewWoker() {
	responseChan := make(chan Job)
	prcs.jobQueue.PopChannel() <- responseChan
	job := <-responseChan
	if job.TaskCreator == nil {
		return
	}
	prcs.processorStatus.workersCount++
	go func() {
		job.process()
		prcs.processorChannels.workerCountChannel <- true
	}()
}

func NewProcessor(jq *JobQueue) *processor {
	return &processor{
		maxroutines:     1000,
		jobQueue:        jq,
		processorStatus: &processorStatus{},
		processorChannels: &processorChannels{
			pauseChannel:       make(chan bool),
			workerCountChannel: make(chan bool, 10000),
		},
	}
}
