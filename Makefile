.PHONY: all
all: test

.PHONY: clean
clean:
	rm -f coverage.out

.PHONY: test
test:
	go test -v -coverprofile=coverage.out ./...

.PHONY: show-cover
show-cover:
	go tool cover -func=coverage.out
