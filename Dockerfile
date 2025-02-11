# Step 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifest (go.mod and go.sum) to the container
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app (binary named 'main')
RUN go build -o main .

# Step 2: Create a smaller runtime image to run the Go application
FROM gcr.io/distroless/base

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled Go binary from the previous stage
COPY --from=builder /app/main .

# Expose the port the Go app runs on
EXPOSE 8080

# Start the Go application
CMD ["/root/main"]
