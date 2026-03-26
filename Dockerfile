FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY app ./app
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /embedded-lab-api ./cmd

FROM alpine:3.22

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /embedded-lab-api /usr/local/bin/embedded-lab-api

USER appuser

EXPOSE 8080

CMD ["embedded-lab-api"]
