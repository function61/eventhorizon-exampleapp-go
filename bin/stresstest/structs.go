package main

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	numOrdersSucceeded int
	numOrdersFailed    int
}

func (s Stats) Sub(other Stats) Stats {
	return Stats{
		numOrdersSucceeded: s.numOrdersSucceeded - other.numOrdersSucceeded,
		numOrdersFailed:    s.numOrdersFailed - other.numOrdersFailed,
	}
}

func (s *Stats) Serialize() string {
	return fmt.Sprintf("failed = %d succeeded = %d", s.numOrdersFailed, s.numOrdersSucceeded)
}

type WorkResponse struct {
	orderId    string
	instanceId string
	err        error
}

type StressTester struct {
	workSubmit     chan bool
	workResponse   chan *WorkResponse
	stopProducer   chan bool
	stopDone       chan bool
	workersStopped *sync.WaitGroup
	stats          Stats
	started        time.Time
	baseUrl        string
	ordersPerSec   int
}
