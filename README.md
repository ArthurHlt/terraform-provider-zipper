# Terraform-provider-zipper [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A terraform provider to create to create a zip file from different kind of source, this sources can be:

- [A git repository](https://github.com/ArthurHlt/zipper#git)
- [An http url](https://github.com/ArthurHlt/zipper#http)
- [A local folder](https://github.com/ArthurHlt/zipper#local)

It uses [zipper](https://github.com/ArthurHlt/zipper) under the hood.

You can access documentation at https://registry.terraform.io/providers/ArthurHlt/zipper

## Installations

**Requirements:** You need, of course, terraform (**>=0.12**) which is available here: https://www.terraform.io/downloads.html

### Automatic

Follow instruction on https://registry.terraform.io/providers/ArthurHlt/zipper

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

See documentation at https://registry.terraform.io/providers/ArthurHlt/zipper