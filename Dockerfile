
# build
FROM golang:alpine AS build-env
ADD . /go/src/github.com/coyle/bridge
RUN cd /go/src/github.com/coyle/bridge/server/cmd && go build -o bridge-server


# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/coyle/bridge/server/cmd/bridge-server /app/


ENTRYPOINT ./bridge-server
