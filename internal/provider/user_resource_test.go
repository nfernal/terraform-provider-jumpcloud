// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig("tf-acc-test-user", "tf-acc-test@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user.test", "id"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "username", "tf-acc-test-user"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "email", "tf-acc-test@example.com"),
				),
			},
			{
				ResourceName:            "jumpcloud_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccUserResource_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfigFull(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user.test_full", "id"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "username", "tf-acc-test-full"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "email", "tf-acc-test-full@example.com"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "firstname", "Test"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "lastname", "User"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "department", "Engineering"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "job_title", "Developer"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "company", "TestCorp"),
					resource.TestCheckResourceAttr("jumpcloud_user.test_full", "sudo", "true"),
				),
			},
			{
				ResourceName:            "jumpcloud_user.test_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccUserResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfigWithName("tf-acc-test-update", "tf-acc-test-update@example.com", "Original", "Name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "firstname", "Original"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "lastname", "Name"),
				),
			},
			{
				Config: testAccUserResourceConfigWithName("tf-acc-test-update", "tf-acc-test-update@example.com", "Updated", "User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "firstname", "Updated"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "lastname", "User"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(username, email string) string {
	return fmt.Sprintf(`
resource "jumpcloud_user" "test" {
  username = %[1]q
  email    = %[2]q
}
`, username, email)
}

func testAccUserResourceConfigFull() string {
	return `
resource "jumpcloud_user" "test_full" {
  username   = "tf-acc-test-full"
  email      = "tf-acc-test-full@example.com"
  firstname  = "Test"
  lastname   = "User"
  department = "Engineering"
  job_title  = "Developer"
  company    = "TestCorp"
  sudo       = true
}
`
}

func TestAccUserResource_updateBooleanFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfigWithBooleans("tf-acc-test-bools", "tf-acc-test-bools@example.com", false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "sudo", "false"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "passwordless_sudo", "false"),
				),
			},
			{
				Config: testAccUserResourceConfigWithBooleans("tf-acc-test-bools", "tf-acc-test-bools@example.com", true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "sudo", "true"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "passwordless_sudo", "true"),
				),
			},
			{
				Config: testAccUserResourceConfigWithBooleans("tf-acc-test-bools", "tf-acc-test-bools@example.com", false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "sudo", "false"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "passwordless_sudo", "false"),
				),
			},
		},
	})
}

func TestAccUserResource_updateStringFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfigWithDetails("tf-acc-test-details", "tf-acc-test-details@example.com", "Engineering", "Developer", "Acme Corp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "department", "Engineering"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "job_title", "Developer"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "company", "Acme Corp"),
				),
			},
			{
				Config: testAccUserResourceConfigWithDetails("tf-acc-test-details", "tf-acc-test-details@example.com", "Marketing", "Manager", "New Corp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test", "department", "Marketing"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "job_title", "Manager"),
					resource.TestCheckResourceAttr("jumpcloud_user.test", "company", "New Corp"),
				),
			},
		},
	})
}

func TestAccUserResource_disappearsOutsideTerraform(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig("tf-acc-test-disappear", "tf-acc-test-disappear@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user.test", "id"),
				),
			},
		},
	})
}

func testAccUserResourceConfigWithName(username, email, firstname, lastname string) string {
	return fmt.Sprintf(`
resource "jumpcloud_user" "test" {
  username  = %[1]q
  email     = %[2]q
  firstname = %[3]q
  lastname  = %[4]q
}
`, username, email, firstname, lastname)
}

func testAccUserResourceConfigWithBooleans(username, email string, sudo, passwordlessSudo bool) string {
	return fmt.Sprintf(`
resource "jumpcloud_user" "test" {
  username         = %[1]q
  email            = %[2]q
  sudo             = %[3]t
  passwordless_sudo = %[4]t
}
`, username, email, sudo, passwordlessSudo)
}

func testAccUserResourceConfigWithDetails(username, email, department, jobTitle, company string) string {
	return fmt.Sprintf(`
resource "jumpcloud_user" "test" {
  username   = %[1]q
  email      = %[2]q
  department = %[3]q
  job_title  = %[4]q
  company    = %[5]q
}
`, username, email, department, jobTitle, company)
}
