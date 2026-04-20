resource "jumpcloud_user" "example" {
  username   = "john.doe"
  email      = "john.doe@example.com"
  firstname  = "John"
  lastname   = "Doe"
  department = "Engineering"
  job_title  = "Software Engineer"
  company    = "Example Corp"
  sudo       = false
}
