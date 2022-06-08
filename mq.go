package ipmq

import (
	"crypto/rand"
	"fmt"
	"sync"
)

type MQ interface {
	Register(c Consumer) (CancelFunc, error)
	Push(Msg) error
}
type Msg interface{}

type CancelFunc func() error
type Consumer func(m Msg) error

type mq struct {
	mu sync.RWMutex
	cm map[string]Consumer
}

func key() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (q *mq) Register(c Consumer) (CancelFunc, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	k := key()

	q.cm[k] = c

	cfn := func() error {
		q.mu.Lock()
		defer q.mu.Unlock()
		delete(q.cm, k)
		return nil
	}

	return cfn, nil
}

func (q *mq) Push(m Msg) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, c := range q.cm {
		go func(c Consumer) {
			c(m)
		}(c)
	}
	return nil
}

func New() MQ {
	return &mq{
		cm: make(map[string]Consumer),
	}
}
