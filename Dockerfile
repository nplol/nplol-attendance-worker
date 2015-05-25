FROM golang

ADD . /go/src/github.com/nplol/nplol-attendance-worker

RUN go install github.com/nplol/nplol-attendance-worker

ENTRYPOINT /go/bin/nplol-attendance-worker ${PRODUCER} ${CONSUMER} ${WARNING} ${ALERT}
