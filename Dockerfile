FROM golang:1.13.3 AS builder

COPY . /go/src/github.com/Traestan/shorturl/
WORKDIR /go/src/github.com/Traestan/shorturl/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/shorturl

FROM alpine:3.10
# Copy our static executable.
COPY --from=builder /go/bin/shorturl /go/bin/shorturl
# Run the hello binary.
ENTRYPOINT ["/go/bin/shorturl"]