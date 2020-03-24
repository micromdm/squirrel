FROM golang:latest as builder

WORKDIR /go/src/github.com/micromdm/squirrel/

ENV CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux

COPY . .

RUN make deps
RUN make


FROM alpine:3.11.5
RUN apk --update add \
    ca-certificates 

COPY --from=builder /go/src/github.com/micromdm/squirrel/build/linux/squirrel /usr/bin/

CMD ["/squirrel", "serve"]

