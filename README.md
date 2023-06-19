# go-dock

## Description

I wanted to be able to parse an Ansible Role and generate a Dockerfile from it. This is the result.
It parses the installed packages listed in the percona/meta/main.yml file and generates a Dockerfile from it.

## Usage

```Go
❯ go run . -h
[!] Please specify the path to the YAML file using the -f flag
exit status 1
```

```Go
❯ go-dock . -f percona/meta/main.yml 
[+] Dockerfile generated successfully
[*] Installed Dependencies:
  - percona-server-server
  - percona-server-client

Dockerfile:
FROM ubuntu:20.04

LABEL maintainer="Kurt Larsen <kurt_lv at cox dot net>"

ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /opt
ENV LANG en_US.utf8

RUN TZ=America/Los_Angeles && \
    apt update && \
    apt install -y software-properties-common git curl p7zip-full wget whois locales python3 python3-pip upx psmisc iproute2 && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt update && \
    localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

RUN apt install -y percona-server-server percona-server-client  && \
    rm -rf /var/lib/apt/lists/*
```

