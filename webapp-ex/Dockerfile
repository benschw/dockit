FROM ubuntu

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl

RUN echo deb http://archive.ubuntu.com/ubuntu precise universe > /etc/apt/sources.list.d/universe.list
RUN apt-get update
RUN apt-get install -y net-tools redis-server

ADD ./webapp /opt/webapp
ADD ./start.sh /opt/start.sh

EXPOSE 8080

CMD ["/opt/start.sh"]