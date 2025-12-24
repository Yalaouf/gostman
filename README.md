# gostman
An API tester written in Go

## Why

I was sick of using bloated electron based apps for simple API testing. I wanted a
lightweight tool that could do the job without consuming too many resources.\
There are already some good alternatives like [HTTPie](https://github.com/httpie/cli),
[curl](https://github.com/curl/curl), or [xh](https://github.com/ducaale/xh) but
I wanted something with a TUI (not GUI to stay in flow mode) that was simple and easy to use.\
Also I wanted no compromise for simplicity over features, so I decided to write my own.

## Why Go

[Go](https://go.dev/) is fast, compiles to a single binary, and can be compiled for multiple platforms easily.
Also Go provides a good [standard library for HTTP requests](https://pkg.go.dev/net/http),
which makes it easy to implement the core functionality, and the community has
awesome libraries for building TUIs like [Bubble Tea](https://github.com/charmbracelet/bubbletea) which I used here.\
Finally, I wanted to improve my Go skills by building a real-world application (which is not a web server).

## Installation

You only need to install Go and to clone the repository and run `go build` to get a binary.

For the Go installation instructions, please refer to the [official Go website](https://go.dev/doc/install).


```bash
git clone https://github.com/Yalaouf/gostman
cd gostman
make build
```

This will create a binary named `gostman` (or `gostman.exe` on Windows) in the current directory.
You can also install it directly using `go install`:

```bash
go install github.com/Yalaouf/gostman@latest
```

## Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful, elegant, and fun TUI framework for Go.
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks that plays nicely with the standard library.

That's it! No other dependencies are required.

## Comming Soon
- Save requests
- Collections
- Environment variables
- Authentication methods (OAuth, Bearer, Basic, etc.)
- Syntax highlighting for JSON and XML responses
- Export requests as curl commands
- More to come...
