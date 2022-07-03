package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	maxQueue   = 100
	maxWorkers = 5
)

var (
	WorkQueue   = make(chan Work, maxQueue)
	WorkerQueue = make(chan chan Work, maxWorkers)
)

type Work struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

type Worker struct {
	ID          int
	Work        chan Work
	WorkerQueue chan chan Work
}

func Dispatcher() {
	WorkerQueue = make(chan chan Work, maxWorkers)

	for i := 1; i <= maxWorkers; i++ {
		fmt.Println("Start worker", i)
		w := Worker{
			ID:          i,
			Work:        make(chan Work),
			WorkerQueue: WorkerQueue,
		}
		go w.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work request")
				go func() {
					worker := <-WorkerQueue

					fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

func (w Worker) Start() {
	for {
		w.WorkerQueue <- w.Work

		select {
		case work := <-w.Work:
			fmt.Printf("worker%d: Received work request, delaying for %d seconds\n", w.ID, work.Delay)

			time.Sleep(time.Duration(work.Delay) * time.Second)
			fmt.Printf("worker%d: title: %s | message: %s!\n", w.ID, work.Title, work.Message)
		}
	}
}

func routeHandler() {
	http.HandleFunc("/work", Collector)
}

func Collector(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var work Work
		err := json.NewDecoder(r.Body).Decode(&work)
		if err != nil {
			fmt.Println(err.Error())
		}
		WorkQueue <- work
		fmt.Println("queued work request")
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	routeHandler()
	Dispatcher()
	http.ListenAndServe(":8000", nil)
}
