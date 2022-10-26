package channels

import (
	"context"
	"fmt"
	"sync"
)

type Result struct {
	JobID       int
	Value       interface{}
	Err         error
	Description string
}

type Task interface {
	Execute(int) Result
}

type Job struct {
	ID      int
	ExecFun func(int) Result
}

func (t *Job) Execute(workerID int) Result {
	//v := fmt.Sprintf("Worker %d execute this task %d", workerID, t.ID)
	//time.Sleep(1 * time.Second)
	result := t.ExecFun(workerID)
	return result
}

type WorkerPool struct {
	MaxWorkers   int
	JobStream    <-chan Job
	ResultStream chan Result
	Done         chan struct{}
}

func New(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		MaxWorkers:   maxWorkers,
		JobStream:    make(<-chan Job),
		ResultStream: make(chan Result),
		Done:         make(chan struct{}),
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	// Fan out worker goroutines
	for i := 0; i < wp.MaxWorkers; i++ {
		wg.Add(1)
		go func(workerID int, ctx context.Context, wg *sync.WaitGroup, jobStream <-chan Job, resultStream chan<- Result) {
			defer wg.Done()
			for {
				select {
				case task, ok := <-jobStream:
					if !ok {
						// Jobs channel is closed
						return
					}
					// execute the job and gather results
					wp.ResultStream <- task.Execute(workerID)
				case <-ctx.Done():
					value := fmt.Sprintf("cancelled worker %d. Error detail: %v\n", workerID, ctx.Err())
					resultStream <- Result{
						Value: value,
						Err:   ctx.Err(),
					}
					return
				}
			}
		}(i, ctx, &wg, wp.JobStream, wp.ResultStream)
	}

	// Fan-in results
	wg.Wait()
	close(wp.ResultStream)
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.ResultStream
}

//func (wp *WorkerPool) AddJobs(ctx context.Context, jobs ...Job) {
//	for _, job := range jobs {
//		select {
//		case <-ctx.Done():
//			fmt.Println("Stop feeding jobs")
//			return
//		case wp.JobStream <- job:
//		}
//	}
//	close(wp.JobStream)
//}

func (wp *WorkerPool) Assign(jobs <-chan Job) {
	wp.JobStream = jobs
}
