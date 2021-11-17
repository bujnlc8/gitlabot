FROM golang:alpine
WORKDIR /bot
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

FROM alpine:3.14.3
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=0  /bot/bot .
CMD ["./bot"]
