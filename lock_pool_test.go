package lp

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitiation(t *testing.T) {
	lp := NewLockPool()
	assert.NotNil(t, lp)
	assert.NotNil(t, lp.pool)
	assert.NotNil(t, lp.waitChannel)
}

func TestAcquire(t *testing.T) {
	lp := &LockPool{
		pool:        make(map[string]*Mutex),
		waitChannel: make(chan bool),
	}

	t.Run("should return success on first lock case", func(t *testing.T) {
		err := lp.Acquire(context.Background(), "some-lock")
		assert.NoError(t, err)

		mu, ok := lp.pool["some-lock"]
		assert.True(t, ok)
		assert.NotNil(t, mu)

		delete(lp.pool, "some-lock")
	})

	t.Run("should return error timeout on get locked pool", func(t *testing.T) {
		lp.mu.Lock()

		err := lp.Acquire(context.Background(), "some-lock")
		assert.Errorf(t, err, "wait timeout to acquire lock pool")

		lp.mu.Unlock()
	})

	t.Run("should return error timeout on get locked instance", func(t *testing.T) {
		lp.Acquire(context.Background(), "some-lock")

		err := lp.Acquire(context.Background(), "some-lock")

		assert.Errorf(t, err, "wait timeout to acquire lock instance from pool")

		delete(lp.pool, "some-lock")
	})
}

func TestRelease(t *testing.T) {
	lp := NewLockPool()

	t.Run("should release lock", func(t *testing.T) {
		lp.Acquire(context.Background(), "some-lock")

		err := lp.Release(context.Background(), "some-lock")
		assert.NoError(t, err)
		assert.True(t, lp.pool["some-lock"].mu.TryLock())

		lp.Release(context.Background(), "some-lock")
	})

	t.Run("should return error timeout waiting pool", func(t *testing.T) {
		lp.mu.Lock()
		err := lp.Release(context.Background(), "some-lock")

		assert.Errorf(t, err, "wait timeout to acquire lock pool")
		lp.mu.Unlock()
	})

	t.Run("should return error if lock instance is not found", func(t *testing.T) {
		err := lp.Release(context.Background(), "doesn't-exist-lock")
		assert.Errorf(t, err, "lock instance not found")
	})

	t.Run("should return error if lock instance is already unlocked", func(t *testing.T) {
		lp.Acquire(context.Background(), "some-lock")

		lp.Release(context.Background(), "some-lock")
		err := lp.Release(context.Background(), "some-lock")

		assert.Errorf(t, err, "mutex is already unlocked")
	})
}

func TestUseCase(t *testing.T) {
	lp := NewLockPool()

	i := 1
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		ctx := context.Background()
		errAcquire := lp.Acquire(ctx, "some-lock")
		assert.NoError(t, errAcquire)
		i += 4
		lp.Release(ctx, "some-lock")
		wg.Done()
	}()

	// time.Sleep(1 * time.Second)
	ctx := context.Background()
	errAcquire := lp.Acquire(ctx, "some-lock")
	assert.NoError(t, errAcquire)

	i += 3

	lp.Release(ctx, "some-lock")
	wg.Wait()
	assert.Equal(t, 8, i)
}
