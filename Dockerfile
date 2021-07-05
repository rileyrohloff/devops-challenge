FROM golang:1.15 as build

WORKDIR /go/src/app

COPY . /go/src/app

RUN CGO_ENABLED=0 go build .

FROM debian:buster-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build go/src/app/devops-challenge .
COPY --from=build go/src/app/input.yaml .

ENTRYPOINT ["./devops-challenge"]