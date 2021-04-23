# Terraform Provider for Idefix

* [Terraform Website](https://www.terraform.io)

## Usage Example

> When using the Idefix Provider with Terraform 0.13 and later, the recommended approach is to declare Provider versions in the root module Terraform configuration, using a `required_providers` block as per the following example. For previous versions, please continue to pin the version within the provider block.

```hcl
# We strongly recommend using the required_providers block to set the
# Idefix Provider source and version being used
terraform {
  required_providers {
    idefix = {
      source = "linkbynet/idefix"
      version = "=0.0.1"
    }
  }
}

# Configure the Idefix Provider
provider "idefix" {
  #login = "..."
  #password = "..."
}
```