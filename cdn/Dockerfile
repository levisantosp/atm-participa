FROM golang:1.27.0-alpine AS build

WORKDIR /app

RUN apk add --no-cache make

COPY go.work go.work.sum /app/
COPY api/go.mod api/go.sum /app/api/
COPY cdn/go.mod cdn/go.sum /app/cdn/
COPY infra/go.mod infra/go.sum /app/infra/

RUN cd cdn && go mod download

COPY . .

RUN make build:cdn

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/cdn/bin/cdn ./bin/cdn

ENTRYPOINT [ "./bin/cdn" ]
