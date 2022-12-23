FROM golang:1.19.4-alpine3.17 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY ./expenses ./expenses
RUN go build -o /app/server .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 2565
CMD ["/app/server"]
