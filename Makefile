GO ?= go
MAP ?= dataset/skyrim.txt
COUNT ?= 2

.PHONY: run
run:
	@$(GO) run ./cmd/invasion -c $(COUNT) -f $(MAP)

.PHONY: test
test:
	@$(GO) test -count=1 ./...
