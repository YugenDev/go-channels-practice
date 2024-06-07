package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

func main() {

	wg := &sync.WaitGroup{}
	idChan := make(chan string)
	fakeChannel := make(chan string)
	closeChans := make(chan int)

	wg.Add(3)

	go generateID(idChan, wg, closeChans)
	go generatefakeID(fakeChannel, wg, closeChans)

	go logID(idChan, wg, fakeChannel, closeChans)

	wg.Wait()
}

func generatefakeID(fakeChannel chan<- string, wg *sync.WaitGroup, closeChannel chan<- int) {
	for i := 0; i < 50; i++ {
		id := uuid.New()
		fakeChannel <- fmt.Sprintf("%d. %s", i+1, id.String())
	}

	close(fakeChannel)
	closeChannel <- 1
	wg.Done()
}

func generateID(channel chan<- string, wg *sync.WaitGroup, closeChannel chan<- int) {

	for i := 0; i < 100; i++ {
		id := uuid.New()
		channel <- fmt.Sprintf("%d. %s", i+1, id.String())
	}

	close(channel)
	closeChannel <- 1

	wg.Done()
}

func logID(channel <-chan string, wg *sync.WaitGroup, fakeChannel <-chan string, closeChannel chan int) {

	closedChannels := 0

	for {
		select {
		case id, ok := <-channel:
			if ok {
				fmt.Println(id)
			}

		case id, ok := <-fakeChannel:
			if ok {
				fmt.Println("fake ->", id)
			}

		case count, ok := <-closeChannel:
			if ok {
				closedChannels += count
			}
		}
		if closedChannels == 2 {
			close(closeChannel)
			wg.Done()
			break
		}
	}
	wg.Done()
}
