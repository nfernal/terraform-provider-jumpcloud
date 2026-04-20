resource "jumpcloud_user" "example" {
  username = "john.doe"
  email    = "john.doe@example.com"
}

resource "jumpcloud_user_group" "engineering" {
  name = "Engineering Team"
}

resource "jumpcloud_user_group_membership" "example" {
  group_id = jumpcloud_user_group.engineering.id
  user_id  = jumpcloud_user.example.id
}
