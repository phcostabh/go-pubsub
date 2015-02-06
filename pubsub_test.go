package pubsub

import "testing"

func TestSubscribeWithOrdinaryFunc(t *testing.T) {
	done := make(chan int)
	ps := New()
	ps.Sub(func(i interface{}) {
		done <- i.(int)
	})
	ps.Pub(1)
	i := <-done
	if i != 1 {
		t.Fatalf("Expected %v, but %d:", 1, i)
	}
}

func TestSubscribeWithAFunThatFullfilsInterface(t *testing.T) {
	done := make(chan int)
	ps := New()
	ps.Sub(Func(func(i interface{}) {
		done <- i.(int)
	}))
	ps.Pub(1)
	i := <-done
	if i != 1 {
		t.Fatalf("Expected %v, but %d:", 1, i)
	}
}

func TestString(t *testing.T) {
	done := make(chan string)
	ps := New()
	ps.Sub(func(s interface{}) {
		done <- s.(string)
	})
	ps.Pub("hello world")
	s := <-done
	if s != "hello world" {
		t.Fatalf("Expected %v, but %d:", "hello world", s)
	}
}

type F struct {
	m string
}

func TestStruct(t *testing.T) {
	done := make(chan *F)
	ps := New()
	ps.Sub(func(f interface{}) {
		done <- f.(*F)
	})
	ps.Pub(&F{"hello world"})
	f := <-done
	if f.m != "hello world" {
		t.Fatalf("Expected %v, but %d:", "hello world", f.m)
	}
}

func TestLeave(t *testing.T) {
	t.Parallel()

	done := make(chan int, 1)
	ps := New()

	ps.Sub(func(i interface{}) {
		done <- i.(int) + 1
	})

	subscription := ps.Sub(func(i interface{}) {
		done <- i.(int) + 2
	})

	ps.Pub(1)

	i := <-done
	if i != 2 {
		t.Fatalf("Expected %v, but %d:", 2, i)
	}

	i = <-done
	if i != 3 {
		t.Fatalf("Expected %v, but %d:", 3, i)
	}

	subscription.Remove()

	ps.Pub(1)
	i = <-done
	if i != 2 {
		t.Fatalf("Expected %v, but %d:", 2, i)
	}

	if len(done) > 0 {
		t.Fatal("Expeted the channel to be empty, but it isn't")
	}

	// This is another way for checking for an empty channel.
	// select {
	// case _, ok := <-done:
	// 	if ok {
	// 		t.Fatal("Expeted the channel to be empty, but it isn't")
	// 	}
	// default:
	// 	t.Log("The channel is empty as expected")
	// }
}
