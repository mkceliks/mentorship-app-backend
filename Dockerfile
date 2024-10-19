# Start with a Golang base image
FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary for Linux
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap -buildvcs=false

# Use a minimal base image for Lambda
FROM amazonlinux:2

# Copy the binary from the builder stage
COPY --from=builder /app/bootstrap /var/task/

# The Lambda runtime already uses the file "bootstrap" as the entrypoint
