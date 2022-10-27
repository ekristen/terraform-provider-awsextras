provider "aws" {
  region = "us-east-2"
}

provider "awsextras" {
  region = "us-east-2"
}

resource "random_string" "testing" {
  length  = 5
  special = false
}

resource "aws_key_pair" "testing" {
  key_name   = "deployer-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 email@example.com"
}

resource "awsextras_remove_key_pairs" "testing" {
  include_regex          = ["student[0-9]+-.*"]
  exclude_key_pair_names = [aws_key_pair.testing.id]
}

terraform {
  required_providers {
    awsextras = {
      source  = "ekristen/awsextras"
      version = ">= 0.1.0"
    }
  }
}
