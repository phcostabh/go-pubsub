package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/phcostabh/go-pubsub"
)

func main() {
	ps := pubsub.New()

	l, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listing", l.Addr().String())
	log.Println("Clients can connect to this server like follow:")
	log.Println("  $ telnet server:5555")
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			buf := bufio.NewReader(c)
			findex := ps.Sub(func(t interface{}) {
				if message, ok := t.(string); ok {
					log.Println(message)
					c.Write([]byte(message + "\n"))
				}
			})
			log.Println("Subscribed", c.RemoteAddr().String())
			for {
				b, _, err := buf.ReadLine()
				if err != nil {
					log.Println("Closed", c.RemoteAddr().String())
					break
				}
				ps.Pub(strings.TrimSpace(string(b)))
			}
			ps.Leave(findex)
		}(c)
	}
}
