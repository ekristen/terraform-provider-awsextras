provider "aws" {
  region = "us-east-2"
}

provider "awsextras" {
  region = "us-east-2"
}

resource "awsextras_secretsmanager_secret_rotation" "testing" {
  secret_id = "rds!db-04cd90a6-1304-435b-81a7-b70107218515"
  automatically_after_days = 1
}

terraform {
  required_providers {
    awsextras = {
      source  = "ekristen/awsextras"
      version = ">= 0.1.0"
    }
  }
}
