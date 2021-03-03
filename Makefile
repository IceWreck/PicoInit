run:
	go run .

build:
	go build -ldflags "-s -w" -o ./bin/PicoInit .