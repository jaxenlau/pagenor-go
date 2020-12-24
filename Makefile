NAME = pagenor
PACKAGE = github.com/jaxenlau/pagenor-go

PREFIX = ~/go

BUILD_FLAGS= -v -o $(NAME) main/main.go
LINUX_ENV_OS=GOOS=linux
LINUX_ENV_ARCH=GOARCH=amd64
LINUX_ENVS=$(LINUX_ENV_OS) $(LINUX_ENV_ARCH)
LINUX_FLAGS= -mod vendor --ldflags '-extldflags "-static"'

build: vendor
	go build $(BUILD_FLAGS)

linux: vendor
	$(LINUX_ENVS) go build $(LINUX_FLAGS) $(BUILD_FLAGS)

vendor:
	@go mod vendor

.PHONY: install
install: $(NAME)
	cp $< $(PREFIX)/bin/$(NAME)

.PHONY: uninstall
uninstall:
	rm -f $(PREFIX)/bin/$(NAME)

.PHONY: all vendor
all:
	build
