/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"time"

	"github.com/beriloqueiroz/stress-test-cli/internal/usecase"
)

type StatusDuration struct {
	Code     int
	Duration time.Duration
}

func main() {

	stressUseCase := usecase.NewStressTest()

	output, err := stressUseCase.Execute(usecase.StressTestUseCaseInputDTO{
		Url:         "http://localhost:8080",
		Requests:    1250,
		Concurrency: 301,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("total requests", output.TotalRequests)
	fmt.Println("total success requests", output.TotalSuccessRequests)
	fmt.Println("total errors requests", output.TotalRequests-output.TotalSuccessRequests)
	fmt.Println("total duration execution", output.TotalTime)
	for k, v := range output.ErrorRequests {
		fmt.Printf("status code: %d = %d\n", k, v)
	}
}
