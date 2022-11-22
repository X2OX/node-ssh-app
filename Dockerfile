FROM golang:alpine AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /build/server github.com/X2OX/node-ssh-app/cmd/server


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/server server

CMD ["/app/server"]