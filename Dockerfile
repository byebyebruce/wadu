FROM golang:1.22-alpine AS builder

RUN apk add --no-cache build-base git make

ENV GOPROXY='https://goproxy.cn'

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=unknown
RUN make build-musl

FROM alpine

RUN apk --no-cache add tzdata

COPY --from=builder /src/bin /app

WORKDIR /app

# 业务http
EXPOSE 8081

#VOLUME /app/log

ENTRYPOINT ["./wadu","--db=/data/wadu.db","--assets=/data"]
