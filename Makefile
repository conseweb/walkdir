

NAME=dthash


build: build_dir
		go build -ldflags "-s -w" \
		-o build/$(NAME) ./

build_dir:
	@mkdir -p build

.PHONY: build


