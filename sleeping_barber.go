package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Customer struct {
	id   int
	done chan bool
}

func barber(waitingRoom <-chan *Customer, ready chan struct{}) {
	for {
		select {
		case customer := <-waitingRoom:
			fmt.Printf("Barber: Starting haircut for Customer %d\n", customer.id)
			haircutTime := time.Duration(rand.Intn(4000)+2000) * time.Millisecond
			fmt.Printf("Barber: Cutting hair for %v ms for Customer %d\n", haircutTime, customer.id)
			time.Sleep(haircutTime)
			fmt.Printf("Barber: Finished haircut for Customer %d\n", customer.id)
			customer.done <- true
		default:
			fmt.Println("Barber: No customers, barber is sleeping...")
			if ready != nil {
				ready <- struct{}{}
				ready = nil
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func receptionist(incoming <-chan *Customer, waitingRoom chan<- *Customer) {
	for customer := range incoming {
		fmt.Printf("Receptionist: Greeting Customer %d\n", customer.id)
		// Try to send the customer into the waiting room without blocking.
		select {
		case waitingRoom <- customer:
			fmt.Printf("Receptionist: Customer %d enters the waiting room\n", customer.id)
		default:
			// The waiting room is full.
			fmt.Printf("Receptionist: Waiting room full, turning away Customer %d\n", customer.id)
			// Signal the customer that they were turned away.
			customer.done <- false
		}
	}
}

func customerProcess(id int, incoming chan<- *Customer) {
	customer := &Customer{
		id:   id,
		done: make(chan bool),
	}
	fmt.Printf("Customer %d: Arriving at the shop\n", customer.id)
	// Send self to the receptionist.
	incoming <- customer

	// Wait for the signal: true means haircut done; false means turned away.
	result := <-customer.done
	if result {
		fmt.Printf("Customer %d: Got a haircut and leaves happy\n", customer.id)
	} else {
		fmt.Printf("Customer %d: Leaves because the shop is full\n", customer.id)
	}
}

func customerGenerator(incoming chan<- *Customer) {
	id := 1
	for {
		go customerProcess(id, incoming)
		id++
		// Customers arrive at random intervals (between 500ms and 5500ms).
		sleepTime := time.Duration(rand.Intn(5000)+500) * time.Millisecond
		time.Sleep(sleepTime)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	waitingRoomCapacity := 6
	waitingRoom := make(chan *Customer, waitingRoomCapacity)
	incoming := make(chan *Customer)

	barberReady := make(chan struct{})

	go barber(waitingRoom, barberReady)

	<-barberReady

	go receptionist(incoming, waitingRoom)
	go customerGenerator(incoming)

	select {}
}
