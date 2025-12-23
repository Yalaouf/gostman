# gostman
An API tester written in Go

## Why

I was sick of using bloated electron based GUI apps for simple API testing. I wanted a
lightweight tool that could do the job without consuming too many resources.

For now, gostman is a command line tool, but I might add a GUI in the future.

## Why Go

Go is fast, compiles to a single binary, and can be compiled for multiple platforms easily.
And I like Go.

## Installation

You only need to clone the repository and run `go build` to get a binary.

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
