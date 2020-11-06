FROM golang:1.15.3

ENV GO111MODULE=off

RUN go get -u github.com/underkode/secret-santa

RUN mkdir -p /data

VOLUME ["/data"]

CMD ["secret-santa", "-d", "/data/", "-t", "$SECRET_SANTA_TELEGRAM_TOKEN"]