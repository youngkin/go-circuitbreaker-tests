/**
Evaluation results:
	1.	Has great support for access to metrics
	2.	Has support for call duration metrics
	3.	Has no support for a half-open circuit breaker. That is, after the circuit breaker opens, and there
		continue to be requests at a rate that exceeds the `NumberOfSecondsToStore`, then the circuit breaker
		will never recover. It will stay open.
	4. 	Has no concept of disabling the timeout (i.e., infinite timeout)
	5. 	Last significant commit was in 2014

	What this library needs is the support for a half-open circuit breaker state.
*/

package cbtests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dahernan/goHystrix"
)

var (
	errThrshd       = 45.0
	minRqsts        = 10
	windowSecs      = 5
	sampleSize      = 50
	rqstTimeoutSecs = 5
)

type MyStringCommand struct {
	message      string
	failRandomly bool
}

func (c *MyStringCommand) Run() (interface{}, error) {
	err := testRandomFailure(c.failRandomly)
	return c.message, err
}

func (c *MyStringCommand) Fallback() (interface{}, error) {
	return "FALLBACK", nil
}

func TestSimpleGoHystrix(t *testing.T) {
	command, circuit, rqstTimeoutSecs := getRandomlyFailingCommand()

	for i := 0; i < 20; i++ {
		open, msg := circuit.IsOpen()
		result, _ := command.Execute()
		fmt.Printf("Round %d; Is circuit open? %t; reason is %s; Result: %s\n", i, open, msg, result)
	}
	printLongStats("RSY: Stats - FIRST time", command)

	// Allow the time-window to move along so circuit breaker will open again (if needed)
	// This doesn't work as expected. Inspecting the library code, the health stats don't age out unless there's no
	// activity AT ALL for the specified `NumberOfSecondsToStore`. So, the only option is to NOT have any requests
	// AT ALL during the `NumberOfSecondsToStore` interval.
	fmt.Printf("\n\nContinue running commands at a slower rate to allow prior stats to roll-out of the time window (of %d seconds)\n\n", rqstTimeoutSecs)
	//time.Sleep(time.Duration(int(float64(rqstTimeoutSecs)*1.5)) * time.Second)
	cmd2, crct2, _ := getSuccessCommand()
	runSlower(cmd2, crct2, rqstTimeoutSecs*1000)

	fmt.Printf("\n\nReturn to randomly failing commands\n\n")
	for i := 0; i < 20; i++ {
		open, msg := circuit.IsOpen()
		result, _ := command.Execute()
		fmt.Printf("Round %d; Is circuit open? %t; reason is %s; Result: %s\n", i, open, msg, result)
	}
	printLongStats("RSY: Stats - SECOND time", command)

	fmt.Printf("\n\n\n")
}

func runSlower(c *goHystrix.Command, circuit *goHystrix.CircuitBreaker, timeToRunMillis int) {
	numLoops := sampleSize + 1
	sleepIntervalMs := int(float32(timeToRunMillis/numLoops) * 2.0)
	for i := 0; i < numLoops; i++ {
		open, msg := circuit.IsOpen()
		result, _ := c.Execute()
		fmt.Printf("Round %d;\t Is circuit open?\t %t; reason is\t %s; Result:\t %s\n", i, open, msg, result)
		time.Sleep(time.Duration(sleepIntervalMs) * time.Millisecond)
	}

}

func getRandomlyFailingCommand() (*goHystrix.Command, *goHystrix.CircuitBreaker, int) {

	// The combination of "stringGroup" & "stringMessage" define a Circuit Breaker and its associated stats
	command := goHystrix.NewCommandWithOptions("stringMessage", "stringGroup", &MyStringCommand{message: "helloooooooo", failRandomly: true},
		goHystrix.CommandOptions{
			ErrorsThreshold:        errThrshd,                                    // Percent errors to trip CB
			MinimumNumberOfRequest: int64(minRqsts),                              // Total number of requests before CB can trip
			NumberOfSecondsToStore: windowSecs,                                   // Time window, stats > setting will roll out of the window
			NumberOfSamplesToStore: sampleSize,                                   // Max number of requests used in the error percentage calc
			Timeout:                time.Duration(rqstTimeoutSecs) * time.Second, // How long to wait before request times out
		})

	circuits := goHystrix.Circuits()
	circuit, ok := circuits.Get("stringGroup", "stringMessage")
	if !ok {
		log.Fatal("Argh, no circuits!!!!")
	}

	fmt.Printf("\nCircuit Breaker settings:\n")
	fmt.Printf("\tError threshold:\t\t: %f \n", errThrshd)
	fmt.Printf("\tMin rqst threshold\t\t: %d\n", minRqsts)
	fmt.Printf("\tWindow size (seconds)\t\t: %d\n", windowSecs)
	fmt.Printf("\tSample size\t\t\t: %d\n", sampleSize)
	fmt.Printf("\tRequest timeout (seconds)\t: %d\n\n", rqstTimeoutSecs)

	return command, circuit, rqstTimeoutSecs
}

func getSuccessCommand() (*goHystrix.Command, *goHystrix.CircuitBreaker, int) {

	// The combination of "stringGroup" & "stringMessage" define a Circuit Breaker and its associated stats
	command := goHystrix.NewCommandWithOptions("stringMessage", "stringGroup", &MyStringCommand{message: "helloooooooo", failRandomly: false},
		goHystrix.CommandOptions{
			ErrorsThreshold:        errThrshd,                                    // Percent errors to trip CB
			MinimumNumberOfRequest: int64(minRqsts),                              // Total number of requests before CB can trip
			NumberOfSecondsToStore: windowSecs,                                   // Time window, stats > setting will roll out of the window
			NumberOfSamplesToStore: sampleSize,                                   // Max number of requests used in the error percentage calc
			Timeout:                time.Duration(rqstTimeoutSecs) * time.Second, // How long to wait before request times out
		})

	circuits := goHystrix.Circuits()
	circuit, ok := circuits.Get("stringGroup", "stringMessage")
	if !ok {
		log.Fatal("Argh, no circuits!!!!")
	}

	fmt.Printf("\nCircuit Breaker settings:\n")
	fmt.Printf("\tError threshold:\t\t: %f \n", errThrshd)
	fmt.Printf("\tMin rqst threshold\t\t: %d\n", minRqsts)
	fmt.Printf("\tWindow size (seconds)\t\t: %d\n", windowSecs)
	fmt.Printf("\tSample size\t\t\t: %d\n", sampleSize)
	fmt.Printf("\tRequest timeout (seconds)\t: %d\n\n", rqstTimeoutSecs)

	return command, circuit, rqstTimeoutSecs
}
