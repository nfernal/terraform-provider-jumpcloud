terraform {
  required_providers {
    jumpcloud = {
      source  = "nfernal/jumpcloud"
      version = "~> 0.1"
    }
  }
}

provider "jumpcloud" {
  # api_key = var.jumpcloud_api_key  # Or set JUMPCLOUD_API_KEY env var
  # org_id  = var.jumpcloud_org_id   # Or set JUMPCLOUD_ORG_ID env var
}
