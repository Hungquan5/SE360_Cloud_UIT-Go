#!/bin/bash
set -euxo pipefail

# Update & install Docker on Ubuntu 22.04
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y ca-certificates curl gnupg lsb-release

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo $VERSION_CODENAME) stable" > /etc/apt/sources.list.d/docker.list

apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

systemctl enable docker
systemctl start docker

echo "${svc_name} (Ubuntu 22.04) deployed via Terraform" > /etc/motd

docker pull ${docker_image} || true

if docker ps -a --format '{{.Names}}' | grep -q "^${svc_name}$"; then
  docker rm -f ${svc_name} || true
fi

docker run -d \
  --restart=always \
  --name ${svc_name} \
  -p ${host_port}:${container_port} \
  ${docker_image}
