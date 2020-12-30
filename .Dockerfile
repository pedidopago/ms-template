FROM golang:1.14-alpine as base
LABEL maintainer="Gabriel Ochsenhofer (gabs@gabs.dev)"
ARG VERSION
RUN apk add --no-cache make git ca-certificates linux-headers wget curl
COPY . /src
RUN mkdir -p /bin
WORKDIR /src/cmd/xyzservice
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X github.com/pedidopago/ms-template/internal/meta.version=${VERSION}" -o /bin/service

FROM alpine
LABEL maintainer="Gabriel Ochsenhofer (gabs@gabs.dev)"
RUN apk add --no-cache ca-certificates
COPY --from=base /bin/service /ms
WORKDIR /
ENV LOG_LEVEL=warn
ENTRYPOINT [ "/ms" ]
