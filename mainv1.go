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
	Rest := Requests % Concurrency
	fmt.Println(Rest)
	waitGroup := sync.WaitGroup{}
	relationship := Requests / Concurrency
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

	for j := 0; j < relationship; j++ {
		waitGroup.Add(Concurrency)
		for i := 0; i < Concurrency; i++ {
			go func() {
				resp, _ := http.DefaultClient.Get("http://localhost:8080")
				chSt <- resp.StatusCode
			}()
		}
		waitGroup.Wait()
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
