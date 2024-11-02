FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG FUNCTION_NAME

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bootstrap handlers/*/${FUNCTION_NAME}/main.go

RUN echo "Contents of /app after Go build:" && ls -la /app

FROM amazonlinux:2

RUN yum install -y zip && yum clean all

WORKDIR /app

ARG FUNCTION_NAME

COPY --from=builder /app/bootstrap /app/bootstrap

RUN zip -j /app/${FUNCTION_NAME}_function.zip /app/bootstrap

RUN echo "Contents of /app after zipping:" && ls -la /app

RUN mkdir -p /app/output
RUN cp /app/${FUNCTION_NAME}_function.zip /app/output/${FUNCTION_NAME}_function.zip

RUN echo "Contents of /app/output after moving zip file:" && ls -la /app/output
