# Step 1: Use Amazon Linux 2 as the base image
FROM amazonlinux:2 as builder

# Step 2: Install Go and other dependencies
RUN yum install -y golang zip

# Step 3: Set the working directory to the handlers/s3 directory inside the container
WORKDIR /go/src/mentorship-app/handlers/s3

# Step 4: Copy the Go module files
COPY go.mod go.sum /go/src/mentorship-app/

# Step 5: Download Go modules
RUN cd /go/src/mentorship-app && go mod download

# Step 6: Copy the rest of the project files
COPY . /go/src/mentorship-app/

# Step 7: Build the Go binary for Linux (compatible with Lambda)
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap upload.go

# Step 8: Package the binary into a zip file
RUN zip function.zip bootstrap

# Step 9: Output the function.zip for deployment
CMD ["cat", "function.zip"]
