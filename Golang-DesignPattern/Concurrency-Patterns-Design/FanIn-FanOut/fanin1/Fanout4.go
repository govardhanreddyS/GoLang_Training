/*

The scenario where the pattern provided might be useful is in a distributed system 
where a central component needs to distribute tasks or messages to multiple worker
 nodes for processing. Let's say we have a system where messages are received from 
 various sources and need to be processed concurrently by multiple workers. 
 Each worker could be responsible for a different type of processing, such as analyzing
  sentiment, extracting keywords, or translating the message.
In this scenario:

The generator could be a component that receives incoming messages from different sources.
Each message could be represented by a struct similar to the Message struct provided earlier, 
containing information like the sender, recipient, content, and timestamp.
The generator would send these messages to a channel.
Multiple worker nodes, represented by the consumers in the example, would read messages 
from the channel concurrently.
Each worker node could perform its specific task on the received message.
After processing, the worker nodes could send the processed data to another channel for further 
processing or storage.
This pattern allows for efficient concurrent processing of incoming messages, improving the overall 
throughput and scalability of the system.
*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Mappable interface {
	Index() int
	Process() Mappable
}

type value struct {
	name  string
	index int
}

// Index returns the original position of the object
func (v *value) Index() int {
	return v.index
}

// Process the value and return the result as the same type but copied
func (v *value) Process() Mappable {
	// Random sleep between 1-5s handled.
	time.Sleep(time.Duration(rand.Intn(5-1)+1) * time.Second)

	return &value{name: v.name, index: v.index}
}

func (v *value) String() string {
	return fmt.Sprintf("{index: %d, value: %s}", v.index, v.name)
}

func main() {
	numOfWorkers := runtime.NumCPU()
	fmt.Fprintln(os.Stdout, "Running on", numOfWorkers, "goroutines.")

	var (
		inChan        = make(chan Mappable)     // Input values
		outChanValues = make(chan Mappable, 10) // Output values
	)

	// Create consumers
	wg := &sync.WaitGroup{} // Waitgroup for workers
	wg.Add(numOfWorkers)

	for s := 0; s < numOfWorkers; s++ {
		go fanoutWorker(wg, inChan, s, outChanValues)
	}

	// Input data
	inputValues := strings.Fields(Names)
	go inputData(inChan, inputValues)

	go func() {
		// Once input data is treated and all workers have returned close the output channel
		wg.Wait()
		close(outChanValues)
	}()

	outputValues := make([]string, len(inputValues)) // Collected Output values

	for v := range outChanValues {
		fmt.Fprintf(os.Stdout, "Success: %v\n", v)
		realValue := v.(*value)
		outputValues[v.Index()] = realValue.name
	}

	if Names == strings.Join(outputValues, " ") {
		fmt.Fprintln(os.Stdout, "Order was respected, input order is the same as output.")
	} else {
		fmt.Fprintln(os.Stdout, "Order was not respected, input order is not the same as output.")
	}
}

// Insert data into the input channel and signal it's done
func inputData(inChan chan<- Mappable, inputValues []string) {
	for i, v := range inputValues {
		inChan <- &value{name: v, index: i}
	}

	close(inChan)
}

func fanoutWorker(wg *sync.WaitGroup, inChan <-chan Mappable,
	routineName int, valOut chan<- Mappable) {
	defer wg.Done()

	for name := range inChan {
		valOut <- name.Process()
	}
}

const Names = "इस उत्तर प्रदेश में एक खुशहाल गाँव था, जहाँ सभी लोग खुशहाल थे। यहाँ के लोग अपने आप को बहुत धन्य मानते थे क्योंकि उनके पास सभी सुविधाएँ थीं। वहाँ के लोगों का जीवन शांत और संतुलित था। वे सभी मिल-जुलकर रहते थे और एक-दूसरे की मदद करते थे।"
