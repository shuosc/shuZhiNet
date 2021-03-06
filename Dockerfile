FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/shuZhiNet
WORKDIR /go/src/shuZhiNet
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/shuZhiNet/shuZhiNet /shuZhiNet
WORKDIR /
RUN export PORT=8000
CMD ./shuZhiNet
EXPOSE 8000