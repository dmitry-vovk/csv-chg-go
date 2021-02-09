package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	w := New(nil).WithWorkersCount(17).WithInterval(time.Second * 37)
	assert.Equal(t, 17, cap(w.limitC))
	assert.Equal(t, 37.0, w.interval.Seconds())
}

func TestShutdown(t *testing.T) {
	w := New(nil).WithInterval(time.Second)
	go w.Run()
	time.Sleep(10 * time.Millisecond) // Wait for goroutine to start
	assert.Eventually(t, func() bool {
		w.Shutdown()
		return true
	}, time.Millisecond*10, time.Millisecond)
}
