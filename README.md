# Gostman
An API tester written in Go

<p align="center">
  <img src="https://github.com/Yalaouf/gostman/blob/main/logo.png" alt="Gostman Logo"/>
</p>

## Disclaimer

This project is in its early stages and is not yet feature-complete. Use it at your own risk.\
This is a personal project and is not affiliated with any company or organization.
The main goal is to learn and have fun. If it becomes useful to others, that's a bonus!

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

> For the Go installation instructions, please refer to the [official Go website](https://go.dev/doc/install).

To install Gostman, run the following commands:

```bash
git clone https://github.com/Yalaouf/gostman
cd gostman
go build
```

This will create a binary named `gostman` (or `gostman.exe` on Windows) in the current directory.
You can also install it directly using `go install`:

```bash
go install github.com/Yalaouf/gostman@latest
```

## Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful, elegant, and fun TUI framework for Go.
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks that plays nicely with the standard library.
- [Chroma](https://github.com/alecthomas/chroma) - A general purpose syntax highlighter in pure Go

That's it! No other dependencies are required. (not like other languages that
need a million packages for simple tasks, looking at you JavaScript/TypeScript).

## Coming Soon
- Save requests and collections
- Environment variables
- Authentication methods (OAuth, Bearer, Basic, etc.)
- Import Postman and Insomnia files
- Adding more APIs (graphQL, gRPC, etc...)
- And more...
