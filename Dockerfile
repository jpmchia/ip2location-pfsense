# Build Stage
#FROM golang:1.21.1-alpine3.18 AS build-stage
FROM golang:alpine AS build-stage

LABEL app="build-ip2location-pfsense"
LABEL REPO="https://github.com/jpmchia/ip2location-pfsense"

ENV PROJPATH=/go/src/github.com/jpmchia/ip2ocation-pfsense/backend

# Because of https://github.com/docker/docker/issues/14914
#ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/jpmchia/ip2location-pfsense
WORKDIR /go/src/github.com/jpmchia/ip2location-pfsense

RUN apk add --no-cache --update curl \
    dumb-init \
    bash \
    grep \
    sed \
    jq \
    ca-certificates \
    openssl \
    git \
    make \
    gcc \
    musl-dev \
    && rm -rf /var/cache/apk/*

RUN make build-alpine

# Final Stage
FROM alpine:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/jpmchia/ip2location-pfsense"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/ip2location-pfsense/bin

WORKDIR /opt/ip2location/bin

COPY --from=build-stage /go/src/github.com/jpmchia/ip2location-pfsense/backend/bin/ip2location-pfsense /opt/ip2location-pfsense/bin/
RUN chmod +x /opt/ip2location-pfsense/bin/ip2location-pfsense

# Create appuser
RUN adduser -D -g '' ip2location
USER ip2location

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/ip2location-pfsense/bin/ip2location-pfsense", "service"]
