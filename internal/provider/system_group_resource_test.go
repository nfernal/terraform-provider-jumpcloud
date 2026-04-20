// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupResourceConfig("tf-acc-test-sysgroup"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_system_group.test", "id"),
					resource.TestCheckResourceAttr("jumpcloud_system_group.test", "name", "tf-acc-test-sysgroup"),
				),
			},
			{
				ResourceName:      "jumpcloud_system_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSystemGroupResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupResourceConfigWithDescription("tf-acc-test-sysgroup-update", "original"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_system_group.test", "description", "original"),
				),
			},
			{
				Config: testAccSystemGroupResourceConfigWithDescription("tf-acc-test-sysgroup-update", "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_system_group.test", "description", "updated"),
				),
			},
		},
	})
}

func testAccSystemGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "jumpcloud_system_group" "test" {
  name = %[1]q
}
`, name)
}

func TestAccSystemGroupResource_createWithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupResourceConfigWithDescription("tf-acc-test-sysgroup-with-desc", "system group description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_system_group.test", "id"),
					resource.TestCheckResourceAttr("jumpcloud_system_group.test", "name", "tf-acc-test-sysgroup-with-desc"),
					resource.TestCheckResourceAttr("jumpcloud_system_group.test", "description", "system group description"),
				),
			},
			{
				ResourceName:      "jumpcloud_system_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSystemGroupResourceConfigWithDescription(name, description string) string {
	return fmt.Sprintf(`
resource "jumpcloud_system_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
