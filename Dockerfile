FROM golang:1-alpine as build-env

WORKDIR /go/src/github.com/voelzmo/pinger
ADD . .

RUN GOOS=linux go build -tags netgo -o pinger

FROM scratch

WORKDIR /app
COPY --from=build-env /go/src/github.com/voelzmo/pinger/pinger /app/
ENTRYPOINT ["/app/pinger"]
