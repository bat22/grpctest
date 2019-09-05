package main

import (
	"context"
	"fmt"
	"github.com/bat22/grpctest/internal/rpc"
	"google.golang.org/grpc"
	"time"
)

const workerCount = 50

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("127.0.0.1:55000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("connected")

	c := rpc.NewTestSvcClient(conn)

	for i := 0; i < workerCount; i++ {
		go testAdd(c)
	}

	var stopCh chan struct{}
	<-stopCh
}

func testAdd(c rpc.TestSvcClient) {
	request := &rpc.AddIntMessage{A: 1, B: 2}
	for {
		t := time.Now()
		_, err := c.Add(context.Background(), request)
		statCh <- time.Since(t)
		if err != nil {
			panic(err)
		}
	}
}

func stat(statCh chan time.Duration) {
	ticker := time.NewTicker(time.Second)
	var c int
	var sum, min, max time.Duration
	for {
		select {
		case <-ticker.C:
			if c > 0 {
				fmt.Println("count:", c, "min/avg/max:", min, sum/time.Duration(c), max)
				c = 0
				sum = 0
			} else {
				fmt.Println("count: 0 min/avg/max: 0 0 0")
			}
		case d := <-statCh:
			if c == 0 {
				min = d
				max = d
			} else {
				if min > d {
					min = d
				}
				if max < d {
					max = d
				}
			}
			c++
			sum += d
		}
	}
}

var statCh chan time.Duration

func init() {
	statCh = make(chan time.Duration, 1000000)
	go stat(statCh)
}
