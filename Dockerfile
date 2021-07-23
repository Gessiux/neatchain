# Build neatchain in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /neatchain
RUN cd /neatchain && make neatchain

# Pull Neatchain into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /neatchain/bin/neatchain /usr/local/bin/

EXPOSE 8556 8555 8554 8553 8552 8551 8550 8550/udp
ENTRYPOINT ["neatchain"]
