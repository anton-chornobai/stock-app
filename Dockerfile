FROM golang:1.25-alpine AS builder

# refers to the directory inside the container where the project lives
WORKDIR /app

ADD go.mod .

COPY . .

RUN go build -o myapp ./cmd/server/main.go

FROM alpine 

EXPOSE 80 8080

WORKDIR /app

COPY --from=builder /app/myapp /app/myapp

CMD ["./myapp"]