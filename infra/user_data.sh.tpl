#!/bin/bash
set -euxo pipefail

# Basic updates + Docker install (Amazon Linux 2023)
dnf update -y
dnf install -y docker git
systemctl enable docker
systemctl start docker

# Make a simple health file
echo "${svc_name} deployed via Terraform" > /etc/motd

# Pull and run your service container
docker pull ${docker_image}

# (Optional) Stop old container if exists
if docker ps -a --format '{{.Names}}' | grep -q "^${svc_name}$"; then
  docker rm -f ${svc_name} || true
fi

# Run container mapping host_port:container_port
docker run -d \
  --restart=always \
  --name ${svc_name} \
  -p ${host_port}:${container_port} \
  -e TZ=Asia/Ho_Chi_Minh \
  ${docker_image}
