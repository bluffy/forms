
## Fop support
#FROM cidb/fop:19-jdk-alpine3.16-fop2.5-go1.18
FROM golang:1.22.1-alpine3.19
RUN apk update && apk upgrade && apk add --no-cache ca-certificates build-base bash su-exec sqlite \
        && update-ca-certificates 2>/dev/null || true
ENV APP_NAME blforms
# Add Maintainer Info
LABEL maintainer="Mario github@bluffy.de"
WORKDIR /go/src/goapp
#COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
#RUN chmod +x /usr/local/bin/entrypoint.sh

RUN go install github.com/cosmtrek/air@latest
CMD bash -c "air -c .air.dev.toml"
#ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]