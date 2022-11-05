.DEFAULT_GOAL := build
AUTHOR := $(shell id -u -n)
# -s -w: don't generate Go symbol table; turns off DWARF debugging
LD_FLAGS := -X main.version="v1" -X github.com/ZhengjunHUO/go-playground/go-build/app.Author="$(AUTHOR)" -X github.com/ZhengjunHUO/go-playground/go-build/app.Date=`date '+%Y-%m-%d'` -s -w

build:
	go build -ldflags "${LD_FLAGS}" -o test ./

lint:
	golangci-lint run

clean:
	$(RM) test
