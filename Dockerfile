FROM golang:alpine as builder

RUN mkdir /build

ADD . /build

WORKDIR /build

RUN go build -o main.go



