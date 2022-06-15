## In-Process Messaging Queue
An extremely simple, in-process messaging queue that allows messages to be pushed to multiple consumers in the same process in a thread-safe way.

The API is very light, here's an example:

```go

import github.com/althk/ipmq

func main() {
	q := ipmq.New()  // get a new instance of type ipmq.MQ

	someConsumerFn := func(msg ipmq.Msg) error {
		// do something with msg
	}

	cancel, err := q.Register(someConsumerFn)
	if err != nil {
		// registration failed, do something
	}

	// if the consumer needs to unregister
	// simply call cancel
	// cancel()  // unregisters the consumer

	// push a msg to all registered consumers
	q.Push("some msg")  // calls all consumers concurrently

}

```

If there is a need for multiple 'topics', simply instantiate multiple instances of `ipmq.MQ` via `ipmq.New`.