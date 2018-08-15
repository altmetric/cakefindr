FROM golang:1.9-alpine as builder
WORKDIR /go/src/github.com/altmetric/cakefindr-keyserver-3/
RUN apk add --update git build-base
COPY . ./
RUN go build -o keyserver-3 ./main.go

FROM alpine:3.7
WORKDIR /app/
COPY --from=builder /go/src/github.com/altmetric/cakefindr-keyserver-3/keyserver-3 .
CMD ["./keyserver-3"]
