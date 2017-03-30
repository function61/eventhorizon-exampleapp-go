package main

// A stress testing / reliability / high availability testing utility which bombards
// the example app with configured number of "place order" requests/second.
import (
	"github.com/function61/pyramid/util/clicommon"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	numWorkers = 20
)

// there are $numWorkers of these workers feeding off of the queue
func StressTesterWorker(workSubmit chan bool, workResponse chan *WorkResponse, wg *sync.WaitGroup, baseUrl string) {
	defer wg.Done()

	for range workSubmit {
		orderId, instanceId, err := placeOrderHttpRequest(baseUrl)
		if err != nil {
			workResponse <- &WorkResponse{"", "", err}
			continue
		}

		workResponse <- &WorkResponse{orderId, instanceId, nil}
	}
}

func NewStressTester(baseUrl string, ordersPerSec int) *StressTester {
	return &StressTester{
		workSubmit:     make(chan bool, numWorkers),
		workResponse:   make(chan *WorkResponse),
		stopProducer:   make(chan bool),
		stopDone:       make(chan bool),
		workersStopped: &sync.WaitGroup{},
		stats:          Stats{},
		baseUrl:        baseUrl,
		ordersPerSec:   ordersPerSec,
	}
}

func (f *StressTester) Run() {
	f.started = time.Now()

	// start workers
	for i := 0; i < numWorkers; i++ {
		f.workersStopped.Add(1)
		go StressTesterWorker(f.workSubmit, f.workResponse, f.workersStopped, f.baseUrl)
	}

	// work producer
	go func() {
		previous := f.stats

		for {
			select {
			case <-time.After(1 * time.Second):
				statsDiff := f.stats.Sub(previous)
				log.Printf("producer: %s", statsDiff.Serialize())

				for i := 0; i < f.ordersPerSec; i++ {
					select {
					case f.workSubmit <- true:
						// noop
					default:
						// channel was full
					}
				}
			case <-f.stopProducer:
				// stop workers
				close(f.workSubmit)

				// so no more results will be sent to workResponse
				f.workersStopped.Wait()

				// stops the result consumer
				close(f.workResponse)
				return
			}

			previous = f.stats
		}
	}()

	// work result consumer
	go func() {
		previousInstance := ""
		for result := range f.workResponse {
			if result.err != nil {
				log.Printf("responseprocessor: FAIL: %s", result.err.Error())
				f.stats.numOrdersFailed++
			} else {
				if result.instanceId != previousInstance && previousInstance != "" {
					log.Printf("responseprocessor: instance change detected: %s", result.instanceId)
				}

				f.stats.numOrdersSucceeded++

				previousInstance = result.instanceId
			}
		}

		f.stopDone <- true
	}()
}

func (f *StressTester) Close() {
	log.Printf("stopping")

	f.stopProducer <- true

	<-f.stopDone

	log.Printf("result: %s", f.stats.Serialize())
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <baseUrl> <ordersPerSec>", os.Args[0])
	}

	baseUrl := os.Args[1]
	ordersPerSec, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	f := NewStressTester(baseUrl, ordersPerSec)

	f.Run()

	clicommon.WaitForInterrupt()

	f.Close()
}
