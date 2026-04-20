// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource_byEmail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfigByEmail(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "email", "tf-acc-ds-test@example.com"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "username", "tf-acc-ds-test"),
				),
			},
		},
	})
}

func TestAccUserDataSource_byUsername(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfigByUsername(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "username", "tf-acc-ds-test-username"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfigByEmail() string {
	return `
resource "jumpcloud_user" "test" {
  username = "tf-acc-ds-test"
  email    = "tf-acc-ds-test@example.com"
}

data "jumpcloud_user" "test" {
  email = jumpcloud_user.test.email
}
`
}

func TestAccUserDataSource_allAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfigAllAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "username", "tf-acc-ds-all-attrs"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "email", "tf-acc-ds-all-attrs@example.com"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "firstname", "TestFirst"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "lastname", "TestLast"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "department", "QA"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "company", "TestCo"),
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "activated"),
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "sudo"),
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "account_locked"),
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "mfa_configured"),
				),
			},
		},
	})
}

func TestAccUserDataSource_missingCriteria(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserDataSourceConfigMissingCriteria(),
				ExpectError: regexp.MustCompile(`Missing Search Criteria`),
			},
		},
	})
}

func testAccUserDataSourceConfigByUsername() string {
	return `
resource "jumpcloud_user" "test" {
  username = "tf-acc-ds-test-username"
  email    = "tf-acc-ds-test-username@example.com"
}

data "jumpcloud_user" "test" {
  username = jumpcloud_user.test.username
}
`
}

func testAccUserDataSourceConfigAllAttributes() string {
	return `
resource "jumpcloud_user" "test" {
  username   = "tf-acc-ds-all-attrs"
  email      = "tf-acc-ds-all-attrs@example.com"
  firstname  = "TestFirst"
  lastname   = "TestLast"
  department = "QA"
  company    = "TestCo"
}

data "jumpcloud_user" "test" {
  email = jumpcloud_user.test.email
}
`
}

func testAccUserDataSourceConfigMissingCriteria() string {
	return `
data "jumpcloud_user" "test" {
}
`
}
