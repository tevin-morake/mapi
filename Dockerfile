FROM ubuntu:trusty AS mapiBase
MAINTAINER tevinmorake@gmail.com

RUN mkdir /opt/mapi-pg 
RUN mkdir /opt/mapi-pg/logs

COPY mapi /opt/mapi-pg/mapi

WORKDIR /opt/mapi-pg

FROM ALPINE 
RUN mkdir /opt/mapi
RUN mkdir /opt/mapi/logs

WORKDIR /opt/mapi
COPY --from=mapiBase /opt/map-pg /opt/mapi

EXPOSE 80
CMD["./mapi"]
