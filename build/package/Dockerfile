# Stage 1 - Binary Build
# BUILD_X args should be passed at build time as docker build args
FROM golang:1.14.0-alpine AS builder
ARG BIN_VERSION
ARG GIT_COMMIT
ARG GIT_SHA
ARG GIT_TAG
ARG GIT_DIRTY
ARG GOPROXY
ENV GOPROXY=$GOPROXY
ENV BIN_OUTDIR=./
ENV BIN_NAME=datastore-exporter
RUN apk update && apk add build-base git libressl-dev
WORKDIR /usr/src/datastore-exporter
# install dependencies in separate docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download
# copy application source and build
COPY ./ .
RUN make static

# Stage 2 - Final Image
# The application should be statically linked
FROM alpine:latest
RUN apk update \
	&& apk add --no-cache ca-certificates \
	&& rm -rf /var/cache/apk/* \
    && addgroup datastore-exporter \
	&& adduser -D -H -G datastore-exporter datastore-exporter
COPY --from=builder /usr/src/datastore-exporter/datastore-exporter /usr/bin/datastore-exporter
VOLUME ["/etc/datastore-exporter"]
ENTRYPOINT ["datastore-exporter"]
EXPOSE 5000
USER datastore-exporter
