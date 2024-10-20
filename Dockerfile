# Start with a Golang base image for building the Lambda binary
FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary for Linux (for AWS Lambda)
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap -buildvcs=false

# Use a minimal base image for Lambda runtime (Amazon Linux 2)
FROM amazonlinux:2

# Install zip for packaging the Lambda function
RUN yum install -y zip

# Copy the Go binary from the builder stage
COPY --from=builder /app/bootstrap /var/task/

# Zip the function for Lambda deployment
WORKDIR /var/task
RUN zip function.zip bootstrap

# The Lambda runtime uses the file "bootstrap" as the entrypoint
