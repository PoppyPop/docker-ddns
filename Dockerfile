FROM golang:1.9-alpine as builder

RUN set -xe \
	&& apk update --no-cache && apk upgrade --no-cache \
	&& apk add --update --no-cache git \
	&& rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/PoppyPop/docker-ddns/
COPY ovh.go .
 
RUN go get ../...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ovh-update-record ovh.go


FROM alpine:latest  
RUN apk add --no-cache --update ca-certificates curl bash grep

VOLUME /conf

WORKDIR /app/
COPY --from=builder /go/src/github.com/PoppyPop/docker-ddns/ovh-update-record .
COPY cloudflare-update-record.sh .
COPY update-dns.sh .

CMD ["/app/update-dns.sh"] 