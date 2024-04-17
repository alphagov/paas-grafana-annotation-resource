FROM golang:1.22-alpine@sha256:cdc86d9f363e8786845bea2040312b4efa321b828acdeb26f393faa864d887b0 AS builder

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

FROM alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b
RUN apk upgrade --no-cache \
  && apk add --no-cache ca-certificates
COPY --from=builder /opt/resource /opt/resource/
