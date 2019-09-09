package worker

// Job ...
type Job interface {
	Handle() error
}
