package usecase

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type StressTestUseCase struct {
}

func NewStressTest() *StressTestUseCase {
	return &StressTestUseCase{}
}

type StressTestUseCaseInputDTO struct {
	Url         string
	Requests    int
	Concurrency int
}

type StressTestUseCaseOutputErrorDTO struct {
	Code  int
	Total int
}

type StressTestUseCaseOutputDTO struct {
	TotalTime            time.Duration
	TotalRequests        int
	TotalSuccessRequests int
	ErrorRequests        map[int]int
}

func (uc *StressTestUseCase) Execute(input StressTestUseCaseInputDTO) (*StressTestUseCaseOutputDTO, error) {
	if err := validate(input); err != nil {
		return nil, err
	}

	start := time.Now()
	var statusErrors map[int]int = map[int]int{}
	totalSuccessRequests := 0
	totalErrorsRequests := 0
	requests := input.Requests
	concurrency := input.Concurrency
	waitGroup := sync.WaitGroup{}

	chSt := make(chan int, concurrency)
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
		concurrency := concurrency
		if (requests - executed) < concurrency {
			concurrency = requests - executed
		}
		if executed >= requests {
			break
		}
		waitGroup.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				resp, _ := http.DefaultClient.Get(input.Url)
				chSt <- resp.StatusCode
			}()
		}
		waitGroup.Wait()
		executed += concurrency

	}
	return &StressTestUseCaseOutputDTO{
		TotalTime:            time.Now().Sub(start),
		TotalRequests:        totalSuccessRequests + totalErrorsRequests,
		TotalSuccessRequests: totalSuccessRequests,
		ErrorRequests:        statusErrors,
	}, nil
}

func validate(input StressTestUseCaseInputDTO) error {
	msg := ""
	if len(input.Url) < 3 {
		msg += "invalid url"
	}

	if input.Requests < input.Concurrency {
		msg += "invalid request or concurrency"
	}

	if input.Requests == 0 {
		msg += "invalid requests"
	}

	if input.Concurrency == 0 {
		msg += "invalid concurrency"
	}

	if msg != "" {
		return errors.New(msg)
	}
	return nil
}
