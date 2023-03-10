FROM golang:1.19 AS builder

WORKDIR /src
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0

COPY . .

RUN go get ./...

RUN go build -ldflags="-s -w" -o ga-event-tracker ./cmd/

FROM alpine:latest

COPY --from=builder /src/ga-event-tracker /

CMD ["/ga-event-tracker"]
