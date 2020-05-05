export GOPATH=$(PWD)

build:
	go build -o cvetracker src/cvetracker/main/main.go
