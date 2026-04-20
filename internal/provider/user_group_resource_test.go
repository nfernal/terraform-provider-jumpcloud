// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupResourceConfig("tf-acc-test-group"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user_group.test", "id"),
					resource.TestCheckResourceAttr("jumpcloud_user_group.test", "name", "tf-acc-test-group"),
				),
			},
			{
				ResourceName:      "jumpcloud_user_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUserGroupResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupResourceConfigWithDescription("tf-acc-test-group-update", "original description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user_group.test", "name", "tf-acc-test-group-update"),
					resource.TestCheckResourceAttr("jumpcloud_user_group.test", "description", "original description"),
				),
			},
			{
				Config: testAccUserGroupResourceConfigWithDescription("tf-acc-test-group-update", "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user_group.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccUserGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "jumpcloud_user_group" "test" {
  name = %[1]q
}
`, name)
}

func TestAccUserGroupResource_createWithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupResourceConfigWithDescription("tf-acc-test-group-with-desc", "initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user_group.test", "id"),
					resource.TestCheckResourceAttr("jumpcloud_user_group.test", "name", "tf-acc-test-group-with-desc"),
					resource.TestCheckResourceAttr("jumpcloud_user_group.test", "description", "initial description"),
				),
			},
			{
				ResourceName:      "jumpcloud_user_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUserGroupResourceConfigWithDescription(name, description string) string {
	return fmt.Sprintf(`
resource "jumpcloud_user_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
