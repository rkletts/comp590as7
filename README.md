# comp590as7

## Design Rationale
For this project, we implemented a simulation of the barbershop problem using Go. It models a scenario where:
- A barber sleeps when there are no customers.
- A receptionist manages incoming customers and either lets them wait or turns them away.
- Customers arrive randomly and either wait to get a haircut or leave the shop if all the chairs in the waiting room are full.

Go’s concurrency tools like goroutines and channels made it easier to handle these interactions and timing challenges in a clean, efficient way.

Some Concepts We Used:
- Buffered channels represent the waiting room, which has limited capacity.
- Unbuffered channels help signal customer arrival and haircut completion.
- select statements allow non-blocking and conditional logic in goroutines.

## Overall Design
1. Customer
Each customer is modeled as a struct containing:
- An id for identification.
- A done channel used to communicate whether the customer was served or turned away.

Customers arrive asynchronously, are greeted by the receptionist, and either wait or leave depending on space in the waiting room.

2. Receptionist
The receptionist is a dedicated goroutine that:
- Accepts incoming customers.
- Attempts to place them into the buffered waitingRoom channel (which has a fixed capacity).
- Signals the customer if they must leave due to a full room.

3. Barber
The barber is a continuous loop that:
- Waits for customers in the waitingRoom channel.
- If no customers are present, prints that the barber is sleeping.
- If a customer is available, simulates a haircut via time.Sleep and then notifies the customer.

An additional ready channel is used to make sure that the barber prints the sleeping message as soon as the program starts, before the first customer arrives.

4. Customer Generator
Simulates a scenario where customers show up at random intervals between 500ms and 5500ms, by spawning new goroutines for each customer.

## What decisions did we make to get it running?
We used non-blocking select in the barber and receptionist functions to prevent deadlocks and allow for proper sleeping and rejecting behavior.

A barberReady channel ensures that the barber is fully initialized and "asleep" before customer arrival begins. This avoids output ordering issues.

The barber function includes a nil check on the ready channel to ensure that it sends the readiness signal only once.

Time durations for haircut and arrival intervals were randomized to simulate variability in a real-world scenario.
