# Builder Image
FROM golang:1.21 as builder

WORKDIR /abs-tract
COPY . .
RUN go mod download
RUN go build -v -o bin/abs-tract

# Ditribution Image
FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder /abs-tract/bin/abs-tract /abs-tract

EXPOSE 5555

ENTRYPOINT ["/abs-tract"]
