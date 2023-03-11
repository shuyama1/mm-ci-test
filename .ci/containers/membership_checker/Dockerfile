FROM golang:1.18-alpine

WORKDIR /tmp

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /membership_checker

EXPOSE 8080

ENTRYPOINT [ "/membership_checker" ]