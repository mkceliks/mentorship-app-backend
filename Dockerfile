FROM amazonlinux:2 as builder

RUN yum install -y golang zip

WORKDIR /go/src/mentorship-app/handlers/s3

COPY go.mod go.sum /go/src/mentorship-app/

RUN cd /go/src/mentorship-app && go mod download

COPY . /go/src/mentorship-app/

RUN GOOS=linux GOARCH=amd64 go build -o bootstrap upload.go

RUN zip function.zip bootstrap

CMD ["cat", "function.zip"]
