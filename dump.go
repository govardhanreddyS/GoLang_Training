package main
/*
./concurrency-sync-benchmark <millions operations per run> <maximum number of workers>

the first argument it's the number of operations per each iteration in millions; this many increments and this many reads will be executed on each iteration
the second argument if maximum number of gouroutines - the same number for readers and writers; the application will iterate starting from 1 reader and 1
 writer and than gradually increase up to this argument but excluding iterations in which the number of operations can't be equally divided to workers
results are all in nanoseconds per operation
ROMutex and WOMutex are actually RWMutexes but only with readers or writers which should be the best case scenario for RLock and WLock
ROAtomic and WOAtomic are atomic primitives with only readers or writers; they actually measure the speed of Load and Add
*/

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type readerCallback func(int, *sync.WaitGroup) int64
type writerCallback func(int, *sync.WaitGroup)

var mutex sync.Mutex
var rwmutex sync.RWMutex
var ch chan int64
var buffCh chan int64
var val int64
var atomicVal atomic.Int64

// Mutexes

func mutexReader(iterations int, wg *sync.WaitGroup) int64 {
	var x int64
	for i := 0; i < iterations; i++ {
		mutex.Lock()
		x = val
		mutex.Unlock()
	}
	wg.Done()
	return x
}

func mutexWriter(iterations int, wg *sync.WaitGroup) {
	for i := 0; i < iterations; i++ {
		mutex.Lock()
		val++
		mutex.Unlock()
	}
	wg.Done()
}

// RWMutexes
func rwMutexReader(iterations int, wg *sync.WaitGroup) int64 {
	var x int64
	for i := 0; i < iterations; i++ {
		rwmutex.RLock()
		x = val
		rwmutex.RUnlock()
	}
	wg.Done()
	return x
}

func rwMutexWriter(iterations int, wg *sync.WaitGroup) {
	for i := 0; i < iterations; i++ {
		rwmutex.Lock()
		val++
		rwmutex.Unlock()
	}
	wg.Done()
}

// Channels

func chanReader(iterations int, wg *sync.WaitGroup) int64 {
	var x int64
	for i := 0; i < iterations; i++ {
		x = <-ch
	}
	wg.Done()
	return x
}

func chanWriter(iterations int, wg *sync.WaitGroup) {
	var i int64
	for i = 0; i < int64(iterations); i++ {
		ch <- i + 1 // we increment because we must be fair to other implementations
	}
	wg.Done()
}

// Buffered Channels

func buffChanReader(iterations int, wg *sync.WaitGroup) int64 {
	var x int64
	for i := 0; i < iterations; i++ {
		x = <-buffCh
	}
	wg.Done()
	return x
}

func buffChanWriter(iterations int, wg *sync.WaitGroup) {
	var i int64
	for i = 0; i < int64(iterations); i++ {
		buffCh <- i + 1 // we increment because we must be fair to other implementations
	}
	wg.Done()
}

// Atomic

func atomicReader(iterations int, wg *sync.WaitGroup) int64 {
	var x int64
	for i := 0; i < iterations; i++ {
		x = atomicVal.Load()
	}
	wg.Done()
	return x
}

func atomicWriter(iterations int, wg *sync.WaitGroup) {
	for i := 0; i < iterations; i++ {
		atomicVal.Add(1)
	}
	wg.Done()
}

// spin workers and measure time

func run(reader readerCallback, writer writerCallback, numReadWorkers int, numWriteWorkers int, iterations int) time.Duration {
	var wg sync.WaitGroup
	wg.Add(numReadWorkers + numWriteWorkers)
	startTime := time.Now()
	for i := 0; i < numReadWorkers; i++ {
		go reader(iterations, &wg)
	}
	for i := 0; i < numWriteWorkers; i++ {
		go writer(iterations, &wg)
	}
	wg.Wait()
	return time.Since(startTime)
}

func printRun(reader readerCallback, writer writerCallback, numReadWorkers int, numWriteWorkers int, numOperations int) {
	numOperationsPerWorker := int(numOperations / (numReadWorkers + numWriteWorkers))
	duration := run(reader, writer, numReadWorkers, numWriteWorkers, numOperationsPerWorker)
	fmt.Printf("%d\t\t", duration.Nanoseconds()/int64(numOperations))
}

// main

