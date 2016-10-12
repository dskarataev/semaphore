package semaphore

import (
	"time"
)

type monitorFunc func(semaphoreInfo HealthCheck, startTime time.Time, errorOccurred bool)

type monitoredSemaphore struct {
	Semaphore
	acqMonFunc monitorFunc
	relMonFunc monitorFunc
}

// NewMonitoredSemaphore creates new monitored semaphore
func NewMonitoredSemaphore(size int, acquireMonitorFunc monitorFunc, releaseMonitorFunc monitorFunc) Semaphore {
	return &monitoredSemaphore{
		Semaphore:  New(size),
		acqMonFunc: acquireMonitorFunc,
		relMonFunc: releaseMonitorFunc,
	}
}

// Acquire tries to acquire resource
func (ms monitoredSemaphore) Acquire(timeout time.Duration) error {
	start := time.Now()
	err := ms.Semaphore.Acquire(timeout)

	// call monitoring function
	ms.acqMonFunc(ms.Semaphore, start, err != nil)

	return err
}

// Release tries to release resource
func (ms monitoredSemaphore) Release() error {
	start := time.Now()
	err := ms.Semaphore.Release()

	// call monitoring function
	ms.relMonFunc(ms.Semaphore, start, err != nil)

	return err
}

