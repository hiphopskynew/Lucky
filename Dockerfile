FROM golang:1.10.0 as builder
# Please set APP_NAME
ARG APP_NAME=lucky

WORKDIR /go/src/${APP_NAME}/

COPY . /go/src/${APP_NAME}

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
# Please set APP_NAME
ARG APP_NAME=lucky

# install curl support docker healthcheck
RUN apk --no-cache add curl

RUN apk --no-cache add ca-certificates

ENV RUN_MODE=production

WORKDIR /root/

ENV RUN_MODE=production

COPY --from=builder /go/src/${APP_NAME}/configs ./configs

COPY --from=builder /go/src/${APP_NAME}/initialize/sql ./initialize/sql

COPY --from=builder /go/src/${APP_NAME}/app .

CMD ["./app"]

HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://localhost:8080/ping -H "X-Activity:Service Health-check" || exit 1 