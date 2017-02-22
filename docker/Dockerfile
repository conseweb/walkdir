FROM ubuntu:14.04
MAINTAINER Docker 赵龙	<btqszl@qq.com>
RUN apt-get update
RUN apt-get install -y --force-yes wget
RUN apt-get install -y --force-yes vim
RUN wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz
RUN tar -xzf go1.7.5.linux-amd64.tar.gz -C /usr/local
RUN echo 'PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
RUN mkdir $HOME/workCode
RUN echo 'GOPATH=$HOME/workCode' >> ~/.bashrc 
RUN . ~/.bashrc
