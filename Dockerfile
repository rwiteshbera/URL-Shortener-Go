FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the binary
RUN go build -o main .


EXPOSE 5000

# Start the application
CMD ["./main"]