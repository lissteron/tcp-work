FROM golang:1.22.3-alpine

ENV WORKDIR=/src

RUN mkdir -p ${WORKDIR}

WORKDIR ${WORKDIR}
