# Terraform Provider - AWS Extras

This is an AWS **Extras** provider. It provides various extras that the default/primary provider does not.

## Resources

- awsextras_terminate_instances

### awsextras_terminate_instances

**Dangerous if not configured properly**

This allows you to identify aws instances in a subnet that are not managed terraform state and remove them on a destroy.

This will current filter on subnet ids and will try and identify instances that should be removed, you and provide a list
of instance ids that should be excluded automatically, otherwise it tries to determine if an instance should be removed
by its lack of tagging.
