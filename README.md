# backoff-go

Backoff is a small and simple library that provides basic backoff functionality for Go. This can be useful for throttling retries and preventing flooding other applications/services/systems with requests.

The core of the library is the `Backoff` interface which has a single method 'Next' that returns the next duration to wait. The library comes with three implementations of Backoff to over common use cases.

* ConstantBackoff - Always returns the same duration
* ExponentialBackoff - Returns an exponentially increasing duration with a 25% jitter up to a max value.
* RandomBackoff - Returns a random duration between a min and max value.

## Installation

```bash
go get github.com/jkratz55/backoff-go
```
