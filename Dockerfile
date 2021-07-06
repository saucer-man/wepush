FROM golang:1.16.4-alpine3.13 AS builder
WORKDIR /src/wepush
ENV GOPROXY="https://goproxy.cn"
COPY src .
RUN go build -ldflags '-w -s' -o wepush

FROM alpine:3.13
WORKDIR /src/wepush
COPY --from=builder /src/wepush/wepush .
CMD ["./wepush"]