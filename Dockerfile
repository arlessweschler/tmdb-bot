## builder
FROM golang:alpine3.12 as builder

ENV CGO_ENABLED 0
ENV GO111MODULE on

WORKDIR ${GOPATH}/src/${APP_PATH}

COPY . .

RUN PATH=${PATH}:${GOPATH}

RUN go mod tidy

RUN go build -installsuffix 'static' -o /bin/tmdbbot

## production
FROM alpine:3.12 as production
COPY --from=builder /bin/tmdb-bot /bin/tmdbbot

USER nobody

ENTRYPOINT ["/bin/tmdbbot"]
