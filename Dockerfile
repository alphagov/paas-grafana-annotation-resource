FROM golang:1.27-alpine@sha256:e9bbdf282b51ac8b34c46e5f31d2d56e7bad60366c35f08d2f295b921b13388b AS builder

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

FROM alpine:3.24.2@sha256:31b6477333eb8257db9e5d7c3a7264fd0467928756f0bbcc27d35bea5d28cdbd
RUN apk upgrade --no-cache \
  && apk add --no-cache ca-certificates
COPY --from=builder /opt/resource /opt/resource/
