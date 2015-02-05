go-pubsub
=========

PubSub model message hub

Usage
-----

To subscribe:

```go
    ps := pubsub.New()
    ps.Sub(func(i interface{} {
        fmt.Println("int subscriber: ", i.(int))
    })
```

To publish:

```go
    ps.Pub(1)
```

Messages are passed to subscriber which have same type argument.

License
-------

MIT: http://mattn.mit-license.org/2013

Author
------

Yasuhiro Matsumoto (mattn.jp@gmail.com)
