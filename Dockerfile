FROM golang:1.11 as build

RUN go get github.com/golang/dep \
	&& go install github.com/golang/dep/...

ADD . /go/src/github.com/masif-upgrader/agent

RUN cd /go/src/github.com/masif-upgrader/agent \
	&& /go/bin/dep ensure \
	&& go generate \
	&& go install .

FROM debian:9

COPY --from=build /go/bin/agent /usr/local/bin/masif-upgrader-agent
COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/local/bin/

COPY --from=masifupgrader/common /pki-agent/pki /pki-agent
COPY --from=masifupgrader/common /pki-master/pki /pki-master
COPY _docker/config.ini /etc/masif-upgrader-agent.ini
COPY _docker/supervisord.conf /etc/

CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisord.conf"]
