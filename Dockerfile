#################
# Base image
#################
FROM alpine:3.12 as echo-redis-base

USER root

RUN addgroup -g 10001 echo-redis && \
    adduser --disabled-password --system --gecos "" --home "/home/echo-redis" --shell "/sbin/nologin" --uid 10001 echo-redis && \
    mkdir -p "/home/echo-redis" && \
    chown echo-redis:0 /home/echo-redis && \
    chmod g=u /home/echo-redis && \
    chmod g=u /etc/passwd

ENV USER=echo-redis
USER 10001
WORKDIR /home/echo-redis

#################
# Builder image
#################
FROM golang:1.16-alpine AS echo-redis-builder
RUN apk add --update --no-cache alpine-sdk
WORKDIR /app
COPY . .
RUN make build

#################
# Final image
#################
FROM echo-redis-base

COPY --from=echo-redis-builder /app/bin/echo-redis /usr/local/bin

# Command to run the executable
ENTRYPOINT ["echo-redis"]
