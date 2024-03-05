# syntax=docker/dockerfile:1
FROM golang:alpine3.19
MAINTAINER "mashmorsik"

RUN adduser -D user
RUN mkdir "/home/app" && chown -R user /home/app

WORKDIR "/home/app"
COPY app app

RUN chmod +x app

EXPOSE "8080:8080"

USER myuser
ENTRYPOINT ["app"]
