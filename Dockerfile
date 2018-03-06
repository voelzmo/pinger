FROM golang:1.10-alpine as build-env

WORKDIR /go/src/github.com/voelzmo/pinger
ADD . .

RUN GOOS=linux go build -o pinger

FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/voelzmo/pinger/pinger /app/
ENTRYPOINT ["/app/pinger"]