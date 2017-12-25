package main

import (
	"fmt"
	"strings"
	"time"
)

func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		fmt.Println("worker", id, "on jobs", j)
		results <- strings.ToUpper(j)
		time.Sleep(1 * time.Second)
	}
}

func DoUpdates(results <-chan string) {
	for result := range results {
		fmt.Println("  [ DoUpdates ] Consumed", result)
	}
}

func main() {

	jobs := make(chan string)
	results := make(chan string)

	for w := 1; w <= 10; w++ {
		go worker(w, jobs, results)
	}

	go func() {
		fmt.Println("Inserting jobs")
		for i := 0; i <= 1000; i++ {
			jobs <- fmt.Sprintf("jobs%v", i)
			fmt.Println("  inserted", i, "to queue. len(jobs)=", len(jobs))
		}
	}()
	go DoUpdates(results)

	for {
	}
}
