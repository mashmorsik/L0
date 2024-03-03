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
# docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223