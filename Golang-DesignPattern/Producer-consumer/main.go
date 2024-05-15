package main

import (
	"fmt"
	"math/rand"
	"time"
)

// sensor simulates a sensor that generates data.
func sensor(sensorID int, dataChannel chan<- string) {
	for i := 0; ; i++ {
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		dataChannel <- fmt.Sprintf("Sensor %d: Data %d", sensorID, i)
	}
}

// mergeSensors merges data from multiple sensors into a single channel.
func mergeSensors(sensorIDs ...int) <-chan string {
	mergedData := make(chan string)

	for _, id := range sensorIDs {
		go func(sensorID int) {
			sensorData := make(chan string)
			go sensor(sensorID, sensorData)
			for {
				mergedData <- <-sensorData
			}
		}(id)
	}

	return mergedData
}

func main() {
	// Merge data from sensors 1 and 2 into a single channel.
	dataStream := mergeSensors(1, 2,3,4,5,6)

	// Read and print data from the merged channel.
	for i := 0; i < 50; i++ {
		fmt.Println(<-dataStream)
	}

	fmt.Println("End of data stream.")
}
