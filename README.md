# Overview
The purpose of these tests is to better understand the various circuit breaker libraries for Go applications. In general, the idea is to see how well each of the libraries supports the [circuit breaker pattern as defined by Martin Fowler](http://martinfowler.com/bliki/CircuitBreaker.html). This is not a rigorous evaluation per Martin Fowler's definition however. My needs were specific to the following capabilities:
1. Support for rolling windows - i.e., stats age out over time
1. Support for tripping circuit when failure rate exceeds a threshold (e.g., 50% over the last 10 calls or last 10 seconds when at least `n` calls have been made)
1. Support for other circuit breaker characteristics such as tripping after a set number of consecutive failures.
1. Support for `half-open` circuit state.
1. Provides access to the current health state of the circuit breaker.
1. Support for integrated fallback behavior.

# Running the tests
There is a test file for the libraries that met an initial evaluation made by looking at README files and code.

Run all tests:

```go test -v```

Run a tests for a specific library

```go test -v rubyist_test.go common.go```

Run a specific test

``` go test -v github.com/youngkin/go-circuitbreaker-tests -run TestComplexRubyist```

# [github.com/dahernan/goHystrix](https://github.com/dahernan/goHystrix)
Evaluation results:

1.	Has great support for access to metrics
2.	Has support for call duration metrics
3.	Has no support for a half-open circuit breaker. That is, after the circuit breaker opens, and there continue to be requests at a rate that exceeds the `NumberOfSecondsToStore`, then the circuit breaker will never recover. It will stay open.
4. 	Has no concept of disabling the timeout (i.e., infinite timeout)
5. 	Last significant commit was in 2014

	What this library needs is the support for a half-open circuit breaker state.

# [github.com/afex/hystrix-go](https://github.com/afex/hystrix-go)
Evaluation results:

1.	Support for closed, half-open, and open states.
2.	Doesn't allow access to any circuit breaker stats except whether or not it's open or closed.
3. 	Has no concept of a disabled, infinite, timeout.
4. 	Is in active development, commits in 2016
5. 	Integrates with any external http client (e.g., Turbine) to publish stats externally.

	What this library needs is to allow access to the circuit breaker stats. Instead they're mixed in with the overall control of the circuit breaker state. This makes it dangerous to expose the states because it also exposes the internal workings of the circuit breaker.

# [github.com/rubyist/circuitbreaker](https://github.com/rubyist/circuitbreaker)
Evaluation results:

1.	No integrated support for a fallback function. Fallback behavior has to be manually integrated into the circuit breaker mechanism.
2.	Nice, can disable timeout (i.e., make infinite).
3. 	Is in active development, commits in 2016
4. 	It only supports sync calls, no asycn support
5. 	Calls to cb.Ready() have side effects that turn a half-open circuit to closed. So a first call to Ready() may return true but calling it again (e.g., in a cb.Call() invocation) will return false.

# [github.com/eapache/go-resiliency](https://github.com/eapache/go-resiliency)
Looked at eapache/go-resiliency as an alternative, but it doesn't allow access to metrics so no further evaluation was performed.

# [github.com/go-kit/git](https://github.com/go-kit/kit)
Uses `github.com/afex/hystrix-go` under the covers, so it has the same characteristics. This library does much much more than just a circuit breaker.
