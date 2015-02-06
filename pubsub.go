/*
Package pubsub provides a portable intaface to pubsub model.
PubSub can publish/subscribe/unsubscribe messages for all.
To subscribe:

	ps := pubsub.New()
	ps.Sub(func(s interface{}) {
		fmt.Println(s.(string))
	})

To publish:

	ps.Pub("hello world")

The message are allowed to pass any types, and passing to subscribers which
can accept the type for the argument of callback.
*/
package pubsub

import (
	"fmt"
	"sync"
)

// Error ...
type Error struct {
	fn interface{}
	e  interface{}
}

func (pse *Error) String() string {
	return fmt.Sprintf("%v: %v", pse.fn, pse.e)
}

func (pse *Error) Error() string {
	return fmt.Sprint(pse.e)
}

// Subscriber ...
func (pse *Error) Subscriber() interface{} {
	return pse.fn
}

// Func is an adapter to have ordinary functions to implement the Subscriber interface.
type Func func(interface{})

// Exec calls fn(i)
func (fn Func) Exec(i interface{}) {
	fn(i)
}

// Subscription ...
type Subscription struct {
	ps    *PubSub
	index int
}

// Remove ...
func (s *Subscription) Remove() {
	s.ps.Leave(s.index)
}

// Subscriber ...
type Subscriber interface {
	Exec(arg interface{})
}

// PubSub contains channel and callbacks.
type PubSub struct {
	sync.Mutex
	ch chan interface{}
	fn []Subscriber
	e  chan error
}

// New return new PubSub intreface.
func New() *PubSub {
	ps := new(PubSub)
	ps.ch = make(chan interface{})
	ps.e = make(chan error)
	go func() {
		for v := range ps.ch {
			ps.Lock()
			for _, fn := range ps.fn {
				go func(fn Subscriber, v interface{}) {
					defer func() {
						if err := recover(); err != nil {
							ps.e <- &Error{fn, err}
						}
					}()
					fn.Exec(v)
				}(fn, v)
			}
			ps.Unlock()
		}
	}()
	return ps
}

func (ps *PubSub) Error() chan error {
	return ps.e
}

// Sub subscribe to the PubSub.
func (ps *PubSub) Sub(fn interface{}) *Subscription {
	ps.Lock()
	defer ps.Unlock()
	findex := len(ps.fn)

	// This is not beautiful, but we don't need relfection though.
	switch fn.(type) {
	case Subscriber:
		ps.fn = append(ps.fn, fn.(Subscriber))
	case func(i interface{}):
		ps.fn = append(ps.fn, Func(fn.(func(i interface{}))))
	default:
		panic(`Either this is not a function or the function doesn't fulfill the signature "func (i interface{})"`)
	}

	return &Subscription{
		ps:    ps,
		index: findex,
	}
}

// Pub publish to the PubSub.
func (ps *PubSub) Pub(v interface{}) {
	ps.ch <- v
}

// Leave unsubscribe to the PubSub.
func (ps *PubSub) Leave(findex int) {
	ps.Lock()
	defer ps.Unlock()
	result := make([]Subscriber, 0, len(ps.fn))
	last := 0
	result = append(result, ps.fn[last:findex]...)
	last = findex + 1
	ps.fn = append(result, ps.fn[last:]...)
}

// Close closes PubSub. To inspect unbsubscribing for another subscruber, you must create message structure to notify them. After publish notifycations, Close should be called.
func (ps *PubSub) Close() {
	close(ps.ch)
	ps.fn = nil
}
