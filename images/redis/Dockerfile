# check=skip=InvalidDefaultArgInFrom
ARG GO_IMAGE
ARG BASE_IMAGE

FROM ghcr.io/databasifylabs/engine-sdk:1.0.0@sha256:cf066886a156193196e9d043e3f75e9084928948e40bdd950b2306f2b340d901 AS sdk

FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS redis-shake
ARG TARGETOS
ARG TARGETARCH
ARG REDIS_SHAKE_VERSION
ARG REDIS_SHAKE_SHA256
WORKDIR /src
RUN set -eux; \
    curl -fsSL -o /tmp/redis-shake.tar.gz "https://github.com/tair-opensource/RedisShake/archive/refs/tags/v${REDIS_SHAKE_VERSION}.tar.gz"; \
    echo "${REDIS_SHAKE_SHA256}  /tmp/redis-shake.tar.gz" | sha256sum -c -; \
    tar -xzf /tmp/redis-shake.tar.gz --strip-components=1; \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath \
        -ldflags "-s -w -X main.Version=v${REDIS_SHAKE_VERSION} -X main.GitCommit=source" \
        -o /out/redis-shake ./cmd/redis-shake

FROM ${BASE_IMAGE}
ARG RELEASE
ARG SDK_VERSION
ARG LINE
ARG PACKAGES
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends $PACKAGES; \
    rm -rf /var/lib/apt/lists/*
COPY --from=redis-shake /out/redis-shake /usr/local/bin/redis-shake
RUN set -eux; \
    groupmod -g 64999 redis; \
    usermod -u 64999 -g 64999 redis; \
    find / -xdev \( -uid 999 -o -gid 999 \) -exec chown -h 64999:64999 {} +
COPY --from=sdk /databasify /usr/local/bin/databasify
LABEL org.opencontainers.image.source="https://github.com/databasifylabs/engine-redis" \
      io.databasify.engine="redis" \
      io.databasify.release="${RELEASE}" \
      io.databasify.line="${LINE}" \
      io.databasify.sdk="${SDK_VERSION}"
ENV DATABASIFY_ENGINE=redis \
    DATABASIFY_RELEASE=${RELEASE} \
    DATABASIFY_LINE=${LINE}
STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/local/bin/databasify", "serve"]
CMD []
