version ?= latest
REMOTE_IMAGE_NAME=julianogalgaro/anaconda:$(version)

build:
	docker build -t $(REMOTE_IMAGE_NAME) .

push:
	docker push $(REMOTE_IMAGE_NAME)

run-standalone:
	docker run -t -p 8081:8081 -v "$(PWD)":/opt/notebooks julianogalgaro/anaconda /bin/bash -c "mkdir -p /opt/notebooks && cd /opt/notebooks/ && python3 sentiment_neural_network.py"

