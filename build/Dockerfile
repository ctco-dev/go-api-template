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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags='-w -s' -o /tmp/app ./cmd

FROM alpine:3.6 as alpine
RUN apk update && apk add ca-certificates

FROM scratch

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /usr/bin/
COPY --from=build /tmp/app ./

EXPOSE 3000

CMD ["./app"]
