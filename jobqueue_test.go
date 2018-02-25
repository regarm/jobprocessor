package jobprocessor

import (
	"fmt"
	"testing"
	"time"
)

//Test for concurrent read and writes for container
//Run with go test -race
func TestRaceConditions(t *testing.T) {
	q := New()
	for i := 0; i < 1000; i++ {
		go func(i int) {
			for i := 0; i < 20; i++ {
				q.Push(Job{
					JobData: i,
				})
			}
		}(i)
	}
	time.Sleep(time.Millisecond * 1000)
	for i := 0; i < 1000; i++ {
		x := q.Pop()
		if x.JobData == nil {
			break
		}
		fmt.Println(x.JobData)
	}
}
