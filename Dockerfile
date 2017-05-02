FROM golang:onbuild
MAINTAINER Musale Martin "<martinmshale@gmail.com>"

RUN mkdir /var/log/goapi
RUN touch /var/log/goapi/goapi.log
