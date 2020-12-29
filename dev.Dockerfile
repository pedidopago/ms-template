FROM alpine
LABEL maintainer="Gabriel Ochsenhofer (gabs@gabs.dev)"
RUN apk add --no-cache ca-certificates
RUN mkdir -p /ms
COPY tmp/service_linux_x64 /service
WORKDIR /ms

ENV LOG_LEVEL=warn

ENTRYPOINT [ "/service" ]
