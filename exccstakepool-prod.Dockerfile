FROM golang:1.12.1-alpine3.9 as builder

RUN apk add git gcc g++ musl-dev curl --update --no-cache
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/EXCCoin/exccstakepool
COPY . .

RUN dep ensure
RUN go build -ldflags='-s -w -X main.appBuild=alpine3.9 -extldflags "-static"' .


FROM alpine:3.9

RUN apk update && apk upgrade && apk add --no-cache ca-certificates && update-ca-certificates 2>/dev/null || true

WORKDIR /app
COPY --from=builder /go/src/github.com/EXCCoin/exccstakepool/exccstakepool .
COPY ./views ./views
COPY ./public ./public

EXPOSE 9113
ENV DATA_DIR=/data
ENV LOG_DIR=/log
ENV CONFIG_FILE=/app/exccstakepool.conf
CMD ["sh", "-c", "/app/exccstakepool --datadir=${DATA_DIR} --logdir=${LOG_DIR} --configfile=${CONFIG_FILE}"]
