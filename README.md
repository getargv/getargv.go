<h1><img src="logo.svg" width="200" alt="getargv"></h1>

[![Go](https://github.com/getargv/getargv.go/actions/workflows/go.yml/badge.svg)](https://github.com/getargv/getargv.go/actions/workflows/go.yml)

This package allows you to query the arguments of other processes on macOS.

## Installation

Install the package and add to the application's `go.mod` file by executing:

    $ go get github.com/getargv/getargv.go@latest

If `go.mod` is not being used to manage dependencies, import the package with:

    import "github.com/getargv/getargv.go"

## Usage

```go
Getargv.asString(some_process_id, 0, false) #=> "arg0\x00arg1\x00"
Getargv.asBytes(some_process_id, 0, false) #=> []byte("arg0\x00arg1\x00")
Getargv.asStrings(some_process_id) #=> ["arg0","arg1"]
```

## Development

After checking out the repo, run `go test` to run the tests.

Go code goes in `getargv.go`. Test code in `getargv_test.go`.

To install this package onto your local machine, run `go install`. To release a new version, create a tag and push to GitHub which should get picked up by [pkg.go.dev](https://pkg.go.dev/).

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/getargv/getargv.go.

## License

The package is available as open source under the terms of the [BSD 3-clause License](https://opensource.org/licenses/BSD-3-Clause).
