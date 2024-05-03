package lp

import (
	"context"
	"errors"
	"sync"
	"time"
)

type ILockPool interface {
	Acquire(ctx context.Context, lockName string) error
	Release(ctx context.Context) error
}

type Mutex struct {
	timestamp int64

	mu sync.Mutex
}

func NewMutex() *Mutex {
	return &Mutex{
		timestamp: time.Now().Unix(),
	}
}

type LockPool struct {
	pool        map[string]*Mutex
	waitChannel chan bool
	mu          sync.Mutex
}

func NewLockPool() *LockPool {
	return &LockPool{
		pool:        make(map[string]*Mutex),
		waitChannel: make(chan bool),
	}
}

func (lp *LockPool) Acquire(ctx context.Context, lockName string) error {
	c1, cancel1 := context.WithTimeout(ctx, 10*time.Second)

	go func(ctx context.Context) {
		lp.mu.Lock()

		lp.waitChannel <- true

		cancel1()
	}(c1)

	select {
	case <-lp.waitChannel:
		break
	case <-c1.Done():
		return errors.New("wait timeout to acquire lock pool")
	}

	if lp.pool[lockName] == nil {
		lp.pool[lockName] = &Mutex{}
	}

	lockWaitChannel := make(chan bool)

	c2, cancel2 := context.WithTimeout(ctx, 10*time.Second)

	go func(ctx context.Context) {
		lp.pool[lockName].mu.Lock()

		lockWaitChannel <- true

		cancel2()
	}(c2)

	select {
	case <-lockWaitChannel:
		lp.pool[lockName].timestamp = time.Now().Unix()
	case <-c2.Done():
		lp.mu.Unlock()
		return errors.New("wait timeout to acquire lock instance from pool")
	}

	lp.mu.Unlock()

	return nil
}

func (lp *LockPool) Release(ctx context.Context, lockName string) error {
	c1, cancel1 := context.WithTimeout(ctx, 10*time.Second)

	go func(ctx context.Context) {
		lp.mu.Lock()

		lp.waitChannel <- true

		cancel1()
	}(c1)

	select {
	case <-lp.waitChannel:
		break
	case <-c1.Done():
		return errors.New("wait timeout to acquire lock pool")
	}

	if _, ok := lp.pool[lockName]; !ok {
		lp.mu.Unlock()
		return errors.New("lock instance not found")
	}

	if lp.pool[lockName].mu.TryLock() {
		lp.mu.Unlock()
		lp.pool[lockName].mu.Unlock()
		return errors.New("mutex is already unlocked")
	}

	lp.mu.Unlock()
	lp.pool[lockName].mu.Unlock()

	return nil
}
