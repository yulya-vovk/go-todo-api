FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-api main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/todo-api .
COPY --from=builder /app/data ./data
EXPOSE 8080
CMD ["./todo-api"]