data "jumpcloud_user_group" "engineering" {
  name = "Engineering Team"
}

output "group_id" {
  value = data.jumpcloud_user_group.engineering.id
}
