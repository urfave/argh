BENCHTIME ?= 10s
STRINGER := .local/bin/stringer
URFAVE_ARGH_TRACING ?= off

export URFAVE_ARGH_TRACING

.PHONY: all
all: generate test

.PHONY: clean
clean:
	rm -f coverage.out

.PHONY: distclean
distclean: clean
	rm -f $(STRINGER)

.PHONY: generate
generate: $(STRINGER)
	PATH=$(PWD)/.local/bin:$(PATH) go generate ./...

.PHONY: test
test:
	go test -v -coverprofile=coverage.out ./...

.PHONY: bench
bench:
	go test -v -bench . -benchtime $(BENCHTIME) ./...

.PHONY: show-cover
show-cover:
	go tool cover -func=coverage.out

$(STRINGER):
	GOBIN=$(PWD)/.local/bin go install golang.org/x/tools/cmd/stringer@latest
