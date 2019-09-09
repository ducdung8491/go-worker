package worker

// Config ...
type Config struct {
	WorkerNumber int
	MaxJobs      int
}

// Manager ...
type Manager struct {
	workers       []*Worker
	workerChannel chan chan Job
	jobChannel    chan Job
}

// NewManager ...
func NewManager(config Config) *Manager {
	workerChannel := make(chan chan Job, config.WorkerNumber)
	workers := []*Worker{}
	for i := 0; i < config.WorkerNumber; i++ {
		workers = append(workers, &Worker{
			workerChannel: workerChannel,
			jobChannel:    make(chan Job),
			quit:          make(chan bool),
		})
	}
	return &Manager{
		workers:       workers,
		workerChannel: workerChannel,
		jobChannel:    make(chan Job),
	}
}

// Start ...
func (m Manager) Start() {
	for _, worker := range m.workers {
		worker.Run()
	}
	go m.listen()
}

// Run ...
func (m Manager) Run(job Job) {
	m.jobChannel <- job
}

func (m Manager) listen() {
	for {
		select {
		case job := <-m.jobChannel:
			go func(job Job) {
				worker := <-m.workerChannel
				worker <- job
			}(job)
		}
	}
}

// Stop ...
func (m Manager) Stop() {
	for _, worker := range m.workers {
		worker.Stop()
	}
}
