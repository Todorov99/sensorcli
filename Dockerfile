FROM golang:1.17-alpine as builder

WORKDIR /sensorcli
ENV CGO_ENABLED=0

COPY . /sensorcli

RUN go build -o sensorcli .

FROM alpine:3.15
RUN apk update && apk add bash

COPY ./model.yaml /
COPY ./docker-entrypoint.sh /
COPY --from=builder /sensorcli /sensorcli

ENTRYPOINT ["./docker-entrypoint.sh"]