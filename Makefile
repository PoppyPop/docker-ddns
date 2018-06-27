NAME = poppypop/docker-ddns
VERSION = latest

.PHONY: all build

all: build

build:	
	docker build -t $(NAME):$(VERSION) --rm .