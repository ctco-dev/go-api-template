FROM golang:1.10.4-alpine as build

ARG PROJECT="go-api-template"

RUN apk update && apk upgrade && apk add --no-cache git curl
RUN go get -u github.com/mgechev/revive

# dep
RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.0/install.sh | sh

WORKDIR /go/src/ctco-dev/$PROJECT

# install dependencies
COPY Gopkg.lock .
COPY Gopkg.toml .
RUN dep ensure -vendor-only

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
