---
page_title: "secretsmanager_secret_rotation Resource - terraform-provider-awsextras"
subcategory: ""
description: |-
    The secretsmanager_secret_rotation allows you to change the rotation rules for a secret that is already in secret manager.
    This is useful for when AWS manages the secret for you but does not provide all the knobs to control the rotation schedule.
---

# Resource `awsextras_secretsmanager_secret_rotation`

The secretsmanager_secret_rotation allows you to change the rotation rules for a secret that is already in secret manager.
This is useful for when AWS manages the secret for you but does not provide all the knobs to control the rotation schedule.

## Example Usage

```terraform
resource "awsextras_secretsmanager_secret_rotation" "rds" {
  secret_id = "rds!db-00000000-0000-0000-0000-000000000000"
  automatically_after_days = 1
}
```

## Argument Reference

- `secret_id` - (Required) The ID of the Secret in Secrets Manager.
- `automatically_after_days` - (Required) The number of days after you want the secret to automatically rotate.

## Attributes Reference

In addition to all the arguments above, there are no additionally exported attributes.
