---
page_title: "remove_key_pairs Resource - terraform-provider-awsextras"
subcategory: ""
description: |-
    The remove_key_pairs resource allows you to remove key pairs from EC2 that are not controlled by terraform. This is
    useful to help cleanup key pairs especially if you have other systems in your infrastructure that dynamically creates
    key pairs for various uses.
---

# Resource `awsextras_remove_key_pairs`

The remove_key_pairs resource allows you to remove key pairs from EC2 that are not controlled by terraform. This is
useful to help cleanup key pairs especially if you have other systems in your infrastructure that dynamically creates
key pairs for various uses.

## Example Usage

```terraform
resource "awsextras_remove_key_pairs" "unmanaged" {
  include_regex = ["student[0-9]+-.*"]
  exclude_names = ["my-known-keypair"]
}
```

## Argument Reference

- `include_regex` - (Optional) Include key pairs based on regex matching on their name
- `exclude_regex` - (Optional) Exclude key pairs based on regex matching on their name
- `exclude_names` - (Optional) Exclude key pairs based on their key name

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `key_pair_names` - Key Pair Names found that should be removed.
