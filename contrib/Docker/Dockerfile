# Build Stage
FROM golang:1.21.1-alpine3.18 AS build-stage

LABEL app="build-ip2location-pfsense"
LABEL REPO="https://github.com/jpmchia/IP2Location-pfSense"

ENV PROJPATH=/go/src/github.com/jpmchia/IP2Location-pfSense/backend

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/jpmchia/IP2Location-pfSense
WORKDIR /go/src/github.com/jpmchia/IP2Location-pfSense

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
FROM jpmchia/alpine-base:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/jpmchia/IP2Location-pfSense"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/IP2Location-pfSense/bin

WORKDIR /opt/IP2LOCATION/bin

COPY --from=build-stage /go/src/github.com/jpmchia/IP2Location-pfSense/bin/ip2location-pfsense /opt/IP2Location-pfSense/bin/
RUN chmod +x /opt/IP2Location-pfSense/bin/ip2location-pfsense

# Create appuser
RUN adduser -D -g '' ip2location
USER ip2location

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/IP2Location-pfSense/bin/ip2location-pfsense", "service"]
