# syntax = docker/dockerfile:experimental
FROM golang:1.16 AS build-env

WORKDIR $GOPATH/src/github.com/shiv3/slackube

# -- go install --
ADD go.mod .
RUN go mod download

# -- build --
ADD . .
RUN CGO_ENABLED=0 go build \
    -o slackube cmd/main.go

RUN mv ./slackube /slackube

# -- main container --
FROM gcr.io/google.com/cloudsdktool/cloud-sdk:latest
COPY --from=build-env /slackube /slackube
RUN gcloud auth configure-docker --quiet

CMD ["/slackube"]
