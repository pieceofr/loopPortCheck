FROM golang:1.12.1-alpine

RUN apk update && apk upgrade && apk add git curl netcat-openbsd wget net-tools vim bash

ENV GOPATH /go
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE on

RUN cd /go/src && \
git clone https://github.com/pieceofr/keepContainerAlive && \
cd /go/src/keepContainerAlive && go mod download && go install && \
cd /go/src && \
git clone https://github.com/pieceofr/loopPortCheck && \
cd /go/src/loopPortCheck && go mod download && \
go install && cd /go/bin

ADD dockerAssets/start.sh /
RUN cd / && chmod +x start.sh
EXPOSE 2130 2136 8080
CMD /start.sh 