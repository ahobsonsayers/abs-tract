# Builder Image
FROM golang:1.21 as builder

WORKDIR /abs-goodreads
COPY . .
RUN go mod download
RUN go build -v -o bin/abs-goodreads

# Ditribution Image
FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder /abs-goodreads/bin/abs-goodreads /abs-goodreads

EXPOSE 5555

ENTRYPOINT ["/abs-goodreads"]
