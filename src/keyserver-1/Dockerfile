FROM golang:1.9-alpine as builder
WORKDIR /go/src/github.com/altmetric/cakefindr-keyserver-1/
RUN apk add --update git build-base
COPY . ./
RUN go build -o keyserver-1 ./main.go

FROM alpine:3.7
WORKDIR /app/
COPY --from=builder /go/src/github.com/altmetric/cakefindr-keyserver-1/keyserver-1 .
CMD ["./keyserver-1"]