func main() {

	if len(os.Args) < 3 {
		panic("Usage: concurrency-sync-benchmark <millions operations per run> <maximum number of workers>")
	}

	numOperations, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Number of operations must be integer")
	}
	numOperations = numOperations * 1e6

	maxWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Maximum number of workers must be integer")
	}

	ch = make(chan int64)
	buffCh = make(chan int64, numOperations)

	fmt.Println("Workers\t\tMutex(ns)\tRWMutex(ns)\tROMutex(ns)\tWOMutext(ns)\tChannels(ns)\tBuffChannels\tAtomic(ns)\tROAtomic(ns)\tWOAtomic(ns)")

	for numWorkers := 1; numWorkers < maxWorkers; numWorkers++ {
		if numOperations%numWorkers != 0 {
			continue
		}
		val = 0
		atomicVal.Store(0)
		fmt.Print(numWorkers, "\t\t")
		printRun(mutexReader, mutexWriter, numWorkers, numWorkers, numOperations)
		printRun(rwMutexReader, rwMutexWriter, numWorkers, numWorkers, numOperations)
		printRun(rwMutexReader, rwMutexWriter, numWorkers, 0, numOperations*2)
		printRun(rwMutexReader, rwMutexWriter, 0, numWorkers, numOperations*2)
		printRun(chanReader, chanWriter, numWorkers, numWorkers, numOperations)
		printRun(buffChanReader, buffChanWriter, numWorkers, numWorkers, numOperations)
		printRun(atomicReader, atomicWriter, numWorkers, numWorkers, numOperations)
		printRun(atomicReader, atomicWriter, numWorkers, 0, numOperations*2)
		printRun(atomicReader, atomicWriter, 0, numWorkers, numOperations*2)
		fmt.Println(" ")
	}

}

/*

Observations
The clear winner (at least in my experiment) is atomic primitives. These are blazingly fast compared to anything else and give consistent results
Atomic writes (actually additions) are pretty expensive compared to atomic reads; atomic reads actually have consistently 0 ns per operation; it's possible that some compiler or CPU optimization kicks in and sees that the value is never actually used nor changed
The infamous mutex comes (pretty much) second especially in high concurrency but an order of magnitude slower than atomic; still much faster than alternatives
The read-write mutex would be third but only for high concurrency or high read-to-write ratio; in low concurrency with equal read-write is actually by far the worst. In a scenarios where there are multiple reads and few writes the RWMutex moves clearly to second place.
Buffered channels would probably be fourth for high concurrency but even better in very low concurrency
Unbuffered channels perform the worst - about 2-3 times worse than buffered channels in high concurrency and much worse for very low concurrency
As a side note timings include the worker spin-up time. So it seems that spinning up about 60k gorroutines takes close to no time at all if we look at ROAtomic unless the compiler was so smart as to detect that we weren't doing anything useful and didn't spin up any goroutines at all. Even in this scenario for other atomic operations it seems that it did spin up goroutines so we would still probably be in the less than a nanosecond territory (results remain consistent regardless of whether we spin up 2 or 60k goroutines) for launching a goroutine which seems pretty impressive.
There are interesting things happening in the first part of the table. Leaving aside atomic primitives (which don't seem to be affected by anything) and a few glitches we see this:

Time per operation for mutexes first increases and than decreases and than it becomes pretty stable
Timings for buffered channels increases until it becomes stable
Unbuffered channels and read-write mutexes get horrible results for low concurrency and decrease as the number of goroutines go up until they become stable as well
In some previous tests stabilization occured at about 10 - 16 goroutines which coincidentally corresponds to (hardware) thread count of the CPU. My assumption is that poor results of low concurrency has something to do with the dynamic frequency scaling of the CPU, which I didn't disable because I wanted results as close to real-world as possible.

Cheatsheet
Use atomic primitives if:

the resources you're trying to protect are of basic types supported by the atomic package
and you want best performance possible no matter what
and you don't care much about the Go mantra, don't go to Go church and generally you don't give a ... about anything except performance
Use simple mutexes if:

the resource you're protecting is not supported by the atomic package but small in size (don't use mutex on a map with a billion items - you'll see why)
and number of readers and writers are pretty balanced but not necessarily equal
or you want to protect a larger portion of code (why would you? don't! it's evil! seriously!)
and concurrency is either very low or pretty big (not just a few goroutines kind of thing)
and you still want the damn performance
and you stil don't give a ... about... - I think you get the picture
Use read-write mutexes if:

if everything said for mutexes applies except that you have a high number of reads and few writes
Use buffered channels if:

the resource you're protecting is not supported by atomic
or you to return from function as soon as possible and the resource you can use in background won't be available anymore and you don't mind to deep-copy it and the possible caveats of this
or you favour consistency and code elegance over pure performance
or you do go to Go church
and you have equal amounts of reads and writes
and your buffer is big enough
Use unbuffered channels if:

all said for buffered channels applies but
for some reason you need to exchange data synchronously
or you have memory constraints
or you don't care about performance at all
or... I can't think right now of anything else
Conclusion
Atomic primitives are the best when it comes to performance but they only support some very basic data types and therefore most of the time they won't fit your needs. Mutexes are still the next best thing when you're sure that protected resources will still be available (eg: global variables) and if (a very big if) you don't mind the caveats (eg: you forget to unlock or unlock at the wrong moment or you lock for too long etc.). So I totally get why Go introduced the channel concept and why people advocate so much for it. It's clean, elegant, modern and fits many scenarios. But still... There's a very strong case for having and using the other synchronization methods when and where they best fit.

You might argue why bother at all. After all we're in nanosecond territory (most of the time). Who cares about nanoseconds? Go the Go way and just write simple elegant code that scales beautifully and consistently. However you should remember that it all adds up. Depending on your particular application you might have thousands, millions or billions of synchronization points in your code. It would add up to microseconds, milliseconds and so on. Depending on your needs you might care or not. But it's always good to know.

About
A benchmark for different concurrency synchronization options available in Go (Golang)

*/
