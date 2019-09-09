package worker

// Worker ...
type Worker struct {
	workerChannel chan chan Job
	jobChannel    chan Job
	quit          chan bool
}

// NewWorker ...
func NewWorker(workerChannel chan chan Job) Worker {
	return Worker{
		workerChannel: workerChannel,
		jobChannel:    make(chan Job),
		quit:          make(chan bool),
	}
}

// Run ...
func (w Worker) Run() {
	go func() {
		for {
			w.workerChannel <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				job.Handle()
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop ...
func (w Worker) Stop() {
	w.quit <- true
}
