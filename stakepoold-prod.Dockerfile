FROM golang:1.12.1-alpine3.9 as builder

RUN apk add git gcc g++ musl-dev curl --update --no-cache
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/EXCCoin/exccstakepool
COPY . .

RUN dep ensure

WORKDIR /go/src/github.com/EXCCoin/exccstakepool/backend/stakepoold
RUN go build -ldflags='-s -w -X main.appBuild=alpine3.9 -extldflags "-static"' .


FROM alpine:3.9

WORKDIR /app
COPY --from=builder /go/src/github.com/EXCCoin/exccstakepool/backend/stakepoold/stakepoold .

EXPOSE 9113
ENV DATA_DIR=/data
ENV CONFIG_FILE=/app/stakepoold.conf
CMD ["sh", "-c", "/app/stakepoold --appdata=${DATA_DIR} --configfile=${CONFIG_FILE}"]
