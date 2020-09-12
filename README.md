[![PkgGoDev](https://pkg.go.dev/badge/github.com/ninedraft/floatring)](https://pkg.go.dev/github.com/ninedraft/floatring) [![report card](https://goreportcard.com/badge/github.com/ninedraft/floatring)](https://goreportcard.com/report/github.com/ninedraft/floatring) [![codecov](https://codecov.io/gh/ninedraft/floatring/branch/master/graph/badge.svg)](https://codecov.io/gh/ninedraft/floatring)

# floatring

Ring buffer for floats

[Wikipedia](https://en.wikipedia.org/wiki/Circular_buffer)

## License

MIT License

## Usage

Add to the dependencies: `go get github.com/ninedraft/floatring`

A circular floating point buffer is a high performance, limited queue of numbers. It is convenient to use for storing the latest time series records for analytics, creating various filters and buffering data streams that are resistant to data loss.
