FROM golang:1.22-alpine AS builder

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /code/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /code/main /usr/local/bin/main

# Copy environment files
COPY .env* ./

EXPOSE 5000

CMD ["/usr/local/bin/main"]