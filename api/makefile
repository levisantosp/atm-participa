dev:
	go run ./cmd/api

fmt:
	find . -path './pg' -prune -o -type f -name '*.go' -exec go tool gofumpt -w {} +
	find . -path './pg' -prune -o -type f -name '*.go' -exec go tool golines -m 80 -w {} +
	find . -path './pg' -prune -o -type f -name '*.go' -exec go tool goimports -w {} +

lint:
	golangci-lint run $(filter-out $@,$(MAKECMDGOALS))

lintf:
	golangci-lint run --fix

create-schema:
ifndef name
	$(error name is required. Usage: make create-schema name=SchemaName)
endif
	go tool ent new $(name)

generate:
	go generate ./ent/...

generate-clean:
	rm -rf ent/generated && go generate ./ent/...

build:
	rm -rf bin && go build -o bin/api ./cmd/api

start:
	bin/api

push:
	go run ./cmd/push

test:
	go test -p 1 ./internal/routes/...
