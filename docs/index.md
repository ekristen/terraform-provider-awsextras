---
page_title: "Provider: AWS Extras"
subcategory: ""
description: |-
Terraform provider for interacting with AWS in ways the default official provider does not.
---

# AWS Extras Provider

The AWS Extras provider is used to interact with AWS in ways that the official provider does not.

Use the navigation to the left to read about the available resources.

## Example Usage

At the moment only supports region and configuration and authentication via environment variables

```terraform
provider "awsextras" {
  region = "us-east-2"
}
```

## Schema

### Optional

- **region** (String, Optional) AWS Region
