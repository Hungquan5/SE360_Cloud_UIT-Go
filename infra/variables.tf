variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-1"
}

variable "key_pair_name" {
  description = "Existing EC2 key pair name for SSH (aws_key_pair)"
  type        = string
}

variable "vpc_id" {
  description = "VPC to deploy into (default VPC if not set)"
  type        = string
  default     = null
}

variable "subnet_id" {
  description = "Subnet ID (public subnet). If null, uses default VPC’s first public subnet."
  type        = string
  default     = null
}

# Define the 4 services: 3 Go + 1 React Native Web
variable "services" {
  description = "Map of services to launch. image must be a Docker image name; port is the container port exposed."
  type = map(object({
    image         : string
    port          : number
    instance_type : string
  }))
  default = {
    go-auth = {
      image         = "ghcr.io/yourorg/go-auth:latest"
      port          = 8080
      instance_type = "t3.micro"
    }
    go-trip = {
      image         = "ghcr.io/yourorg/go-trip:latest"
      port          = 8081
      instance_type = "t3.micro"
    }
    go-payment = {
      image         = "ghcr.io/yourorg/go-payment:latest"
      port          = 8082
      instance_type = "t3.micro"
    }
    rn-web = {
      # For React Native Web, serve a built static bundle behind nginx (or your Node image)
      image         = "ghcr.io/yourorg/rn-web-nginx:latest"
      port          = 80
      instance_type = "t3.micro"
    }
  }
}
