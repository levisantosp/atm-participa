fmt:
	find . -type f -name '*.go' -exec go tool gofumpt -w {} +
	find . -type f -name '*.go' -exec go tool golines -m 80 -w {} +
	find . -type f -name '*.go' -exec go tool goimports -w {} +

lint:
	golangci-lint run $(filter-out $@,$(MAKECMDGOALS))

lintf:
	golangci-lint run --fix

build:
	go build ./...
