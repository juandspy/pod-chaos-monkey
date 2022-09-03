# Build
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /pod-chaos-monkey

# Deploy
FROM alpine:3.16
COPY --from=builder /pod-chaos-monkey .
ENTRYPOINT [ "./pod-chaos-monkey" ]