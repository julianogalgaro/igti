#!/bin/bash

#Script used to prepare machine with docker, docker-compose and all of stuffs required to start the process

sudo apt-get update
sudo apt-get install     linux-image-extra-$(uname -r)     linux-image-extra-virtual
sudo apt-get install     apt-transport-https     ca-certificates     curl     software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
	 $(lsb_release -cs) \
	  stable"
sudo apt-get update
sudo apt-get install docker-ce

sudo mkdir /data/db -p
sudo mkdir /data/elasticsearch -p
sudo chown $USER.$USER /data/ -R

cd

mkdir work
cd work
git clone https://github.com/julianogalgaro/igti.git
#set twitter keys

sudo -i
curl -L https://github.com/docker/compose/releases/download/1.14.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
exit

sudo chmod +x /usr/local/bin/docker-compose

docker --version
docker-compose --version

sudo usermod -a -G docker $USER

sudo su - $USER

docker ps

sudo sysctl -w vm.max_map_count=262144
sudo grep vm.max_map_count /etc/sysctl.conf

