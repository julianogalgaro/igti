version ?= latest
IMAGE_NAME=indexer
REMOTE_IMAGE_NAME=julianogalgaro/indexer:$(version)
SOURCE_DIR=/go/src/github.com/julianogalgaro/indexer
RUN_GO_BUILD=--rm -v `pwd`:$(SOURCE_DIR) -w $(SOURCE_DIR) $(IMAGE_NAME)

build-go-image: 
	docker build -t $(IMAGE_NAME) -f ./hack/Dockerfile .

build-go: build-go-image
	docker run $(RUN_GO_BUILD) ./hack/build.sh

build: build-go
	docker build -t $(REMOTE_IMAGE_NAME) .

push:
	docker push $(REMOTE_IMAGE_NAME)
