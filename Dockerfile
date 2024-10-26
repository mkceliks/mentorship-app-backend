FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG FUNCTION_NAME
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bootstrap handlers/s3/${FUNCTION_NAME}/main.go

FROM amazonlinux:2

RUN yum install -y zip && yum clean all

WORKDIR /app

ARG FUNCTION_NAME
COPY --from=builder /app/bootstrap /app/bootstrap
RUN zip -j /app/${FUNCTION_NAME}_function.zip /app/bootstrap

CMD ["cp", "/app/*.zip", "/app/output/"]
