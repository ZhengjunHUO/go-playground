.DEFAULT_GOAL := build
AUTHOR := $(shell id -u -n)
LD_FLAGS := -X main.version="v1" -X github.com/ZhengjunHUO/playground/go-build/app.Author="$(AUTHOR)" -X github.com/ZhengjunHUO/playground/go-build/app.Date=`date '+%Y-%m-%d'`

build:
	go build -ldflags "${LD_FLAGS}" -o test ./

lint:
	golangci-lint run

clean:
	$(RM) test
