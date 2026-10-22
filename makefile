.PHONY: dev fmt lint lintf create-schema generate generate-clean build start push test
.PHONY: dev\:api fmt\:api lint\:api lintf\:api build\:api start\:api push\:api test\:api
.PHONY: dev\:cdn fmt\:cdn lint\:cdn lintf\:cdn build\:cdn start\:cdn test\:cdn
.PHONY: fmt\:infra lint\:infra lintf\:infra build\:infra

dev:
	$(MAKE) -j2 'dev:api' 'dev:cdn'

dev\:api:
	$(MAKE) -C api dev

dev\:cdn:
	$(MAKE) -C cdn dev

fmt: fmt\:api fmt\:cdn fmt\:infra

fmt\:api:
	$(MAKE) -C api fmt

fmt\:cdn:
	$(MAKE) -C cdn fmt

fmt\:infra:
	$(MAKE) -C infra fmt

lint: lint\:api lint\:cdn lint\:infra

lint\:api:
	$(MAKE) -C api lint

lint\:cdn:
	$(MAKE) -C cdn lint

lint\:infra:
	$(MAKE) -C infra lint

lintf: lintf\:api lintf\:cdn lintf\:infra

lintf\:api:
	$(MAKE) -C api lintf

lintf\:cdn:
	$(MAKE) -C cdn lintf

lintf\:infra:
	$(MAKE) -C infra lintf

create-schema:
	$(MAKE) -C api create-schema name=$(name)

generate:
	$(MAKE) -C api generate

generate-clean:
	$(MAKE) -C api generate-clean

build: build\:api build\:cdn build\:infra

build\:api:
	$(MAKE) -C api build

build\:cdn:
	$(MAKE) -C cdn build

build\:infra:
	$(MAKE) -C infra build

start:
	$(MAKE) -j2 'start:api' 'start:cdn'

start\:api:
	$(MAKE) -C api start

start\:cdn:
	$(MAKE) -C cdn start

push: push\:api

push\:api:
	$(MAKE) -C api push

test: test\:api test\:cdn

test\:api:
	$(MAKE) -C api test

test\:cdn:
	$(MAKE) -C cdn test
