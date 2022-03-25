#Build
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY main.go ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Deploy
FROM alpine:latest

RUN apk --no-cache add curl ca-certificates

WORKDIR /
COPY --from=builder /main /main

EXPOSE 8080
CMD [ "/main" ]