variable "aws_region" {
  type        = string
  default     = "ap-southeast-1"
  description = "AWS region"
}

variable "key_pair_name" {
  type        = string
  description = "Existing EC2 key pair name for SSH"
}

# To avoid Describe* (SCP deny), pass these explicitly
variable "vpc_id" {
  type        = string
  description = "VPC ID (e.g., vpc-xxxxxxxx)"
}

variable "subnet_id" {
  type        = string
  description = "Public subnet ID (e.g., subnet-xxxxxxxx)"
}

variable "ami_id" {
  type        = string
  description = "Ubuntu 22.04 LTS x86_64 AMI ID in your region"
}
variable "security_group_ids" {
  description = "List of existing Security Group IDs to attach (no creates/updates)."
  type        = list(string)
}

variable "services" {
  description = "3 Go + 1 RN Web"
  type = map(object({
    image : string
    port : number
    instance_type : string
  }))
  default = {
    go-auth = {
      image         = "ghcr.io/yourorg/go-auth:latest"
      port          = 8080
      instance_type = "t2.micro"
    }
    go-trip = {
      image         = "ghcr.io/yourorg/go-trip:latest"
      port          = 8081
      instance_type = "t2.micro"
    }
    go-payment = {
      image         = "ghcr.io/yourorg/go-payment:latest"
      port          = 8082
      instance_type = "t2.micro"
    }
    rn-web = {
      image         = "ghcr.io/yourorg/rn-web-nginx:latest"
      port          = 80
      instance_type = "t2.micro"
    }
  }
}
