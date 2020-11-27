FROM golang:1.14-alpine AS builder

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

FROM alpine:3.12
RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /opt/resource /opt/resource/
