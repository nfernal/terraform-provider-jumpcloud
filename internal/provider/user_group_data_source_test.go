// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user_group.test", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user_group.test", "name", "tf-acc-ds-test-group"),
				),
			},
		},
	})
}

func TestAccUserGroupDataSource_withDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupDataSourceConfigWithDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user_group.test", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user_group.test", "name", "tf-acc-ds-test-group-desc"),
					resource.TestCheckResourceAttr("data.jumpcloud_user_group.test", "description", "A detailed description"),
				),
			},
		},
	})
}

func testAccUserGroupDataSourceConfig() string {
	return `
resource "jumpcloud_user_group" "test" {
  name        = "tf-acc-ds-test-group"
  description = "Test group for data source"
}

data "jumpcloud_user_group" "test" {
  name = jumpcloud_user_group.test.name
}
`
}

func testAccUserGroupDataSourceConfigWithDescription() string {
	return `
resource "jumpcloud_user_group" "test" {
  name        = "tf-acc-ds-test-group-desc"
  description = "A detailed description"
}

data "jumpcloud_user_group" "test" {
  name = jumpcloud_user_group.test.name
}
`
}
