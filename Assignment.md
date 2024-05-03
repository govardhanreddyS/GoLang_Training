Assignment 1: Goroutine with Channel
Problem:
Write a Go program that calculates the sum of numbers from 1 to N concurrently using goroutines and channels.
The program should take the value of N as input from the user.

Assignment 2: Producer-Consumer with Channel
Problem:
Implement the producer-consumer problem using goroutines and channels. 
The producer should generate numbers from 1 to 100 and send them to a channel, and the consumer should print those numbers.

Assignment 3: Mutex for Synchronization
Problem:
Write a Go program that uses mutexes to synchronize access to a shared variable.
Multiple goroutines should increment the variable concurrently, and the final value should be printed.

Assignment 4: Goroutine with Timeout
Problem:
Write a Go program that performs a task in a goroutine and waits for it to finish. However,
if the task takes more than 3 seconds, the program should print a timeout message and exit.

Assignment 5: Pipeline Pattern
Problem:
Implement a pipeline pattern in Go where one goroutine generates numbers, another squares them, and a third prints the squared numbers.
