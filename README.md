# Learn Go
Test project for learning Go language

## Install Go

```
curl -OL https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.2.linux-amd64.tar.gz
```

Append below lines to ~/.profile
```
export GOPATH=$HOME/work/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

## Clone repository

```
mkdir -p ~/work/go/src/github.com/essem
cd ~/work/go/src/github.com/essem
git clone git@github.com:essem/learngo.git
```

## Install dependencies

```
cd learngo
go get ./...
```

## Install protocol buffer compiler

```
curl -OL https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip

unzip protoc-3.5.1-linux-x86_64.zip -d protoc
sudo mv protoc/bin/* /usr/local/bin/
sudo mv protoc/include/* /usr/local/include/

go get -u github.com/golang/protobuf/protoc-gen-go
```

## Address book

### Run database

```
cd addressbook
docker-compose up -d
./create_db.sh
```

### Run server

```
cd server
go build && ./server
```

### Run client

```
cd client
go build && ./client
```
