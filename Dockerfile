FROM amazonlinux:2 as builder

RUN yum install -y golang zip

WORKDIR /mentorship-app-backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p build && \
    for dir in handlers/s3/*; do \
        if [ -d "$dir" ]; then \
            handler_name=$(basename "$dir"); \
            cd "/mentorship-app-backend/$dir" && \
            echo "Current directory: $(pwd)" && \
            echo "Files in $(pwd):" && ls -l && \
            if [ -f "main.go" ]; then \
                GOOS=linux GOARCH=amd64 go build -o bootstrap main.go && \
                zip "/mentorship-app-backend/build/${handler_name}.zip" bootstrap; \
            else \
                echo "Error: main.go not found in $dir"; exit 1; \
            fi; \
        fi; \
    done

CMD ["ls", "-R", "/mentorship-app-backend/build"]
