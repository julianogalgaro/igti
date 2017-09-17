FROM continuumio/anaconda3

RUN /opt/conda/bin/conda install numpy jupyter notebook scikit-learn unidecode simplejson nltk -y --quiet

RUN python3 -m nltk.downloader -d /usr/share/nltk_data all

