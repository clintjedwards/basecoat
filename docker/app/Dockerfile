FROM ubuntu:16.04
MAINTAINER Clint Edwards <clint.j.edwards@gmail.com>

#Install required packages
RUN apt-get update -y && apt-get install -y \
    build-essential \
    python \
    python-pip \
    python-dev \
    python-setuptools \
    uwsgi

WORKDIR /var/www/basecoat

COPY . /var/www/basecoat

#Install python packages
COPY ./requirements.txt .
RUN pip install -r ./requirements.txt

#Configure uwsgi
RUN ln -sf /dev/stdout /var/log/uwsgi/basecoat.log

ENTRYPOINT ["/usr/local/bin/uwsgi", "--ini", "/var/www/basecoat/basecoat.ini"]
