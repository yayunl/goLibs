package main

import (
	"context"
	"fanOutFanIn/channels"
	"fmt"
	"time"
)

//type PrintTask struct {
//	channels.Job
//}
//
//func (pt *PrintTask) Execute(i int) channels.Result {
//	time.Sleep(time.Second)
//	result := pt.Job.ExecFun(i)
//	return result
//}

func jobGenerator(done <-chan struct{}, tasks ...channels.Job) <-chan channels.Job {
	dataStream := make(chan channels.Job)
	// a pipeline stage for creating jobs
	go func() {
		defer close(dataStream)
		for _, task := range tasks {
			select {
			case <-done:
				return
			case dataStream <- task:
			}
		}
	}()
	return dataStream
}

func main() {
	start := time.Now()
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.TODO())
	//defer cancel()

	// Create a batch of tasks
	tasks := []channels.Job{}
	for i := 0; i < 200; i++ {
		taskID := i
		newTask := channels.Job{
			ID: i,
			// This is the func that each worker would run
			ExecFun: func(workerID int) channels.Result {
				time.Sleep(time.Second)
				value := fmt.Sprintf("Worker %d execute the task %d.", workerID, taskID)
				return channels.Result{Value: value, Err: nil, JobID: i}
			},
		}

		tasks = append(tasks, newTask)
	}
	// Create a pool of workers
	wp := channels.New(10)
	// Assign the jobs to the workers
	wp.Assign(jobGenerator(done, tasks...))
	// Start the workers to run the assigned tasks
	go wp.Run(ctx)

	// Gather the results
	results := wp.Results()
	i := 0
	for r := range results {
		fmt.Println(r.Value)
		if i == 20 {
			cancel()
		}
		i += 1
	}
	fmt.Println("Time elapsed: ", time.Since(start))
}
