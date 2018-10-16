#BUILD BINARY
FROM golang:1.10.3-alpine as builder

ENV REPO_PATH ${GOPATH}/src/github.com/alexrios/challenge-api
ENV BIN_NAME main

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep

# Dependencies to build
RUN apk add --no-cache gcc musl-dev git && chmod +x /usr/bin/dep

RUN mkdir -p $REPO_PATH
COPY ./  $REPO_PATH
WORKDIR $REPO_PATH

#-vendor-only  populate vendor/ from Gopkg.lock without updating it first (default: false)
#2 cp -Rp vendor/* /go/src/ && \
RUN dep ensure --vendor-only && \
    go get ./... && \
    go build -a -ldflags "-s -w" -o $BIN_NAME

#BUILD IMAGE
FROM alpine:3.8
ENV REPO_PATH /go/src/github.com/alexrios/challenge-api
ENV BIN_NAME main
RUN apk add --no-cache ca-certificates
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder $REPO_PATH/$BIN_NAME .
EXPOSE 8080
CMD ["./main"]