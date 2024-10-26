FROM amazonlinux:2

RUN yum install -y \
    zip \
    golang \
    && yum clean all

WORKDIR /app

COPY main.go .

RUN GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

RUN zip function.zip bootstrap

CMD ["cp", "function.zip", "/app/function.zip"]
