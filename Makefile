build:
	go build -o gostman main.go

.PHONY: build run clean
.SILENT:

run:
	if [ -f gostman ]; then ./gostman; else go run main.go; fi

clean:
	rm gostman

test:
	go test -coverprofile=coverage.out ./...

coverage:
	go tool cover -html=coverage.out
