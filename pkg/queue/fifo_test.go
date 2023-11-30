package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sergeizaitcev/gophermart/pkg/queue"
)

func TestFIFO(t *testing.T) {
	var fifo queue.FIFO[int]

	waitCh := make(chan struct{})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		for i := 0; i < 10; i++ {
			v, err := fifo.Dequeue(ctx)
			if assert.NoError(t, err) {
				assert.Equal(t, i+1, v)
			}
		}

		_, err := fifo.Dequeue(ctx)
		assert.Error(t, err)

		close(waitCh)
	}()

	for i := 0; i < 10; i++ {
		_ = fifo.Enqueue(context.Background(), i+1)
	}

	<-waitCh

	assert.Empty(t, fifo.Size())
}
