package job

import (
	"context"
	"sync"
)

var (
	activeJobs = make(map[string]context.CancelFunc)
	mu         sync.Mutex
)

func AddJob(id string, cancel context.CancelFunc) {
	mu.Lock()
	defer mu.Unlock()
	activeJobs[id] = cancel
}

func CancelJob(id string) bool {
	mu.Lock()
	defer mu.Unlock()
	cancel, exists := activeJobs[id]
	if exists {
		cancel()
		delete(activeJobs, id)
		return true
	}
	return false
}
