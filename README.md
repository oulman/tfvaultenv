# tfvaultenv

## Overview

tfvaultenv reads secrets from HashiCorp Vault and outputes environment variables for various Terraform providers with those secrets.

This project is a work in progress and additional Secrets Engines, Providers, and features are planned.

Currently supported are:

### Secrets Engines

- [Active Directory](https://www.vaultproject.io/docs/secrets/ad) (Password Rotation)
- Kv2
- AWS STS

### Terraform Providers

- vSphere
- F5 BIG IP
- Infoblox
- AWS

## Installation

- Download the release for your platform from [Releases](https://github.com/oulman/tfvaultenv/releases)
- untar or unzip the file and move tfvaultenv into your $PATH
- Create a .tfvaultenv.config.hcl file in your Terraform project. (see Configuration below and the examples directory)

## Configuration

The configuration is written in hcl in `.tfvaultenv.config.hcl`. By default tfvaultenv will look in the current working directory for the config file. You can optionally set the `TFVAULTENV_CONFIG_DEPTH` environment variable to search up to N parent directories. This is useful in nested Terraform directory structure scenarios.

### Active Directory

#### Example

```hcl
ad "vsphere" {
   role = "rolename"
   target_provider = "vsphere"
   extra_env_vars = {
       "VSPHERE_SERVER" = "vcenter.example.com"
   }
}
```

#### Arguments

- `role`: (Required) Name of the [Vault Active Directory Secrets Engine role name](https://www.vaultproject.io/docs/secrets/ad)
- `target_provider`: (Required) Name of the Terraform provider to generate environment variables for
- `extra_env_vars`: (Optional) Map of additional environment variables to set

## Usage

### Persisting environment variables

```
$ export `tfvaultenv`
$ env | grep AWS_
AWS_ACCESS_KEY_ID=ASIA<SNIP>
AWS_ACCESS_SECRET_KEY=nJJFD/<SNIP>
AWS_ACCESS_SESSION_TOKEN=<SNIP>
```

### Piping to Terraform

```
$ tfvaultenv | terraform apply
<SNIP>
```

### Printing to stdout

```
$ tfvaultenv
AWS_ACCESS_KEY_ID=ASIA<SNIP>
AWS_ACCESS_SECRET_KEY=nJJFD/<SNIP>
AWS_ACCESS_SESSION_TOKEN=<SNIP>
```
