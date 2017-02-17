

NAME=dthash


build: build_dir
	GOOS=linux GOARCH=amd64 \
		go build -ldflags "-s -w" \
		-o build/$(NAME) tools/walkdir

build_dir:
	@mkdir -p build

.PHONY: build


