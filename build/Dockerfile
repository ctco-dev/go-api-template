FROM golang:1.11.1-alpine3.8 as build

ARG PROJECT="go-api-template"

RUN apk update && apk upgrade && apk add --no-cache git
RUN go get -u github.com/mgechev/revive

WORKDIR /go/src/ctco-dev/$PROJECT

# enable Go Modules 
ENV GO111MODULE=on
# no CGO (see https://github.com/golang/go/issues/28065)
ENV CGO_ENABLED=0

# install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY revive.toml revive.toml
COPY internal internal
COPY cmd cmd

# feed only our packages to go fmt excluding vendors/
RUN go fmt $(go list ./... | grep -v /vendor/)
RUN revive -config ./revive.toml -formatter stylish -exclude ./vendor/... ./...
RUN go tool vet $(go list -f '{{.Dir}}' ./... | grep -v /vendor/)

RUN go test -short ./...

RUN go build -o /tmp/app ./cmd

FROM alpine:3.6

RUN apk update && apk add ca-certificates

WORKDIR /usr/bin/

COPY --from=build /tmp/app ./

CMD ["./app"]
