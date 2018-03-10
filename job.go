package jobprocessor

type TaskCreator interface {
	createTask() func(interface{})
}

type Job struct {
	TaskCreator TaskCreator
	JobData     interface{}
}

func (j Job) process() {
	task := j.TaskCreator.createTask()
	task(j.JobData)
}
