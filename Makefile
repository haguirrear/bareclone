
BIN=bareclone

.PHONY: build clean
build:
	go build -o $(BIN) main.go

clean:
	rm -f $(BIN)
