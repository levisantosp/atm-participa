dev:
  just --jobs 2 dev-api dev-cdn

[parallel]
dev-api:
  cd apps/api && just dev

[parallel]
dev-cdn:
  cd apps/cdn && just dev

fmt: fmt-api fmt-cdn fmt-infra

fmt-api:
  cd apps/api && just fmt

fmt-cdn:
  cd apps/cdn && just fmt

fmt-infra:
  cd pkg/infra && just fmt

lint *args:
  just lint-api {{ args }}
  just lint-cdn {{ args }}
  just lint-infra {{ args }}

lint-api *args:
  cd apps/api && just lint {{ args }}

lint-cdn *args:
  cd apps/cdn && just lint {{ args }}

lint-infra *args:
  cd pkg/infra && just lint {{ args }}

lintf: lintf-api lintf-cdn lintf-infra

lintf-api:
  cd apps/api && just lintf

lintf-cdn:
  cd apps/cdn && just lintf

lintf-infra:
  cd pkg/infra && just lintf

create-schema name:
  cd apps/api && just create-schema {{ name }}

generate:
  cd apps/api && just generate

generate-clean:
  cd apps/api && just generate-clean

generate-migrations:
  cd apps/api && just generate-migrations

migrate:
  cd apps/api && just migrate

build: build-api build-cdn build-infra

build-api:
  cd apps/api && just build

build-cdn:
  cd apps/cdn && just build

build-infra:
  cd pkg/infra && just build

start:
  just --jobs 2 start-api start-cdn

[parallel]
start-api:
  cd apps/api && just start

[parallel]
start-cdn:
  cd apps/cdn && just start

push: push-api

push-api:
  cd apps/api && just push

test: test-api test-cdn

test-api:
  cd apps/api && just test

test-cdn:
  cd apps/cdn && just test
