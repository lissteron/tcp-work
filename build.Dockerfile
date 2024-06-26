FROM golang:1.22.3-alpine

ENV WORKDIR=/src
#RUN apk add --no-cache git

RUN mkdir -p ${WORKDIR}

WORKDIR ${WORKDIR}
