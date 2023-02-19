FROM golang:1.20.0-alpine3.17 AS builder

WORKDIR $GOPATH/src/http/
COPY go/* ./
COPY cert/* ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/proxy

FROM scratch
COPY --from=builder /go/bin/proxy /go/bin/proxy
EXPOSE 8181

CMD ["/go/bin/proxy"]