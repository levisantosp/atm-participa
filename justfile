dev:
  just --jobs 2 dev-api dev-cdn

[parallel]
dev-api:
  cd api && just dev

[parallel]
dev-cdn:
  cd cdn && just dev

fmt: fmt-api fmt-cdn fmt-infra

fmt-api:
  cd api && just fmt

fmt-cdn:
  cd cdn && just fmt

fmt-infra:
  cd infra && just fmt

lint *args:
  just lint-api {{ args }}
  just lint-cdn {{ args }}
  just lint-infra {{ args }}

lint-api *args:
  cd api && just lint {{ args }}

lint-cdn *args:
  cd cdn && just lint {{ args }}

lint-infra *args:
  cd infra && just lint {{ args }}

lintf: lintf-api lintf-cdn lintf-infra

lintf-api:
  cd api && just lintf

lintf-cdn:
  cd cdn && just lintf

lintf-infra:
  cd infra && just lintf

create-schema name:
  cd api && just create-schema {{ name }}

generate:
  cd api && just generate

generate-clean:
  cd api && just generate-clean

generate-migrations:
  cd api && just generate-migrations

migrate:
  cd api && just migrate

build: build-api build-cdn build-infra

build-api:
  cd api && just build

build-cdn:
  cd cdn && just build

build-infra:
  cd infra && just build

start:
  just --jobs 2 start-api start-cdn

[parallel]
start-api:
  cd api && just start

[parallel]
start-cdn:
  cd cdn && just start

push: push-api

push-api:
  cd api && just push

test: test-api test-cdn

test-api:
  cd api && just test

test-cdn:
  cd cdn && just test
