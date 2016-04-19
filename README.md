# Overview
The purpose of these tests is to better understand the various circuit breaker libraries for Go applications. In general, the idea is to see how well each of the libraries supports the [circuit breaker pattern as defined by Martin Fowler](http://martinfowler.com/bliki/CircuitBreaker.html). This is not a rigorous evaluation per Martin Fowler's definition however. My needs were specific to the following capabilities:

1. Support for rolling windows - i.e., stats age out over time
1. Support for tripping circuit when failure rate exceeds a threshold (e.g., 50% over the last 10 calls or last 10 seconds when at least `n` calls have been made)
1. Support for other circuit breaker characteristics such as tripping after a set number of consecutive failures.
1. Support for `half-open` circuit state.
1. Support for call timeout and ability to disable said support, i.e., set the timeout to infinity.
1. Provides access to the current health state of the circuit breaker.
1. Support for integrated fallback behavior. I.e., accepts a "fallback" function to be called by the library in-place of the original function if the original function fails.
1. Is relatively current with regard to recent commit history (as of April 17, 2016)

**Caveats:**

1. I inspected library code where necessary to better understand the intended behavior of the libraries evaluated. This was primarily because the documentation across the libraries was a little light. My reading of the code may or may not match the intended behavior of one or more libraries in this evaluation. It's also possible that I misread the code. I'll happily accept updates, clarifications, and corrections to any and all of this review.
1. The evaluation code served my purposes. It is likely that additional comments and/or code changes could be made to make the code more generally useful and/or illuminate how the libriaries are intended to be used.

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

1. Has support for rolling windows - i.e., stats age out over time
1. Has support for tripping circuit when failure rate exceeds a threshold (e.g., 50% over the last 10 calls or last 10 seconds when at least `n` calls have been made)
1. Has support for other circuit breaker characteristics such as tripping after a set number of consecutive failures.
1. DOES NOT support `half-open` circuit state.
    1. After the circuit breaker opens, and there continue to be requests at a rate that exceeds the `NumberOfSecondsToStore`, then the circuit breaker will never recover. It will stay open.
	1. It has what looks like a "leaky bucket" approach to letting the circuit breaker reset itself giving it behavior similar to a "half-open" approach. However, I couldn't get this to work in my tests. Once the circuit breaker tripped it stayed tripped regardless of how many attempts to run successful commands were tried over a long time period (i.e., longer than `NumberOfSecondsToStore` and `NumberOfSamplesToStore`).
1. Supports call timeouts, doesn't allow for completely disabling it though.
1. Provides excellent access to the current health state of the circuit breaker.
    1. Unique among all the libraries was access to call duration metrics.
	1. It also supports http and statsd interfaces.
1. Has support for integrated fallback behavior.
1. The last commit, to the README, was made in January 2016. The last significant commit was made June 6, 2014.

# [github.com/dahernan/breaker](https://github.com/dahernan/breaker)
This is a followup to goHystrix by the same author. I didn't have time to test it, but it looks similar to goHystrix with the following differences:

1. It has a simpler API
1. It doesn't have integrated fallback capability
1. It doesn't support async calls
1. It doesn't support timeouts
1. The metrics capabilities aren't as rich as goHystrix, but they are more than adequate.
    1. Metrics aren't exposed via statsd or http.
1. Like goHystrix it uses a "leaky bucket" approach to letting the circuit breaker reset itself. Since I didn't have time to fully test the library I can't say if I would get different results with this library than with goHystrix.


# [github.com/afex/hystrix-go](https://github.com/afex/hystrix-go)
Evaluation results:

1. Has support for rolling windows - i.e., stats age out over time
1. Has support for tripping circuit when failure rate exceeds a threshold (e.g., 50% over the last 10 calls or last 10 seconds when at least `n` calls have been made)
1. Has support for other circuit breaker characteristics such as tripping after a set number of consecutive failures.
1. Has support `half-open` circuit state.
1. Supports call timeouts, doesn't allow for completely disabling it though.
1. DOES NOT provide programmatic access to the current health state of the circuit breaker.
    1. It does provide HTTP and Statsd interfaces for metrics though.
1. Has support for integrated fallback behavior.
1. The last commit was made September 8, 2015.

# [github.com/rubyist/circuitbreaker](https://github.com/rubyist/circuitbreaker)
Evaluation results:

1. Has support for rolling windows - i.e., stats age out over time
1. Has support for tripping circuit when failure rate exceeds a threshold (e.g., 50% over the last 10 calls or last 10 seconds when at least `n` calls have been made)
1. Has support for other circuit breaker characteristics such as tripping after a set number of consecutive failures.
1. Has support `half-open` circuit state.
1. Supports call timeouts INCLUDING the ability to disable call timeouts.
1. Provides access to the current health state of the circuit breaker.
1. DOES NOT support integrated fallback behavior.
    1. Although not ideal, it is possible to manually implement a fallback strategy using other public functions (primitives) of the library. This does impose a responsibility on the user to correctly use these primitives to build up a working circuit breaker.
1. The last commit was made February 2, 2016.

Other items of note include:

1. 	It only supports sync calls, no async support
2. 	Calls to `Breaker.Ready()` have side effects that turn a half-open circuit to closed. So a first call to `Ready()` may return true but calling it again (e.g., in a `Breaker.Call()` invocation) will return false. Lesson learned, don't use `Breaker.Ready()` to test the status of a circuit breaker if you are also using `Breaker.Call()`.

# [github.com/eapache/go-resiliency](https://github.com/eapache/go-resiliency)
Looked at eapache/go-resiliency as an alternative, but it doesn't allow access to metrics so no further evaluation was performed.

# [github.com/go-kit/git](https://github.com/go-kit/kit)
Uses `github.com/afex/hystrix-go` under the covers, so it has the same characteristics. This library does much much more than just a circuit breaker.
