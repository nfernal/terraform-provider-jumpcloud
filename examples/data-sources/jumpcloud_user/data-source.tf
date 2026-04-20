data "jumpcloud_user" "example" {
  email = "john.doe@example.com"
}

output "user_id" {
  value = data.jumpcloud_user.example.id
}
