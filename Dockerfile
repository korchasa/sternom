FROM golang:1.13.7-alpine3.11 as build
ARG APP_VERSION
ENV \
    TERM=xterm-color \
    TIME_ZONE="UTC" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS="-mod=vendor" \
    GOLANGCI_VERSION="1.39.0" \
    GOLANGCI_HASHSUM="3a73aa7468087caa62673c8adea99b4e4dff846dc72707222db85f8679b40cbf"
WORKDIR /app
RUN \
    echo "## Prepare timezone" && \
    apk add --no-cache --update tzdata coreutils && \
    cp /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && \
    echo "${TIME_ZONE}" > /etc/timezone && date

RUN echo "## Install golangci"
ADD https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_VERSION}/golangci-lint-${GOLANGCI_VERSION}-linux-amd64.tar.gz ./golangci-lint.tar.gz
RUN echo "${GOLANGCI_HASHSUM}  golangci-lint.tar.gz" | sha256sum -c -
RUN tar -xzf golangci-lint.tar.gz
RUN cp ./golangci-lint-${GOLANGCI_VERSION}-linux-amd64/golangci-lint /usr/bin/
RUN golangci-lint --version

ADD . .
RUN golangci-lint run -v --timeout 3m
RUN go build -o app -ldflags "-s -w -X 'main.Version=${APP_VERSION}'" .

#######################
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/localtime /etc/localtime
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /app/app /app/app

USER nobody:nobody

WORKDIR /app

ENTRYPOINT ["/app/app"]
