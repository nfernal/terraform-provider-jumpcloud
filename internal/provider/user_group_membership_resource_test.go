// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupMembershipResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupMembershipResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user_group_membership.test", "id"),
					resource.TestCheckResourceAttrPair(
						"jumpcloud_user_group_membership.test", "group_id",
						"jumpcloud_user_group.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"jumpcloud_user_group_membership.test", "user_id",
						"jumpcloud_user.test", "id",
					),
				),
			},
			{
				ResourceName:      "jumpcloud_user_group_membership.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUserGroupMembershipResource_multipleMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupMembershipResourceConfigMultiple(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jumpcloud_user_group_membership.test1", "id"),
					resource.TestCheckResourceAttrSet("jumpcloud_user_group_membership.test2", "id"),
					resource.TestCheckResourceAttrPair(
						"jumpcloud_user_group_membership.test1", "group_id",
						"jumpcloud_user_group_membership.test2", "group_id",
					),
				),
			},
		},
	})
}

func TestAccUserGroupMembershipResource_invalidImportID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupMembershipResourceConfig(),
				Check:  resource.TestCheckResourceAttrSet("jumpcloud_user_group_membership.test", "id"),
			},
			{
				ResourceName:  "jumpcloud_user_group_membership.test",
				ImportState:   true,
				ImportStateId: "invalid-no-slash",
				ExpectError:   regexp.MustCompile(`Invalid Import ID`),
			},
		},
	})
}

func testAccUserGroupMembershipResourceConfig() string {
	return `
resource "jumpcloud_user" "test" {
  username = "tf-acc-membership-test"
  email    = "tf-acc-membership-test@example.com"
}

resource "jumpcloud_user_group" "test" {
  name = "tf-acc-membership-test-group"
}

resource "jumpcloud_user_group_membership" "test" {
  group_id = jumpcloud_user_group.test.id
  user_id  = jumpcloud_user.test.id
}
`
}

func testAccUserGroupMembershipResourceConfigMultiple() string {
	return `
resource "jumpcloud_user" "test1" {
  username = "tf-acc-membership-multi-1"
  email    = "tf-acc-membership-multi-1@example.com"
}

resource "jumpcloud_user" "test2" {
  username = "tf-acc-membership-multi-2"
  email    = "tf-acc-membership-multi-2@example.com"
}

resource "jumpcloud_user_group" "test" {
  name = "tf-acc-membership-multi-group"
}

resource "jumpcloud_user_group_membership" "test1" {
  group_id = jumpcloud_user_group.test.id
  user_id  = jumpcloud_user.test1.id
}

resource "jumpcloud_user_group_membership" "test2" {
  group_id = jumpcloud_user_group.test.id
  user_id  = jumpcloud_user.test2.id
}
`
}
