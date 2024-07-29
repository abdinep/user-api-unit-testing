FROM golang:1.22-alpine AS builder

WORKDIR /mini_project

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /mini_project

COPY --from=builder /mini_project/main .

EXPOSE 8080

CMD ["./main"]