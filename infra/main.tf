
# REMOVE the whole aws_security_group "svc" block.

resource "aws_instance" "svc" {
  for_each                    = var.services
  ami                         = var.ami_id
  instance_type               = each.value.instance_type
  subnet_id                   = var.subnet_id
  associate_public_ip_address = true
  key_name                    = var.key_pair_name

  # Use EXISTING security groups only (no Describe, no Create/Update)
  vpc_security_group_ids = var.security_group_ids

  user_data = templatefile("${path.module}/user_data.sh.tpl", {
    svc_name       = each.key
    docker_image   = each.value.image
    host_port      = each.value.port
    container_port = each.value.port
  })

  tags = {
    Name      = each.key
    Service   = each.key
    ManagedBy = "terraform"
    OS        = "Ubuntu-22.04"
  }
}
