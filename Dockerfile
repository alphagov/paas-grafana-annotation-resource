FROM golang:1.23-alpine@sha256:13aaa4b92fd4dc81683816b4b62041442e9f685deeb848897ce78c5e2fb03af7 AS builder

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

FROM alpine:3.21.0@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45
RUN apk upgrade --no-cache \
  && apk add --no-cache ca-certificates
COPY --from=builder /opt/resource /opt/resource/
