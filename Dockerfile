# Use Amazon Linux 2 as the base image
FROM amazonlinux:2

# Install Go and zip
RUN yum install -y golang zip

# Set the working directory
WORKDIR /go/src/mentorship-app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go binary for AWS Lambda
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap ./handlers/s3/upload.go

# Zip the Lambda function for deployment
RUN zip function.zip bootstrap

# The final command (optional, to inspect the zip file)
CMD ["cat", "function.zip"]
