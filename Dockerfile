FROM alpine:3.4
RUN apk --update add \
    ca-certificates 

COPY ./build/squirrel-linux-amd64 /squirrel

CMD ["/squirrel", "serve"]

