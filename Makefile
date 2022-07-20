.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -v ./...

SRC := $(shell find -name main.go)

.PHONY: build
build: $(SRC)
	for main in $(SRC) ; do \
		go build $$main ; \
	done
	rm main