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

    # Replace this with your bucket name!
    bucket         = "terraform-up-and-running-state-dhmoney"
    key            = "stage/data-stores/mysql/terraform.tfstate"
    region         = "us-east-2"

    # # Replace this with your DynamoDB table name!
    dynamodb_table = "terraform-up-and-running-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = "us-east-2"
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

module "mysql" {
  source = "../../../modules/data-stores/mysql"

  db_name     = var.db_name

  # # Pass the secrets to the resource
  # username = local.db_creds.username
  # password = local.db_creds.password

  db_username = var.db_username
  db_password = var.db_password
}
