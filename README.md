# Gostman
An API tester written in Go

<p align="center">
  <img width="500" height="500" src="https://github.com/Yalaouf/gostman/blob/main/logo.png" alt="Gostman Logo"/>
</p>

## Disclaimer

This project is in its early stages and is not yet feature-complete.\
I wrote this project just to have fun and get rid of the bloated Electron equivalent.\
If it becomes useful to others, that's a bonus!

## Why

I was sick of using bloated electron based apps for simple API testing. I wanted a
lightweight tool that could do the job without consuming too many resources.\
There are already some good alternatives like [HTTPie](https://github.com/httpie/cli),
[curl](https://github.com/curl/curl), or [xh](https://github.com/ducaale/xh) but
I wanted something with a TUI (not GUI to stay in terminal) that was simple and easy to use.\
Also I wanted no compromise for simplicity over features, so I decided to write my own for fun.

## Why Go

[Go](https://go.dev/) is fast, compiles to a single binary, and can be compiled for multiple platforms easily.
Also Go provides a good [standard library for HTTP requests](https://pkg.go.dev/net/http),
which makes it easy to implement the core functionality, and the community has
awesome libraries for building TUIs like [Bubble Tea](https://github.com/charmbracelet/bubbletea) which I used here.\

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
go install github.com/Yalaouf/gostman/gostman@v1.0.1
```

## Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful, elegant, and fun TUI framework for Go.
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks that plays nicely with the standard library.
- [Chroma](https://github.com/alecthomas/chroma) - A general purpose syntax highlighter in pure Go
- [WordWrap](https://github.com/muesli/reflow) - A collection of ANSI-aware methods and io.Writers helping you to transform blocks of text.
- [Uuid](https://www.github.com/google/uuid) - The uuid package generates and inspects UUIDs based on RFC 9562 and DCE 1.1: Authentication and Security Services.

## Demo

https://github.com/user-attachments/assets/8cd3bf7c-4537-4a01-9b5f-a8bffc83b306
## Coming Soon
- Unit tests on TUI
- Environment variables
- Authentication methods (OAuth, Bearer, Basic, etc.)
- Import Postman and Insomnia files
- Adding more protocols (graphQL, gRPC, etc...)
- More themes (only `catppuccin` for now)
- And more...
