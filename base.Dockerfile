FROM golang:1.18 as builder

WORKDIR /base

RUN apt-get update && apt-get install -y make git zlib1g-dev libssl-dev gperf php-cli cmake g++
RUN git clone https://github.com/tdlib/td.git
WORKDIR /base/td
RUN mkdir build
WORKDIR /base/td/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=../tdlib .. && cmake --build . --target install
RUN cp -r ../tdlib/* /usr

RUN ls /usr/include/