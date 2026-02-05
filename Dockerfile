FROM golang:1.25-alpine@sha256:f4622e3bed9b03190609db905ac4b02bba2368ba7e62a6ad4ac6868d2818d314 AS builder

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

FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659
RUN apk upgrade --no-cache \
  && apk add --no-cache ca-certificates
COPY --from=builder /opt/resource /opt/resource/
