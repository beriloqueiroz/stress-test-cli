/*
Copyright © 2024 NAME HERE berilo.queiroz@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/beriloqueiroz/stress-test-cli/internal/usecase"
	"github.com/spf13/cobra"
)

func GetUseCase() *usecase.StressTestUseCase {
	return usecase.NewStressTest()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress-test-cli",
	Short: "Cli to stress test",
	Long: `This application runs a stress test using the HTTP GET method. Use the flags --url, --requests, and --concurrency to execute the stress test. 
	The --url flag specifies the endpoint to test, --requests indicates the number of requests, and --concurrency determines the number of simultaneous requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println(err)
		}
		requests, err := cmd.Flags().GetInt("requests")
		if err != nil {
			fmt.Println(err)
		}
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			fmt.Println(err)
		}
		output, err := GetUseCase().Execute(usecase.StressTestUseCaseInputDTO{
			Url:         url,
			Requests:    requests,
			Concurrency: concurrency,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("total requests:", output.TotalRequests)
		fmt.Println("total success requests:", output.TotalSuccessRequests)
		fmt.Println("total execution time:", output.TotalTime)
		for k, v := range output.ErrorRequests {
			fmt.Printf("total de erros %d: %d\n", k, v)
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("url", "", "url to get test (example: http://localhost:8080)")
	rootCmd.Flags().Int("requests", 1, "requests quantity to get test (example: 100)")
	rootCmd.Flags().Int("concurrency", 1, "requests simultaneous to get test (example: 50)")
}
