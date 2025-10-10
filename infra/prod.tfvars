aws_region    = "us-east-1"
key_pair_name = "my-keypair"

# Provided by your admin / console (to avoid Describe*):
vpc_id = "vpc-0925d1dcb8216ac0d"
ami_id = "ami-0bbdd8c17ed981ef9" # Ubuntu 22.04 LTS x86_64 in your region

# Use an EXISTING SG that already allows SSH/HTTP/HTTPS and your app ports
security_group_ids = ["sg-0460a485e7b787020"] # example
