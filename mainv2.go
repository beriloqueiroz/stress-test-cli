/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type StatusDuration struct {
	Code     int
	Duration time.Duration
}

func main() {
	start := time.Now()
	var statusErrors map[int]int = map[int]int{}
	totalSuccessRequests := 0
	totalErrorsRequests := 0
	Requests := 500
	Concurrency := 300
	waitGroup := sync.WaitGroup{}

	chSt := make(chan int, Concurrency)
	go func() {
		for x := range chSt {
			if x == 200 {
				totalSuccessRequests++
				waitGroup.Done()
				continue
			}
			totalErrorsRequests++
			statusError, ok := statusErrors[x]
			if !ok {
				statusErrors[x] = 1
				waitGroup.Done()
				continue
			}
			statusErrors[x] = statusError + 1
			waitGroup.Done()
		}
		close(chSt)
	}()

	executed := 0
	for {
		concurrency := Concurrency
		if (Requests - executed) < concurrency {
			concurrency = Requests - executed
		}
		if executed >= Requests {
			break
		}
		waitGroup.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				resp, _ := http.DefaultClient.Get("http://localhost:8080")
				chSt <- resp.StatusCode
			}()
		}
		waitGroup.Wait()
		executed += concurrency

	}
	totalDuration := time.Now().Sub(start)
	fmt.Println("total requests", totalSuccessRequests+totalErrorsRequests)
	fmt.Println("total success requests", totalSuccessRequests)
	fmt.Println("total errors requests", totalErrorsRequests)
	fmt.Println("total duration execution", totalDuration)
	for k, v := range statusErrors {
		fmt.Printf("status code: %d = %d\n", k, v)
	}
}
