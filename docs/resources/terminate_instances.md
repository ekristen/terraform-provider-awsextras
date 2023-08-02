---
page_title: "terminate_instances Resource - terraform-provider-awsextras"
subcategory: ""
description: |-
    The terminate_instances resource allows you to terminate instances not controlled by terraform, this is useful if you
    are building infrastructure where something else can build instances in your subnets and you want to be able to simply
    be able to run a destroy.
---

# Resource `awsextras_terminate_instances`

The terminate_instances resource allows you to terminate instances not controlled by terraform, this is useful if you
are building infrastructure where something else can build instances in your subnets and you want to be able to simply
be able to run a destroy.

## Example Usage

```terraform
resource "awsextras_terminate_instances" "unmanaged" {
  subnet_ids = [aws_subnet.testing.id]
  exclude_instances_ids = [aws_instance.managed.id]
  include_untagged = true
}
```

## Argument Reference

- `subnet_ids` - (Required) Limits instance search to provided list of Subnet IDs.
- `include_untagged` - (Optional) [default: true] Includes any instance found with no tags or only `Name` tag present.
- `exclude_instance_ids` - (Optional) Used to whitelist or exclude instance ids known to be managed by terraform.
- `exclude_tags` - (Optional) Used to whitelist or exclude instances if they contain any tag of these tags.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `instance_ids` - Instance ids found that should be terminated.
