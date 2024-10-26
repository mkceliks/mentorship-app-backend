FROM amazonlinux:2 as builder

RUN yum install -y \
      golang \
      zip && \
    yum clean all

WORKDIR /app

COPY . .

RUN go mod download

RUN for dir in $(find handlers/s3/* -type d); do \
      if [ -f "$dir/main.go" ]; then \
        function_name=$(basename "$dir"); \
        echo "Building Lambda function: $function_name"; \
        cd "$dir" && \
        GOOS=linux GOARCH=amd64 go build -o bootstrap main.go && \
        zip "/app/handlers/s3/${function_name}/${function_name}_function.zip" bootstrap && \
        cd -; \
      fi; \
    done

FROM amazonlinux:2
WORKDIR /app

COPY --from=builder /app/handlers/s3/upload/upload_function.zip ./handlers/s3/upload/
COPY --from=builder /app/handlers/s3/download/download_function.zip ./handlers/s3/download/
