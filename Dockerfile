FROM golang:1.13.8-alpine3.11 as builder

WORKDIR /sensor

COPY . /sensor

RUN go build -o sensorcli .

FROM alpine:3.11
RUN apk update && apk add bash

COPY ./model.yaml /
COPY ./docker-entrypoint.sh /
COPY --from=builder /sensor /sensor

ENTRYPOINT ["./docker-entrypoint.sh"]