terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    # This backend configuration is filled in automatically at test time by Terratest. If you wish to run this example
    # manually, uncomment and fill in the config below.

    bucket         = "terraform-up-and-running-state-dhmoney"
    key            = "prod/data-stores/mysql/terraform.tfstate"
    region         = "us-east-2"
    dynamodb_table = "terraform-up-and-running-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = "us-east-2"
  alias  = "primary"
}

provider "aws" {
  region = "us-west-1"
  alias  = "replica"
}

##Only for use with AWS KMS
# data "aws_secretsmanager_secret_version" "creds" {
#   secret_id = "db-creds"
# }

# locals {
#   db_creds = jsondecode(
#     data.aws_secretsmanager_secret_version.creds.secret_string
#   )
# }

module "mysql_primary" {
  source = "../../../modules/data-stores/mysql"

  providers = {
    aws = aws.primary
  }

  db_name     = var.db_name

  # # Pass the secrets to the resource
  # db_username = local.db_creds.username
  # db_password = local.db_creds.password

  db_username = var.db_username
  db_password = var.db_password

  # Must be enabled to support replication
  backup_retention_period = 1
}

module "mysql_replica" {
  source = "../../../modules/data-stores/mysql"

  providers = {
    aws = aws.replica
  }

  # Make this a replica of the primary
  replicate_source_db = module.mysql_primary.arn
}
