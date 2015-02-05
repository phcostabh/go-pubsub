package main

import (
	"fmt"
	"log"
	"time"

	"github.com/phcostabh/go-pubsub"
)

type foo struct {
	bar string
}

func main() {
	ps := pubsub.New()

	go func() {
		for err := range ps.Error() {
			if pse, ok := err.(*pubsub.Error); ok {
				log.Println(pse.String())
				// ps.Leave(pse.Subscriber())
			} else {
				log.Println(err.Error())
			}
		}
	}()

	ps.Sub(func(i interface{}) {
		fmt.Println("int subscriber: ", i.(int))
	})
	ps.Sub(func(s interface{}) {
		fmt.Println("string subscriber: ", s.(string))
	})
	ps.Sub(func(f interface{}) {
		fmt.Println("foo subscriber1: ", f.(*foo).bar)
	})
	ps.Sub(func(f interface{}) {
		fmt.Println("foo subscriber2: ", f.(*foo).bar)
	})
	var f3 func(f interface{})
	f3 = func(f interface{}) {
		fmt.Println("foo subscriber3: ", f.(*foo).bar)
	}
	ps.Sub(f3)
	ps.Sub(func(f interface{}) {
		fmt.Println("float64 subscriber: ", f.(float64))
		panic("Crash!")
	})

	ps.Pub(1)
	ps.Pub("hello")
	ps.Pub(2)
	ps.Pub(2.4)
	ps.Pub(&foo{"bar!"})
	time.Sleep(1 * time.Second)
	ps.Pub(2.4)
	ps.Pub(&foo{"bar!"})

	time.Sleep(5 * time.Second)
}
