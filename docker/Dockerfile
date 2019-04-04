FROM golang:alpine AS builder

RUN apk --no-cache add git; \
    go get -u github.com/targodan/piio/piio

FROM alpine

COPY --from=builder /go/bin/piio /usr/bin

WORKDIR /opt/piio

VOLUME /opt/piio/data

ENTRYPOINT ["piio"]
CMD ["serve", "--pi", "data/pi.bin", "--addr", "0.0.0.0:80"]
