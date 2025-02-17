FROM golang:1.23-alpine@sha256:b6224579413159e9d4e82e1f160e1e551e947c3fe79db50d1431678374c5dbef AS builder

RUN apk add make
RUN mkdir -p /opt/resource
RUN mkdir -p /opt/code/bin

ADD go.mod /opt/code/
ADD go.sum /opt/code/
WORKDIR /opt/code
RUN go mod download

ADD ./ /opt/code/
RUN make compile

RUN cp /opt/code/bin/* /opt/resource/

FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c
RUN apk upgrade --no-cache \
  && apk add --no-cache ca-certificates
COPY --from=builder /opt/resource /opt/resource/
