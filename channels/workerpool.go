package channels

//
//type T = interface{}
//
//type WorkerPool interface {
//	Run()                // Run a worker
//	AddTask(task func()) // Add a task to the pool
//}
//
//// Define a custom type `workerPool` and let it satisfy the `WorkerPool` interface
//type workerPool struct {
//	maxWorker      int
//	queuedTaskChan chan func()
//}
//
//func New(maxWorkers int) *workerPool {
//	wp := workerPool{maxWorker: maxWorkers,
//	}
//}
//
//// Run function implements the `WorkerPool` interface.
//func (wp *workerPool) Run() {
//	for i := 0; i < wp.maxWorker; i++ {
//		// Spawn the goroutine based on the maxWorker
//		go func(workerID int) {
//			// Each worker pulls tasks which are anonymous functions out of the channel `wp.queuedTaskChan`
//			for task := range wp.queuedTaskChan {
//				// Execute the task
//				task()
//			}
//		}(i + 1)
//	}
//}
//
//// AddTask function implements the `WorkerPool` interface.
//func (wp *workerPool) AddTask(task func()) {
//	// Add a new task to the workerPool
//	wp.queuedTaskChan <- task
//}
