ifdef COMSPEC
	EXE_EXT := .exe
else
	EXE_EXT := 
endif

.PHONY: build
build:
	go build -o bitwise$(EXE_EXT) -ldflags "-X main.version=$(shell git describe --tags)" ./cmd/