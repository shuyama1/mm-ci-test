FROM golang:1.18-alpine

WORKDIR /tmp
COPY . ./

RUN go build

ENTRYPOINT ["/tmp/membership_checker"]