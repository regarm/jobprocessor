package jobprocessor

var Done = true

type TaskCreator interface {
	createTask() func(JobData, chan bool)
}
type JobData interface{}

type Job struct {
	StatusChannel chan bool
	TaskCreator   TaskCreator
	JobData       JobData
}

func (j Job) process() {
	task := j.TaskCreator.createTask()
	task(j.JobData, j.StatusChannel)
}
