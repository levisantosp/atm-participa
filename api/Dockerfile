FROM golang:1.27.3-alpine AS build

WORKDIR /app

RUN apk add --no-cache just

COPY go.work go.work.sum /app/
COPY api/go.mod api/go.sum /app/api/
COPY cdn/go.mod cdn/go.sum /app/cdn/
COPY infra/go.mod infra/go.sum /app/infra/

RUN cd api && go mod download

COPY . .

RUN just build-api

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/api/bin/api ./bin/api

ENTRYPOINT [ "./bin/api" ]
