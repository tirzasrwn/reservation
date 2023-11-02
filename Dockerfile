FROM golang:1.20.5-alpine3.18 AS builder
LABEL stage=builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
USER 0:0
RUN go build -o ./reservation ./cmd/web
RUN chmod -R 777 ./reservation

FROM alpine:3.18 AS production
WORKDIR /app
COPY --from=builder /app ./
CMD ["./reservation"]
EXPOSE 4545 
