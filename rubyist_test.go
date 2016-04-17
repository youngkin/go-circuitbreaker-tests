/**
Evaluation results:
	1.	No integrated support for a fallback function. Fallback behavior has to be manually integrated into
		the circuit breaker mechanism.
	2.	Nice, can disable timeout (i.e., make infinite).
	3. 	Is in active development, commits in 2016
	4. 	It only supports sync calls, no asycn support
	5. 	Calls to cb.Ready() have side effects that turn a half-open circuit to closed. So a first call to Ready()
		may return true but calling it again (e.g., in a cb.Call() invocation) will return false.

	Looked at eapache/go-resiliency as an alternative, but it doesn't allow any access to metrics so no use
	doing a further evaluation.
*/

package cbtests

import (
	"fmt"
	"testing"

	"time"

	"github.com/rubyist/circuitbreaker"
)

/**
Tests:
	1.	Simple test just to see how it works
	2.	Test with simulated errors and config
	3.	Access/Report stats on circuit breaker.
	4.	Test with simulated errors and test recovery (e.g., protected function starts working again)
	5. 	Test with fallback function to simulate alternative behavior
	6.	Test with fallback function to simulate failing program when circuit opens
	7.	If time, look into async behavior.

	Some other things to check:
	1.	Does the library log? This is bad
	2.	Do the stats look right? goHystrix seems to have some variability in the reported stats results. Sometimes
		they're dead-on, sometimes they seem to be off by 1.
	3. 	Understand the stats
*/

var errorThreshold = 0.35

func TestSimpleRubyist(t *testing.T) {
	// Creates a circuit breaker based on the error rate
	cb := circuit.NewBreakerWithOptions(&circuit.Options{
		ShouldTrip: circuit.RateTripFunc(errorThreshold, 10),
	}) // trip when error rate hits 45%, with at least 10 samples

	fmt.Printf("\nCircuit Breaker settings:\n")
	fmt.Printf("\tError threshold:\t\t: %f \n", errorThreshold)
	fmt.Printf("\tMin rqst threshold\t\t: %d\n", 10)
	fmt.Printf("\tWindow size (seconds)\t\t: %d\n", circuit.DefaultWindowTime/time.Second)
	fmt.Printf("\tSample size\t\t\t: %d\n", circuit.DefaultWindowBuckets)
	fmt.Printf("\tHalf-open window\t\t: %d (ms) \n\n", 500)
	fmt.Printf("\tRequest timeout (seconds)\t: %d\n\n", 0)

	for i := 0; i < 20; i++ {
		open := cb.Tripped()
		err := cb.Call(testDependency, 0)
		fmt.Printf("Round %d;\tIs circuit ok? %t;\tError: %v\n", i, !open, err)
	}

	successes := int(cb.Successes())
	failures := int(cb.Failures())
	consecFailures := int(cb.ConsecFailures())
	errorRate := cb.ErrorRate()

	printStats("Rubyist Test: Round 1", successes+failures, successes, failures, consecFailures, errorRate)

	// Give the circuit breaker time to transition to half-open or closed
	fmt.Print("\n\nContinue running commands at a slower rate to allow circuit breaker to maybe transition to half-open/closed\n\n")
	runSlowerCB(cb, 5000, 5)

	// Try again with a, hopefully, recovered/closed circuit breaker
	for i := 0; i < 20; i++ {
		open := cb.Tripped()
		err := cb.Call(testDependency, 0)
		fmt.Printf("Round %d;\tIs circuit ok? %t;\tError: %v\n", i, !open, err)
	}

	successes = int(cb.Successes())
	failures = int(cb.Failures())
	consecFailures = int(cb.ConsecFailures())
	errorRate = cb.ErrorRate()

	printStats("Rubyist Test: Round 2 (half-open)", successes+failures, successes, failures, consecFailures, errorRate)

	// Give the circuit breaker time to transition to closed
	fmt.Print("\n\nSleep until CB transitions to closed\n\n")
	time.Sleep(2 * circuit.DefaultWindowTime)

	// Try again with a, hopefully, closed circuit breaker
	for i := 0; i < 20; i++ {
		open := cb.Tripped()
		err := cb.Call(testDependency, 0)
		fmt.Printf("Round %d;\tIs circuit ok? %t;\tError: %v\n", i, !open, err)
	}

	successes = int(cb.Successes())
	failures = int(cb.Failures())
	consecFailures = int(cb.ConsecFailures())
	errorRate = cb.ErrorRate()

	printStats("Rubyist Test: Round 3 (fully closed)", successes+failures, successes, failures, consecFailures, errorRate)
}

func TestComplexRubyist(t *testing.T) {
	// Creates a circuit breaker based on the error rate
	cb := circuit.NewBreakerWithOptions(&circuit.Options{
		ShouldTrip: circuit.RateTripFunc(errorThreshold, 10),
	}) // trip when error rate hits 45%, with at least 10 samples

	fmt.Printf("\nCircuit Breaker settings:\n")
	fmt.Printf("\tError threshold:\t\t: %f \n", errorThreshold)
	fmt.Printf("\tMin rqst threshold\t\t: %d\n", 10)
	fmt.Printf("\tWindow size (seconds)\t\t: %d\n", circuit.DefaultWindowTime/time.Second)
	fmt.Printf("\tSample size\t\t\t: %d\n", circuit.DefaultWindowBuckets)
	fmt.Printf("\tHalf-open window\t\t: %d (ms) \n", 500)
	fmt.Printf("\tRequest timeout (seconds)\t: %d\n\n", 0)

	for i := 0; i < 20; i++ {
		open := cb.Tripped()
		err := cb.Call(testDependency, 0)
		fmt.Printf("Round %d;\tIs circuit ok? %t;\tError: %v\n", i, !open, err)
	}

	successes := int(cb.Successes())
	failures := int(cb.Failures())
	consecFailures := int(cb.ConsecFailures())
	errorRate := cb.ErrorRate()

	printStats("Rubyist Test: Round 1", successes+failures, successes, failures, consecFailures, errorRate)

	// Let circuit transition to half-open
	fmt.Print("\n\nSleep until CB transitions to half-open/closed\n\n")
	time.Sleep(1000 + circuit.DefaultWindowTime)

	// These should all work
	for i := 0; i < 20; i++ {
		tripped := cb.Tripped()
		var err error
		err = cb.Call(func() error {
			return nil
		}, 0)
		fmt.Printf("Round %d;\tIs circuit tripped? %v;\tError: %v\n", i, tripped, err)
		if tripped && err == nil {
			fmt.Printf("\tIf tripped is true and Error is nil, then the circuit was half-open\n")
		}
	}

	successes = int(cb.Successes())
	failures = int(cb.Failures())
	consecFailures = int(cb.ConsecFailures())
	errorRate = cb.ErrorRate()

	printStats("Rubyist Test: Round 2 (half-open, all successes)", successes+failures, successes, failures, consecFailures, errorRate)
}

func runSlowerCB(cb *circuit.Breaker, timeToRunMillis, numLoops int) {
	sleepIntervalMs := int(float32(timeToRunMillis/numLoops) * 1.5)
	for i := 0; i < numLoops; i++ {
		cb.Call(testDependency, 0)
		time.Sleep(time.Duration(sleepIntervalMs) * time.Millisecond)
	}

}
