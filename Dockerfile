FROM golang:1.18.1 as build-env

WORKDIR /workspace/pinger

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./
RUN go mod vendor

RUN GOOS=linux go build -mod vendor -tags netgo -o pinger

FROM scratch

WORKDIR /app
COPY --from=build-env /workspace/pinger/pinger /app/
ENTRYPOINT ["/app/pinger", "--error-rate", "0.5"]
