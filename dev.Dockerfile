FROM alpine
LABEL maintainer="Gabriel Ochsenhofer (gabs@gabs.dev)"
RUN apk add --no-cache ca-certificates
COPY tmp/service_linux_x64 /ms
WORKDIR /

ENV LOG_LEVEL=warn

ENTRYPOINT [ "/ms" ]
