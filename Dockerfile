# Use Golang base image for building the binary
FROM golang:1.21 as builder

WORKDIR /app
COPY . .

# Build the Go binary for Lambda with Amazon Linux 2 compatibility
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap ./handlers/s3/upload.go

# Use a lightweight base image for Lambda
FROM amazonlinux:2
WORKDIR /var/task
COPY --from=builder /app/bootstrap .
