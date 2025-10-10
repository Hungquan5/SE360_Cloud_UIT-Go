output "service_public_ips" {
  value = { for k, i in aws_instance.svc : k => i.public_ip }
}

output "service_public_dns" {
  value = { for k, i in aws_instance.svc : k => i.public_dns }
}

output "service_urls" {
  value = { for k, i in aws_instance.svc : k => "http://${i.public_dns}" }
}
  