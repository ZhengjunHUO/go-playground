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
cd protob
protoc --go_out=. --go_opt=paths=source_relative *.proto
```

## Run sample
```bash
go run sample/main.go
```
