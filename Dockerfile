# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .


# Stage 2: Create a minimal runtime image
FROM golang:1.20-alpine

WORKDIR /app

COPY --from=builder /app .

EXPOSE 5000

CMD ["./main"]
