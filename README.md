# Terraform-provider-zipper [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A terraform provider to create to create a zip file from different kind of source, this sources can be:

- [A git repository](https://github.com/ArthurHlt/zipper#git)
- [An http url](https://github.com/ArthurHlt/zipper#http)
- [A local folder](https://github.com/ArthurHlt/zipper#local)

It uses [zipper](https://github.com/ArthurHlt/zipper) under the hood.

## Installations

**Requirements:** You need, of course, terraform (**>=0.12**) which is available here: https://www.terraform.io/downloads.html

### Automatic

To install a specific version, set PROVIDER_CLOUDFOUNDRY_VERSION before executing the following command

```bash
$ export PROVIDER_ZIPPER_VERSION="v0.12.0"
```

#### via curl

```bash
$ bash -c "$(curl -fsSL https://raw.github.com/ArthurHlt/terraform-provider-zipper/master/bin/install.sh)"
```

#### via wget

```bash
$ bash -c "$(wget https://raw.github.com/ArthurHlt/terraform-provider-zipper/master/bin/install.sh -O -)"
```

### Manually

1. Get the build for your system in releases: https://github.com/ArthurHlt/terraform-provider-zipper/releases/latest
2. Create a `providers` directory inside terraform user folder: `mkdir -p ~/.terraform.d/providers`
3. Move the provider previously downloaded in this folder: `mv /path/to/download/directory/terraform-provider-zipper ~/.terraform.d/providers`
4. Ensure provider is executable: `chmod +x ~/.terraform.d/providers/terraform-provider-zipper`
5. add `providers` path to your `.terraformrc`:

```bash
cat <<EOF > ~/.terraformrc
providers {
    zipper = "/full/path/to/.terraform.d/providers/terraform-provider-zipper"
}
EOF
```

6. you can now performs any terraform action on zipper resources

## Usage

Simple example to retrieve a git repository and use it on [AWS Lambda function](https://www.terraform.io/docs/providers/aws/r/lambda_function.html):

```tf
provider "zipper" {
  skip_ssl_validation = false
}

resource "zipper_file" "fixture" {
  source      = "https://github.com/ArthurHlt/go-lambda-ping.git"
  output_path = "path/to/lambda/function.zip"
}

resource "aws_iam_role" "lambda_exec_role" {
  name               = "lambda_exec_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "demo_lambda" {
  function_name    = "demo_lambda"
  handler          = "main"
  runtime          = "go1.x"
  filename         = "${zipper_file.fixture.output_path}"
  source_code_hash = "${zipper_file.fixture.output_sha}"
  role             = "${aws_iam_role.lambda_exec_role.arn}"
}
```

## Resource zipper_file

Zip file from different kind of source as resource by using [zipper](https://github.com/ArthurHlt/zipper). Using as resource permit to not download when not necessary.

### Example Usage

Basic usage

```tf
resource "zipper_file" "fixture" {
  type               = "git"
  source             = "https://github.com/ArthurHlt/go-lambda-ping.git"
  output_path        = "path/to/lambda/function.zip"
  not_when_nonexists = false
}
```

## Argument Reference

The following arguments are supported:

- `source` - (Required) Target source for zipper written in uri style ([see zipper doc for more information](https://github.com/ArthurHlt/zipper)).
- `output_path` - (Required) The output of the archive file.
- `type` - (Optional) Source type to use to create zip, e.g.: http, local or git. (if omitted type will be auto-detected)
- `not_when_nonexists` - (Optional) Set to true to not create zip when not exists at output_path if sources files didn't change. (to earn time if not necessary)

## Attributes Reference

The following attributes are exported:

- `id` - Id is actually equivalent to `output_sha`.
- `output_sha` - SHA1 checksum made by zipper.
- `output_size` - Size of the zip file.

## Datasource zipper_file

Equivalent to resource but less smart to know when to download or not source (less restrictive than a resource).