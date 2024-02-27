FROM golang:1.20.5 as builder
WORKDIR /app/codingsilver
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go 
RUN go build -o authn-webhook ./main.go 

FROM alpine:3.19.0
RUN apk add gcompat
WORKDIR /root
COPY cert.pem cert.pem
COPY key.pem key.pem
COPY --from=builder /app/codingsilver/authn-webhook /usr/local/bin/authn-webhook

