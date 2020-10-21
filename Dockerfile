# Build image
FROM golang:1.15-alpine as builder

ENV CGO_ENABLED=0

WORKDIR /exporter

COPY . /exporter

RUN go build -o /go/bin/exporter main.go etherscan.go parity.go

# Executable image
FROM alpine

WORKDIR /

COPY --from=builder /go/bin/exporter /exporter

CMD ["/exporter"]
