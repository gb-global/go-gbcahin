# Build Sipe in a stock Go builder container
FROM golang:1.13-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-gbchian
RUN cd /go-gbchian && make gbchain


# Pull Sipe into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-gbchian/build/bin/gbchian /usr/local/bin/

EXPOSE 8545 8546 30312 30312/udp
ENTRYPOINT ["gbchian"]
