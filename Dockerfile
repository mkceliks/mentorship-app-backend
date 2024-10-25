FROM amazonlinux:2 as builder

RUN yum install -y golang zip

WORKDIR /go/src/mentorship-app-backend
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir -p build && \
    for dir in handlers/s3/*; do \
        if [ -d "$dir" ]; then \
            handler_name=$(basename "$dir"); \
            cd "/go/src/mentorship-app-backend/$dir" && \
            GOOS=linux GOARCH=amd64 go build -o bootstrap main.go && \
            zip "/go/src/mentorship-app-backend/build/${handler_name}.zip" bootstrap; \
        fi; \
    done

CMD ["ls", "-R", "/go/src/mentorship-app-backend/build"]
