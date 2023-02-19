FROM golang:1.20.1 AS builder

WORKDIR $GOPATH/src/proxy/
COPY go/* ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" -o /go/bin/proxy

 FROM scratch
 COPY --from=builder /go/bin/proxy /go/bin/proxy
 COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/go/bin/proxy"]
