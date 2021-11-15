# tfvaultenv

## Overview

tfvaultenv reads secrets from HashiCorp Vault and outputs environment variables for various Terraform providers with those secrets.

This project is a work in progress and additional Secrets Engines, Providers, and features are planned. Please see the project [roadmap](ROADMAP.md) for more details.

Currently supported are:

### Secrets Engines

- [Active Directory](https://www.vaultproject.io/docs/secrets/ad) (Password Rotation)
- [Kv2](https://www.vaultproject.io/docs/secrets/kv/kv-v2)
- [AWS](https://www.vaultproject.io/docs/secrets/aws) (STS only currently)

### Terraform Providers

- [vSphere](https://registry.terraform.io/providers/hashicorp/vsphere/latest/docs)
- [F5 BIG IP](https://registry.terraform.io/providers/F5Networks/bigip/latest)
- [Infoblox](https://registry.terraform.io/providers/infobloxopen/infoblox/latest)
- [AWS](https://registry.terraform.io/providers/hashicorp/aws/latest)

## Installation

- Download the release for your platform from [Releases](https://github.com/oulman/tfvaultenv/releases)
- untar or unzip the file and move tfvaultenv into your $PATH
- Create a .tfvaultenv.config.hcl file in your Terraform project. (see Configuration below and the examples directory)

## Configuration

The configuration is written in [HCL](https://github.com/hashicorp/hcl) and the default name is `.tfvaultenv.config.hcl`. Unless overridden, tfvaultenv will look in the current working directory for the config file. You can optionally set the `--config` and `--configdepth` arguments to change the config file name or search up to N parent directories. This is useful in nested Terraform directory structure scenarios.

Configuration is set in blocks representing supported secrets engines and authentication methods.

### Secrets Engines

#### AWS

##### Example

```hcl
aws "sts" {
   method = "assumed_role"
   role = "rolename"
   role_arn = "arn:aws:iam::00000000000:role/TerraformRole"
   extra_env_vars = {
       "AWS_DEFAULT_REGION" = "us-east-2"
   }
   ttl = 900
}
```

##### Arguments

- `method`: (Required) Name of the [AWS Secrets Engine Method](https://www.vaultproject.io/docs/secrets/aws) Currently only `assumed_role` is supported
- `role`: (Required) AWS Secrets Engine role name
- `role_arn`: (Optional) Role ARN to assume when method is set to `assumed_role`
- `extra_env_vars`: (Optional) Map of additional environment variables to set
- `mount`: (Optional) Path to the mounted AWS secrets engine. Default: `aws`
- `ttl`: (Optional) TTL to set on the token or iam_user

#### Azure
##### Example

```hcl
azure "sub1" {
   role = "sub1-rw"
   extra_env_vars = {
       "ARM_TENANT_ID" = "194dd302-295b-4993-b29e-2ca2d37b9031"
       "ARM_SUBSCRIPTION_ID" = "9b9c4322-74a2-474e-ad94-c5e6890713c9"
   }
}
```

##### Arguments

- `role`: (Required) Azure Secrets Engine role name
- `extra_env_vars`: (Optional) Map of additional environment variables to set
- `mount`: (Optional) Path to the mounted Azure secrets engine. Default: `azure`

#### Active Directory

##### Example

```hcl
ad "vsphere" {
   role = "rolename"
   target_provider = "vsphere"
   extra_env_vars = {
       "VSPHERE_SERVER" = "vcenter.example.com"
   }
}
```

##### Arguments

- `role`: (Required) Name of the [Vault Active Directory Secrets Engine role name](https://www.vaultproject.io/docs/secrets/ad)
- `target_provider`: (Required) Name of the Terraform provider to generate environment variables for
- `extra_env_vars`: (Optional) Map of additional environment variables to set
- `path`: (Optional) Path to the mounted AD secrets engine. Default: `ad`

#### Kv2 Secret

##### Example

```hcl
kv_secret "infoblox" {
   path = "infoblox/terraform"
   target_provider = "infoblox"
   attribute_map = {
       "ib_user"     = "username"
       "ib_password" = "password"
   }
   extra_env_vars = {
       "FOO" = "bar"
   }
}
```

##### Arguments

- `path`: (Required) Path to the secret under the secrets engine mount
- `mount`: (Optional) Mount name of the secrets engine. Default: "secrets"
- `attribute_map`: (Optional) Map of kv2 secret attribute names to provider values. Defaults to username and password
- `target_provider`: (Required) Name of the Terraform provider to generate environment variables for
- `extra_env_vars`: (Optional) Map of additional environment variables to set

### Auth Methods

By default `tfvaultenv` creates an implicit auth method that supports token based authentication in the form of VAULT_TOKEN, ~/.vault-token, and token helpers. Supported auth methods such as JWT (see below) can be used and can override token auth by configuring a priority of 1 or above. Auth methods can be conditionally activated using `when {}` blocks based on environment variables or other supported conditions. When multiple auth methods are defined you can specify priorities to ensure that the preferred fallback auth method is used.

#### Common arguments

- `method`: (Required) Name of the Vault authentication method
- `path`: (Required) Path to the auth engine mount
- `priority`: (Required) Priority - set > 0 to override implicit token based auth
- `when`: (Optional) Conditional block methods to determine if the auth method should be used. Currently only `env_present` is supported.

#### JWT

```hcl
auth "gitlab" {
  method = "jwt"
  path = "gitlab"
  priority = 100

  jwt {
    role = env("VAULT_ROLE")
    token = env("CI_JOB_JWT")
  }

  when {
    env_present = "CI_JOB_JWT"
  }
}
```

##### Arguments

- `role`: (Required) Name of the JWT auth engine role
- `token`: (Required) JWT token to pass to Vault API

## Usage

### Setting environment variables

```
$ export `tfvaultenv get`
$ env | grep AWS_
AWS_ACCESS_KEY_ID=ASIA<SNIP>
AWS_ACCESS_SECRET_KEY=nJJFD/<SNIP>
AWS_ACCESS_SESSION_TOKEN=<SNIP>
```

### Piping to Terraform

```
$ tfvaultenv get | terraform apply
<SNIP>
```

### Printing to stdout

```
$ tfvaultenv get
AWS_ACCESS_KEY_ID=ASIA<SNIP>
AWS_ACCESS_SECRET_KEY=nJJFD/<SNIP>
AWS_ACCESS_SESSION_TOKEN=<SNIP>
```

### Specifying an alternate configuration file

```
$ tfvaultenv get --config /path/to/config.hcl
```
