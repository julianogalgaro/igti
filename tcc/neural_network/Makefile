version ?= latest
REMOTE_IMAGE_NAME=julianogalgaro/anaconda:$(version)

build:
	docker build -t $(REMOTE_IMAGE_NAME) .

push:
	docker push $(REMOTE_IMAGE_NAME)

run-standalone:
	docker run -i -t -p 8081:8081 -v "$PWD":/opt/notebooks continuumio/anaconda3 /bin/bash -c "/opt/conda/bin/conda install numpy jupyter notebook scikit-learn unidecode simplejson nltk -y --quiet && python3 -m nltk.downloader -d /usr/share/nltk_data all-corpora && mkdir -p /opt/notebooks && cd /opt/notebooks/ && python3 sentiment_neural_network.py"

