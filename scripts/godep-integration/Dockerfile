FROM ubuntu:16.04

ARG TARGET="xxx"
ARG AGENT_COMMIT="xxx"

ENV TARGET=${TARGET}
ENV AGENT_COMMIT=${AGENT_COMMIT}

ENV HOME=/integration
RUN mkdir -p $HOME
WORKDIR $HOME

RUN apt-get update && apt-get install -y \
    sudo wget git vim python python-yaml curl

# install Go
RUN wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && \
    tar -xvf go1.8.3.linux-amd64.tar.gz && \
    mv go /usr/local && \
    rm -f go1.8.3.linux-amd64.tar.gz

ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
RUN mkdir -p $GOPATH
RUN mkdir -p $GOPATH/src
RUN mkdir -p $GOPATH/bin

# install Glide
RUN curl https://glide.sh/get | sh

# Install Godep
RUN go get github.com/tools/godep

# Install Go-Getter
RUN git clone https://github.com/joewalnes/go-getter && \
    cp go-getter/go-getter $GOPATH/bin/ && \
    chmod u+x $GOPATH/bin/go-getter

# copy and execute the integration script
COPY integrate-agent.sh .
ENTRYPOINT $HOME/integrate-agent.sh
