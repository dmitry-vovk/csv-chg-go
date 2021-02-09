package worker

import (
	"sync"
	"time"

	"github.com/dmitry-vovk/csv-chg-go/api"
)

type APIClient interface {
	GetItem(uuid string) (*api.Item, error)
	PostAlert(uuid string) error
}

type Worker struct {
	client   APIClient            // API client instance
	interval time.Duration        // Delay between requests cycles
	uuids    map[compact]struct{} // List of UUIDs
	deleteC  chan compact         // UUIDs to delete
	wg       sync.WaitGroup       // Used to track request completion for graceful shutdown
	doneC    chan struct{}        // Closed when requested to shut down
	stoppedC chan struct{}        // Closed when shutdown has completed
	limitC   chan struct{}        // Limits number of parallel requests
}

const defaultWorkers = 1

// New returns an instance of Worker
func New(client APIClient) *Worker {
	return &Worker{
		client:   client,
		uuids:    make(map[compact]struct{}),
		deleteC:  make(chan compact),
		doneC:    make(chan struct{}),
		stoppedC: make(chan struct{}),
		limitC:   make(chan struct{}, defaultWorkers),
	}
}

// WithWorkersCount sets the limit of parallel API requests
func (w *Worker) WithWorkersCount(n int) *Worker {
	w.limitC = make(chan struct{}, n)
	return w
}

// WithInterval sets the interval between series of requests
func (w *Worker) WithInterval(interval time.Duration) *Worker {
	w.interval = interval
	return w
}

// Shutdown initiates worker stop and blocks until it finishes
func (w *Worker) Shutdown() {
	close(w.doneC)
	<-w.stoppedC
}
