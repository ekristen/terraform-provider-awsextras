provider "aws" {
  region = "us-east-2"
}

provider "awsextras" {
  region = "us-east-2"
}

resource "random_string" "testing" {
  length = 5
  special = false
}

resource "aws_vpc" "testing" {
  cidr_block = "10.0.0.0/16"
  tags = {
    name = "tf-provider-awsextras-testing-${random_string.testing.id}"
  }
}

resource "aws_subnet" "testing" {
  vpc_id     = aws_vpc.testing.id
  cidr_block = "10.0.1.0/24"
  tags = {
    name = "tf-provider-awsextras-testing-${random_string.testing.id}"
  }
}

resource "awsextras_terminate_instances" "testing" {
  subnet_ids = [aws_subnet.testing.id]
}

terraform {
  required_providers {
    awsextras = {
      source  = "ekristen/awsextras"
      version = ">= 0.1.0"
    }
  }
}