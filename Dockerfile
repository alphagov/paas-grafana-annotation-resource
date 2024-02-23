# FROM golang:1.14-alpine@sha256:ef409ff24dd3d79ec313efe88153d703fee8b80a522d294bb7908216dc7aa168 AS builder
FROM golang:1.22-alpine@sha256:3325c5593767d8f1fd580e224707ca5e847a1679470a027aaa3c71711ce16613 AS builder

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

FROM alpine:3.19@sha256:6457d53fb065d6f250e1504b9bc42d5b6c65941d57532c072d929dd0628977d0
RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /opt/resource /opt/resource/
