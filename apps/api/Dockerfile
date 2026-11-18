FROM golang:1.27.4-alpine AS build

WORKDIR /app

RUN apk add --no-cache just

COPY go.work go.work.sum /app/
COPY apps/api/go.mod apps/api/go.sum /app/apps/api/
COPY apps/cdn/go.mod apps/cdn/go.sum /app/apps/cdn/
COPY pkg/infra/go.mod pkg/infra/go.sum /app/pkg/infra/

RUN cd apps/api && go mod download

COPY . .

RUN just build-api

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/apps/api/bin/api ./bin/api

ENTRYPOINT [ "./bin/api" ]
