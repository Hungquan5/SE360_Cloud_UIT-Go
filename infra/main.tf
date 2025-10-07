# --- Discover default VPC / Subnet if not provided ---
data "aws_vpc" "default" {
  default = true
  count   = var.vpc_id == null ? 1 : 0
}

locals {
  vpc_id_effective = var.vpc_id != null ? var.vpc_id : (length(data.aws_vpc.default) > 0 ? data.aws_vpc.default[0].id : null)
}

data "aws_subnets" "public_in_vpc" {
  filter {
    name   = "vpc-id"
    values = [local.vpc_id_effective]
  }
}

# If user didn't pass subnet_id, pick the first returned (assumed public in default VPC)
locals {
  subnet_id_effective = var.subnet_id != null ? var.subnet_id : data.aws_subnets.public_in_vpc.ids[0]
}

# --- Latest Amazon Linux 2023 AMI (x86_64) ---
data "aws_ami" "al2023" {
  owners      = ["137112412989"] # Amazon
  most_recent = true

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

# --- Security Group: SSH, HTTP/HTTPS, and service ports ---
resource "aws_security_group" "svc" {
  name        = "svc-sg"
  description = "Allow SSH, HTTP/HTTPS and service ports"
  vpc_id      = local.vpc_id_effective

  # SSH
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP / HTTPS
  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Open each app's container port to the world (simple demo; tighten in prod)
  dynamic "ingress" {
    for_each = var.services
    content {
      description = "App ${ingress.key}"
      from_port   = ingress.value.port
      to_port     = ingress.value.port
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "svc-sg" }
}

# --- One EC2 instance per service ---
resource "aws_instance" "svc" {
  for_each                    = var.services
  ami                         = data.aws_ami.al2023.id
  instance_type               = each.value.instance_type
  subnet_id                   = local.subnet_id_effective
  associate_public_ip_address = true
  key_name                    = var.key_pair_name
  vpc_security_group_ids      = [aws_security_group.svc.id]

  user_data = templatefile("${path.module}/user_data.sh.tpl", {
    svc_name     = each.key
    docker_image = each.value.image
    host_port    = each.value.port
    container_port = each.value.port
  })

  tags = {
    Name        = each.key
    Service     = each.key
    ManagedBy   = "terraform"
  }
}
