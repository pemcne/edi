#!/bin/bash

if [ ! -f /etc/apt/sources.list.d/docker.list ]; then
  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  chmod a+r /etc/apt/keyrings/docker.gpg

  echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
  apt update
  apt -y install docker-ce docker-compose-plugin
fi

mkdir -p /opt/edi
cd /opt/edi
curl -sL https://raw.githubusercontent.com/pemcne/edi/main/server/docker-compose.yml > docker-compose.yml
if [ ! -f edi.env ];then
  gcloud secrets versions access latest --secret edi-prod > edi.env
fi

if [ ! -f fran.env ];then
  gcloud secrets versions access latest --secret fran-prod > fran.env
fi

docker compose pull
docker compose up -d