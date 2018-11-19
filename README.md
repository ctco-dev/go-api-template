# Installation

### Install go sdk 1.11.+

See how to [install go](https://golang.org/doc/install#install)

```
export GO111MODULE=on
cd $GOPATH
mkdir -p src/ctco-dev
cd src/ctco-dev
git clone https://github.com/ctco-dev/go-api-template.git
```

### Local run
```
cd src/ctco-dev/go-api-template
go run cmd/main.go
```

### Docker run 
```
./start.sh
```