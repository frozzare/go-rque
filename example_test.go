package rque_test

import (
	"fmt"
	"log"
	"time"

	"github.com/frozzare/go-rque"
)

// Example illustrates the use of rque
func Example() {
	que, err := rque.New(rque.Config{
		Address:  "localhost:28015",
		Database: "test",
		Table:    "queue",
	})

	if err != nil {
		log.Fatalf("Failed to create que: %s", err)
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		que.PostJob(rque.Job{
			ID:   "1",
			Name: "World",
		})
		time.Sleep(500 * time.Millisecond)
		que.PostJob(rque.Job{
			ID:   "1",
			Name: "Goodbye",
		})
		time.Sleep(500 * time.Millisecond)
		que.Quit()
	}()

	for job := range que.Jobs() {
		fmt.Printf("Hello %s\n", job.Name)
	}
	// Output:
	// Hello World
	// Hello Goodbye
}
