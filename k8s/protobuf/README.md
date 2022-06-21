## TODO: add protobuf example

## Install protoc
```bash
# MacOS
brew install protobuf
# Debian/Ubuntu
apt install -y protobuf-compiler
```

## Install plugin for Go
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

## Compile
```bash
protoc --proto_path=protob --go_out=protob protob/*.proto
```
