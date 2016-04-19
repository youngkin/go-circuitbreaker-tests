/**
Evaluation results:
	1.	Support for closed, half-open, and open states.
	2.	Doesn't allow access to any circuit breaker stats except whether or not it's open or closed.
	3. 	Has no concept of a disabled, infinite, timeout.
	4. 	Is in active development, commits in 2016
	5. 	Integrates with any external statsd or http client (e.g., Turbine) to publish stats externally.

	What would make this library great is programmatic access to circuit breaker state information.
*/

package cbtests

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/afex/hystrix-go/hystrix"
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

func TestSimpleSyncAfex(t *testing.T) {
	setup()

	hCmd := "my_command"
	hystrix.ConfigureCommand(hCmd, hystrix.CommandConfig{
		Timeout:                1000, // Milliseconds allowed before a request will timeout
		MaxConcurrentRequests:  100,  // Max number of concurrent requests to a dependency
		RequestVolumeThreshold: 10,   // Minimum number of requests needed before circuit breaker state is evaluated
		SleepWindow:            5000, // Time to allow an open circuit breaker to be tested to see if the dependency has recovered
		ErrorPercentThreshold:  45,   // Percentage of failed requests before circuit breaker opens
	})

	cb, ok, _ := hystrix.GetCircuit(hCmd)
	if !ok {
		t.Fatalf("Circuit named %s not found\n", hCmd)
	}

	for i := 0; i < 20; i++ {
		err := hystrix.Do(hCmd, testDependency, nil)
		fmt.Printf("Round %d: Is circuit open? %v; Result is %v\n", i, cb.IsOpen(), err)
	}

	//
	// PROBLEM: No way to access circuit breaker stats, they're all private
	//
	fmt.Printf("\n\n")
}

func setup() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
}
