FROM golang:1.14.1-alpine3.11 as build

WORKDIR /go/proxy

COPY . .

RUN go install github.com/im7mortal/proxySearchEngine/cmd/proxy

CMD ["proxy"]