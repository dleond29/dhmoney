# MySQL on RDS example (prod environment)

This folder contains an example [Terraform](https://www.terraform.io/) configuration that deploys a MySQL database  (using
[RDS](https://aws.amazon.com/rds/) in an [Amazon Web Services (AWS) account](http://aws.amazon.com/).

## Quick start

Configure your [AWS access 
keys](http://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys) as 
environment variables:

```
export AWS_ACCESS_KEY_ID=(your access key id)
export AWS_SECRET_ACCESS_KEY=(your secret access key)
```

Configure the database credentials as environment variables:

```
export TF_VAR_db_username=(desired database username)
export TF_VAR_db_password=(desired database password)
```

Open `main.tf`, uncomment the `backend` configuration, and fill in the name of your S3 bucket, DynamoDB table, and
the path to use for the Terraform state file.

Deploy the code:

```
terraform init
terraform apply
```

Clean up when you're done:

```
terraform destroy
```
