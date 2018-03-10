package jobprocessor

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type taskHandler struct{}

type jobData struct {
	statusChannel chan bool
	data          interface{}
}

var sum int32 = 0

func (t taskHandler) createTask() func(interface{}) {
	return func(jd interface{}) {
		var jda jobData
		jda = jd.(jobData)
		atomic.AddInt32(&sum, jda.data.(int32))
		jda.statusChannel <- true
	}
}

//Test for concurrent read and writes for container
//Run with GORACE="halt_on_error=1" go test -race
func TestConcurrencyAndRaceConditions(t *testing.T) {
	fmt.Println("Testing for concurrency and race conditions")
	q := New()
	q.Sched()
	notifyChannel := make(chan bool)
	for i := 0; i < 1000; i++ {
		go func(j int) {
			for j := 0; j < 200; j++ {
				jd := jobData{
					statusChannel: notifyChannel,
					data:          int32(j),
				}
				j := Job{
					JobData:     jd,
					TaskCreator: taskHandler{},
				}
				q.PushChannel() <- j
			}
		}(i)
	}
	ticker := time.NewTicker(time.Second * 1)
	counter := 0
	done := false
	for {
		select {
		case <-ticker.C:
			fmt.Println(runtime.NumGoroutine())
		case <-notifyChannel:
			counter++
			if counter == 200000 {
				done = true
				break
			}
		}
		if done {
			break
		}
	}
	assert.Equal(t, 19900000, int(sum))
	assert.Equal(t, 200000, counter)
}
