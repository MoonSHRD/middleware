FROM golang:1.11-stretch

RUN apt-get update && apt-get install -y build-essential cmake git

WORKDIR /go/src/middleware
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8000

