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

generate-migrations:
	GOWORK=off atlas migrate diff --env local

migrate:
	@if [ -z "$$DATABASE_URL" ] && [ -f .env ]; then set -a; . ./.env; set +a; fi; \
	case "$$DATABASE_URL" in *sslmode=*) ;; *\?*) DATABASE_URL="$$DATABASE_URL&sslmode=disable" ;; *) DATABASE_URL="$$DATABASE_URL?sslmode=disable" ;; esac; \
	case "$$DATABASE_URL" in *search_path=*) ;; *\?*) DATABASE_URL="$$DATABASE_URL&search_path=public" ;; *) DATABASE_URL="$$DATABASE_URL?search_path=public" ;; esac; \
	export DATABASE_URL; \
	atlas migrate apply --env local --revisions-schema atlas_schema_revisions

build:
	rm -rf bin && go build -o bin/api ./cmd/api

start:
	bin/api

push:
	go run ./cmd/push

test:
	go test -p 1 ./internal/routes/...
