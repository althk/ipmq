package ipmq

import (
	"sync"
	"sync/atomic"
	"testing"
)

var testCounter int32
var wg sync.WaitGroup

func testConsumer(m Msg) error {
	atomic.AddInt32(&testCounter, 1)
	wg.Done()
	return nil
}

func Test_mq_Register_Adds_ConsumerFn(t *testing.T) {
	testCounter = 0
	m := &mq{
		cm: make(map[string]Consumer),
	}
	_, err := m.Register(testConsumer)
	if err != nil {
		t.Fatalf("mq.Register error, got: %v, want: nil", err)
	}
	if len(m.cm) != 1 {
		t.Errorf("mq.Register fail, want 1 consumer enqueued, got %d", len(m.cm))
	}
}

func Test_mq_Register_Cancels_Consumer(t *testing.T) {
	testCounter = 0
	m := &mq{
		cm: make(map[string]Consumer),
	}
	cancel, err := m.Register(testConsumer)
	if err != nil {
		t.Fatalf("mq.Register error, got: %v, want: nil", err)
	}
	cancel()
	if len(m.cm) != 0 {
		t.Errorf("mq.Register fail, want 0 consumer after cancel, got %d", len(m.cm))
	}

}
func Test_mq_Push_Consumers(t *testing.T) {
	testCounter = 0
	q := New()
	q.Register(testConsumer)
	q.Register(testConsumer)
	wg.Add(2)
	q.Push("some msg")
	wg.Wait()
	if testCounter != 2 {
		t.Errorf("mq.Push fail, want testCounter: 2, got: %d", testCounter)
	}
}
